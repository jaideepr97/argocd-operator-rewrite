package notifications

import (
	"github.com/go-logr/logr"
	"github.com/jaideepr97/argocd-operator-rewrite/api/v1alpha1"
	"github.com/jaideepr97/argocd-operator-rewrite/pkg/argoutil"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type NotificationsReconciler struct {
	Client   *client.Client
	Scheme   *runtime.Scheme
	Instance *v1alpha1.ArgoCD
	Logger   logr.Logger
}

func (nr *NotificationsReconciler) Reconcile() error {
	nr.Logger = ctrl.Log.WithName(ArgoCDNotificationsControllerComponent).WithValues("instance", nr.Instance.Name, "instance-namespace", nr.Instance.Namespace)

	if nr.Instance.Spec.Notifications == nil {
		// nr.Logger.Info("Notifications controller disabled, deleting resources")
		// if err := nr.DeleteResources(); err != nil {
		// 	nr.Logger.Error(err, "one or more resources could not be deleted")
		// }
		return nil
	}

	if err := nr.reconcileRole(); err != nil {
		nr.Logger.Error(err, "error reconciling roles")
		return err
	}

	// if err := nr.reconcileServiceAccount(); err != nil {
	// 	nr.Logger.Error(err, "error reconciling serviceaccount")
	// 	return err
	// }

	// if err := nr.reconcileRoleBinding(); err != nil {
	// 	nr.Logger.Error(err, "error reconciling rolebinding")
	// 	return err
	// }

	// if err := nr.reconcileSecret(); err != nil {
	// 	nr.Logger.Error(err, "error reconciling secret")
	// 	return err
	// }

	// if err := nr.reconcileConfigmap(); err != nil {
	// 	nr.Logger.Error(err, "error reconciling configmap")
	// 	return err
	// }

	// if err := nr.reconcileDeployment(); err != nil {
	// 	nr.Logger.Error(err, "error reconciling deployment")
	// 	return err
	// }
	return nil
}

func (nr *NotificationsReconciler) DeleteResources() error {
	name := argoutil.GenerateResourceName(nr.Instance.Name, ArgoCDNotificationsControllerComponent)
	var deletionError error

	// if err := nr.DeleteDeployment(name, nr.Instance.Namespace); err != nil {
	// 	nr.Logger.Error(err, "DeleteResources: failed to delete role")
	// 	deletionError = err
	// }

	// if err := nr.DeleteConfigMap(name, nr.Instance.Namespace); err != nil {
	// 	nr.Logger.Error(err, "DeleteResources: failed to delete configmap")
	// 	deletionError = err
	// }

	// if err := nr.DeleteSecret(name, nr.Instance.Namespace); err != nil {
	// 	nr.Logger.Error(err, "DeleteResources: failed to delete secret")
	// 	deletionError = err
	// }

	// if err := nr.DeleteRoleBinding(name, nr.Instance.Namespace); err != nil {
	// 	nr.Logger.Error(err, "DeleteResources: failed to delete roleBinding")
	// 	deletionError = err
	// }

	// if err := nr.DeleteServiceaccount(name, nr.Instance.Namespace); err != nil {
	// 	nr.Logger.Error(err, "DeleteResources: failed to delete serviceaccount")
	// 	deletionError = err
	// }

	if err := nr.DeleteRole(name, nr.Instance.Namespace); err != nil {
		nr.Logger.Error(err, "DeleteResources: failed to delete role")
		deletionError = err
	}

	if deletionError != nil {
		return deletionError
	}

	return nil
}
