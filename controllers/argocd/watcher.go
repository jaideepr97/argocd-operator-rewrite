package argocd

import (
	"github.com/jaideepr97/argocd-operator-rewrite/api/v1alpha1"
	"github.com/jaideepr97/argocd-operator-rewrite/common"
	"github.com/jaideepr97/argocd-operator-rewrite/pkg/monitoring"
	"github.com/jaideepr97/argocd-operator-rewrite/pkg/networking"
	"github.com/jaideepr97/argocd-operator-rewrite/pkg/workloads"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	oappsv1 "github.com/openshift/api/apps/v1"
	routev1 "github.com/openshift/api/route/v1"
	monitoringv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// setResourceWatches will register Watches for each of the managed Resources.
func (r *ArgoCDReconciler) setResourceWatches(bldr *builder.Builder, clusterResourceMapper, tlsSecretMapper, namespaceMapper, clusterSecretMapper handler.MapFunc) *builder.Builder {

	clusterResourceHandler := handler.EnqueueRequestsFromMapFunc(clusterResourceMapper)
	clusterSecretResourceHandler := handler.EnqueueRequestsFromMapFunc(clusterSecretMapper)
	tlsSecretHandler := handler.EnqueueRequestsFromMapFunc(tlsSecretMapper)
	namespaceHandler := handler.EnqueueRequestsFromMapFunc(namespaceMapper)

	// Watch for changes to primary resource Argo CD
	// bldr.For(&v1alpha1.ArgoCD{}, builder.WithPredicates(ignoreStatusAndDeletionPredicate()))
	bldr.For(&v1alpha1.ArgoCD{})

	// Watch for changes to ConfigMap sub-resources owned by Argo CD instances.
	bldr.Owns(&corev1.ConfigMap{})

	// Watch for changes to Deployment sub-resources owned by Argo CD instances.
	bldr.Owns(&appsv1.Deployment{})

	// Watch for changes to StatefulSet sub-resources owned by Argo CD instances.
	bldr.Owns(&appsv1.StatefulSet{})

	// Watch for changes to Secret sub-resources owned by Argo CD instances.
	bldr.Owns(&corev1.Secret{})

	// Watch for changes to secrets of type TLS that might be created by external processes
	bldr.Watches(&corev1.Secret{Type: corev1.SecretTypeTLS}, tlsSecretHandler)

	// Watch for changes to cluster secrets added to the argocd instance
	bldr.Watches(&corev1.Secret{ObjectMeta: metav1.ObjectMeta{
		Labels: map[string]string{
			common.ArgoCDSecretTypeLabel: "cluster",
		}}}, clusterSecretResourceHandler)

	// Watch for changes to Role sub-resources owned by Argo CD instances.
	bldr.Owns(&rbacv1.Role{})

	// Watch for changes to RoleBinding sub-resources owned by Argo CD instances.
	bldr.Owns(&rbacv1.RoleBinding{})

	// Watch for changes to clusterRole sub-resources owned by Argo CD instances.
	bldr.Watches(&rbacv1.ClusterRoleBinding{}, clusterResourceHandler)

	// Watch for changes to clusterRoleBinding sub-resources owned by Argo CD instances.
	bldr.Watches(&rbacv1.ClusterRole{}, clusterResourceHandler)

	// Watch for changes to Service sub-resources owned by Argo CD instances.
	bldr.Owns(&corev1.Service{})

	// Watch for changes to Ingress sub-resources owned by Argo CD instances.
	bldr.Owns(&networkingv1.Ingress{})

	// Watch for changes to namespaces managed by Argo CD instances.
	bldr.Watches(&corev1.Namespace{}, namespaceHandler, builder.WithPredicates(namespaceFilterPredicate()))

	// Inspect cluster to verify availability of extra features
	// This sets the flags that are used in subsequent checks
	InspectCluster()

	if networking.IsRouteAPIAvailable() {
		// Watch OpenShift Route sub-resources owned by ArgoCD instances.
		bldr.Owns(&routev1.Route{})
	}

	if monitoring.IsPrometheusAPIAvailable() {
		// Watch PrometheusRule sub-resources owned by ArgoCD instances.
		bldr.Owns(&monitoringv1.PrometheusRule{})

		// Watch Prometheus ServiceMonitor sub-resources owned by ArgoCD instances.
		bldr.Owns(&monitoringv1.ServiceMonitor{})
	}

	if workloads.IsTemplateAPIAvailable() {
		// Watch for the changes to Deployment Config
		bldr.Owns(&oappsv1.DeploymentConfig{})
	}

	// bldr.WithEventFilter(ignoreStatusAndDeletionPredicate())

	return bldr
}

// namespaceFilterPredicate defines how we filter events on namespaces to decide if a new round of reconciliation should be triggered or not.
func namespaceFilterPredicate() predicate.Predicate {
	return predicate.Funcs{
		UpdateFunc: func(e event.UpdateEvent) bool {
			triggerReconciliation := false

			InitializeScheduledForRBACDeletion()

			// This checks if either ArgoCDManaged lables exist in newMeta, if yes then trigger reconciliation - but first
			// 1. Check if oldMeta had the label or not. If not ==> new managed namespace! Return true to schedule reconciliation
			// 2. if yes, check if the old and new values are different, if yes ==> namespace is now managed by a different instance.
			// we must delete old rbac resources from this namespace, and potentially remove it from the cluster secret
			// Event is then handled by the reconciler, which would create new RBAC resources appropriately

			// handle resource management label
			if valNew, ok := e.ObjectNew.GetLabels()[common.ArgoCDResourcesManagedByLabel]; ok {
				if valOld, ok := e.ObjectOld.GetLabels()[common.ArgoCDResourcesManagedByLabel]; ok && valOld != valNew {
					managedNamespace := e.ObjectOld.GetName()
					previouslyManagingNamespace := valOld

					// create ManagedNsOpts in rbacDeletionMap if not present already. Add previouslyManagingNs and resource management label to the managedNs options so that resource management role can be deleted from this namespace, and previously managing namespace can be removed from cluster secret
					managedNsOpts, ok := ScheduledForRBACDeletion[managedNamespace]
					if !ok {
						ScheduledForRBACDeletion[managedNamespace] = &ManagedNamespaceOptions{
							ManagingNamespace:           previouslyManagingNamespace,
							ResourceDeletionLabelValues: []string{common.ArgoCDRBACTypeResourceMananagement},
						}
					} else {
						managedNsOpts.ManagingNamespace = previouslyManagingNamespace
						managedNsOpts.ResourceDeletionLabelValues = append(managedNsOpts.ResourceDeletionLabelValues, common.ArgoCDRBACTypeResourceMananagement)
					}
				}
				triggerReconciliation = true
			} else {
				// newMeta does not carry ArgoCDResourcesManagedByLabel label, so check if oldMeta did carry it. If yes, label was removed. Add this namespace to scheduledForRBACDeletion
				// and add namespace in oldMeta as the value

				if valOld, ok := e.ObjectOld.GetLabels()[common.ArgoCDResourcesManagedByLabel]; ok && valOld != "" {
					managedNamespace := e.ObjectOld.GetName()
					previouslyManagingNamespace := valOld

					// create ManagedNsOpts in rbacDeletionMap if not present already. Add previouslyManagingNs and resource management annotation to the managedNs options so that resource management role can be deleted from this namespace, and previously managing namespace can be removed from cluster secret
					managedNsOpts, ok := ScheduledForRBACDeletion[managedNamespace]
					if !ok {
						ScheduledForRBACDeletion[managedNamespace] = &ManagedNamespaceOptions{
							ManagingNamespace:           previouslyManagingNamespace,
							ResourceDeletionLabelValues: []string{common.ArgoCDRBACTypeResourceMananagement},
						}
					} else {
						managedNsOpts.ManagingNamespace = previouslyManagingNamespace
						managedNsOpts.ResourceDeletionLabelValues = append(managedNsOpts.ResourceDeletionLabelValues, common.ArgoCDRBACTypeResourceMananagement)
					}

					triggerReconciliation = true
				}
			}

			// handle application management label
			if valNew, ok := e.ObjectNew.GetLabels()[common.ArgoCDAppsManagedByLabel]; ok {
				if valOld, ok := e.ObjectOld.GetLabels()[common.ArgoCDAppsManagedByLabel]; ok && valOld != valNew {
					managedNamespace := e.ObjectOld.GetName()

					// create ManagedNsOpts in rbacDeletionMap if not present already. add app management annotation to the managedNs options so that app management role can be deleted from this namespace
					managedNsOpts, ok := ScheduledForRBACDeletion[managedNamespace]
					if !ok {
						ScheduledForRBACDeletion[managedNamespace] = &ManagedNamespaceOptions{
							ResourceDeletionLabelValues: []string{common.ArgoCDRBACTypeAppManagement},
						}
					} else {
						managedNsOpts.ResourceDeletionLabelValues = append(managedNsOpts.ResourceDeletionLabelValues, common.ArgoCDRBACTypeAppManagement)
					}
				}

				triggerReconciliation = true
			} else {
				// newMeta does not carry ArgoCDAppsManagedByLabel label, so check if oldMeta did carry it. If yes, label was removed. Add this namespace to scheduledForRBACDeletion

				if valOld, ok := e.ObjectOld.GetLabels()[common.ArgoCDAppsManagedByLabel]; ok && valOld != "" {
					managedNamespace := e.ObjectOld.GetName()

					// create ManagedNsOpts in rbacDeletionMap if not present already. add app management annotation to the managedNs options so that app management role can be deleted from this namespace
					managedNsOpts, ok := ScheduledForRBACDeletion[managedNamespace]
					if !ok {
						ScheduledForRBACDeletion[managedNamespace] = &ManagedNamespaceOptions{
							ResourceDeletionLabelValues: []string{common.ArgoCDRBACTypeAppManagement},
						}
					} else {
						managedNsOpts.ResourceDeletionLabelValues = append(managedNsOpts.ResourceDeletionLabelValues, common.ArgoCDRBACTypeAppManagement)
					}

					triggerReconciliation = true
				}
			}
			return triggerReconciliation
		},
		DeleteFunc: func(e event.DeleteEvent) bool {
			triggerReconciliation := false
			// If namespace is scheduled for deletion, check if it carries either ArgoCDManaged lables. If yes, add this namespace to scheduledForRBACDeletion and add namespace in meta as the value (if required)

			// handle resource management label
			if valOld, ok := e.Object.GetLabels()[common.ArgoCDResourcesManagedByLabel]; ok && valOld != "" {
				managedNamespace := e.Object.GetName()
				previouslyManagingNamespace := valOld

				// create ManagedNsOpts in rbacDeletionMap if not present already. Add previouslyManagingNs and resource management annotation to the managedNs options so that resource management role can be deleted from this namespace, and previously managing namespace can be removed from cluster secret
				managedNsOpts, ok := ScheduledForRBACDeletion[managedNamespace]
				if !ok {
					ScheduledForRBACDeletion[managedNamespace] = &ManagedNamespaceOptions{
						ManagingNamespace:           previouslyManagingNamespace,
						ResourceDeletionLabelValues: []string{common.ArgoCDRBACTypeResourceMananagement},
					}
				} else {
					managedNsOpts.ManagingNamespace = previouslyManagingNamespace
					managedNsOpts.ResourceDeletionLabelValues = append(managedNsOpts.ResourceDeletionLabelValues, common.ArgoCDRBACTypeResourceMananagement)
				}

				triggerReconciliation = true
			}

			// handle application management label
			if valOld, ok := e.Object.GetLabels()[common.ArgoCDAppsManagedByLabel]; ok && valOld != "" {
				managedNamespace := e.Object.GetName()

				// create ManagedNsOpts in rbacDeletionMap if not present already. add app management annotation to the managedNs options so that app management role can be deleted from this namespace
				managedNsOpts, ok := ScheduledForRBACDeletion[managedNamespace]
				if !ok {
					ScheduledForRBACDeletion[managedNamespace] = &ManagedNamespaceOptions{
						ResourceDeletionLabelValues: []string{common.ArgoCDRBACTypeAppManagement},
					}
				} else {
					managedNsOpts.ResourceDeletionLabelValues = append(managedNsOpts.ResourceDeletionLabelValues, common.ArgoCDRBACTypeAppManagement)
				}

				triggerReconciliation = true
			}
			return triggerReconciliation
		},
	}
}

func ignoreStatusAndDeletionPredicate() predicate.Predicate {
	return predicate.Funcs{
		UpdateFunc: func(e event.UpdateEvent) bool {
			triggerReconciliation := false
			// Ignore updates to status in which case metadata.Generation does not change
			if e.ObjectOld.GetGeneration() != e.ObjectNew.GetGeneration() {
				triggerReconciliation = true
			}
			return triggerReconciliation
		},
		DeleteFunc: func(e event.DeleteEvent) bool {
			triggerReconciliation := false
			// Evaluates to false if the object has been confirmed deleted.
			if e.DeleteStateUnknown {
				triggerReconciliation = true
			}
			return triggerReconciliation
		},
	}
}
