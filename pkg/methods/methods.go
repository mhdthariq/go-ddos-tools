package methods

import "slices"

// Layer7Methods contains all Layer 7 attack methods
var Layer7Methods = []string{
	"CFB", "BYPASS", "GET", "POST", "OVH", "STRESS", "DYN", "SLOW", "HEAD",
	"NULL", "COOKIE", "PPS", "EVEN", "GSB", "DGB", "AVB", "CFBUAM",
	"APACHE", "XMLRPC", "BOT", "BOMB", "DOWNLOADER", "KILLER", "TOR", "RHEX", "STOMP",
}

// Layer4Methods contains all Layer 4 attack methods
var Layer4Methods = []string{
	"TCP", "UDP", "SYN", "VSE", "MINECRAFT",
	"MCBOT", "CONNECTION", "CPS", "FIVEM", "FIVEM-TOKEN",
	"TS3", "MCPE", "ICMP", "OVH-UDP",
}

// AmplificationMethods contains all amplification attack methods
var AmplificationMethods = []string{
	"MEM", "NTP", "DNS", "ARD",
	"CLDAP", "CHAR", "RDP",
}

// AllMethods combines all attack methods
var AllMethods = append(append(Layer7Methods, Layer4Methods...), AmplificationMethods...)

// IsValidMethod checks if a method is valid
func IsValidMethod(method string) bool {
	return slices.Contains(AllMethods, method)
}

// IsLayer7Method checks if a method is a Layer 7 method
func IsLayer7Method(method string) bool {
	return slices.Contains(Layer7Methods, method)
}

// IsLayer4Method checks if a method is a Layer 4 method
func IsLayer4Method(method string) bool {
	return slices.Contains(Layer4Methods, method) || slices.Contains(AmplificationMethods, method)
}

// IsAmplificationMethod checks if a method is an amplification method
func IsAmplificationMethod(method string) bool {
	return slices.Contains(AmplificationMethods, method)
}
