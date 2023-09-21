package arm

import (
	_ "embed" // required for embed
)

//go:embed resources/azurerm_virtual_machine.template
var azurermVirtualMachine []byte

//go:embed resources/azurerm_virtual_machine_extension.template
var azurermVirtualMachineExtension []byte

//go:embed resources/azurerm_network_interface.template
var azurermNetworkInterface []byte

//go:embed resources/azurerm_network_security_group.template
var azurermNetworkSecurityGroup []byte

//go:embed resources/azurerm_public_ip.template
var azurermPublicIP []byte

//go:embed resources/azurerm_virtual_network.template
var azurermVirtualNetwork []byte

//go:embed resources/azurerm_storage_account.template
var azurermStorageAccount []byte

//go:embed resources/azurerm_subnet.template
var azurermSubnet []byte

//go:embed resources/azurerm_analysis_services_server.template
var azurermAnalysisServicesServer []byte

//go:embed resources/azurerm_api_management.template
var azurermAPIManagement []byte

//go:embed resources/azurerm_container_app.template
var azurermContainerApp []byte

//go:embed resources/azurerm_container_app_environment.template
var azurermContainerAppEnvironment []byte

//go:embed resources/azurerm_template_deployment.template
var azurermTemplateDeployment []byte

//go:embed resources/azurerm_role_assignment.template
var azurermRoleAssignment []byte

//go:embed resources/azurerm_user_assigned_identity.template
var azurermUserAssignedIdentity []byte

//go:embed resources/azurerm_log_analytics_workspace.template
var azurermLogAnalyticsWorkspace []byte

//go:embed resources/azurerm_role_definition.template
var azurermRoleDefinition []byte

//go:embed resources/azurerm_servicebus_namespace.template
var azurermServicebusNamespace []byte

//go:embed resources/azurerm_servicebus_namespace_authorization_rule.template
var azurermServicebusNamespaceAuthorizationRule []byte

//go:embed resources/azurerm_servicebus_queue.template
var azurermServicebusQueue []byte

//go:embed resources/azurerm_active_directory_domain_service.template
var azurermActiveDirectoryDomainService []byte

//go:embed resources/azurerm_bastion_host.template
var azurermBastionHost []byte

//go:embed resources/azurerm_key_vault.template
var azurermKeyVault []byte

//go:embed resources/azurerm_container_registry.template
var azurermContainerRegistry []byte

//go:embed resources/azurerm_kubernetes_cluster.template
var azurermKubernetesCluster []byte

//go:embed resources/azurerm_log_analytics_solution.template
var azurermLogAnalyticsSolution []byte

//go:embed resources/azurerm_private_dns_zone.template
var azurermPrivateDNSZone []byte

//go:embed resources/azurerm_private_dns_zone_virtual_network_link.template
var azurermPrivateDNSZoneVirtualNetworkLink []byte

//go:embed resources/azurerm_private_endpoint.template
var azurermPrivateEndpoint []byte

//go:embed resources/azurerm_monitor_activity_log_alert.template
var azurermMonitorActivityLogAlert []byte

//go:embed resources/azurerm_web_application_firewall_policy.template
var azurermWebApplicationFirewallPolicy []byte

//go:embed resources/azurerm_application_gateway.template
var azurermApplicationGateway []byte

//go:embed resources/azurerm_network_interface_application_gateway_backend_address_pool_association.template
var azurermNetworkInterfaceApplicationGatewayBackendAddressPoolAssociation []byte

//go:embed resources/azurerm_availability_set.template
var azurermAvailabilitySet []byte

//go:embed resources/azurerm_managed_disk.template
var azurermComputeDisks []byte

//go:embed resources/azurerm_network_security_rule.template
var azurermNetworkNSGSecurityRules []byte
