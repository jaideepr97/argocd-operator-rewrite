package secret

import (
	"fmt"
	"sort"
	"strings"

	"github.com/jaideepr97/argocd-operator-rewrite/api/v1alpha1"
	"github.com/jaideepr97/argocd-operator-rewrite/common"
	"github.com/jaideepr97/argocd-operator-rewrite/pkg/workloads"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	ctrlClient "sigs.k8s.io/controller-runtime/pkg/client"
)

type SecretReconciler struct {
	Client        *client.Client
	Scheme        *runtime.Scheme
	Instance      *v1alpha1.ArgoCD
	ClusterScoped bool
}

func (sr *SecretReconciler) Reconcile() error {
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
