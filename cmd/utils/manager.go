package utils

import (
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/metrics/filters"
	metricsserver "sigs.k8s.io/controller-runtime/pkg/metrics/server"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

func (opts *AppOptions) GetManager(scheme *runtime.Scheme) (manager.Manager, error) {
	webhookServerOptions := webhook.Options{}
	if len(opts.webhookCertPath) > 0 {
		opts.Log.Info("Initializing webhook certificate watcher using provided certificates",
			"webhook-cert-path", opts.webhookCertPath, "webhook-cert-name", opts.webhookCertName, "webhook-cert-key", opts.webhookCertKey)

		webhookServerOptions.CertDir = opts.webhookCertPath
		webhookServerOptions.CertName = opts.webhookCertName
		webhookServerOptions.KeyName = opts.webhookCertKey
	}

	webhookServer := webhook.NewServer(webhookServerOptions)

	// Metrics endpoint is enabled in 'config/default/kustomization.yaml'. The Metrics options configure the server.
	// More info:
	// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.23.3/pkg/metrics/server
	// - https://book.kubebuilder.io/reference/metrics.html
	metricsServerOptions := metricsserver.Options{
		BindAddress:   opts.metricsAddr,
		SecureServing: opts.secureMetrics,
		TLSOpts:       opts.tlsOpts,
	}

	if opts.secureMetrics {
		// FilterProvider is used to protect the metrics endpoint with authn/authz.
		// These configurations ensure that only authorized users and service accounts
		// can access the metrics endpoint. The RBAC are configured in 'config/rbac/kustomization.yaml'. More info:
		// https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.23.3/pkg/metrics/filters#WithAuthenticationAndAuthorization
		metricsServerOptions.FilterProvider = filters.WithAuthenticationAndAuthorization
	}

	// If the certificate is not specified, controller-runtime will automatically
	// generate self-signed certificates for the metrics server. While convenient for development and testing,
	// this setup is not recommended for production.
	//
	// TODO(user): If you enable certManager, uncomment the following lines:
	// - [METRICS-WITH-CERTS] at config/default/kustomization.yaml to generate and use certificates
	// managed by cert-manager for the metrics server.
	// - [PROMETHEUS-WITH-CERTS] at config/prometheus/kustomization.yaml for TLS certification.
	if len(opts.metricsCertPath) > 0 {
		opts.Log.Info("Initializing metrics certificate watcher using provided certificates",
			"metrics-cert-path", opts.metricsCertPath, "metrics-cert-name", opts.metricsCertName, "metrics-cert-key", opts.metricsCertKey)

		metricsServerOptions.CertDir = opts.metricsCertPath
		metricsServerOptions.CertName = opts.metricsCertName
		metricsServerOptions.KeyName = opts.metricsCertKey
	}

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:                 scheme,
		Metrics:                metricsServerOptions,
		WebhookServer:          webhookServer,
		HealthProbeBindAddress: opts.probeAddr,
		LeaderElection:         opts.enableLeaderElection,
		LeaderElectionID:       "31eb3360.my.domain",
		// LeaderElectionReleaseOnCancel defines if the leader should step down voluntarily
		// when the Manager ends. This requires the binary to immediately end when the
		// Manager is stopped, otherwise, this setting is unsafe. Setting this significantly
		// speeds up voluntary leader transitions as the new leader don't have to wait
		// LeaseDuration time first.
		//
		// In the default scaffold provided, the program ends immediately after
		// the manager stops, so would be fine to enable this option. However,
		// if you are doing or is intended to do any operation such as perform cleanups
		// after the manager stops then its usage might be unsafe.
		// LeaderElectionReleaseOnCancel: true,
	})
	if err != nil {
		opts.Log.Error(err, "Failed to start manager")
		return mgr, err
	}

	if err := mgr.AddHealthzCheck("healthz", healthz.Ping); err != nil {
		opts.Log.Error(err, "Failed to set up health check")
		return mgr, err
	}
	if err := mgr.AddReadyzCheck("readyz", healthz.Ping); err != nil {
		opts.Log.Error(err, "Failed to set up ready check")
		return mgr, err
	}

	opts.Log.Info("Starting manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		opts.Log.Error(err, "Failed to run manager")
		return mgr, err
	}

	return mgr, nil
}
