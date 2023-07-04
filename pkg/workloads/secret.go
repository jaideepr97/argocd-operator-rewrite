package workloads

import (
	"context"

	"github.com/jaideepr97/argocd-operator-rewrite/pkg/argoutil"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrlClient "sigs.k8s.io/controller-runtime/pkg/client"
)

type SecretRequest struct {
	Name         string
	InstanceName string
	Namespace    string
	Component    string
}

func newSecret(name, instanceName, namespace, component string) *corev1.Secret {
	secretName := argoutil.GenerateResourceName(instanceName, component)
	if name != "" {
		secretName = name
	}
	return &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      secretName,
			Namespace: namespace,
			Labels:    argoutil.LabelsForCluster(instanceName, component),
		},
		Type: corev1.SecretTypeOpaque,
	}
}

func RequestSecret(request SecretRequest) *corev1.Secret {
	return newSecret(request.Name, request.InstanceName, request.Namespace, request.Component)
}

func CreateSecret(secret *corev1.Secret, client ctrlClient.Client) error {
	return client.Create(context.TODO(), secret)
}

func GetSecret(name, namespace string, client ctrlClient.Client) (*corev1.Secret, error) {
	existingSecret := &corev1.Secret{}
	err := client.Get(context.TODO(), types.NamespacedName{Name: name, Namespace: namespace}, existingSecret)
	if err != nil {
		return nil, err
	}
	return existingSecret, nil
}

func ListSecrets(namespace string, client ctrlClient.Client, listOptions ctrlClient.ListOption) (*corev1.SecretList, error) {
	existingSecrets := &corev1.SecretList{}
	err := client.List(context.TODO(), existingSecrets, listOptions)
	if err != nil {
		return nil, err
	}
	return existingSecrets, nil
}

func UpdateSecret(secret *corev1.Secret, client ctrlClient.Client) error {
	_, err := GetSecret(secret.Name, secret.Namespace, client)
	if err != nil {
		return err
	}

	if err = client.Update(context.TODO(), secret); err != nil {
		return err
	}

	return nil
}

func DeleteSecret(name, namespace string, client ctrlClient.Client) error {
	existingSecret, err := GetSecret(name, namespace, client)
	if err != nil {
		if !errors.IsNotFound(err) {
			return err
		}
		return nil
	}

	if err := client.Delete(context.TODO(), existingSecret); err != nil {
		return err
	}
	return nil
}
