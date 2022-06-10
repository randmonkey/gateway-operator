/*
Copyright 2022 Kong Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package manager

import (
	"fmt"
	"os"

	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/cluster"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	gatewayv1alpha2 "sigs.k8s.io/gateway-api/apis/v1alpha2"

	operatorv1alpha1 "github.com/kong/gateway-operator/apis/v1alpha1"
	"github.com/kong/gateway-operator/controllers"
	"github.com/kong/gateway-operator/internal/admission"
	"github.com/kong/gateway-operator/pkg/vars"
)

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(operatorv1alpha1.AddToScheme(scheme))
	utilruntime.Must(gatewayv1alpha2.AddToScheme(scheme))
	//+kubebuilder:scaffold:scheme
}

type Config struct {
	MetricsAddr     string
	ProbeAddr       string
	WebhookPort     int
	LeaderElection  bool
	DevelopmentMode bool
	Out             *os.File
	NewClientFunc   cluster.NewClientFunc
	ControllerName  string
}

var DefaultConfig = Config{
	MetricsAddr:     ":8080",
	ProbeAddr:       ":8081",
	WebhookPort:     9443,
	DevelopmentMode: false,
	LeaderElection:  true,
}

func Run(cfg Config) error {
	if cfg.ControllerName != "" {
		setupLog.Info(fmt.Sprintf("custom controller name provided: %s", cfg.ControllerName))
		vars.ControllerName = cfg.ControllerName
	}

	opts := zap.Options{
		Development: cfg.DevelopmentMode,
	}

	if cfg.Out != nil {
		opts.DestWriter = cfg.Out
	}

	ctrl.SetLogger(zap.New(zap.UseFlagOptions(&opts)))

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:                 scheme,
		MetricsBindAddress:     cfg.MetricsAddr,
		Port:                   cfg.WebhookPort,
		HealthProbeBindAddress: cfg.ProbeAddr,
		LeaderElection:         cfg.LeaderElection,
		LeaderElectionID:       "a7feedc84.konghq.com",
		NewClient:              cfg.NewClientFunc,
	})
	if err != nil {
		return fmt.Errorf("unable to start manager: %w", err)
	}

	if err = (&controllers.DataPlaneReconciler{
		Client: mgr.GetClient(),
		Scheme: mgr.GetScheme(),
	}).SetupWithManager(mgr); err != nil {
		return fmt.Errorf("unable to create controller DataPlane: %w", err)
	}
	if err = (&controllers.ControlPlaneReconciler{
		Client: mgr.GetClient(),
		Scheme: mgr.GetScheme(),
	}).SetupWithManager(mgr); err != nil {
		return fmt.Errorf("unable to create controller ControlPlane: %w", err)
	}
	if err = (&controllers.GatewayReconciler{
		Client: mgr.GetClient(),
		Scheme: mgr.GetScheme(),
	}).SetupWithManager(mgr); err != nil {
		return fmt.Errorf("unable to create controller Gateway: %w", err)
	}
	//+kubebuilder:scaffold:builder

	if err := mgr.AddHealthzCheck("healthz", healthz.Ping); err != nil {
		return fmt.Errorf("unable to set up health check: %w", err)
	}
	if err := mgr.AddReadyzCheck("readyz", healthz.Ping); err != nil {
		return fmt.Errorf("unable to set up ready check: %w", err)
	}

	tlsCertPath := "/tmp/k8s-webhook-server/serving-certs/tls.crt"
	tlsKeyPath := "/tmp/k8s-webhook-server/serving-certs/tls.key"
	if _, certFileErr := os.Stat(tlsCertPath); certFileErr == nil {
		if _, keyFileErr := os.Stat(tlsKeyPath); keyFileErr == nil {
			hookServer := admission.NewWebhookServerFromManager(mgr)
			setupLog.Info("start webhook at port" + fmt.Sprintf("%s:%d", hookServer.Host, hookServer.Port))
		} else {
			setupLog.Info("TLS key file does not exist, do not start webhook, path:" + tlsKeyPath)
		}
	} else {
		setupLog.Info("TLS certificate file does not exist, do not start webhook, path:" + tlsCertPath)
	}

	setupLog.Info("starting manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		return fmt.Errorf("problem running manager: %w", err)
	}

	return nil
}
