package dataplane

import (
	"context"
	"encoding/base64"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	k8stypes "k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	operatorv1alpha1 "github.com/kong/gateway-operator/api/v1alpha1"
	"github.com/kong/gateway-operator/internal/consts"
)

type Validator struct {
	c client.Client
}

func NewValidator(c client.Client) *Validator {
	return &Validator{c: c}
}

func (v *Validator) Validate(dataplane *operatorv1alpha1.DataPlane) error {
	err := v.ValidateDeployOptions(dataplane.Namespace, &dataplane.Spec.DeploymentOptions)
	if err != nil {
		return err
	}
	return nil
}

func (v *Validator) ValidateDeployOptions(namespace string, opts *operatorv1alpha1.DeploymentOptions) error {

	// validata db mode.
	dbMode := ""
	dbModeFound := false
	for _, envVar := range opts.Env {
		// use the last appearance of the same key as the result since k8s takes this precedence.
		if envVar.Name == consts.EnvVarKongDatabase {
			// value is non-empty.
			if envVar.Value != "" {
				dbMode = envVar.Value
				dbModeFound = true
			} else if envVar.ValueFrom != nil {
				// value is empty,get from ValueFrom from configmap/secret.
				if envVar.ValueFrom.ConfigMapKeyRef != nil {
					cmKeyRef := envVar.ValueFrom.ConfigMapKeyRef
					cm := &corev1.ConfigMap{}
					namespacedName := k8stypes.NamespacedName{Namespace: namespace, Name: cmKeyRef.Name}
					err := v.c.Get(context.Background(), namespacedName, cm)
					if err != nil {
						return fmt.Errorf("failed to get configMap %s in configMapKeyRef, error %w", cmKeyRef.Name, err)
					}
					if cm.Data != nil && cm.Data[cmKeyRef.Key] != "" {
						dbMode = cm.Data[cmKeyRef.Key]
						dbModeFound = true
					}
				}

				if envVar.ValueFrom.SecretKeyRef != nil {
					secretRef := envVar.ValueFrom.SecretKeyRef
					secret := &corev1.Secret{}
					namespacedName := k8stypes.NamespacedName{Namespace: namespace, Name: secretRef.Name}
					err := v.c.Get(context.Background(), namespacedName, secret)
					if err != nil {
						return fmt.Errorf("failed to get secret %s in secretRef, error %w", secretRef.Name, err)
					}
					if secret.Data != nil && len(secret.Data[secretRef.Key]) > 0 {
						decoded, err := base64.StdEncoding.DecodeString(string(secret.Data[secretRef.Key]))
						if err == nil {
							dbMode = string(decoded)
							dbModeFound = true
						}
					}
				}
			}
		}
	}

	// TODO: if dbMode not found in envVar, search for it in EnvVarFrom.
	_ = dbModeFound

	// only support dbless mode.
	if dbMode != "" && dbMode != "off" {
		return fmt.Errorf("database backend %s of dataplane not supported currently", dbMode)
	}
	return nil
}
