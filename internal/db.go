package internal

import (
	"database/sql"
	"fmt"
	"github.com/mattn/go-sqlite3"
	"io"
	"regexp"
	"strconv"
	"strings"
	"sync"

	_ "github.com/mattn/go-sqlite3"
)

// TODO/NOTE: Errors are not handled for most sqlite statements.

var registerSqliteExtended sync.Once

func getDbCxn(dbpath string) (*sql.DB, error) {
	Log.Printf("dbpath %v", dbpath)

	registerSqliteExtended.Do(func() {
		regex_match := func(re, s string) (bool, error) {
			return regexp.MatchString("^"+re+"$", s)
		}
		sql.Register("sqlite3_extended",
			&sqlite3.SQLiteDriver{
				ConnectHook: func(conn *sqlite3.SQLiteConn) error {
					return conn.RegisterFunc("regex_match", regex_match, true)
				},
			})
	})

	db, err := sql.Open("sqlite3_extended", dbpath)
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

func safeClose(c io.Closer, ctx string) {
	err := c.Close()
	if err != nil {
		Log.Println(ctx, err)
	}
}

func panicOnErr(err error) {
	if err != nil {
		Log.Fatal(err)
	}
}

func panicIf(b bool) {
	if b {
		Log.Fatal("panicIf(false)")
	}
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

	if rows.Next() {
		var orig string
		rows.Scan(&orig)
		Log.Println(orig)
		return &orig
	} else if ret := s.URLFromAliasTemplate("template:" + alias); ret != nil {
		return ret
	} else {
		return s.URLFromAliasTemplate("_template:" + alias)
	}
}

func (s *StorageInterface) URLFromAliasTemplate(expandedTemplate string) *string {
	panicIf(!strings.HasPrefix(expandedTemplate, "template:") && !strings.HasPrefix(expandedTemplate, "_template:"))
	funcStr := "URLFromAliasTemplate(" + expandedTemplate + ")"

	stmt, err := s.db.Prepare("select alias, orig from aliases where regex_match(alias, ?)")
	defer safeClose(stmt, funcStr + ".stmt")
	panicOnErr(err)

	rows, err := stmt.Query(expandedTemplate)
	defer safeClose(rows, funcStr + ".rows")
	panicOnErr(err)

	if !rows.Next() {
		return nil
	}
	var longUrl, templateReStr string

	panicOnErr(rows.Scan(&templateReStr, &longUrl))
	Log.Println(templateReStr, longUrl)

	templateRe, err := regexp.Compile(templateReStr)
	panicOnErr(err)

	ret := templateRe.ReplaceAllString(expandedTemplate, longUrl)
	return &ret
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
