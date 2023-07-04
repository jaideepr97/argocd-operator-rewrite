package appcontroller

import (
	"fmt"
	"strings"

	"github.com/jaideepr97/argocd-operator-rewrite/api/v1alpha1"
	"github.com/jaideepr97/argocd-operator-rewrite/controllers/argocd/redis"
	"github.com/jaideepr97/argocd-operator-rewrite/controllers/argocd/reposerver"
	"github.com/jaideepr97/argocd-operator-rewrite/controllers/argocd/server"
	"github.com/jaideepr97/argocd-operator-rewrite/pkg/argoutil"
	corev1 "k8s.io/api/core/v1"
)

// getApplicationControllerResources will return the ResourceRequirements for the Argo CD application controller container.
func getApplicationControllerResources(cr *v1alpha1.ArgoCD) corev1.ResourceRequirements {
	resources := corev1.ResourceRequirements{}

	// Allow override of resource requirements from CR
	if cr.Spec.Controller.Resources != nil {
		resources = *cr.Spec.Controller.Resources
	}

	return resources
}

// getApplicationControllerParellismLimit returns the parallelism limit for the application controller
func getApplicationControllerParellismLimit(cr *v1alpha1.ArgoCD) int32 {
	pl := ArgoCDDefaultControllerParallelismLimit
	if cr.Spec.Controller.ParallelismLimit > 0 {
		pl = cr.Spec.Controller.ParallelismLimit
	}
	return pl
}

// getApplicationControllerCommands will return the commands for the ArgoCD Application Controller component.
func getApplicationControllerCommands(cr *v1alpha1.ArgoCD, useTLSForRedis bool) []string {
	cmd := []string{
		"argocd-application-controller",
		"--operation-processors", fmt.Sprint(server.GetServerOperationProcessors(cr)),
		"--redis", redis.GetRedisServerAddress(cr),
	}

	if useTLSForRedis {
		cmd = append(cmd, "--redis-use-tls")
		if redis.IsRedisTLSVerificationDisabled(cr) {
			cmd = append(cmd, "--redis-insecure-skip-tls-verify")
		} else {
			cmd = append(cmd, "--redis-ca-certificate", "/app/config/controller/tls/redis/tls.crt")
		}
	}

	cmd = append(cmd, "--repo-server", reposerver.GetRepoServerAddress(cr))
	cmd = append(cmd, "--status-processors", fmt.Sprint(server.GetArgoServerStatusProcessors(cr)))
	cmd = append(cmd, "--kubectl-parallelism-limit", fmt.Sprint(getApplicationControllerParellismLimit(cr)))

	if cr.Spec.SourceNamespaces != nil && len(cr.Spec.SourceNamespaces) > 0 {
		cmd = append(cmd, "--application-namespaces", fmt.Sprint(strings.Join(cr.Spec.SourceNamespaces, ",")))
	}

	cmd = append(cmd, "--loglevel")
	cmd = append(cmd, argoutil.GetLogLevel(cr.Spec.Controller.LogLevel))

	cmd = append(cmd, "--logformat")
	cmd = append(cmd, argoutil.GetLogFormat(cr.Spec.Controller.LogFormat))

	return cmd
}

func getSourceNamespaceRBACName(instanceName, instaceNamespace string) string {
	return argoutil.GenerateUniqueResourceName(instanceName, instaceNamespace, ArgoCDApplicationControllerComponent)
}
