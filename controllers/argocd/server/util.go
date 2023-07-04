package server

import (
	"github.com/jaideepr97/argocd-operator-rewrite/api/v1alpha1"
	"github.com/jaideepr97/argocd-operator-rewrite/pkg/argoutil"
)

// GetServerOperationProcessors will return the numeric Operation Processors value for the ArgoCD Server.
func GetServerOperationProcessors(cr *v1alpha1.ArgoCD) int32 {
	op := ArgoCDDefaultServerOperationProcessors
	if cr.Spec.Controller.Processors.Operation > op {
		op = cr.Spec.Controller.Processors.Operation
	}
	return op
}

// GetArgoServerStatusProcessors will return the numeric Status Processors value for the ArgoCD Server.
func GetArgoServerStatusProcessors(cr *v1alpha1.ArgoCD) int32 {
	sp := ArgoCDDefaultServerStatusProcessors
	if cr.Spec.Controller.Processors.Status > sp {
		sp = cr.Spec.Controller.Processors.Status
	}
	return sp
}

func getSourceNamespaceRBACName(instanceName, instaceNamespace string) string {
	return argoutil.GenerateUniqueResourceName(instanceName, instaceNamespace, ArgoCDServerComponent)
}
