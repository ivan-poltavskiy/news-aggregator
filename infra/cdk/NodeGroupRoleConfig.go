package main

// NodeGroupRoleConfig defines the configuration for a node group IAM role.
type NodeGroupRoleConfig struct {
	// Name is the name of the node group role.
	Name string

	// AssumedBy is the entity that is allowed to assume the role.
	AssumedBy string
}

// NewNodeGroupRoleConfig initializes a new instance of NodeGroupRoleConfig with default values.
// Returns a pointer to the newly created NodeGroupRoleConfig struct.
func NewNodeGroupRoleConfig() *NodeGroupRoleConfig {
	return &NodeGroupRoleConfig{
		Name:      "node-group-role",
		AssumedBy: "ec2.amazonaws.com",
	}
}
