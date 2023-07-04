package networking

// Keys
const (
	// ArgoCDKeyIngressBackendProtocol is the backend-protocol key for labels.
	ArgoCDKeyIngressBackendProtocol = "nginx.ingress.kubernetes.io/backend-protocol"

	// ArgoCDKeyIngressClass is the ingress class key for labels.
	ArgoCDKeyIngressClass = "kubernetes.io/ingress.class"

	// ArgoCDKeyIngressSSLRedirect is the ssl force-redirect key for labels.
	ArgoCDKeyIngressSSLRedirect = "nginx.ingress.kubernetes.io/force-ssl-redirect"

	// ArgoCDKeyIngressSSLPassthrough is the ssl passthrough key for labels.
	ArgoCDKeyIngressSSLPassthrough = "nginx.ingress.kubernetes.io/ssl-passthrough"
)
