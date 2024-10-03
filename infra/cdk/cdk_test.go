package main

import (
	"github.com/aws/aws-cdk-go/awscdk/v2/assertions"
	"github.com/aws/jsii-runtime-go"
	"testing"

	"github.com/aws/aws-cdk-go/awscdk/v2"
)

func TestCdkIvanStack(t *testing.T) {
	// GIVEN

	app := awscdk.NewApp(nil)

	// WHEN
	stack := NewCdkIvanStack(app, "MyStack", &CdkStackProps{})

	//THEN
	template := assertions.Template_FromStack(stack, nil)

	template.HasResourceProperties(jsii.String("AWS::EC2::VPC"), map[string]interface{}{
		"CidrBlock": "10.0.0.0/16",
	})

	template.HasResourceProperties(jsii.String("AWS::EC2::Subnet"), map[string]interface{}{
		"CidrBlock":           "10.0.0.0/20",
		"MapPublicIpOnLaunch": true,
		"Tags": assertions.Match_ArrayWith(&[]interface{}{
			map[string]interface{}{
				"Key":   "aws-cdk:subnet-name",
				"Value": "ivan-cdk-public-subnet-1",
			},
		}),
	})

	template.HasResourceProperties(jsii.String("AWS::EC2::SecurityGroup"), map[string]interface{}{
		"GroupDescription": "MyStack/ivan-cdk-sg",
		"GroupName":        "ivan-cdk-sg",
	})

	template.HasResourceProperties(jsii.String("Custom::AWSCDK-EKS-Cluster"), map[string]interface{}{})

	template.HasResourceProperties(jsii.String("AWS::IAM::Role"), map[string]interface{}{
		"AssumeRolePolicyDocument": map[string]interface{}{
			"Statement": assertions.Match_ArrayWith(&[]interface{}{
				map[string]interface{}{
					"Action": "sts:AssumeRole",
					"Effect": "Allow",
					"Principal": map[string]interface{}{
						"Service": "eks.amazonaws.com",
					},
				},
			}),
		},
	})

	template.HasResourceProperties(jsii.String("AWS::EKS::Nodegroup"), map[string]interface{}{
		"DiskSize":      20,
		"InstanceTypes": assertions.Match_ArrayWith(&[]interface{}{"t2.medium"}),
		"ScalingConfig": map[string]interface{}{
			"MinSize":     1,
			"MaxSize":     10,
			"DesiredSize": 2,
		},
	})

	template.HasResourceProperties(jsii.String("AWS::EC2::RouteTable"), map[string]interface{}{
		"VpcId": map[string]interface{}{
			"Ref": "ivancdkvpcE8FEAA9F",
		},
	})

	template.HasResourceProperties(jsii.String("AWS::EC2::Route"), map[string]interface{}{
		"DestinationCidrBlock": "0.0.0.0/0",
	})
}
