// This file is generated by /hack/generators/kic-clusterrole-generator. DO NOT EDIT.

package controllers

// -----------------------------------------------------------------------------
// Kong Ingress Controller - RBAC
// -----------------------------------------------------------------------------

//+kubebuilder:rbac:groups=core,resources=endpoints,verbs=list;watch
//+kubebuilder:rbac:groups=core,resources=endpoints/status,verbs=get;patch;update
//+kubebuilder:rbac:groups=core,resources=events,verbs=create;patch
//+kubebuilder:rbac:groups=core,resources=nodes,verbs=list;watch
//+kubebuilder:rbac:groups=core,resources=pods,verbs=get;list;watch
//+kubebuilder:rbac:groups=core,resources=secrets,verbs=list;watch
//+kubebuilder:rbac:groups=core,resources=secrets/status,verbs=get;patch;update
//+kubebuilder:rbac:groups=core,resources=services,verbs=get;list;watch
//+kubebuilder:rbac:groups=core,resources=services/status,verbs=get;patch;update
//+kubebuilder:rbac:groups=configuration.konghq.com,resources=ingressclassparameterses,verbs=get;list;watch
//+kubebuilder:rbac:groups=configuration.konghq.com,resources=kongclusterplugins,verbs=get;list;watch
//+kubebuilder:rbac:groups=configuration.konghq.com,resources=kongclusterplugins/status,verbs=get;patch;update
//+kubebuilder:rbac:groups=configuration.konghq.com,resources=kongconsumers,verbs=get;list;watch
//+kubebuilder:rbac:groups=configuration.konghq.com,resources=kongconsumers/status,verbs=get;patch;update
//+kubebuilder:rbac:groups=configuration.konghq.com,resources=kongingresses,verbs=get;list;watch
//+kubebuilder:rbac:groups=configuration.konghq.com,resources=kongingresses/status,verbs=get;patch;update
//+kubebuilder:rbac:groups=configuration.konghq.com,resources=kongplugins,verbs=get;list;watch
//+kubebuilder:rbac:groups=configuration.konghq.com,resources=kongplugins/status,verbs=get;patch;update
//+kubebuilder:rbac:groups=configuration.konghq.com,resources=tcpingresses,verbs=get;list;watch
//+kubebuilder:rbac:groups=configuration.konghq.com,resources=tcpingresses/status,verbs=get;patch;update
//+kubebuilder:rbac:groups=configuration.konghq.com,resources=udpingresses,verbs=get;list;watch
//+kubebuilder:rbac:groups=configuration.konghq.com,resources=udpingresses/status,verbs=get;patch;update
//+kubebuilder:rbac:groups=extensions,resources=ingresses,verbs=get;list;watch
//+kubebuilder:rbac:groups=extensions,resources=ingresses/status,verbs=get;patch;update
//+kubebuilder:rbac:groups=networking.k8s.io,resources=ingressclasses,verbs=get;list;watch
//+kubebuilder:rbac:groups=networking.k8s.io,resources=ingresses,verbs=get;list;watch
//+kubebuilder:rbac:groups=networking.k8s.io,resources=ingresses/status,verbs=get;patch;update

//+kubebuilder:rbac:groups=gateway.networking.k8s.io,resources=gatewayclasses,verbs=get;list;watch
//+kubebuilder:rbac:groups=gateway.networking.k8s.io,resources=gatewayclasses/status,verbs=get;update
//+kubebuilder:rbac:groups=gateway.networking.k8s.io,resources=gateways,verbs=get;list;update;watch
//+kubebuilder:rbac:groups=gateway.networking.k8s.io,resources=gateways/status,verbs=get;update
//+kubebuilder:rbac:groups=gateway.networking.k8s.io,resources=httproutes,verbs=get;list;watch
//+kubebuilder:rbac:groups=gateway.networking.k8s.io,resources=httproutes/status,verbs=get;update
//+kubebuilder:rbac:groups=gateway.networking.k8s.io,resources=referencegrants,verbs=get;list;watch
//+kubebuilder:rbac:groups=gateway.networking.k8s.io,resources=referencegrants/status,verbs=get
//+kubebuilder:rbac:groups=gateway.networking.k8s.io,resources=tcproutes,verbs=get;list;watch
//+kubebuilder:rbac:groups=gateway.networking.k8s.io,resources=tcproutes/status,verbs=get;update
//+kubebuilder:rbac:groups=gateway.networking.k8s.io,resources=tlsroutes,verbs=get;list;watch
//+kubebuilder:rbac:groups=gateway.networking.k8s.io,resources=tlsroutes/status,verbs=get;update
//+kubebuilder:rbac:groups=gateway.networking.k8s.io,resources=udproutes,verbs=get;list;watch
//+kubebuilder:rbac:groups=gateway.networking.k8s.io,resources=udproutes/status,verbs=get;update

//+kubebuilder:rbac:groups=networking.internal.knative.dev,resources=ingresses,verbs=get;list;watch
//+kubebuilder:rbac:groups=networking.internal.knative.dev,resources=ingresses/status,verbs=get;patch;update

//+kubebuilder:rbac:groups=core,resources=configmaps,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=coordination.k8s.io,resources=leases,verbs=get;list;watch;create;update;patch;delete
