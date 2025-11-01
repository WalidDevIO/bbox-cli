package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

// FirewallInterface provides methods to manage firewall rules
type FirewallInterface struct {
	Client *BboxClient
}

// GetFirewallRules retrieves all firewall rules from the device
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

// DeleteFirewallRule removes a firewall rule by its ID
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

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to delete rule: status %d", resp.StatusCode)
	}

	return nil
}

// AddFirewallRule creates a new firewall rule
func (fi *FirewallInterface) AddFirewallRule(rule FirewallRule) error {
	if fi.Client.Bearer == nil {
		return errors.New("no bearer token available")
	}

	url := fmt.Sprintf("/firewall/rules?btoken=%s", fi.Client.Bearer.Token)
	data := rule.RuleAsString()

	resp, err := fi.Client.Post(
		url,
		"application/x-www-form-urlencoded",
		strings.NewReader(data),
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

// UpdateFirewallRule modifies an existing firewall rule
func (fi *FirewallInterface) UpdateFirewallRule(rule FirewallRule) error {
	if fi.Client.Bearer == nil {
		return errors.New("no bearer token available")
	}

	// Find the rule ID by description
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

	// Prepare the update request
	url := fmt.Sprintf("/firewall/rules/%s", ruleID)
	data := rule.RuleAsString()

	r, err := fi.Client.NewRequest("PUT", url, strings.NewReader(data))
	if err != nil {
		return err
	}

	// Add bearer token to query parameters
	q := r.URL.Query()
	q.Add("btoken", fi.Client.Bearer.Token)
	r.URL.RawQuery = q.Encode()

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

// GenerateUniqueDescription creates a unique description for firewall rules
// by appending a UUID to the base description. This ensures uniqueness since
// rule IDs are not predictable.
func GenerateUniqueDescription(base string) string {
	return fmt.Sprintf("%s-bbcli-%s", base, uuid.New().String())
}

// RuleAsString converts the firewall rule to URL-encoded form data
// for API requests
func (r *FirewallRule) RuleAsString() string {
	return fmt.Sprintf(
		"enable=%d&action=%s&srcipnot=%v&srcip=%v&dstipnot=%v&dstip=%v&srcportnot=%v&srcports=%v&dstportnot=%v&dstports=%v&order=%d&protocols=%v&ipprotocol=%v&description=%v",
		r.Enable,
		r.Action,
		r.SrcIPNot,
		r.SrcIP,
		r.DstIPNot,
		r.DstIP,
		r.SrcPortNot,
		r.SrcPorts,
		r.DstPortNot,
		r.DstPorts,
		r.Order,
		r.Protocols,
		r.IPProtocol,
		r.Description,
	)
}
