package appcontroller

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

func (acr *AppControllerReconciler) reconcileRoles() error {

	var reconciliationError error = nil
	acr.Logger.Info("reconciling roles")

	roleRequest := permissions.RoleRequest{
		InstanceName: acr.Instance.Name,
		Component:    ArgoCDApplicationControllerComponent,
		Client:       acr.Client,
		Mutations:    []mutation.MutateFunc{mutation.ApplyReconcilerMutation},
	}

	// reconcile Roles for managed namespaces
	if err := acr.reconcileManagedRoles(&roleRequest); err != nil {
		reconciliationError = err
	}

	// reconcile Roles for source namespaces
	if err := acr.reconcileSourceRoles(&roleRequest); err != nil {
		reconciliationError = err
	}

	return reconciliationError
}

func (acr *AppControllerReconciler) reconcileManagedRoles(rr *permissions.RoleRequest) error {
	var reconciliationError error = nil

	for namespace := range acr.ManagedNamespaces {
		// Skip namespace if can't be retrieved or in terminating state
		ns, err := cluster.GetNamespace(namespace, *acr.Client)
		if err != nil {
			acr.Logger.Error(err, "reconcileManagedRoles: unable to retrieve namesapce", "name", namespace)
			continue
		}
		if ns.DeletionTimestamp != nil {
			acr.Logger.V(1).Info("reconcileManagedRoles: skipping namespace in terminating state", "name", namespace)
			continue
		}

		rules := policyRuleForManagedNamespace()
		// control-plane rules for namespace scoped instance
		if namespace == acr.Instance.Namespace {
			rules = policyRuleForNamespaceScope()
		}
		rr.Rules = rules
		rr.Namespace = namespace

		desiredRole, err := permissions.RequestRole(*rr)
		if err != nil {
			acr.Logger.Error(err, "reconcileManagedRoles: failed to request role", "name", desiredRole.Name, "namespace", desiredRole.Namespace)
			acr.Logger.V(1).Info("reconcileManagedRoles: one or more mutations could not be applied")
			reconciliationError = err
			continue
		}

		if namespace != acr.Instance.Namespace {
			// add special label for resource management to role
			if len(desiredRole.Labels) == 0 {
				desiredRole.Labels = make(map[string]string)
			}
			desiredRole.Labels[common.ArgoCDKeyRBACType] = common.ArgoCDRBACTypeResourceMananagement
		}

		customRoleName := getCustomRoleName()
		existingRole, err := permissions.GetRole(desiredRole.Name, desiredRole.Namespace, *acr.Client)
		if err != nil {
			if !errors.IsNotFound(err) {
				acr.Logger.Error(err, "reconcileManagedRoles: failed to retrieve role", "name", existingRole.Name, "namespace", existingRole.Namespace)
				reconciliationError = err
				continue
			}

			if customRoleName != "" {
				// skip default role creation/reconciliation
				acr.Logger.V(1).Info("reconcileManagedRoles: skip reconciliation in favor of custom role", "name", customRoleName)
				continue
			}

			// setting owner reference only for control-plane role
			if namespace == acr.Instance.Namespace {
				if err = controllerutil.SetControllerReference(acr.Instance, desiredRole, acr.Scheme); err != nil {
					acr.Logger.Error(err, "reconcileManagedRoles: failed to set owner reference for role", "name", desiredRole.Name, "namespace", desiredRole.Namespace)
					reconciliationError = err
				}
			}

			if err = permissions.CreateRole(desiredRole, *acr.Client); err != nil {
				acr.Logger.Error(err, "reconcileManagedRoles: failed to create role", "name", desiredRole.Name, "namespace", desiredRole.Namespace)
				reconciliationError = err
				continue
			}
			acr.Logger.V(0).Info("reconcileManagedRoles: role created", "name", desiredRole.Name, "namespace", desiredRole.Namespace)
			continue
		}

		if customRoleName != "" {
			// Delete default role in namespace
			if err := acr.DeleteRole(desiredRole.Name, desiredRole.Namespace); err != nil {
				acr.Logger.Error(err, "reconcileManagedRoles: failed to delete role", "name", desiredRole.Name, "namespace", desiredRole.Namespace)
				reconciliationError = err
				continue
			}
			acr.Logger.V(0).Info("reconcileManagedRoles: role deleted", "name", desiredRole.Name, "namespace", desiredRole.Namespace)
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
			if err = permissions.UpdateRole(existingRole, *acr.Client); err != nil {
				acr.Logger.Error(err, "reconcileManagedRoles: failed to update role", "name", existingRole.Name, "namespace", existingRole.Namespace)
				reconciliationError = err
				continue
			}
			acr.Logger.V(0).Info("reconcileManagedRoles: role updated", "name", existingRole.Name, "namespace", existingRole.Namespace)
			continue
		}
	}

	return reconciliationError
}

func (acr *AppControllerReconciler) reconcileSourceRoles(rr *permissions.RoleRequest) error {
	var reconciliationError error = nil

	for namespace, _ := range acr.SourceNamespaces {
		// Skip namespace if can't be retrieved or in terminating state
		ns, err := cluster.GetNamespace(namespace, *acr.Client)
		if err != nil {
			acr.Logger.Error(err, "reconcileSourceRoles: unable to retrieve namesapce", "name", namespace)
			continue
		}
		if ns.DeletionTimestamp != nil {
			acr.Logger.V(1).Info("reconcileSourceRoles: skipping namespace in terminating state", "name", namespace)
			continue
		}

		rr.Rules = policyRuleForSourceNamespace()
		rr.Namespace = namespace
		rr.Name = getSourceNamespaceRBACName(acr.Instance.Name, acr.Instance.Namespace)

		desiredRole, err := permissions.RequestRole(*rr)
		if err != nil {
			acr.Logger.Error(err, "reconcileSourceRoles: failed to request role", "name", desiredRole.Name, "namespace", desiredRole.Namespace)
			acr.Logger.V(1).Info("reconcileSourceRoles: one or more mutations could not be applied")
			reconciliationError = err
			continue
		}

		// add special label for source namespace role
		if len(desiredRole.Labels) == 0 {
			desiredRole.Labels = make(map[string]string)
		}
		desiredRole.Labels[common.ArgoCDKeyRBACType] = common.ArgoCDRBACTypeAppManagement

		existingRole, err := permissions.GetRole(desiredRole.Name, desiredRole.Namespace, *acr.Client)
		if err != nil {
			if !errors.IsNotFound(err) {
				acr.Logger.Error(err, "reconcileSourceRoles: failed to retrieve role", "name", existingRole.Name, "namespace", existingRole.Namespace)
				reconciliationError = err
				continue
			}

			if err = permissions.CreateRole(desiredRole, *acr.Client); err != nil {
				acr.Logger.Error(err, "reconcileSourceRoles: failed to create role", "name", existingRole.Name, "namespace", existingRole.Namespace)
				reconciliationError = err
				continue
			}
			acr.Logger.V(0).Info("reconcileSourceRoles: role created", "name", desiredRole.Name, "namespace", desiredRole.Namespace)
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
			if err = permissions.UpdateRole(existingRole, *acr.Client); err != nil {
				acr.Logger.Error(err, "reconcileSourceRoles: failed to update role", "name", existingRole.Name, "namespace", existingRole.Namespace)
				reconciliationError = err
				continue
			}
			acr.Logger.V(0).Info("reconcileSourceRoles: role updated", "name", existingRole.Name, "namespace", existingRole.Namespace)
			continue
		}
	}
	return reconciliationError
}

func (acr *AppControllerReconciler) DeleteRole(name, namespace string) error {
	if err := permissions.DeleteRole(name, namespace, *acr.Client); err != nil {
		acr.Logger.Error(err, "DeleteRole: failed to delete role", "name", name, "namespace", namespace)
		return err
	}
	acr.Logger.V(0).Info("DeleteRole: role deleted", "name", name, "namespace", namespace)
	return nil
}

func getCustomRoleName() string {
	return argoutil.GetEnvVar(ArgoCDControllerClusterRoleEnvName)
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
				"create",
				"get",
				"list",
				"patch",
				"update",
				"watch",
				"delete",
			},
		},
	}
}
