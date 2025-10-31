package cli

import (
	"fmt"
	"log"

	bboxclient "bbox-cli/client"
)

func handleNat(client *bboxclient.BboxClient, args []string) {
	if len(args) < 1 {
		fmt.Println("nat usage: bboxcli nat <show> [id]")
		return
	}

	action := args[0]

	switch action {
	case "show":
		if len(args) > 1 {
			panic("Detailed view not implemented yet")
		} else {
			// Show list view
			showNatList(client)
		}
	default:
		fmt.Printf("Unknown nat action: %s\n", action)
		fmt.Println("nat usage: bboxcli nat <show> [id]")
	}
}

func showNatList(client *bboxclient.BboxClient) {
	nat := client.Nat()
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
