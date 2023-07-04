package configmap

import (
	"github.com/jaideepr97/argocd-operator-rewrite/api/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type ConfigMapReconciler struct {
	Client   *client.Client
	Scheme   *runtime.Scheme
	Instance *v1alpha1.ArgoCD
}

func (cmr *ConfigMapReconciler) Reconcile() error {
	return nil
}
