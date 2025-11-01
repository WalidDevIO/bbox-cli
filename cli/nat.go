package cli

import (
	"fmt"
	"log"

	bboxclient "bbox-cli/client"
)

func handleNat(client *bboxclient.BboxClient, args []string) {
	if len(args) < 1 {
		PrintUsage()
		return
	}

	nat := client.Nat()
	action := args[0]

	switch action {
	case "show":
		if len(args) > 1 {
			showNatDetail(nat, args[1])
		} else {
			// Show list view
			showNatList(nat)
		}
	case "enable":
		if len(args) < 2 {
			PrintUsage()
			return
		}
		nat.EnableNatRule(args[1])
	case "disable":
		if len(args) < 2 {
			PrintUsage()
			return
		}
		nat.DisableNatRule(args[1])
	default:
		fmt.Printf("Unknown nat action: %s\n", action)
		PrintUsage()
	}
}

func showNatList(nat *bboxclient.NatInterface) {
	rules, err := nat.GetNatRules()
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	if len(rules) == 0 {
		fmt.Println("No NAT rules found")
		return
	}

	fmt.Printf("%-4s %-3s %-15s %-15s %-15s %-15s %-15s\n",
		"", "ID", "DESCRIPTION", "DST IP", "DST PORTS", "SRC IP", "SRC PORTS")
	fmt.Println(repeatString("-", 85))

	for _, rule := range rules {
		status := "❌"
		if rule.Enable == bboxclient.Enabled {
			status = "✅"
		}

		// Truncate description to 15 chars
		desc := truncate(rule.Description, 15)

		// Truncate IPs and ports to 15 chars
		srcIP := truncate(defaultIfEmpty(rule.SrcIP.String(), "ANY"), 15)
		srcPorts := truncate(defaultIfEmpty(rule.SrcPorts.String(), "ANY"), 15)
		dstIP := truncate(defaultIfEmpty(rule.TargetIP.String(), "ANY"), 15)
		dstPorts := truncate(defaultIfEmpty(rule.TargetPorts.String(), "ANY"), 15)

		fmt.Printf("[%s] %-3d %-15s %-15s %-15s %-15s %-15s\n",
			status,
			rule.ID,
			desc,
			dstIP,
			dstPorts,
			srcIP,
			srcPorts,
		)
	}
}

func showNatDetail(nat *bboxclient.NatInterface, id string) {
	rules, err := nat.GetNatRules()
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	var ruleFound *bboxclient.NatRule
	for _, rule := range rules {
		if fmt.Sprintf("%d", rule.ID) == id {
			ruleFound = &rule
			break
		}
	}

	if ruleFound == nil {
		fmt.Printf("NAT rule with ID %s not found\n", id)
		return
	}

	fmt.Printf("NAT Rule ID: %d\n", ruleFound.ID)
	fmt.Printf("Description: %s\n", ruleFound.Description)
	fmt.Printf("Enabled: %t\n", ruleFound.Enable == bboxclient.Enabled)
	fmt.Printf("Source IP: %s\n", ruleFound.SrcIP.String())
	fmt.Printf("Source Ports: %s\n", ruleFound.SrcPorts.String())
	fmt.Printf("Target IP: %s\n", ruleFound.TargetIP.String())
	fmt.Printf("Target Ports: %s\n", ruleFound.TargetPorts.String())
}
