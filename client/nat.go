package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type NatInterface struct {
	Client *BboxClient
}

type NatResponse struct {
	Nat NatRules `json:"nat"`
}

type NatRules struct {
	Enable EnableState `json:"enable"`
	Rules  []NatRule   `json:"rules"`
}

type NatRule struct {
	ID          int         `json:"id"`
	Enable      EnableState `json:"enable"`
	Description string      `json:"description"`
	Protocol    Protocol    `json:"protocol"`
	SrcIP       StringOrInt `json:"externalip"`
	SrcPorts    StringOrInt `json:"externalport"`
	TargetIP    StringOrInt `json:"internalip"`
	TargetPorts StringOrInt `json:"internalport"`
}

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

func (ni *NatInterface) changeNatRuleState(ruleID string, enable EnableState) error {
	path := "/nat/rules/" + ruleID
	data := fmt.Sprintf("enable=%d", enable)
	req, err := http.NewRequest("PUT", path, io.Reader(strings.NewReader(data)))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	_, err = ni.Client.Do(req)
	return err
}

func (ni *NatInterface) EnableNatRule(ruleID string) error {
	return ni.changeNatRuleState(ruleID, Enabled)
}

func (ni *NatInterface) DisableNatRule(ruleID string) error {
	return ni.changeNatRuleState(ruleID, Disabled)
}
