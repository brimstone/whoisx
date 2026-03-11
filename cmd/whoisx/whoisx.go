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
		Use:   "whoisx",
		Short: "Get owner info for an IP or two (default)",
		Long:  `Get owner info for an IP or set of IPs from one of the four whois servers.`,
		RunE:  Run,
	}

	whoisxCmd.Flags().BoolP("json", "j", false, "output in JSON")

	return whoisxCmd
}

func Run(cmd *cobra.Command, args []string) error {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		w, err := whoisx.IPLookup(line)
		if err != nil {
			log.Printf("ERROR: failed to lookup IP %s: %v", line, err)

			continue
		}

		j, err := json.Marshal(w)
		if err != nil {
			log.Printf("ERROR: failed to marshal response for IP %s: %v", line, err)

			continue
		}

		fmt.Printf("%s\n", j)
	}

	err := scanner.Err()
	if err != nil {
		return err
	}

	return nil
}
