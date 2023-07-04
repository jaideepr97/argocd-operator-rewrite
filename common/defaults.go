// Copyright 2020 ArgoCD Operator Developers
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// 	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package common

const (
	// ArgoCDDefaultLogLevel is the default log level to be used by all ArgoCD components.
	ArgoCDDefaultLogLevel = "info"

	// ArgoCDDefaultLogFormat is the default log format to be used by all ArgoCD components.
	ArgoCDDefaultLogFormat = "text"

	// ArgoCDDefaultAdminPasswordLength is the length of the generated default admin password.
	ArgoCDDefaultAdminPasswordLength = 32

	// ArgoCDDefaultAdminPasswordNumDigits is the number of digits to use for the generated default admin password.
	ArgoCDDefaultAdminPasswordNumDigits = 5

	// ArgoCDDefaultAdminPasswordNumSymbols is the number of symbols to use for the generated default admin password.
	ArgoCDDefaultAdminPasswordNumSymbols = 0

	// ArgoCDDefaultApplicationInstanceLabelKey is the default app name as a tracking label.
	ArgoCDDefaultApplicationInstanceLabelKey = "app.kubernetes.io/instance"

	// ArgoCDDefaultArgoImage is the ArgoCD container image to use when not specified.
	ArgoCDDefaultArgoImage = "quay.io/argoproj/argocd"

	// ArgoCDDefaultArgoVersion is the Argo CD container image digest to use when version not specified.
	ArgoCDDefaultArgoVersion = "sha256:0fd690bd7b89bd6f947b4000de33abd53ebcd36b57216f1c675a1127707b5eef" // v2.6.3

	// ArgoCDDefaultBackupKeyLength is the length of the generated default backup key.
	ArgoCDDefaultBackupKeyLength = 32

	// ArgoCDDefaultBackupKeyNumDigits is the number of digits to use for the generated default backup key.
	ArgoCDDefaultBackupKeyNumDigits = 5

	// ArgoCDDefaultBackupKeyNumSymbols is the number of symbols to use for the generated default backup key.
	ArgoCDDefaultBackupKeyNumSymbols = 5

	// ArgoCDDefaultConfigManagementPlugins is the default configuration value for the config management plugins.
	ArgoCDDefaultConfigManagementPlugins = ""

	// ArgoCDDefaultExportJobImage is the export job container image to use when not specified.
	ArgoCDDefaultExportJobImage = "quay.io/argoprojlabs/argocd-operator-util"

	// ArgoCDDefaultExportJobVersion is the export job container image tag to use when not specified.
	ArgoCDDefaultExportJobVersion = "sha256:6f80965a2bef1c80875be0995b18d9be5a6ad4af841cbc170ed3c60101a7deb2" // 0.5.0

	// ArgoCDDefaultExportLocalCapicity is the default capacity to use for local export.
	ArgoCDDefaultExportLocalCapicity = "2Gi"

	// ArgoCDDefaultGATrackingID is the default Google Analytics tracking ID.
	ArgoCDDefaultGATrackingID = ""

	// ArgoCDDefaultGAAnonymizeUsers is the default value for anonymizing google analytics users.
	ArgoCDDefaultGAAnonymizeUsers = false

	// ArgoCDDefaultHelpChatURL is the default help chat URL.
	ArgoCDDefaultHelpChatURL = ""

	// ArgoCDDefaultHelpChatText is the default help chat text.
	ArgoCDDefaultHelpChatText = ""

	// ArgoCDDefaultIngressPath is the path to use for the Ingress when not specified.
	ArgoCDDefaultIngressPath = "/"

	// ArgoCDDefaultKustomizeBuildOptions is the default kustomize build options.
	ArgoCDDefaultKustomizeBuildOptions = ""

	// ArgoCDDefaultOIDCConfig is the default OIDC configuration.
	ArgoCDDefaultOIDCConfig = ""

	// ArgoCDDefaultRBACPolicy is the default RBAC policy CSV data.
	ArgoCDDefaultRBACPolicy = ""

	// ArgoCDDefaultRBACDefaultPolicy is the default Argo CD RBAC policy.
	ArgoCDDefaultRBACDefaultPolicy = "role:readonly"

	// ArgoCDDefaultRBACScopes is the default Argo CD RBAC scopes.
	ArgoCDDefaultRBACScopes = "[groups]"

	// ArgoCDDefaultRepositories is the default repositories.
	ArgoCDDefaultRepositories = ""

	// ArgoCDDefaultRepositoryCredentials is the default repository credentials
	ArgoCDDefaultRepositoryCredentials = ""

	// ArgoCDDefaultResourceCustomizations is the default resource customizations.
	ArgoCDDefaultResourceCustomizations = ""

	// ArgoCDDefaultResourceExclusions is the default resource exclusions.
	ArgoCDDefaultResourceExclusions = ""

	// ArgoCDDefaultResourceInclusions is the default resource inclusions.
	ArgoCDDefaultResourceInclusions = ""

	// ArgoCDDefaultRSAKeySize is the default RSA key size when not specified.
	ArgoCDDefaultRSAKeySize = 2048

	// ArgoCDDefaultSSHKnownHosts is the default SSH Known hosts data.
	ArgoCDDefaultSSHKnownHosts = `bitbucket.org ssh-rsa AAAAB3NzaC1yc2EAAAABIwAAAQEAubiN81eDcafrgMeLzaFPsw2kNvEcqTKl/VqLat/MaB33pZy0y3rJZtnqwR2qOOvbwKZYKiEO1O6VqNEBxKvJJelCq0dTXWT5pbO2gDXC6h6QDXCaHo6pOHGPUy+YBaGQRGuSusMEASYiWunYN0vCAI8QaXnWMXNMdFP3jHAJH0eDsoiGnLPBlBp4TNm6rYI74nMzgz3B9IikW4WVK+dc8KZJZWYjAuORU3jc1c/NPskD2ASinf8v3xnfXeukU0sJ5N6m5E8VLjObPEO+mN2t/FZTMZLiFqPWc/ALSqnMnnhwrNi2rbfg/rd/IpL8Le3pSBne8+seeFVBoGqzHM9yXw==
github.com ssh-rsa AAAAB3NzaC1yc2EAAAABIwAAAQEAq2A7hRGmdnm9tUDbO9IDSwBK6TbQa+PXYPCPy6rbTrTtw7PHkccKrpp0yVhp5HdEIcKr6pLlVDBfOLX9QUsyCOV0wzfjIJNlGEYsdlLJizHhbn2mUjvSAHQqZETYP81eFzLQNnPHt4EVVUh7VfDESU84KezmD5QlWpXLmvU31/yMf+Se8xhHTvKSCZIFImWwoG6mbUoWf9nzpIoaSjB+weqqUUmpaaasXVal72J+UX2B+2RPW3RcT0eOzQgqlJL3RKrTJvdsjE3JEAvGq3lGHSZXy28G3skua2SmVi/w4yCE6gbODqnTWlg7+wC604ydGXA8VJiS5ap43JXiUFFAaQ==
gitlab.com ecdsa-sha2-nistp256 AAAAE2VjZHNhLXNoYTItbmlzdHAyNTYAAAAIbmlzdHAyNTYAAABBBFSMqzJeV9rUzU4kWitGjeR4PWSa29SPqJ1fVkhtj3Hw9xjLVXVYrU9QlYWrOLXBpQ6KWjbjTDTdDkoohFzgbEY=
gitlab.com ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIAfuCHKVTjquxvt6CM6tdG4SLp1Btn/nOeHHE5UOzRdf
gitlab.com ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCsj2bNKTBSpIYDEGk9KxsGh3mySTRgMtXL583qmBpzeQ+jqCMRgBqB98u3z++J1sKlXHWfM9dyhSevkMwSbhoR8XIq/U0tCNyokEi/ueaBMCvbcTHhO7FcwzY92WK4Yt0aGROY5qX2UKSeOvuP4D6TPqKF1onrSzH9bx9XUf2lEdWT/ia1NEKjunUqu1xOB/StKDHMoX4/OKyIzuS0q/T1zOATthvasJFoPrAjkohTyaDUz2LN5JoH839hViyEG82yB+MjcFV5MU3N1l1QL3cVUCh93xSaua1N85qivl+siMkPGbO5xR/En4iEY6K2XPASUEMaieWVNTRCtJ4S8H+9
ssh.dev.azure.com ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC7Hr1oTWqNqOlzGJOfGJ4NakVyIzf1rXYd4d7wo6jBlkLvCA4odBlL0mDUyZ0/QUfTTqeu+tm22gOsv+VrVTMk6vwRU75gY/y9ut5Mb3bR5BV58dKXyq9A9UeB5Cakehn5Zgm6x1mKoVyf+FFn26iYqXJRgzIZZcZ5V6hrE0Qg39kZm4az48o0AUbf6Sp4SLdvnuMa2sVNwHBboS7EJkm57XQPVU3/QpyNLHbWDdzwtrlS+ez30S3AdYhLKEOxAG8weOnyrtLJAUen9mTkol8oII1edf7mWWbWVf0nBmly21+nZcmCTISQBtdcyPaEno7fFQMDD26/s0lfKob4Kw8H
vs-ssh.visualstudio.com ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC7Hr1oTWqNqOlzGJOfGJ4NakVyIzf1rXYd4d7wo6jBlkLvCA4odBlL0mDUyZ0/QUfTTqeu+tm22gOsv+VrVTMk6vwRU75gY/y9ut5Mb3bR5BV58dKXyq9A9UeB5Cakehn5Zgm6x1mKoVyf+FFn26iYqXJRgzIZZcZ5V6hrE0Qg39kZm4az48o0AUbf6Sp4SLdvnuMa2sVNwHBboS7EJkm57XQPVU3/QpyNLHbWDdzwtrlS+ez30S3AdYhLKEOxAG8weOnyrtLJAUen9mTkol8oII1edf7mWWbWVf0nBmly21+nZcmCTISQBtdcyPaEno7fFQMDD26/s0lfKob4Kw8H
github.com ecdsa-sha2-nistp256 AAAAE2VjZHNhLXNoYTItbmlzdHAyNTYAAAAIbmlzdHAyNTYAAABBBEmKSENjQEezOmxkZMy7opKgwFB9nkt5YRrYMjNuG5N87uRgg6CLrbo5wAdT/y6v0mKV0U2w0WZ2YB/++Tpockg=
github.com ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIOMqqnkVzrm0SdG6UOoqKLsabgH5C9okWi0dh2l9GKJl
`
)

// DefaultLabels returns the default set of labels for controllers.
func DefaultLabels(name, component string) map[string]string {
	return map[string]string{
		ArgoCDKeyName:      name,
		ArgoCDKeyPartOf:    ArgoCDAppName,
		ArgoCDKeyManagedBy: name,
		ArgoCDKeyComponent: component,
	}
}

// DefaultAnnotations returns the default set of annotations for child resources of ArgoCD
func DefaultAnnotations(name string, namespace string) map[string]string {
	return map[string]string{
		AnnotationName:      name,
		AnnotationNamespace: namespace,
	}
}

// DefaultNodeSelector returns the defult nodeSelector for ArgoCD workloads
func DefaultNodeSelector() map[string]string {
	return map[string]string{
		"kubernetes.io/os": "linux",
	}
}
