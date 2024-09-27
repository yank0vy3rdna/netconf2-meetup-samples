package proxy

import (
	"errors"
	"fmt"
	"math/rand"
	"net"
	"regexp"
	"strings"
)

// Функция для преобразования wildcard IP в регулярное выражение
func wildcardToRegex(wildcardIP string) string {
	parts := strings.Split(wildcardIP, ".")

	partToRegex := func(part string) string {
		if part == "*" {
			return `(25[0-5]|2[0-4][0-9]|1[0-9]{2}|[1-9]?[0-9])`
		}
		return regexp.QuoteMeta(part)
	}

	var regexParts []string
	for _, part := range parts {
		regexParts = append(regexParts, partToRegex(part))
	}

	regex := "^" + strings.Join(regexParts, `\.`) + "$"

	return regex
}

var ErrNoRulesMatched = errors.New("no rules matched")

func (u *usecase) ProxyTo(from net.IP) (uint16, error) {
	u.mu.RLock()
	config := u.currentConfig
	u.mu.RUnlock()

	for _, rule := range config.Rules.Rule {
		regex := wildcardToRegex(rule.From)
		r, err := regexp.Compile(regex)
		if err != nil {
			return 0, err
		}
		if r.MatchString(from.String()) {
			var port uint16
			if len(rule.To.Group) > 0 {

				ports := config.PortGroups.GroupById(rule.To.Group)
				port = ports[rand.Intn(len(ports))]
			} else {
				port = rule.To.Port
			}
			fmt.Printf("Request from %s forwarded to port %d\n", from.String(), port)
			return port, nil
		}
	}

	return 0, ErrNoRulesMatched
}
