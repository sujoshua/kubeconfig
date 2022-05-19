/*
Copyright 2019 The Kubernetes Authors.

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

package v1beta2

import (
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// InitConfiguration contains a list of elements that is specific "kubeadm init"-only runtime
// information.
type InitConfiguration struct {
	metav1.TypeMeta `json:",inline"`

	// `kubeadm init`-only information. These fields are solely used the first time `kubeadm init` runs.
	// After that, the information in the fields IS NOT uploaded to the `kubeadm-config` ConfigMap
	// that is used by `kubeadm upgrade` for instance. These fields must be omitempty.

	// `bootstrapTokens` is respected at `kubeadm init` time and describes a set of bootstrap tokens to create.
	// This information IS NOT uploaded to the kubeadm cluster ConfigMap, partly because of its sensitive nature.
	BootstrapTokens []BootstrapToken `json:"bootstrapTokens,omitempty"`

	// `nodeRegistration` holds fields that relate to registering the new control-plane node to the cluster.
	NodeRegistration NodeRegistrationOptions `json:"nodeRegistration,omitempty"`

	// `localAPIEndpoint` represents the endpoint of the API server instance that's deployed on this control plane node.
	// In HA setups, this differs from `ClusterConfiguration.controlPlaneEndpoint` in the sense that ControlPlaneEndpoint
	// is the global endpoint for the cluster, which then load-balances the requests to each individual API server. This
	// configuration object lets you customize what IP/DNS name and port the local API server advertises it's accessible
	// on. By default, kubeadm tries to auto-detect the IP of the default interface and use that, but in case that process
	// fails you may set the desired value here.
	LocalAPIEndpoint APIEndpoint `json:"localAPIEndpoint,omitempty"`

	// `certificateKey` sets the key with which certificates and keys are encrypted prior to being uploaded in
	// a secret in the cluster during the `uploadcerts init` phase.
	CertificateKey string `json:"certificateKey,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ClusterConfiguration contains cluster-wide configuration for a kubeadm cluster
type ClusterConfiguration struct {
	metav1.TypeMeta `json:",inline"`

	// `etcd` holds configuration for etcd.
	Etcd Etcd `json:"etcd,omitempty"`

	// `networking` holds configuration for the networking topology of the cluster.
	Networking Networking `json:"networking,omitempty"`

	// `kubernetesVersion` is the target version of the control plane.
	KubernetesVersion string `json:"kubernetesVersion,omitempty"`

	// `controlPlaneEndpoint` sets a stable IP address or DNS name for the control plane; it
	// can be a valid IP address or a RFC-1123 DNS subdomain, both with optional TCP port.
	// In case the `controlPlaneEndpoint` is not specified, the `advertiseAddress` + `bindPort`
	// are used; in case the `controlPlaneEndpoint` is specified but without a TCP port,
	// the `bindPort` is used.
	// Possible usages are:
	//
	// - In a cluster with more than one control plane instances, this field should be
	//   assigned the address of the external load balancer in front of the
	//   control plane instances.
	// - In environments with enforced node recycling, the `controlPlaneEndpoint`
	//   could be used for assigning a stable DNS to the control plane.
	ControlPlaneEndpoint string `json:"controlPlaneEndpoint,omitempty"`

	// `apiServer` contains extra settings for the API server.
	APIServer APIServer `json:"apiServer,omitempty"`

	// `controllerManager` contains extra settings for the controller manager.
	ControllerManager ControlPlaneComponent `json:"controllerManager,omitempty"`

	// `scheduler` contains extra settings for the scheduler.
	Scheduler ControlPlaneComponent `json:"scheduler,omitempty"`

	// `dns` defines the options for the DNS add-on installed in the cluster.
	DNS DNS `json:"dns,omitempty"`

	// `certificatesDir` specifies where to store or look for all required certificates.
	CertificatesDir string `json:"certificatesDir,omitempty"`

	// `imageRepository` sets the container registry to pull images from.
	// If empty, `k8s.gcr.io` will be used by default; in case of kubernetes version is
	// a CI build (kubernetes version starts with `ci/`) `gcr.io/k8s-staging-ci-images`
	// is used as a default for control plane components and for kube-proxy, while
	// `k8s.gcr.io` will be used for all the other images.
	ImageRepository string `json:"imageRepository,omitempty"`

	// `useHyperKubeImage` controls if hyperkube should be used for Kubernetes components
	// instead of their respective separate images.
	// DEPRECATED: As `hyperkube` is itself deprecated, this fields is too. It will be
	// removed in future kubeadm config versions, kubeadm will print multiple warnings
	// when this set to true, and at some point it may become ignored.
	UseHyperKubeImage bool `json:"useHyperKubeImage,omitempty"`

	// `featureGates` contains the feature gates enabled by the user.
	FeatureGates map[string]bool `json:"featureGates,omitempty"`

	// The cluster name.
	ClusterName string `json:"clusterName,omitempty"`
}

// ControlPlaneComponent holds settings common to control plane component of the cluster
type ControlPlaneComponent struct {
	// `extraArgs` is an extra set of flags to pass to a control plane component.
	// A key in this map is the flag name as it appears on the command line except
	// without leading dash(es).
	ExtraArgs map[string]string `json:"extraArgs,omitempty"`

	// `extraVolumes` is an extra set of host volumes mounted to the control plane
	// component.
	ExtraVolumes []HostPathMount `json:"extraVolumes,omitempty"`
}

// APIServer holds settings necessary for API server deployments in the cluster.
type APIServer struct {
	ControlPlaneComponent `json:",inline"`

	// `certSANs` sets extra Subject Alternative Names (SANs) for the API Server
	// signing certificate.
	CertSANs []string `json:"certSANs,omitempty"`

	// `timeoutForControlPlane` controls the timeout that we wait for the API server
	// to appear.
	TimeoutForControlPlane *metav1.Duration `json:"timeoutForControlPlane,omitempty"`
}

// DNSAddOnType defines string identifying DNS add-on types.
type DNSAddOnType string

const (
	// CoreDNS add-on type.
	CoreDNS DNSAddOnType = "CoreDNS"
)

// DNS defines the DNS addon that should be used in the cluster
type DNS struct {
	// `type` defines the DNS add-on to be used.
	Type DNSAddOnType `json:"type"`

	// ImageMeta allows to customize the image used for the DNS component
	ImageMeta `json:",inline"`
}

// ImageMeta allows to customize the image used for components that are not
// originated from the Kubernetes/Kubernetes release process
type ImageMeta struct {
	// `imageRepository` sets the container registry to pull images from.
	// If not set, the `imageRepository` defined in ClusterConfiguration will be used.
	ImageRepository string `json:"imageRepository,omitempty"`

	// `imageTag` allows for specifying a tag for the image.
	// In case this value is set, kubeadm does not change automatically the
	// version of the above components during upgrades.
	ImageTag string `json:"imageTag,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ClusterStatus contains the cluster status. The ClusterStatus will be stored in
// the kubeadm-config ConfigMap in the cluster, and then updated by kubeadm when
// additional control plane instance joins or leaves the cluster.
type ClusterStatus struct {
	metav1.TypeMeta `json:",inline"`

	// `apiEndpoints` currently available in the cluster, one for each control
	// plane/API server instance.
	// The key of the map is the IP of the host's default interface.
	APIEndpoints map[string]APIEndpoint `json:"apiEndpoints"`
}

// APIEndpoint struct contains elements of API server instance deployed on a node.
type APIEndpoint struct {
	// `advertiseAddress` sets the IP address for the API server to advertise.
	AdvertiseAddress string `json:"advertiseAddress,omitempty"`

	// `bindPort` sets the secure port for the API Server to bind to.
	// Defaults to 6443.
	BindPort int32 `json:"bindPort,omitempty"`
}

// NodeRegistrationOptions holds fields that relate to registering a new control-plane
// or node to the cluster, either via "kubeadm init" or "kubeadm join".
type NodeRegistrationOptions struct {

	// `name` is the `.Metadata.Name` field of the Node API object that will be created
	// in this `kubeadm init` or `kubeadm join` operation.
	// This field is also used in the `CommonName` field of the kubelet's client certificate
	// to the API server.
	// Defaults to the hostname of the node if not provided.
	Name string `json:"name,omitempty"`

	// `criSocket` is used to retrieve container runtime information. This information will
	// be annotated to the Node API object, for later re-use.
	CRISocket string `json:"criSocket,omitempty"`

	// `taints` specifies the taints the Node API object should be registered with.
	// If this field is unset, i.e. nil, in the `kubeadm init` process it will be defaulted with
	// a control-plane taint for control-plane nodes. If you don't want to taint your control-plane
	// node, set this field to an empty list, i.e. `taints: []`, in the YAML file. This field is
	// solely used for Node registration.
	Taints []v1.Taint `json:"taints"`

	// `kubeletExtraArgs` passes through extra arguments to the kubelet. The arguments here are
	// passed to the kubelet command line via the environment file kubeadm writes at runtime for
	// the kubelet to source. This overrides the generic base-level configuration in the
	// 'kubelet-config-1.X' ConfigMap.
	// Flags have higher priority when parsing. These values are local and specific to the node
	// kubeadm is executing on.
	// A key in this map is the flag name as it appears on the command line except without leading dash(es).
	KubeletExtraArgs map[string]string `json:"kubeletExtraArgs,omitempty"`

	// `ignorePreflightErrors` provides a list of pre-flight errors to be ignored when the
	// current node is registered.
	IgnorePreflightErrors []string `json:"ignorePreflightErrors,omitempty"`
}

// Networking contains elements describing cluster's networking configuration
type Networking struct {
	// `serviceSubnet` is the subnet used by kubernetes Services. Defaults to "10.96.0.0/12".
	ServiceSubnet string `json:"serviceSubnet,omitempty"`
	// `podSubnet` is the subnet used by Pods.
	PodSubnet string `json:"podSubnet,omitempty"`
	// `dnsDomain` is the DNS domain used by kubernetes Services. Defaults to "cluster.local".
	DNSDomain string `json:"dnsDomain,omitempty"`
}

// BootstrapToken describes one bootstrap token, stored as a Secret in the cluster
type BootstrapToken struct {
	// `token` is used for establishing bidirectional trust between nodes and control-planes.
	// Used for joining nodes in the cluster.
	Token *BootstrapTokenString `json:"token"`
	// `description` sets a human-friendly message why this token exists and what it's used
	// for, so other administrators can know its purpose.
	Description string `json:"description,omitempty"`
	// `ttl` defines the time to live for this token. Defaults to '24h'.
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

	// `local` provides configuration knobs for configuring the local etcd instance.
	// `local` and `external` are mutually exclusive.
	Local *LocalEtcd `json:"local,omitempty"`

	// `external` describes how to connect to an external etcd cluster.
	// `local` and `external` are mutually exclusive.
	External *ExternalEtcd `json:"external,omitempty"`
}

// LocalEtcd describes that kubeadm should run an etcd cluster locally.
type LocalEtcd struct {
	// ImageMeta allows to customize the container used for etcd.
	ImageMeta `json:",inline"`

	// `dataDir` is the directory etcd will place its data.
	// Defaults to "/var/lib/etcd".
	DataDir string `json:"dataDir"`

	// `extraArgs` are extra arguments provided to the etcd binary when run
	// inside a static pod.
	// A key in this map is the flag name as it appears on the
	// command line except without leading dash(es).
	ExtraArgs map[string]string `json:"extraArgs,omitempty"`

	// `serverCertSANs` sets extra Subject Alternative Names (SANs) for the
	// etcd server signing certificate.
	ServerCertSANs []string `json:"serverCertSANs,omitempty"`
	// `peerCertSANs` sets extra Subject Alternative Names (SANs) for the
	// etcd peer signing certificate.
	PeerCertSANs []string `json:"peerCertSANs,omitempty"`
}

// ExternalEtcd describes an external etcd cluster.
// Kubeadm has no knowledge of where certificate files live and they must be supplied.
type ExternalEtcd struct {
	// `endpoints` of etcd members.
	Endpoints []string `json:"endpoints"`

	// `caFile` is an SSL Certificate Authority (CA) file used to secure etcd communication.
	// Required if using a TLS connection.
	CAFile string `json:"caFile"`

	// `certFile` is an SSL certification file used to secure etcd communication.
	// Required if using a TLS connection.
	CertFile string `json:"certFile"`

	// `keyFile` is an SSL key file used to secure etcd communication.
	// Required if using a TLS connection.
	KeyFile string `json:"keyFile"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// JoinConfiguration contains elements describing a particular node.
type JoinConfiguration struct {
	metav1.TypeMeta `json:",inline"`

	// `nodeRegistration` holds fields that relate to registering the new
	// control-plane node to the cluster
	NodeRegistration NodeRegistrationOptions `json:"nodeRegistration,omitempty"`

	// `caCertPath` is the path to the SSL certificate authority used to
	// secure comunications between a node and the control-plane.
	// Defaults to "/etc/kubernetes/pki/ca.crt".
	CACertPath string `json:"caCertPath,omitempty"`

	// `discovery` specifies the options for the kubelet to use during the TLS
	// bootstrap process.
	Discovery Discovery `json:"discovery"`

	// `controlPlane` defines the additional control plane instance to be deployed
	// on the joining node. If nil, no additional control plane instance will be deployed.
	ControlPlane *JoinControlPlane `json:"controlPlane,omitempty"`
}

// JoinControlPlane contains elements describing an additional control plane instance
// to be deployed on the joining node.
type JoinControlPlane struct {
	// `localAPIEndpoint` represents the endpoint of the API server instance
	// to be deployed on this node.
	LocalAPIEndpoint APIEndpoint `json:"localAPIEndpoint,omitempty"`

	// `certificateKey` is the key that is used for decryption of certificates after
	// they are downloaded from the secret upon joining a new control plane node.
	// The corresponding encryption key is in the InitConfiguration.
	CertificateKey string `json:"certificateKey,omitempty"`
}

// Discovery specifies the options for the kubelet to use during the TLS Bootstrap process
type Discovery struct {
	// `bootstrapToken` is used to set the options for bootstrap token based discovery.
	// `bootstrapToken` and `file` are mutually exclusive.
	BootstrapToken *BootstrapTokenDiscovery `json:"bootstrapToken,omitempty"`

	// `file` is used to specify a file or URL to a kubeconfig file from which to load
	// cluster information.
	// `bootstrapToken` and `file` are mutually exclusive.
	File *FileDiscovery `json:"file,omitempty"`

	// `tlsBootstrapToken` is a token used for TLS bootstrapping.
	// If `bootstrapToken` is set, this field is defaulted to `.bootstrapToken.token`,
	// but can be overridden.
	// If `file` is set, this field **must be set** in case the KubeConfigFile does not
	// contain any other authentication information.
	TLSBootstrapToken string `json:"tlsBootstrapToken,omitempty"`

	// `timeout` modifies the discovery timeout.
	Timeout *metav1.Duration `json:"timeout,omitempty"`
}

// BootstrapTokenDiscovery is used to set the options for bootstrap token based discovery
type BootstrapTokenDiscovery struct {
	// `token` is a token used to validate cluster information fetched from
	// the control-plane.
	Token string `json:"token"`

	// `apiServerEndpoint` is an IP or domain name to the API server from which information
	// will be fetched.
	APIServerEndpoint string `json:"apiServerEndpoint,omitempty"`

	// `caCertHashes` specifies a set of public key pins to verify when token-based discovery
	// is used. The root CA found during discovery must match one of these values.
	// Specifying an empty set disables root CA pinning, which can be unsafe.
	// Each hash is specified as "<type>:<value>", where the only currently supported type is "sha256".
	// This is a hex-encoded SHA-256 hash of the Subject Public Key Info (SPKI) object in
	// DER-encoded ASN.1. These hashes can be calculated using, for example, OpenSSL.
	CACertHashes []string `json:"caCertHashes,omitempty"`

	// `unsafeSkipCAVerification` allows token-based discovery without CA verification via
	// `caCertHashes`. This can weaken the security of kubeadm since other nodes can
	// impersonate the control-plane.
	UnsafeSkipCAVerification bool `json:"unsafeSkipCAVerification,omitempty"`
}

// FileDiscovery is used to specify a file or URL to a kubeconfig file from which to load cluster information
type FileDiscovery struct {
	// `kubeConfigPath` is used to specify the actual file path or URL to the kubeconfig file
	// from which to load cluster information.
	KubeConfigPath string `json:"kubeConfigPath"`
}

// HostPathMount contains elements describing volumes that are mounted from the host.
type HostPathMount struct {
	// `name` of the volume inside the Pod template.
	Name string `json:"name"`
	// `hostPath` is the path in the host that will be mounted inside the Pod.
	HostPath string `json:"hostPath"`
	// `mountPath`is the path inside the Pod where hostPath volume will be mounted.
	MountPath string `json:"mountPath"`
	// `readOnly` controls write access to the volume.
	ReadOnly bool `json:"readOnly,omitempty"`
	// `pathType` is the type of the HostPath.
	PathType v1.HostPathType `json:"pathType,omitempty"`
}
