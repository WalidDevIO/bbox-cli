package cli

import (
	bboxclient "bbox-cli/client"
	"fmt"
)

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-2] + ".."
}

func defaultIfEmpty(s, def string) string {
	if s == "" {
		return def
	}
	return s
}

func repeatString(s string, count int) string {
	result := ""
	for i := 0; i < count; i++ {
		result += s
	}
	return result
}

func readInput(prompt string) string {
	fmt.Print(prompt)
	var input string
	fmt.Scanln(&input)
	return input
}

func parseIPOrPort(input string) (bboxclient.StringOrInt, bboxclient.EnableState) {
	if input == "" {
		return bboxclient.StringOrInt(""), bboxclient.Disabled
	}
	if input[0] == '!' {
		return bboxclient.StringOrInt(input[1:]), bboxclient.Enabled
	}
	return bboxclient.StringOrInt(input), bboxclient.Disabled
}

func parseProtocols(input string) bboxclient.Protocol {
	if input == "" {
		return bboxclient.ProtocolAny
	}
	return bboxclient.Protocol(input)
}

func parseEnable(input string) bboxclient.EnableState {
	if input == "y" || input == "Y" {
		return bboxclient.Enabled
	}
	return bboxclient.Disabled
}
