package cli

import (
	"fmt"
	"log"
	"strconv"

	bboxclient "bbox-cli/client"
)

func handleFirewall(client *bboxclient.BboxClient, args []string) {
	if len(args) < 1 {
		fmt.Println("firewall usage: bboxcli firewall <show> [id]")
		return
	}

	action := args[0]

	switch action {
	case "show":
		if len(args) > 1 {
			// Show detailed view for specific ID
			showFirewallDetail(client, args[1])
		} else {
			// Show list view
			showFirewallList(client)
		}
	case "add":
		panic("Add firewall rule not implemented yet")
	case "delete":
		if len(args) < 2 {
			fmt.Println("firewall usage: bboxcli firewall delete <id>")
			return
		}
		deleteFirewallRule(client, args[1])
	default:
		fmt.Printf("Unknown firewall action: %s\n", action)
		fmt.Println("firewall usage: bboxcli firewall <show> [id]")
	}
}

func showFirewallList(client *bboxclient.BboxClient) {
	fw := client.Firewall()
	rules, err := fw.GetFirewallRules()
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	if len(rules) == 0 {
		fmt.Println("No firewall rules found")
		return
	}

	fmt.Printf("%-5s %-3s %-15s %-10s %-15s %-15s %-15s %-15s\n",
		"", "ID", "DESCRIPTION", "ACTION", "DST IP", "DST PORTS", "SRC IP", "SRC PORTS")
	fmt.Println(repeatString("-", 100))

	for _, rule := range rules {
		status := "❌"
		if rule.Enable == bboxclient.Enabled {
			status = "✅"
		}

		// Truncate description to 15 chars
		desc := truncate(rule.Description, 15)

		// Truncate IPs and ports to 15 chars
		dstIP := truncate(defaultIfEmpty(rule.DstIP.String(), "ANY"), 15)
		dstPorts := truncate(defaultIfEmpty(rule.DstPorts.String(), "ANY"), 15)
		srcIP := truncate(defaultIfEmpty(rule.SrcIP.String(), "ANY"), 15)
		srcPorts := truncate(defaultIfEmpty(rule.SrcPorts.String(), "ANY"), 15)

		fmt.Printf("[%s] %-3d %-15s %-10s %-15s %-15s %-15s %-15s\n",
			status,
			rule.ID,
			desc,
			rule.Action,
			dstIP,
			dstPorts,
			srcIP,
			srcPorts,
		)
	}
}

func showFirewallDetail(client *bboxclient.BboxClient, idStr string) {
	ruleID, err := strconv.Atoi(idStr)
	if err != nil {
		fmt.Printf("Error: invalid ID '%s'\n", idStr)
		return
	}

	fw := client.Firewall()
	rules, err := fw.GetFirewallRules()
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	var rule *bboxclient.FirewallRule
	for i := range rules {
		if rules[i].ID == ruleID {
			rule = &rules[i]
			break
		}
	}

	if rule == nil {
		fmt.Printf("Error: Rule with ID %d not found\n", ruleID)
		return
	}

	// Display each field on a separate line
	status := "Disabled"
	if rule.Enable == 1 {
		status = "Enabled"
	}

	fmt.Println("\nFirewall Rule Details")
	fmt.Println(repeatString("=", 50))
	fmt.Printf("ID:          %d\n", rule.ID)
	fmt.Printf("Status:      %s\n", status)
	fmt.Printf("Description: %s\n", rule.Description)
	fmt.Printf("Action:      %s\n", rule.Action)
	fmt.Printf("Order:       %d\n", rule.Order)
	fmt.Println(repeatString("-", 50))
	fmt.Printf("Source IP:      %s\n", defaultIfEmpty(rule.SrcIP.String(), "ANY"))
	fmt.Printf("Source Ports:   %s\n", defaultIfEmpty(rule.SrcPorts.String(), "ANY"))
	fmt.Printf("Dest IP:        %s\n", defaultIfEmpty(rule.DstIP.String(), "ANY"))
	fmt.Printf("Dest Ports:     %s\n", defaultIfEmpty(rule.DstPorts.String(), "ANY"))
	fmt.Printf("Protocols:      %s\n", defaultIfEmpty(string(rule.Protocols), "ANY"))
	fmt.Println(repeatString("=", 50))
}

func deleteFirewallRule(client *bboxclient.BboxClient, ruleID string) {
	fw := client.Firewall()
	err := fw.DeleteFirewallRule(ruleID)
	if err != nil {
		log.Fatalf("Error deleting firewall rule: %v", err)
	}

	fmt.Printf("Firewall rule with ID %s deleted successfully\n", ruleID)
}
