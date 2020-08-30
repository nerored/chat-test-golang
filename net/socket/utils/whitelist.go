/*
	访问白名单
	支持解析v4，v6地址
*/
package utils

import (
	"net"
	"strconv"
	"strings"

	"github.com/nerored/chat-test-golang/log"
)

type WhiteList struct {
	rules map[string]bool
}

func (w *WhiteList) AddCanAccessRules(rules ...string) {
	for _, rule := range rules {
		if len(rule) == 0 {
			continue
		}

		w.rules[rule] = true
	}
}

func (w *WhiteList) AccessCheck(addr string) (canAccess bool) {
	if len(w.rules) == 0 {
		return true
	}

	ipAddr, err := net.ResolveIPAddr("", addr)

	if err != nil || ipAddr == nil {
		log.Erro("[utils-access] resolve %v failed err %v", addr, err)
		return
	}

	ipv4, ipv6 := ipAddr.IP.To4(), ipAddr.IP.To16()

	switch {
	case ipv4 != nil:
		hostFields := strings.Split(ipv4.String(), ".")

		if len(hostFields) != 4 {
			return
		}

		for rule, value := range w.rules {
			if strings.Count(rule, ".") != 4 {
				continue
			}

			if IPFieldMarch(strings.Split(rule, "."), hostFields) {
				return value
			}
		}

		return
	case ipv6 != nil:
		hostFields := IPv6Complete(strings.Split(ipv6.String(), ":"))

		if len(hostFields) != 8 {
			return
		}

		for rule, value := range w.rules {
			marchFields := IPv6Complete(strings.Split(rule, ":"))

			if len(marchFields) != 8 {
				continue
			}

			if IPFieldMarch(marchFields, hostFields) {
				return value
			}
		}

		return
	}

	return
}

func IPv6Complete(fields []string) (result []string) {
	if len(fields) <= 0 {
		return
	}

	fullCount := 8 - len(fields)
	for _, field := range fields {
		if len(field) > 0 {
			result = append(result, field)
			continue
		}

		result = append(result, "0")

		for ; fullCount > 0; fullCount-- {
			result = append(result, "0")
		}
	}

	return
}

func IPFieldMarch(str1, str2 []string) (ok bool) {
	if len(str1) != len(str2) {
		return
	}

	var turnInt func(s string) (result int64, err error)

	if len(str1) == 4 {
		turnInt = func(s string) (result int64, err error) {
			result, err = strconv.ParseInt(s, 10, 8)
			return
		}
	} else {
		turnInt = func(s string) (result int64, err error) {
			result, err = strconv.ParseInt(s, 16, 16)
			return
		}
	}

	for i := 0; i < len(str1); i++ {
		if str1[i] == "*" {
			continue
		}

		n1, err := turnInt(str1[i])

		if err != nil {
			return
		}

		n2, err := turnInt(str2[i])

		if err != nil {
			return
		}

		if n1 != n2 {
			return
		}
	}

	return true
}
