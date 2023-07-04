package server

import (
	"reflect"

	"github.com/jaideepr97/argocd-operator-rewrite/common"
	"github.com/jaideepr97/argocd-operator-rewrite/pkg/argoutil"
	"github.com/jaideepr97/argocd-operator-rewrite/pkg/cluster"
	"github.com/jaideepr97/argocd-operator-rewrite/pkg/permissions"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (sr *ServerReconciler) reconcileRoleBindings() error {
	var reconciliationError error = nil
	sr.Logger.Info("reconciling rolebindings")

	roleBindingRequest := permissions.RoleBindingRequest{
		InstanceName: sr.Instance.Name,
		Component:    ArgoCDServerComponent,
		Client:       sr.Client,
	}

	saName := argoutil.GenerateResourceName(sr.Instance.Name, ArgoCDServerComponent)
	sa, err := permissions.GetServiceAccount(saName, sr.Instance.Namespace, *sr.Client)
	if err != nil {
		sr.Logger.Error(err, "reconcileRoleBindings: failed to retrieve serviceaccount")
		return err
	}

	// reconcile roleBindings for managed namespaces
	if err = sr.reconcileManagedRoleBindings(&roleBindingRequest, sa); err != nil {
		reconciliationError = err
	}

	// reconcile roleBindings for source namespaces
	if err = sr.reconcileSourceRoleBindings(&roleBindingRequest, sa); err != nil {
		reconciliationError = err
	}

	return reconciliationError
}

func (sr *ServerReconciler) reconcileManagedRoleBindings(rbr *permissions.RoleBindingRequest, sa *corev1.ServiceAccount) error {
	var reconciliationError error = nil

	for namespace := range sr.ManagedNamespaces {
		// Skip namespace if can't be retrieved or in terminating state
		ns, err := cluster.GetNamespace(namespace, *sr.Client)
		if err != nil {
			sr.Logger.Error(err, "reconcileManagedRoleBindings: unable to retrieve namesapce", "name", namespace)
			continue
		}
		if ns.DeletionTimestamp != nil {
			sr.Logger.V(1).Info("reconcileManagedRoleBindings: skipping namespace in terminating state", "name", namespace)
			continue
		}

		rbr.Namespace = namespace
		desiredRB := permissions.RequestRoleBinding(*rbr)

		if namespace != sr.Instance.Namespace {
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
				Name:     argoutil.GenerateResourceName(sr.Instance.Name, ArgoCDServerComponent),
			}
		}

		existingRB, err := permissions.GetRoleBinding(desiredRB.Name, desiredRB.Namespace, *sr.Client)
		if err != nil {
			if !errors.IsNotFound(err) {
				sr.Logger.Error(err, "reconcileRoleBinding: failed to retrieve roleBinding", "name", existingRB.Name, "namespace", existingRB.Namespace)
				reconciliationError = err
				continue
			}

			// setting owner reference only to control-plane rolebinding
			if namespace == sr.Instance.Namespace {
				if err = controllerutil.SetControllerReference(sr.Instance, desiredRB, sr.Scheme); err != nil {
					sr.Logger.Error(err, "reconcileRoleBinding: failed to set owner reference for roleBinding", "name", desiredRB.Name, "namespace", desiredRB.Namespace)
					reconciliationError = err
				}
			}

			if err = permissions.CreateRoleBinding(desiredRB, *sr.Client); err != nil {
				sr.Logger.Error(err, "reconcileRoleBinding: failed to create roleBinding", "name", desiredRB.Name, "namespace", desiredRB.Namespace)
				reconciliationError = err
				continue
			}
			sr.Logger.V(0).Info("reconcileRoleBinding: roleBinding created", "name", desiredRB.Name, "namespace", desiredRB.Namespace)
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
			if err = permissions.UpdateRoleBinding(existingRB, *sr.Client); err != nil {
				sr.Logger.Error(err, "reconcileRoleBinding: failed to update roleBinding", "name", existingRB.Name, "namespace", existingRB.Namespace)
				reconciliationError = err
				continue
			}
			sr.Logger.V(0).Info("reconcileRoleBinding: roleBinding updated", "name", existingRB.Name, "namespace", existingRB.Namespace)
			continue
		}
	}
	return reconciliationError
}

func (sr *ServerReconciler) reconcileSourceRoleBindings(rbr *permissions.RoleBindingRequest, sa *corev1.ServiceAccount) error {
	var reconciliationError error = nil

	for namespace := range sr.SourceNamespaces {
		// Skip namespace if can't be retrieved or in terminating state
		ns, err := cluster.GetNamespace(namespace, *sr.Client)
		if err != nil {
			sr.Logger.Error(err, "reconcileSourceRoleBindings: unable to retrieve namesapce", "name", namespace)
			continue
		}
		if ns.DeletionTimestamp != nil {
			sr.Logger.V(1).Info("reconcileSourceRoleBindings: skipping namespace in terminating state", "name", namespace)
			continue
		}

		rbr.Namespace = namespace
		rbr.Name = getSourceNamespaceRBACName(sr.Instance.Name, sr.Instance.Namespace)

		desiredRB := permissions.RequestRoleBinding(*rbr)

		if namespace != sr.Instance.Namespace {
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
			Name:     getSourceNamespaceRBACName(sr.Instance.Name, sr.Instance.Namespace),
		}

		existingRB, err := permissions.GetRoleBinding(desiredRB.Name, desiredRB.Namespace, *sr.Client)
		if err != nil {
			if !errors.IsNotFound(err) {
				sr.Logger.Error(err, "reconcileRoleBinding: failed to retrieve roleBinding", "name", existingRB.Name, "namespace", existingRB.Namespace)
				reconciliationError = err
				continue
			}

			if err = permissions.CreateRoleBinding(desiredRB, *sr.Client); err != nil {
				sr.Logger.Error(err, "reconcileRoleBinding: failed to create roleBinding", "name", desiredRB.Name, "namespace", desiredRB.Namespace)
				reconciliationError = err
				continue
			}
			sr.Logger.V(0).Info("reconcileRoleBinding: roleBinding created", "name", desiredRB.Name, "namespace", desiredRB.Namespace)
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
			if err = permissions.UpdateRoleBinding(existingRB, *sr.Client); err != nil {
				sr.Logger.Error(err, "reconcileRoleBinding: failed to update roleBinding", "name", existingRB.Name, "namespace", existingRB.Namespace)
				reconciliationError = err
				continue
			}
			sr.Logger.V(0).Info("reconcileRoleBinding: roleBinding updated", "name", existingRB.Name, "namespace", existingRB.Namespace)
			continue
		}
	}

	return reconciliationError
}

func (sr *ServerReconciler) DeleteRoleBinding(name, namespace string) error {
	if err := permissions.DeleteRoleBinding(name, namespace, *sr.Client); err != nil {
		sr.Logger.Error(err, "DeleteRoleBinding: failed to delete roleBinding", "name", name, "namespace", namespace)
		return err
	}
	sr.Logger.V(0).Info("DeleteRoleBinding: roleBinding deleted", "name", name, "namespace", namespace)
	return nil
}
