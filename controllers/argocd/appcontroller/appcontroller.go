package appcontroller

import (
	"github.com/go-logr/logr"
	"github.com/jaideepr97/argocd-operator-rewrite/api/v1alpha1"
	"github.com/jaideepr97/argocd-operator-rewrite/pkg/argoutil"
	"github.com/jaideepr97/argocd-operator-rewrite/pkg/cluster"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type AppControllerReconciler struct {
	Client            *client.Client
	Scheme            *runtime.Scheme
	Instance          *v1alpha1.ArgoCD
	ClusterScoped     bool
	Logger            logr.Logger
	ManagedNamespaces map[string]string
	SourceNamespaces  map[string]string
}

func (acr *AppControllerReconciler) Reconcile() error {

	acr.Logger = ctrl.Log.WithName(ArgoCDApplicationControllerComponent).WithValues("instance", acr.Instance.Name, "instance-namespace", acr.Instance.Namespace)

	// if instance namespace is deleted, remove app controller resources
	namespace, err := cluster.GetNamespace(acr.Instance.Namespace, *acr.Client)
	if err != nil {
		acr.Logger.Error(err, "Reconcile: failed to retrieve namespace", "name", acr.Instance.Namespace)
	}

	if namespace.DeletionTimestamp != nil {
		if err := acr.DeleteResources(); err != nil {
			acr.Logger.Error(err, "failed to delete resources")
		}
		return err
	}

	// reconcile app controller resources
	if err := acr.reconcileServiceAccount(); err != nil {
		acr.Logger.Error(err, "error reconciling serviceaccount")
		return err
	}

	if acr.ClusterScoped {
		err = acr.reconcileClusterRole()
		if err != nil {
			acr.Logger.Error(err, "error reconciling clusterRole")
			return err
		}
	}

	if err := acr.reconcileRoles(); err != nil {
		acr.Logger.Error(err, "error reconciling roles")
	}

	if err := acr.reconcileRoleBindings(); err != nil {
		acr.Logger.Error(err, "error reconciling rolebindings")
	}

	return nil
}

func (acr *AppControllerReconciler) DeleteResources() error {
	name := argoutil.GenerateResourceName(acr.Instance.Name, ArgoCDApplicationControllerComponent)
	var deletionError error = nil

	if err := acr.DeleteRoleBinding(name, acr.Instance.Namespace); err != nil {
		acr.Logger.Error(err, "DeleteResources: failed to delete roleBinding")
		deletionError = err
	}

	if err := acr.DeleteRole(name, acr.Instance.Namespace); err != nil {
		acr.Logger.Error(err, "DeleteResources: failed to delete role")
		deletionError = err
	}

	if err := acr.DeleteServiceAccount(name, acr.Instance.Namespace); err != nil {
		acr.Logger.Error(err, "DeleteResources: failed to delete serviceaccount")
		deletionError = err
	}

	return deletionError
}
