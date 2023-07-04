package notifications

import (
	"reflect"

	"github.com/jaideepr97/argocd-operator-rewrite/pkg/cluster"
	"github.com/jaideepr97/argocd-operator-rewrite/pkg/mutation"
	"github.com/jaideepr97/argocd-operator-rewrite/pkg/permissions"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (nr *NotificationsReconciler) reconcileRole() error {

	nr.Logger.Info("reconciling roles")

	roleRequest := permissions.RoleRequest{
		InstanceName: nr.Instance.Name,
		Component:    ArgoCDNotificationsControllerComponent,
		Client:       nr.Client,
		Mutations:    []mutation.MutateFunc{mutation.ApplyReconcilerMutation},
		Rules:        getPolicyRules(),
		Namespace:    nr.Instance.Namespace,
	}

	desiredRole, err := permissions.RequestRole(roleRequest)
	if err != nil {
		nr.Logger.Error(err, "reconcileRole: failed to request role", "name", desiredRole.Name, "namespace", desiredRole.Namespace)
		nr.Logger.V(1).Info("reconcileRole: one or more mutations could not be applied")
		return err
	}

	namespace, err := cluster.GetNamespace(nr.Instance.Namespace, *nr.Client)
	if err != nil {
		nr.Logger.Error(err, "reconcileRole: failed to retrieve namespace", "name", nr.Instance.Namespace)
		return err
	}
	if namespace.DeletionTimestamp != nil {
		if err := nr.DeleteRole(desiredRole.Name, desiredRole.Namespace); err != nil {
			nr.Logger.Error(err, "reconcileRole: failed to delete role", "name", desiredRole.Name, "namespace", desiredRole.Namespace)
		}
		return err
	}

	existingRole, err := permissions.GetRole(desiredRole.Name, desiredRole.Namespace, *nr.Client)
	if err != nil {
		if !errors.IsNotFound(err) {
			nr.Logger.Error(err, "reconcileRole: failed to retrieve role", "name", existingRole.Name, "namespace", existingRole.Namespace)
			return err
		}

		if err = controllerutil.SetControllerReference(nr.Instance, desiredRole, nr.Scheme); err != nil {
			nr.Logger.Error(err, "reconcileRole: failed to set owner reference for role", "name", desiredRole.Name, "namespace", desiredRole.Namespace)
		}

		if err = permissions.CreateRole(desiredRole, *nr.Client); err != nil {
			nr.Logger.Error(err, "reconcileRole: failed to create role", "name", desiredRole.Name, "namespace", desiredRole.Namespace)
			return err
		}
		nr.Logger.V(0).Info("reconcileRole: role created", "name", desiredRole.Name, "namespace", desiredRole.Namespace)
		return nil
	}

	if !reflect.DeepEqual(existingRole.Rules, desiredRole.Rules) {
		existingRole.Rules = desiredRole.Rules
		if err = permissions.UpdateRole(existingRole, *nr.Client); err != nil {
			nr.Logger.Error(err, "reconcileRole: failed to update role", "name", existingRole.Name, "namespace", existingRole.Namespace)
			return err
		}
	}
	nr.Logger.V(0).Info("reconcileRole: role updated", "name", existingRole.Name, "namespace", existingRole.Namespace)
	return nil
}

func (nr *NotificationsReconciler) DeleteRole(name, namespace string) error {
	if err := permissions.DeleteRole(name, namespace, *nr.Client); err != nil {
		nr.Logger.Error(err, "DeleteRole: failed to delete role", "name", name, "namespace", namespace)
		return err
	}
	nr.Logger.V(0).Info("DeleteRole: role deleted", "name", name, "namespace", namespace)
	return nil
}

func getPolicyRules() []rbacv1.PolicyRule {
	return []rbacv1.PolicyRule{
		{
			APIGroups: []string{
				"argoproj.io",
			},
			Resources: []string{
				"applications",
				"appprojects",
			},
			Verbs: []string{
				"get",
				"list",
				"patch",
				"update",
				"watch",
			},
		},
		{
			APIGroups: []string{
				"",
			},
			Resources: []string{
				"configmaps",
				"secrets",
			},
			Verbs: []string{
				"list",
				"watch",
			},
		},
		{
			APIGroups: []string{
				"",
			},
			ResourceNames: []string{
				"argocd-notifications-cm",
			},
			Resources: []string{
				"configmaps",
			},
			Verbs: []string{
				"get",
			},
		},
		{
			APIGroups: []string{
				"",
			},
			ResourceNames: []string{
				"argocd-notifications-secret",
			},
			Resources: []string{
				"secrets",
			},
			Verbs: []string{
				"get",
			},
		},
	}
}
