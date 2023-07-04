package workloads

import (
	"context"
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	ctrlClient "sigs.k8s.io/controller-runtime/pkg/client"
)

func CreateDeployment(deployment *appsv1.Deployment, client ctrlClient.Client) error {
	// TO DO: log creation of deployment here (info)
	if err := client.Create(context.TODO(), deployment); err != nil {
		// TO DO: log error here (error)
		return fmt.Errorf("CreateDeployment: failed to create deployment %s in namespace %s: %w", deployment.Name, deployment.Namespace, err)
	}
	return nil
}

func GetDeployment(name string, namespace string, client ctrlClient.Client) (*appsv1.Deployment, error) {
	existingDeployment := &appsv1.Deployment{}
	err := client.Get(context.TODO(), types.NamespacedName{
		Namespace: namespace,
		Name:      name,
	}, existingDeployment)
	if err != nil {
		if errors.IsNotFound(err) {
			// TO DO: log deployment not found (debug)
			return nil, fmt.Errorf("GetDeployment: unable to find deployment %s: %w", name, err)
		}
		return nil, err
	}
	return existingDeployment, nil
}

func ListDeployments(client ctrlClient.Client, listOptions []ctrlClient.ListOption) (*appsv1.DeploymentList, error) {
	existingDeployments := &appsv1.DeploymentList{}
	err := client.List(context.TODO(), existingDeployments, listOptions...)
	if err != nil {
		return nil, fmt.Errorf("ListDeployments: unable to list deployments: %w", err)
	}
	return existingDeployments, nil
}

func UpdateDeployment(deployment *appsv1.Deployment, client ctrlClient.Client) error {
	_, err := GetDeployment(deployment.Name, deployment.Namespace, client)
	if err != nil {
		// TO DO: log Deployment not found (error)
		return fmt.Errorf("UpdateDeployment: unable to find Deployment %s in namespace %s: %w", deployment.Name, deployment.Namespace, err)
	}

	if err = client.Update(context.TODO(), deployment); err != nil {
		return fmt.Errorf("UpdateDeployment: unable to update Deployment %s in namespace %s: %w", deployment.Name, deployment.Namespace, err)
	}

	return nil
}

func DeleteDeployment(name, namespace string, client ctrlClient.Client) error {
	existingDeployment, err := GetDeployment(name, namespace, client)
	if err != nil {
		if !errors.IsNotFound(err) {
			return fmt.Errorf("DeleteDeployment: unable to retrieve Deployment %s in namespace %s: %w", name, namespace, err)
		}
		// TO DO: log Deployment was not found
		return nil
	}

	// TO DO: log deletion of Deployment here (info)
	if err := client.Delete(context.TODO(), existingDeployment); err != nil {
		// TO DO: log error here (warn)
		return fmt.Errorf("DeleteDeployment: failed to delete Deployment %s in namespace %s: %w", existingDeployment.Name, existingDeployment.Namespace, err)
	}
	return nil
}
