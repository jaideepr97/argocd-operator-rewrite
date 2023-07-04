package sso

import (
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type SSOReconciler struct {
	Client *client.Client
	Scheme *runtime.Scheme
}
