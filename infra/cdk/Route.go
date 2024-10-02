package main

// RouteConfig defines the configuration for a network route.
type RouteConfig struct {
	// Name is the name of the route.
	Name string

	// DestinationCidrBlock is the CIDR block for the destination network.
	DestinationCidrBlock string
}

// NewRouteConfig initializes a new instance of RouteConfig with default values.
// Returns a pointer to the newly created RouteConfig struct.
func NewRouteConfig() *RouteConfig {
	return &RouteConfig{
		Name:                 "ivan-cdk-route",
		DestinationCidrBlock: "0.0.0.0/0",
	}
}
