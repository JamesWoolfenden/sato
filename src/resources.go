package sato

import (
	_ "embed" //required for embed
)

//go:embed resources/aws_sns_topic.template
var awsSNSTopic []byte

//go:embed resources/aws_iam_role.template
var awsIamRole []byte

//go:embed resources/aws_route_table_association.template
var awsRouteTableAssociation []byte

//go:embed resources/aws_route_table.template
var awsRouteTable []byte

//go:embed resources/aws_route.template
var awsRoute []byte

//go:embed resources/aws_nat_gateway.template
var awsNatGateway []byte

//go:embed resources/aws_vpn_gateway_attachment.template
var awsVpnGatewayAttachment []byte

//go:embed resources/aws_subnet.template
var awsSubnet []byte

//go:embed resources/aws_network_acl_rule.template
var awsNetworkACLRule []byte

//go:embed resources/aws_network_acl.template
var awsNetworkACL []byte

//go:embed resources/aws_eip.template
var awsEIP []byte

//go:embed resources/aws_vpc.template
var awsVpc []byte

//go:embed resources/aws_cloudwatch_log_group.template
var awsCloudwatchLogGroup []byte

//go:embed resources/aws_vpc_dhcp_options_association.template
var awsVpcDhcpOptionsAssociation []byte

//go:embed resources/aws_vpc_dhcp_options.template
var awsVpcDhcpOptions []byte

//go:embed resources/aws_network_acl_association.template
var awsNetworkACLAssociation []byte

//go:embed resources/aws_vpc_endpoint.template
var awsVpcEndpoint []byte

//go:embed resources/aws_flow_log.template
var awsFlowLog []byte

//go:embed resources/aws_internet_gateway.template
var awsInternetGateway []byte

//go:embed resources/aws_s3_bucket.template
var awsS3Bucket []byte

//go:embed resources/aws_lambda_function.template
var awsLambdaFunction []byte

//go:embed resources/aws_sfn_state_machine.template
var awsStepfunctionStateMachine []byte

//go:embed resources/aws_dynamodb_table.template
var awsDynamodbTable []byte

//go:embed resources/aws_iam_instance_profile.template
var awsIamInstanceProfile []byte

//go:embed resources/aws_cloudformation_stack.template
var awsCloudformationStack []byte

//go:embed resources/aws_secretsmanager_secret.template
var awsSecretsManagerSecret []byte

//go:embed resources/aws_security_group.template
var awsSecurityGroup []byte

//go:embed resources/aws_instance.template
var awsInstance []byte

//go:embed resources/aws_kms_key.template
var awsKmsKey []byte

//go:embed resources/aws_kms_alias.template
var awskmsAlias []byte

//go:embed resources/aws_s3_bucket_policy.template
var awsS3BucketPolicy []byte

//go:embed resources/aws_iam_policy.template
var awsIamManagedPolicy []byte

//go:embed resources/aws_ssm_association.template
var awsSsmAssociation []byte

//go:embed resources/aws_ssm_document.template
var awsSsmDocument []byte

//go:embed resources/aws_autoscaling_group.template
var awsAutoscalingGroup []byte

//go:embed resources/aws_launch_configuration.template
var awsLaunchConfiguration []byte

//go:embed resources/aws_lambda_permission.template
var awsLambdaPermission []byte

//go:embed resources/aws_elasticache_parameter_group.template
var awsElasticacheParameterGroup []byte

//go:embed resources/aws_elasticache_replication_group.template
var awsElasticacheReplicationGroup []byte

//go:embed resources/aws_elasticache_subnet_group.template
var awsElasticacheSubnetGroup []byte

//go:embed resources/aws_directory_service_directory.template
var awsDirectoryServiceDirectory []byte

//go:embed resources/aws_codebuild_project.template
var awsCodebuildProject []byte

//go:embed resources/aws_codepipeline.template
var awsCodepipeline []byte

//go:embed resources/aws_security_group_rule_ingress.template
var awsSecurityGroupRuleIngress []byte

//go:embed resources/aws_security_group_rule_egress.template
var awsSecurityGroupRuleEgress []byte
