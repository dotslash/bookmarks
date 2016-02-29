package main
import (
    "strings"
)


func makeResponse(aliases []AliasInfo, serverPrefix string) response {
    if (!strings.HasSuffix(serverPrefix, "/")) {
        serverPrefix = serverPrefix + "/"
    }
    serverPrefix = serverPrefix + "red/"
    var rows []row
    for _, info := range aliases {
        rows = append(rows, info.toRow(serverPrefix))
    }
    // log.Println(rows)
    return response{Data:rows, Metadata:md}
}

type AliasInfo struct {
    Alias string
    Orig  string
    id    string
}
func (a *AliasInfo) toRow(prefix string) row {
    info := rowInfo{Fullurl:a.Orig, Alias:a.Alias, Shorturl:prefix+a.Alias}
    return row{Id:a.Alias, Values:info}
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
    Id     string  `json:"id"`
    Values rowInfo `json:"values"`
}

type response struct {
    Data     []row    `json:"data"`
    Metadata []colInfo `json:"metadata"`
}

type aliasInfos []AliasInfo

var md = []colInfo{
    colInfo{Name:"fullurl", Label:"Full Url", Datatype:"url", Bar:false, Editable:true},
    colInfo{Name:"alias", Label:"Alias", Datatype:"string", Bar:false, Editable:true},
    colInfo{Name:"shorturl", Label:"Short Url", Datatype:"url", Bar:false, Editable:false},
    colInfo{Name:"action", Label:"Actions", Datatype:"html", Bar:true, Editable:false},
}

//func main() {
//	d,_ := json.Marshal(md)
//	log.Println(string(d))
//	log.Println(md)
//}
