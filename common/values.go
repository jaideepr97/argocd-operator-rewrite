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

import "time"

const (
	// ArgoCDAppName is the application name for labels.
	ArgoCDAppName = "argocd"

	// ArgoCDCASuffix is the name suffix for ArgoCD CA resources.
	ArgoCDCASuffix = "ca"

	// ArgoCDConfigMapName is the upstream hard-coded ArgoCD ConfigMap name.
	ArgoCDConfigMapName = "argocd-cm"

	// ArgoCDGPGKeysConfigMapName is the upstream hard-coded ArgoCD gpg-keys ConfigMap name.
	ArgoCDGPGKeysConfigMapName = "argocd-gpg-keys-cm"

	// ArgoCDDuration365Days is a duration representing 365 days.
	ArgoCDDuration365Days = time.Hour * 24 * 365

	// ArgoCDExportName is the export name for labels.
	ArgoCDExportName = "argocd.export"

	// ArgoCDExportStorageBackendAWS is the value for the AWS storage backend.
	ArgoCDExportStorageBackendAWS = "aws"

	// ArgoCDExportStorageBackendAzure is the value for the Azure storage backend.
	ArgoCDExportStorageBackendAzure = "azure"

	// ArgoCDExportStorageBackendGCP is the value for the GCP storage backend.
	ArgoCDExportStorageBackendGCP = "gcp"

	// ArgoCDExportStorageBackendLocal is the value for the local storage backend.
	ArgoCDExportStorageBackendLocal = "local"

	// ArgoCDKnownHostsConfigMapName is the upstream hard-coded SSH known hosts data ConfigMap name.
	ArgoCDKnownHostsConfigMapName = "argocd-ssh-known-hosts-cm"

	// ArgoCDRBACConfigMapName is the upstream hard-coded RBAC ConfigMap name.
	ArgoCDRBACConfigMapName = "argocd-rbac-cm"

	// ArgoCDSecretName is the upstream hard-coded ArgoCD Secret name.
	ArgoCDSecretName = "argocd-secret"

	// ArgoCDStatusCompleted is the completed status value.
	ArgoCDStatusCompleted = "Completed"

	// ArgoCDTLSCertsConfigMapName is the upstream hard-coded TLS certificate data ConfigMap name.
	ArgoCDTLSCertsConfigMapName = "argocd-tls-certs-cm"

	// ArgoCDRBACTypeAppManagement is the value used when an rbac resource is targeted for applicaton management
	ArgoCDRBACTypeAppManagement = "app-management"

	// ArgoCDRBACTypeAppManagement is the value used when an rbac resource is targeted for resource management
	ArgoCDRBACTypeResourceMananagement = "resource-management"

	// ArgoCDSecretTypeCluster is the value used for cluster secret labels
	ArgoCDSecretTypeCluster = "cluster"
)
