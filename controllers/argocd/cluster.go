package argocd

import (
	"os"

	"github.com/jaideepr97/argocd-operator-rewrite/common"
	"github.com/jaideepr97/argocd-operator-rewrite/pkg/argoutil"
	"github.com/jaideepr97/argocd-operator-rewrite/pkg/monitoring"
	"github.com/jaideepr97/argocd-operator-rewrite/pkg/networking"
	"github.com/jaideepr97/argocd-operator-rewrite/pkg/workloads"

	configv1 "github.com/openshift/api/config/v1"
)

var (
	versionAPIFound = false
)

// IsVersionAPIAvailable returns true if the version api is present
func IsVersionAPIAvailable() bool {
	return versionAPIFound
}

// VerifyVersionAPI will verify that the cluster version API is present.
func VerifyVersionAPI() error {
	found, err := argoutil.VerifyAPI(configv1.GroupName, configv1.GroupVersion.Version)
	if err != nil {
		return err
	}
	versionAPIFound = found
	return nil
}

// InspectCluster will verify the availability of extra features on the cluster, such as Prometheus and OpenShift Routes.
func InspectCluster() {
	if err := monitoring.VerifyPrometheusAPI(); err != nil {
		// TO DO: log error verifying prometheus API (warn)
	}

	if err := networking.VerifyRouteAPI(); err != nil {
		// TO DO: log error verifying prometheus API (warn)
	}

	if err := workloads.VerifyTemplateAPI(); err != nil {
		// TO DO: log error verifying prometheus API (warn)
	}

	if err := VerifyVersionAPI(); err != nil {
		// TO DO: log error verifying prometheus API (warn)
	}
}

func GetClusterConfigNamespaces() string {
	return os.Getenv(common.ArgoCDClusterConfigNamespacesEnvVar)
}

func AllowedNamespace(current string, namespaces string) bool {
	clusterConfigNamespaces := argoutil.SplitList(namespaces)
	if len(clusterConfigNamespaces) > 0 {
		if clusterConfigNamespaces[0] == "*" {
			return true
		}

		for _, n := range clusterConfigNamespaces {
			if n == current {
				return true
			}
		}
	}
	return false
}
