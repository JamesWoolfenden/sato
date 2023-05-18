package arm

import "github.com/rs/zerolog/log"

func lookup(myType string) []byte {
	TFLookup := map[string]interface{}{
		"Microsoft.Compute/virtualMachines":            azurermVirtualMachine,
		"Microsoft.Compute/virtualMachines/extensions": azurermVirtualMachineExtension,
		"Microsoft.Network/networkInterfaces":          azurermNetworkInterface,
		"Microsoft.Network/networkSecurityGroups":      azurermNetworkSecurityGroup,
		"Microsoft.Network/publicIPAddresses":          azurermPublicIP,
		"Microsoft.Network/virtualNetworks":            azurermVirtualNetwork,
		"Microsoft.Storage/storageAccounts":            azurermStorageAccount,
		"Microsoft.Network/virtualNetworks/subnets":    azurermSubnet,
		"Microsoft.AnalysisServices/servers":           azurermAnalysisServicesServer,
		"Microsoft.ApiManagement/service":              azurermApiManagement,
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
