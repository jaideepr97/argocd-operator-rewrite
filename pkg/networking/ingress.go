package networking

import (
	"context"
	"fmt"

	networkingv1 "k8s.io/api/networking/v1"
	ctrlClient "sigs.k8s.io/controller-runtime/pkg/client"
)

func CreateIngress(ingress networkingv1.Ingress, client ctrlClient.Client) error {
	if err := client.Create(context.TODO(), &ingress); err != nil {
		return fmt.Errorf("CreateIngress: failed to create ingress %s in namespace %s: %w", ingress.Name, ingress.Namespace, err)
	}
	return nil
}

func ListIngresses(namespace string, client ctrlClient.Client, listOptions []ctrlClient.ListOption) (*networkingv1.IngressList, error) {
	existingIngresses := &networkingv1.IngressList{}
	err := client.List(context.TODO(), existingIngresses, listOptions...)
	if err != nil {
		return nil, fmt.Errorf("ListIngresses: unable to list ingresses in namespace %s: %w", namespace, err)
	}
	return existingIngresses, nil
}

func GetIngress(name, namespace string, client ctrlClient.Client) (*networkingv1.Ingress, error) {
	existingIngress := &networkingv1.Ingress{}
	err := client.Get(context.TODO(), ctrlClient.ObjectKey{Namespace: namespace, Name: name}, existingIngress)
	if err != nil {
		return nil, fmt.Errorf("GetIngress: unable to get ingress %s in namespace %s: %w", name, namespace, err)
	}
	return existingIngress, nil
}

func UpdateIngress(ingress networkingv1.Ingress, client ctrlClient.Client) error {
	_, err := GetIngress(ingress.Name, ingress.Namespace, client)
	if err != nil {
		return fmt.Errorf("UpdateIngress: unable to find ingress %s in namespace %s: %w", ingress.Name, ingress.Namespace, err)
	}

	if err := client.Update(context.TODO(), &ingress); err != nil {
		return fmt.Errorf("UpdateIngress: unable to update ingress %s in namespace %s: %w", ingress.Name, ingress.Namespace, err)
	}
	return nil
}

func DeleteIngress(name, namespace string, client ctrlClient.Client) error {
	existingIngress, err := GetIngress(name, namespace, client)
	if err != nil {
		return fmt.Errorf("DeleteIngress: unable to get ingress %s in namespace %s: %w", name, namespace, err)
	}

	err = client.Delete(context.TODO(), existingIngress, &ctrlClient.DeleteOptions{})
	if err != nil {
		return fmt.Errorf("DeleteIngress: unable to delete ingress %s in namespace %s: %w", name, namespace, err)
	}
	return nil
}
