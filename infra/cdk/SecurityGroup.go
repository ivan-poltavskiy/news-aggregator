package main

// SecurityGroup defines the configuration for a security group.
type SecurityGroup struct {
	// Name is the name of the security group.
	Name string
}

// NewSecurityGroupConfig initializes a new instance of SecurityGroup with default values.
// Returns a pointer to the newly created SecurityGroup struct.
func NewSecurityGroupConfig() *SecurityGroup {
	return &SecurityGroup{
		Name: "ivan-cdk-sg",
	}
}
