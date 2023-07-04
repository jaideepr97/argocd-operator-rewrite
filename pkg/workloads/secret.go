package workloads

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	ctrlClient "sigs.k8s.io/controller-runtime/pkg/client"
)

func CreateSecret(secret *corev1.Secret, client ctrlClient.Client) error {
	// TO DO: log creation of secret here (info)
	if err := client.Create(context.TODO(), secret); err != nil {
		// TO DO: log error here (error)
		return fmt.Errorf("CreateSecret: failed to create secret %s in namespace %s: %w", secret.Name, secret.Namespace, err)
	}
	return nil
}

func ListSecrets(namespace string, client ctrlClient.Client, listOptions ctrlClient.ListOption) (*corev1.SecretList, error) {
	existingSecrets := &corev1.SecretList{}
	err := client.List(context.TODO(), existingSecrets, listOptions)
	if err != nil {
		return nil, fmt.Errorf("ListSecrets: unable to list secrets in namespace %s: %w", namespace, err)
	}

	return existingSecrets, nil
}

func GetSecret(name, namespace string, client ctrlClient.Client) (*corev1.Secret, error) {
	existingSecret := &corev1.Secret{}
	err := client.Get(context.TODO(), types.NamespacedName{Name: name, Namespace: namespace}, existingSecret)
	if err != nil {
		if errors.IsNotFound(err) {
			// TO DO: log secret not found (debug)
			return nil, fmt.Errorf("GetSecret: unable to find secret %s in namespace %s: %w", name, namespace, err)
		}
	}

	return existingSecret, nil
}

func UpdateSecret(secret *corev1.Secret, client ctrlClient.Client) error {
	_, err := GetSecret(secret.Name, secret.Namespace, client)
	if err != nil {
		// TO DO: log secret not found (error)
		return fmt.Errorf("UpdateSecret: unable to find secret %s in namespace %s: %w", secret.Name, secret.Namespace, err)
	}

	if err = client.Update(context.TODO(), secret); err != nil {
		return fmt.Errorf("UpdateSecret: unable to update secret %s in namespace %s: %w", secret.Name, secret.Namespace, err)
	}

	return nil
}

// DeleteSecret deletes the specified Secret in the specified namespace
func DeleteSecret(name, namespace string, client ctrlClient.Client) error {
	existingSecret, err := GetSecret(name, namespace, client)
	if err != nil {
		return fmt.Errorf("DeleteSecret: unable to get secret %s in namespace %s: %w", name, namespace, err)
	}

	err = client.Delete(context.TODO(), existingSecret, &ctrlClient.DeleteOptions{})
	if err != nil {
		return fmt.Errorf("DeleteSecret: unable to delete secret %s in namespace %s: %w", name, namespace, err)
	}

	return nil
}
