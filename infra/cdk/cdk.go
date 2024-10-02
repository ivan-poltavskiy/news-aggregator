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
	vpcCidr              = "10.0.0.0/16"
	vpcName              = "ivan-cdk-vpc"
	maxAzs               = 2
	publicSubnetName     = "ivan-cdk-public-subnet-1"
	privateSubnetName    = "ivan-cdk-private-subnet-1"
	subnetCidrMask       = 20
	accountID            = "406477933661"
	region               = "eu-west-2"
	routeTableName       = "ivan-cdk-rt"
	nodeGroupMinSize     = 1
	nodeGroupMaxSize     = 10
	nodeGroupDesiredSize = 2
	sshKeyName           = "ivan-kp"
	diskSize             = 20
)

// NewCdkIvanStack Returns the new stack with the cluster, nodegroup and additional resources
func NewCdkIvanStack(scope constructs.Construct, id string, props *CdkStackProps) awscdk.Stack {
	stack := awscdk.NewStack(scope, &id, &props.StackProps)

	//creates new vpc with parameters from vpcConfig
	vpc := awsec2.NewVpc(stack, jsii.String(vpcName), &awsec2.VpcProps{
		IpAddresses: awsec2.IpAddresses_Cidr(jsii.String(vpcCidr)),
		MaxAzs:      jsii.Number(maxAzs),
		SubnetConfiguration: &[]*awsec2.SubnetConfiguration{
			{
				Name:                jsii.String(publicSubnetName),
				SubnetType:          awsec2.SubnetType_PUBLIC,
				CidrMask:            jsii.Number(subnetCidrMask),
				MapPublicIpOnLaunch: jsii.Bool(true),
			},
			{
				Name:       jsii.String(privateSubnetName),
				SubnetType: awsec2.SubnetType_PRIVATE_WITH_EGRESS,
				CidrMask:   jsii.Number(subnetCidrMask),
			},
		},
		VpcName: jsii.String(vpcName),
	})

	// creates new route table with id of the vpc
	routeTable := awsec2.NewCfnRouteTable(stack, jsii.String(routeTableName), &awsec2.CfnRouteTableProps{
		VpcId: vpc.VpcId(),
		Tags: &[]*awscdk.CfnTag{
			{
				Key:   jsii.String("name"),
				Value: jsii.String(routeTableName),
			},
		},
	})

	// Adds a route to the internet gateway
	awsec2.NewCfnRoute(stack, jsii.String("ivan-cdk-route"), &awsec2.CfnRouteProps{
		RouteTableId:         routeTable.AttrRouteTableId(),
		DestinationCidrBlock: jsii.String("0.0.0.0/0"),
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
		Version:         awseks.KubernetesVersion_Of(jsii.String("1.30")),
		ClusterName:     jsii.String("eks-cdk-cluster"),
		Vpc:             vpc,
		SecurityGroup:   securityGroup,
		DefaultCapacity: jsii.Number(0),
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
			awsec2.NewInstanceType(jsii.String("t2.medium")),
		},
		NodeRole:    nodeGroupRole,
		MinSize:     jsii.Number(nodeGroupMinSize),
		MaxSize:     jsii.Number(nodeGroupMaxSize),
		DesiredSize: jsii.Number(nodeGroupDesiredSize),
		RemoteAccess: &awseks.NodegroupRemoteAccess{
			SshKeyName: jsii.String(sshKeyName),
		},
		Subnets: &awsec2.SubnetSelection{
			Subnets: vpc.PublicSubnets(),
		},
		AmiType:  awseks.NodegroupAmiType_AL2_X86_64,
		DiskSize: jsii.Number(diskSize),
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

	flag.StringVar(&vpcCidr, "vpc-cidr", vpcCidr, "CIDR block for the VPC")
	flag.StringVar(&vpcName, "vpc-name", vpcName, "Name of the VPC")
	flag.IntVar(&maxAzs, "max-azs", maxAzs, "Maximum Availability Zones")
	flag.StringVar(&publicSubnetName, "public-subnet-name", publicSubnetName, "Name of the public subnet")
	flag.StringVar(&privateSubnetName, "private-subnet-name", privateSubnetName, "Name of the private subnet")
	flag.IntVar(&subnetCidrMask, "subnet-cidr-mask", subnetCidrMask, "Subnet CIDR mask")
	flag.StringVar(&accountID, "account-id", accountID, "AWS Account ID")
	flag.StringVar(&region, "region", region, "AWS Region")
	flag.StringVar(&routeTableName, "route-table-name", routeTableName, "Name of the route table")
	flag.IntVar(&nodeGroupMinSize, "node-group-min-size", nodeGroupMinSize, "Minimum size of the node group")
	flag.IntVar(&nodeGroupMaxSize, "node-group-max-size", nodeGroupMaxSize, "Maximum size of the node group")
	flag.IntVar(&nodeGroupDesiredSize, "node-group-desired-size", nodeGroupDesiredSize, "Desired size of the node group")
	flag.StringVar(&sshKeyName, "ssh-key-name", sshKeyName, "SSH Key name for EC2 instances")
	flag.IntVar(&diskSize, "disk-size", diskSize, "Disk size for EC2 instances")
	flag.Parse()

	app := awscdk.NewApp(nil)

	NewCdkIvanStack(app, "CdkStack", &CdkStackProps{
		awscdk.StackProps{
			Env: env(),
		},
	})

	app.Synth(nil)
}

// env determines the AWS environment (account+region) in which our stack is to
// be deployed. For more information see: https://docs.aws.amazon.com/cdk/latest/guide/environments.html
func env() *awscdk.Environment {
	return &awscdk.Environment{
		Account: jsii.String(accountID),
		Region:  jsii.String(region),
	}
}
