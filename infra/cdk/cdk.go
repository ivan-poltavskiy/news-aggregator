package main

import (
	"flag"
	"fmt"
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsec2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awseks"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type CdkStackProps struct {
	awscdk.StackProps
}

var (
	defaultVpcCidr              = "10.0.0.0/16"
	defaultVpcName              = "ivan-cdk-vpc"
	defaultMaxAzs               = 2
	defaultPublicSubnetName     = "ivan-cdk-public-subnet-1"
	defaultPrivateSubnetName    = "ivan-cdk-private-subnet-1"
	defaultSubnetCidrMask       = 20
	defaultAccountID            = "406477933661"
	defaultRegion               = "eu-west-2"
	defaultRouteTableName       = "ivan-cdk-rt"
	defaultNodeGroupMinSize     = 1
	defaultNodeGroupMaxSize     = 10
	defaultNodeGroupDesiredSize = 2
	defaultSshKeyName           = "ivan-kp"
	defaultDiskSize             = 20
	defaultCapacity             = 0
	defaultDestinationCidrBlock = "0.0.0.0/0"
	defaultK8SVersion           = "1.30"
	defaultInstanceType         = "t2.medium"
)

// NewCdkIvanStack Returns the new stack with the cluster, nodegroup and additional resources
func NewCdkIvanStack(scope constructs.Construct, id string, props *CdkStackProps) awscdk.Stack {
	stack := awscdk.NewStack(scope, &id, &props.StackProps)

	//creates new vpc with parameters from vpcConfig
	vpc := awsec2.NewVpc(stack, jsii.String(defaultVpcName), &awsec2.VpcProps{
		IpAddresses: awsec2.IpAddresses_Cidr(jsii.String(defaultVpcCidr)),
		MaxAzs:      jsii.Number(defaultMaxAzs),
		SubnetConfiguration: &[]*awsec2.SubnetConfiguration{
			{
				Name:                jsii.String(defaultPublicSubnetName),
				SubnetType:          awsec2.SubnetType_PUBLIC,
				CidrMask:            jsii.Number(defaultSubnetCidrMask),
				MapPublicIpOnLaunch: jsii.Bool(true),
			},
			{
				Name:       jsii.String(defaultPrivateSubnetName),
				SubnetType: awsec2.SubnetType_PRIVATE_WITH_EGRESS,
				CidrMask:   jsii.Number(defaultSubnetCidrMask),
			},
		},
		VpcName: jsii.String(defaultVpcName),
	})

	// creates new route table with id of the vpc
	routeTable := awsec2.NewCfnRouteTable(stack, jsii.String(defaultRouteTableName), &awsec2.CfnRouteTableProps{
		VpcId: vpc.VpcId(),
		Tags: &[]*awscdk.CfnTag{
			{
				Key:   jsii.String("name"),
				Value: jsii.String(defaultRouteTableName),
			},
		},
	})

	// Adds a route to the internet gateway
	awsec2.NewCfnRoute(stack, jsii.String("ivan-cdk-route"), &awsec2.CfnRouteProps{
		RouteTableId:         routeTable.AttrRouteTableId(),
		DestinationCidrBlock: jsii.String(defaultDestinationCidrBlock),
		GatewayId:            vpc.InternetGatewayId(),
	})

	// Associates public subnets with the route table
	publicSubnets := vpc.PublicSubnets()
	for i, subnet := range *publicSubnets {
		awsec2.NewCfnSubnetRouteTableAssociation(stack, jsii.String(fmt.Sprintf("SubnetAssociation%d", i)), &awsec2.CfnSubnetRouteTableAssociationProps{
			SubnetId:     subnet.SubnetId(),
			RouteTableId: routeTable.Ref(),
		})
	}

	// Security group
	securityGroup := awsec2.NewSecurityGroup(stack, jsii.String("ivan-cdk-sg"), &awsec2.SecurityGroupProps{
		Vpc:               vpc,
		SecurityGroupName: jsii.String("ivan-cdk-sg"),
		AllowAllOutbound:  jsii.Bool(true),
	})

	// Ingress rules
	securityGroup.AddIngressRule(awsec2.Peer_AnyIpv4(), awsec2.Port_Tcp(jsii.Number(80)), jsii.String("HTTP"), jsii.Bool(true))
	securityGroup.AddIngressRule(awsec2.Peer_AnyIpv4(), awsec2.Port_Tcp(jsii.Number(443)), jsii.String("HTTPS"), jsii.Bool(true))
	securityGroup.AddIngressRule(awsec2.Peer_AnyIpv4(), awsec2.Port_Tcp(jsii.Number(22)), jsii.String("SSH"), jsii.Bool(true))

	// IAM role for EKS cluster
	clusterRole := awsiam.NewRole(stack, jsii.String("eks-cdk-cluster-role"), &awsiam.RoleProps{
		AssumedBy: awsiam.NewServicePrincipal(jsii.String("eks.amazonaws.com"), nil),
		ManagedPolicies: &[]awsiam.IManagedPolicy{
			awsiam.ManagedPolicy_FromAwsManagedPolicyName(jsii.String("AmazonEKSClusterPolicy")),
			awsiam.ManagedPolicy_FromAwsManagedPolicyName(jsii.String("AmazonEKSServicePolicy")),
		},
	})

	// EKS cluster
	eksCluster := awseks.NewCluster(stack, jsii.String("eks-cdk-cluster"), &awseks.ClusterProps{
		Version:         awseks.KubernetesVersion_Of(jsii.String(defaultK8SVersion)),
		ClusterName:     jsii.String("eks-cdk-cluster"),
		Vpc:             vpc,
		SecurityGroup:   securityGroup,
		DefaultCapacity: jsii.Number(defaultCapacity),
		Role:            clusterRole,
	})

	// IAM role for node group
	nodeGroupRole := awsiam.NewRole(stack, jsii.String("node-group-role"), &awsiam.RoleProps{
		AssumedBy:   awsiam.NewServicePrincipal(jsii.String("ec2.amazonaws.com"), nil),
		Description: jsii.String("Role for EKS Node Group"),
		ManagedPolicies: &[]awsiam.IManagedPolicy{
			awsiam.ManagedPolicy_FromAwsManagedPolicyName(jsii.String("AmazonEC2ContainerRegistryReadOnly")),
			awsiam.ManagedPolicy_FromAwsManagedPolicyName(jsii.String("AmazonEC2FullAccess")),
			awsiam.ManagedPolicy_FromAwsManagedPolicyName(jsii.String("AmazonEKS_CNI_Policy")),
			awsiam.ManagedPolicy_FromAwsManagedPolicyName(jsii.String("AmazonEKSWorkerNodePolicy")),
		},
		RoleName: jsii.String("node-group-role"),
	})

	// Node group
	eksCluster.AddNodegroupCapacity(jsii.String("node-group"), &awseks.NodegroupOptions{
		InstanceTypes: &[]awsec2.InstanceType{
			awsec2.NewInstanceType(jsii.String(defaultInstanceType)),
		},
		NodeRole:    nodeGroupRole,
		MinSize:     jsii.Number(defaultNodeGroupMinSize),
		MaxSize:     jsii.Number(defaultNodeGroupMaxSize),
		DesiredSize: jsii.Number(defaultNodeGroupDesiredSize),
		RemoteAccess: &awseks.NodegroupRemoteAccess{
			SshKeyName: jsii.String(defaultSshKeyName),
		},
		Subnets: &awsec2.SubnetSelection{
			Subnets: vpc.PublicSubnets(),
		},
		AmiType:  awseks.NodegroupAmiType_AL2_X86_64,
		DiskSize: jsii.Number(defaultDiskSize),
	})

	// Outputs
	awscdk.NewCfnOutput(stack, jsii.String("EKSClusterName"), &awscdk.CfnOutputProps{
		Value: eksCluster.ClusterName(),
	})
	awscdk.NewCfnOutput(stack, jsii.String("VPCId"), &awscdk.CfnOutputProps{
		Value: vpc.VpcId(),
	})
	awscdk.NewCfnOutput(stack, jsii.String("ClusterSecurityGroupId"), &awscdk.CfnOutputProps{
		Value: securityGroup.SecurityGroupId(),
	})
	return stack
}

func main() {
	defer jsii.Close()

	flag.StringVar(&defaultVpcCidr, "vpc-cidr", defaultVpcCidr, "CIDR block for the VPC")
	flag.StringVar(&defaultVpcName, "vpc-name", defaultVpcName, "Name of the VPC")
	flag.IntVar(&defaultMaxAzs, "max-azs", defaultMaxAzs, "Maximum Availability Zones")
	flag.StringVar(&defaultPublicSubnetName, "public-subnet-name", defaultPublicSubnetName, "Name of the public subnet")
	flag.StringVar(&defaultPrivateSubnetName, "private-subnet-name", defaultPrivateSubnetName, "Name of the private subnet")
	flag.IntVar(&defaultSubnetCidrMask, "subnet-cidr-mask", defaultSubnetCidrMask, "Subnet CIDR mask")
	flag.StringVar(&defaultAccountID, "account-id", defaultAccountID, "AWS Account ID")
	flag.StringVar(&defaultRegion, "defaultRegion", defaultRegion, "AWS Region")
	flag.StringVar(&defaultRouteTableName, "route-table-name", defaultRouteTableName, "Name of the route table")
	flag.IntVar(&defaultNodeGroupMinSize, "node-group-min-size", defaultNodeGroupMinSize, "Minimum size of the node group")
	flag.IntVar(&defaultNodeGroupMaxSize, "node-group-max-size", defaultNodeGroupMaxSize, "Maximum size of the node group")
	flag.IntVar(&defaultNodeGroupDesiredSize, "node-group-desired-size", defaultNodeGroupDesiredSize, "Desired size of the node group")
	flag.StringVar(&defaultSshKeyName, "ssh-key-name", defaultSshKeyName, "SSH Key name for EC2 instances")
	flag.StringVar(&defaultDestinationCidrBlock, "destination-cidr-block", defaultDestinationCidrBlock, "The IPv4 CIDR address block used for the destination match")
	flag.StringVar(&defaultK8SVersion, "default-k8s-version", defaultK8SVersion, "Version of the Kubernetes")
	flag.StringVar(&defaultInstanceType, "default-instance-type", defaultInstanceType, "Type of instance for the Node")
	flag.IntVar(&defaultDiskSize, "disk-size", defaultDiskSize, "Disk size for EC2 instances")
	flag.IntVar(&defaultCapacity, "default-capacity", defaultCapacity, "Number of instances to allocate as an initial capacity for this cluster")
	flag.Parse()

	app := awscdk.NewApp(nil)

	NewCdkIvanStack(app, "CdkStack", &CdkStackProps{
		awscdk.StackProps{
			Env: env(),
		},
	})

	app.Synth(nil)
}

// env determines the AWS environment (account+defaultRegion) in which our stack is to
// be deployed. For more information see: https://docs.aws.amazon.com/cdk/latest/guide/environments.html
func env() *awscdk.Environment {
	return &awscdk.Environment{
		Account: jsii.String(defaultAccountID),
		Region:  jsii.String(defaultRegion),
	}
}
