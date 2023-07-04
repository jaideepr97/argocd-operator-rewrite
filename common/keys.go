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

import (
	corev1 "k8s.io/api/core/v1"
)

const (
	// ArgoCDKeyAdminEnabled is the configuration key for the admin enabled setting..
	ArgoCDKeyAdminEnabled = "admin.enabled"

	// ArgoCDKeyApplicationInstanceLabelKey is the configuration key for the application instance label.
	ArgoCDKeyApplicationInstanceLabelKey = "application.instanceLabelKey"

	// ArgoCDKeyAdminPassword is the admin password key for labels.
	ArgoCDKeyAdminPassword = "admin.password"

	// ArgoCDKeyAdminPasswordMTime is the admin password last modified key for labels.
	ArgoCDKeyAdminPasswordMTime = "admin.passwordMtime"

	// ArgoCDKeyBackupKey is the "backup key" key for ConfigMaps.
	ArgoCDKeyBackupKey = "backup.key"

	// ArgoCDKeyConfigManagementPlugins is the configuration key for config management plugins.
	ArgoCDKeyConfigManagementPlugins = "configManagementPlugins"

	// ArgoCDKeyDexOAuthRedirectURI is the key for the OAuth Redirect URI annotation.
	ArgoCDKeyDexOAuthRedirectURI = "serviceaccounts.openshift.io/oauth-redirecturi.argocd"

	// ArgoCDKeyDexConfig is the key for dex configuration.
	ArgoCDKeyDexConfig = "dex.config"

	// ArgoCDKeyFailureDomainZone is the failure-domain zone key for labels.
	ArgoCDKeyFailureDomainZone = "failure-domain.beta.kubernetes.io/zone"

	// ArgoCDKeyGATrackingID is the configuration key for the Google  Analytics Tracking ID.
	ArgoCDKeyGATrackingID = "ga.trackingid"

	// ArgoCDKeyGAAnonymizeUsers is the configuration key for the Google Analytics user anonymization.
	ArgoCDKeyGAAnonymizeUsers = "ga.anonymizeusers"

	// ArgoCDKeyHelpChatURL is the congifuration key for the help chat URL.
	ArgoCDKeyHelpChatURL = "help.chatUrl"

	// ArgoCDKeyHelpChatText is the congifuration key for the help chat text.
	ArgoCDKeyHelpChatText = "help.chatText"

	// ArgoCDKeyHostname is the resource hostname key for labels.
	ArgoCDKeyHostname = "kubernetes.io/hostname"

	// ArgoCDKeyKustomizeBuildOptions is the configuration key for the kustomize build options.
	ArgoCDKeyKustomizeBuildOptions = "kustomize.buildOptions"

	// ArgoCDKeyMetrics is the resource metrics key for labels.
	ArgoCDKeyMetrics = "metrics"

	// ArgoCDKeyOIDCConfig is the configuration key for the OIDC configuration.
	ArgoCDKeyOIDCConfig = "oidc.config"

	// ArgoCDKeyName is the resource name key for labels.
	ArgoCDKeyName = "app.kubernetes.io/name"

	// ArgoCDKeyPartOf is the resource part-of key for labels.
	ArgoCDKeyPartOf = "app.kubernetes.io/part-of"

	// ArgoCDKeyManagedBy is the managed-by key for labels.
	ArgoCDKeyManagedBy = "app.kubernetes.io/managed-by"

	// ArgoCDKeyRBACType is the label to describe if the rbac resource is meant for resource management
	// or application management
	ArgoCDKeyRBACType = "argocds.argoproj.io/rbac-type"

	// ArgoCDKeyComponent is the resource component key for labels.
	ArgoCDKeyComponent = "app.kubernetes.io/component"

	// ArgoCDKeyStatefulSetPodName is the resource StatefulSet Pod Name key for labels.
	ArgoCDKeyStatefulSetPodName = "statefulset.kubernetes.io/pod-name"

	// ArgoCDKeyPrometheus is the resource prometheus key for labels.
	ArgoCDKeyPrometheus = "prometheus"

	// ArgoCDKeyRBACPolicyCSV is the configuration key for the Argo CD RBAC policy CSV.
	ArgoCDKeyRBACPolicyCSV = "policy.csv"

	// ArgoCDKeyRBACPolicyDefault is the configuration key for the Argo CD RBAC default policy.
	ArgoCDKeyRBACPolicyDefault = "policy.default"

	// ArgoCDKeyRBACScopes is the configuration key for the Argo CD RBAC scopes.
	ArgoCDKeyRBACScopes = "scopes"

	// ArgoCDKeyRelease is the prometheus release key for labels.
	ArgoCDKeyRelease = "release"

	// ArgoCDKeyResourceCustomizations is the configuration key for resource customizations.
	ArgoCDKeyResourceCustomizations = "resource.customizations"

	// ArgoCDKeyResourceExclusions is the configuration key for resource exclusions.
	ArgoCDKeyResourceExclusions = "resource.exclusions"

	// ArgoCDKeyResourceInclusions is the configuration key for resource inclusions.
	ArgoCDKeyResourceInclusions = "resource.inclusions"

	// ArgoCDKeyResourceTrackingMethod is the configuration key for resource tracking method
	ArgoCDKeyResourceTrackingMethod = "application.resourceTrackingMethod"

	// ArgoCDKeyRepositories is the configuration key for repositories.
	ArgoCDKeyRepositories = "repositories"

	// ArgoCDKeyRepositoryCredentials is the configuration key for repository.credentials.
	ArgoCDKeyRepositoryCredentials = "repository.credentials"

	// ArgoCDKeyServerURL is the key for server url.
	ArgoCDKeyServerURL = "url"

	// ArgoCDKeySSHKnownHosts is the resource ssh_known_hosts key for labels.
	ArgoCDKeySSHKnownHosts = "ssh_known_hosts"

	// ArgoCDKeyStatusBadgeEnabled is the configuration key for enabling the status badge.
	ArgoCDKeyStatusBadgeEnabled = "statusbadge.enabled"

	// ArgoCDKeyBannerContent is the configuration key for a banner message content.
	ArgoCDKeyBannerContent = "ui.bannercontent"

	// ArgoCDKeyBannerURL is the configuration key for a banner message URL.
	ArgoCDKeyBannerURL = "ui.bannerurl"

	// ArgoCDKeyTLSCACert is the key for TLS CA certificates.
	ArgoCDKeyTLSCACert = "ca.crt"

	// ArgoCDKeyTLSCert is the key for TLS certificates.
	ArgoCDKeyTLSCert = corev1.TLSCertKey

	// ArgoCDKeyTLSPrivateKey is the key for TLS private keys.
	ArgoCDKeyTLSPrivateKey = corev1.TLSPrivateKeyKey

	// ArgoCDPolicyMatcherMode is the key for matchers function for casbin.
	// There are two options for this, 'glob' for glob matcher or 'regex' for regex matcher.
	ArgoCDPolicyMatcherMode = "policy.matchMode"

	// ArgoCDKeyTolerateUnreadyEndpounts is the resource tolerate unready endpoints key for labels.
	ArgoCDKeyTolerateUnreadyEndpounts = "service.alpha.kubernetes.io/tolerate-unready-endpoints"

	// ArgoCDKeyUsersAnonymousEnabled is the configuration key for anonymous user access.
	ArgoCDKeyUsersAnonymousEnabled = "users.anonymous.enabled"

	// ArgoCDImageEnvName is the environment variable used to get the image
	// to used for the argocd container.
	ArgoCDImageEnvName = "ARGOCD_IMAGE"

	// ArgoCDDeletionFinalizer is a finalizer to implement pre-delete hooks
	ArgoCDDeletionFinalizer = "argoproj.io/finalizer"

	// ArgoCDDefaultServer is the default server address
	ArgoCDDefaultServer = "https://kubernetes.default.svc"

	// ArgoCDSecretTypeLabel is needed for cluster secrets
	ArgoCDSecretTypeLabel = "argocd.argoproj.io/secret-type"

	// ArgoCDResourcesManagedByLabel is needed to identify namespace managed by an instance on ArgoCD
	ArgoCDResourcesManagedByLabel = "argocd.argoproj.io/resources-managed-by"

	// ArgoCDAppsManagedByLabel is needed to identify namespace mentioned as sourceNamespace on ArgoCD
	ArgoCDAppsManagedByLabel = "argocd.argoproj.io/apps-managed-by"

	// ArgoCDClusterConfigNamespacesEnvVar is the environment variable that contains the list of namespaces allowed to host cluster config
	// instances
	ArgoCDClusterConfigNamespacesEnvVar = "ARGOCD_CLUSTER_CONFIG_NAMESPACES"
)
