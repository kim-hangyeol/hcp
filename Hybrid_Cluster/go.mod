module Hybrid_Cluster

go 1.16

require (
	admiralty.io/multicluster-controller v0.6.0
	github.com/aws/aws-sdk-go v1.40.8
	github.com/ghodss/yaml v1.0.0
	github.com/jinzhu/copier v0.3.2
	github.com/spf13/cobra v1.2.1
	github.com/spf13/viper v1.8.1
	golang.org/x/net v0.0.0-20210614182718-04defd469f4e
	google.golang.org/api v0.44.0
	gopkg.in/yaml.v2 v2.4.0
	k8s.io/api v0.21.2
	k8s.io/apimachinery v0.21.2
	k8s.io/client-go v0.21.2
	k8s.io/klog v1.0.0
	sigs.k8s.io/aws-iam-authenticator v0.5.3
	sigs.k8s.io/controller-runtime v0.9.3
	sigs.k8s.io/kubefed v0.8.1
)
