package see

import (
	"fmt"
	"strings"
)

// Lookup converts from cf/arm to terraform resource name.
func Lookup(resource string, reverse bool) (*string, error) {
	var result string

	Lookup := map[string]string{
		"aws::EFS::filesystem":                                                "aws_efs_file_system",
		"aws::EFS::mounttarget":                                               "aws_efs_mount_target",
		"aws::EKS::cluster":                                                   "aws_eks_cluster",
		"aws::EKS::nodegroup":                                                 "aws_eks_node_group",
		"aws::applicationautoscaling::scalabletarget":                         "aws_appautoscaling_target",
		"aws::applicationautoscaling::scalingpolicy":                          "aws_appAutoscaling_policy",
		"aws::autoscaling::autoscalinggroup":                                  "aws_autoscaling_group",
		"aws::autoscaling::launchconfiguration":                               "aws_launch_configuration",
		"aws::autoscaling::lifecyclehook":                                     "aws_autoscaling_lifecycle_hook",
		"aws::autoscaling::scalingpolicy":                                     "aws_autoscaling_policy",
		"aws::autoscaling::scheduledaction":                                   "aws_autoscaling_schedule",
		"aws::backup::backupplan":                                             "aws_backup_plan",
		"aws::backup::backupselection":                                        "aws_backup_selection",
		"aws::backup::backupvault":                                            "aws_backup_vault",
		"aws::cloud9::environmentec2":                                         "aws_cloud9_environment_ec2",
		"aws::cloudformation::Stack":                                          "aws_cloudformation_stack",
		"aws::cloudfront::cloudFrontOriginAccessIdentity":                     "aws_cloudfront_origin_access_identity",
		"aws::cloudfront::distribution":                                       "aws_cloudfront_distribution",
		"aws::cloudwatch::alarm":                                              "aws_cloudwatch_metric_alarm",
		"aws::cloudwatch::dashboard":                                          "aws_cloudwatch_dashboard",
		"aws::codebuild::project":                                             "aws_codebuild_project",
		"aws::codecommit::repository":                                         "aws_codecommit_repository",
		"aws::codepipeline::pipeline":                                         "aws_codepipeline",
		"aws::config::configrule":                                             "aws_config_config_rule",
		"aws::config::configurationrecorder":                                  "aws_config_configuration_recorder",
		"aws::config::deliverychannel":                                        "aws_config_delivery_channel",
		"aws::directoryservice::microsoftad":                                  "aws_directory_service_directory",
		"aws::dns::endpoint":                                                  "aws_dms_endpoint",
		"aws::dns::replicationinstance":                                       "aws_dms_replication_instance",
		"aws::dns::replicationsubnetgroup":                                    "aws_dms_replication_subnet_group",
		"aws::dns::replicationtask":                                           "aws_dms_replication_task",
		"aws::dynamodb::table":                                                "aws_dynamodb_table",
		"aws::ec2::dhcpoptions":                                               "aws_vpc_dhcp_options",
		"aws::ec2::eip":                                                       "aws_eip",
		"aws::ec2::eipassociation":                                            "aws_eip_association",
		"aws::ec2::flowlog":                                                   "aws_flow_log",
		"aws::ec2::instance":                                                  "aws_instance",
		"aws::ec2::internetgateway":                                           "aws_Internet_gateway",
		"aws::ec2::launchtemplate":                                            "aws_launch_template",
		"aws::ec2::natgateway":                                                "aws_nat_gateway",
		"aws::ec2::networkacl":                                                "aws_network_acl",
		"aws::ec2::networkaclentry":                                           "aws_network_acl_rule",
		"aws::ec2::networkinterface":                                          "aws_network_interface",
		"aws::ec2::route":                                                     "aws_route",
		"aws::ec2::routetable":                                                "aws_route_table",
		"aws::ec2::securitygroup":                                             "aws_security_group",
		"aws::ec2::securitygroupegress":                                       "aws_security_group_rule_egress",
		"aws::ec2::securitygroupingress":                                      "aws_security_group_rule_ingress",
		"aws::ec2::subnet":                                                    "aws_subnet",
		"aws::ec2::subnetnetworkaclassociation":                               "aws_network_acl_association",
		"aws::ec2::subnetroutetableassociation":                               "aws_route_table_association",
		"aws::ec2::volume":                                                    "aws_ebs_volume",
		"aws::ec2::vpc":                                                       "aws_vpc",
		"aws::ec2::vpcdhcpoptionsassociation":                                 "aws_vpc_dhcp_options_association",
		"aws::ec2::vpcendpoint":                                               "aws_vpc_endpoint",
		"aws::ec2::vpcgatewayattachment":                                      "aws_vpn_gatewayAttachment",
		"aws::ecs::cluster":                                                   "aws_ecs_cluster",
		"aws::ecs::service":                                                   "aws_ecs_service",
		"aws::ecs::taskdefinition":                                            "aws_ecs_task_definition",
		"aws::elasticache::parametergroup":                                    "aws_elasticache_parameter_group",
		"aws::elasticache::replicationgroup":                                  "aws_elasticache_replication_group",
		"aws::elasticache::subnetgroup":                                       "aws_elasticache_subnet_group",
		"aws::elasticloadbalancing::loadbalancer":                             "aws_elb",
		"aws::elasticloadbalancingv2::listener":                               "aws_lb_listener",
		"aws::elasticloadbalancingv2::listenerrule":                           "aws_lb_listener_rule",
		"aws::elasticloadbalancingv2::loadbalancer":                           "aws_lb",
		"aws::elasticloadbalancingv2::targetgroup":                            "aws_lb_target_group",
		"aws::events::rule":                                                   "aws_cloudwatch_event_rule",
		"aws::iam::accessKey":                                                 "aws_iam_access_key",
		"aws::iam::group":                                                     "aws_iam_group",
		"aws::iam::instanceprofile":                                           "aws_iam_instance_profile",
		"aws::iam::managedpolicy":                                             "aws_iam_managed_policy",
		"aws::iam::policy":                                                    "aws_iam_policy",
		"aws::iam::role":                                                      "aws_iam_role",
		"aws::iam::user":                                                      "aws_iam_user",
		"aws::iam::usertogroupaddition":                                       "aws_iam_group_membership",
		"aws::kms::alias":                                                     "aws_kms_alias",
		"aws::kms::key":                                                       "aws_Kms_key",
		"aws::lambda::eventsourcemapping":                                     "aws_lambda_event_source_mapping",
		"aws::lambda::function":                                               "aws_lambda_function",
		"aws::lambda::permission":                                             "aws_lambda_permission",
		"aws::lambda::version":                                                "aws_lambda_version",
		"aws::logs::loggroup":                                                 "aws_cloudwatch_loggroup",
		"aws::logs::metricfilter":                                             "aws_cloudwatch_logMetricFilter",
		"aws::neptune::dbcluster":                                             "aws_neptune_cluster",
		"aws::neptune::dbclusterparametergroup":                               "aws_neptune_cluster_parameter_group",
		"aws::neptune::dbinstance":                                            "aws_neptune_cluster_instance",
		"aws::neptune::dbparametergroup":                                      "aws_neptune_parameter_group",
		"aws::neptune::dbsubnetgroup":                                         "aws_neptune_subnet_group",
		"aws::rds::dbcluster":                                                 "aws_rds_cluster",
		"aws::rds::dbclusterparametergroup":                                   "aws_rds_cluster_parameter_group",
		"aws::rds::dbinstance":                                                "aws_db_instance",
		"aws::rds::dbparametergroup":                                          "aws_db_parameter_group",
		"aws::rds::dbsubnetgroup":                                             "aws_db_subnet_group",
		"aws::route53::recordset":                                             "aws_route53_record",
		"aws::s3::bucket":                                                     "aws_s3_bucket",
		"aws::s3::bucketpolicy":                                               "aws_s3_bucket_policy",
		"aws::secretsmanager::secret":                                         "aws_secrets_manager_secret",
		"aws::servicecatalog::portfolio":                                      "aws_service_catalog_portfolio",
		"aws::servicecatalog::portfolioproductassociation":                    "aws_service_catalog_product_portfolio_association",
		"aws::servicecatalog::portfolioshare":                                 "aws_service_catalog_portfolio_share",
		"aws::servicecatalog::tagoption":                                      "aws_service_catalog_tag_option",
		"aws::servicecatalog::tagoptionassociation":                           "aws_service_catalog_tag_option_association",
		"aws::servicediscovery::service":                                      "aws_service_discovery_service",
		"aws::sns::subscription":                                              "aws_sns_subscription",
		"aws::sns::topic":                                                     "aws_sns_topic",
		"aws::sns::topicpolicy":                                               "aws_sns_topic_policy",
		"aws::sqs::queue":                                                     "aws_sqs_queue",
		"aws::ssm::association":                                               "aws_ssm_association",
		"aws::ssm::document":                                                  "aws_ssm_document",
		"aws::ssm::maintenancewindow":                                         "aws_ssm_maintenance_window",
		"aws::ssm::maintenancewindowtarget":                                   "aws_ssm_maintenance_window_target",
		"aws::ssm::maintenancewindowtask":                                     "aws_ssm_maintenance_window_task",
		"aws::stepfunctions::statemachine":                                    "aws_sfn_state_machine",
		"aws::wafv2::webaclassociation":                                       "aws_wafv2_webacl_association",
		"microsoft.aad/domainservices":                                        "azurerm_active_directory_domain_service",
		"microsoft.analysisservices/servers":                                  "azurerm_analysis_services_server",
		"microsoft.apimanagement/service":                                     "azurerm_api_management",
		"microsoft.app/containerapps":                                         "azurerm_container_app",
		"microsoft.app/managedenvironments":                                   "azurerm_container_app_environment",
		"microsoft.authorization/roleassignments":                             "azurerm_role_assignment",
		"microsoft.authorization/roledefinitions":                             "azurerm_role_definition",
		"microsoft.cognitiveservices/accounts":                                "azurerm_cognitive_account",
		"microsoft.compute/availabilitysets":                                  "azurerm_availability_set",
		"microsoft.compute/disks":                                             "azurerm_managed_disk",
		"microsoft.compute/virtualmachines":                                   "azurerm_virtual_machine",
		"microsoft.compute/virtualmachines/extensions":                        "azurerm_virtual_machine_extension",
		"microsoft.compute/virtualmachinescalesets":                           "azurerm_linux_virtual_machine_scale_set",
		"microsoft.containerregistry/registries":                              "azurerm_container_registry",
		"microsoft.containerservice/managedclusters":                          "azurerm_kubernetes_cluster",
		"microsoft.documentdb/databaseaccounts":                               "azurerm_cosmosdb_account",
		"microsoft.insights/activitylogalerts":                                "azurerm_monitor_activity_log_alert",
		"microsoft.keyvault/vaults":                                           "azurerm_key_vault",
		"microsoft.managedidentity/userassignedidentities":                    "azurerm_user_assigned_identity",
		"microsoft.network/applicationgateways":                               "azurerm_application_gateway",
		"microsoft.network/applicationgateways/authenticationcertificates":    "azurerm_application_gateway",
		"microsoft.network/applicationgateways/backendaddresspools":           "azurerm_network_interface_application_gateway_backend_address_pool_association",
		"microsoft.network/applicationgateways/backendhttpsettingscollection": "azurerm_application_gateway",
		"microsoft.network/applicationgateways/frontendipconfigurations":      "azurerm_application_gateway",
		"microsoft.network/applicationgateways/frontendports":                 "azurerm_application_gateway",
		"microsoft.network/applicationgateways/httplisteners":                 "azurerm_application_gateway",
		"microsoft.network/applicationgateways/sslcertificates":               "azurerm_application_gateway",
		"microsoft.network/applicationgatewaywebapplicationfirewallpolicies":  "azurerm_web_application_firewall_policy",
		"microsoft.network/bastionhosts":                                      "azurerm_bastion_host",
		"microsoft.network/networkinterfaces":                                 "azurerm_network_interface",
		"microsoft.network/networksecuritygroups":                             "azurerm_network_security_group",
		"microsoft.network/networksecuritygroups/securityrules":               "azurerm_network_security_rule",
		"microsoft.network/privatednszones":                                   "azurerm_private_dns_zone",
		"microsoft.network/privateendpoints":                                  "azurerm_private_endpoint",
		"microsoft.network/privateendpoints/privatednszonegroups":             "azurerm_private_endpoint",
		"microsoft.network/publicipaddresses":                                 "azurerm_public_ip",
		"microsoft.network/virtualnetworks":                                   "azurerm_virtual_network",
		"microsoft.network/virtualnetworks/subnets":                           "azurerm_subnet",
		"microsoft.operationalinsights/workspaces":                            "azurerm_log_analytics_workspace",
		"microsoft.operationsmanagement/solutions":                            "azurerm_log_analytics_solution",
		"microsoft.resources/deployments":                                     "azurerm_template_deployment",
		"microsoft.servicebus/namespaces":                                     "azurerm_servicebus_namespace",
		"microsoft.servicebus/namespaces/authorizationRules":                  "azurerm_servicebus_namespace_authorization_rule",
		"microsoft.servicebus/namespaces/queues":                              "azurerm_servicebus_queue",
		"microsoft.storage/storageaccounts":                                   "azurerm_storage_account",
	}

	if reverse {
		Reversed := reverseMap(Lookup)
		result = Reversed[resource]
	} else {
		result = Lookup[strings.TrimSuffix(strings.ToLower(resource), "/")]
	}

	var err error

	if result == "" {
		err = fmt.Errorf("resource %s not found", resource)

		return nil, err
	}

	return &result, err
}

func reverseMap(m map[string]string) map[string]string {
	n := make(map[string]string, len(m))
	for k, v := range m {
		n[v] = k
	}

	return n
}
