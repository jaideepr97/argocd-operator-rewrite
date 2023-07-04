package appcontroller

import (
	"github.com/jaideepr97/argocd-operator-rewrite/pkg/permissions"
	"k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (acr *AppControllerReconciler) reconcileServiceAccount() error {
	acr.Logger.V(0).Info("reconciling serviceaccount")

	saRequest := permissions.ServiceaccountRequest{
		InstanceName: acr.Instance.Name,
		Namespace:    acr.Instance.Namespace,
		Component:    ArgoCDApplicationControllerComponent,
		Client:       acr.Client,
	}

	desiredSA := permissions.RequestServiceaccount(saRequest)

	existingSA, err := permissions.GetServiceAccount(desiredSA.Name, desiredSA.Namespace, *acr.Client)
	if err != nil {
		if !errors.IsNotFound(err) {
			acr.Logger.Error(err, "reconcileServiceAccount: failed to retrieve sserviceaccount", "name", existingSA.Name, "namespace", existingSA.Namespace)
			return err
		}

		if err = controllerutil.SetControllerReference(acr.Instance, desiredSA, acr.Scheme); err != nil {
			acr.Logger.Error(err, "reconcileServiceAccount: failed to set owner reference for serviceaccount", "name", desiredSA.Name, "namespace", desiredSA.Namespace)
		}

		if err = permissions.CreateServiceAccount(desiredSA, *acr.Client); err != nil {
			acr.Logger.Error(err, "reconcileServiceAccount: failed to create serviceaccount", "name", desiredSA.Name, "namespace", desiredSA.Namespace)
			return err
		}
		acr.Logger.V(0).Info("reconcileServiceAccount: serviceaccount created", "name", desiredSA.Name, "namespace", desiredSA.Namespace)
		return err

	}

	// serviceaccount exists, do nothing
	return nil
}

func (acr *AppControllerReconciler) DeleteServiceAccount(name, namespace string) error {
	if err := permissions.DeleteServiceAccount(name, namespace, *acr.Client); err != nil {
		acr.Logger.Error(err, "DeleteServiceAccount: failed to delete serviceaccount", "name", name, "namespace", namespace)
		return err
	}
	acr.Logger.V(0).Info("DeleteServiceAccount: serviceaccount deleted", "name", name, "namespace", namespace)
	return nil
}
