package appcontroller

import (
	"reflect"

	"github.com/jaideepr97/argocd-operator-rewrite/common"
	"github.com/jaideepr97/argocd-operator-rewrite/pkg/argoutil"
	"github.com/jaideepr97/argocd-operator-rewrite/pkg/permissions"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (acr *AppControllerReconciler) reconcileRoleBindings() error {
	var reconciliationError error = nil
	acr.Logger.Info("reconciling rolebindings")

	roleBindingRequest := permissions.RoleBindingRequest{
		InstanceName: acr.Instance.Name,
		Component:    ArgoCDApplicationControllerComponent,
		Client:       acr.Client,
	}

	saName := argoutil.GenerateResourceName(acr.Instance.Name, ArgoCDApplicationControllerComponent)
	sa, err := permissions.GetServiceAccount(saName, acr.Instance.Namespace, *acr.Client)
	if err != nil {
		acr.Logger.Error(err, "reconcileRoleBindings: failed to retrieve serviceaccount")
		return err
	}

	// reconcile roleBindings for managed namespaces
	if err = acr.reconcileManagedRoleBindings(&roleBindingRequest, sa); err != nil {
		reconciliationError = err
	}

	// reconcile roleBindings for source namespaces
	if err = acr.reconcileSourceRoleBindings(&roleBindingRequest, sa); err != nil {
		reconciliationError = err
	}

	return reconciliationError
}

func (acr *AppControllerReconciler) reconcileManagedRoleBindings(rbr *permissions.RoleBindingRequest, sa *corev1.ServiceAccount) error {
	var reconciliationError error = nil

	for namespace := range acr.ManagedNamespaces {
		rbr.Namespace = namespace
		desiredRB := permissions.RequestRoleBinding(*rbr)

		if namespace != acr.Instance.Namespace {
			// add special label for resource management to roleBinding
			if len(desiredRB.Labels) == 0 {
				desiredRB.Labels = make(map[string]string)
			}
			desiredRB.Labels[common.ArgoCDKeyRBACType] = common.ArgoCDRBACTypeResourceMananagement
		}

		desiredRB.Subjects = []rbacv1.Subject{
			{
				Kind:      rbacv1.ServiceAccountKind,
				Name:      sa.Name,
				Namespace: sa.Namespace,
			},
		}

		customRoleName := getCustomRoleName()
		if customRoleName != "" {
			desiredRB.RoleRef = rbacv1.RoleRef{
				APIGroup: rbacv1.GroupName,
				Kind:     "CusterRole",
				Name:     customRoleName,
			}
		} else {
			desiredRB.RoleRef = rbacv1.RoleRef{
				APIGroup: rbacv1.GroupName,
				Kind:     "Role",
				Name:     argoutil.GenerateResourceName(acr.Instance.Name, ArgoCDApplicationControllerComponent),
			}
		}

		existingRB, err := permissions.GetRoleBinding(desiredRB.Name, desiredRB.Namespace, *acr.Client)
		if err != nil {
			if !errors.IsNotFound(err) {
				acr.Logger.Error(err, "reconcileRoleBinding: failed to retrieve roleBinding", "name", existingRB.Name, "namespace", existingRB.Namespace)
				reconciliationError = err
				continue
			}

			// setting owner reference only to control-plane rolebinding
			if namespace == acr.Instance.Namespace {
				if err = controllerutil.SetControllerReference(acr.Instance, desiredRB, acr.Scheme); err != nil {
					acr.Logger.Error(err, "reconcileRoleBinding: failed to set owner reference for roleBinding", "name", desiredRB.Name, "namespace", desiredRB.Namespace)
					reconciliationError = err
				}
			}

			if err = permissions.CreateRoleBinding(desiredRB, *acr.Client); err != nil {
				acr.Logger.Error(err, "reconcileRoleBinding: failed to create roleBinding", "name", desiredRB.Name, "namespace", desiredRB.Namespace)
				reconciliationError = err
				continue
			}
			acr.Logger.V(0).Info("reconcileRoleBinding: roleBinding created", "name", desiredRB.Name, "namespace", desiredRB.Namespace)
			continue
		}

		rbChanged := false
		if !reflect.DeepEqual(existingRB.RoleRef, desiredRB.RoleRef) {
			existingRB.RoleRef = desiredRB.RoleRef
			rbChanged = true
		} else if !reflect.DeepEqual(existingRB.Subjects, desiredRB.Subjects) {
			existingRB.Subjects = desiredRB.Subjects
			rbChanged = true
		} else if !reflect.DeepEqual(existingRB.Labels, desiredRB.Labels) {
			existingRB.Labels = desiredRB.Labels
			rbChanged = true
		}

		if rbChanged {
			if err = permissions.UpdateRoleBinding(existingRB, *acr.Client); err != nil {
				acr.Logger.Error(err, "reconcileRoleBinding: failed to update roleBinding", "name", existingRB.Name, "namespace", existingRB.Namespace)
				reconciliationError = err
				continue
			}
			acr.Logger.V(0).Info("reconcileRoleBinding: roleBinding updated", "name", existingRB.Name, "namespace", existingRB.Namespace)
			continue
		}
	}
	return reconciliationError
}

func (acr *AppControllerReconciler) reconcileSourceRoleBindings(rbr *permissions.RoleBindingRequest, sa *corev1.ServiceAccount) error {
	var reconciliationError error = nil

	for namespace := range acr.SourceNamespaces {
		rbr.Namespace = namespace
		desiredRB := permissions.RequestRoleBinding(*rbr)

		desiredRB.Name = getSourceNamespaceRBACName(acr.Instance.Name, acr.Instance.Namespace)
		if namespace != acr.Instance.Namespace {
			// add special label for app management to roleBinding
			if len(desiredRB.Labels) == 0 {
				desiredRB.Labels = make(map[string]string)
			}
			desiredRB.Labels[common.ArgoCDKeyRBACType] = common.ArgoCDRBACTypeAppManagement
		}

		desiredRB.Subjects = []rbacv1.Subject{
			{
				Kind:      rbacv1.ServiceAccountKind,
				Name:      sa.Name,
				Namespace: sa.Namespace,
			},
		}

		desiredRB.RoleRef = rbacv1.RoleRef{
			APIGroup: rbacv1.GroupName,
			Kind:     "Role",
			Name:     getSourceNamespaceRBACName(acr.Instance.Name, acr.Instance.Namespace),
		}

		existingRB, err := permissions.GetRoleBinding(desiredRB.Name, desiredRB.Namespace, *acr.Client)
		if err != nil {
			if !errors.IsNotFound(err) {
				acr.Logger.Error(err, "reconcileRoleBinding: failed to retrieve roleBinding", "name", existingRB.Name, "namespace", existingRB.Namespace)
				reconciliationError = err
				continue
			}

			if err = permissions.CreateRoleBinding(desiredRB, *acr.Client); err != nil {
				acr.Logger.Error(err, "reconcileRoleBinding: failed to create roleBinding", "name", desiredRB.Name, "namespace", desiredRB.Namespace)
				reconciliationError = err
				continue
			}
			acr.Logger.V(0).Info("reconcileRoleBinding: roleBinding created", "name", desiredRB.Name, "namespace", desiredRB.Namespace)
			continue
		}

		rbChanged := false
		if !reflect.DeepEqual(existingRB.RoleRef, desiredRB.RoleRef) {
			existingRB.RoleRef = desiredRB.RoleRef
			rbChanged = true
		} else if !reflect.DeepEqual(existingRB.Subjects, desiredRB.Subjects) {
			existingRB.Subjects = desiredRB.Subjects
			rbChanged = true
		} else if !reflect.DeepEqual(existingRB.Labels, desiredRB.Labels) {
			existingRB.Labels = desiredRB.Labels
			rbChanged = true
		}

		if rbChanged {
			if err = permissions.UpdateRoleBinding(existingRB, *acr.Client); err != nil {
				acr.Logger.Error(err, "reconcileRoleBinding: failed to update roleBinding", "name", existingRB.Name, "namespace", existingRB.Namespace)
				reconciliationError = err
				continue
			}
			acr.Logger.V(0).Info("reconcileRoleBinding: roleBinding updated", "name", existingRB.Name, "namespace", existingRB.Namespace)
			continue
		}
	}

	return reconciliationError
}

func (acr *AppControllerReconciler) DeleteRoleBinding(name, namespace string) error {
	if err := permissions.DeleteRoleBinding(name, namespace, *acr.Client); err != nil {
		acr.Logger.Error(err, "DeleteRoleBinding: failed to delete roleBinding", "name", name, "namespace", namespace)
		return err
	}
	acr.Logger.V(0).Info("DeleteRoleBinding: roleBinding deleted", "name", name, "namespace", namespace)
	return nil
}
