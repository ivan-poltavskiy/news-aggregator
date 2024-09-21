package main

import (
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

func NewCdkStack(scope constructs.Construct, id string, props *CdkStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	return stack
}

func NewCdkIvanStack(scope constructs.Construct, id string, props *CdkStackProps) awscdk.Stack {
	stack := awscdk.NewStack(scope, &id, &props.StackProps)

	vpc := awsec2.NewVpc(stack, jsii.String("ivan-cdk-vpc"), &awsec2.VpcProps{
		IpAddresses: awsec2.IpAddresses_Cidr(jsii.String("10.0.0.0/16")),
		MaxAzs:      jsii.Number(2),
		SubnetConfiguration: &[]*awsec2.SubnetConfiguration{
			{
				Name:                jsii.String("ivan-cdk-public-subnet-1"),
				SubnetType:          awsec2.SubnetType_PUBLIC,
				CidrMask:            jsii.Number(20),
				MapPublicIpOnLaunch: jsii.Bool(true),
			},
			{
				Name:       jsii.String("ivan-cdk-private-subnet-1"),
				SubnetType: awsec2.SubnetType_PRIVATE_WITH_EGRESS,
				CidrMask:   jsii.Number(20),
			},
		},
		VpcName: jsii.String("ivan-cdk-vpc"),
	})

	routeTable := awsec2.NewCfnRouteTable(stack, jsii.String("ivan-cdk-rt"), &awsec2.CfnRouteTableProps{
		VpcId: vpc.VpcId(),
		Tags: &[]*awscdk.CfnTag{
			{
				Key:   jsii.String("name"),
				Value: jsii.String("ivan-cdk-rt"),
			},
		},
	})

	// Add a route to the route table that directs traffic to the internet gateway
	awsec2.NewCfnRoute(stack, jsii.String("ivan-cdk-route"), &awsec2.CfnRouteProps{
		RouteTableId:         routeTable.AttrRouteTableId(),
		DestinationCidrBlock: jsii.String("0.0.0.0/0"),
		GatewayId:            vpc.InternetGatewayId(),
	})

	publicSubnets := vpc.PublicSubnets()
	for i, subnet := range *publicSubnets {
		awsec2.NewCfnSubnetRouteTableAssociation(stack, jsii.String(fmt.Sprintf("SubnetAssociation%d", i)), &awsec2.CfnSubnetRouteTableAssociationProps{
			SubnetId:     subnet.SubnetId(),
			RouteTableId: routeTable.Ref(),
		})
	}
	securityGroup := awsec2.NewSecurityGroup(stack, jsii.String("ivan-cdk-sg"), &awsec2.SecurityGroupProps{
		Vpc:               vpc,
		SecurityGroupName: jsii.String("ivan-cdk-sg"),
		AllowAllOutbound:  jsii.Bool(true),
	})

	securityGroup.AddIngressRule(awsec2.Peer_AnyIpv4(), awsec2.Port_Tcp(jsii.Number(80)), jsii.String("HTTP"), jsii.Bool(true))
	securityGroup.AddIngressRule(awsec2.Peer_AnyIpv4(), awsec2.Port_Tcp(jsii.Number(443)), jsii.String("HTTPS"), jsii.Bool(true))
	securityGroup.AddIngressRule(awsec2.Peer_AnyIpv4(), awsec2.Port_Tcp(jsii.Number(22)), jsii.String("SSH"), jsii.Bool(true))

	clusterRole := awsiam.NewRole(stack, jsii.String("eks-cdk-cluster-role"), &awsiam.RoleProps{
		AssumedBy: awsiam.NewServicePrincipal(jsii.String("eks.amazonaws.com"), nil),
		ManagedPolicies: &[]awsiam.IManagedPolicy{
			awsiam.ManagedPolicy_FromAwsManagedPolicyName(jsii.String("AmazonEKSClusterPolicy")),
			awsiam.ManagedPolicy_FromAwsManagedPolicyName(jsii.String("AmazonEKSServicePolicy")),
		},
	})

	eksCluster := awseks.NewCluster(stack, jsii.String("eks-cdk-cluster"), &awseks.ClusterProps{
		Version:         awseks.KubernetesVersion_Of(jsii.String("1.30")),
		ClusterName:     jsii.String("eks-cdk-cluster"),
		Vpc:             vpc,
		SecurityGroup:   securityGroup,
		DefaultCapacity: jsii.Number(0),
		Role:            clusterRole,
	})

	eksCluster.AwsAuth().AddUserMapping(awsiam.User_FromUserArn(stack, jsii.String("ivan"), jsii.String("arn:aws:iam::406477933661:user/ivan")), &awseks.AwsAuthMapping{
		Username: jsii.String("ivan"),
		Groups: &[]*string{
			jsii.String("system:masters"),
		},
	})

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

	eksCluster.AddNodegroupCapacity(jsii.String("node-group"), &awseks.NodegroupOptions{
		InstanceTypes: &[]awsec2.InstanceType{
			awsec2.NewInstanceType(jsii.String("t2.medium")),
		},
		NodeRole:    nodeGroupRole,
		MinSize:     jsii.Number(1),
		MaxSize:     jsii.Number(10),
		DesiredSize: jsii.Number(2),
		RemoteAccess: &awseks.NodegroupRemoteAccess{
			SshKeyName: jsii.String("ivan-kp"),
		},
		Subnets: &awsec2.SubnetSelection{
			Subnets: vpc.PublicSubnets(),
		},
		AmiType:  awseks.NodegroupAmiType_AL2_X86_64,
		DiskSize: jsii.Number(20),
	})

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
		Account: jsii.String("406477933661"),
		Region:  jsii.String("eu-west-2"),
	}
}
