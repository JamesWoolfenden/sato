package sato

import (
	"bytes"
	"fmt"
	"strings"
	tftemplate "text/template"

	"github.com/awslabs/goformation/v7/cloudformation"
	"github.com/rs/zerolog/log"
)

// ParseResources converts resource to Terraform
func ParseResources(resources cloudformation.Resources, funcMap tftemplate.FuncMap, destination string) error {
	for item, resource := range resources {
		var output bytes.Buffer

		myType := resources[item].AWSCloudFormationType()

		myContent := lookup(myType)

		//needs to pivot on policy template from resource
		tmpl, err := tftemplate.New("sato").Funcs(funcMap).Parse(string(myContent))

		if err != nil {
			return err
		}

		_ = tmpl.Execute(&output, M{
			"resource": resource,
			"item":     item,
		})

		err = Write(ReplaceDependant(ReplaceVariables(output.String())), destination, fmt.Sprint(ToTFName(myType), ".", strings.ToLower(item)))
		if err != nil {
			return err
		}
	}
	return nil
}

func lookup(myType string) []byte {
	TFLookup := map[string]interface{}{
		"AWS::ApplicationAutoScaling::ScalableTarget":      awsAppAutoscalingTarget,
		"AWS::ApplicationAutoScaling::ScalingPolicy":       awsAppAutoscalingPolicy,
		"AWS::AutoScaling::AutoScalingGroup":               awsAutoscalingGroup,
		"AWS::AutoScaling::LaunchConfiguration":            awsLaunchConfiguration,
		"AWS::AutoScaling::LifecycleHook":                  awsAutoscalingLifecycleHook,
		"AWS::AutoScaling::ScalingPolicy":                  awsAutoscalingPolicy,
		"AWS::AutoScaling::ScheduledAction":                awsAutoscalingSchedule,
		"AWS::Backup::BackupPlan":                          awsBackupPlan,
		"AWS::Backup::BackupSelection":                     awsBackupSelection,
		"AWS::Backup::BackupVault":                         awsBackupVault,
		"AWS::Cloud9::EnvironmentEC2":                      awsCloud9EnvironmentEc2,
		"AWS::CloudFormation::Stack":                       awsCloudformationStack,
		"AWS::CloudFront::CloudFrontOriginAccessIdentity":  awsCloudfrontOriginAccessIdentity,
		"AWS::CloudFront::Distribution":                    awsCloudfrontDistribution,
		"AWS::CloudWatch::Alarm":                           awsCloudwatchMetricAlarm,
		"AWS::CloudWatch::Dashboard":                       awsCloudwatchDashboard,
		"AWS::CodeBuild::Project":                          awsCodebuildProject,
		"AWS::CodeCommit::Repository":                      awsCodecommitRepository,
		"AWS::CodePipeline::Pipeline":                      awsCodepipeline,
		"AWS::Config::ConfigRule":                          awsConfigConfigRule,
		"AWS::Config::ConfigurationRecorder":               awsConfigConfigurationRecorder,
		"AWS::Config::DeliveryChannel":                     awsConfigDeliveryChannel,
		"AWS::DMS::Endpoint":                               awsDmsEndpoint,
		"AWS::DMS::ReplicationInstance":                    awsDmsReplicationinstance,
		"AWS::DMS::ReplicationSubnetGroup":                 awsDmsReplicationSubnetGroup,
		"AWS::DMS::ReplicationTask":                        awsDmsReplicationTask,
		"AWS::DirectoryService::MicrosoftAD":               awsDirectoryServiceDirectory,
		"AWS::DynamoDB::Table":                             awsDynamodbTable,
		"AWS::EC2::DHCPOptions":                            awsVpcDhcpOptions,
		"AWS::EC2::EIP":                                    awsEIP,
		"AWS::EC2::EIPAssociation":                         awsEipAssociation,
		"AWS::EC2::FlowLog":                                awsFlowLog,
		"AWS::EC2::Instance":                               awsInstance,
		"AWS::EC2::InternetGateway":                        awsInternetGateway,
		"AWS::EC2::LaunchTemplate":                         awsLaunchTemplate,
		"AWS::EC2::NatGateway":                             awsNatGateway,
		"AWS::EC2::NetworkAcl":                             awsNetworkACL,
		"AWS::EC2::NetworkAclEntry":                        awsNetworkACLRule,
		"AWS::EC2::NetworkInterface":                       awsNetworkInterface,
		"AWS::EC2::Route":                                  awsRoute,
		"AWS::EC2::RouteTable":                             awsRouteTable,
		"AWS::EC2::SecurityGroup":                          awsSecurityGroup,
		"AWS::EC2::SecurityGroupEgress":                    awsSecurityGroupRuleEgress,
		"AWS::EC2::SecurityGroupIngress":                   awsSecurityGroupRuleIngress,
		"AWS::EC2::Subnet":                                 awsSubnet,
		"AWS::EC2::SubnetNetworkAclAssociation":            awsNetworkACLAssociation,
		"AWS::EC2::SubnetRouteTableAssociation":            awsRouteTableAssociation,
		"AWS::EC2::VPC":                                    awsVpc,
		"AWS::EC2::VPCDHCPOptionsAssociation":              awsVpcDhcpOptionsAssociation,
		"AWS::EC2::VPCEndpoint":                            awsVpcEndpoint,
		"AWS::EC2::VPCGatewayAttachment":                   awsVpnGatewayAttachment,
		"AWS::EC2::Volume":                                 awsEbsVolume,
		"AWS::ECS::Cluster":                                awsEcsCluster,
		"AWS::ECS::Service":                                awsEcsService,
		"AWS::ECS::TaskDefinition":                         awsEcsTaskDefinition,
		"AWS::EFS::FileSystem":                             awsEfsFileSystem,
		"AWS::EFS::MountTarget":                            awsEfsMountTarget,
		"AWS::EKS::Cluster":                                awsEksCluster,
		"AWS::EKS::Nodegroup":                              awsEksNodeGroup,
		"AWS::ElastiCache::ParameterGroup":                 awsElasticacheParameterGroup,
		"AWS::ElastiCache::ReplicationGroup":               awsElasticacheReplicationGroup,
		"AWS::ElastiCache::SubnetGroup":                    awsElasticacheSubnetGroup,
		"AWS::ElasticLoadBalancing::LoadBalancer":          awsElb,
		"AWS::ElasticLoadBalancingV2::Listener":            awsLbListener,
		"AWS::ElasticLoadBalancingV2::ListenerRule":        awsLbListenerRule,
		"AWS::ElasticLoadBalancingV2::LoadBalancer":        awsLb,
		"AWS::ElasticLoadBalancingV2::TargetGroup":         awsLbTargetGroup,
		"AWS::Events::Rule":                                awsCloudwatchEventRule,
		"AWS::IAM::AccessKey":                              awsIamAccessKey,
		"AWS::IAM::Group":                                  awsIamGroup,
		"AWS::IAM::InstanceProfile":                        awsIamInstanceProfile,
		"AWS::IAM::ManagedPolicy":                          awsIamManagedPolicy,
		"AWS::IAM::Policy":                                 awsIamPolicy,
		"AWS::IAM::Role":                                   awsIamRole,
		"AWS::IAM::User":                                   awsIamUser,
		"AWS::IAM::UserToGroupAddition":                    awsIamGroupMembership,
		"AWS::KMS::Alias":                                  awskmsAlias,
		"AWS::KMS::Key":                                    awsKmsKey,
		"AWS::Lambda::EventSourceMapping":                  awsLambdaEventSourceMapping,
		"AWS::Lambda::Function":                            awsLambdaFunction,
		"AWS::Lambda::Permission":                          awsLambdaPermission,
		"AWS::Lambda::Version":                             awsLambdaVersion,
		"AWS::Logs::LogGroup":                              awsCloudwatchLogGroup,
		"AWS::Logs::MetricFilter":                          awsCloudwatchLogMetricFilter,
		"AWS::Neptune::DBCluster":                          awsNeptuneCluster,
		"AWS::Neptune::DBClusterParameterGroup":            awsNeptuneClusterDBParameterGroup,
		"AWS::Neptune::DBInstance":                         awsNeptuneDbInstance,
		"AWS::Neptune::DBParameterGroup":                   awsNeptuneDBParameterGroup,
		"AWS::Neptune::DBSubnetGroup":                      awsNeptuneDnSubnetGroup,
		"AWS::RDS::DBCluster":                              awsRdsCluster,
		"AWS::RDS::DBClusterParameterGroup":                awsDbParameterGroup,
		"AWS::RDS::DBInstance":                             awsDbInstance,
		"AWS::RDS::DBParameterGroup":                       awsDbParameterGroup,
		"AWS::RDS::DBSubnetGroup":                          awsDbSubnetGroup,
		"AWS::Route53::RecordSet":                          awsRoute53Record,
		"AWS::S3::Bucket":                                  awsS3Bucket,
		"AWS::S3::BucketPolicy":                            awsS3BucketPolicy,
		"AWS::SNS::Subscription":                           awsSNSSubscription,
		"AWS::SNS::Topic":                                  awsSNSTopic,
		"AWS::SNS::TopicPolicy":                            awsSNSTopicPolicy,
		"AWS::SQS::Queue":                                  awsSqsQueue,
		"AWS::SSM::Association":                            awsSsmAssociation,
		"AWS::SSM::Document":                               awsSsmDocument,
		"AWS::SSM::MaintenanceWindow":                      awsSsmMaintenanceWindow,
		"AWS::SSM::MaintenanceWindowTarget":                awsSsmMaintenanceWindowTarget,
		"AWS::SSM::MaintenanceWindowTask":                  awsSsmMaintenanceWindowTask,
		"AWS::SecretsManager::Secret":                      awsSecretsManagerSecret,
		"AWS::ServiceCatalog::Portfolio":                   awsServiceCatalogPortfolio,
		"AWS::ServiceCatalog::PortfolioProductAssociation": awsServiceCatalogProductPortfolioAssociation,
		"AWS::ServiceCatalog::PortfolioShare":              awsServiceCatalogPortfolioShare,
		"AWS::ServiceCatalog::TagOption":                   awsServiceCatalogTagOption,
		"AWS::ServiceCatalog::TagOptionAssociation":        awsServiceCatalogTagOptionAssociation,
		"AWS::ServiceDiscovery::Service":                   awsServiceDiscoveryService,
		"AWS::StepFunctions::StateMachine":                 awsStepfunctionStateMachine,
		"AWS::WAFv2::WebACLAssociation":                    awsWAFv2WebACLAssociation,
	}

	var myContent []byte
	if TFLookup[myType] != nil {
		myContent = TFLookup[myType].([]byte)
	} else {
		//we don't want to half the parsing so just log it.
		log.Warn().Msgf("%s not found", myType)
	}
	return myContent
}
