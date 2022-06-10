package admission

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/sirupsen/logrus"
	admissionv1 "k8s.io/api/admission/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/webhook"

	operatorv1alpha1 "github.com/kong/gateway-operator/api/v1alpha1"
)

var (
	scheme = runtime.NewScheme()
	codecs = serializer.NewCodecFactory(scheme)
)

func NewWebhookServerFromManager(mgr ctrl.Manager) *webhook.Server {
	hookServer := mgr.GetWebhookServer()
	handler := NewRequestHandler(logrus.New())
	hookServer.Register("/validate", handler)
	return hookServer
}

type Validator interface {
	ValidateControlPlane(context.Context, operatorv1alpha1.ControlPlane) error
	ValidateDataPlane(context.Context, operatorv1alpha1.DataPlane) error
}

type RequestHandler struct {
	// Validator validates the entities that the k8s API-server asks
	// it the server to validate.
	Validator Validator

	Logger logrus.FieldLogger
}

func NewRequestHandler(logger logrus.FieldLogger) *RequestHandler {
	return &RequestHandler{
		Validator: &validator{},
		Logger:    logger.WithField("a", "b"),
	}
}

func (h *RequestHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Body == nil {
		h.Logger.Error("received request with empty body")
		http.Error(w, "admission review object is missing",
			http.StatusBadRequest)
		return
	}
	data, err := io.ReadAll(r.Body)
	if err != nil {
		h.Logger.WithError(err).Error("failed to read request from client")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	review := admissionv1.AdmissionReview{}
	if err := json.Unmarshal(data, &review); err != nil {
		h.Logger.WithError(err).Error("failed to parse AdmissionReview object")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response, err := h.handleValidation(r.Context(), review.Request)
	if err != nil {
		h.Logger.WithError(err).Error("failed to run validation")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	review.Response = response
	data, err = json.Marshal(review)
	if err != nil {
		h.Logger.WithError(err).Error("failed to marshal response")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, err = w.Write(data)
	if err != nil {
		h.Logger.WithError(err).Error("failed to write response")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

var (
	controlPlaneGVResource = metav1.GroupVersionResource{
		Group:    operatorv1alpha1.GroupVersion.Group,
		Version:  operatorv1alpha1.GroupVersion.Version,
		Resource: "controlplanes",
	}
	dataPlaneGVResource = metav1.GroupVersionResource{
		Group:    operatorv1alpha1.GroupVersion.Group,
		Version:  operatorv1alpha1.GroupVersion.Version,
		Resource: "dataplane",
	}
)

func (h *RequestHandler) handleValidation(ctx context.Context, req *admissionv1.AdmissionRequest) (
	*admissionv1.AdmissionResponse, error) {

	deserializer := codecs.UniversalDeserializer()

	var response admissionv1.AdmissionResponse
	ok := true
	msg := ""

	switch req.Resource {
	case controlPlaneGVResource:
		controlPlane := operatorv1alpha1.ControlPlane{}
		if req.Operation == admissionv1.Create || req.Operation == admissionv1.Update {
			_, _, err := deserializer.Decode(req.Object.Raw, nil, &controlPlane)
			if err != nil {
				return nil, err
			}
			err = h.Validator.ValidateControlPlane(ctx, controlPlane)
			if err != nil {
				ok = false
				msg = err.Error()
			}
		}
	case dataPlaneGVResource:
		dataPlane := operatorv1alpha1.DataPlane{}
		if req.Operation == admissionv1.Create || req.Operation == admissionv1.Update {
			_, _, err := deserializer.Decode(req.Object.Raw, nil, &dataPlane)
			if err != nil {
				return nil, err
			}
			err = h.Validator.ValidateDataPlane(ctx, dataPlane)
			if err != nil {
				ok = false
				msg = err.Error()
			}
		}
	}

	response.UID = req.UID
	response.Allowed = ok
	response.Result = &metav1.Status{
		Message: msg,
	}
	if !ok {
		response.Result.Code = 400
	}
	return &response, nil
}
