package workloads

import (
	"context"
	"fmt"

	"github.com/jaideepr97/argocd-operator-rewrite/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrlClient "sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// CreateConfigMap creates the given ConfigMap in the specified namespace
func CreateConfigMap(cr *v1alpha1.ArgoCD, cm *corev1.ConfigMap, scheme *runtime.Scheme, client ctrlClient.Client) error {
	if err := controllerutil.SetControllerReference(cr, cm, scheme); err != nil {
		// TO DO: log error here (error)
		return fmt.Errorf("CreateConfigMap: failed to set instance %s in namespace %s as owner of ConfigMap %s: %w", cr.Name, cr.Namespace, cm.Name, err)
	}

	// TO DO: log creation of ConfigMap here (info)
	if err := client.Create(context.TODO(), cm); err != nil {
		// TO DO: log error here (error)
		return fmt.Errorf("CreateConfigMap: failed to create ConfigMap %s in namespace %s: %w", cm.Name, cm.Namespace, err)
	}
	return nil
}

// ListConfigMaps lists the ConfigMaps in the specified namespace that match the specified options
func ListConfigMaps(namespace string, client ctrlClient.Client, listOptions []ctrlClient.ListOption) (*corev1.ConfigMapList, error) {
	existingConfigMaps := &corev1.ConfigMapList{}
	err := client.List(context.TODO(), existingConfigMaps, listOptions...)
	if err != nil {
		return nil, fmt.Errorf("ListConfigMaps: unable to list ConfigMaps in namespace %s: %w", namespace, err)
	}

	return existingConfigMaps, nil
}

// GetConfigMap returns the specified ConfigMap in the specified namespace
func GetConfigMap(name, namespace string, client ctrlClient.Client) (*corev1.ConfigMap, error) {
	existingConfigMap := &corev1.ConfigMap{}
	err := client.Get(context.TODO(), types.NamespacedName{Name: name, Namespace: namespace}, existingConfigMap)
	if err != nil {
		if errors.IsNotFound(err) {
			// TO DO: log ConfigMap not found (debug)
			return nil, fmt.Errorf("GetConfigMap: unable to find ConfigMap %s in namespace %s: %w", name, namespace, err)
		}
	}

	return existingConfigMap, nil
}

// UpdateConfigMap updates the specified ConfigMap in the specified namespace
func UpdateConfigMap(cm *corev1.ConfigMap, client ctrlClient.Client) error {
	_, err := GetConfigMap(cm.Name, cm.Namespace, client)
	if err != nil {
		// TO DO: log ConfigMap not found (error)
		return fmt.Errorf("UpdateConfigMap: unable to find ConfigMap %s in namespace %s: %w", cm.Name, cm.Namespace, err)
	}

	if err = client.Update(context.TODO(), cm); err != nil {
		return fmt.Errorf("UpdateConfigMap: unable to update ConfigMap %s in namespace %s: %w", cm.Name, cm.Namespace, err)
	}

	return nil
}

// DeleteConfigMap deletes the specified ConfigMap in the specified namespace
func DeleteConfigMap(name, namespace string, client ctrlClient.Client) error {
	existingConfigMap, err := GetConfigMap(name, namespace, client)
	if err != nil {
		return fmt.Errorf("DeleteConfigMap: unable to get ConfigMap %s in namespace %s: %w", name, namespace, err)
	}

	err = client.Delete(context.TODO(), existingConfigMap, &ctrlClient.DeleteOptions{})
	if err != nil {
		return fmt.Errorf("DeleteConfigMap: unable to delete ConfigMap %s in namespace %s: %w", name, namespace, err)
	}

	return nil
}
