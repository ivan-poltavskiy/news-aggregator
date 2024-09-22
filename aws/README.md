This CloudFormation template deploys an Amazon EKS cluster along with the required AWS
services such as a VPC, subnets, IAM roles, security groups, and 
several EKS addons.

The following services are deployed through this CloudFormation template:

- VPC: Virtual Private Cloud for networking.
- Public Subnets: Two public subnets across different availability zones.
- Internet Gateway: To provide access to the internet for resources within the public subnets.
- Route Table and Route: For routing traffic through the internet gateway.
- Security Group: Manages inbound traffic for EKS nodes and control plane (HTTP, HTTPS, SSH).
- IAM Roles: For EKS Cluster, Node Group, and related policies.
- EKS Cluster: The core of the Kubernetes cluster.
- EKS Node Group: A scalable group of EC2 worker nodes to run your workloads.
- Addons

`Before deploying this stack, ensure that you have:`

- An AWS account with the necessary permissions to create the resources listed above.

- A configured EC2 SSH key pair for SSH access to worker nodes.

- AWS CLI installed and configured with proper access credentials.

- A default region set up in your AWS configuration.
