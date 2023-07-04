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
	"sigs.k8s.io/controller-runtime/pkg/client"
	ctrlClient "sigs.k8s.io/controller-runtime/pkg/client"
)

type ClusterRoleRequest struct {
	Name                string
	InstanceName        string
	InstanceNamespace   string
	InstanceAnnotations map[string]string
	Component           string
	Rules               []rbacv1.PolicyRule

	// array of functions to mutate role before returning to requester
	Mutations []mutation.MutateFunc
	Client    ctrlClient.Client
}

// newClusterRole returns a new clusterRole instance.
func newClusterRole(name, instanceName, instanceNamespace, component string, instanceAnnotations map[string]string, rules []rbacv1.PolicyRule) *rbacv1.ClusterRole {
	crName := argoutil.GenerateUniqueResourceName(instanceName, instanceNamespace, component)
	if name != "" {
		crName = name
	}

	return &rbacv1.ClusterRole{
		ObjectMeta: metav1.ObjectMeta{
			Name:        crName,
			Labels:      argoutil.LabelsForCluster(instanceName, component),
			Annotations: argoutil.AnnotationsForCluster(instanceName, instanceNamespace, instanceAnnotations),
		},
		Rules: rules,
	}
}

func RequestClusterRole(request ClusterRoleRequest) (*rbacv1.ClusterRole, error) {
	var errCount int
	dlusterROle := newClusterRole(request.Name, request.InstanceName, request.InstanceNamespace, request.Component, request.InstanceAnnotations, request.Rules)

	if len(request.Mutations) > 0 {
		for _, mutation := range request.Mutations {
			err := mutation(nil, *dlusterROle, &request.Client)
			if err != nil {
				// TO DO: log error while invoking mutation
				errCount++
			}
		}
		if errCount > 0 {
			return dlusterROle, fmt.Errorf("RequestRole: one or more mutation functions could not be applied")
		}
	}
	return dlusterROle, nil
}

func CreateClusterRole(clusterRole *rbacv1.ClusterRole, client client.Client) error {
	// TO DO: log creation of clusterRole here (info)
	if err := client.Create(context.TODO(), clusterRole); err != nil {
		// TO DO: log error here (error)
		return fmt.Errorf("CreateClusterRole: failed to create clusterRole %s: %w", clusterRole.Name, err)
	}
	return nil
}

func GetClusterRole(name string, client ctrlClient.Client) (*rbacv1.ClusterRole, error) {
	existingClusterRole := &rbacv1.ClusterRole{}
	err := client.Get(context.TODO(), types.NamespacedName{Name: name}, existingClusterRole)
	if err != nil {
		if errors.IsNotFound(err) {
			// TO DO: log clusterRole not found (debug)
			return nil, fmt.Errorf("GetRole: unable to find clusterRole %s: %w", name, err)
		}
	}
	return existingClusterRole, nil
}

func ListClusterRoles(client ctrlClient.Client, listOptions []ctrlClient.ListOption) (*rbacv1.ClusterRoleList, error) {
	existingClusterRoles := &rbacv1.ClusterRoleList{}
	err := client.List(context.TODO(), existingClusterRoles, listOptions...)
	if err != nil {
		return nil, fmt.Errorf("ListClusterRoles: unable to list clusterRoles: %w", err)
	}

	return existingClusterRoles, nil
}

func UpdateClusterRole(clusterRole *rbacv1.ClusterRole, client ctrlClient.Client) error {
	_, err := GetClusterRole(clusterRole.Name, client)
	if err != nil {
		// TO DO: log role not found (error)
		return fmt.Errorf("UpdateRole: unable to find role %s: %w", clusterRole.Name, err)
	}

	if err = client.Update(context.TODO(), clusterRole); err != nil {
		return fmt.Errorf("UpdateRole: unable to update role %s: %w", clusterRole.Name, err)
	}

	return nil
}

func DeleteClusterRole(name string, client ctrlClient.Client) error {
	existingClusterRole, err := GetClusterRole(name, client)
	if err != nil {
		if !errors.IsNotFound(err) {
			return fmt.Errorf("DeleteClusterRole: unable to retrieve clusterRole %s: %w", name, err)
		}
		// TO DO: log role was not found
		return nil
	}

	// TO DO: log deletion of role here (info)
	if err := client.Delete(context.TODO(), existingClusterRole); err != nil {
		// TO DO: log error here (warn)
		return fmt.Errorf("DeleteClusterRole: failed to delete clusterRole %s: %w", existingClusterRole.Name, err)
	}
	return nil
}
