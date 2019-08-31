package main

import (
	"strings"
)

func makeViewResponse(aliases []AliasInfo, serverPrefix string) viewResponse {
	if !strings.HasSuffix(serverPrefix, "/") {
		serverPrefix = serverPrefix + "/"
	}
	serverPrefix = serverPrefix + "r/"
	var rows []row
	for _, info := range aliases {
		rows = append(rows, info.toRow(serverPrefix))
	}
	// log.Println(rows)
	return viewResponse{Data: rows, Metadata: md}
}

func makeRevLookUpResponse(shortUrls []string) revLookupResponse {
	return revLookupResponse{ShortUrls: shortUrls}
}

// AliasInfo contains info about a shortened alias.
type AliasInfo struct {
	Alias string
	Orig  string
	id    string
}

func (a *AliasInfo) toRow(prefix string) row {
	info := rowInfo{Fullurl: a.Orig, Alias: a.Alias, Shorturl: prefix + a.Alias}
	return row{ID: a.Alias, Values: info}
}

// Column Info : Metadata
type colInfo struct {
	Name     string  `json:"name"`
	Label    string  `json:"label"`
	Datatype string  `json:"datatype"`
	Bar      bool    `json:"bar"`
	Editable bool    `json:"editable"`
	Values   *string `json:"values"`
	// dummy: this will always be nil
}

type rowInfo struct {
	Fullurl  string `json:"fullurl"`
	Alias    string `json:"alias"`
	Shorturl string `json:"shorturl"`
	Action   string `json:"action"`
}

type row struct {
	ID     string  `json:"id"`
	Values rowInfo `json:"values"`
}

type viewResponse struct {
	Data     []row     `json:"data"`
	Metadata []colInfo `json:"metadata"`
}

type revLookupResponse struct {
	ShortUrls []string `json:"shorturls"`
}

type aliasInfos []AliasInfo

var md = []colInfo{
	{Name: "fullurl", Label: "Full Url", Datatype: "url", Bar: false, Editable: true},
	{Name: "alias", Label: "Alias", Datatype: "string", Bar: false, Editable: true},
	{Name: "shorturl", Label: "Short Url", Datatype: "url", Bar: false, Editable: false},
	{Name: "action", Label: "Actions", Datatype: "html", Bar: true, Editable: false},
}

//func main() {
//	d,_ := json.Marshal(md)
//	log.Println(string(d))
//	log.Println(md)
//}
