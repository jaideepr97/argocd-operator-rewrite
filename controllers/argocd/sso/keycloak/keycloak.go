package sso

import (
	"fmt"

	"github.com/jaideepr97/argocd-operator-rewrite/pkg/workloads"
	oappsv1 "github.com/openshift/api/apps/v1"
)

// HandleKeycloakPodDeletion resets the Realm Creation Status to false when keycloak pod is deleted.
func HandleKeycloakPodDeletion(dc *oappsv1.DeploymentConfig) error {
	dcClient, err := GetOappsClient()
	if err != nil {
		return fmt.Errorf("handleKeycloakPodDeletion: unable to get openshift apps client: %w", err)
	}

	existingDC, err := workloads.GetDeploymentConfig(KeycloakIdentifier, dc.Namespace, dcClient)
	if err != nil {
		return fmt.Errorf("handleKeycloakPodDeletion: unable to get deploymentConfig: %w", err)
	}

	existingDC.Annotations[ArgoCDKeycloakRealmCreatedAnnotation] = "false"

	if err = workloads.UpdateDeploymentConfig(existingDC, dcClient); err != nil {
		return fmt.Errorf("handleKeycloakPodDeletion: unable to update deploymentConfig: %w", err)
	}
	return nil
}
