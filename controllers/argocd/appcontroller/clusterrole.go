package appcontroller

import (
	"reflect"

	"github.com/jaideepr97/argocd-operator-rewrite/pkg/mutation"
	"github.com/jaideepr97/argocd-operator-rewrite/pkg/permissions"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/errors"
)

func (acr *AppControllerReconciler) reconcileClusterRole() error {

	acr.Logger.V(0).Info("reconciling clusterRole")

	clusterRoleRequest := permissions.ClusterRoleRequest{
		InstanceName:      acr.Instance.Name,
		InstanceNamespace: acr.Instance.Namespace,
		Annotations:       acr.Instance.Annotations,
		Component:         ArgoCDApplicationControllerComponent,
		Rules:             policyRuleForClusterScope(),
		Mutations:         []mutation.MutateFunc{mutation.ApplyReconcilerMutation},
		Client:            *acr.Client,
	}

	desiredClusterRole, err := permissions.RequestClusterRole(clusterRoleRequest)
	if err != nil {
		acr.Logger.Error(err, "reconcileClusterRole: failed to request clusterRole")
		acr.Logger.V(1).Info("reconcileClusterRole: one or more mutations could not be applied")
	}

	existingClusterRole, err := permissions.GetClusterRole(desiredClusterRole.Name, *acr.Client)
	if err != nil {
		if !errors.IsNotFound(err) {
			acr.Logger.Error(err, "reconcileClusterRole: failed to retrieve clusterRole")
			return err
		}

		if err = permissions.CreateClusterRole(desiredClusterRole, *acr.Client); err != nil {
			acr.Logger.Error(err, "reconcileClusterRole: failed to create clusterRole")
			return err
		}
		acr.Logger.V(0).Info("reconcileClusterRole: clusterRole created")
		return nil
	}

	// if !acr.ClusterScoped {
	// 	if err := acr.DeleteClusterRole(desiredClusterRole.Name); err != nil {
	// 		acr.Logger.Error(err, "reconcileClusterRole:failed to delete clusterRole")
	// 		return err
	// 	}
	// 	acr.Logger.V(0).Info("reconcileClusterRole: clusterRole deleted")
	// 	return nil
	// }

	if !reflect.DeepEqual(existingClusterRole.Rules, desiredClusterRole.Rules) {
		existingClusterRole.Rules = desiredClusterRole.Rules
		if err = permissions.UpdateClusterRole(existingClusterRole, *acr.Client); err != nil {
			acr.Logger.Error(err, "reconcileClusterRole: failed to update clusterRole")
			return err
		}
		acr.Logger.V(0).Info("reconcileClusterRole: clusterRole updated")
	}
	return nil
}

func (acr *AppControllerReconciler) DeleteClusterRole(name string) error {
	return permissions.DeleteClusterRole(name, *acr.Client)
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
				"*",
			},
		},
	}
}
