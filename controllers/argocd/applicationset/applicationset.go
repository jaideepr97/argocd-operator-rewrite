package applicationset

import (
	"github.com/go-logr/logr"
	"github.com/jaideepr97/argocd-operator-rewrite/api/v1alpha1"
	"github.com/jaideepr97/argocd-operator-rewrite/pkg/argoutil"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type ApplicationSetReconciler struct {
	Client   *client.Client
	Scheme   *runtime.Scheme
	Instance *v1alpha1.ArgoCD
	Logger   logr.Logger
}

func (asr *ApplicationSetReconciler) Reconcile() error {
	asr.Logger = ctrl.Log.WithName(ArgoCDApplicationSetControllerComponent).WithValues("instance", asr.Instance.Name, "instance-namespace", asr.Instance.Namespace)

	if asr.Instance.Spec.ApplicationSet == nil {
		// asr.Logger.Info("ApplicationSet disabled, deleting resources")
		// if err := asr.DeleteResources(); err != nil {
		// 	asr.Logger.Error(err, "one or more resources could not be deleted")
		// }
		return nil
	}

	if err := asr.reconcileRole(); err != nil {
		asr.Logger.Error(err, "error reconciling roles")
		return err
	}

	// if err := asr.reconcileServiceAccount(); err != nil {
	// 	asr.Logger.Error(err, "error reconciling serviceaccount")
	// 	return err
	// }

	// if err := asr.reconcileRoleBinding(); err != nil {
	// 	asr.Logger.Error(err, "error reconciling rolebinding")
	// 	return err
	// }

	// if err := asr.reconcileDeployment(); err != nil {
	// 	asr.Logger.Error(err, "error reconciling deployment")
	// 	return err
	// }

	// if err := asr.reconcileService(); err != nil {
	// 	asr.Logger.Error(err, "error reconciling service")
	// 	return err
	// }

	return nil
}

func (asr *ApplicationSetReconciler) DeleteResources() error {
	name := argoutil.GenerateResourceName(asr.Instance.Name, ArgoCDApplicationSetControllerComponent)
	var deletionError error

	// if err := asr.DeleteService(name, asr.Instance.Namespace); err != nil {
	// 	asr.Logger.Error(err, "DeleteResources: failed to delete service")
	// 	deletionError = err
	// }

	// if err := asr.DeleteDeployment(name, asr.Instance.Namespace); err != nil {
	// 	asr.Logger.Error(err, "DeleteResources: failed to delete deployment")
	// 	deletionError = err
	// }

	// if err := asr.DeleteRoleBinding(name, asr.Instance.Namespace); err != nil {
	// 	asr.Logger.Error(err, "DeleteResources: failed to delete roleBinding")
	// 	deletionError = err
	// }

	// if err := asr.DeleteServiceaccount(name, asr.Instance.Namespace); err != nil {
	// 	asr.Logger.Error(err, "DeleteResources: failed to delete serviceaccount")
	// 	deletionError = err
	// }

	if err := asr.DeleteRole(name, asr.Instance.Namespace); err != nil {
		asr.Logger.Error(err, "DeleteResources: failed to delete role")
		deletionError = err
	}

	if deletionError != nil {
		return deletionError
	}

	return nil
}
