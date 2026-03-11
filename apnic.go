// Copyright (c) 2026 Matt Robinson brimstone@the.narro.ws

package whoisx

func apnicCheck(search, whoisServer string, inetnum string) (*WhoisInfo, error) {
	whoisData, err := getTcpContent(search, whoisServer)
	if err != nil {
		return nil, err
	}

	wi := WhoisInfo{
		Address:      search,
		Inetnum:      parseRPSLValue(whoisData, inetnum, "inetnum"),
		Netname:      parseRPSLValue(whoisData, inetnum, "netname"),
		AdminC:       parseRPSLValue(whoisData, inetnum, "admin-c"),
		Country:      parseRPSLValue(whoisData, inetnum, "country"),
		Descr:        parseRPSLValue(whoisData, inetnum, "descr"),
		LastModified: parseRPSLValue(whoisData, inetnum, "changed"),
		MntBy:        parseRPSLValue(whoisData, inetnum, "mnt-by"),
		MntLower:     parseRPSLValue(whoisData, inetnum, "mnt-lower"),
		MntRoutes:    parseRPSLValue(whoisData, inetnum, "mnt-routes"),
		Source:       parseRPSLValue(whoisData, inetnum, "source"),
		TechC:        parseRPSLValue(whoisData, inetnum, "tech-c"),
		Organization: parseRPSLValue(whoisData, "irt", "irt"),

		Person: WhoisPerson{
			Name:         parseRPSLValue(whoisData, "role", "role"),
			AbuseMailbox: parseRPSLValue(whoisData, "irt", "abuse-mailbox"),
			Address:      parseRPSLValue(whoisData, "role", "address"),
			LastModified: parseRPSLValue(whoisData, "role", "changed"),
			MntBy:        parseRPSLValue(whoisData, "role", "mnt-by"),
			NicHdl:       parseRPSLValue(whoisData, "role", "nic-hdl"),
			Phone:        parseRPSLValue(whoisData, "role", "phone"),
			Source:       parseRPSLValue(whoisData, "role", "source"),
		},
	}

	return &wi, err
}
func ApnicCheck4(search, whoisServer string) (*WhoisInfo, error) {
	return apnicCheck(search, whoisServer, "inetnum")
}

func ApnicCheck6(search, whoisServer string) (*WhoisInfo, error) {
	return apnicCheck(search, whoisServer, "inet6num")
}
