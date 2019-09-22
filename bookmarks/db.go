package main

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

func getDbCxn(dbpath string) (*sql.DB, error) {
	Log.Printf("dbpath %v", dbpath)
	db, err := sql.Open("sqlite3", dbpath)
	if err != nil {
		Log.Fatal(err)
		return nil, err
	}
	Log.Info("got DB cxn")
	return db, nil
}

func getSecret(db *sql.DB) string {
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

// StorageInterface is the interface to interact with the "storage layer"
// (i.e sqlite)
type StorageInterface struct {
	db     *sql.DB
	secret string
}

// NewStorageInterface creates and initializes a new storage interface object.
func NewStorageInterface(dbFile string) *StorageInterface {
	db, err := getDbCxn(dbFile)
	if err != nil {
		panic(fmt.Errorf("Failed to open sqlite connection: %v", err))
	}
	return &StorageInterface{db: db, secret: getSecret(db)}
}

// URLFromAlias returns the full url for the given `alias`.
func (s *StorageInterface) URLFromAlias(alias string) *string {
	stmt, err := s.db.Prepare("select orig from aliases where alias = ?")
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

// GetShortUrls gets the short urls for the given full url. If secret matches
// the secret in db, then returns secret aliases as well.
func (s *StorageInterface) GetShortUrls(secret string, orig string) []string {
	showHidden := secret == s.secret
	statement := "select alias from aliases where orig = ? and alias not like '\\_%' escape '\\'"
	if showHidden {
		statement = "select alias from aliases where orig = ?"
	}
	stmt, err := s.db.Prepare(statement)
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

// GetAllAliases gets all aliases. If the given secret matches the secret
// stored in db, then returns secret aliases as well.
func (s *StorageInterface) GetAllAliases(secret string) []AliasInfo {
	showHidden := secret == s.secret
	stmt, err := s.db.Prepare("select orig, alias, rec_id from aliases")
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

// AddAlias adds an `alias` for `orig` url. Returns "ok" on success.
func (s *StorageInterface) AddAlias(
	orig string, alias string, secret string) string {
	if secret != s.secret {
		return "Secret Did not match"
	}
	success := s.addAliasInternal(alias, orig, false)
	if success {
		return "ok"
	}
	return "Could not write"
}

func (s *StorageInterface) addAliasInternal(
	alias string, orig string, overwrite bool) bool {
	prevVal := s.URLFromAlias(alias)
	if !overwrite && prevVal != nil {
		Log.Println("returning because alias is already there")
		return false
	}
	tx, err := s.db.Begin()
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

// DelAlias deletes the given `alias`. Returns "ok" on success.
func (s *StorageInterface) DelAlias(alias string, secret string) string {
	if secret != s.secret {
		return "Secret Did not match"
	}
	success := s.delAliasInternal(alias)
	if success {
		return "ok"
	}
	return "Could not write"
}

func (s *StorageInterface) delAliasInternal(alias string) bool {
	tx, err := s.db.Begin()
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

// UpdateAlias updates the alias for the given long url. Returns
// "ok" on success.
func (s *StorageInterface) UpdateAlias(
	presAlias, oldVal, newVal, colname, secret string) string {
	if secret != s.secret {
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
	tx, err := s.db.Begin()
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
