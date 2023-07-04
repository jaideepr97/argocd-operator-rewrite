package server

import (
	"reflect"

	"github.com/jaideepr97/argocd-operator-rewrite/common"
	"github.com/jaideepr97/argocd-operator-rewrite/pkg/argoutil"
	"github.com/jaideepr97/argocd-operator-rewrite/pkg/cluster"
	"github.com/jaideepr97/argocd-operator-rewrite/pkg/mutation"
	"github.com/jaideepr97/argocd-operator-rewrite/pkg/permissions"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (sr *ServerReconciler) reconcileRoles() error {
	var reconciliationError error = nil

	sr.Logger.V(0).Info("reconciling roles")
	roleRequest := permissions.RoleRequest{
		InstanceName: sr.Instance.Name,
		Component:    ArgoCDServerComponent,
		Client:       sr.Client,
		Mutations:    []mutation.MutateFunc{mutation.ApplyReconcilerMutation},
	}
	// reconcile Roles for managed namespaces
	if err := sr.reconcileManagedRoles(&roleRequest); err != nil {
		reconciliationError = err
	}

	// reconcile Roles for source namespaces
	if err := sr.reconcileSourceRoles(&roleRequest); err != nil {
		reconciliationError = err
	}

	return reconciliationError
}

func (sr *ServerReconciler) reconcileManagedRoles(rr *permissions.RoleRequest) error {
	var reconciliationError error = nil

	for namespace := range sr.ManagedNamespaces {
		// Skip namespace if can't be retrieved or in terminating state
		ns, err := cluster.GetNamespace(namespace, *sr.Client)
		if err != nil {
			sr.Logger.Error(err, "reconcileManagedRoles: unable to retrieve namesapce", "name", namespace)
			continue
		}
		if ns.DeletionTimestamp != nil {
			sr.Logger.V(1).Info("reconcileManagedRoles: skipping namespace in terminating state", "name", namespace)
			continue
		}

		rules := policyRuleForManagedNamespace()
		// control-plane rules for namespace scoped instance
		if namespace == sr.Instance.Namespace {
			rules = policyRuleForNamespaceScope()
		}
		rr.Rules = rules
		rr.Namespace = namespace

		desiredRole, err := permissions.RequestRole(*rr)
		if err != nil {
			sr.Logger.Error(err, "reconcileRole: failed to request role", "name", desiredRole.Name, "namespace", desiredRole.Namespace)
			sr.Logger.V(1).Info("reconcileRole: one or more mutations could not be applied")
			reconciliationError = err
			continue
		}

		if namespace != sr.Instance.Namespace {
			// add special label for resource management to role
			if len(desiredRole.Labels) == 0 {
				desiredRole.Labels = make(map[string]string)
			}
			desiredRole.Labels[common.ArgoCDKeyRBACType] = common.ArgoCDRBACTypeResourceMananagement
		}

		customRoleName := getCustomRoleName()
		existingRole, err := permissions.GetRole(desiredRole.Name, desiredRole.Namespace, *sr.Client)
		if err != nil {
			if !errors.IsNotFound(err) {
				sr.Logger.Error(err, "reconcileRole: failed to retrieve role", "name", existingRole.Name, "namespace", existingRole.Namespace)
				reconciliationError = err
				continue
			}

			if customRoleName != "" {
				// skip default role creation/reconciliation
				sr.Logger.V(1).Info("reconcileRole: skip reconciliation in favor of custom role", "name", customRoleName)
				continue
			}

			// setting owner reference only for control-plane role
			if namespace == sr.Instance.Namespace {
				if err = controllerutil.SetControllerReference(sr.Instance, desiredRole, sr.Scheme); err != nil {
					sr.Logger.Error(err, "reconcileRole: failed to set owner reference for role", "name", desiredRole.Name, "namespace", desiredRole.Namespace)
					reconciliationError = err
				}
			}

			if err = permissions.CreateRole(desiredRole, *sr.Client); err != nil {
				sr.Logger.Error(err, "reconcileRole: failed to create role", "name", desiredRole.Name, "namespace", desiredRole.Namespace)
				reconciliationError = err
				continue
			}
			sr.Logger.V(0).Info("reconcileRole: role created", "name", desiredRole.Name, "namespace", desiredRole.Namespace)
			continue
		}

		if customRoleName != "" {
			// Delete default role in namespace
			if err := sr.DeleteRole(desiredRole.Name, desiredRole.Namespace); err != nil {
				sr.Logger.Error(err, "reconcileRole: failed to delete role", "name", desiredRole.Name, "namespace", desiredRole.Namespace)
				reconciliationError = err
				continue
			}
			sr.Logger.V(0).Info("reconcileRole: role deleted", "name", desiredRole.Name, "namespace", desiredRole.Namespace)
			continue
		}

		roleChanged := false
		if !reflect.DeepEqual(existingRole.Rules, desiredRole.Rules) {
			existingRole.Rules = desiredRole.Rules
			roleChanged = true
		} else if !reflect.DeepEqual(existingRole.Labels, desiredRole.Labels) {
			existingRole.Labels = desiredRole.Labels
			roleChanged = true
		}

		if roleChanged {
			if err = permissions.UpdateRole(existingRole, *sr.Client); err != nil {
				sr.Logger.Error(err, "reconcileRole: failed to update role", "name", existingRole.Name, "namespace", existingRole.Namespace)
				reconciliationError = err
				continue
			}
			sr.Logger.V(0).Info("reconcileRole: role updated", "name", existingRole.Name, "namespace", existingRole.Namespace)
			continue
		}
	}

	return reconciliationError
}

func (sr *ServerReconciler) reconcileSourceRoles(rr *permissions.RoleRequest) error {
	var reconciliationError error = nil

	for namespace, _ := range sr.SourceNamespaces {
		// Skip namespace if can't be retrieved or in terminating state
		ns, err := cluster.GetNamespace(namespace, *sr.Client)
		if err != nil {
			sr.Logger.Error(err, "reconcileSourceRoles: unable to retrieve namesapce", "name", namespace)
			continue
		}
		if ns.DeletionTimestamp != nil {
			sr.Logger.V(1).Info("reconcileSourceRoles: skipping namespace in terminating state", "name", namespace)
			continue
		}

		rr.Rules = policyRuleForSourceNamespace()
		rr.Namespace = namespace
		rr.Name = getSourceNamespaceRBACName(sr.Instance.Name, sr.Instance.Namespace)

		desiredRole, err := permissions.RequestRole(*rr)
		if err != nil {
			sr.Logger.Error(err, "reconcileRole: failed to request role", "name", desiredRole.Name, "namespace", desiredRole.Namespace)
			sr.Logger.V(1).Info("reconcileRole: one or more mutations could not be applied")
			reconciliationError = err
			continue
		}

		// add special label for source namespace role
		if len(desiredRole.Labels) == 0 {
			desiredRole.Labels = make(map[string]string)
		}
		desiredRole.Labels[common.ArgoCDKeyRBACType] = common.ArgoCDRBACTypeAppManagement

		existingRole, err := permissions.GetRole(desiredRole.Name, desiredRole.Namespace, *sr.Client)
		if err != nil {
			if !errors.IsNotFound(err) {
				sr.Logger.Error(err, "reconcileRole: failed to retrieve role", "name", existingRole.Name, "namespace", existingRole.Namespace)
				reconciliationError = err
				continue
			}

			if err = permissions.CreateRole(desiredRole, *sr.Client); err != nil {
				sr.Logger.Error(err, "reconcileRole: failed to create role", "name", existingRole.Name, "namespace", existingRole.Namespace)
				reconciliationError = err
				continue
			}
			sr.Logger.V(0).Info("reconcileRole: role created", "name", desiredRole.Name, "namespace", desiredRole.Namespace)
			continue
		}

		roleChanged := false
		if !reflect.DeepEqual(existingRole.Rules, desiredRole.Rules) {
			existingRole.Rules = desiredRole.Rules
			roleChanged = true
		} else if !reflect.DeepEqual(existingRole.Labels, desiredRole.Labels) {
			existingRole.Labels = desiredRole.Labels
			roleChanged = true
		}

		if roleChanged {
			if err = permissions.UpdateRole(existingRole, *sr.Client); err != nil {
				sr.Logger.Error(err, "reconcileRole: failed to update role", "name", existingRole.Name, "namespace", existingRole.Namespace)
				reconciliationError = err
				continue
			}
			sr.Logger.V(0).Info("reconcileRole: role updated", "name", existingRole.Name, "namespace", existingRole.Namespace)
			continue
		}
	}
	return reconciliationError
}

func (sr *ServerReconciler) DeleteRole(name, namespace string) error {
	if err := permissions.DeleteRole(name, namespace, *sr.Client); err != nil {
		sr.Logger.Error(err, "DeleteRole: failed to delete role", "name", name, "namespace", namespace)
		return err
	}
	sr.Logger.V(0).Info("DeleteRole: role deleted", "name", name, "namespace", namespace)
	return nil
}

func getCustomRoleName() string {
	return argoutil.GetEnvVar(ArgoCDServerClusterRoleEnvName)
}

// TO DO: restrict this by specifying exactly what is needed
func policyRuleForNamespaceScope() []rbacv1.PolicyRule {
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

// TO DO: restrict this by specifying exactly what is needed (no app permissions)
func policyRuleForManagedNamespace() []rbacv1.PolicyRule {
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

func policyRuleForSourceNamespace() []rbacv1.PolicyRule {
	return []rbacv1.PolicyRule{
		{
			APIGroups: []string{
				"argoproj.io",
			},
			Resources: []string{
				"applications",
			},
			Verbs: []string{
				"get",
				"list",
				"watch",
			},
		},
	}
}
