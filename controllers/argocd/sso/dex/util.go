package sso

import (
	"fmt"

	"github.com/jaideepr97/argocd-operator-rewrite/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
)

// getDexOAuthClientID will return the OAuth client ID for the given ArgoCD.
func getDexOAuthClientID(cr *v1alpha1.ArgoCD) string {
	return fmt.Sprintf("system:serviceaccount:%s:%s", cr.Namespace, fmt.Sprintf("%s-%s", cr.Name, ArgoCDDefaultDexServiceAccountName))
}

// getDexResources will return the ResourceRequirements for the Dex container.
func getDexResources(cr *v1alpha1.ArgoCD) corev1.ResourceRequirements {
	resources := corev1.ResourceRequirements{}

	// Allow override of resource requirements from CR
	if cr.Spec.SSO.Dex.Resources != nil {
		resources = *cr.Spec.SSO.Dex.Resources
	}
	return resources
}

func getDexConfig(cr *v1alpha1.ArgoCD) string {
	config := ArgoCDDefaultDexConfig

	// Allow override of config from CR
	if cr.Spec.ExtraConfig["dex.config"] != "" {
		config = cr.Spec.ExtraConfig["dex.config"]
	} else if len(cr.Spec.SSO.Dex.Config) > 0 {
		config = cr.Spec.SSO.Dex.Config
	}
	return config
}
