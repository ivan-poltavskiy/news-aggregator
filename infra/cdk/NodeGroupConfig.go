package main

import "github.com/aws/aws-cdk-go/awscdk/v2/awseks"

// NodeGroupConfig defines the configuration for an EKS node group.
type NodeGroupConfig struct {
	// Name is the name of the node group.
	Name string

	// InstanceTypes specifies the instance type for the node group.
	InstanceTypes string

	// MinSize specifies the minimum number of nodes in the node group.
	MinSize int

	// MaxSize specifies the maximum number of nodes in the node group.
	MaxSize int

	// DesiredSize specifies the desired number of nodes in the node group.
	DesiredSize int

	// SshKeyName is the name of the SSH key pair to use for nodes in the group.
	SshKeyName string

	// AmiType specifies the AMI type to use for the node group.
	AmiType awseks.NodegroupAmiType

	// DiskSize specifies the disk size (in GB) for the nodes.
	DiskSize int
}

// NewNodeGroupConfig initializes a new instance of NodeGroupConfig with default values.
// Returns a pointer to the newly created NodeGroupConfig struct.
func NewNodeGroupConfig() *NodeGroupConfig {
	return &NodeGroupConfig{
		Name:          "node-group",
		InstanceTypes: "t2.medium",
		MinSize:       1,
		MaxSize:       10,
		DesiredSize:   2,
		SshKeyName:    "ivan-kp",
		AmiType:       awseks.NodegroupAmiType_AL2_X86_64,
		DiskSize:      20,
	}
}
