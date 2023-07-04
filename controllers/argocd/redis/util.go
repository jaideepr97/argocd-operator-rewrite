package redis

import (
	"github.com/jaideepr97/argocd-operator-rewrite/api/v1alpha1"
	"github.com/jaideepr97/argocd-operator-rewrite/pkg/argoutil"
)

// GetRedisServerAddress will return the Redis service address for the given ArgoCD.
func GetRedisServerAddress(cr *v1alpha1.ArgoCD) string {
	if cr.Spec.HA != nil {
		return getRedisHAProxyAddress(cr)
	}
	return argoutil.FqdnServiceRef(ArgoCDDefaultRedisSuffix, cr.Name, cr.Namespace, ArgoCDDefaultRedisPort)
}

// getRedisHAProxyAddress will return the Redis HA Proxy service address for the given ArgoCD.
func getRedisHAProxyAddress(cr *v1alpha1.ArgoCD) string {
	return argoutil.FqdnServiceRef(ArgoCDRedhisHAProxysuffix, cr.Name, cr.Namespace, ArgoCDDefaultRedisPort)
}

func IsRedisTLSVerificationDisabled(cr *v1alpha1.ArgoCD) bool {
	return cr.Spec.Redis.DisableTLSVerification
}
