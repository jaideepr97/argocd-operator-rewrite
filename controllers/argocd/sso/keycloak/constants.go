package sso

// Defaults
const (
	// Default is the default Keycloak Image used for the non-openshift platforms when not specified.
	ArgoCDDefaultKeycloakImage = "quay.io/keycloak/keycloak"

	// ArgoCDDefaultKeycloakVersion is the default Keycloak version used for the non-openshift platform when not specified.
	// Version: 15.0.2
	ArgoCDDefaultKeycloakVersion = "sha256:64fb81886fde61dee55091e6033481fa5ccdac62ae30a4fd29b54eb5e97df6a9"

	// ArgoCDDefaultKeycloakImageForOpenShift is the default Keycloak Image used for the OpenShift platform when not specified.
	ArgoCDDefaultKeycloakImageForOpenShift = "registry.redhat.io/rh-sso-7/sso75-openshift-rhel8"

	// ArgoCDDefaultKeycloakVersionForOpenShift is the default Keycloak version used for the OpenShift platform when not specified.
	// Version: 7.5.1
	ArgoCDDefaultKeycloakVersionForOpenShift = "sha256:720a7e4c4926c41c1219a90daaea3b971a3d0da5a152a96fed4fb544d80f52e3"

	// Default name for Keycloak broker.
	ArgoCDDefaultKeycloakBrokerName = "keycloak-broker"

	// Default Keycloak Instance Admin user.
	ArgoCDDefaultKeycloakAdminUser = "admin"

	// Default Keycloak Instance Admin password.
	ArgoCDDefaultKeycloakAdminPassword = "admin"
)

// Keys
const (
	// ArgoCDKeycloakImageEnvName is the environment variable used to get the image
	// to used for the Keycloak container.
	ArgoCDKeycloakImageEnvName = "ARGOCD_KEYCLOAK_IMAGE"

	// ArgoCDKeycloakRealmCreatedAnnotation is the annotation used to track if keycloak realm is created
	ArgoCDKeycloakRealmCreatedAnnotation = "argocd.argoproj.io/realm-created"

	// Identifier for Keycloak.
	KeycloakIdentifier = "keycloak"
	// Identifier for TemplateInstance and Template.
	TemplateIdentifier = "rhsso"
)
