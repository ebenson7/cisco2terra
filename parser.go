package main

/*

***IMPORTANT TO NOTE***
I haven't worked at all with Meraki so all my knowledge comes from looking through the Cisco Meraki Terraform provider
and Meraki APIv1. The following code should be seen as a start to an interesting problem of implementing a homegrown ASA config parser
and creating Terraform HCL configuration, not as a full solution yet.

All of these also really depend on if you're using Cisco Defense Orchestrator as well. You'd use that to push out configuration
by using Meraki Templates instead of using Terraform, I would assume.

A lot of things don't map from ASA to Meraki MX. I've made a list below of things I've found online that do or don't. This isn't an entirely exhaustive
list either. Things could map from ASA to Meraki as a whole, but my scope for this is just firewall-to-firewall appliances.

Things that map from ASA to Meraki MX (IN PROGRESS):
- Access Lists -> NAT/Port Forwarding
- Interfaces -> LAN Ports //Probably not implementing at this time since this really depends on context of its created or not
- VLANs -> VLANs //Probably not implementing at this time since this really depends on context of its created or not
- object network -> Network Object
- object group network -> Network Groups

Things that don't seem to map from ASA to Meraki MX, I don't think would map well, or I can't find any documentation on currently
- Names for IP addresses
- ip local pools for VPN clients. //This is handled at the Network-wide configuration, which is currently outside of the scope of the project.
- Pager
- Logging
- MTU
- ICMP timeouts
- Timeouts
- User Identity
- AAA
- HTTP
- SNMP
- Crypto
- SSH
- DHCP

Things that might be able to map, but I can't find anything on it.
- object group services // Probably can generate ACL with any/any to any/destination port but I'd rather be safe that sorry.
- policy-maps
- xlate -> NAT/Port Forwarding // Not sure if Meraki API allows allowed IPs of "None". I know there's an implicit deny, but not sure.
*/

import (
	"bufio"
	"log"
	"os"
	"strings"
)

func parseFile(path string, command string) [][]string {
	asa_file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer asa_file.Close()

	var block []string
	var body  [][]string

	scanner := bufio.NewScanner(asa_file)
	for scanner.Scan() {
		cfgLine := scanner.Text()

		//Since we're parsing only the beginning of the group of strings we want, we don't need to parse the whole line
		//Regex took 640.291µs
		//HasPrefix took 352.084µs

		//Checks if we're currently working on a block or not. If not, create a new block.
		switch {
		case block == nil && strings.HasPrefix(cfgLine, command):
			block = append(block, cfgLine)
		case block != nil && strings.HasPrefix(cfgLine, " "):
			block = append(block, cfgLine)
		case block != nil && strings.HasPrefix(cfgLine, command):
			body = append(body, block)
			block = nil
			block = append(block, cfgLine)
		case block != nil && !strings.HasPrefix(cfgLine, " "):
			body = append(body, block)
			block = nil
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return body
}