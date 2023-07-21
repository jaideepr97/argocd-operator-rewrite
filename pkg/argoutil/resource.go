package argoutil

import (
	"fmt"

	"github.com/jaideepr97/argocd-operator-rewrite/common"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
)

// GenerateResourceName generates names for namespace scoped scoped resources
func GenerateResourceName(instanceName, component string) string {
	return fmt.Sprintf("%s-%s", instanceName, component)
}

// GenerateUniqueResourceName generates unique names for cluster scoped resources
func GenerateUniqueResourceName(instanceName, instanceNamespace, component string) string {
	return fmt.Sprintf("%s-%s-%s", instanceName, instanceNamespace, component)
}

// LabelsForCluster returns the labels for all cluster resources.
func LabelsForCluster(instanceName, component string) map[string]string {
	return common.DefaultLabels(instanceName, component)
}

// annotationsForCluster returns the annotations for all cluster resources.
func AnnotationsForCluster(instanceName, instanceNamespace string) map[string]string {
	return common.DefaultAnnotations(instanceName, instanceNamespace)
}

// nameWithSuffix will return a name with the given suffix appended to it
func NameWithSuffix(name, suffix string) string {
	return fmt.Sprintf("%s-%s", name, suffix)
}

// ConvertLabelSelector takes a metav1.LabelSelector as input and returns a labels.Selector
func ConvertLabelSelector(LabelSelector metav1.LabelSelector) (labels.Selector, error) {
	selector, err := metav1.LabelSelectorAsSelector(&LabelSelector)
	if err != nil {
		return nil, fmt.Errorf("ConvertLabelSelector: could not convert to labels.Selector: %w", err)
	}

	return selector, nil
}
