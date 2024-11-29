/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

package endpoint

import (
	"net/netip"
	"strconv"
	"strings"

	"github.com/go-kratos/kratos/v2/log"
)

func Parse(scheme, host, addr, endpoint string) string {
	scheme = strings.TrimSuffix(scheme, "://")
	naip, _ := netip.ParseAddrPort(addr)
	if endpoint == "" {
		endpoint = scheme + "://" + host + ":" + strconv.Itoa(int(naip.Port()))
	} else {
		prefix, suffix, ok := strings.Cut(endpoint, "://")
		if !ok {
			args := strings.SplitN(prefix, ":", 2)
			if len(args) == 2 {
				args[1] = strconv.Itoa(int(naip.Port()))
			} else if len(args) == 1 {
				args = append(args, strconv.Itoa(int(naip.Port())))
			} else {
				// unknown
				log.Infow("unknown http endpoint", endpoint)
			}
			endpoint = scheme + "://" + strings.Join(args, ":")
		} else {
			args := strings.SplitN(suffix, ":", 2)
			if len(args) == 2 {
				args[1] = strconv.Itoa(int(naip.Port()))
			} else if len(args) == 1 {
				args = append(args, strconv.Itoa(int(naip.Port())))
			} else {
				// unknown
				log.Infow("unknown http endpoint", endpoint)
			}
			endpoint = prefix + "://" + strings.Join(args, ":")
		}
	}
	return endpoint
}
