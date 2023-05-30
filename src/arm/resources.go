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
