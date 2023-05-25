package arm

import "github.com/rs/zerolog/log"

func lookup(myType string) []byte {
	TFLookup := map[string]interface{}{
		"Microsoft.AnalysisServices/servers":               azurermAnalysisServicesServer,
		"Microsoft.ApiManagement/service":                  azurermAPIManagement,
		"Microsoft.App/containerApps":                      azurermContainerApp,
		"Microsoft.App/managedEnvironments":                azurermContainerAppEnvironment,
		"Microsoft.Authorization/roleAssignments":          azurermRoleAssignment,
		"Microsoft.Authorization/roleDefinitions":          azurermRoleDefinition,
		"Microsoft.Compute/virtualMachines":                azurermVirtualMachine,
		"Microsoft.Compute/virtualMachines/extensions":     azurermVirtualMachineExtension,
		"Microsoft.ManagedIdentity/userAssignedIdentities": azurermUserAssignedIdentity,
		"Microsoft.Network/networkInterfaces":              azurermNetworkInterface,
		"Microsoft.Network/networkSecurityGroups":          azurermNetworkSecurityGroup,
		"Microsoft.Network/publicIPAddresses":              azurermPublicIP,
		"Microsoft.Network/virtualNetworks":                azurermVirtualNetwork,
		"Microsoft.Network/virtualNetworks/subnets":        azurermSubnet,
		"Microsoft.OperationalInsights/workspaces":         azurermLogAnalyticsWorkspace,
		"Microsoft.Storage/storageAccounts":                azurermStorageAccount,
		"Microsoft.Resources/deployments":                  azurermTemplateDeployment,
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
