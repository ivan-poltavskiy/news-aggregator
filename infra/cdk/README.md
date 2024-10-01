This is a CDK template for deploying an EKS cluster and additional services in AWS.
A list of all services used:
- VPC
- Public and Private Subnets
- Route Table
- Route
- Security Group (along with inbound rules for HTTP, HTTPS and SSH)
- IAM roles for EKS Cluster and Node Group
- Node group
- Addons

For deploying it to the stack, first of all execute the `cdk bootstrap` command for creating additional 
and necessary resources. 

After it, by using the command `cdk diff`, a comparison deployed stack with current state will be displayed.

Use the `cdk deploy` for deploying the stack to your default AWS account with region, which provided like default region in your AWS configuration.

## Useful commands

 * `cdk deploy`      deploy this stack to your default AWS account/region
 * `cdk diff`        compare deployed stack with current state
 * `cdk synth`       emits the synthesized CloudFormation template
 * `go test`         run unit tests
