package utils

import (
	"crypto/tls"

	"github.com/go-logr/logr"
	"github.com/spf13/cobra"
)

type AppOptions struct {
	metricsAddr          string
	metricsCertPath      string
	metricsCertName      string
	metricsCertKey       string
	webhookCertPath      string
	webhookCertName      string
	webhookCertKey       string
	enableLeaderElection bool
	probeAddr            string
	secureMetrics        bool
	enableHTTP2          bool
	tlsOpts              []func(*tls.Config)
	Log                  logr.Logger
}

func AddCommandFlag(cmd *cobra.Command, opts *AppOptions) {
	cmd.Flags().StringVar(&opts.metricsAddr, "metrics-bind-address", "0", "The address the metrics endpoint binds to. "+
		"Use :8443 for HTTPS or :8080 for HTTP, or leave as 0 to disable the metrics service.")
	cmd.Flags().StringVar(&opts.probeAddr, "health-probe-bind-address", ":8081", "The address the probe endpoint binds to.")
	cmd.Flags().BoolVar(&opts.enableLeaderElection, "leader-elect", false,
		"Enable leader election for controller manager. "+
			"Enabling this will ensure there is only one active controller manager.")
	cmd.Flags().BoolVar(&opts.secureMetrics, "metrics-secure", true,
		"If set, the metrics endpoint is served securely via HTTPS. Use --metrics-secure=false to use HTTP instead.")
	cmd.Flags().StringVar(&opts.webhookCertPath, "webhook-cert-path", "", "The directory that contains the webhook certificate.")
	cmd.Flags().StringVar(&opts.webhookCertName, "webhook-cert-name", "tls.crt", "The name of the webhook certificate file.")
	cmd.Flags().StringVar(&opts.webhookCertKey, "webhook-cert-key", "tls.key", "The name of the webhook key file.")
	cmd.Flags().StringVar(&opts.metricsCertPath, "metrics-cert-path", "",
		"The directory that contains the metrics server certificate.")
	cmd.Flags().StringVar(&opts.metricsCertName, "metrics-cert-name", "tls.crt", "The name of the metrics server certificate file.")
	cmd.Flags().StringVar(&opts.metricsCertKey, "metrics-cert-key", "tls.key", "The name of the metrics server key file.")
	cmd.Flags().BoolVar(&opts.enableHTTP2, "enable-http2", false,
		"If set, HTTP/2 will be enabled for the metrics and webhook servers")

	// if the enable-http2 flag is false (the default), http/2 should be disabled
	// due to its vulnerabilities. More specifically, disabling http/2 will
	// prevent from being vulnerable to the HTTP/2 Stream Cancellation and
	// Rapid Reset CVEs. For more information see:
	// - https://github.com/advisories/GHSA-qppj-fm5r-hxr3
	// - https://github.com/advisories/GHSA-4374-p667-p6c8
	disableHTTP2 := func(c *tls.Config) {
		opts.Log.Info("Disabling HTTP/2")
		c.NextProtos = []string{"http/1.1"}
	}

	if !opts.enableHTTP2 {
		opts.tlsOpts = append(opts.tlsOpts, disableHTTP2)
	}
}
