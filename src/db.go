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
var actual_secret = *getSecret()

func getDbCxn() *sql.DB {
	Log.Info("starting DB cxn")
	dir, err := os.Getwd()
	if err != nil {
		Log.Error(err)
	}
	var dbpath = dir + "/foo.db"
	Log.Println("dbpath %s", dbpath)
	db, err := sql.Open("sqlite3", dbpath)
	if err != nil {
		Log.Fatal(err)
	}
	Log.Info("got DB cxn")
	return db
}

func getSecret() *string {
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
		return &orig
	}
	return nil
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

func addUrlAndAlias(alias string, orig string, overwrite bool) bool {
	prev_val := urlFromAlias(alias)
	if !overwrite && prev_val != nil {
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
	show_hidden := secret == actual_secret
	statement := "select alias from aliases where orig = ? and alias not like '\\_%' escape '\\'"
	if show_hidden {
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

func getAllAliases(secret string) aliasInfos {
	show_hidden := secret == actual_secret
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
	var ret aliasInfos
	for rows.Next() {
		var orig string
		var alias string
		var rec_id int
		rows.Scan(&orig, &alias, &rec_id)
		if strings.HasPrefix(alias, "_") && !show_hidden {
			continue
		}
		str_id := strconv.Itoa(rec_id)
		ret = append(ret, AliasInfo{Alias: alias, Orig: orig, id: str_id})
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
	if secret != actual_secret {
		return "Secret Did not match"
	}
	success := addUrlAndAlias(alias, orig, false)
	if success {
		return "ok"
	} else {
		return "Could not write"
	}
}
func delAlias(alias string, secret string) string {
	if secret != actual_secret {
		return "Secret Did not match"
	}
	success := delByAlias(alias)
	if success {
		return "ok"
	} else {
		return "Could not write"
	}
}
func updateAlias(presAlias, oldVal, newVal, colname, secret string) string {
	if secret != actual_secret {
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

//func main() {
//	Log.Println(addUrlAndAlias("abs", "abs.com", false))
//	x := urlFromAlias("abs")
//	Log.Println("url %s", *x)
//	Log.Println(*x)
//	Log.Println(addUrlAndAlias("abs", "abs.com", false))
//	Log.Println(addUrlAndAlias("abs", *x + "1", true))
//	Log.Println(secret)
//	Log.Println(getAllAliases())
//}
