package workloads

import (
	"github.com/jaideepr97/argocd-operator-rewrite/pkg/argoutil"
	template "github.com/openshift/api/template/v1"
)

var (
	templateAPIFound = false
)

// IsTemplateAPIAvailable returns true if the template API is present.
func IsTemplateAPIAvailable() bool {
	return templateAPIFound
}

// VerifyTemplateAPI will verify that the template API is present.
func VerifyTemplateAPI() error {
	found, err := argoutil.VerifyAPI(template.SchemeGroupVersion.Group, template.SchemeGroupVersion.Version)
	if err != nil {
		return err
	}
	templateAPIFound = found
	return nil
}
