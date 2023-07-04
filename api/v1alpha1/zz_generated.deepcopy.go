//go:build !ignore_autogenerated
// +build !ignore_autogenerated

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

// Code generated by controller-gen. DO NOT EDIT.

package v1alpha1

import (
	routev1 "github.com/openshift/api/route/v1"
	autoscalingv1 "k8s.io/api/autoscaling/v1"
	"k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ApplicationControllerSpec) DeepCopyInto(out *ApplicationControllerSpec) {
	*out = *in
	if in.Processors != nil {
		in, out := &in.Processors, &out.Processors
		*out = new(ProcessorsSpec)
		**out = **in
	}
	if in.Resources != nil {
		in, out := &in.Resources, &out.Resources
		*out = new(v1.ResourceRequirements)
		(*in).DeepCopyInto(*out)
	}
	if in.Sharding != nil {
		in, out := &in.Sharding, &out.Sharding
		*out = new(ShardingSpec)
		(*in).DeepCopyInto(*out)
	}
	if in.Env != nil {
		in, out := &in.Env, &out.Env
		*out = make([]v1.EnvVar, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ApplicationControllerSpec.
func (in *ApplicationControllerSpec) DeepCopy() *ApplicationControllerSpec {
	if in == nil {
		return nil
	}
	out := new(ApplicationControllerSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ApplicationSetSpec) DeepCopyInto(out *ApplicationSetSpec) {
	*out = *in
	if in.Env != nil {
		in, out := &in.Env, &out.Env
		*out = make([]v1.EnvVar, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.ExtraCommandArgs != nil {
		in, out := &in.ExtraCommandArgs, &out.ExtraCommandArgs
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.Resources != nil {
		in, out := &in.Resources, &out.Resources
		*out = new(v1.ResourceRequirements)
		(*in).DeepCopyInto(*out)
	}
	if in.WebhookServer != nil {
		in, out := &in.WebhookServer, &out.WebhookServer
		*out = new(WebhookServerSpec)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ApplicationSetSpec.
func (in *ApplicationSetSpec) DeepCopy() *ApplicationSetSpec {
	if in == nil {
		return nil
	}
	out := new(ApplicationSetSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ArgoCD) DeepCopyInto(out *ArgoCD) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ArgoCD.
func (in *ArgoCD) DeepCopy() *ArgoCD {
	if in == nil {
		return nil
	}
	out := new(ArgoCD)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ArgoCD) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ArgoCDList) DeepCopyInto(out *ArgoCDList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]ArgoCD, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ArgoCDList.
func (in *ArgoCDList) DeepCopy() *ArgoCDList {
	if in == nil {
		return nil
	}
	out := new(ArgoCDList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ArgoCDList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ArgoCDSpec) DeepCopyInto(out *ArgoCDSpec) {
	*out = *in
	if in.ApplicationSet != nil {
		in, out := &in.ApplicationSet, &out.ApplicationSet
		*out = new(ApplicationSetSpec)
		(*in).DeepCopyInto(*out)
	}
	if in.Controller != nil {
		in, out := &in.Controller, &out.Controller
		*out = new(ApplicationControllerSpec)
		(*in).DeepCopyInto(*out)
	}
	if in.ExtraConfig != nil {
		in, out := &in.ExtraConfig, &out.ExtraConfig
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.HA != nil {
		in, out := &in.HA, &out.HA
		*out = new(HASpec)
		(*in).DeepCopyInto(*out)
	}
	if in.Import != nil {
		in, out := &in.Import, &out.Import
		*out = new(ImportSpec)
		(*in).DeepCopyInto(*out)
	}
	if in.InitialSSHKnownHosts != nil {
		in, out := &in.InitialSSHKnownHosts, &out.InitialSSHKnownHosts
		*out = new(SSHHostsSpec)
		**out = **in
	}
	if in.KustomizeVersions != nil {
		in, out := &in.KustomizeVersions, &out.KustomizeVersions
		*out = make([]KustomizeVersionSpec, len(*in))
		copy(*out, *in)
	}
	out.Monitoring = in.Monitoring
	if in.NodePlacement != nil {
		in, out := &in.NodePlacement, &out.NodePlacement
		*out = new(NodePlacementSpec)
		(*in).DeepCopyInto(*out)
	}
	if in.Notifications != nil {
		in, out := &in.Notifications, &out.Notifications
		*out = new(NotificationsSpec)
		(*in).DeepCopyInto(*out)
	}
	if in.RBAC != nil {
		in, out := &in.RBAC, &out.RBAC
		*out = new(RBACSpec)
		(*in).DeepCopyInto(*out)
	}
	if in.Redis != nil {
		in, out := &in.Redis, &out.Redis
		*out = new(RedisSpec)
		(*in).DeepCopyInto(*out)
	}
	if in.Repo != nil {
		in, out := &in.Repo, &out.Repo
		*out = new(RepoSpec)
		(*in).DeepCopyInto(*out)
	}
	if in.ResourceHealthChecks != nil {
		in, out := &in.ResourceHealthChecks, &out.ResourceHealthChecks
		*out = make([]ResourceHealthCheck, len(*in))
		copy(*out, *in)
	}
	if in.ResourceIgnoreDifferences != nil {
		in, out := &in.ResourceIgnoreDifferences, &out.ResourceIgnoreDifferences
		*out = new(ResourceIgnoreDifference)
		(*in).DeepCopyInto(*out)
	}
	if in.ResourceActions != nil {
		in, out := &in.ResourceActions, &out.ResourceActions
		*out = make([]ResourceAction, len(*in))
		copy(*out, *in)
	}
	if in.Server != nil {
		in, out := &in.Server, &out.Server
		*out = new(ServerSpec)
		(*in).DeepCopyInto(*out)
	}
	if in.SourceNamespaces != nil {
		in, out := &in.SourceNamespaces, &out.SourceNamespaces
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.SSO != nil {
		in, out := &in.SSO, &out.SSO
		*out = new(SSOSpec)
		(*in).DeepCopyInto(*out)
	}
	if in.TLS != nil {
		in, out := &in.TLS, &out.TLS
		*out = new(TLSSpec)
		(*in).DeepCopyInto(*out)
	}
	if in.Banner != nil {
		in, out := &in.Banner, &out.Banner
		*out = new(Banner)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ArgoCDSpec.
func (in *ArgoCDSpec) DeepCopy() *ArgoCDSpec {
	if in == nil {
		return nil
	}
	out := new(ArgoCDSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ArgoCDStatus) DeepCopyInto(out *ArgoCDStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ArgoCDStatus.
func (in *ArgoCDStatus) DeepCopy() *ArgoCDStatus {
	if in == nil {
		return nil
	}
	out := new(ArgoCDStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AutoscalingSpec) DeepCopyInto(out *AutoscalingSpec) {
	*out = *in
	if in.HPA != nil {
		in, out := &in.HPA, &out.HPA
		*out = new(autoscalingv1.HorizontalPodAutoscalerSpec)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AutoscalingSpec.
func (in *AutoscalingSpec) DeepCopy() *AutoscalingSpec {
	if in == nil {
		return nil
	}
	out := new(AutoscalingSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Banner) DeepCopyInto(out *Banner) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Banner.
func (in *Banner) DeepCopy() *Banner {
	if in == nil {
		return nil
	}
	out := new(Banner)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CASpec) DeepCopyInto(out *CASpec) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CASpec.
func (in *CASpec) DeepCopy() *CASpec {
	if in == nil {
		return nil
	}
	out := new(CASpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DexSpec) DeepCopyInto(out *DexSpec) {
	*out = *in
	if in.Groups != nil {
		in, out := &in.Groups, &out.Groups
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.Resources != nil {
		in, out := &in.Resources, &out.Resources
		*out = new(v1.ResourceRequirements)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DexSpec.
func (in *DexSpec) DeepCopy() *DexSpec {
	if in == nil {
		return nil
	}
	out := new(DexSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DynamicScalingSpec) DeepCopyInto(out *DynamicScalingSpec) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DynamicScalingSpec.
func (in *DynamicScalingSpec) DeepCopy() *DynamicScalingSpec {
	if in == nil {
		return nil
	}
	out := new(DynamicScalingSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GRPCSpec) DeepCopyInto(out *GRPCSpec) {
	*out = *in
	if in.Ingress != nil {
		in, out := &in.Ingress, &out.Ingress
		*out = new(IngressSpec)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GRPCSpec.
func (in *GRPCSpec) DeepCopy() *GRPCSpec {
	if in == nil {
		return nil
	}
	out := new(GRPCSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *HASpec) DeepCopyInto(out *HASpec) {
	*out = *in
	if in.Resources != nil {
		in, out := &in.Resources, &out.Resources
		*out = new(v1.ResourceRequirements)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new HASpec.
func (in *HASpec) DeepCopy() *HASpec {
	if in == nil {
		return nil
	}
	out := new(HASpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *IgnoreDifferenceCustomization) DeepCopyInto(out *IgnoreDifferenceCustomization) {
	*out = *in
	if in.JqPathExpressions != nil {
		in, out := &in.JqPathExpressions, &out.JqPathExpressions
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.JsonPointers != nil {
		in, out := &in.JsonPointers, &out.JsonPointers
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.ManagedFieldsManagers != nil {
		in, out := &in.ManagedFieldsManagers, &out.ManagedFieldsManagers
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new IgnoreDifferenceCustomization.
func (in *IgnoreDifferenceCustomization) DeepCopy() *IgnoreDifferenceCustomization {
	if in == nil {
		return nil
	}
	out := new(IgnoreDifferenceCustomization)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ImportSpec) DeepCopyInto(out *ImportSpec) {
	*out = *in
	if in.Namespace != nil {
		in, out := &in.Namespace, &out.Namespace
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ImportSpec.
func (in *ImportSpec) DeepCopy() *ImportSpec {
	if in == nil {
		return nil
	}
	out := new(ImportSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *IngressSpec) DeepCopyInto(out *IngressSpec) {
	*out = *in
	if in.Annotations != nil {
		in, out := &in.Annotations, &out.Annotations
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.IngressClassName != nil {
		in, out := &in.IngressClassName, &out.IngressClassName
		*out = new(string)
		**out = **in
	}
	if in.TLS != nil {
		in, out := &in.TLS, &out.TLS
		*out = make([]networkingv1.IngressTLS, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new IngressSpec.
func (in *IngressSpec) DeepCopy() *IngressSpec {
	if in == nil {
		return nil
	}
	out := new(IngressSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *KeycloakSpec) DeepCopyInto(out *KeycloakSpec) {
	*out = *in
	if in.Resources != nil {
		in, out := &in.Resources, &out.Resources
		*out = new(v1.ResourceRequirements)
		(*in).DeepCopyInto(*out)
	}
	if in.VerifyTLS != nil {
		in, out := &in.VerifyTLS, &out.VerifyTLS
		*out = new(bool)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new KeycloakSpec.
func (in *KeycloakSpec) DeepCopy() *KeycloakSpec {
	if in == nil {
		return nil
	}
	out := new(KeycloakSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *KustomizeVersionSpec) DeepCopyInto(out *KustomizeVersionSpec) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new KustomizeVersionSpec.
func (in *KustomizeVersionSpec) DeepCopy() *KustomizeVersionSpec {
	if in == nil {
		return nil
	}
	out := new(KustomizeVersionSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MonitoringSpec) DeepCopyInto(out *MonitoringSpec) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MonitoringSpec.
func (in *MonitoringSpec) DeepCopy() *MonitoringSpec {
	if in == nil {
		return nil
	}
	out := new(MonitoringSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NodePlacementSpec) DeepCopyInto(out *NodePlacementSpec) {
	*out = *in
	if in.NodeSelector != nil {
		in, out := &in.NodeSelector, &out.NodeSelector
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.Tolerations != nil {
		in, out := &in.Tolerations, &out.Tolerations
		*out = make([]v1.Toleration, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NodePlacementSpec.
func (in *NodePlacementSpec) DeepCopy() *NodePlacementSpec {
	if in == nil {
		return nil
	}
	out := new(NodePlacementSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NotificationsSpec) DeepCopyInto(out *NotificationsSpec) {
	*out = *in
	if in.Replicas != nil {
		in, out := &in.Replicas, &out.Replicas
		*out = new(int32)
		**out = **in
	}
	if in.Env != nil {
		in, out := &in.Env, &out.Env
		*out = make([]v1.EnvVar, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.Resources != nil {
		in, out := &in.Resources, &out.Resources
		*out = new(v1.ResourceRequirements)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NotificationsSpec.
func (in *NotificationsSpec) DeepCopy() *NotificationsSpec {
	if in == nil {
		return nil
	}
	out := new(NotificationsSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ProcessorsSpec) DeepCopyInto(out *ProcessorsSpec) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ProcessorsSpec.
func (in *ProcessorsSpec) DeepCopy() *ProcessorsSpec {
	if in == nil {
		return nil
	}
	out := new(ProcessorsSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RBACSpec) DeepCopyInto(out *RBACSpec) {
	*out = *in
	if in.DefaultPolicy != nil {
		in, out := &in.DefaultPolicy, &out.DefaultPolicy
		*out = new(string)
		**out = **in
	}
	if in.Policy != nil {
		in, out := &in.Policy, &out.Policy
		*out = new(string)
		**out = **in
	}
	if in.Scopes != nil {
		in, out := &in.Scopes, &out.Scopes
		*out = new(string)
		**out = **in
	}
	if in.PolicyMatcherMode != nil {
		in, out := &in.PolicyMatcherMode, &out.PolicyMatcherMode
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RBACSpec.
func (in *RBACSpec) DeepCopy() *RBACSpec {
	if in == nil {
		return nil
	}
	out := new(RBACSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RedisSpec) DeepCopyInto(out *RedisSpec) {
	*out = *in
	if in.Resources != nil {
		in, out := &in.Resources, &out.Resources
		*out = new(v1.ResourceRequirements)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RedisSpec.
func (in *RedisSpec) DeepCopy() *RedisSpec {
	if in == nil {
		return nil
	}
	out := new(RedisSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RepoSpec) DeepCopyInto(out *RepoSpec) {
	*out = *in
	if in.ExtraRepoCommandArgs != nil {
		in, out := &in.ExtraRepoCommandArgs, &out.ExtraRepoCommandArgs
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.Replicas != nil {
		in, out := &in.Replicas, &out.Replicas
		*out = new(int32)
		**out = **in
	}
	if in.Resources != nil {
		in, out := &in.Resources, &out.Resources
		*out = new(v1.ResourceRequirements)
		(*in).DeepCopyInto(*out)
	}
	if in.ExecTimeout != nil {
		in, out := &in.ExecTimeout, &out.ExecTimeout
		*out = new(int)
		**out = **in
	}
	if in.Env != nil {
		in, out := &in.Env, &out.Env
		*out = make([]v1.EnvVar, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.Volumes != nil {
		in, out := &in.Volumes, &out.Volumes
		*out = make([]v1.Volume, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.VolumeMounts != nil {
		in, out := &in.VolumeMounts, &out.VolumeMounts
		*out = make([]v1.VolumeMount, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.InitContainers != nil {
		in, out := &in.InitContainers, &out.InitContainers
		*out = make([]v1.Container, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.SidecarContainers != nil {
		in, out := &in.SidecarContainers, &out.SidecarContainers
		*out = make([]v1.Container, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RepoSpec.
func (in *RepoSpec) DeepCopy() *RepoSpec {
	if in == nil {
		return nil
	}
	out := new(RepoSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ResourceAction) DeepCopyInto(out *ResourceAction) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ResourceAction.
func (in *ResourceAction) DeepCopy() *ResourceAction {
	if in == nil {
		return nil
	}
	out := new(ResourceAction)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ResourceHealthCheck) DeepCopyInto(out *ResourceHealthCheck) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ResourceHealthCheck.
func (in *ResourceHealthCheck) DeepCopy() *ResourceHealthCheck {
	if in == nil {
		return nil
	}
	out := new(ResourceHealthCheck)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ResourceIdentifiers) DeepCopyInto(out *ResourceIdentifiers) {
	*out = *in
	in.Customization.DeepCopyInto(&out.Customization)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ResourceIdentifiers.
func (in *ResourceIdentifiers) DeepCopy() *ResourceIdentifiers {
	if in == nil {
		return nil
	}
	out := new(ResourceIdentifiers)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ResourceIgnoreDifference) DeepCopyInto(out *ResourceIgnoreDifference) {
	*out = *in
	if in.All != nil {
		in, out := &in.All, &out.All
		*out = new(IgnoreDifferenceCustomization)
		(*in).DeepCopyInto(*out)
	}
	if in.ResourceIdentifiers != nil {
		in, out := &in.ResourceIdentifiers, &out.ResourceIdentifiers
		*out = make([]ResourceIdentifiers, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ResourceIgnoreDifference.
func (in *ResourceIgnoreDifference) DeepCopy() *ResourceIgnoreDifference {
	if in == nil {
		return nil
	}
	out := new(ResourceIgnoreDifference)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RouteSpec) DeepCopyInto(out *RouteSpec) {
	*out = *in
	if in.Annotations != nil {
		in, out := &in.Annotations, &out.Annotations
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.Labels != nil {
		in, out := &in.Labels, &out.Labels
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.TLS != nil {
		in, out := &in.TLS, &out.TLS
		*out = new(routev1.TLSConfig)
		**out = **in
	}
	if in.WildcardPolicy != nil {
		in, out := &in.WildcardPolicy, &out.WildcardPolicy
		*out = new(routev1.WildcardPolicyType)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RouteSpec.
func (in *RouteSpec) DeepCopy() *RouteSpec {
	if in == nil {
		return nil
	}
	out := new(RouteSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SSHHostsSpec) DeepCopyInto(out *SSHHostsSpec) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SSHHostsSpec.
func (in *SSHHostsSpec) DeepCopy() *SSHHostsSpec {
	if in == nil {
		return nil
	}
	out := new(SSHHostsSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SSOSpec) DeepCopyInto(out *SSOSpec) {
	*out = *in
	if in.Resources != nil {
		in, out := &in.Resources, &out.Resources
		*out = new(v1.ResourceRequirements)
		(*in).DeepCopyInto(*out)
	}
	if in.VerifyTLS != nil {
		in, out := &in.VerifyTLS, &out.VerifyTLS
		*out = new(bool)
		**out = **in
	}
	if in.Dex != nil {
		in, out := &in.Dex, &out.Dex
		*out = new(DexSpec)
		(*in).DeepCopyInto(*out)
	}
	if in.Keycloak != nil {
		in, out := &in.Keycloak, &out.Keycloak
		*out = new(KeycloakSpec)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SSOSpec.
func (in *SSOSpec) DeepCopy() *SSOSpec {
	if in == nil {
		return nil
	}
	out := new(SSOSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ServerSpec) DeepCopyInto(out *ServerSpec) {
	*out = *in
	if in.Autoscale != nil {
		in, out := &in.Autoscale, &out.Autoscale
		*out = new(AutoscalingSpec)
		(*in).DeepCopyInto(*out)
	}
	if in.GRPC != nil {
		in, out := &in.GRPC, &out.GRPC
		*out = new(GRPCSpec)
		(*in).DeepCopyInto(*out)
	}
	if in.Ingress != nil {
		in, out := &in.Ingress, &out.Ingress
		*out = new(IngressSpec)
		(*in).DeepCopyInto(*out)
	}
	if in.Replicas != nil {
		in, out := &in.Replicas, &out.Replicas
		*out = new(int32)
		**out = **in
	}
	if in.Resources != nil {
		in, out := &in.Resources, &out.Resources
		*out = new(v1.ResourceRequirements)
		(*in).DeepCopyInto(*out)
	}
	if in.Route != nil {
		in, out := &in.Route, &out.Route
		*out = new(RouteSpec)
		(*in).DeepCopyInto(*out)
	}
	if in.Service != nil {
		in, out := &in.Service, &out.Service
		*out = new(ServiceSpec)
		**out = **in
	}
	if in.Env != nil {
		in, out := &in.Env, &out.Env
		*out = make([]v1.EnvVar, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.ExtraCommandArgs != nil {
		in, out := &in.ExtraCommandArgs, &out.ExtraCommandArgs
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ServerSpec.
func (in *ServerSpec) DeepCopy() *ServerSpec {
	if in == nil {
		return nil
	}
	out := new(ServerSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ServiceSpec) DeepCopyInto(out *ServiceSpec) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ServiceSpec.
func (in *ServiceSpec) DeepCopy() *ServiceSpec {
	if in == nil {
		return nil
	}
	out := new(ServiceSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ShardingSpec) DeepCopyInto(out *ShardingSpec) {
	*out = *in
	if in.DynamicScaling != nil {
		in, out := &in.DynamicScaling, &out.DynamicScaling
		*out = new(DynamicScalingSpec)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ShardingSpec.
func (in *ShardingSpec) DeepCopy() *ShardingSpec {
	if in == nil {
		return nil
	}
	out := new(ShardingSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TLSSpec) DeepCopyInto(out *TLSSpec) {
	*out = *in
	if in.CA != nil {
		in, out := &in.CA, &out.CA
		*out = new(CASpec)
		**out = **in
	}
	if in.InitialCerts != nil {
		in, out := &in.InitialCerts, &out.InitialCerts
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TLSSpec.
func (in *TLSSpec) DeepCopy() *TLSSpec {
	if in == nil {
		return nil
	}
	out := new(TLSSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *WebhookServerSpec) DeepCopyInto(out *WebhookServerSpec) {
	*out = *in
	if in.Ingress != nil {
		in, out := &in.Ingress, &out.Ingress
		*out = new(IngressSpec)
		(*in).DeepCopyInto(*out)
	}
	if in.Route != nil {
		in, out := &in.Route, &out.Route
		*out = new(RouteSpec)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new WebhookServerSpec.
func (in *WebhookServerSpec) DeepCopy() *WebhookServerSpec {
	if in == nil {
		return nil
	}
	out := new(WebhookServerSpec)
	in.DeepCopyInto(out)
	return out
}
