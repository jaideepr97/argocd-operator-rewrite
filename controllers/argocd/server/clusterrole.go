package server

import (
	"reflect"

	"github.com/jaideepr97/argocd-operator-rewrite/pkg/mutation"
	"github.com/jaideepr97/argocd-operator-rewrite/pkg/permissions"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/errors"
)

func (sr *ServerReconciler) reconcileClusterRole() error {

	sr.Logger.V(0).Info("reconciling clusterRole")

	clusterRoleRequest := permissions.ClusterRoleRequest{
		InstanceName:      sr.Instance.Name,
		InstanceNamespace: sr.Instance.Namespace,
		Annotations:       sr.Instance.Annotations,
		Component:         ArgoCDServerComponent,
		Rules:             policyRuleForClusterScope(),
		Mutations:         []mutation.MutateFunc{mutation.ApplyReconcilerMutation},
		Client:            *sr.Client,
	}

	desiredClusterRole, err := permissions.RequestClusterRole(clusterRoleRequest)
	if err != nil {
		sr.Logger.Error(err, "reconcileClusterRole: failed to request clusterRole")
		sr.Logger.V(1).Info("reconcileClusterRole: one or more mutations could not be applied")
	}

	existingClusterRole, err := permissions.GetClusterRole(desiredClusterRole.Name, *sr.Client)
	if err != nil {
		if !errors.IsNotFound(err) {
			sr.Logger.Error(err, "reconcileClusterRole: failed to retrieve clusterRole")
			return err
		}

		if err = permissions.CreateClusterRole(desiredClusterRole, *sr.Client); err != nil {
			sr.Logger.Error(err, "reconcileClusterRole: failed to create clusterRole")
			return err
		}
		sr.Logger.V(0).Info("reconcileClusterRole: clusterRole created")
		return nil
	}

	// if !sr.ClusterScoped {
	// 	if err := sr.DeleteClusterRole(desiredClusterRole.Name); err != nil {
	// 		sr.Logger.Error(err, "reconcileClusterRole:failed to delete clusterRole")
	// 		return err
	// 	}
	// 	sr.Logger.V(0).Info("reconcileClusterRole: clusterRole deleted")
	// 	return nil
	// }

	if !reflect.DeepEqual(existingClusterRole.Rules, desiredClusterRole.Rules) {
		existingClusterRole.Rules = desiredClusterRole.Rules
		if err = permissions.UpdateClusterRole(existingClusterRole, *sr.Client); err != nil {
			sr.Logger.Error(err, "reconcileClusterRole: failed to update clusterRole")
			return err
		}
		sr.Logger.V(0).Info("reconcileClusterRole: clusterRole updated")
	}
	return nil
}

func (sr *ServerReconciler) DeleteClusterRole(name string) error {
	return permissions.DeleteClusterRole(name, *sr.Client)
}

func policyRuleForClusterScope() []rbacv1.PolicyRule {
	return []rbacv1.PolicyRule{
		{
			APIGroups: []string{
				"*",
			},
			Resources: []string{
				"*",
			},
			Verbs: []string{
				"get",
				"delete",
				"patch",
			},
		},
		{
			APIGroups: []string{
				"argoproj.io",
			},
			Resources: []string{
				"applications",
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
			Resources: []string{
				"events",
			},
			Verbs: []string{
				"list",
			},
		},
	}
}
