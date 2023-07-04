package reposerver

import (
	"github.com/jaideepr97/argocd-operator-rewrite/api/v1alpha1"
	"github.com/jaideepr97/argocd-operator-rewrite/pkg/argoutil"
)

// GetRepoServerAddress will return the Argo CD repo server address.
func GetRepoServerAddress(cr *v1alpha1.ArgoCD) string {
	return argoutil.FqdnServiceRef(ArgoCDRepoServerSuffix, cr.Name, cr.Namespace, ArgoCDDefaultRepoServerPort)
}
