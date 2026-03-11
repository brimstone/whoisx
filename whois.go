// Copyright (c) 2026 Matt Robinson brimstone@the.narro.ws

package whoisx

import (
	"errors"
	"net"
	"strings"
)

//go:generate go run gen.go

type lookupFunc func(string, string) (*WhoisInfo, error)

var lookupFunctions = map[string]lookupFunc{
	"whois.afrinic.net4": AfrinicCheck,
	"whois.apnic.net4":   ApnicCheck4,
	"whois.arin.net4":    ArinCheck,
	"whois.lacnic.net4":  LacnicCheck,
	"whois.ripe.net4":    RipeCheck4,

	"whois.afrinic.net6": AfrinicCheck,
	"whois.apnic.net6":   ApnicCheck6,
	"whois.arin.net6":    ArinCheck,
	"whois.lacnic.net6":  LacnicCheck,
	"whois.ripe.net6":    RipeCheck6,
}

func iplookuprecurse(ipaddr string, whoisServer string, ipVersion string) (*WhoisInfo, error) {
	lf := lookupFunctions[whoisServer+ipVersion]
	if lf == nil {
		return nil, errors.New("Unable to find whois function for: " + whoisServer)
	}

	info, err := lf(ipaddr, whoisServer)
	if err != nil {
		return info, err
	}

	// If RIPE is asked, try ARIN
	if info.Netname == "NON-RIPE-NCC-MANAGED-ADDRESS-BLOCK" {
		return iplookuprecurse(ipaddr, "whois.arin.net", ipVersion)
	}
	// If AFRINIC is asked, try ARIN
	if info.Netname == "IANA-BLK" {
		return iplookuprecurse(ipaddr, "whois.arin.net", ipVersion)
	}
	// If APNIC is asked, try ARIN
	if strings.HasPrefix(info.Netname, "IANA-NETBLOCK") {
		return iplookuprecurse(ipaddr, "whois.arin.net", ipVersion)
	}
	// If ARIN is asked, they'll tell you where to go
	if whoisServer == "whois.arin.net" {
		switch info.Organization {
		case "Asia Pacific Network Information Centre":
			return iplookuprecurse(ipaddr, "whois.apnic.net", ipVersion)
		case "RIPE Network Coordination Centre":
			return iplookuprecurse(ipaddr, "whois.ripe.net", ipVersion)
		case "African Network Information Center":
			return iplookuprecurse(ipaddr, "whois.afiric.net", ipVersion)
		}
	}

	return info, err
}

// IPLookup function that returns IP information at provider and returns information.
func IPLookup(ipaddr string) (*WhoisInfo, error) {
	initIPv6()

	var ipVersion string

	var whoisServer string

	ip := net.ParseIP(ipaddr)

	if ip4 := ip.To4(); ip4 != nil {
		ipVersion = "4"
		whoisServer = getIPv4Server(ip4)
	} else if ip6 := ip.To16(); ip6 != nil {
		ipVersion = "6"
		whoisServer = getIPv6Server(ip6)
	} else {
		return nil, errors.New("Invalid IP address: " + ipaddr)
	}

	if whoisServer == "" {
		return nil, errors.New("unable to find WhoIs-Server for: " + ipaddr)
	}

	return iplookuprecurse(ipaddr, whoisServer, ipVersion)
}

// getIPv4Server returns the whois server fot the given IPv4 address.
func getIPv4Server(ip net.IP) string {
	return ipv4prefixes[ip[0]]
}

// getIPv6Server returns the whois server fot the given IPv6 address.
func getIPv6Server(ip net.IP) string {
	for i := range ipv6prefixes {
		entry := &ipv6prefixes[i]
		if entry.net.Contains(ip) {
			return entry.whois
		}
	}

	return ""
}

// WhoisInfo struct with information on IP address range.
type WhoisInfo struct {
	Address      string      `json:"address"`
	Noc          string      `json:"noc,omitempty"`
	Inetnum      string      `json:"inetnum,omitempty"`
	Netname      string      `json:"netname,omitempty"`
	Descr        string      `json:"descr,omitempty"`
	Country      string      `json:"country,omitempty"`
	Organization string      `json:"organization,omitempty"`
	AdminC       string      `json:"adminc,omitempty"`
	TechC        string      `json:"techc,omitempty"`
	MntLower     string      `json:"mnt_lower,omitempty"`
	Status       string      `json:"status,omitempty"`
	MntBy        string      `json:"mnt_by,omitempty"`
	Created      string      `json:"created,omitempty"`
	LastModified string      `json:"last_modified,omitempty"`
	Source       string      `json:"source,omitempty"`
	MntRoutes    string      `json:"mnt_routes,omitempty"`
	Person       WhoisPerson `json:"person,omitzero"`
	Route        WhoisRoute  `json:"route,omitzero"`
}

// WhoisPerson struct for Person information from provider.
type WhoisPerson struct {
	Name         string `json:"name,omitempty"`
	Address      string `json:"address,omitempty"`
	Phone        string `json:"phone,omitempty"`
	AbuseMailbox string `json:"abuse_mailbox,omitempty"`
	NicHdl       string `json:"nichdl,omitempty"`
	MntBy        string `json:"mnt_by,omitempty"`
	Created      string `json:"created,omitempty"`
	LastModified string `json:"last_modified,omitempty"`
	Source       string `json:"source,omitempty"`
}

// WhoisRoute struct for Route and Network information from provider.
type WhoisRoute struct {
	Route        string `json:"route,omitempty"`
	Descr        string `json:"descr,omitempty"`
	Origin       string `json:"origin,omitempty"`
	MntBy        string `json:"mnt_by,omitempty"`
	Created      string `json:"created,omitempty"`
	LastModified string `json:"last_modified,omitempty"`
	Source       string `json:"source,omitempty"`
}
