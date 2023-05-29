package arm

import (
	"strings"

	"github.com/rs/zerolog/log"
)

func lookup(myType string) []byte {
	TFLookup := map[string]interface{}{
		"Microsoft.AnalysisServices/servers":                 azurermAnalysisServicesServer,
		"Microsoft.ApiManagement/service":                    azurermAPIManagement,
		"Microsoft.App/containerApps":                        azurermContainerApp,
		"Microsoft.App/managedEnvironments":                  azurermContainerAppEnvironment,
		"Microsoft.Authorization/roleAssignments":            azurermRoleAssignment,
		"Microsoft.Authorization/roleDefinitions":            azurermRoleDefinition,
		"Microsoft.Compute/virtualMachines":                  azurermVirtualMachine,
		"Microsoft.Compute/virtualMachines/extensions":       azurermVirtualMachineExtension,
		"Microsoft.ManagedIdentity/userAssignedIdentities":   azurermUserAssignedIdentity,
		"Microsoft.Network/networkInterfaces":                azurermNetworkInterface,
		"Microsoft.Network/networkSecurityGroups":            azurermNetworkSecurityGroup,
		"Microsoft.Network/publicIPAddresses":                azurermPublicIP,
		"Microsoft.Network/virtualNetworks":                  azurermVirtualNetwork,
		"Microsoft.Network/virtualNetworks/subnets":          azurermSubnet,
		"Microsoft.OperationalInsights/workspaces":           azurermLogAnalyticsWorkspace,
		"Microsoft.Storage/storageAccounts":                  azurermStorageAccount,
		"Microsoft.Resources/deployments":                    azurermTemplateDeployment,
		"Microsoft.ServiceBus/namespaces":                    azurermServicebusNamespace,
		"Microsoft.ServiceBus/namespaces/authorizationRules": azurermServicebusNamespaceAuthorizationRule,
		"Microsoft.ServiceBus/namespaces/queues":             azurermServicebusQueue,
	}

	var myContent []byte
	var ok bool

	if TFLookup[myType] != nil {
		myContent, ok = TFLookup[strings.TrimSuffix(myType, "/")].([]byte)
		if !ok {
			log.Warn().Msg("Failed to cast lookup")
		}
	} else {
		// we don't want to half the parsing so just log it.
		log.Warn().Msgf("%s not found", myType)
	}

	return myContent
}
