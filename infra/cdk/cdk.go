package main

import (
	"flag"
	"fmt"
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsec2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awseks"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
	// "github.com/aws/aws-cdk-go/awscdk/v2/awssqs"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type CdkStackProps struct {
	awscdk.StackProps
}

// Initialization structure witch describe the AWS resources and set up default values for these resources.
var vpcConfig = NewVpcConfig()
var routeTableConfig = NewRouteTableConfig()
var routeConfig = NewRouteConfig()
var securityGroupConfig = NewSecurityGroupConfig()
var clusterRoleConfig = NewClusterRoleConfig()
var eksClusterConfig = NewEksClusterConfig()
var nodeGroupRoleConfig = NewNodeGroupRoleConfig()
var nodeGroupConfig = NewNodeGroupConfig()

// Initialization default values for config EKS cluster in AWS
var cdkStackId = "IvanCdkStack"
var accountId = "406477933661"
var region = "eu-west-2"

// NewCdkIvanStack Returns the new stack with the cluster, nodegroup and another additional resources
func NewCdkIvanStack(scope constructs.Construct, id string, props *CdkStackProps) awscdk.Stack {
	stack := awscdk.NewStack(scope, &id, &props.StackProps)

	//creates new vpc with parameters from vpcConfig
	vpc := awsec2.NewVpc(stack, jsii.String(vpcConfig.Name), &awsec2.VpcProps{
		IpAddresses: awsec2.IpAddresses_Cidr(jsii.String(vpcConfig.Cidr)),
		MaxAzs:      jsii.Number(vpcConfig.MaxAZs),
		SubnetConfiguration: &[]*awsec2.SubnetConfiguration{
			{
				Name:                jsii.String(vpcConfig.PublicSubnetName),
				SubnetType:          awsec2.SubnetType_PUBLIC,
				CidrMask:            jsii.Number(vpcConfig.PublicMask),
				MapPublicIpOnLaunch: jsii.Bool(true),
			},
			{
				Name:       jsii.String(vpcConfig.PrivateSubnetName),
				SubnetType: awsec2.SubnetType_PRIVATE_WITH_EGRESS,
				CidrMask:   jsii.Number(vpcConfig.PrivateMask),
			},
		},
		VpcName: jsii.String(vpcConfig.Name),
	})

	// creates new route table with id of the vpc
	routeTable := awsec2.NewCfnRouteTable(stack, jsii.String(routeTableConfig.Name), &awsec2.CfnRouteTableProps{
		VpcId: vpc.VpcId(),
		Tags: &[]*awscdk.CfnTag{
			{
				Key:   jsii.String(routeTableConfig.TagsKey),
				Value: jsii.String(routeTableConfig.TagsValue),
			},
		},
	})

	// add a route to the route table that directs traffic to the internet gateway
	awsec2.NewCfnRoute(stack, jsii.String(routeConfig.Name), &awsec2.CfnRouteProps{
		RouteTableId:         routeTable.AttrRouteTableId(),
		DestinationCidrBlock: jsii.String(routeConfig.DestinationCidrBlock),
		GatewayId:            vpc.InternetGatewayId(),
	})

	// creates route table association for each public subnets in the vpc
	publicSubnets := vpc.PublicSubnets()
	for i, subnet := range *publicSubnets {
		awsec2.NewCfnSubnetRouteTableAssociation(stack, jsii.String(fmt.Sprintf("SubnetAssociation%d", i)), &awsec2.CfnSubnetRouteTableAssociationProps{
			SubnetId:     subnet.SubnetId(),
			RouteTableId: routeTable.Ref(),
		})
	}

	// creates new security group for the current vpc
	securityGroup := awsec2.NewSecurityGroup(stack, jsii.String(securityGroupConfig.Name), &awsec2.SecurityGroupProps{
		Vpc:               vpc,
		SecurityGroupName: jsii.String(securityGroupConfig.Name),
		AllowAllOutbound:  jsii.Bool(true),
	})

	// adds new ingress rule for security group for HTTP traffic
	securityGroup.AddIngressRule(awsec2.Peer_AnyIpv4(), awsec2.Port_Tcp(jsii.Number(80)), jsii.String("HTTP"), jsii.Bool(true))
	// adds new ingress rule for security group for HTTPS traffic
	securityGroup.AddIngressRule(awsec2.Peer_AnyIpv4(), awsec2.Port_Tcp(jsii.Number(443)), jsii.String("HTTPS"), jsii.Bool(true))
	// adds new ingress rule for security group for SSH traffic
	securityGroup.AddIngressRule(awsec2.Peer_AnyIpv4(), awsec2.Port_Tcp(jsii.Number(22)), jsii.String("SSH"), jsii.Bool(true))

	// creates new IAM role for the cluster with provided policies
	clusterRole := awsiam.NewRole(stack, jsii.String(clusterRoleConfig.Name), &awsiam.RoleProps{
		AssumedBy: awsiam.NewServicePrincipal(jsii.String(clusterRoleConfig.AssumedBy), nil),
		ManagedPolicies: &[]awsiam.IManagedPolicy{
			awsiam.ManagedPolicy_FromAwsManagedPolicyName(jsii.String("AmazonEKSClusterPolicy")),
			awsiam.ManagedPolicy_FromAwsManagedPolicyName(jsii.String("AmazonEKSServicePolicy")),
		},
	})

	// creates new EKS cluster with provided version of K8S, provided VPC, Security Group, and Cluster Role.
	eksCluster := awseks.NewCluster(stack, jsii.String(eksClusterConfig.Name), &awseks.ClusterProps{
		Version:         awseks.KubernetesVersion_Of(jsii.String(eksClusterConfig.K8SVersion)),
		ClusterName:     jsii.String(eksClusterConfig.Name),
		Vpc:             vpc,
		SecurityGroup:   securityGroup,
		DefaultCapacity: jsii.Number(eksClusterConfig.DefaultCapacity),
		Role:            clusterRole,
	})

	// adds user to the mapping in the AWS for managing this cluster
	eksCluster.AwsAuth().AddUserMapping(awsiam.User_FromUserArn(stack, jsii.String(eksClusterConfig.UserId), jsii.String(eksClusterConfig.UserArn)), &awseks.AwsAuthMapping{
		Username: jsii.String(eksClusterConfig.UserId),
		Groups: &[]*string{
			jsii.String("system:masters"),
		},
	})

	// creates the IAM role for node group with the provided policies
	nodeGroupRole := awsiam.NewRole(stack, jsii.String(nodeGroupRoleConfig.Name), &awsiam.RoleProps{
		AssumedBy:   awsiam.NewServicePrincipal(jsii.String(nodeGroupRoleConfig.AssumedBy), nil),
		Description: jsii.String("Role for EKS Node Group"),
		ManagedPolicies: &[]awsiam.IManagedPolicy{
			awsiam.ManagedPolicy_FromAwsManagedPolicyName(jsii.String("AmazonEC2ContainerRegistryReadOnly")),
			awsiam.ManagedPolicy_FromAwsManagedPolicyName(jsii.String("AmazonEC2FullAccess")),
			awsiam.ManagedPolicy_FromAwsManagedPolicyName(jsii.String("AmazonEKS_CNI_Policy")),
			awsiam.ManagedPolicy_FromAwsManagedPolicyName(jsii.String("AmazonEKSWorkerNodePolicy")),
		},
		RoleName: jsii.String(nodeGroupRoleConfig.Name),
	})

	// creates new node group
	eksCluster.AddNodegroupCapacity(jsii.String(nodeGroupConfig.Name), &awseks.NodegroupOptions{
		InstanceTypes: &[]awsec2.InstanceType{
			awsec2.NewInstanceType(jsii.String(nodeGroupConfig.InstanceTypes)),
		},
		NodeRole:    nodeGroupRole,
		MinSize:     jsii.Number(nodeGroupConfig.MinSize),
		MaxSize:     jsii.Number(nodeGroupConfig.MaxSize),
		DesiredSize: jsii.Number(nodeGroupConfig.DesiredSize),
		RemoteAccess: &awseks.NodegroupRemoteAccess{
			SshKeyName: jsii.String(nodeGroupConfig.SshKeyName),
		},
		Subnets: &awsec2.SubnetSelection{
			Subnets: vpc.PublicSubnets(),
		},
		AmiType:  nodeGroupConfig.AmiType,
		DiskSize: jsii.Number(nodeGroupConfig.DiskSize),
	})

	// Addons
	awseks.NewCfnAddon(stack, jsii.String("VPCCNIAddon"), &awseks.CfnAddonProps{
		ClusterName:      eksCluster.ClusterName(),
		AddonName:        jsii.String("vpc-cni"),
		ResolveConflicts: jsii.String("OVERWRITE"),
	})

	awseks.NewCfnAddon(stack, jsii.String("CoreDNSAddon"), &awseks.CfnAddonProps{
		ClusterName:      eksCluster.ClusterName(),
		AddonName:        jsii.String("coredns"),
		ResolveConflicts: jsii.String("OVERWRITE"),
	})

	awseks.NewCfnAddon(stack, jsii.String("KubeProxyAddon"), &awseks.CfnAddonProps{
		ClusterName:      eksCluster.ClusterName(),
		AddonName:        jsii.String("kube-proxy"),
		ResolveConflicts: jsii.String("OVERWRITE"),
	})

	awseks.NewCfnAddon(stack, jsii.String("PodIdentityAddon"), &awseks.CfnAddonProps{
		ClusterName:      eksCluster.ClusterName(),
		AddonName:        jsii.String("eks-pod-identity-agent"),
		ResolveConflicts: jsii.String("OVERWRITE"),
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

	parseFlags()

	app := awscdk.NewApp(nil)

	NewCdkIvanStack(app, cdkStackId, &CdkStackProps{
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
		Account: jsii.String(accountId),
		Region:  jsii.String(region),
	}
}

func parseFlags() {

	// VPC config flags
	flag.StringVar(&vpcConfig.Name, "vpc-name", vpcConfig.Name, "VPC name")
	flag.StringVar(&vpcConfig.Cidr, "vpc-cidr", vpcConfig.Cidr, "VPC CIDR")
	flag.IntVar(&vpcConfig.MaxAZs, "max-azs", vpcConfig.MaxAZs, "Maximum number of AZs")
	flag.StringVar(&vpcConfig.PublicSubnetName, "public-subnet-name", vpcConfig.PublicSubnetName, "Public Subnet Name")
	flag.StringVar(&vpcConfig.PrivateSubnetName, "private-subnet-name", vpcConfig.PrivateSubnetName, "Private Subnet Name")
	flag.IntVar(&vpcConfig.PublicMask, "public-mask", vpcConfig.PublicMask, "Public Subnet Mask")
	flag.IntVar(&vpcConfig.PrivateMask, "private-mask", vpcConfig.PrivateMask, "Private Subnet Mask")

	// Route table config flags
	flag.StringVar(&routeTableConfig.Name, "rt-name", routeTableConfig.Name, "Route Table Name")
	flag.StringVar(&routeTableConfig.TagsKey, "rt-tags-key", routeTableConfig.TagsKey, "Route Table Key in Tags")
	flag.StringVar(&routeTableConfig.TagsValue, "rt-tags-value", routeTableConfig.TagsValue, "Route Table Value in Tags")

	// Route config flags
	flag.StringVar(&routeConfig.Name, "route-name", routeConfig.Name, "Route Name")
	flag.StringVar(&routeConfig.DestinationCidrBlock, "route-cidr-block", routeConfig.DestinationCidrBlock, "Destination CIDR Block for Route")

	// Security group config flags
	flag.StringVar(&securityGroupConfig.Name, "sg-name", securityGroupConfig.Name, "Security Group Name")

	// Cluster role config flags
	flag.StringVar(&clusterRoleConfig.Name, "cluster-role-name", clusterRoleConfig.Name, "Cluster Role Name")
	flag.StringVar(&clusterRoleConfig.AssumedBy, "cluster-role-assumed-by", clusterRoleConfig.AssumedBy, "IAM Principal for Cluster Role")

	// EKS cluster config flags
	flag.StringVar(&eksClusterConfig.Name, "eks-cluster-name", eksClusterConfig.Name, "EKS Cluster Name")
	flag.StringVar(&eksClusterConfig.K8SVersion, "k8s-version", eksClusterConfig.K8SVersion, "Kubernetes Version for EKS Cluster")
	flag.IntVar(&eksClusterConfig.DefaultCapacity, "default-capacity", eksClusterConfig.DefaultCapacity, "Default Capacity for EKS Cluster")
	flag.StringVar(&eksClusterConfig.UserId, "eks-user-id", eksClusterConfig.UserId, "User ID for AWS Auth")
	flag.StringVar(&eksClusterConfig.UserArn, "eks-user-arn", eksClusterConfig.UserArn, "User ARN for AWS Auth")

	// Node group role config flags
	flag.StringVar(&nodeGroupRoleConfig.Name, "node-group-role-name", nodeGroupRoleConfig.Name, "Node Group Role Name")
	flag.StringVar(&nodeGroupRoleConfig.AssumedBy, "node-group-role-assumed-by", nodeGroupRoleConfig.AssumedBy, "IAM Principal for Node Group Role")

	// Node group config flags
	flag.StringVar(&nodeGroupConfig.Name, "node-group-name", nodeGroupConfig.Name, "Node Group Name")
	flag.StringVar(&nodeGroupConfig.InstanceTypes, "node-group-instance-types", nodeGroupConfig.InstanceTypes, "Instance Types for Node Group")
	flag.IntVar(&nodeGroupConfig.MinSize, "node-group-min-size", nodeGroupConfig.MinSize, "Minimum Size for Node Group")
	flag.IntVar(&nodeGroupConfig.MaxSize, "node-group-max-size", nodeGroupConfig.MaxSize, "Maximum Size for Node Group")
	flag.IntVar(&nodeGroupConfig.DesiredSize, "node-group-desired-size", nodeGroupConfig.DesiredSize, "Desired Size for Node Group")
	flag.StringVar(&nodeGroupConfig.SshKeyName, "ssh-key-name", nodeGroupConfig.SshKeyName, "SSH Key Name for Node Group")
	flag.IntVar(&nodeGroupConfig.DiskSize, "node-group-disk-size", nodeGroupConfig.DiskSize, "Disk Size for Node Group")

	// General stack and AWS settings
	flag.StringVar(&cdkStackId, "cdk-stack-id", cdkStackId, "Id of CDK Stack")
	flag.StringVar(&accountId, "account-id", accountId, "Id of AWS account")
	flag.StringVar(&region, "region", region, "Region in AWS for stack")

	flag.Parse()
}
