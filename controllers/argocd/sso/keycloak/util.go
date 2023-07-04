package sso

import (
	"fmt"

	oappsv1client "github.com/openshift/client-go/apps/clientset/versioned/typed/apps/v1"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
)

func GetOappsClient() (*oappsv1client.AppsV1Client, error) {
	cfg, err := config.GetConfig()
	if err != nil {
		return nil, fmt.Errorf("GetOappsClient: unable to get config: %w", err)
	}

	dcClient, err := oappsv1client.NewForConfig(cfg)
	if err != nil {
		return nil, fmt.Errorf("GetOappsClient: unable to create openshift apps client: %w", err)
	}

	return dcClient, nil
}
