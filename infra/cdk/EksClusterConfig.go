package main

// EksClusterConfig defines the configuration for an EKS cluster.
type EksClusterConfig struct {
	// Name is the name of the EKS cluster.
	Name string

	// K8SVersion specifies the version of Kubernetes to use.
	K8SVersion string

	// DefaultCapacity is the default capacity of the EKS cluster.
	DefaultCapacity int

	// UserArn is the ARN of the IAM user associated with the cluster.
	UserArn string

	// UserId is the ID of the IAM user associated with the cluster.
	UserId string
}

// NewEksClusterConfig initializes a new instance of EksClusterConfig with default values.
// Returns a pointer to the newly created EksClusterConfig struct.
func NewEksClusterConfig() *EksClusterConfig {
	return &EksClusterConfig{
		Name:            "eks-cdk-cluster",
		K8SVersion:      "1.30",
		DefaultCapacity: 0,
		UserArn:         "arn:aws:iam::406477933661:user/ivan",
		UserId:          "ivan",
	}
}
