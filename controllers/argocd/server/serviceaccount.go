package server

import (
	"github.com/jaideepr97/argocd-operator-rewrite/pkg/permissions"
	"k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (sr *ServerReconciler) reconcileServiceAccount() error {
	sr.Logger.V(0).Info("reconciling serviceaccount")

	saRequest := permissions.ServiceaccountRequest{
		InstanceName: sr.Instance.Name,
		Namespace:    sr.Instance.Namespace,
		Component:    ArgoCDServerComponent,
		Client:       sr.Client,
	}

	desiredSA := permissions.RequestServiceaccount(saRequest)

	existingSA, err := permissions.GetServiceAccount(desiredSA.Name, desiredSA.Namespace, *sr.Client)
	if err != nil {
		if !errors.IsNotFound(err) {
			sr.Logger.Error(err, "reconcileServiceAccount: failed to retrieve sserviceaccount", "name", existingSA.Name, "namespace", existingSA.Namespace)
			return err
		}

		if err = controllerutil.SetControllerReference(sr.Instance, desiredSA, sr.Scheme); err != nil {
			sr.Logger.Error(err, "reconcileServiceAccount: failed to set owner reference for serviceaccount", "name", desiredSA.Name, "namespace", desiredSA.Namespace)
		}

		if err = permissions.CreateServiceAccount(desiredSA, *sr.Client); err != nil {
			sr.Logger.Error(err, "reconcileServiceAccount: failed to create serviceaccount", "name", desiredSA.Name, "namespace", desiredSA.Namespace)
			return err
		}
		sr.Logger.V(0).Info("reconcileServiceAccount: serviceaccount created", "name", desiredSA.Name, "namespace", desiredSA.Namespace)
		return err

	}

	// serviceaccount exists, do nothing
	return nil
}

func (sr *ServerReconciler) DeleteServiceAccount(name, namespace string) error {
	if err := permissions.DeleteServiceAccount(name, namespace, *sr.Client); err != nil {
		sr.Logger.Error(err, "DeleteServiceAccount: failed to delete serviceaccount", "name", name, "namespace", namespace)
		return err
	}
	sr.Logger.V(0).Info("DeleteServiceAccount: serviceaccount deleted", "name", name, "namespace", namespace)
	return nil
}
