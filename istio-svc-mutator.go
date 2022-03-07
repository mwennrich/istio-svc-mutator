package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/sirupsen/logrus"
	kwhhttp "github.com/slok/kubewebhook/v2/pkg/http"
	kwhlogrus "github.com/slok/kubewebhook/v2/pkg/log/logrus"
	kwhmodel "github.com/slok/kubewebhook/v2/pkg/model"
	kwhmutating "github.com/slok/kubewebhook/v2/pkg/webhook/mutating"
)

type config struct {
	certFile string
	keyFile  string
}

func initFlags() *config {
	cfg := &config{}

	fl := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	fl.StringVar(&cfg.certFile, "tls-cert-file", "/etc/istio-svc-mutator/cert.pem", "TLS certificate file")
	fl.StringVar(&cfg.keyFile, "tls-key-file", "/etc/istio-svc-mutator/key.pem", "TLS key file")

	fl.Parse(os.Args[1:])
	return cfg
}

func run() error {
	logrusLogEntry := logrus.NewEntry(logrus.New())
	logrusLogEntry.Logger.SetLevel(logrus.DebugLevel)
	logger := kwhlogrus.NewLogrus(logrusLogEntry)

	cfg := initFlags()

	// Create mutator.
	mt := kwhmutating.MutatorFunc(func(_ context.Context, _ *kwhmodel.AdmissionReview, obj metav1.Object) (*kwhmutating.MutatorResult, error) {
		svc, ok := obj.(*corev1.Service)
		if !ok {
			return &kwhmutating.MutatorResult{}, nil
		}

		fmt.Print("Mutating service: ", svc.Name)
		// Mutate our object with the required annotations.
		if svc.Name == "istio-ingressgateway" && svc.Namespace == "istio-ingress" {
			svc.Spec.ExternalTrafficPolicy = corev1.ServiceExternalTrafficPolicyTypeLocal
			svc.Spec.HealthCheckNodePort = 31397
		}

		return &kwhmutating.MutatorResult{MutatedObject: svc}, nil
	})

	// Create webhook.
	mcfg := kwhmutating.WebhookConfig{
		ID:      "istio-svc-mutator.metal-stack.dev",
		Mutator: mt,
		Logger:  logger,
	}
	wh, err := kwhmutating.NewWebhook(mcfg)
	if err != nil {
		return fmt.Errorf("error creating webhook: %w", err)
	}

	// Get HTTP handler from webhook.
	whHandler, err := kwhhttp.HandlerFor(kwhhttp.HandlerConfig{Webhook: wh, Logger: logger})
	if err != nil {
		return fmt.Errorf("error creating webhook handler: %w", err)
	}

	// Serve.
	logger.Infof("Listening on :8080")
	err = http.ListenAndServeTLS(":8080", cfg.certFile, cfg.keyFile, whHandler)
	if err != nil {
		return fmt.Errorf("error serving webhook: %w", err)
	}

	return nil
}

func main() {
	err := run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error running app: %s", err)
		os.Exit(1)
	}
}
