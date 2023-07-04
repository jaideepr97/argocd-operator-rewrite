package reposerver

import (
	"github.com/jaideepr97/argocd-operator-rewrite/api/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type RepoServerReconciler struct {
	Client   *client.Client
	Scheme   *runtime.Scheme
	Instance *v1alpha1.ArgoCD
}

func (rsr *RepoServerReconciler) Reconcile() error {
	return nil
}
