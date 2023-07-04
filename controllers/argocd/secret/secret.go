package secret

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/go-logr/logr"
	"github.com/jaideepr97/argocd-operator-rewrite/api/v1alpha1"
	"github.com/jaideepr97/argocd-operator-rewrite/common"
	"github.com/jaideepr97/argocd-operator-rewrite/pkg/argoutil"
	"github.com/jaideepr97/argocd-operator-rewrite/pkg/workloads"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	ctrlClient "sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

type SecretReconciler struct {
	Client            *client.Client
	Scheme            *runtime.Scheme
	Instance          *v1alpha1.ArgoCD
	ClusterScoped     bool
	Logger            logr.Logger
	ManagedNamespaces map[string]string
}

func (sr *SecretReconciler) Reconcile() error {
	sr.Logger = ctrl.Log.WithName(ArgoCDSecretControllerComponent).WithValues("instance", sr.Instance.Name, "instance-namespace", sr.Instance.Namespace)

	if err := sr.reconcileClusterPermissionsSecret(); err != nil {
		sr.Logger.Error(err, "failed to reconcile cluster permissions secret")
		return err
	}

	return nil
}

func (sr *SecretReconciler) reconcileClusterPermissionsSecret() error {
	secretName := argoutil.NameWithSuffix(sr.Instance.Name, DefaultClusterConfigSecretSuffix)

	secretRequest := workloads.SecretRequest{
		Name:         secretName,
		InstanceName: sr.Instance.Name,
		Namespace:    sr.Instance.Namespace,
	}

	desiredSecret := workloads.RequestSecret(secretRequest)
	desiredSecret.Labels[common.ArgoCDSecretTypeLabel] = common.ArgoCDSecretTypeCluster

	dataBytes, err := json.Marshal(map[string]interface{}{
		"tlsClientConfig": map[string]interface{}{
			"insecure": false,
		},
	})
	if err != nil {
		sr.Logger.Error(err, "reconcileClusterPermissionsSecret: failed to marshal json")
		return err
	}

	desiredManagedNamespaces := make([]string, 0)
	for ns, _ := range sr.ManagedNamespaces {
		desiredManagedNamespaces = append(desiredManagedNamespaces, ns)
	}
	sort.Strings(desiredManagedNamespaces)

	desiredSecret.Data = map[string][]byte{
		"config": dataBytes,
		"name":   []byte("in-cluster"),
		"server": []byte(common.ArgoCDDefaultServer),
	}

	if !sr.ClusterScoped {
		desiredSecret.Data["namespaces"] = []byte(strings.Join(desiredManagedNamespaces, ","))
	}

	existingSecret, err := workloads.GetSecret(secretName, sr.Instance.Namespace, *sr.Client)
	if err != nil {
		if !errors.IsNotFound(err) {
			sr.Logger.Error(err, "reconcileClusterPermissionsSecret: failed to retrieve secret", "name", secretName)
			return err
		}

		if err = controllerutil.SetControllerReference(sr.Instance, desiredSecret, sr.Scheme); err != nil {
			sr.Logger.Error(err, "reconcileClusterPermissionsSecret: failed to set owner reference for secret", "name", desiredSecret.Name, "namespace", desiredSecret.Namespace)
		}

		if err = workloads.CreateSecret(desiredSecret, *sr.Client); err != nil {
			sr.Logger.Error(err, "reconcileClusterPermissionsSecret: failed to create secret", "name", desiredSecret.Name, "namespace", desiredSecret.Namespace)
			return err
		}
		sr.Logger.V(0).Info("reconcileClusterPermissionsSecret: secret created")
		return nil
	}

	secretChanged := false

	if sr.ClusterScoped {
		if _, ok := existingSecret.Data["namespaces"]; ok {
			delete(existingSecret.Data, "namespaces")
			secretChanged = true
		}
	} else {
		existingManagedNamespaces := argoutil.SplitList(string(existingSecret.Data["namespaces"]))
		if !argoutil.Equal(existingManagedNamespaces, desiredManagedNamespaces) {
			existingSecret.Data["namespaces"] = desiredSecret.Data["namespaces"]
			secretChanged = true
		}
	}

	if secretChanged {
		if err = workloads.UpdateSecret(existingSecret, *sr.Client); err != nil {
			sr.Logger.Error(err, "reconcileClusterPermissionsSecret: failed to update secret", "name", existingSecret.Name, "namespace", existingSecret.Namespace)
			return err
		}
		sr.Logger.V(0).Info("reconcileClusterPermissionsSecret: secret upadted")
	}

	return nil
}

func (sr *SecretReconciler) DeleteManagedNamespaceFromClusterSecret(managedNamespace, previouslyManagingNamespace string) error {
	listOptions := ctrlClient.MatchingLabels{
		common.ArgoCDSecretTypeLabel: "cluster",
	}

	// List all the secrets created for ArgoCD using the label selector, and retrieve them
	secrets, err := workloads.ListSecrets(previouslyManagingNamespace, *sr.Client, listOptions)
	if err != nil {
		// TO DO: log error retrieving list of secrets to be evaluated in namespace
		return fmt.Errorf("deleteManagedNamespaceFromClusterSecret: unable to list cluster secrets in namespace %s: %w", previouslyManagingNamespace, err)
	}

	for _, secret := range secrets.Items {
		if string(secret.Data["server"]) != common.ArgoCDDefaultServer {
			continue
		}
		if namespaces, ok := secret.Data["namespaces"]; ok {
			namespaceList := strings.Split(string(namespaces), ",")
			var result []string

			for _, n := range namespaceList {
				// remove the namespace from the list of managed namespaces
				if strings.TrimSpace(n) == managedNamespace {
					continue
				}
				result = append(result, strings.TrimSpace(n))
				sort.Strings(result)
				secret.Data["namespaces"] = []byte(strings.Join(result, ","))
			}
			// Update the secret with the updated list of namespaces
			if err = workloads.UpdateSecret(&secret, *sr.Client); err != nil {
				// TO DO: log error updating cluster secret (warn)
				return fmt.Errorf("deleteManagedNamespaceFromClusterSecret: unable to update cluster secret %s in namespace %s: %w", secret.Name, previouslyManagingNamespace, err)
			}
		}
	}

	return nil
}
