package argoutil

import (
	"fmt"

	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/discovery"
)

// VerifyAPI returns true if the given group/version is present in the cluster.
func VerifyAPI(group string, version string) (bool, error) {
	k8s, err := GetK8sClient()
	if err != nil {
		return false, fmt.Errorf("VerifyAPI: unable to get client: %w", err)
	}

	gv := schema.GroupVersion{
		Group:   group,
		Version: version,
	}

	if err = discovery.ServerSupportsVersion(k8s, gv); err != nil {
		// log that API is missing (trace/debug)
		return false, nil
	}

	// log that API was verified (trace/debug)
	return true, nil
}
