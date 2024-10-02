package main

// ClusterRoleConfig defines the configuration for an EKS cluster role.
type ClusterRoleConfig struct {
	// Name is the name of the cluster role.
	Name string

	// AssumedBy is the entity that is allowed to assume the role.
	AssumedBy string
}

// NewClusterRoleConfig initializes a new instance of ClusterRoleConfig with default values.
// Returns a pointer to the newly created ClusterRoleConfig struct.
func NewClusterRoleConfig() *ClusterRoleConfig {
	return &ClusterRoleConfig{
		Name:      "eks-cdk-cluster-role",
		AssumedBy: "eks.amazonaws.com",
	}
}
