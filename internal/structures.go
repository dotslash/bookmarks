package internal

import "strings"

// AliasInfo contains info about a shortened alias.
type AliasInfo struct {
	Alias string
	Orig  string
	id    string
}

func (a *AliasInfo) toRow(prefix string) Row {
	info := RowInfo{Fullurl: a.Orig, Alias: a.Alias, Shorturl: prefix + a.Alias}
	return Row{ID: a.Alias, Values: info}
}

// ColInfo holds metadata for columns in the alias table showed in the UI.
type ColInfo struct {
	Name     string  `json:"name"`
	Label    string  `json:"label"`
	Datatype string  `json:"datatype"`
	Bar      bool    `json:"bar"`
	Editable bool    `json:"editable"`
	Values   *string `json:"values"`
}

// RowInfo contains information for a singe alias row in the alias table.
type RowInfo struct {
	Fullurl  string `json:"fullurl"`
	Alias    string `json:"alias"`
	Shorturl string `json:"shorturl"`
	Action   string `json:"action"`
}

// Row is a wrapper around RowInfo with an ID field.
type Row struct {
	ID     string  `json:"id"`
	Values RowInfo `json:"values"`
}

// ViewResponse is response for View call.
type ViewResponse struct {
	Data     []Row     `json:"data"`
	Metadata []ColInfo `json:"metadata"`
}

// RevLookupResponse is response for RevLookup call.
type RevLookupResponse struct {
	ShortUrls []string `json:"shorturls"`
}

// CreateViewResponse creates viewResponse from the given aliases.
func CreateViewResponse(aliases []AliasInfo, serverPrefix string) ViewResponse {
	if !strings.HasSuffix(serverPrefix, "/") {
		serverPrefix = serverPrefix + "/"
	}
	serverPrefix = serverPrefix + "r/"
	rows := make([]Row, 0)
	for _, info := range aliases {
		rows = append(rows, info.toRow(serverPrefix))
	}
	// log.Println(rows)
	metadata := []ColInfo{
		{Name: "fullurl", Label: "Full Url", Datatype: "string", Bar: false, Editable: true},
		{Name: "alias", Label: "Alias", Datatype: "string", Bar: false, Editable: true},
		{Name: "shorturl", Label: "Short Url", Datatype: "url", Bar: false, Editable: false},
		{Name: "action", Label: "Actions", Datatype: "html", Bar: true, Editable: false},
	}
	return ViewResponse{Data: rows, Metadata: metadata}
}

// ErrStruct holds an error message that can be converted to a json.
type ErrStruct struct {
	Code int    `json:"code"`
	Text string `json:"text"`
}
