/*
Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package argocd

import (
	"context"

	"github.com/go-logr/logr"
	"github.com/jaideepr97/argocd-operator-rewrite/api/v1alpha1"
	"github.com/jaideepr97/argocd-operator-rewrite/common"
	"github.com/jaideepr97/argocd-operator-rewrite/controllers/argocd/appcontroller"
	"github.com/jaideepr97/argocd-operator-rewrite/controllers/argocd/applicationset"
	"github.com/jaideepr97/argocd-operator-rewrite/controllers/argocd/configmap"
	"github.com/jaideepr97/argocd-operator-rewrite/controllers/argocd/notifications"
	"github.com/jaideepr97/argocd-operator-rewrite/controllers/argocd/redis"
	"github.com/jaideepr97/argocd-operator-rewrite/controllers/argocd/reposerver"
	"github.com/jaideepr97/argocd-operator-rewrite/controllers/argocd/secret"
	"github.com/jaideepr97/argocd-operator-rewrite/controllers/argocd/server"
	"github.com/jaideepr97/argocd-operator-rewrite/pkg/argoutil"
	"github.com/jaideepr97/argocd-operator-rewrite/pkg/cluster"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	ctrlClient "sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// blank assignment to verify that ReconcileArgoCD implements reconcile.Reconciler
var _ reconcile.Reconciler = &ArgoCDReconciler{}

// ArgoCDReconciler reconciles a ArgoCD object
type ArgoCDReconciler struct {
	client.Client
	Scheme        *runtime.Scheme
	Instance      *v1alpha1.ArgoCD
	ClusterScoped bool
	Logger        logr.Logger

	ManagedNamespaces map[string]string
	SourceNamespaces  map[string]string

	SecretController        *secret.SecretReconciler
	ConfigMapController     *configmap.ConfigMapReconciler
	RedisController         *redis.RedisReconciler
	ReposerverController    *reposerver.RepoServerReconciler
	ServerController        *server.ServerReconciler
	NotificationsController *notifications.NotificationsReconciler
	AppController           *appcontroller.AppControllerReconciler
	AppsetController        *applicationset.ApplicationSetReconciler
}

//+kubebuilder:rbac:groups=argoproj.io,resources=argocds,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=argoproj.io,resources=argocds/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=argoproj.io,resources=argocds/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the ArgoCD object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.1/pkg/reconcile
func (r *ArgoCDReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	argocdControllerLog := ctrl.Log.WithName("argocd-controller")

	argocd := &v1alpha1.ArgoCD{}
	err := r.Client.Get(ctx, req.NamespacedName, argocd)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	r.Instance = argocd
	r.Logger = argocdControllerLog
	if r.Instance.GetDeletionTimestamp() != nil {
		// TO DO: Argo CD clean up logic here
		if argocd.IsDeletionFinalizerPresent() {
			if err := r.removeDeletionFinalizer(argocd); err != nil {
				return reconcile.Result{}, err
			}
		}
		return reconcile.Result{}, nil
	}

	if !r.Instance.IsDeletionFinalizerPresent() {
		if err := r.addDeletionFinalizer(); err != nil {
			return reconcile.Result{}, err
		}
	}

	r.ClusterScoped = AllowedNamespace(r.Instance.Namespace, GetClusterConfigNamespaces())
	r.Logger = r.Logger.WithValues("instance", r.Instance.Name, "instance-namespace", r.Instance.Namespace)

	if err = r.setManagedNamespaces(); err != nil {
		return reconcile.Result{}, err
	}

	if err = r.setSourceNamespaces(); err != nil {
		return reconcile.Result{}, err
	}

	r.InitializeControllerReconcilers()

	if err = r.reconcileControllers(); err != nil {
		return reconcile.Result{}, err
	}

	return ctrl.Result{}, nil
}

func (r *ArgoCDReconciler) reconcileControllers() error {

	// core components, return reconciliation errors
	if err := r.AppController.Reconcile(); err != nil {
		r.Logger.Error(err, "failed to reconcile application controller")
		return err
	}

	if err := r.ServerController.Reconcile(); err != nil {
		r.Logger.Error(err, "failed to reconcile server")
		return err
	}

	if err := r.SecretController.Reconcile(); err != nil {
		r.Logger.Error(err, "failed to reconcile secret controller")
		return err
	}

	// non-core components, don't return reconciliation errors
	if err := r.AppsetController.Reconcile(); err != nil {
		r.Logger.Error(err, "failed to reconcile applicationset controller")
	}

	if err := r.NotificationsController.Reconcile(); err != nil {
		r.Logger.Error(err, "failed to reconcile notifications controller")
	}

	return nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ArgoCDReconciler) SetupWithManager(mgr ctrl.Manager) error {
	bldr := ctrl.NewControllerManagedBy(mgr)
	r.setResourceWatches(bldr, clusterResourceMapper, r.tlsSecretMapper, r.namespaceMapper, r.clusterSecretMapper)
	return bldr.Complete(r)
}

// setManagedNamespaces finds all namespaces being managed by a namespace-scoped Argo CD instance
func (r *ArgoCDReconciler) setManagedNamespaces() error {
	r.ManagedNamespaces = make(map[string]string)
	listOptions := []ctrlClient.ListOption{
		client.MatchingLabels{
			common.ArgoCDResourcesManagedByLabel: r.Instance.Namespace,
		},
	}

	// get the list of namespaces managed by the Argo CD instance
	Managednamespaces, err := cluster.ListNamespaces(r.Client, listOptions)
	if err != nil {
		r.Logger.Error(err, "failed to retrieve list of managed namespaces")
		return err
	}

	r.Logger.Info("processing namespaces for resource management")

	for _, namespace := range Managednamespaces.Items {
		r.ManagedNamespaces[namespace.Name] = ""
	}

	// get control plane namespace
	_, err = cluster.GetNamespace(r.Instance.Namespace, r.Client)
	if err != nil {
		r.Logger.Error(err, "failed to retrieve control plane namespace")
		return err
	}

	// append control-plane namespace to this map
	r.ManagedNamespaces[r.Instance.Namespace] = ""
	return nil
}

// setSourceNamespaces sets a list of namespaces that a cluster-scoped Argo CD
// instance is allowed to source Applications from
func (r *ArgoCDReconciler) setSourceNamespaces() error {
	r.SourceNamespaces = make(map[string]string)
	allowedSourceNamespaces := make(map[string]string)

	if !r.ClusterScoped {
		r.Logger.V(1).Info("setSourceNamespaces: instance is not cluster scoped, skip processing namespaces for application management")
		return nil
	}

	r.Logger.Info("processing namespaces for application management")

	// Get list of existing namespaces currently carrying the ArgoCDAppsManagedBy label and convert to a map
	listOptions := []ctrlClient.ListOption{
		client.MatchingLabels{
			common.ArgoCDAppsManagedByLabel: r.Instance.Namespace,
		},
	}

	existingManagedNamesapces, err := cluster.ListNamespaces(r.Client, listOptions)
	if err != nil {
		r.Logger.Error(err, "setSourceNamespaces: failed to list namespaces")
		return err
	}
	existingManagedNsMap := make(map[string]string)
	for _, ns := range existingManagedNamesapces.Items {
		existingManagedNsMap[ns.Name] = ""
	}

	// Get list of desired namespaces that should be carrying the ArgoCDAppsManagedBy label and convert to a map
	desiredManagedNsMap := make(map[string]string)
	for _, ns := range r.Instance.Spec.SourceNamespaces {
		desiredManagedNsMap[ns] = ""
	}

	// check if any of the desired namespaces are missing the label. If yes, add ArgoCDAppsManagedByLabel to it
	for _, desiredNs := range r.Instance.Spec.SourceNamespaces {
		if _, ok := existingManagedNsMap[desiredNs]; !ok {
			ns, err := cluster.GetNamespace(desiredNs, r.Client)
			if err != nil {
				r.Logger.Error(err, "setSourceNamespaces: failed to retrieve namespace", "name", ns.Name)
				continue
			}

			// sanity check
			if len(ns.Labels) == 0 {
				ns.Labels = make(map[string]string)
			}
			// check if desired namespace is already being managed by a different cluster scoped Argo CD instance. If yes, skip it
			// If not, add ArgoCDAppsManagedByLabel to it and add it to allowedSourceNamespaces
			if val, ok := ns.Labels[common.ArgoCDAppsManagedByLabel]; ok && val != r.Instance.Namespace {
				r.Logger.V(1).Info("setSourceNamespaces: skipping namespace as it is already managed by a different instance", "namespace", ns.Name, "managing-instance-namespace", val)
				continue
			} else {
				ns.Labels[common.ArgoCDAppsManagedByLabel] = r.Instance.Namespace
				allowedSourceNamespaces[desiredNs] = ""
			}
			err = cluster.UpdateNamespace(ns, r.Client)
			if err != nil {
				r.Logger.Error(err, "setSourceNamespaces: failed to update namespace", "namespace", ns.Name)
				continue
			}
			r.Logger.V(1).Info("setSourceNamespaces: labeled namespace", "namespace", ns.Name)
			continue
		}
		allowedSourceNamespaces[desiredNs] = ""
		continue
	}

	// check if any of the exisiting namespaces are carrying the label when they should not be. If yes, remove it
	for existingNs, _ := range existingManagedNsMap {
		if _, ok := desiredManagedNsMap[existingNs]; !ok {
			ns, err := cluster.GetNamespace(existingNs, r.Client)
			if err != nil {
				r.Logger.Error(err, "setSourceNamespaces: failed to retrieve namespace", "name", ns.Name)
				continue
			}
			delete(ns.Labels, common.ArgoCDAppsManagedByLabel)
			err = cluster.UpdateNamespace(ns, r.Client)
			if err != nil {
				r.Logger.Error(err, "setSourceNamespaces: failed to update namespace", "namespace", ns.Name)
				continue
			}
			r.Logger.V(1).Info("setSourceNamespaces: unlabeled namespace", "namespace", ns.Name)
			continue
		}
	}

	r.SourceNamespaces = allowedSourceNamespaces
	return nil
}

func (r *ArgoCDReconciler) addDeletionFinalizer() error {
	r.Instance.Finalizers = append(r.Instance.Finalizers, common.ArgoCDDeletionFinalizer)
	if err := r.Client.Update(context.TODO(), r.Instance); err != nil {
		r.Logger.Error(err, "addDeletionFinalizer: failed to add deletion finalizer")
		return err
	}
	return nil
}

func (r *ArgoCDReconciler) removeDeletionFinalizer(argocd *v1alpha1.ArgoCD) error {
	argocd.Finalizers = argoutil.RemoveString(argocd.GetFinalizers(), common.ArgoCDDeletionFinalizer)
	if err := r.Client.Update(context.TODO(), argocd); err != nil {
		r.Logger.Error(err, "removeDeletionFinalizer: failed to remove deletion finalizer")
		return err
	}
	return nil
}

func (r *ArgoCDReconciler) InitializeControllerReconcilers() {
	r.SecretController = &secret.SecretReconciler{
		Client:            &r.Client,
		Scheme:            r.Scheme,
		Instance:          r.Instance,
		ClusterScoped:     r.ClusterScoped,
		ManagedNamespaces: r.ManagedNamespaces,
	}

	r.ConfigMapController = &configmap.ConfigMapReconciler{
		Client:   &r.Client,
		Scheme:   r.Scheme,
		Instance: r.Instance,
	}

	r.RedisController = &redis.RedisReconciler{
		Client:   &r.Client,
		Scheme:   r.Scheme,
		Instance: r.Instance,
	}

	r.ReposerverController = &reposerver.RepoServerReconciler{
		Client:   &r.Client,
		Scheme:   r.Scheme,
		Instance: r.Instance,
	}

	r.ServerController = &server.ServerReconciler{
		Client:            &r.Client,
		Scheme:            r.Scheme,
		Instance:          r.Instance,
		ClusterScoped:     r.ClusterScoped,
		ManagedNamespaces: r.ManagedNamespaces,
		SourceNamespaces:  r.SourceNamespaces,
	}

	r.NotificationsController = &notifications.NotificationsReconciler{
		Client:   &r.Client,
		Scheme:   r.Scheme,
		Instance: r.Instance,
	}

	r.AppController = &appcontroller.AppControllerReconciler{
		Client:            &r.Client,
		Scheme:            r.Scheme,
		Instance:          r.Instance,
		ClusterScoped:     r.ClusterScoped,
		ManagedNamespaces: r.ManagedNamespaces,
		SourceNamespaces:  r.SourceNamespaces,
	}

	r.AppsetController = &applicationset.ApplicationSetReconciler{
		Client:   &r.Client,
		Scheme:   r.Scheme,
		Instance: r.Instance,
	}
}
