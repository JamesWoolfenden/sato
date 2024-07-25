package arm

import (
	"strings"

	"github.com/rs/zerolog/log"
)

func lookup(myType string) []byte {
	TFLookup := map[string]interface{}{
		"microsoft.aad/domainservices":                                       azurermActiveDirectoryDomainService,
		"microsoft.analysisservices/servers":                                 azurermAnalysisServicesServer,
		"microsoft.apiManagement/service":                                    azurermAPIManagement,
		"microsoft.app/containerapps":                                        azurermContainerApp,
		"microsoft.app/managedenvironments":                                  azurermContainerAppEnvironment,
		"microsoft.authorization/roleassignments":                            azurermRoleAssignment,
		"microsoft.authorization/roledefinitions":                            azurermRoleDefinition,
		"microsoft.containerregistry/registries":                             azurermContainerRegistry,
		"microsoft.containerservice/managedclusters":                         azurermKubernetesCluster,
		"microsoft.compute/virtualmachines":                                  azurermVirtualMachine,
		"microsoft.compute/virtualmachines/extensions":                       azurermVirtualMachineExtension,
		"microsoft.keyvault/vaults":                                          azurermKeyVault,
		"microsoft.managedidentity/userassignedidentities":                   azurermUserAssignedIdentity,
		"microsoft.network/applicationgateways":                              azurermApplicationGateway,
		"Microsoft.Network/applicationGateways/backendAddressPools":          azurermNetworkInterfaceApplicationGatewayBackendAddressPoolAssociation,
		"microsoft.network/bastionhosts":                                     azurermBastionHost,
		"microsoft.network/networkinterfaces":                                azurermNetworkInterface,
		"microsoft.network/networksecuritygroups":                            azurermNetworkSecurityGroup,
		"microsoft.network/publicipaddresses":                                azurermPublicIP,
		"microsoft.network/virtualnetworks":                                  azurermVirtualNetwork,
		"microsoft.network/virtualnetworks/subnets":                          azurermSubnet,
		"microsoft.network/privatednszones":                                  azurermPrivateDNSZone,
		"microsoft.network/privatednszones/virtualnetworklinks":              azurermPrivateDNSZoneVirtualNetworkLink,
		"microsoft.network/privateendpoints":                                 azurermPrivateEndpoint,
		"microsoft.network/applicationgatewaywebapplicationfirewallpolicies": azurermWebApplicationFirewallPolicy,
		"microsoft.operationalinsights/workspaces":                           azurermLogAnalyticsWorkspace,
		"microsoft.insights/activitylogalerts":                               azurermMonitorActivityLogAlert,
		"microsoft.operationsmanagement/solutions":                           azurermLogAnalyticsSolution,
		"microsoft.resources/deployments":                                    azurermTemplateDeployment,
		"microsoft.servicebus/namespaces":                                    azurermServicebusNamespace,
		"microsoft.servicebus/namespaces/authorizationRules":                 azurermServicebusNamespaceAuthorizationRule,
		"microsoft.servicebus/namespaces/queues":                             azurermServicebusQueue,
		"microsoft.storage/storageaccounts":                                  azurermStorageAccount,
		"microsoft.compute/availabilitysets":                                 azurermAvailabilitySet,
		"microsoft.compute/disks":                                            azurermComputeDisks,
		"microsoft.network/networksecuritygroups/securityrules":              azurermNetworkNSGSecurityRules,
	}

	var myContent []byte

	//goland:noinspection GoLinter
	var ok bool

	myType = strings.ToLower(strings.TrimSuffix(myType, "/"))

	if TFLookup[myType] != nil {
		myContent, ok = TFLookup[myType].([]byte)
		if !ok {
			log.Warn().Msg("Failed to cast lookup")
		}
	} else {
		// we don't want to half the parsing so just log it.
		log.Warn().Msgf("%s not found", myType)
	}

	return myContent
}
