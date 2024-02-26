# Terraform Configuration for SMS Receive API

This Terraform configuration sets up an AWS-based infrastructure
for a messaging application. It includes resources for networking,
database, AWS Lambda, API Gateway, and logging. Below is a detailed
overview of each main resource and how to set them up.

## Main Resources

### AWS VPC and Subnets

- `aws_vpc.my_vpc`: Creates a Virtual Private Cloud (VPC) witha CIDR
- block of `10.0.0.0/16`. This VPC acts as a virtual network
for your AWS resources.
- `aws_subnet.subnet_1` and `aws_subnet.subnet_2`: Two subnets
are created within the VPC in different availability
zones (`us-west-2a` and `us-west-2b`) for high availability and
fault tolerance.

### RDS MySQL Instance

- `aws_db_instance.mysql`: Deploys an RDS MySQL instance. It is
publicly accessible and resides within the created VPC and subnets.

### AWS Lambda Function

- `aws_lambda_function.receive_sms_lambda`: A Lambda function setup
to handle incoming SMS messages. It's integrated with the VPC
and utilizes environment variables for configuration.

### API Gateway

- `aws_api_gateway_rest_api.api_gateway`: Creates an API Gateway to
expose HTTP endpoints.
- `aws_api_gateway_resource.api_resource` and
`aws_api_gateway_method.post_method`: Configures a resource and method
for receiving SMS messages via HTTP POST requests.
- `aws_api_gateway_integration.lambda_integration`: Integrates the
API Gateway with the AWS Lambda function.

### Network Load Balancer and VPC Link

- `aws_lb.lb`: A Network Load Balancer to route external traffic.
- `aws_api_gateway_vpc_link.my_vpc_link`: Links the API Gateway with
the VPC through the NLB.

### Security Groups and Internet Gateway

- `aws_security_group.mysql_sg`: Security group for the RDS instance,
allowing inbound MySQL traffic.
- `aws_internet_gateway.equilibria_internet_gateway`: Internet Gateway
attached to the VPC for internet access.

### CloudWatch Log Group

- `aws_cloudwatch_log_group.api_gateway_logs`: Log group for storing
logs from the API Gateway.

### IAM Roles and Policies

- `aws_iam_role.lambda_execution_role`: Execution role for the Lambda
function with necessary permissions.
- `aws_iam_role_policy.lambda_vpc_policy`: Policy granting the Lambda
function permissions to work within a VPC.
- `aws_iam_role.api_gateway_cloudwatch_role` and
`aws_iam_role_policy.api_gateway_cloudwatch_policy`: Role and policy
for API Gateway to log to CloudWatch.

### SSM Parameters

- Multiple `aws_ssm_parameter` resources: These parameters store
configuration and secrets like database credentials and API keys
securely in AWS Systems Manager Parameter Store.

## Setting Up Environment Variables

Environment variables should be set up as per the `variable` definitions
in the Terraform files. See `/terraform/.tfenv` for more details on
setting these variables.

## Outputs

- `db_endpoint`, `db_port`, `db_name`: Output the connection details
of the RDS instance.
- `api_invoke_url`: Provides the URL to invoke the Lambda function
via the API Gateway.

## Instructions for Use

1. **Set Environment Variables**: Configure your environment variables
as per the requirements in `/terraform/.tfenv`.
2. **Initialize Terraform**: Run `terraform init` in your Terraform
directory to initialize the workspace.
3. **Apply Configuration**: Execute `terraform apply` to create there
sources as per the configuration. Confirm the action by typing
`yes` when prompted.
4. **Access Outputs**: After successful Terraform execution, access
the outputs to get information like the RDS endpoint and API Gateway
URL.

## Important Notes

- **Security**: Review and restrict security group rules as necessary,
especially for the RDS instance.
- **Public Accessibility**: Be cautious with publicly accessible resources.
Ensure strong passwords are used and consider limiting access to specific IPs.
- **Terraform State**: Manage your Terraform state carefully, especially in a
team environment or when working with CI/CD pipelines.
