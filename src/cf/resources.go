package cf

import (
	_ "embed" // required for embed
)

//go:embed resources/aws_sns_topic.template
var awsSNSTopic []byte

//go:embed resources/aws_sns_topic_subscription.template
var awsSNSSubscription []byte

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

//go:embed resources/aws_iam_managed_policy.template
var awsIamManagedPolicy []byte

//go:embed resources/aws_iam_policy.template
var awsIamPolicy []byte

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

//go:embed resources/aws_lb.template
var awsLb []byte

//go:embed resources/aws_lb_listener_rule.template
var awsLbListenerRule []byte

//go:embed resources/aws_lb_target_group.template
var awsLbTargetGroup []byte

//go:embed resources/aws_lb_listener.template
var awsLbListener []byte

//go:embed resources/aws_cloudfront_distribution.template
var awsCloudfrontDistribution []byte

//go:embed resources/aws_iam_user.template
var awsIamUser []byte

//go:embed resources/aws_cloud9_environment_ec2.template
var awsCloud9EnvironmentEc2 []byte

//go:embed resources/aws_codecommit_repository.template
var awsCodecommitRepository []byte

//go:embed resources/aws_cloudwatch_metric_alarm.template
var awsCloudwatchMetricAlarm []byte

//go:embed resources/aws_route53_record.template
var awsRoute53Record []byte

//go:embed resources/aws_db_subnet_group.template
var awsDBSubnetGroup []byte

//go:embed resources/aws_rds_cluster.template
var awsRdsCluster []byte

//go:embed resources/aws_cloudwatch_event_rule.template
var awsCloudwatchEventRule []byte

//go:embed resources/aws_cloudfront_origin_access_identity.template
var awsCloudfrontOriginAccessIdentity []byte

//go:embed resources/aws_elb.template
var awsElb []byte

//go:embed resources/aws_autoscaling_policy.template
var awsAutoscalingPolicy []byte

//go:embed resources/aws_autoscaling_schedule.template
var awsAutoscalingSchedule []byte

//go:embed resources/aws_ebs_volume.template
var awsEbsVolume []byte

//go:embed resources/aws_sns_topic_policy.template
var awsSNSTopicPolicy []byte

//go:embed resources/aws_config_delivery_channel.template
var awsConfigDeliveryChannel []byte

//go:embed resources/aws_config_configuration_recorder.template
var awsConfigConfigurationRecorder []byte

//go:embed resources/aws_config_config_rule.template
var awsConfigConfigRule []byte

//go:embed resources/aws_dms_replication_task.template
var awsDmsReplicationTask []byte

//go:embed resources/aws_db_parameter_group.template
var awsDBParameterGroup []byte

//go:embed resources/aws_dms_endpoint.template
var awsDmsEndpoint []byte

//go:embed resources/aws_dms_replication_instance.template
var awsDmsReplicationinstance []byte

//go:embed resources/aws_dms_replication_subnet_group.template
var awsDmsReplicationSubnetGroup []byte

//go:embed resources/aws_db_instance.template
var awsDBInstance []byte

//go:embed resources/aws_eip_association.template
var awsEipAssociation []byte

//go:embed resources/aws_network_interface.template
var awsNetworkInterface []byte

//go:embed resources/aws_appautoscaling_policy.template
var awsAppAutoscalingPolicy []byte

//go:embed resources/aws_appautoscaling_target.template
var awsAppAutoscalingTarget []byte

//go:embed resources/aws_ecs_task_definition.template
var awsEcsTaskDefinition []byte

//go:embed resources/aws_ecs_cluster.template
var awsEcsCluster []byte

//go:embed resources/aws_ecs_service.template
var awsEcsService []byte

//go:embed resources/aws_efs_file_system.template
var awsEfsFileSystem []byte

//go:embed resources/aws_efs_mount_target.template
var awsEfsMountTarget []byte

//go:embed resources/aws_eks_node_group.template
var awsEksNodeGroup []byte

//go:embed resources/aws_launch_template.template
var awsLaunchTemplate []byte

//go:embed resources/aws_eks_cluster.template
var awsEksCluster []byte

//go:embed resources/aws_sqs_queue.template
var awsSqsQueue []byte

//go:embed resources/aws_neptune_cluster.template
var awsNeptuneCluster []byte

//go:embed resources/aws_neptune_cluster_instance.template
var awsNeptuneDBInstance []byte

//go:embed resources/aws_neptune_parameter_group.template
var awsNeptuneDBParameterGroup []byte

//go:embed resources/aws_neptune_cluster_parameter_group.template
var awsNeptuneClusterDBParameterGroup []byte

//go:embed resources/aws_neptune_subnet_group.template
var awsNeptuneDnSubnetGroup []byte

//go:embed resources/aws_iam_group.template
var awsIamGroup []byte

//go:embed resources/aws_iam_access_key.template
var awsIamAccessKey []byte

//go:embed resources/aws_iam_group_membership.template
var awsIamGroupMembership []byte

//go:embed resources/aws_servicecatalog_tag_option_resource_association.template
var awsServiceCatalogTagOptionAssociation []byte

//go:embed resources/aws_servicecatalog_product_portfolio_association.template
var awsServiceCatalogProductPortfolioAssociation []byte

//go:embed resources/aws_servicecatalog_portfolio.template
var awsServiceCatalogPortfolio []byte

//go:embed resources/aws_servicecatalog_portfolio_share.template
var awsServiceCatalogPortfolioShare []byte

//go:embed resources/aws_servicecatalog_tag_option.template
var awsServiceCatalogTagOption []byte

//go:embed resources/aws_backup_plan.template
var awsBackupPlan []byte

//go:embed resources/aws_backup_selection.template
var awsBackupSelection []byte

//go:embed resources/aws_backup_vault.template
var awsBackupVault []byte

//go:embed resources/aws_ssm_maintenance_window_task.template
var awsSsmMaintenanceWindowTask []byte

//go:embed resources/aws_ssm_maintenance_window_target.template
var awsSsmMaintenanceWindowTarget []byte

//go:embed resources/aws_wafv2_web_acl_association.template
var awsWAFv2WebACLAssociation []byte

//go:embed resources/aws_autoscaling_lifecycle_hook.template
var awsAutoscalingLifecycleHook []byte

//go:embed resources/aws_ssm_maintenance_window.template
var awsSsmMaintenanceWindow []byte

//go:embed resources/aws_lambda_version.template
var awsLambdaVersion []byte

//go:embed resources/aws_lambda_event_source_mapping.template
var awsLambdaEventSourceMapping []byte

//go:embed resources/aws_service_discovery_service.template
var awsServiceDiscoveryService []byte

//go:embed resources/aws_cloudwatch_dashboard.template
var awsCloudwatchDashboard []byte

//go:embed resources/aws_cloudwatch_log_metric_filter.template
var awsCloudwatchLogMetricFilter []byte

//go:embed resources/aws_athena_workgroup.template
var awsAthenaWorkGroup []byte

//go:embed resources/aws_athena_named_query.template
var awsAthenaNamedQuery []byte

//go:embed resources/aws_kinesis_firehose_delivery_stream.template
var awsKinesisFirehoseDeliveryStream []byte
