package main

// RouteTableConfig defines the configuration for a route table.
type RouteTableConfig struct {
	// Name is the name of the route table.
	Name string

	// TagsKey specifies the key for the route table's tags.
	TagsKey string

	// TagsValue specifies the value for the route table's tags.
	TagsValue string
}

// NewRouteTableConfig initializes a new instance of RouteTableConfig with default values.
// Returns a pointer to the newly created RouteTableConfig struct.
func NewRouteTableConfig() *RouteTableConfig {
	return &RouteTableConfig{
		Name:      "ivan-cdk-rt",
		TagsKey:   "name",
		TagsValue: "ivan-cdk-rt",
	}
}
