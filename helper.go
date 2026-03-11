// Copyright (c) 2026 Matt Robinson brimstone@the.narro.ws

package whoisx

import (
	"bufio"
	"io"
	"net"
	"regexp"
	"strings"
	"time"
)

var (
	rpslLinePattern = regexp.MustCompile(`(.+):\W+(.+)`)
)

func getTcpContent(search string, host string) (string, error) {
	var s string

	conn, err := net.DialTimeout("tcp", net.JoinHostPort(host, "43"), time.Second*28)

	defer func() { _ = conn.Close() }()

	if err != nil {
		return s, err
	}

	_, err = conn.Write([]byte(search + "\r\n"))
	if err != nil {
		return s, err
	}

	buffer, err := io.ReadAll(conn)
	if err != nil {
		return s, err
	}

	s = string(buffer[:])

	return s, err
}

func parseRPSLValue(whoisText string, class string, section string) string {
	var sectionValue = ""

	var hasIn = false

	sc := bufio.NewScanner(strings.NewReader(whoisText))

	for sc.Scan() {
		var line = sc.Text()
		if strings.HasPrefix(line, class) {
			if !hasIn {
				hasIn = true
			}
		}

		if hasIn {
			if strings.HasPrefix(line, section) {
				sectionValue = parseRPSLine(line)

				break
			}
		}
	}

	return sectionValue
}

func parseRPSLine(whoisLine string) string {
	s := rpslLinePattern.FindAllStringSubmatch(whoisLine, -1)

	if len(s) >= 1 {
		return s[0][2]
	}

	return ""
}
