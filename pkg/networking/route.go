package networking

import (
	"context"
	"fmt"

	"github.com/jaideepr97/argocd-operator-rewrite/pkg/argoutil"
	routev1 "github.com/openshift/api/route/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	ctrlClient "sigs.k8s.io/controller-runtime/pkg/client"
)

var routeAPIFound = false

// IsRouteAPIAvailable returns true if the Route API is present.
func IsRouteAPIAvailable() bool {
	return routeAPIFound
}

// VerifyRouteAPI will verify that the Route API is present.
func VerifyRouteAPI() error {
	found, err := argoutil.VerifyAPI(routev1.GroupName, routev1.SchemeGroupVersion.Version)
	if err != nil {
		return err
	}
	routeAPIFound = found
	return nil
}

func CreateRoute(route *routev1.Route, client ctrlClient.Client) error {
	if err := client.Create(context.TODO(), route); err != nil {
		return fmt.Errorf("CreateRoute: failed to create route %s in namespace %s: %w", route.Name, route.Namespace, err)
	}
	return nil
}

func ListRoutes(namespace string, client ctrlClient.Client, listOptions []ctrlClient.ListOption) (*routev1.RouteList, error) {
	existingRoutes := &routev1.RouteList{}
	err := client.List(context.TODO(), existingRoutes, listOptions...)
	if err != nil {
		return nil, fmt.Errorf("ListRoutes: unable to list routes in namespace %s: %w", namespace, err)
	}

	return existingRoutes, nil
}

func GetRoute(name, namespace string, client ctrlClient.Client) (*routev1.Route, error) {
	existingRoute := &routev1.Route{}
	err := client.Get(context.TODO(), types.NamespacedName{Name: name, Namespace: namespace}, existingRoute)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, fmt.Errorf("GetRoute: unable to find route %s in namespace %s: %w", name, namespace, err)
		}
	}

	return existingRoute, nil
}

func UpdateRoute(route *routev1.Route, client ctrlClient.Client) error {
	_, err := GetRoute(route.Name, route.Namespace, client)
	if err != nil {
		return fmt.Errorf("UpdateRoute: unable to find route %s in namespace %s: %w", route.Name, route.Namespace, err)
	}

	if err = client.Update(context.TODO(), route); err != nil {
		return fmt.Errorf("UpdateRoute: unable to update route %s in namespace %s: %w", route.Name, route.Namespace, err)
	}

	return nil
}

func DeleteRoute(name, namespace string, client ctrlClient.Client) error {
	existingRoute, err := GetRoute(name, namespace, client)
	if err != nil {
		return fmt.Errorf("DeleteRoute: unable to get route %s in namespace %s: %w", name, namespace, err)
	}

	err = client.Delete(context.TODO(), existingRoute, &ctrlClient.DeleteOptions{})
	if err != nil {
		return fmt.Errorf("DeleteRoute: unable to delete route %s in namespace %s: %w", name, namespace, err)
	}

	return nil
}
