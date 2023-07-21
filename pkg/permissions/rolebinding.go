package permissions

import (
	"context"

	"github.com/jaideepr97/argocd-operator-rewrite/pkg/argoutil"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrlClient "sigs.k8s.io/controller-runtime/pkg/client"
)

type RoleBindingRequest struct {
	Name         string
	InstanceName string
	Namespace    string
	Component    string
	Labels       map[string]string
	Annotations  map[string]string
	RoleRef      rbacv1.RoleRef
	Subjects     []rbacv1.Subject
}

// newRoleBinding returns a new RoleBinding instance.
func newRoleBinding(name, instanceName, namespace, component string, labels, annotations map[string]string, roleRef rbacv1.RoleRef, subjects []rbacv1.Subject) *rbacv1.RoleBinding {
	rbName := argoutil.GenerateResourceName(instanceName, component)
	if name != "" {
		rbName = name
	}
	return &rbacv1.RoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:        rbName,
			Namespace:   namespace,
			Labels:      argoutil.MergeMaps(argoutil.LabelsForCluster(instanceName, component), labels),
			Annotations: annotations,
		},
		RoleRef:  roleRef,
		Subjects: subjects,
	}
}

func RequestRoleBinding(request RoleBindingRequest) *rbacv1.RoleBinding {
	return newRoleBinding(request.Name, request.InstanceName, request.Namespace, request.Component, request.Labels, request.Annotations, request.RoleRef, request.Subjects)
}

func CreateRoleBinding(rb *rbacv1.RoleBinding, client ctrlClient.Client) error {
	return client.Create(context.TODO(), rb)
}

func GetRoleBinding(name, namespace string, client ctrlClient.Client) (*rbacv1.RoleBinding, error) {
	existingRB := &rbacv1.RoleBinding{}
	err := client.Get(context.TODO(), types.NamespacedName{Name: name, Namespace: namespace}, existingRB)
	if err != nil {
		return nil, err
	}
	return existingRB, nil
}

func ListRoleBindings(namespace string, client ctrlClient.Client, listOptions []ctrlClient.ListOption) (*rbacv1.RoleBindingList, error) {
	existingRBs := &rbacv1.RoleBindingList{}
	err := client.List(context.TODO(), existingRBs, listOptions...)
	if err != nil {
		return nil, err
	}
	return existingRBs, nil
}

func UpdateRoleBinding(rb *rbacv1.RoleBinding, client ctrlClient.Client) error {
	_, err := GetRoleBinding(rb.Name, rb.Namespace, client)
	if err != nil {
		return err
	}

	if err = client.Update(context.TODO(), rb); err != nil {
		return err
	}

	return nil
}

func DeleteRoleBinding(name, namespace string, client ctrlClient.Client) error {
	existingRB, err := GetRoleBinding(name, namespace, client)
	if err != nil {
		if !errors.IsNotFound(err) {
			return err
		}
		return nil
	}

	if err := client.Delete(context.TODO(), existingRB); err != nil {
		return err
	}
	return nil
}
