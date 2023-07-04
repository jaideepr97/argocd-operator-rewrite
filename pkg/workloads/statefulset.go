package workloads

import (
	"context"
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	ctrlClient "sigs.k8s.io/controller-runtime/pkg/client"
)

func CreateStatefulSet(client ctrlClient.Client, statefulSet *appsv1.StatefulSet) error {
	if err := client.Create(context.Background(), statefulSet); err != nil {
		return fmt.Errorf("CreateStatefulSet: failed to create StatefulSet %s in namespace %s: %w", statefulSet.Name, statefulSet.Namespace, err)
	}
	return nil
}

func GetStatefulSet(client ctrlClient.Client, name, namespace string) (*appsv1.StatefulSet, error) {
	existingStatefulSet := &appsv1.StatefulSet{}
	err := client.Get(context.TODO(), ctrlClient.ObjectKey{
		Namespace: namespace,
		Name:      name,
	}, existingStatefulSet)
	if err != nil {
		return nil, fmt.Errorf("GetStatefulSet: unable to get StatefulSet %s in namespace %s: %w", name, namespace, err)
	}
	return existingStatefulSet, nil
}

func ListStatefulSets(client ctrlClient.Client, listOptions ctrlClient.ListOption) (*appsv1.StatefulSetList, error) {
	existingStatefulSets := &appsv1.StatefulSetList{}
	err := client.List(context.Background(), existingStatefulSets, listOptions)
	if err != nil {
		return nil, fmt.Errorf("ListStatefulSet: unable to list StatefulSets: %w", err)
	}
	return existingStatefulSets, nil
}

func UpdateStatefulSet(client ctrlClient.Client, statefulSet *appsv1.StatefulSet) error {
	_, err := GetStatefulSet(client, statefulSet.Name, statefulSet.Namespace)
	if err != nil {
		return fmt.Errorf("UpdateStatefulSet: unable to find StatefulSet %s in namespace %s: %w",
			statefulSet.Name, statefulSet.Namespace, err)
	}

	if err = client.Update(context.Background(), statefulSet); err != nil {
		return fmt.Errorf("UpdateStatefulSet: unable to update StatefulSet %s in namespace %s: %w",
			statefulSet.Name, statefulSet.Namespace, err)
	}

	return nil
}

func DeleteStatefulSet(client ctrlClient.Client, name, namespace string) error {
	existingStatefulSet, err := GetStatefulSet(client, name, namespace)
	if err != nil {
		if !errors.IsNotFound(err) {
			return fmt.Errorf("DeleteStatefulSet: unable to retrieve StatefulSet %s in namespace %s: %w",
				name, namespace, err)
		}
		return nil
	}

	if err := client.Delete(context.Background(), existingStatefulSet); err != nil {
		return fmt.Errorf("DeleteStatefulSet: failed to delete StatefulSet %s in namespace %s: %w",
			existingStatefulSet.Name, existingStatefulSet.Namespace, err)
	}

	return nil
}
