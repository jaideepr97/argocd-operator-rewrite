package reposerver

// Defaults
const (
	// ArgoCDDefaultRepoMetricsPort is the default listen port for the Argo CD repo server metrics.
	ArgoCDDefaultRepoMetricsPort = 8084

	// ArgoCDDefaultRepoServerPort is the default listen port for the Argo CD repo server.
	ArgoCDDefaultRepoServerPort = 8081
)

// Keys
const (
	// ArgoCDRepoImageEnvName is the environment variable used to get the image
	// to used for the Dex container.
	ArgoCDRepoImageEnvName = "ARGOCD_REPOSERVER_IMAGE"

	// ArgoCDRepoServerComponent is the name of the repo server control plane component
	ArgoCDRepoServerComponent = "repo-server"

	// ArgoCDRepoServerSuffix is the default suffix for repo-server resources
	ArgoCDRepoServerSuffix = "repo-server"

	// ArgoCDRepoServerTLSSuffix is the default suffix for repo-server tls resources
	ArgoCDRepoServerTLSSuffix = "repo-server-tls"
)

// Values
const (
	// ArgoCDRepoServerTLSSecretName is the name of the TLS secret for the repo-server
	ArgoCDRepoServerTLSSecretName = "argocd-repo-server-tls"
)
