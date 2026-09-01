// Copyright (c) 2026 Matt Robinson brimstone@the.narro.ws

package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/brimstone/whoisx"
	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {
	var whoisxCmd = &cobra.Command{
		Use:   "whoisx [IP ...]",
		Short: "Get owner info for an IP or two (default)",
		Long: `Get owner info for an IP or set of IPs from one of the RIR whois servers.

IP addresses can be given as positional arguments, or piped in on stdin, one per line.`,
		RunE:  Run,
	}

	whoisxCmd.Flags().BoolP("json", "j", false, "output in JSON")

	return whoisxCmd
}

func Run(cmd *cobra.Command, args []string) error {
	var ips []string

	if len(args) > 0 {
		ips = args
	} else {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line != "" {
				ips = append(ips, line)
			}
		}
		if err := scanner.Err(); err != nil {
			return err
		}
	}

	for _, ip := range ips {
		w, err := whoisx.IPLookup(ip)
		if err != nil {
			log.Printf("ERROR: failed to lookup IP %s: %v", ip, err)
			continue
		}

		j, err := json.Marshal(w)
		if err != nil {
			log.Printf("ERROR: failed to marshal response for IP %s: %v", ip, err)
			continue
		}

		fmt.Printf("%s\n", j)
	}

	return nil
}
