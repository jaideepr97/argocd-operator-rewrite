package networking

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	ctrlClient "sigs.k8s.io/controller-runtime/pkg/client"
)

func CreateService(service *corev1.Service, client ctrlClient.Client) error {
	if err := client.Create(context.TODO(), service); err != nil {
		return fmt.Errorf("CreateService: failed to create service %s in namespace %s: %w", service.Name, service.Namespace, err)
	}
	return nil
}

func ListServices(namespace string, client ctrlClient.Client, listOptions []ctrlClient.ListOption) (*corev1.ServiceList, error) {
	existingServices := &corev1.ServiceList{}
	err := client.List(context.TODO(), existingServices, listOptions...)
	if err != nil {
		return nil, fmt.Errorf("ListServices: unable to list services in namespace %s: %w", namespace, err)
	}

	return existingServices, nil
}

func GetService(name, namespace string, client ctrlClient.Client) (*corev1.Service, error) {
	existingService := &corev1.Service{}
	err := client.Get(context.TODO(), types.NamespacedName{Name: name, Namespace: namespace}, existingService)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, fmt.Errorf("GetService: unable to find service %s in namespace %s: %w", name, namespace, err)
		}
	}

	return existingService, nil
}

func UpdateService(service *corev1.Service, client ctrlClient.Client) error {
	_, err := GetService(service.Name, service.Namespace, client)
	if err != nil {
		return fmt.Errorf("UpdateService: unable to find service %s in namespace %s: %w", service.Name, service.Namespace, err)
	}

	if err = client.Update(context.TODO(), service); err != nil {
		return fmt.Errorf("UpdateService: unable to update service %s in namespace %s: %w", service.Name, service.Namespace, err)
	}

	return nil
}

func DeleteService(name, namespace string, client ctrlClient.Client) error {
	existingService, err := GetService(name, namespace, client)
	if err != nil {
		return fmt.Errorf("DeleteService: unable to get service %s in namespace %s: %w", name, namespace, err)
	}

	err = client.Delete(context.TODO(), existingService, &ctrlClient.DeleteOptions{})
	if err != nil {
		return fmt.Errorf("DeleteService: unable to delete service %s in namespace %s: %w", name, namespace, err)
	}

	return nil
}
