/*
Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	"github.com/jaideepr97/argocd-operator-rewrite/common"
	routev1 "github.com/openshift/api/route/v1"
	autoscalingv1 "k8s.io/api/autoscaling/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// ApplicationSetSpec defines the desired state of the ApplicationSet controller
type ApplicationSetSpec struct {
	// Env lets you specify environment for applicationSet controller pods
	Env []corev1.EnvVar `json:"env,omitempty"`

	// ExtraCommandArgs allows users to pass command line arguments to ApplicationSet controller.
	// They get added to default command line arguments provided by the operator.
	// Please note that the command line arguments provided as part of ExtraCommandArgs
	// will not overwrite the default command line arguments.
	ExtraCommandArgs []string `json:"extraCommandArgs,omitempty"`

	// Image is the Argo CD ApplicationSet image (optional)
	Image string `json:"image,omitempty"`

	// Version is the Argo CD ApplicationSet image tag. (optional)
	Version string `json:"version,omitempty"`

	// Resources defines the Compute Resources required by the container for ApplicationSet.
	Resources *corev1.ResourceRequirements `json:"resources,omitempty"`

	// LogLevel describes the log level that should be used by the ApplicationSet controller. Defaults to ArgoCDDefaultLogLevel if not set.  Valid options are debug, info, error, and warn.
	LogLevel string `json:"logLevel,omitempty"`

	WebhookServer *WebhookServerSpec `json:"webhookServer,omitempty"`
}

// WebhookServerSpec defines the options for the ApplicationSet Webhook Server component.
type WebhookServerSpec struct {
	// Host is the hostname to use for Ingress/Route resources.
	Host string `json:"host,omitempty"`

	Ingress *IngressSpec `json:"ingress,omitempty"`

	Route *RouteSpec `json:"route,omitempty"`
}

// ArgoCDIngressSpec defines the desired state for the Ingress resources.
type IngressSpec struct {
	// Annotations is the map of annotations to apply to the Ingress.
	Annotations map[string]string `json:"annotations,omitempty"`

	// IngressClassName for the Ingress resource.
	IngressClassName *string `json:"ingressClassName,omitempty"`

	// Path used for the Ingress resource.
	Path string `json:"path,omitempty"`

	// TLS configuration. Currently the Ingress only supports a single TLS
	// port, 443. If multiple members of this list specify different hosts, they
	// will be multiplexed on the same port according to the hostname specified
	// through the SNI TLS extension, if the ingress controller fulfilling the
	// ingress supports SNI.
	// +optional
	TLS []networkingv1.IngressTLS `json:"tls,omitempty"`
}

// ArgoCDRouteSpec defines the desired state for an OpenShift Route.
type RouteSpec struct {
	// Annotations is the map of annotations to use for the Route resource.
	Annotations map[string]string `json:"annotations,omitempty"`

	// Labels is the map of labels to use for the Route resource
	Labels map[string]string `json:"labels,omitempty"`

	// Path the router watches for, to route traffic for to the service.
	Path string `json:"path,omitempty"`

	// TLS provides the ability to configure certificates and termination for the Route.
	TLS *routev1.TLSConfig `json:"tls,omitempty"`

	// WildcardPolicy if any for the route. Currently only 'Subdomain' or 'None' is allowed.
	WildcardPolicy *routev1.WildcardPolicyType `json:"wildcardPolicy,omitempty"`
}

// ApplicationControllerSpec defines the desired state for the ArgoCD Application Controller component.
type ApplicationControllerSpec struct {
	Processors *ProcessorsSpec `json:"processors,omitempty"`

	// LogLevel refers to the log level used by the Application Controller component. Defaults to ArgoCDDefaultLogLevel if not configured. Valid options are debug, info, error, and warn.
	LogLevel string `json:"logLevel,omitempty"`

	// LogFormat refers to the log format used by the Application Controller component. Defaults to ArgoCDDefaultLogFormat if not configured. Valid options are text or json.
	LogFormat string `json:"logFormat,omitempty"`

	// Resources defines the Compute Resources required by the container for the Application Controller.
	Resources *corev1.ResourceRequirements `json:"resources,omitempty"`

	// ParallelismLimit defines the limit for parallel kubectl operations
	ParallelismLimit int32 `json:"parallelismLimit,omitempty"`

	// AppSync is used to control the sync frequency, by default the ArgoCD
	// controller polls Git every 3m.
	//
	// Set this to a duration, e.g. 10m or 600s to control the synchronisation
	// frequency.
	// +optional
	AppSync string `json:"appSync,omitempty"`

	Sharding *ShardingSpec `json:"sharding,omitempty"`

	// Env lets you specify environment for application controller pods
	Env []corev1.EnvVar `json:"env,omitempty"`
}

// ShardingSpec defines the options available for enabling sharding for the Application Controller component.
type ShardingSpec struct {
	// Replicas defines the number of shards of Application controller to run.
	Replicas int32 `json:"replicas,omitempty"`

	DynamicScaling *DynamicScalingSpec `json:"dynamicScaling,omitempty"`
}

// DynamicScalingSpec defines the options available for dynamically scaling up/down the no. of Application Controller shards.
type DynamicScalingSpec struct {
	// MinShards defines the minimum number of shards at any given point
	// +kubebuilder:validation:Minimum=1
	MinShards int32 `json:"minShards,omitempty"`

	// MaxShards defines the maximum number of shards at any given point
	MaxShards int32 `json:"maxShards,omitempty"`

	// ClustersPerShard defines the maximum number of clusters managed by each argocd shard
	// +kubebuilder:validation:Minimum=1
	ClustersPerShard int32 `json:"clustersPerShard,omitempty"`
}

// ProcessorsSpec defines the desired state for ArgoCD Application Controller processors.
type ProcessorsSpec struct {
	// Operation is the number of application operation processors.
	Operation int32 `json:"operation,omitempty"`

	// Status is the number of application status processors.
	Status int32 `json:"status,omitempty"`
}

// ImportSpec defines the desired state for the ArgoCD import/restore process.
type ImportSpec struct {
	// Name of an ArgoCDExport from which to import data.
	Name string `json:"name"`

	// Namespace for the ArgoCDExport, defaults to the same namespace as the ArgoCD.
	Namespace *string `json:"namespace,omitempty"`
}

// HASpec defines the desired state for High Availability support for Argo CD.
type HASpec struct {
	// RedisProxyImage is the Redis HAProxy container image.
	RedisProxyImage string `json:"redisProxyImage,omitempty"`

	// RedisProxyVersion is the Redis HAProxy container image tag.
	RedisProxyVersion string `json:"redisProxyVersion,omitempty"`

	// Resources defines the Compute Resources required by the container for HA.
	Resources *corev1.ResourceRequirements `json:"resources,omitempty"`
}

// InitialSSHKnownHosts defines the SSH known hosts data upon creation of the cluster for connecting Git repositories via SSH.
type SSHHostsSpec struct {
	// ExcludeDefaultHosts describes whether you would like to include the default
	// list of SSH Known Hosts provided by ArgoCD.
	ExcludeDefaultHosts bool `json:"excludedefaulthosts,omitempty"`

	// Keys describes a custom set of SSH Known Hosts that you would like to
	// have included in your ArgoCD server.
	Keys string `json:"keys,omitempty"`
}

// KustomizeVersionSpec is used to specify information about a kustomize version to be used within ArgoCD.
type KustomizeVersionSpec struct {
	// Version is a configured kustomize version in the format of vX.Y.Z
	Version string `json:"version,omitempty"`
	// Path is the path to a configured kustomize version on the filesystem of your repo server.
	Path string `json:"path,omitempty"`
}

// MonitoringSpec is used to configure workload status monitoring for a given Argo CD instance.
// It triggers creation of serviceMonitor and PrometheusRules that alert users when a given workload
// status meets a certain criteria. For e.g, it can fire an alert if the application controller is
// pending for x mins consecutively.
type MonitoringSpec struct {
	// Enabled defines whether workload status monitoring is enabled for this instance or not
	Enabled bool `json:"enabled"`
}

// NodePlacementSpec is used to specify NodeSelector and Tolerations for Argo CD workloads
type NodePlacementSpec struct {
	// NodeSelector is a field of PodSpec, it is a map of key value pairs used for node selection
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`
	// Tolerations allow the pods to schedule onto nodes with matching taints
	Tolerations []corev1.Toleration `json:"tolerations,omitempty"`
}

// NotificationsSpec defines the desired state of Argo CD Notifications controller
type NotificationsSpec struct {
	// Replicas defines the number of replicas to run for notifications-controller
	Replicas *int32 `json:"replicas,omitempty"`

	// Env let you specify environment variables for Notifications pods
	Env []corev1.EnvVar `json:"env,omitempty"`

	// Image is the Argo CD Notifications image (optional)
	Image string `json:"image,omitempty"`

	// Version is the Argo CD Notifications image tag. (optional)
	Version string `json:"version,omitempty"`

	// Resources defines the Compute Resources required by the container for Argo CD Notifications.
	Resources *corev1.ResourceRequirements `json:"resources,omitempty"`

	// LogLevel describes the log level that should be used by the argocd-notifications. Defaults to ArgoCDDefaultLogLevel if not set.  Valid options are debug,info, error, and warn.
	LogLevel string `json:"logLevel,omitempty"`
}

// RBACSpec defines the desired state for the Argo CD RBAC configuration.
type RBACSpec struct {
	// DefaultPolicy is the name of the default role which Argo CD will falls back to, when
	// authorizing API requests (optional). If omitted or empty, users may be still be able to login,
	// but will see no apps, projects, etc...
	DefaultPolicy *string `json:"defaultPolicy,omitempty"`

	// Policy is CSV containing user-defined RBAC policies and role definitions.
	// Policy rules are in the form:
	//   p, subject, resource, action, object, effect
	// Role definitions and bindings are in the form:
	//   g, subject, inherited-subject
	// See https://github.com/argoproj/argo-cd/blob/master/docs/operator-manual/rbac.md for additional information.
	Policy *string `json:"policy,omitempty"`

	// Scopes controls which OIDC scopes to examine during rbac enforcement (in addition to `sub` scope).
	// If omitted, defaults to: '[groups]'.
	Scopes *string `json:"scopes,omitempty"`

	// PolicyMatcherMode configures the matchers function mode for casbin.
	// There are two options for this, 'glob' for glob matcher or 'regex' for regex matcher.
	PolicyMatcherMode *string `json:"policyMatcherMode,omitempty"`
}

// RedisSpec defines the desired state for the Redis server component.
type RedisSpec struct {
	// Image is the Redis container image.
	Image string `json:"image,omitempty"`

	// Resources defines the Compute Resources required by the container for Redis.
	Resources *corev1.ResourceRequirements `json:"resources,omitempty"`

	// Version is the Redis container image tag.
	Version string `json:"version,omitempty"`

	// DisableTLSVerification defines whether redis server API should be accessed using strict TLS validation
	DisableTLSVerification bool `json:"disableTLSVerification,omitempty"`

	// AutoTLS specifies the method to use for automatic TLS configuration for the redis server
	// The value specified here can currently be:
	// - openshift - Use the OpenShift service CA to request TLS config
	AutoTLS string `json:"autotls,omitempty"`
}

// RepoSpec defines the desired state for the Argo CD repo server component.
type RepoSpec struct {
	// Extra Command arguments allows users to pass command line arguments to repo server workload. They get added to default command line arguments provided
	// by the operator.
	// Please note that the command line arguments provided as part of ExtraRepoCommandArgs will not overwrite the default command line arguments.
	ExtraRepoCommandArgs []string `json:"extraRepoCommandArgs,omitempty"`

	// LogLevel describes the log level that should be used by the Repo Server. Defaults to ArgoCDDefaultLogLevel if not set.  Valid options are debug, info, error, and warn.
	LogLevel string `json:"logLevel,omitempty"`

	// LogFormat describes the log format that should be used by the Repo Server. Defaults to ArgoCDDefaultLogFormat if not configured. Valid options are text or json.
	LogFormat string `json:"logFormat,omitempty"`

	// MountSAToken describes whether you would like to have the Repo server mount the service account token
	MountSAToken bool `json:"mountsatoken,omitempty"`

	// Replicas defines the number of replicas for argocd-repo-server. Value should be greater than or equal to 0. Default is nil.
	Replicas *int32 `json:"replicas,omitempty"`

	// Resources defines the Compute Resources required by the container for Redis.
	Resources *corev1.ResourceRequirements `json:"resources,omitempty"`

	// ServiceAccount defines the ServiceAccount user that you would like the Repo server to use
	ServiceAccount string `json:"serviceaccount,omitempty"`

	// VerifyTLS defines whether repo server API should be accessed using strict TLS validation
	VerifyTLS bool `json:"verifytls,omitempty"`

	// AutoTLS specifies the method to use for automatic TLS configuration for the repo server
	// The value specified here can currently be:
	// - openshift - Use the OpenShift service CA to request TLS config
	AutoTLS string `json:"autotls,omitempty"`

	// Image is the ArgoCD Repo Server container image.
	Image string `json:"image,omitempty"`

	// Version is the ArgoCD Repo Server container image tag.
	Version string `json:"version,omitempty"`

	// ExecTimeout specifies the timeout in seconds for tool execution
	ExecTimeout *int `json:"execTimeout,omitempty"`

	// Env lets you specify environment for repo server pods
	Env []corev1.EnvVar `json:"env,omitempty"`

	// Volumes adds volumes to the repo server deployment
	Volumes []corev1.Volume `json:"volumes,omitempty"`

	// VolumeMounts adds volumeMounts to the repo server container
	VolumeMounts []corev1.VolumeMount `json:"volumeMounts,omitempty"`

	// InitContainers defines the list of initialization containers for the repo server deployment
	InitContainers []corev1.Container `json:"initContainers,omitempty"`

	// SidecarContainers defines the list of sidecar containers for the repo server deployment
	SidecarContainers []corev1.Container `json:"sidecarContainers,omitempty"`
}

// Resource Customization for custom health check
type ResourceHealthCheck struct {
	Group string `json:"group,omitempty"`
	Kind  string `json:"kind,omitempty"`
	Check string `json:"check,omitempty"`
}

// Resource Customization for ignore difference
type ResourceIgnoreDifference struct {
	All                 *IgnoreDifferenceCustomization `json:"all,omitempty"`
	ResourceIdentifiers []ResourceIdentifiers          `json:"resourceIdentifiers,omitempty"`
}

// Resource Customization fields for ignore difference
type ResourceIdentifiers struct {
	Group         string                        `json:"group,omitempty"`
	Kind          string                        `json:"kind,omitempty"`
	Customization IgnoreDifferenceCustomization `json:"customization,omitempty"`
}

type IgnoreDifferenceCustomization struct {
	JqPathExpressions     []string `json:"jqPathExpressions,omitempty"`
	JsonPointers          []string `json:"jsonPointers,omitempty"`
	ManagedFieldsManagers []string `json:"managedFieldsManagers,omitempty"`
}

// Resource Customization for custom action
type ResourceAction struct {
	Group  string `json:"group,omitempty"`
	Kind   string `json:"kind,omitempty"`
	Action string `json:"action,omitempty"`
}

// AutoscalingSpec defines the desired state for autoscaling the Argo CD Server component.
type AutoscalingSpec struct {
	// Enabled will toggle autoscaling support for the Argo CD Server component.
	Enabled bool `json:"enabled"`

	// HPA defines the HorizontalPodAutoscaler options for the Argo CD Server component.
	HPA *autoscalingv1.HorizontalPodAutoscalerSpec `json:"hpa,omitempty"`
}

// GRPCSpec defines the desired state for the Argo CD Server GRPC options.
type GRPCSpec struct {
	// Host is the hostname to use for Ingress/Route resources.
	Host string `json:"host,omitempty"`

	Ingress *IngressSpec `json:"ingress,omitempty"`
}

// ServiceSpec defines the Service options for Argo CD Server component.
type ServiceSpec struct {
	// Type is the ServiceType to use for the Service resource.
	Type corev1.ServiceType `json:"type"`
}

// ServerSpec defines the options for the ArgoCD Server component.
type ServerSpec struct {
	Autoscale *AutoscalingSpec `json:"autoscale,omitempty"`

	GRPC *GRPCSpec `json:"grpc,omitempty"`

	// Host is the hostname to use for Ingress/Route resources.
	Host string `json:"host,omitempty"`

	Ingress *IngressSpec `json:"ingress,omitempty"`

	// Insecure toggles the insecure flag.
	Insecure bool `json:"insecure,omitempty"`

	// LogLevel refers to the log level to be used by the ArgoCD Server component. Defaults to ArgoCDDefaultLogLevel if not set.  Valid options are debug, info, error, and warn.
	LogLevel string `json:"logLevel,omitempty"`

	// LogFormat refers to the log level to be used by the ArgoCD Server component. Defaults to ArgoCDDefaultLogFormat if not configured. Valid options are text or json.
	LogFormat string `json:"logFormat,omitempty"`

	// Replicas defines the number of replicas for argocd-server. Default is nil. Value should be greater than or equal to 0. Value will be ignored if Autoscaler is enabled.
	Replicas *int32 `json:"replicas,omitempty"`

	// Resources defines the Compute Resources required by the container for the Argo CD server component.
	Resources *corev1.ResourceRequirements `json:"resources,omitempty"`

	Route *RouteSpec `json:"route,omitempty"`

	Service *ServiceSpec `json:"service,omitempty"`

	// Env lets you specify environment for API server pods
	Env []corev1.EnvVar `json:"env,omitempty"`

	// Extra Command arguments that would append to the Argo CD server command.
	// ExtraCommandArgs will not be added, if one of these commands is already part of the server command
	// with same or different value.
	ExtraCommandArgs []string `json:"extraCommandArgs,omitempty"`
}

// SSOProviderType string defines the type of SSO provider.
type SSOProviderType string

const (
	// SSOProviderTypeKeycloak means keycloak will be Installed and Integrated with Argo CD. A new realm with name argocd
	// will be created in this keycloak. This realm will have a client with name argocd that uses OpenShift v4 as Identity Provider.
	SSOProviderTypeKeycloak SSOProviderType = "keycloak"

	// SSOProviderTypeDex means dex will be Installed and Integrated with Argo CD.
	SSOProviderTypeDex SSOProviderType = "dex"
)

// DexSpec defines the desired state for the Dex server component.
type DexSpec struct {
	//Config is the dex connector configuration.
	Config string `json:"config,omitempty"`

	// Optional list of required groups a user must be a member of
	Groups []string `json:"groups,omitempty"`

	// Image is the Dex container image.
	Image string `json:"image,omitempty"`

	// OpenShiftOAuth enables OpenShift OAuth authentication for the Dex server.
	OpenShiftOAuth bool `json:"openShiftOAuth,omitempty"`

	// Resources defines the Compute Resources required by the container for Dex.
	Resources *corev1.ResourceRequirements `json:"resources,omitempty"`

	// Version is the Dex container image tag.
	Version string `json:"version,omitempty"`
}

// KeycloakSpec defines the desired state for the Keycloak component.
type KeycloakSpec struct {
	// Image is the Keycloak container image.
	Image string `json:"image,omitempty"`

	// Resources defines the Compute Resources required by the container for Keycloak.
	Resources *corev1.ResourceRequirements `json:"resources,omitempty"`

	// Custom root CA certificate for communicating with the Keycloak OIDC provider
	RootCA string `json:"rootCA,omitempty"`

	// Version is the Keycloak container image tag.
	Version string `json:"version,omitempty"`

	// VerifyTLS set to false disables strict TLS validation.
	VerifyTLS *bool `json:"verifyTLS,omitempty"`
}

// SSOSpec defines SSO provider.
type SSOSpec struct {
	// Image is the SSO container image.
	Image string `json:"image,omitempty"`
	// Provider installs and configures the given SSO Provider with Argo CD.
	Provider SSOProviderType `json:"provider,omitempty"`
	// Resources defines the Compute Resources required by the container for SSO.
	Resources *corev1.ResourceRequirements `json:"resources,omitempty"`
	// VerifyTLS set to false disables strict TLS validation.
	VerifyTLS *bool `json:"verifyTLS,omitempty"`
	// Version is the SSO container image tag.
	Version string `json:"version,omitempty"`

	Dex *DexSpec `json:"dex,omitempty"`

	Keycloak *KeycloakSpec `json:"keycloak,omitempty"`
}

// CASpec defines the CA options for ArgoCD.
type CASpec struct {
	// ConfigMapName is the name of the ConfigMap containing the CA Certificate.
	ConfigMapName string `json:"configMapName,omitempty"`

	// SecretName is the name of the Secret containing the CA Certificate and Key.
	SecretName string `json:"secretName,omitempty"`
}

// TLSSpec defines the TLS options for ArgCD.
type TLSSpec struct {
	CA *CASpec `json:"ca,omitempty"`

	// InitialCerts defines custom TLS certificates upon creation of the cluster for connecting Git repositories via HTTPS.
	InitialCerts map[string]string `json:"initialCerts,omitempty"`
}

// Banner defines an additional banner message to be displayed in Argo CD UI
// https://argo-cd.readthedocs.io/en/stable/operator-manual/custom-styles/#banners
type Banner struct {
	// Content defines the banner message content to display
	Content string `json:"content"`
	// URL defines an optional URL to be used as banner message link
	URL string `json:"url,omitempty"`
}

// ArgoCDSpec defines the desired state of ArgoCD
type ArgoCDSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	ApplicationSet *ApplicationSetSpec `json:"applicationSet,omitempty"`

	// ApplicationInstanceLabelKey is the key name where Argo CD injects the app name as a tracking label.
	ApplicationInstanceLabelKey string `json:"applicationInstanceLabelKey,omitempty"`

	// ConfigManagementPlugins is used to specify additional config management plugins.
	ConfigManagementPlugins string `json:"configManagementPlugins,omitempty"`

	Controller *ApplicationControllerSpec `json:"controller,omitempty"`

	// DisableAdmin will disable the admin user.
	DisableAdmin bool `json:"disableAdmin,omitempty"`

	// ExtraConfig can be used to add fields to Argo CD configmap that are not supported by Argo CD CRD.
	//
	// Note: ExtraConfig takes precedence over Argo CD CRD.
	// For example, A user sets `argocd.Spec.DisableAdmin` = true and also
	// `a.Spec.ExtraConfig["admin.enabled"]` = true. In this case, operator updates
	// Argo CD Configmap as follows -> argocd-cm.Data["admin.enabled"] = true.
	ExtraConfig map[string]string `json:"extraConfig,omitempty"`

	// GATrackingID is the google analytics tracking ID to use.
	GATrackingID string `json:"gaTrackingID,omitempty"`

	// GAAnonymizeUsers toggles user IDs being hashed before sending to google analytics.
	GAAnonymizeUsers bool `json:"gaAnonymizeUsers,omitempty"`

	HA *HASpec `json:"ha,omitempty"`

	// HelpChatURL is the URL for getting chat help, this will typically be your Slack channel for support.
	HelpChatURL string `json:"helpChatURL,omitempty"`

	// HelpChatText is the text for getting chat help, defaults to "Chat now!"
	HelpChatText string `json:"helpChatText,omitempty"`

	// Image is the ArgoCD container image for all ArgoCD components.
	Image string `json:"image,omitempty"`

	Import *ImportSpec `json:"import,omitempty"`

	InitialSSHKnownHosts *SSHHostsSpec `json:"initialSSHKnownHosts,omitempty"`

	// KustomizeBuildOptions is used to specify build options/parameters to use with `kustomize build`.
	KustomizeBuildOptions string `json:"kustomizeBuildOptions,omitempty"`

	// KustomizeVersions is a listing of configured versions of Kustomize to be made available within ArgoCD.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Kustomize Build Options'",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:text","urn:alm:descriptor:com.tectonic.ui:advanced"}
	KustomizeVersions []KustomizeVersionSpec `json:"kustomizeVersions,omitempty"`

	// OIDCConfig is the OIDC configuration as an alternative to dex.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="OIDC Config'",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:text","urn:alm:descriptor:com.tectonic.ui:advanced"}
	OIDCConfig string `json:"oidcConfig,omitempty"`

	Monitoring MonitoringSpec `json:"monitoring,omitempty"`

	NodePlacement *NodePlacementSpec `json:"nodePlacement,omitempty"`

	Notifications *NotificationsSpec `json:"notifications,omitempty"`

	RBAC *RBACSpec `json:"rbac,omitempty"`

	Redis *RedisSpec `json:"redis,omitempty"`

	Repo *RepoSpec `json:"repo,omitempty"`

	// ResourceHealthChecks customizes resource health check behavior.
	ResourceHealthChecks []ResourceHealthCheck `json:"resourceHealthChecks,omitempty"`

	ResourceIgnoreDifferences *ResourceIgnoreDifference `json:"resourceIgnoreDifferences,omitempty"`

	// ResourceActions customizes resource action behavior.
	ResourceActions []ResourceAction `json:"resourceActions,omitempty"`

	// ResourceExclusions is used to completely ignore entire classes of resource group/kinds.
	ResourceExclusions string `json:"resourceExclusions,omitempty"`

	// ResourceInclusions is used to only include specific group/kinds in the reconciliation process.
	ResourceInclusions string `json:"resourceInclusions,omitempty"`

	// ResourceTrackingMethod defines how Argo CD should track resources that it manages
	ResourceTrackingMethod string `json:"resourceTrackingMethod,omitempty"`

	Server *ServerSpec `json:"server,omitempty"`

	// SourceNamespaces defines the namespaces application resources are allowed to be created in
	SourceNamespaces []string `json:"sourceNamespaces,omitempty"`

	SSO *SSOSpec `json:"sso,omitempty"`

	// StatusBadgeEnabled toggles application status badge feature.
	StatusBadgeEnabled bool `json:"statusBadgeEnabled,omitempty"`

	TLS *TLSSpec `json:"tls,omitempty"`

	// UsersAnonymousEnabled toggles anonymous user access. The anonymous users get default role permissions specified argocd-rbac-cm.
	UsersAnonymousEnabled bool `json:"usersAnonymousEnabled,omitempty"`

	// Version is the tag to use with the ArgoCD container image for all ArgoCD components.
	Version string `json:"version,omitempty"`

	Banner *Banner `json:"banner,omitempty"`
}

// ArgoCDStatus defines the observed state of ArgoCD
type ArgoCDStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// ApplicationController is a simple, high-level summary of where the Argo CD application controller component is in its lifecycle.
	// There are four possible ApplicationController values:
	// Pending: The Argo CD application controller component has been accepted by the Kubernetes system, but one or more of the required resources have not been created.
	// Running: All of the required Pods for the Argo CD application controller component are in a Ready state.
	// Failed: At least one of the  Argo CD application controller component Pods had a failure.
	// Unknown: The state of the Argo CD application controller component could not be obtained.
	ApplicationController string `json:"applicationController,omitempty"`

	// ApplicationSetController is a simple, high-level summary of where the Argo CD applicationSet controller component is in its lifecycle.
	// There are four possible ApplicationSetController values:
	// Pending: The Argo CD applicationSet controller component has been accepted by the Kubernetes system, but one or more of the required resources have not been created.
	// Running: All of the required Pods for the Argo CD applicationSet controller component are in a Ready state.
	// Failed: At least one of the  Argo CD applicationSet controller component Pods had a failure.
	// Unknown: The state of the Argo CD applicationSet controller component could not be obtained.
	ApplicationSetController string `json:"applicationSetController,omitempty"`

	// Dex is a simple, high-level summary of where the Argo CD Dex component is in its lifecycle.
	// There are four possible dex values:
	// Pending: The Argo CD Dex component has been accepted by the Kubernetes system, but one or more of the required resources have not been created.
	// Running: All of the required Pods for the Argo CD Dex component are in a Ready state.
	// Failed: At least one of the  Argo CD Dex component Pods had a failure.
	// Unknown: The state of the Argo CD Dex component could not be obtained.
	Dex string `json:"dex,omitempty"`

	// NotificationsController is a simple, high-level summary of where the Argo CD notifications controller component is in its lifecycle.
	// There are four possible NotificationsController values:
	// Pending: The Argo CD notifications controller component has been accepted by the Kubernetes system, but one or more of the required resources have not been created.
	// Running: All of the required Pods for the Argo CD notifications controller component are in a Ready state.
	// Failed: At least one of the  Argo CD notifications controller component Pods had a failure.
	// Unknown: The state of the Argo CD notifications controller component could not be obtained.
	NotificationsController string `json:"notificationsController,omitempty"`

	// Phase is a simple, high-level summary of where the ArgoCD is in its lifecycle.
	// There are four possible phase values:
	// Pending: The ArgoCD has been accepted by the Kubernetes system, but one or more of the required resources have not been created.
	// Available: All of the resources for the ArgoCD are ready.
	// Failed: At least one resource has experienced a failure.
	// Unknown: The state of the ArgoCD phase could not be obtained.
	Phase string `json:"phase,omitempty"`

	// Redis is a simple, high-level summary of where the Argo CD Redis component is in its lifecycle.
	// There are four possible redis values:
	// Pending: The Argo CD Redis component has been accepted by the Kubernetes system, but one or more of the required resources have not been created.
	// Running: All of the required Pods for the Argo CD Redis component are in a Ready state.
	// Failed: At least one of the  Argo CD Redis component Pods had a failure.
	// Unknown: The state of the Argo CD Redis component could not be obtained.
	Redis string `json:"redis,omitempty"`

	// Repo is a simple, high-level summary of where the Argo CD Repo component is in its lifecycle.
	// There are four possible repo values:
	// Pending: The Argo CD Repo component has been accepted by the Kubernetes system, but one or more of the required resources have not been created.
	// Running: All of the required Pods for the Argo CD Repo component are in a Ready state.
	// Failed: At least one of the  Argo CD Repo component Pods had a failure.
	// Unknown: The state of the Argo CD Repo component could not be obtained.
	Repo string `json:"repo,omitempty"`

	// Server is a simple, high-level summary of where the Argo CD server component is in its lifecycle.
	// There are four possible server values:
	// Pending: The Argo CD server component has been accepted by the Kubernetes system, but one or more of the required resources have not been created.
	// Running: All of the required Pods for the Argo CD server component are in a Ready state.
	// Failed: At least one of the  Argo CD server component Pods had a failure.
	// Unknown: The state of the Argo CD server component could not be obtained.
	Server string `json:"server,omitempty"`

	// RepoTLSChecksum contains the SHA256 checksum of the latest known state of tls.crt and tls.key in the argocd-repo-server-tls secret.
	RepoTLSChecksum string `json:"repoTLSChecksum,omitempty"`

	// RedisTLSChecksum contains the SHA256 checksum of the latest known state of tls.crt and tls.key in the argocd-operator-redis-tls secret.
	RedisTLSChecksum string `json:"redisTLSChecksum,omitempty"`

	// Host is the hostname of the server route/ingress.
	Host string `json:"host,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// ArgoCD is the Schema for the argocds API
type ArgoCD struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ArgoCDSpec   `json:"spec,omitempty"`
	Status ArgoCDStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ArgoCDList contains a list of ArgoCD
type ArgoCDList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ArgoCD `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ArgoCD{}, &ArgoCDList{})
}

// IsDeletionFinalizerPresent checks if the instance has deletion finalizer
func (argocd *ArgoCD) IsDeletionFinalizerPresent() bool {
	for _, finalizer := range argocd.GetFinalizers() {
		if finalizer == common.ArgoCDDeletionFinalizer {
			return true
		}
	}
	return false
}

// WantsAutoTLS returns true if user configured a route with reencryption
// termination policy.
func (s *ServerSpec) WantsAutoTLS() bool {
	return s.Route.TLS != nil && s.Route.TLS.Termination == routev1.TLSTerminationReencrypt
}

// WantsAutoTLS returns true if the repository server configuration has set
// the autoTLS toggle to a supported provider.
func (r *RepoSpec) WantsAutoTLS() bool {
	return r.AutoTLS == "openshift"
}

// WantsAutoTLS returns true if the redis server configuration has set
// the autoTLS toggle to a supported provider.
func (r *RedisSpec) WantsAutoTLS() bool {
	return r.AutoTLS == "openshift"
}

// ApplicationInstanceLabelKey returns either the custom application instance
// label key if set, or the default value.
func (a *ArgoCD) ApplicationInstanceLabelKey() string {
	if a.Spec.ApplicationInstanceLabelKey != "" {
		return a.Spec.ApplicationInstanceLabelKey
	} else {
		return common.ArgoCDDefaultApplicationInstanceLabelKey
	}
}

// ResourceTrackingMethod represents the Argo CD resource tracking method to use
type ResourceTrackingMethod int

const (
	ResourceTrackingMethodInvalid            ResourceTrackingMethod = -1
	ResourceTrackingMethodLabel              ResourceTrackingMethod = 0
	ResourceTrackingMethodAnnotation         ResourceTrackingMethod = 1
	ResourceTrackingMethodAnnotationAndLabel ResourceTrackingMethod = 2
)

const (
	stringResourceTrackingMethodLabel              string = "label"
	stringResourceTrackingMethodAnnotation         string = "annotation"
	stringResourceTrackingMethodAnnotationAndLabel string = "annotation+label"
)

// String returns the string representation for a ResourceTrackingMethod
func (r ResourceTrackingMethod) String() string {
	switch r {
	case ResourceTrackingMethodLabel:
		return stringResourceTrackingMethodLabel
	case ResourceTrackingMethodAnnotation:
		return stringResourceTrackingMethodAnnotation
	case ResourceTrackingMethodAnnotationAndLabel:
		return stringResourceTrackingMethodAnnotationAndLabel
	}

	// Default is to use label
	return stringResourceTrackingMethodLabel
}

// ParseResourceTrackingMethod parses a string into a resource tracking method
func ParseResourceTrackingMethod(name string) ResourceTrackingMethod {
	switch name {
	case stringResourceTrackingMethodLabel, "":
		return ResourceTrackingMethodLabel
	case stringResourceTrackingMethodAnnotation:
		return ResourceTrackingMethodAnnotation
	case stringResourceTrackingMethodAnnotationAndLabel:
		return ResourceTrackingMethodAnnotationAndLabel
	}

	return ResourceTrackingMethodInvalid
}
