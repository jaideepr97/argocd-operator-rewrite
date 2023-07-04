package permissions

import (
	"context"
	"fmt"

	"github.com/jaideepr97/argocd-operator-rewrite/api/v1alpha1"
	"github.com/jaideepr97/argocd-operator-rewrite/pkg/argoutil"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// newClusterclusterRoleBinding returns a new clusterclusterRoleBinding instance.
func newClusterclusterRoleBinding(instanceName, instanceNamespace, component string, instanceAnnotations map[string]string, rules []rbacv1.PolicyRule) *rbacv1.ClusterRoleBinding {
	return &rbacv1.ClusterRoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:        argoutil.GenerateUniqueResourceName(instanceName, instanceNamespace, component),
			Labels:      argoutil.LabelsForCluster(instanceName, component),
			Annotations: argoutil.AnnotationsForCluster(instanceName, instanceNamespace, instanceAnnotations),
		},
	}
}

func CreateClusterRoleBinding(cr *v1alpha1.ArgoCD, clusterRoleBinding *rbacv1.ClusterRole, scheme *runtime.Scheme, client client.Client) error {
	if err := controllerutil.SetControllerReference(cr, clusterRoleBinding, scheme); err != nil {
		// TO DO: log error here (error)
		return fmt.Errorf("CreateClusterRole: failed to set instance %s in namespace %s as owner of clusterRoleBinding %s: %w", cr.Name, cr.Namespace, clusterRoleBinding.Name, err)
	}

	// TO DO: log creation of clusterRoleBinding here (info)
	if err := client.Create(context.TODO(), clusterRoleBinding); err != nil {
		// TO DO: log error here (error)
		return fmt.Errorf("CreateClusterRole: failed to create clusterRoleBinding %s: %w", clusterRoleBinding.Name, err)
	}
	return nil
}
