package server

import (
	"github.com/go-logr/logr"
	"github.com/jaideepr97/argocd-operator-rewrite/api/v1alpha1"
	"github.com/jaideepr97/argocd-operator-rewrite/pkg/argoutil"
	"github.com/jaideepr97/argocd-operator-rewrite/pkg/cluster"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type ServerReconciler struct {
	Client            *client.Client
	Scheme            *runtime.Scheme
	Instance          *v1alpha1.ArgoCD
	ClusterScoped     bool
	Logger            logr.Logger
	ManagedNamespaces map[string]string
	SourceNamespaces  map[string]string
}

func (sr *ServerReconciler) Reconcile() error {

	sr.Logger = ctrl.Log.WithName(ArgoCDServerComponent).WithValues("instance", sr.Instance.Name, "instance-namespace", sr.Instance.Namespace)

	// if instance namespace is deleted, remove server resources
	namespace, err := cluster.GetNamespace(sr.Instance.Namespace, *sr.Client)
	if err != nil {
		sr.Logger.Error(err, "Reconcile: failed to retrieve namespace", "name", sr.Instance.Namespace)
	}

	if namespace.DeletionTimestamp != nil {
		if err := sr.DeleteResources(); err != nil {
			sr.Logger.Error(err, "failed to delete resources")
		}
		return err
	}

	// reconcile server resources
	if err := sr.reconcileServiceAccount(); err != nil {
		sr.Logger.Error(err, "error reconciling serviceaccount")
		return err
	}

	if sr.ClusterScoped {
		err = sr.reconcileClusterRole()
		if err != nil {
			sr.Logger.Error(err, "error reconciling clusterRole")
			return err
		}
	}

	if err := sr.reconcileRoles(); err != nil {
		sr.Logger.Error(err, "error reconciling roles")
	}

	if err := sr.reconcileRoleBindings(); err != nil {
		sr.Logger.Error(err, "error reconciling rolebindings")
	}

	return nil
}

func (sr *ServerReconciler) DeleteResources() error {
	name := argoutil.GenerateResourceName(sr.Instance.Name, ArgoCDServerComponent)
	var deletionError error = nil

	if err := sr.DeleteRoleBinding(name, sr.Instance.Namespace); err != nil {
		sr.Logger.Error(err, "DeleteResources: failed to delete roleBinding")
		deletionError = err
	}

	if err := sr.DeleteRole(name, sr.Instance.Namespace); err != nil {
		sr.Logger.Error(err, "DeleteResources: failed to delete role")
		deletionError = err
	}

	if err := sr.DeleteServiceAccount(name, sr.Instance.Namespace); err != nil {
		sr.Logger.Error(err, "DeleteResources: failed to delete serviceaccount")
		deletionError = err
	}

	return deletionError
}
