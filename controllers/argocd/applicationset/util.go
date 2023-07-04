package applicationset

import "github.com/jaideepr97/argocd-operator-rewrite/api/v1alpha1"

// getApplicationSetHTTPServerHost will return the host for the given ArgoCD.
func getApplicationSetHTTPServerHost(cr *v1alpha1.ArgoCD) string {
	host := cr.Name
	if len(cr.Spec.ApplicationSet.WebhookServer.Host) > 0 {
		host = cr.Spec.ApplicationSet.WebhookServer.Host
	}
	return host
}
