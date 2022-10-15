package sato

import (
	_ "embed" //required for embed"
)

//go:embed resources/aws_sns_topic.policy.template
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
var awsNetworkAclRule []byte

//go:embed resources/aws_network_acl.template
var awsNetworkAcl []byte

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

//go:embed resources/aws_vpc_dhcp_options_association.template
var awsNetworkAclAssociation []byte

//go:embed resources/aws_vpc_endpoint.template
var awsVpcEndpoint []byte

//go:embed resources/aws_flow_log.template
var awsFlowLog []byte

//go:embed resources/aws_internet_gateway.template
var awsInternetGateway []byte
