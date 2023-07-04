package applicationset

import (
	"reflect"

	"github.com/jaideepr97/argocd-operator-rewrite/pkg/cluster"
	"github.com/jaideepr97/argocd-operator-rewrite/pkg/mutation"
	"github.com/jaideepr97/argocd-operator-rewrite/pkg/permissions"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (asr *ApplicationSetReconciler) reconcileRole() error {

	asr.Logger.V(0).Info("reconciling roles")

	roleRequest := permissions.RoleRequest{
		InstanceName: asr.Instance.Name,
		Component:    ArgoCDApplicationSetControllerComponent,
		Client:       asr.Client,
		Mutations:    []mutation.MutateFunc{mutation.ApplyReconcilerMutation},
		Rules:        getPolicyRules(),
		Namespace:    asr.Instance.Namespace,
	}

	desiredRole, err := permissions.RequestRole(roleRequest)
	if err != nil {
		asr.Logger.Error(err, "reconcileRole: failed to request role", "name", desiredRole.Name, "namespace", desiredRole.Namespace)
		asr.Logger.V(1).Info("reconcileRole: one or more mutations could not be applied")
		return err
	}

	namespace, err := cluster.GetNamespace(asr.Instance.Namespace, *asr.Client)
	if err != nil {
		asr.Logger.Error(err, "reconcileRole: failed to retrieve namespace", "name", asr.Instance.Namespace)
		return err
	}

	if namespace.DeletionTimestamp != nil {
		if err := asr.DeleteRole(desiredRole.Name, desiredRole.Namespace); err != nil {
			asr.Logger.Error(err, "reconcileRole: failed to delete role", "name", desiredRole.Name, "namespace", desiredRole.Namespace)
		}
		return err
	}

	existingRole, err := permissions.GetRole(desiredRole.Name, desiredRole.Namespace, *asr.Client)
	if err != nil {
		if !errors.IsNotFound(err) {
			asr.Logger.Error(err, "reconcileRole: failed to retrieve role", "name", existingRole.Name, "namespace", existingRole.Namespace)
			return err
		}

		if err = controllerutil.SetControllerReference(asr.Instance, desiredRole, asr.Scheme); err != nil {
			asr.Logger.Error(err, "reconcileRole: failed to set owner reference for role", "name", desiredRole.Name, "namespace", desiredRole.Namespace)
		}

		if err = permissions.CreateRole(desiredRole, *asr.Client); err != nil {
			asr.Logger.Error(err, "reconcileRole: failed to create role", "name", desiredRole.Name, "namespace", desiredRole.Namespace)
			return err
		}
		asr.Logger.V(0).Info("reconcileRole: role created", "name", desiredRole.Name, "namespace", desiredRole.Namespace)
		return nil
	}

	if !reflect.DeepEqual(existingRole.Rules, desiredRole.Rules) {
		existingRole.Rules = desiredRole.Rules
		if err = permissions.UpdateRole(existingRole, *asr.Client); err != nil {
			asr.Logger.Error(err, "reconcileRole: failed to update role", "name", existingRole.Name, "namespace", existingRole.Namespace)
			return err
		}
	}
	asr.Logger.V(0).Info("reconcileRole: role updated", "name", existingRole.Name, "namespace", existingRole.Namespace)
	return nil
}

func (asr *ApplicationSetReconciler) DeleteRole(name, namespace string) error {
	if err := permissions.DeleteRole(name, namespace, *asr.Client); err != nil {
		asr.Logger.Error(err, "DeleteRole: failed to delete role", "name", name, "namespace", namespace)
		return err
	}
	asr.Logger.V(0).Info("DeleteRole: role deleted", "name", name, "namespace", namespace)
	return nil
}

func getPolicyRules() []rbacv1.PolicyRule {
	return []rbacv1.PolicyRule{

		// ApplicationSet
		{
			APIGroups: []string{"argoproj.io"},
			Resources: []string{
				"applications",
				"applicationsets",
				"appprojects",
				"applicationsets/finalizers",
			},
			Verbs: []string{
				"create",
				"delete",
				"get",
				"list",
				"patch",
				"update",
				"watch",
			},
		},
		// ApplicationSet Status
		{
			APIGroups: []string{"argoproj.io"},
			Resources: []string{
				"applicationsets/status",
			},
			Verbs: []string{
				"get",
				"patch",
				"update",
			},
		},

		// Events
		{
			APIGroups: []string{""},
			Resources: []string{
				"events",
			},
			Verbs: []string{
				"create",
				"delete",
				"get",
				"list",
				"patch",
				"update",
				"watch",
			},
		},

		// Read Secrets/ConfigMaps
		{
			APIGroups: []string{""},
			Resources: []string{
				"secrets",
				"configmaps",
			},
			Verbs: []string{
				"get",
				"list",
				"watch",
			},
		},

		// Read Deployments
		{
			APIGroups: []string{"apps", "extensions"},
			Resources: []string{
				"deployments",
			},
			Verbs: []string{
				"get",
				"list",
				"watch",
			},
		},
	}

}
