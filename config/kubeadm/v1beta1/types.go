/*
Copyright 2018 The Kubernetes Authors.

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

package v1beta1

import (
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// DEPRECATED - This group version of InitConfiguration is deprecated by apis/kubeadm/v1beta2/InitConfiguration.
// InitConfiguration contains a list of elements that is specific "kubeadm init"-only runtime
// information.
type InitConfiguration struct {
	metav1.TypeMeta `json:",inline"`

	// ClusterConfiguration holds the cluster-wide information, and embeds that struct
	// (which can be (un)marshalled separately as well). When InitConfiguration is
	// marshalled to bytes in the external version, this information IS NOT preserved
	// (which can be seen from the `json:"-"` tag).  This is due to that when
	// InitConfiguration is (un)marshalled, it turns into two YAML documents, one for the
	// InitConfiguration and ClusterConfiguration. Hence, the information must not be
	// duplicated, and is therefore omitted here.
	ClusterConfiguration `json:"-"`

	// `kubeadm init`-only information. These fields are solely used the first time `kubeadm init` runs.
	// After that, the information in the fields IS NOT uploaded to the `kubeadm-config` ConfigMap
	// that is used by `kubeadm upgrade` for instance. These fields must be omitempty.

	// `bootstrapTokens` is respected at `kubeadm init` time and describes a set of
	// Bootstrap Tokens to create. This information IS NOT uploaded to the kubeadm cluster
	// configmap, partly because of its sensitive nature
	BootstrapTokens []BootstrapToken `json:"bootstrapTokens,omitempty"`

	// Fields that relate to registering the new control-plane node to the cluster
	NodeRegistration NodeRegistrationOptions `json:"nodeRegistration,omitempty"`

	// `ocalAPIEndpoint` represents the endpoint of the API server instance that's deployedon this control plane node.
	// In HA setups, this differs from ClusterConfiguration.ControlPlaneEndpoint in the sense that ControlPlaneEndpoint
	// is the global endpoint for the cluster, which then loadbalances the requests to each individual API server. This
	// configuration object lets you customize what IP/DNS name and port the local API server advertises it's accessible
	// on. By default, kubeadm tries to auto-detect the IP of the default interface and use that, but in case that
	// process fails you may set the desired value here.
	LocalAPIEndpoint APIEndpoint `json:"localAPIEndpoint,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// DEPRECATED - This group version of ClusterConfiguration is deprecated by apis/kubeadm/v1beta2/ClusterConfiguration.
// ClusterConfiguration contains cluster-wide configuration for a kubeadm cluster
type ClusterConfiguration struct {
	metav1.TypeMeta `json:",inline"`

	// The configuration for etcd.
	Etcd Etcd `json:"etcd"`

	// `networking` holds configuration for the networking topology of the cluster.
	Networking Networking `json:"networking"`

	// `kubernetesVersion` is the target version of the control plane.
	KubernetesVersion string `json:"kubernetesVersion"`

	// `controlPlaneEndpoint` sets a stable IP address or DNS name for the control plane; it
	// can be a valid IP address or a RFC-1123 DNS subdomain, both with optional TCP port.
	// In case the ControlPlaneEndpoint is not specified, the AdvertiseAddress + BindPort
	// are used; in case the ControlPlaneEndpoint is specified but without a TCP port,
	// the BindPort is used.
	//
	// Possible usages are:
	//
	// - e.g. In a cluster with more than one control plane instances, this field should be
	// assigned the address of the external load balancer in front of the
	// control plane instances.
	// - e.g.  in environments with enforced node recycling, the ControlPlaneEndpoint
	// could be used for assigning a stable DNS to the control plane.
	ControlPlaneEndpoint string `json:"controlPlaneEndpoint"`

	// Extra settings for the API server control plane component
	APIServer APIServer `json:"apiServer,omitempty"`

	// Extra settings for the controller manager control plane component.
	ControllerManager ControlPlaneComponent `json:"controllerManager,omitempty"`

	// Extra settings for the scheduler control plane component.
	Scheduler ControlPlaneComponent `json:"scheduler,omitempty"`

	// The options for the DNS add-on installed in the cluster.
	DNS DNS `json:"dns"`

	// Where to store or look for all required certificates.
	CertificatesDir string `json:"certificatesDir"`

	// `imageRepository` sets the container registry to pull images from.
	// If empty, `k8s.gcr.io` will be used by default; in case of kubernetes
	// version is a CI build (kubernetes version starts with `ci/` or `ci-cross/`)
	// `gcr.io/k8s-staging-ci-images` will be used as a default for control plane
	// components and for kube-proxy, while `k8s.gcr.io` will be used for all
	// the other images.
	ImageRepository string `json:"imageRepository"`

	// `useHyperKubeImage` controls if hyperkube should be used for Kubernetes
	// components instead of their respective separate images.
	// *DEPRECATED*: As hyperkube is itself deprecated, this fields is too. It will
	// be removed in future kubeadm config versions, kubeadm will print multiple
	// warnings when set to true, and at some point it may become ignored.
	UseHyperKubeImage bool `json:"useHyperKubeImage,omitempty"`

	// Feature gates enabled by the user.
	FeatureGates map[string]bool `json:"featureGates,omitempty"`

	// The cluster name.
	ClusterName string `json:"clusterName,omitempty"`
}

// ControlPlaneComponent holds settings common to control plane component of the cluster
type ControlPlaneComponent struct {
	// An extra set of flags to pass to the control plane component.
	ExtraArgs map[string]string `json:"extraArgs,omitempty"`

	// An extra set of host volumes, mounted to the control plane component.
	ExtraVolumes []HostPathMount `json:"extraVolumes,omitempty"`
}

// APIServer holds settings necessary for API server deployments in the cluster
type APIServer struct {
	ControlPlaneComponent `json:",inline"`

	// `certSANs` sets extra Subject Alternative Names for the API Server signing cert.
	CertSANs []string `json:"certSANs,omitempty"`

	// `timeoutForControlPlane` controls the timeout that we use for API server to appear.
	TimeoutForControlPlane *metav1.Duration `json:"timeoutForControlPlane,omitempty"`
}

// DNSAddOnType defines string identifying DNS add-on types
type DNSAddOnType string

const (
	// CoreDNS add-on type
	CoreDNS DNSAddOnType = "CoreDNS"
)

// DNS defines the DNS addon that should be used in the cluster
type DNS struct {
	// `type` defines the DNS add-on to be used.
	Type DNSAddOnType `json:"type"`

	// `imageMeta` allows to customize the image used for the DNS component.
	ImageMeta `json:",inline"`
}

// ImageMeta allows to customize the image used for components that are not
// originated from the Kubernetes/Kubernetes release process
type ImageMeta struct {
	// `imageRepository` sets the container registry to pull images from.
	// If not set, the `imageRepository` defined in ClusterConfiguration will
	// be used instead.
	ImageRepository string `json:"imageRepository,omitempty"`

	// `imageTag` allows to specify a tag for the image.
	// In case this value is set, kubeadm does not change automatically the version
	// of the above components during upgrades.
	ImageTag string `json:"imageTag,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ClusterStatus contains the cluster status. The ClusterStatus will be stored
// in the "kubeadm-config" ConfigMap in the cluster, and then updated by kubeadm
// when additional control plane instance joins or leaves the cluster.
type ClusterStatus struct {
	metav1.TypeMeta `json:",inline"`

	// `apiEndpoints` currently available in the cluster, one for each control
	// plane/api server instance. The key of the map is the IP of the host's default interface.
	APIEndpoints map[string]APIEndpoint `json:"apiEndpoints"`
}

// APIEndpoint struct contains elements of API server instance deployed on a node.
type APIEndpoint struct {
	// `advertiseAddress` sets the IP address for the API server to advertise.
	AdvertiseAddress string `json:"advertiseAddress"`

	// The secure port for the API Server to bind to.
	// Defaults to 6443.
	BindPort int32 `json:"bindPort"`
}

// NodeRegistrationOptions holds fields that relate to registering a new control-plane
// or node to the cluster, either via "kubeadm init" or "kubeadm join"
type NodeRegistrationOptions struct {

	// `name` is the `.Metadata.Name` field of the Node API object that will be
	// created in this `kubeadm init` or `kubeadm join` operation.
	// This field is also used in the CommonName field of the kubelet's client
	// certificate to the API server.
	// Defaults to the hostname of the node if not provided.
	Name string `json:"name,omitempty"`

	// `criSocket` is used to retrieve container runtime info. This information will
	// be annotated to the Node API object, for later re-use
	CRISocket string `json:"criSocket,omitempty"`

	// `taints` specifies the taints the Node API object should be registered with.
	// If this field is unset, i.e. nil, in the `kubeadm init` process it will be
	// defaulted to `["node-role.kubernetes.io/master"=""]`. If you don't want to
	// taint your control-plane node, set this field to an empty list, i.e. `[]`
	// in the YAML file. This field is solely used for Node registration.
	Taints []v1.Taint `json:"taints,omitempty"`

	// `kubeletExtraArgs` passes through extra arguments to the kubelet. The arguments
	// here are passed to the kubelet command line via the environment file kubeadm
	// writes at runtime for the kubelet to source. This overrides the generic
	// base-level configuration in the "kubelet-config-1.X" ConfigMap flags have
	// higher priority when parsing.
	// These values are local and specific to the node kubeadm is executing on.
	KubeletExtraArgs map[string]string `json:"kubeletExtraArgs,omitempty"`
}

// Networking contains elements describing cluster's networking configuration
type Networking struct {
	// `serviceSubnet` is the subnet used by k8s services. Defaults to "10.96.0.0/12".
	ServiceSubnet string `json:"serviceSubnet"`
	// `podSubnet` is the subnet used by Pods.
	PodSubnet string `json:"podSubnet"`
	// `dnsDomain` is the dns domain used by k8s services. Defaults to "cluster.local".
	DNSDomain string `json:"dnsDomain"`
}

// BootstrapToken describes one bootstrap token, stored as a Secret in the cluster
type BootstrapToken struct {
	// `token` is used for establishing bidirectional trust between nodes and control-planes.
	// Used for joining nodes in the cluster.
	Token *BootstrapTokenString `json:"token"`
	// `description` sets a human-friendly message why this token exists and what it's used
	// for, so other administrators can know its purpose.
	Description string `json:"description,omitempty"`
	// `ttl` defines the time to live for this token. Defaults to "24h".
	// `expires` and `ttl` are mutually exclusive.
	TTL *metav1.Duration `json:"ttl,omitempty"`
	// `expires` specifies the timestamp when this token expires. Defaults to being set
	// dynamically at runtime based on the `ttl`. `expires` and `ttl` are mutually exclusive.
	Expires *metav1.Time `json:"expires,omitempty"`
	// `usages` describes the ways in which this token can be used. Can by default be used
	// for establishing bidirectional trust, but that can be changed here.
	Usages []string `json:"usages,omitempty"`
	// `groups` specifies the extra groups that this token will authenticate as when/if
	// used for authentication.
	Groups []string `json:"groups,omitempty"`
}

// Etcd contains elements describing Etcd configuration.
type Etcd struct {

	// `local` provides configuration knobs for configuring the local etcd instance
	// `local` and `external` are mutually exclusive.
	Local *LocalEtcd `json:"local,omitempty"`

	// `external` defines how to connect to an external etcd cluster.
	// `local` and `external` are mutually exclusive.
	External *ExternalEtcd `json:"external,omitempty"`
}

// LocalEtcd describes that kubeadm should run an etcd cluster locally
type LocalEtcd struct {
	// This allows to customize the container used for etcd.
	ImageMeta `json:",inline"`

	// The directory etcd will place its data.
	// Defaults to "/var/lib/etcd".
	DataDir string `json:"dataDir"`

	// Extra arguments provided to the etcd binary when run inside a static Pod.
	ExtraArgs map[string]string `json:"extraArgs,omitempty"`

	// Extra Subject Alternative Names for the etcd server signing cert.
	ServerCertSANs []string `json:"serverCertSANs,omitempty"`
	// Extra Subject Alternative Names for the etcd peer signing cert.
	PeerCertSANs []string `json:"peerCertSANs,omitempty"`
}

// ExternalEtcd describes an external etcd cluster.
// Kubeadm has no knowledge of where certificate files live and they must be supplied.
type ExternalEtcd struct {
	// Endpoints of etcd members. Required for `ExternalEtcd`.
	Endpoints []string `json:"endpoints"`

	// A SSL Certificate Authority file used to secure etcd communication.
	// Required if using a TLS connection.
	CAFile string `json:"caFile"`

	// A SSL certification file used to secure etcd communication.
	// Required if using a TLS connection.
	CertFile string `json:"certFile"`

	// A SSL key file used to secure etcd communication.
	// Required if using a TLS connection.
	KeyFile string `json:"keyFile"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// DEPRECATED - This group version of JoinConfiguration is deprecated by v1beta2.JoinConfiguration.
// JoinConfiguration contains elements describing a particular node.
type JoinConfiguration struct {
	metav1.TypeMeta `json:",inline"`

	// Fields that relate to registering the new control-plane node to the cluster.
	NodeRegistration NodeRegistrationOptions `json:"nodeRegistration"`

	// The path to the SSL certificate authority used to
	// secure comunications between node and control-plane.
	// Defaults to "/etc/kubernetes/pki/ca.crt".
	CACertPath string `json:"caCertPath"`

	// The options for the kubelet to use during the TLS Bootstrap process.
	Discovery Discovery `json:"discovery"`

	// The additional control plane instance to be deployed on the joining node.
	// If nil, no additional control plane instance will be deployed.
	ControlPlane *JoinControlPlane `json:"controlPlane,omitempty"`
}

// JoinControlPlane contains elements describing an additional control plane instance to be deployed on the joining node.
type JoinControlPlane struct {
	// The endpoint of the API server instance to be deployed on this node.
	LocalAPIEndpoint APIEndpoint `json:"localAPIEndpoint,omitempty"`
}

// Discovery specifies the options for the kubelet to use during the TLS Bootstrap process
type Discovery struct {
	// BootstrapToken is used to set the options for bootstrap token based discovery
	// BootstrapToken and File are mutually exclusive
	BootstrapToken *BootstrapTokenDiscovery `json:"bootstrapToken,omitempty"`

	// A file or URL to a kubeconfig file from which to load cluster information.
	// `bootstrapToken` and `file` are mutually exclusive
	File *FileDiscovery `json:"file,omitempty"`

	// A token used for TLS bootstrapping.
	// If `bootstrapToken` is set, this field is defaulted to `bootstrapToken.token`, but can be overridden.
	// If `file` is set, this field **must be set** in case the KubeConfigFile does
	// not contain any other authentication information
	TLSBootstrapToken string `json:"tlsBootstrapToken"`

	// Timeout modifies the discovery timeout
	Timeout *metav1.Duration `json:"timeout,omitempty"`
}

// BootstrapTokenDiscovery is used to set the options for bootstrap token based discovery
type BootstrapTokenDiscovery struct {
	// A token used to validate cluster information fetched from the control-plane.
	Token string `json:"token"`

	// APIServerEndpoint is an IP or domain name to the API server from which info will be fetched.
	APIServerEndpoint string `json:"apiServerEndpoint,omitempty"`

	// This specifies a set of public key pins to verify
	// when token-based discovery is used. The root CA found during discovery
	// must match one of these values. Specifying an empty set disables root CA
	// pinning, which can be unsafe. Each hash is specified as "<type>:<value>",
	// where the only currently supported type is "sha256". This is a hex-encoded
	// SHA-256 hash of the Subject Public Key Info (SPKI) object in DER-encoded
	// ASN.1. These hashes can be calculated using, for example, OpenSSL.
	CACertHashes []string `json:"caCertHashes,omitempty"`

	// This allows token-based discovery without CA verification via CACertHashes.
	// This can weaken the security of kubeadm since other nodes can impersonate the control-plane.
	UnsafeSkipCAVerification bool `json:"unsafeSkipCAVerification"`
}

// FileDiscovery is used to specify a file or URL to a kubeconfig file from which to load cluster information
type FileDiscovery struct {
	// The actual file path or URL to the kubeconfig file from which to load cluster information.
	KubeConfigPath string `json:"kubeConfigPath"`
}

// HostPathMount contains elements describing volumes that are mounted from the
// host.
type HostPathMount struct {
	// Name of the volume inside the pod template.
	Name string `json:"name"`
	// The path in the host that will be mounted inside the pod.
	HostPath string `json:"hostPath"`
	// The path inside the Pod where hostPath will be mounted.
	MountPath string `json:"mountPath"`
	// This controls write access to the volume
	ReadOnly bool `json:"readOnly,omitempty"`
	// The type of the hostPath.
	PathType v1.HostPathType `json:"pathType,omitempty"`
}
