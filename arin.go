// Copyright (c) 2026 Matt Robinson brimstone@the.narro.ws

package whoisx

import (
	"encoding/xml"
	"io"
	"net/http"
	"time"
)

type Org struct {
	Code string `xml:"handle,attr"`
	Name string `xml:"name,attr"`
}

type Net struct {
	Name         string `xml:"name"`
	EndAddress   string `xml:"endAddress"`
	StartAddress string `xml:"startAddress"`
	OrgInfo      Org    `xml:"orgRef"`
	LastModified string `xml:"updateDate"`
	Created      string `xml:"registrationDate"`
}

func ArinCheck(search string, whoisServer string) (*WhoisInfo, error) {
	var resp *http.Response

	var err error

	// Try up to 3 times with exponential backoff
	for i := range 3 {
		resp, err = http.Get("https://" + whoisServer + "/rest/ip/" + search)
		if err == nil {
			break
		}

		// Wait before retrying (exponential backoff: 100ms, 200ms, 400ms)
		time.Sleep(time.Duration(i+1) * 100 * time.Millisecond)
	}

	if err != nil {
		return nil, err
	}

	defer func() { _ = resp.Body.Close() }()

	return parseArinReply(search, resp.Body)
}

func parseArinReply(search string, r io.Reader) (*WhoisInfo, error) {
	v := Net{}

	err := xml.NewDecoder(r).Decode(&v)
	if err != nil {
		return nil, err
	}

	return &WhoisInfo{
		Address:      search,
		Noc:          "ARIN",
		Inetnum:      v.StartAddress + " - " + v.EndAddress,
		Netname:      v.Name,
		Organization: v.OrgInfo.Name,
		Created:      v.Created,
		LastModified: v.LastModified,
	}, nil
}
