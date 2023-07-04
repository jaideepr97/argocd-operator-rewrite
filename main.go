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

package main

import (
	"flag"
	"fmt"
	"os"
	goruntime "runtime"

	// Import all Kubernetes client auth plugins (e.g. Azure, GCP, OIDC, etc.)
	// to ensure that exec-entrypoint and run can make use of them.

	"go.uber.org/zap/zapcore"
	_ "k8s.io/client-go/plugin/pkg/client/auth"

	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	argoprojiov1alpha1 "github.com/jaideepr97/argocd-operator-rewrite/api/v1alpha1"
	"github.com/jaideepr97/argocd-operator-rewrite/common"
	argocd "github.com/jaideepr97/argocd-operator-rewrite/controllers/argocd"
	"github.com/jaideepr97/argocd-operator-rewrite/pkg/monitoring"
	"github.com/jaideepr97/argocd-operator-rewrite/pkg/networking"
	"github.com/jaideepr97/argocd-operator-rewrite/pkg/workloads"
	"github.com/jaideepr97/argocd-operator-rewrite/version"

	appsv1 "github.com/openshift/api/apps/v1"
	configv1 "github.com/openshift/api/config/v1"
	oauthv1 "github.com/openshift/api/oauth/v1"
	routev1 "github.com/openshift/api/route/v1"
	templatev1 "github.com/openshift/api/template/v1"
	monitoringv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	//+kubebuilder:scaffold:imports
)

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))

	utilruntime.Must(argoprojiov1alpha1.AddToScheme(scheme))
	//+kubebuilder:scaffold:scheme
}

func printVersion() {
	setupLog.Info(fmt.Sprintf("Go Version: %s", goruntime.Version()))
	setupLog.Info(fmt.Sprintf("Go OS/Arch: %s/%s", goruntime.GOOS, goruntime.GOARCH))
	// setupLog.Info(fmt.Sprintf("Version of operator-sdk: %v", sdkVersion.Version))
	setupLog.Info(fmt.Sprintf("Version of %s-operator: %v", common.ArgoCDAppName, version.Version))
}

func main() {

	var (
		metricsAddr          string
		probeAddr            string
		enableLeaderElection bool
	)

	flag.StringVar(&metricsAddr, "metrics-bind-address", ":8080", "The address the metric endpoint binds to.")
	flag.StringVar(&probeAddr, "health-probe-bind-address", ":8081", "The address the probe endpoint binds to.")
	flag.BoolVar(&enableLeaderElection, "leader-elect", false,
		"Enable leader election for controller manager. "+
			"Enabling this will ensure there is only one active controller manager.")
	opts := zap.Options{
		Development: true,
		Level:       zapcore.Level(-1),
	}
	opts.BindFlags(flag.CommandLine)
	flag.Parse()

	ctrl.SetLogger(zap.New(zap.UseFlagOptions(&opts)))

	printVersion()

	// Inspect cluster to verify availability of extra features
	argocd.InspectCluster()

	// namespace, err := k8sutil.GetWatchNamespace()
	// if err != nil {
	// 	setupLog.Error(err, "Failed to get watch namespace, defaulting to all namespace mode")
	// }
	// setupLog.Info(fmt.Sprintf("Watching namespace \"%s\"", namespace))

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:                 scheme,
		MetricsBindAddress:     metricsAddr,
		Port:                   9443,
		HealthProbeBindAddress: probeAddr,
		LeaderElection:         enableLeaderElection,
		LeaderElectionID:       "c60f8c25.argoproj.io",
		// LeaderElectionReleaseOnCancel defines if the leader should step down voluntarily
		// when the Manager ends. This requires the binary to immediately end when the
		// Manager is stopped, otherwise, this setting is unsafe. Setting this significantly
		// speeds up voluntary leader transitions as the new leader don't have to wait
		// LeaseDuration time first.
		//
		// In the default scaffold provided, the program ends immediately after
		// the manager stops, so would be fine to enable this option. However,
		// if you are doing or is intended to do any operation such as perform cleanups
		// after the manager stops then its usage might be unsafe.
		// LeaderElectionReleaseOnCancel: true,
	})
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	setupLog.Info("Registering Components.")

	// Setup Scheme for all resources
	if err := argoprojiov1alpha1.AddToScheme(mgr.GetScheme()); err != nil {
		setupLog.Error(err, "")
		os.Exit(1)
	}

	// Setup Scheme for Prometheus if available.
	if monitoring.IsPrometheusAPIAvailable() {
		if err := monitoringv1.AddToScheme(mgr.GetScheme()); err != nil {
			setupLog.Error(err, "")
			os.Exit(1)
		}
	}

	// Setup Scheme for OpenShift Routes if available.
	if networking.IsRouteAPIAvailable() {
		if err := routev1.Install(mgr.GetScheme()); err != nil {
			setupLog.Error(err, "")
			os.Exit(1)
		}
	}

	// Set up the scheme for openshift config if available
	if argocd.IsVersionAPIAvailable() {
		if err := configv1.Install(mgr.GetScheme()); err != nil {
			setupLog.Error(err, "")
			os.Exit(1)
		}
	}

	// Setup Schemes for SSO if template instance is available.
	if workloads.IsTemplateAPIAvailable() {
		if err := templatev1.Install(mgr.GetScheme()); err != nil {
			setupLog.Error(err, "")
			os.Exit(1)
		}
		if err := appsv1.Install(mgr.GetScheme()); err != nil {
			setupLog.Error(err, "")
			os.Exit(1)
		}
		if err := oauthv1.Install(mgr.GetScheme()); err != nil {
			setupLog.Error(err, "")
			os.Exit(1)
		}
	}

	if err = (&argocd.ArgoCDReconciler{
		Client: mgr.GetClient(),
		Scheme: mgr.GetScheme(),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "ArgoCD")
		os.Exit(1)
	}
	//+kubebuilder:scaffold:builder

	if err := mgr.AddHealthzCheck("healthz", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up health check")
		os.Exit(1)
	}
	if err := mgr.AddReadyzCheck("readyz", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up ready check")
		os.Exit(1)
	}

	setupLog.Info("starting manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
	}
}
