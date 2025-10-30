package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

type FirewallInterface struct {
	Client *BboxClient
}

const NotConcerned = "0"
const Any = ""

type EnableState int

const (
	Disabled EnableState = 0
	Enabled  EnableState = 1
)

type Action string

const (
	ActionAllow Action = "Accept"
	ActionDeny  Action = "Drop"
)

type Protocol string

const (
	ProtocolAny Protocol = "tcp,udp"
	ProtocolTCP Protocol = "tcp"
	ProtocolUDP Protocol = "udp"
)

type IPProtocol string

const (
	IPProtocolIPv4 IPProtocol = "IPv4"
	IPProtocolIPv6 IPProtocol = "IPv6"
	IPProtocolBoth IPProtocol = "IPv4+IPv6"
)

type FirewallRule struct {
	ID          int         `json:"id"`
	Description string      `json:"description"`
	Enable      EnableState `json:"enable"`
	Action      Action      `json:"action"`
	SrcIPNot    StringOrInt `json:"srcipnot"`
	SrcIP       StringOrInt `json:"srcip"`
	DstIPNot    StringOrInt `json:"dstipnot"`
	DstIP       StringOrInt `json:"dstip"`
	SrcPortNot  StringOrInt `json:"srcportnot"`
	SrcPorts    StringOrInt `json:"srcports"`
	DstPortNot  StringOrInt `json:"dstportnot"`
	DstPorts    StringOrInt `json:"dstports"`
	Order       int         `json:"order"`
	Protocols   Protocol    `json:"protocols"`
	IPProtocol  IPProtocol  `json:"ipprotocol"`
	Utilisation int         `json:"utilisation"`
}

type Firewall struct {
	Rules []FirewallRule `json:"rules"`
}

type FirewallResponse struct {
	Firewall Firewall `json:"firewall"`
}

func (fi *FirewallInterface) GetFirewallRules() ([]FirewallRule, error) {
	resp, err := fi.Client.Get("/firewall/rules")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var firewallResp []FirewallResponse
	if err := json.NewDecoder(resp.Body).Decode(&firewallResp); err != nil {
		return nil, err
	}

	if len(firewallResp) == 0 {
		return nil, errors.New("no firewall rules in response")
	}

	return firewallResp[0].Firewall.Rules, nil
}

func (fi *FirewallInterface) DeleteFirewallRule(ruleID string) error {
	if fi.Client.Bearer == nil {
		return errors.New("no bearer token available")
	}

	url := fmt.Sprintf("/firewall/rules/%s?btoken=%s", ruleID, fi.Client.Bearer.Token)
	r, err := http.NewRequest("DELETE", fi.Client.Url.String()+url, nil)
	if err != nil {
		return err
	}

	resp, err := fi.Client.Do(r)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("failed to delete rule: status %d", resp.StatusCode)
	}

	return nil
}

func (fi *FirewallInterface) AddFirewallRule(rule FirewallRule) error {
	if fi.Client.Bearer == nil {
		return errors.New("no bearer token available")
	}

	url := fmt.Sprintf("/firewall/rules?btoken=%s", fi.Client.Bearer.Token)
	data := fmt.Sprintf(
		"enable=%d&action=%s&srcipnot=%v&srcip=%v&dstipnot=%v&dstip=%v&srcportnot=%v&srcports=%v&dstportnot=%v&dstports=%v&order=%d&protocols=%v&ipprotocol=%v&description=%v",
		rule.Enable,
		rule.Action,
		rule.SrcIPNot,
		rule.SrcIP,
		rule.DstIPNot,
		rule.DstIP,
		rule.SrcPortNot,
		rule.SrcPorts,
		rule.DstPortNot,
		rule.DstPorts,
		rule.Order,
		rule.Protocols,
		rule.IPProtocol,
		rule.Description,
	)

	resp, err := fi.Client.Post(
		url, "application/x-www-form-urlencoded", io.Reader(strings.NewReader(data)),
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("failed to add rule: status %d", resp.StatusCode)
	}

	return nil
}

func (fi *FirewallInterface) UpdateFirewallRule(rule FirewallRule) error {
	if fi.Client.Bearer == nil {
		return errors.New("no bearer token available")
	}

	rules, err := fi.GetFirewallRules()
	if err != nil {
		return err
	}

	var ruleID string
	for _, r := range rules {
		if r.Description == rule.Description {
			ruleID = fmt.Sprintf("%d", r.ID)
			break
		}
	}

	url := fmt.Sprintf("/firewall/rules/%s?btoken=%s", ruleID, fi.Client.Bearer.Token)
	data := fmt.Sprintf(
		"enable=%d&action=%s&srcipnot=%v&srcip=%v&dstipnot=%v&dstip=%v&srcportnot=%v&srcports=%v&dstportnot=%v&dstports=%v&order=%d&protocols=%v&ipprotocol=%v&description=%v",
		rule.Enable,
		rule.Action,
		rule.SrcIPNot,
		rule.SrcIP,
		rule.DstIPNot,
		rule.DstIP,
		rule.SrcPortNot,
		rule.SrcPorts,
		rule.DstPortNot,
		rule.DstPorts,
		rule.Order,
		rule.Protocols,
		rule.IPProtocol,
		rule.Description,
	)

	r, err := http.NewRequest("PUT", fi.Client.Url.String()+url, io.Reader(strings.NewReader(data)))
	if err != nil {
		return err
	}

	resp, err := fi.Client.Do(r)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to update rule: status %d", resp.StatusCode)
	}

	return nil
}

func GenerateUniqueDescription(base string) string {
	// Generate a unique description by appending a UUID
	// Rule IDs are not predictable, so we use UUIDs for uniqueness
	return fmt.Sprintf("%s-bbcli-%s", base, uuid.New().String())
}
