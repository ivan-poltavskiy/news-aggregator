package main

// VpcConfig defines the configuration for a VPC (Virtual Private Cloud).
type VpcConfig struct {
	// Name is the name of the VPC.
	Name string

	// Cidr specifies the CIDR block for the VPC.
	Cidr string

	// MaxAZs specifies the maximum number of availability zones for the VPC.
	MaxAZs int

	// PublicSubnetName is the name of the public subnet.
	PublicSubnetName string

	// PrivateSubnetName is the name of the private subnet.
	PrivateSubnetName string

	// PublicMask specifies the subnet mask for the public subnet.
	PublicMask int

	// PrivateMask specifies the subnet mask for the private subnet.
	PrivateMask int
}

// NewVpcConfig initializes a new instance of VpcConfig with default values.
// Returns a pointer to the newly created VpcConfig struct.
func NewVpcConfig() *VpcConfig {
	return &VpcConfig{
		Name:              "ivan-cdk-vpc",
		Cidr:              "10.0.0.0/16",
		MaxAZs:            2,
		PublicSubnetName:  "ivan-cdk-public-subnet-1",
		PrivateSubnetName: "ivan-cdk-private-subnet-1",
		PublicMask:        20,
		PrivateMask:       20,
	}
}
