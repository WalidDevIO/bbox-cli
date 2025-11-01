package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// NatInterface provides methods to interact with NAT rules on the Bbox device.
type NatInterface struct {
	Client *BboxClient
}

// GetNatRules retrieves all NAT rules from the Bbox device.
func (ni *NatInterface) GetNatRules() ([]NatRule, error) {
	var result []NatResponse
	r, err := ni.Client.Get("/nat/rules")
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	if err := json.NewDecoder(r.Body).Decode(&result); err != nil {
		return nil, err
	}
	return result[0].Nat.Rules, nil
}

// GetNatRuleByID retrieves a specific NAT rule by its ID.
func (ni *NatInterface) GetNatRuleByID(ruleID int) (NatRule, error) {
	var result NatRule
	r, err := ni.Client.Get("/nat/rules/" + string(rune(ruleID)))
	if err != nil {
		return result, err
	}
	defer r.Body.Close()

	if err := json.NewDecoder(r.Body).Decode(&result); err != nil {
		return result, err
	}
	return result, nil
}

// changeNatRuleState enables or disables a NAT rule based on the provided state.
func (ni *NatInterface) changeNatRuleState(ruleID string, enable EnableState) error {
	path := "/nat/rules/" + ruleID
	data := fmt.Sprintf("enable=%d", enable)
	req, err := http.NewRequest("PUT", path, strings.NewReader(data))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	_, err = ni.Client.Do(req)
	return err
}

// EnableNatRule enables a NAT rule by its ID.
func (ni *NatInterface) EnableNatRule(ruleID string) error {
	return ni.changeNatRuleState(ruleID, Enabled)
}

// DisableNatRule disables a NAT rule by its ID.
func (ni *NatInterface) DisableNatRule(ruleID string) error {
	return ni.changeNatRuleState(ruleID, Disabled)
}
