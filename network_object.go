package main

import (
	"github.com/zclconf/go-cty/cty"
	"regexp"
	"strings"
)

var networkObjects []ObjectNetwork

type ObjectNetwork struct {
	name string

	//Usually, these would be net.IP values, but since we're passing them as a string anyway, I don't think it matters.
	host string
	mask string
}

func getNetworkObject(path string) {
	var MaskStr string
	
	parsedObjects := parseFile(path, "object network")

	//Using regex to pick up any things like obj_ or obj-. Without it, it'll drop them.
	networkObjectRegex := regexp.MustCompile(`^object network`)

	for _, slices := range parsedObjects {
		objectNetworkName := strings.TrimSpace(networkObjectRegex.ReplaceAllString(slices[0], ""))
		objectHostName := strings.Split(strings.TrimSpace(slices[1]), " ")
		hostName := objectHostName[1]

		//TODO Figure out what to do for nat blocks inside object network blocks
		if strings.Contains(objectHostName[0], "subnet") {
			MaskStr = objectHostName[2]
		} else {
			MaskStr = "255.255.255.255"
		}

		network_object := ObjectNetwork{
			name: objectNetworkName,
			host: hostName,
			mask: MaskStr,
		}

		networkObjects = append(networkObjects, network_object)
	}
}

func GenerateNetworObjectskHCL(organization_id string, path string) {
	getNetworkObject(path)

	for _, obj := range networkObjects {
		netObj := rootBody.AppendNewBlock("resource", []string{"meraki_organizations_policy_objects", obj.name})
		netObjBody := netObj.Body()
		netObjBody.SetAttributeValue("organization_id", cty.StringVal(organization_id))
		netObjBody.SetAttributeValue("name", cty.StringVal(obj.name))
		netObjBody.SetAttributeValue("ip", cty.StringVal(obj.host))
		netObjBody.SetAttributeValue("mask", cty.StringVal(obj.mask))
		netObjBody.SetAttributeValue("type", cty.StringVal("ipAndMask"))
		rootBody.AppendNewline()
	}
}
