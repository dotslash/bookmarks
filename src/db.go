package main
import (
	"database/sql"
	"log"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"strconv"
	"fmt"
)
var db = getDbCxn()
var actual_secret = *getSecret()

func getDbCxn() *sql.DB {
	dir, err := os.Getwd(); if err != nil {
		log.Fatal(err)
	}
	var dbpath = dir + "/foo.db"
	log.Println("dbpath %s", dbpath)
	db, err := sql.Open("sqlite3", dbpath)
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func getSecret() *string {
	stmt, err := db.Prepare("select value from config where key = 'bm_secret'"); if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	rows, err := stmt.Query(); if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var orig string
		rows.Scan(&orig)
		log.Println(orig)
		return &orig
	}
	return nil
}

func urlFromAlias(alias string) *string {
	stmt, err := db.Prepare("select orig from aliases where alias = ?"); if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	rows, err := stmt.Query(alias); if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var orig string
		rows.Scan(&orig)
		log.Println(orig)
		return &orig
	}
	return nil
}

func addUrlAndAlias(alias string, orig string, overwrite bool) bool {
	prev_val := urlFromAlias(alias)
	if !overwrite && prev_val != nil {
		log.Println("returning because alias is already there")
		return false
	}
	tx, err := db.Begin(); if err != nil {
		log.Fatal(err)
		return false
	}
	query := "insert or replace into aliases(alias,orig) values(?,?)"
	stmt, err := tx.Prepare(query); if err != nil {
		log.Fatal(err)
		return false
	}
	defer stmt.Close()

	// We dont care about the result!
	_, err = stmt.Exec(alias, orig); if err != nil {
		log.Fatal(err)
		tx.Rollback()
		return false
	}
	tx.Commit()
	return true
}

func getAllAliases() aliasInfos {
	stmt, err := db.Prepare("select orig, alias, rec_id from aliases"); if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	rows, err := stmt.Query(); if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var ret aliasInfos;
	for rows.Next() {
		var orig string
		var alias string
		var rec_id int
		rows.Scan(&orig, &alias, &rec_id)
		str_id := strconv.Itoa(rec_id)
		ret = append(ret, AliasInfo{Alias:alias, Orig:orig, id:str_id})
	}
	return ret

}


func delByAlias(alias string) bool {
	tx, err := db.Begin(); if err != nil {
		log.Fatal(err)
		return false
	}
	stmt, err := tx.Prepare("delete from aliases where alias = ?"); if err != nil {
		log.Fatal(err)
		return false
	}
	defer stmt.Close()
	_, err = stmt.Exec(alias); if err != nil {
		log.Fatal(err)
		tx.Rollback()
		return false
	}
	tx.Commit()
	return true
}

func addAlias(orig string, alias string, secret string) string {
	if (secret != actual_secret) {
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
	if (secret != actual_secret) {
		return "Secret Did not match"
	}
	success := delByAlias(alias)
	if success {
		return "ok"
	} else {
		return "Could not write";
	}
	return "ok"
}
func updateAlias(presAlias, oldVal, newVal, colname, secret string) string{
	query := `
		update aliases
		set %s = ?
		where alias = ? and %s = ?`
	if colname != "alias" {colname = "orig"}
	query = fmt.Sprintf(query, colname, colname)
	log.Println(query)
	tx, err := db.Begin(); if err != nil {
		log.Fatal(err)
		return "Fail"
	}
	stmt, err := tx.Prepare(query); if err != nil {
		log.Fatal(err)
		return "Fail"
	}
	defer stmt.Close()
	_, err = stmt.Exec(newVal, presAlias, oldVal); if err != nil {
		log.Fatal(err)
		tx.Rollback()
		return "Fail"
	}
	tx.Commit()
	return "ok"
}
//func main() {
//	log.Println(addUrlAndAlias("abs", "abs.com", false))
//	x := urlFromAlias("abs")
//	log.Println("url %s", *x)
//	log.Println(*x)
//	log.Println(addUrlAndAlias("abs", "abs.com", false))
//	log.Println(addUrlAndAlias("abs", *x + "1", true))
//	log.Println(secret)
//	log.Println(getAllAliases())
//}
