package redis

// Defaults
const (
	// ArgoCDDefaultRedisConfigPath is the default Redis configuration directory when not specified.
	ArgoCDDefaultRedisConfigPath = "/var/lib/redis"

	// ArgoCDDefaultRedisHAReplicas is the defaul number of replicas for Redis when rinning in HA mode.
	ArgoCDDefaultRedisHAReplicas = int32(3)

	// ArgoCDDefaultRedisHAProxyImage is the default Redis HAProxy image to use when not specified.
	ArgoCDDefaultRedisHAProxyImage = "haproxy"

	// ArgoCDDefaultRedisHAProxyVersion is the default Redis HAProxy image tag to use when not specified.
	ArgoCDDefaultRedisHAProxyVersion = "sha256:7392fbbbb53e9e063ca94891da6656e6062f9d021c0e514888a91535b9f73231" // 2.0.25-alpine

	// ArgoCDDefaultRedisImage is the Redis container image to use when not specified.
	ArgoCDDefaultRedisImage = "redis"

	// ArgoCDDefaultRedisPort is the default listen port for Redis.
	ArgoCDDefaultRedisPort = 6379

	// ArgoCDDefaultRedisSentinelPort is the default listen port for Redis sentinel.
	ArgoCDDefaultRedisSentinelPort = 26379

	//ArgoCDDefaultRedisSuffix is the default suffix to use for Redis resources.
	ArgoCDDefaultRedisSuffix = "redis"

	// ArgoCDRedhisHAProxysuffix is the suffix to use for the Redis HA proxy resources
	ArgoCDRedhisHAProxysuffix = "redis-ha-haproxy"

	// ArgoCDDefaultRedisVersion is the Redis container image tag to use when not specified.
	ArgoCDDefaultRedisVersion = "sha256:8061ca607db2a0c80010aeb5fc9bed0253448bc68711eaa14253a392f6c48280" // 6.2.4-alpine

	// ArgoCDDefaultRedisVersionHA is the Redis container image tag to use when not specified in HA mode.
	ArgoCDDefaultRedisVersionHA = "sha256:8061ca607db2a0c80010aeb5fc9bed0253448bc68711eaa14253a392f6c48280" // 6.2.4-alpine
)

// Keys
const (
	// ArgoCDRedisHAImageEnvName is the environment variable used to get the image
	// to used for the the Redis container in HA mode.
	ArgoCDRedisHAImageEnvName = "ARGOCD_REDIS_HA_IMAGE"

	// ArgoCDRedisImageEnvName is the environment variable used to get the image
	// to used for the Redis container.
	ArgoCDRedisImageEnvName = "ARGOCD_REDIS_IMAGE"

	// ArgoCDRedisHAProxyImageEnvName is the environment variable used to get the image
	// to used for the Redis HA Proxy container.
	ArgoCDRedisHAProxyImageEnvName = "ARGOCD_REDIS_HA_PROXY_IMAGE"
)

// Values
const (
	// ArgoCDRedisComponent is the name of the Redis control plane component
	ArgoCDRedisComponent = "redis"

	// ArgoCDRedisHAComponent is the name of the Redis HA control plane component
	ArgoCDRedisHAComponent = "redis-ha"
)

// Miscellaneous constants
const (
	// ArgoCDRedisHAConfigMapName is the upstream ArgoCD Redis HA ConfigMap name.
	ArgoCDRedisHAConfigMapName = "argocd-redis-ha-configmap"

	// ArgoCDRedisHAHealthConfigMapName is the upstream ArgoCD Redis HA Health ConfigMap name.
	ArgoCDRedisHAHealthConfigMapName = "argocd-redis-ha-health-configmap"

	// ArgoCDRedisProbesConfigMapName is the upstream ArgoCD Redis Probes ConfigMap name.
	ArgoCDRedisProbesConfigMapName = "argocd-redis-ha-probes"

	// ArgoCDRedisServerTLSSecretName is the name of the TLS secret for the redis-server
	ArgoCDRedisServerTLSSecretName = "argocd-operator-redis-tls"
)
