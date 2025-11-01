package client

import (
	"errors"
)

var (
	ErrFirewallRuleNotFound = errors.New("firewall rule not found")
	ErrNatRuleNotFound      = errors.New("NAT rule not found")
)

// Constants for special values
const (
	Any = ""
)

// EnableState represents whether a rule or feature is enabled
type EnableState int

const (
	Disabled EnableState = 0
	Enabled  EnableState = 1
)

// Action represents the action to take when a rule matches
type Action string

const (
	ActionAllow Action = "Accept"
	ActionDeny  Action = "Drop"
)

// Protocol represents network protocols
type Protocol string

const (
	ProtocolAny Protocol = "tcp,udp"
	ProtocolTCP Protocol = "tcp"
	ProtocolUDP Protocol = "udp"
)

// IPProtocol represents IP protocol versions
type IPProtocol string

const (
	IPProtocolIPv4 IPProtocol = "IPv4"
	IPProtocolIPv6 IPProtocol = "IPv6"
	IPProtocolBoth IPProtocol = "IPv4+IPv6"
)

// FirewallRule represents a single firewall rule configuration
type FirewallRule struct {
	ID          int         `json:"id"`
	Description string      `json:"description"`
	Enable      EnableState `json:"enable"`
	Action      Action      `json:"action"`

	// Source configuration
	SrcIPNot   EnableState `json:"srcipnot"`
	SrcIP      StringOrInt `json:"srcip"`
	SrcPortNot EnableState `json:"srcportnot"`
	SrcPorts   StringOrInt `json:"srcports"`

	// Destination configuration
	DstIPNot   EnableState `json:"dstipnot"`
	DstIP      StringOrInt `json:"dstip"`
	DstPortNot EnableState `json:"dstportnot"`
	DstPorts   StringOrInt `json:"dstports"`

	// Protocol and ordering
	Order       int        `json:"order"`
	Protocols   Protocol   `json:"protocols"`
	IPProtocol  IPProtocol `json:"ipprotocol"`
	Utilisation int        `json:"utilisation"`
}

// Firewall represents a collection of firewall rules
type Firewall struct {
	Rules []FirewallRule `json:"rules"`
}

// FirewallResponse wraps the firewall data from API responses
type FirewallResponse struct {
	Firewall Firewall `json:"firewall"`
}

// NatResponse wraps the NAT rules data from API responses
type NatResponse struct {
	Nat NatRules `json:"nat"`
}

// NatRules represents a collection of NAT rules
type NatRules struct {
	Enable EnableState `json:"enable"`
	Rules  []NatRule   `json:"rules"`
}

// NatRule represents a single NAT rule configuration
type NatRule struct {
	ID          int         `json:"id"`
	Enable      EnableState `json:"enable"`
	Description string      `json:"description"`

	// Protocol configuration
	Protocol Protocol `json:"protocol"`

	// Source configuration
	SrcIP    StringOrInt `json:"externalip"`
	SrcPorts StringOrInt `json:"externalport"`

	// Target configuration
	TargetIP    StringOrInt `json:"internalip"`
	TargetPorts StringOrInt `json:"internalport"`
}

// StringOrInt is a custom type to handle fields that can be either string or int wrapped as strings
type StringOrInt string
