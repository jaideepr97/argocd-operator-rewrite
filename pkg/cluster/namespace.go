package cluster

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	ctrlClient "sigs.k8s.io/controller-runtime/pkg/client"
)

func GetNamespace(name string, client ctrlClient.Client) (*corev1.Namespace, error) {
	existingNamespace := &corev1.Namespace{}
	err := client.Get(context.TODO(), types.NamespacedName{Name: name}, existingNamespace)
	if err != nil {
		if errors.IsNotFound(err) {
			// TO DO: log namespace not found (debug)
			return nil, fmt.Errorf("GetNamespace: unable to find namespace %s: %w", name, err)
		}
	}
	return existingNamespace, nil
}

func ListNamespaces(client ctrlClient.Client, listOptions []ctrlClient.ListOption) (*corev1.NamespaceList, error) {
	existingNamespaces := &corev1.NamespaceList{}
	err := client.List(context.TODO(), existingNamespaces, listOptions...)
	if err != nil {
		return nil, fmt.Errorf("ListNamespaces: unable to list namespaces: %w", err)
	}

	return existingNamespaces, nil
}

func UpdateNamespace(namespace *corev1.Namespace, client ctrlClient.Client) error {
	_, err := GetNamespace(namespace.Name, client)
	if err != nil {
		// TO DO: log namespace not found (error)
		return fmt.Errorf("UpdateNamespace: unable to find namespace %s: %w", namespace.Name, err)
	}

	if err = client.Update(context.TODO(), namespace); err != nil {
		return fmt.Errorf("UpdateNamespace: unable to update namespace %s: %w", namespace.Name, err)
	}
	return nil
}
