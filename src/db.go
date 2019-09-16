package main

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

var db = getDbCxn()
var actualSecret = getSecret()

func getDbCxn() *sql.DB {
	Log.Info("starting DB cxn")
	dir, err := os.Getwd()
	if err != nil {
		Log.Error(err)
	}
	var dbpath = dir + "/foo.db"
	Log.Printf("dbpath %v", dbpath)
	db, err := sql.Open("sqlite3", dbpath)
	if err != nil {
		Log.Fatal(err)
	}
	Log.Info("got DB cxn")
	return db
}

func getSecret() string {
	stmt, err := db.Prepare("select value from config where key = 'bm_secret'")
	if err != nil {
		Log.Fatal(err)
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		Log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var orig string
		rows.Scan(&orig)
		Log.Println(orig)
		return orig
	}
	// No secret => return ""
	return ""
}

func urlFromAlias(alias string) *string {
	stmt, err := db.Prepare("select orig from aliases where alias = ?")
	if err != nil {
		Log.Fatal(err)
	}
	defer stmt.Close()

	rows, err := stmt.Query(alias)
	if err != nil {
		Log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var orig string
		rows.Scan(&orig)
		Log.Println(orig)
		return &orig
	}
	return nil
}

func addURLAndAlias(alias string, orig string, overwrite bool) bool {
	prevVal := urlFromAlias(alias)
	if !overwrite && prevVal != nil {
		Log.Println("returning because alias is already there")
		return false
	}
	tx, err := db.Begin()
	if err != nil {
		Log.Fatal(err)
		return false
	}
	query := "insert or replace into aliases(alias,orig) values(?,?)"
	stmt, err := tx.Prepare(query)
	if err != nil {
		Log.Fatal(err)
		return false
	}
	defer stmt.Close()

	// We dont care about the result!
	_, err = stmt.Exec(alias, orig)
	if err != nil {
		Log.Fatal(err)
		tx.Rollback()
		return false
	}
	tx.Commit()
	return true
}

func getShortUrls(secret string, orig string) []string {
	showHidden := secret == actualSecret
	statement := "select alias from aliases where orig = ? and alias not like '\\_%' escape '\\'"
	if showHidden {
		statement = "select alias from aliases where orig = ?"
	}
	stmt, err := db.Prepare(statement)
	if err != nil {
		Log.Fatal(err)
	}
	defer stmt.Close()

	rows, err := stmt.Query(orig)
	if err != nil {
		Log.Fatal(err)
	}
	defer rows.Close()
	ret := []string{}
	for rows.Next() {
		var short string
		rows.Scan(&short)
		Log.Println(short)
		ret = append(ret, short)
	}
	return ret
}

func getAllAliases(secret string) []AliasInfo {
	showHidden := secret == actualSecret
	stmt, err := db.Prepare("select orig, alias, rec_id from aliases")
	if err != nil {
		Log.Fatal(err)
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		Log.Fatal(err)
	}
	defer rows.Close()
	var ret []AliasInfo
	for rows.Next() {
		var orig string
		var alias string
		var recID int
		rows.Scan(&orig, &alias, &recID)
		if strings.HasPrefix(alias, "_") && !showHidden {
			continue
		}
		strID := strconv.Itoa(recID)
		ret = append(ret, AliasInfo{Alias: alias, Orig: orig, id: strID})
	}
	return ret

}

func delByAlias(alias string) bool {
	tx, err := db.Begin()
	if err != nil {
		Log.Fatal(err)
		return false
	}
	stmt, err := tx.Prepare("delete from aliases where alias = ?")
	if err != nil {
		Log.Fatal(err)
		return false
	}
	defer stmt.Close()
	_, err = stmt.Exec(alias)
	if err != nil {
		Log.Fatal(err)
		tx.Rollback()
		return false
	}
	tx.Commit()
	return true
}

func addAlias(orig string, alias string, secret string) string {
	if secret != actualSecret {
		return "Secret Did not match"
	}
	success := addURLAndAlias(alias, orig, false)
	if success {
		return "ok"
	}
	return "Could not write"
}

func delAlias(alias string, secret string) string {
	if secret != actualSecret {
		return "Secret Did not match"
	}
	success := delByAlias(alias)
	if success {
		return "ok"
	}
	return "Could not write"
}

func updateAlias(presAlias, oldVal, newVal, colname, secret string) string {
	if secret != actualSecret {
		return "secret did not match"
	}
	query := `
		update aliases
		set %s = ?
		where alias = ? and %s = ?`
	if colname != "alias" {
		colname = "orig"
	}
	query = fmt.Sprintf(query, colname, colname)
	// Log.Println(query)
	tx, err := db.Begin()
	if err != nil {
		Log.Fatal(err)
		return "Fail"
	}
	stmt, err := tx.Prepare(query)
	if err != nil {
		Log.Fatal(err)
		return "Fail"
	}
	defer stmt.Close()
	_, err = stmt.Exec(newVal, presAlias, oldVal)
	if err != nil {
		Log.Fatal(err)
		tx.Rollback()
		return "Fail"
	}
	tx.Commit()
	return "ok"
}

// func main() {
// 	Log.Println(addUrlAndAlias("abs", "abs.com", false))
// 	x := urlFromAlias("abs")
// 	Log.Println("url %s", *x)
// 	Log.Println(*x)
// 	Log.Println(addUrlAndAlias("abs", "abs.com", false))
// 	Log.Println(addUrlAndAlias("abs", *x + "1", true))
// 	Log.Println(secret)
// 	Log.Println(getAllAliases())
// }
