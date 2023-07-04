package permissions

import (
	"context"
	"fmt"

	"github.com/jaideepr97/argocd-operator-rewrite/pkg/argoutil"
	"github.com/jaideepr97/argocd-operator-rewrite/pkg/mutation"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrlClient "sigs.k8s.io/controller-runtime/pkg/client"
)

type RoleRequest struct {
	Name         string
	InstanceName string
	Namespace    string
	Component    string
	Rules        []rbacv1.PolicyRule

	// array of functions to mutate role before returning to requester
	Mutations []mutation.MutateFunc
	Client    *ctrlClient.Client
}

// newRole returns a new Role instance.
func newRole(name, instanceName, namespace, component string,
	rules []rbacv1.PolicyRule) *rbacv1.Role {
	roleName := argoutil.GenerateResourceName(instanceName, component)
	if name != "" {
		roleName = name
	}
	return &rbacv1.Role{
		ObjectMeta: metav1.ObjectMeta{
			Name:      roleName,
			Namespace: namespace,
			Labels:    argoutil.LabelsForCluster(instanceName, component),
		},
		Rules: rules,
	}
}

func RequestRole(request RoleRequest) (*rbacv1.Role, error) {
	var errCount int
	role := newRole(request.Name, request.InstanceName, request.Namespace, request.Component, request.Rules)

	if len(request.Mutations) > 0 {
		for _, mutation := range request.Mutations {
			err := mutation(nil, *role, request.Client)
			if err != nil {
				errCount++
			}
		}
		if errCount > 0 {
			return role, fmt.Errorf("RequestRole: one or more mutation functions could not be applied")
		}
	}

	return role, nil
}

func CreateRole(role *rbacv1.Role, client ctrlClient.Client) error {
	return client.Create(context.TODO(), role)
}

func GetRole(name, namespace string, client ctrlClient.Client) (*rbacv1.Role, error) {
	existingRole := &rbacv1.Role{}
	err := client.Get(context.TODO(), types.NamespacedName{Name: name, Namespace: namespace}, existingRole)
	if err != nil {
		return nil, err
	}
	return existingRole, nil
}

func ListRoles(namespace string, client ctrlClient.Client, listOptions []ctrlClient.ListOption) (*rbacv1.RoleList, error) {
	existingRoles := &rbacv1.RoleList{}
	err := client.List(context.TODO(), existingRoles, listOptions...)
	if err != nil {
		return nil, err
	}
	return existingRoles, nil
}

func UpdateRole(role *rbacv1.Role, client ctrlClient.Client) error {
	_, err := GetRole(role.Name, role.Namespace, client)
	if err != nil {
		return err
	}

	if err = client.Update(context.TODO(), role); err != nil {
		return err
	}

	return nil
}

func DeleteRole(name, namespace string, client ctrlClient.Client) error {
	existingRole, err := GetRole(name, namespace, client)
	if err != nil {
		if !errors.IsNotFound(err) {
			return err
		}
		return nil
	}

	if err := client.Delete(context.TODO(), existingRole); err != nil {
		return err
	}
	return nil
}
