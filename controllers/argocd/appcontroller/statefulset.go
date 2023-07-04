package appcontroller

import (
	"github.com/jaideepr97/argocd-operator-rewrite/pkg/argoutil"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// newStatefulSet returns a new StatefulSet instance for the given ArgoCD instance.
func newStatefulSet(instanceName, namespace, component string) *appsv1.StatefulSet {
	return &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      argoutil.GenerateResourceName(instanceName, component),
			Namespace: namespace,
			Labels:    argoutil.LabelsForCluster(instanceName, component),
		},
	}
}
