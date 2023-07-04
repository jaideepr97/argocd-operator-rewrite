package argocd

import (
	"context"
	"strings"

	"github.com/jaideepr97/argocd-operator-rewrite/api/v1alpha1"
	"github.com/jaideepr97/argocd-operator-rewrite/common"
	"github.com/jaideepr97/argocd-operator-rewrite/controllers/argocd/appcontroller"
	"github.com/jaideepr97/argocd-operator-rewrite/controllers/argocd/redis"
	"github.com/jaideepr97/argocd-operator-rewrite/controllers/argocd/reposerver"
	"github.com/jaideepr97/argocd-operator-rewrite/controllers/argocd/server"
	"github.com/jaideepr97/argocd-operator-rewrite/pkg/permissions"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrlClient "sigs.k8s.io/controller-runtime/pkg/client"
)

// ManagedNamespaceOptions tracks the following for a given namespace that is no longer managed:
// 1. The namespace of the previously managing Argo CD instance
// 2. The label value for type of rbac resources to be deleted from the managed namespace: resource-management/app-management/both
type ManagedNamespaceOptions struct {
	ManagingNamespace           string
	ResourceDeletionLabelValues []string
}

// scheduledForRBACDeletion tracks namespaces that need to have roles/rolebindings deleted from them
// when they are no longer being managed by a given Argo CD instance. These namespaces also need to be
// removed from the cluster secret of the previously managing Argo CD instance.
// The key is the managed namespace, and the value is the namespace options associated with this namespace
var ScheduledForRBACDeletion map[string]*ManagedNamespaceOptions

// clusterResourceMapper maps a watch event on any cluster scoped object carrying apprpriate annotations
// back to the ArgoCD object we want to reconcile
func clusterResourceMapper(ctx context.Context, o client.Object) []reconcile.Request {
	clusterResourceAnnotations := o.GetAnnotations()
	namespacedArgoCDObject := client.ObjectKey{}

	if name, ok := clusterResourceAnnotations[common.AnnotationName]; ok {
		namespacedArgoCDObject.Name = name
	}

	if namespace, ok := clusterResourceAnnotations[common.AnnotationNamespace]; ok {
		namespacedArgoCDObject.Namespace = namespace
	}

	var result = []reconcile.Request{}
	if namespacedArgoCDObject.Name != "" && namespacedArgoCDObject.Namespace != "" {
		result = []reconcile.Request{
			{NamespacedName: namespacedArgoCDObject},
		}
	}
	return result
}

// tlsSecretMapper maps a watch event on a secret of type TLS back to the
// ArgoCD object that we want to reconcile.
func (r *ArgoCDReconciler) tlsSecretMapper(ctx context.Context, o client.Object) []reconcile.Request {
	var result = []reconcile.Request{}

	if !isSecretOfInterest(o) {
		return result
	}
	namespacedArgoCDObject := client.ObjectKey{}

	secretOwnerRefs := o.GetOwnerReferences()
	if len(secretOwnerRefs) > 0 {
		// OpenShift service CA makes the owner reference for the TLS secret to the
		// service, which in turn is owned by the controller. This method performs
		// a lookup of the controller through the intermediate owning service.
		for _, secretOwner := range secretOwnerRefs {
			if isOwnerOfInterest(secretOwner) {
				key := client.ObjectKey{Name: secretOwner.Name, Namespace: o.GetNamespace()}
				svc := &corev1.Service{}

				// Get the owning object of the secret
				err := r.Client.Get(context.TODO(), key, svc)
				if err != nil {
					// TO DO: log error here
					return result
				}

				// If there's an object of kind ArgoCD in the owner's list,
				// this will be our reconciled object.
				serviceOwnerRefs := svc.GetOwnerReferences()
				for _, serviceOwner := range serviceOwnerRefs {
					if serviceOwner.Kind == "ArgoCD" {
						namespacedArgoCDObject.Name = serviceOwner.Name
						namespacedArgoCDObject.Namespace = svc.ObjectMeta.Namespace
						result = []reconcile.Request{
							{NamespacedName: namespacedArgoCDObject},
						}
						return result
					}
				}
			}
		}
	} else {
		// For secrets without owner (i.e. manually created), we apply some
		// heuristics. This may not be as accurate (e.g. if the user made a
		// typo in the resource's name), but should be good enough for now.
		secret, ok := o.(*corev1.Secret)
		if !ok {
			return result
		}
		if owner, ok := secret.Annotations[common.AnnotationName]; ok {
			namespacedArgoCDObject.Name = owner
			namespacedArgoCDObject.Namespace = o.GetNamespace()
			result = []reconcile.Request{
				{NamespacedName: namespacedArgoCDObject},
			}
		}
	}

	return result
}

// clusterSecretMapper maps a watch event on a secret with cluster label, back to the
// ArgoCD object that we want to reconcile.
func (r *ArgoCDReconciler) clusterSecretMapper(ctx context.Context, o client.Object) []reconcile.Request {
	var result = []reconcile.Request{}

	labels := o.GetLabels()
	if v, ok := labels[common.ArgoCDSecretTypeLabel]; ok && v == "cluster" {
		argocds := &v1alpha1.ArgoCDList{}
		if err := r.Client.List(context.TODO(), argocds, &client.ListOptions{Namespace: o.GetNamespace()}); err != nil {
			return result
		}

		if len(argocds.Items) != 1 {
			return result
		}

		argocd := argocds.Items[0]
		namespacedName := client.ObjectKey{
			Name:      argocd.Name,
			Namespace: argocd.Namespace,
		}
		result = []reconcile.Request{
			{NamespacedName: namespacedName},
		}
	}

	return result
}

// namespaceMapper maps a watch event on a namespace, back to the
// ArgoCD object that we want to reconcile.
func (r *ArgoCDReconciler) namespaceMapper(ctx context.Context, obj client.Object) []reconcile.Request {
	var result = []reconcile.Request{}

	managedNamespace := obj.GetName()
	if managedNsOptions, ok := ScheduledForRBACDeletion[managedNamespace]; ok {

		err := r.deleteRBACFromPreviouslyManagedNamespace(managedNamespace, managedNsOptions.ResourceDeletionLabelValues)
		if err != nil {
			r.Logger.Error(err, "namespaceMapper: failed to delete resources from previously managed namespace: %s", managedNamespace)
		}

		if managedNsOptions.ManagingNamespace != "" {
			// This means namespace was managed for resources, and cluster secret in managing namespace needs to be updated
			// Delegate handling of cluster secret update to the secret controller
			err = r.SecretController.DeleteManagedNamespaceFromClusterSecret(managedNsOptions.ManagingNamespace, managedNamespace)
			if err != nil {
				r.Logger.Error(err, "namespaceMapper: failed to delete previously managed namespace %s from cluster secret in namespace %s", managedNamespace, managedNsOptions.ManagingNamespace)
			}
		}

		delete(ScheduledForRBACDeletion, managedNamespace)
	}

	if ns, ok := obj.GetLabels()[common.ArgoCDResourcesManagedByLabel]; ok {
		argocds := &v1alpha1.ArgoCDList{}
		if err := r.Client.List(context.TODO(), argocds, &client.ListOptions{Namespace: ns}); err != nil {
			return result
		}

		if len(argocds.Items) != 1 {
			return result
		}

		argocd := argocds.Items[0]
		namespacedName := client.ObjectKey{
			Name:      argocd.Name,
			Namespace: argocd.Namespace,
		}
		result = []reconcile.Request{
			{NamespacedName: namespacedName},
		}
	}

	return result
}

// isSecretOfInterest returns true if the name of the given secret matches one of the
// well-known tls secrets used to secure communication amongst the Argo CD components.
func isSecretOfInterest(o client.Object) bool {
	if strings.HasSuffix(o.GetName(), reposerver.ArgoCDRepoServerTLSSuffix) {
		return true
	}
	if o.GetName() == redis.ArgoCDRedisServerTLSSecretName {
		return true
	}
	return false
}

// isOwnerOfInterest returns true if the given owner is one of the Argo CD services that
// may have been made the owner of the tls secret created by the OpenShift service CA, used
// to secure communication amongst the Argo CD components.
func isOwnerOfInterest(owner metav1.OwnerReference) bool {
	if owner.Kind != "Service" {
		return false
	}
	if strings.HasSuffix(owner.Name, reposerver.ArgoCDRepoServerSuffix) {
		return true
	}
	if strings.HasSuffix(owner.Name, redis.ArgoCDDefaultRedisSuffix) {
		return true
	}
	return false
}

// deleteRBACFromPreviouslyManagedCluster deletes roles and rolebindings from a namespace that is either
// 1. no longer managed by an Argo CD instance or
// 2. managed by a different Argo CD instance
// 3. About to be deleted
// This includes roles for both app-controller and server, for resource management as well as app management
// In any case, we must delete all roles and rolebindings from this namespace at this point.
// Any other roles and rolebindings that need to be created will be handled by the reconciler
func (r *ArgoCDReconciler) deleteRBACFromPreviouslyManagedNamespace(namespace string, ResourceDeletionLabelValues []string) error {

	listOptions := getResourceSelectionOptions(ResourceDeletionLabelValues)
	listOptions = append(listOptions, &ctrlClient.ListOptions{
		Namespace: namespace,
	})

	// List all the roles created for ArgoCD using listOptions, and request to have them deleted
	existingRoles, err := permissions.ListRoles(namespace, r.Client, listOptions)
	if err != nil {
		r.Logger.Error(err, "deleteRBACFromPreviouslyManagedCluster: failed to list roles", "namespace", namespace)
		return err
	}

	// List all the roleBindings created for ArgoCD using the label selector, and request to have them deleted
	existingRoleBindings, err := permissions.ListRoleBindings(namespace, r.Client, listOptions)
	if err != nil {
		r.Logger.Error(err, "deleteRBACFromPreviouslyManagedCluster: failed to list roleBindings", "namespace", namespace)
		return err
	}

	for _, rb := range existingRoleBindings.Items {
		// delegate deletion to appController roleBindings to appController
		if strings.HasSuffix(rb.Name, appcontroller.ArgoCDApplicationControllerComponent) {
			if err := r.AppController.DeleteRoleBinding(rb.Name, rb.Namespace); err != nil {
				r.Logger.Error(err, "deleteRBACFromPreviouslyManagedCluster: failed to delete app controller roleBinding", "name", rb.Name, "namespace", namespace)
			}

		}

		// delegate deletion to server roleBindings to server
		if strings.HasSuffix(rb.Name, server.ArgoCDServerComponent) {
			if err := r.ServerController.DeleteRoleBinding(rb.Name, rb.Namespace); err != nil {
				r.Logger.Error(err, "deleteRBACFromPreviouslyManagedCluster: failed to delete server roleBinding", "name", rb.Name, "namespace", namespace)
			}
		}
	}

	for _, role := range existingRoles.Items {
		// delegate deletion to appController roles to appController
		if strings.HasSuffix(role.Name, appcontroller.ArgoCDApplicationControllerComponent) {
			if err := r.AppController.DeleteRole(role.Name, role.Namespace); err != nil {
				r.Logger.Error(err, "deleteRBACFromPreviouslyManagedCluster: failed to delete app controller role", "name", role.Name, "namespace", namespace)
			}

		}

		// delegate deletion to server roles to server
		if strings.HasSuffix(role.Name, server.ArgoCDServerComponent) {
			if err := r.ServerController.DeleteRole(role.Name, role.Namespace); err != nil {
				r.Logger.Error(err, "deleteRBACFromPreviouslyManagedCluster: failed to delete server role", "name", role.Name, "namespace", namespace)
			}
		}
	}

	return nil
}

// getResourceSelectionOptions sets the rbac-type label in the listOptions based on which managed by label is involved
func getResourceSelectionOptions(ResourceDeletionLabelValues []string) []ctrlClient.ListOption {
	var listOptions []ctrlClient.ListOption = make([]ctrlClient.ListOption, 0)
	for _, val := range ResourceDeletionLabelValues {
		listOptions = append(listOptions, ctrlClient.MatchingLabels{
			common.ArgoCDKeyRBACType: val,
		})
	}
	return listOptions
}

func InitializeScheduledForRBACDeletion() {
	if len(ScheduledForRBACDeletion) == 0 {
		ScheduledForRBACDeletion = make(map[string]*ManagedNamespaceOptions)
	}
}
