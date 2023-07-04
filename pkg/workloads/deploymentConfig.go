package workloads

import (
	"context"
	"fmt"

	oappsv1 "github.com/openshift/api/apps/v1"
	oappsv1client "github.com/openshift/client-go/apps/clientset/versioned/typed/apps/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func GetDeploymentConfig(name, namespace string, client *oappsv1client.AppsV1Client) (*oappsv1.DeploymentConfig, error) {
	existingDC, err := client.DeploymentConfigs(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			// TO DO: log dc not found (debug)
			return nil, fmt.Errorf("GetDeploymentConfig: unable to find deploymentConfig %s in namespace %s: %w", name, namespace, err)
		}
	}

	return existingDC, nil
}

// CreateDeploymentConfig creates a new DeploymentConfig in the specified namespace
func CreateDeploymentConfig(dc *oappsv1.DeploymentConfig, client *oappsv1client.AppsV1Client) error {
	_, err := client.DeploymentConfigs(dc.Namespace).Create(context.TODO(), dc, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("CreateDeploymentConfig: unable to create deploymentConfig %s in namespace %s: %w", dc.Name, dc.Namespace, err)
	}
	return nil
}

// DeleteDeploymentConfig deletes an existing DeploymentConfig in the specified namespace
func DeleteDeploymentConfig(dcName, namespace string, client *oappsv1client.AppsV1Client) error {
	err := client.DeploymentConfigs(namespace).Delete(context.TODO(), dcName, metav1.DeleteOptions{})
	if err != nil {
		return fmt.Errorf("DeleteDeploymentConfig: unable to delete deploymentConfig %s in namespace %s: %w", dcName, namespace, err)
	}
	return nil
}

// ListDeploymentConfigs lists all DeploymentConfigs in the specified namespace
func ListDeploymentConfigs(namespace string, client *oappsv1client.AppsV1Client) (*oappsv1.DeploymentConfigList, error) {
	dcList, err := client.DeploymentConfigs(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("ListDeploymentConfigs: unable to get DeploymentConfig list in namespace %s: %w", namespace, err)
	}
	return dcList, nil
}

func UpdateDeploymentConfig(dc *oappsv1.DeploymentConfig, client *oappsv1client.AppsV1Client) error {
	existingDC, err := GetDeploymentConfig(dc.Name, dc.Namespace, client)
	if err != nil {
		// TO DO: log dc not found (debug)
		return fmt.Errorf("GetDeploymentConfig: unable to find deploymentConfig %s in namespace %s: %w", dc.Name, dc.Namespace, err)
	}

	_, err = client.DeploymentConfigs(dc.Namespace).Update(context.TODO(), existingDC, metav1.UpdateOptions{})
	if err != nil {
		return fmt.Errorf("UpdateDeploymentConfig: unable to update deploymentConfig %s in namespace %s: %w", dc.Name, dc.Namespace, err)
	}

	return nil
}
