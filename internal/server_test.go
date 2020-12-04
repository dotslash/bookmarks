package internal

import (
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gavv/httpexpect/v2"
	"github.com/stretchr/testify/require"
)

// TODO(dotslash): Try out a more modern test framework, that displays test
// exection summaries better. As of now og output looks like this -
// https://gist.github.com/dotslash/56bed5506ff96c2f7740c54c451d5986

type Outcome bool

// Failure outcome.
var FAILURE = Outcome(false)

// Bookmarks web server address used to generate short urls.
var testBookMarksServer = "https://dotslash.com"

// https://stackoverflow.com/q/22892120
func RandStringBytesMaskImprSrc(n int) string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyz0123456789"
	const (
		letterIdxBits = 6                    // 6 bits to represent a letter index
		letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
		letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
	)
	src := rand.NewSource(time.Now().UnixNano())
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}
	return string(b)
}

// Without this:
//   x, err := foo()
//   if err != nil {
//	   panic(err)
//   }
// With this:
//   x := PanicIfErr1(foo()).(string)
func PanicIfErr1(inp interface{}, err error) interface{} {
	if err != nil {
		panic(err)
	}
	return inp
}

// Without this:
//   err := foo()
//   if err != nil {
//	   panic(err)
//   }
// With this:
//   PanicIfErr(foo())
func PanicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}

// An implementation of http.Handler interface, that returns the text "redirected"
// in the response and saves the input request to `requests` slice.
type SaveRequestsHandler struct {
	requests []*http.Request
}

func (h *SaveRequestsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.requests = append(h.requests, r)
	w.Write([]byte("redirected"))
	w.WriteHeader(http.StatusOK)
}

// Creates sqlite db with the correct schema into a random file in /tmp/ and
// returns the path of the sqlite file.
func CreateTmpSqlite() string {
	tmpSqliteFile := "/tmp/bookmarks_test_" + RandStringBytesMaskImprSrc(20)
	sourceFile := "testdata/test_db"
	input := PanicIfErr1(ioutil.ReadFile(sourceFile)).([]byte)
	PanicIfErr(ioutil.WriteFile(tmpSqliteFile, input, 0644))
	return tmpSqliteFile
}

func NewTestHelper(t *testing.T) TestHelper {
	sqliteFile := CreateTmpSqlite()
	saveReqHandler := &SaveRequestsHandler{}
	server := httptest.NewServer(
		NewRouter(testBookMarksServer, sqliteFile))
	return TestHelper{
		server:         server,
		sqliteFile:     sqliteFile,
		redirectReqs:   &saveReqHandler.requests,
		redirectServer: httptest.NewServer(saveReqHandler),
		httpexpect: httpexpect.WithConfig(httpexpect.Config{
			BaseURL:  server.URL,
			Reporter: httpexpect.NewRequireReporter(t),
			Printers: []httpexpect.Printer{
				httpexpect.NewDebugPrinter(t, true),
			},
		}),
		require: require.New(t),
	}
}

type TestHelper struct {
	// bookmarks HTTP server.
	server *httptest.Server
	// A server that will be used as target for the redirects.
	// This saves the incoming requests to `redirectReqs` field.
	redirectServer *httptest.Server
	redirectReqs   *[]*http.Request
	// Path of sqlite file used for test.
	sqliteFile string
	// httpexpect and require are test utilites for checking expectations
	// and actual values returned.
	// httpexpect is tied to requests for `server`.
	httpexpect *httpexpect.Expect
	require    *require.Assertions
}

func (th *TestHelper) cleanup() {
	th.redirectServer.Close()
	th.server.Close()
}

func (th *TestHelper) AddAlias(short, url, secret string, ok ...Outcome) {
	th.require.LessOrEqual(len(ok), 1)
	expectSuccess := len(ok) == 0 || ok[0]
	res := th.httpexpect.POST("/actions/add").
		WithFormField("short", short).
		WithFormField("url", url).
		WithFormField("secret", secret).Expect()
	if expectSuccess {
		res.Status(http.StatusOK).Text().Equal("ok")
	} else {
		res.Status(http.StatusOK).Text().NotEqual("ok")
	}
}

func (th *TestHelper) UpdateOrigURL(alias, oldLong, newLong, secret string, ok ...Outcome) {
	th.require.LessOrEqual(len(ok), 1)
	expectSuccess := len(ok) == 0 || ok[0]

	res := th.httpexpect.POST("/actions/update").
		WithFormField("id", alias).
		WithFormField("newvalue", newLong).
		WithFormField("oldvalue", oldLong).
		WithFormField("colname", "orig").
		WithFormField("secret", secret).Expect()
	if expectSuccess {
		res.Status(http.StatusOK).Text().Equal("ok")
	} else {
		res.Status(http.StatusOK).Text().NotEqual("ok")
	}

}

func (th *TestHelper) UpdateAlias(oldAlias, newAlias, secret string, ok ...Outcome) {
	th.require.LessOrEqual(len(ok), 1)
	expectSuccess := len(ok) == 0 || ok[0]

	res := th.httpexpect.POST("/actions/update").
		WithFormField("id", oldAlias).
		WithFormField("newvalue", newAlias).
		WithFormField("oldvalue", oldAlias).
		WithFormField("colname", "alias").
		WithFormField("secret", secret).Expect()
	if expectSuccess {
		res.Status(http.StatusOK).Text().Equal("ok")
	} else {
		res.Status(http.StatusOK).Text().NotEqual("ok")
	}

}

func (th *TestHelper) RemoveAlias(short, secret string, ok ...Outcome) {
	th.require.LessOrEqual(len(ok), 1)
	expectSuccess := len(ok) == 0 || ok[0]
	res := th.httpexpect.POST("/actions/delete").
		WithFormField("short", short).WithFormField("secret", secret).Expect()

	if expectSuccess {
		res.Status(http.StatusOK).Text().Equal("ok")
	} else {
		res.Status(http.StatusOK).Text().NotEqual("ok")
	}
}

func (th *TestHelper) Lookup(short, long string, ok ...Outcome) {
	th.require.LessOrEqual(len(ok), 1)
	expectSuccess := len(ok) == 0 || ok[0]
	res := th.httpexpect.POST("/actions/lookup").
		WithFormField("short", short).Expect()
	if expectSuccess {
		res.Status(http.StatusOK).Body().Equal(long)
	} else {
		res.Status(http.StatusNotFound)
	}
}

func (th *TestHelper) RevLookup(long string, short []interface{}, secret string) {
	th.httpexpect.POST("/actions/revlookup").
		WithFormField("long", long).WithFormField("secret", secret).
		Expect().Status(http.StatusOK).
		JSON().Path("$.shorturls").Array().ContainsOnly(short...)
}
func (th *TestHelper) RedirectIs404(short string) {
	prev := len(*th.redirectReqs)
	th.httpexpect.GET(short).Expect().StatusRange(httpexpect.Status4xx)
	th.require.Equal(prev, len(*th.redirectReqs))
}

func (th *TestHelper) RedirectWorks(short, urlPath string) {
	prev := len(*th.redirectReqs)
	th.httpexpect.GET("/r/" + short).Expect().Text().Equal("redirected")
	th.require.Equal(prev+1, len(*th.redirectReqs))
	apath := (*th.redirectReqs)[prev].URL.Path
	th.require.Equal(urlPath, apath)
}

func (th *TestHelper) ViewBookmarks(aliases []AliasInfo, secret string) {
	expRespStruct := CreateViewResponse(aliases, testBookMarksServer)
	actRespStr := th.httpexpect.POST("/actions/view").
		WithFormField("secret", secret).Expect().
		Status(http.StatusOK).Body().Raw()
	actRespStruct := ViewResponse{}
	th.require.Nil(json.Unmarshal([]byte(actRespStr), &actRespStruct))
	// This is slightly fragile. This fails if the order of rows returned in
	// the response does not match the order of `alises` array. Using insertion
	// order seems to work.
	// If this is getting problematic, sort the rows in the struct before
	// comparing
	th.require.Equal(expRespStruct, actRespStruct)
}

func TestMutationsAndRedirects(t *testing.T) {
	th := NewTestHelper(t)
	defer th.cleanup()

	// Begin: s1, _s1 redirects dont exist.
	th.RedirectIs404("s1")
	th.RedirectIs404("_s1")

	// Add redirects for s1, _s1.
	th.AddAlias("s1", th.redirectServer.URL+"/long1", "strong-secret")
	th.AddAlias("_s1", th.redirectServer.URL+"/_long1", "strong-secret")
	// - Check redirects.
	// NOTE: Eventhough _s1 is a private alias, it will still be redirected
	//       without needing the secret. Secret is needed for lookups involving
	//       private alises and all mutations.
	th.RedirectWorks("s1", "/long1")
	th.RedirectWorks("_s1", "/_long1")

	// Update long urls for s1, _s1.
	th.UpdateOrigURL(
		"s1",                           // alias
		th.redirectServer.URL+"/long1", // oldLong
		th.redirectServer.URL+"/long2", // newLong
		"strong-secret",                // secret
	)
	th.UpdateOrigURL(
		"_s1",                           // alias
		th.redirectServer.URL+"/_long1", // oldLong
		th.redirectServer.URL+"/_long2", // newLong
		"strong-secret",                 // secret
	)
	// - Check that redirect urls are updated.
	th.RedirectWorks("s1", "/long2")
	th.RedirectWorks("_s1", "/_long2")

	// Update alises; s1->s2, _s1->_s2.
	th.UpdateAlias("s1", "s2", "strong-secret")
	th.UpdateAlias("_s1", "_s2", "strong-secret")
	// - Check redirects for _s2, s2
	th.RedirectWorks("s2", "/long2")
	th.RedirectWorks("_s2", "/_long2")
	// - Check no redirects for _s1, s1
	th.RedirectIs404("s1")
	th.RedirectIs404("_s1")

	// Remove aliases for _s2, s2
	th.RemoveAlias("s2", "strong-secret")
	th.RemoveAlias("_s2", "strong-secret")
	// - Check no redirects for _s2, s2
	th.RedirectIs404("_s2")
	th.RedirectIs404("s2")
}

func TestLookupsForPublicBookmarks(t *testing.T) {
	th := NewTestHelper(t)
	defer th.cleanup()

	// Add aliases for _s, s1, s2.
	// _s is a private alias.
	th.AddAlias("_s", th.redirectServer.URL+"/long", "strong-secret")
	th.AddAlias("s1", th.redirectServer.URL+"/long", "strong-secret")
	th.AddAlias("s2", th.redirectServer.URL+"/long", "strong-secret")
	// Alias to long url should work for all the aliases. (even the private one)
	th.Lookup("_s", th.redirectServer.URL+"/long")
	th.Lookup("s1", th.redirectServer.URL+"/long")
	th.Lookup("s2", th.redirectServer.URL+"/long")
	// RevLookup with no secret => no private alises returned.
	th.RevLookup(th.redirectServer.URL+"/long",
		[]interface{}{"s1", "s2"}, "")
	// RevLookup with secret => private alises returned.
	th.RevLookup(th.redirectServer.URL+"/long",
		[]interface{}{"s1", "s2", "_s"}, "strong-secret")
	// RevLookup with random url => nothing returned.
	th.RevLookup(th.redirectServer.URL+"/does_not_exist",
		[]interface{}{}, "strong-secret")
	// View bookmarks with correct secret => all aliases returned.
	th.ViewBookmarks(
		[]AliasInfo{
			{"_s", th.redirectServer.URL + "/long", "_s"},
			{"s1", th.redirectServer.URL + "/long", "s1"},
			{"s2", th.redirectServer.URL + "/long", "s2"},
		},
		"strong-secret")
	// View bookmarks with wrong secret => only public aliases returned.
	th.ViewBookmarks(
		[]AliasInfo{
			{"s1", th.redirectServer.URL + "/long", "s1"},
			{"s2", th.redirectServer.URL + "/long", "s2"},
		},
		"wrong-secret")
}

func TestWrongPasswordFailsMutations(t *testing.T) {
	th := NewTestHelper(t)
	defer th.cleanup()

	// Setup: Add alias and check it works.
	th.AddAlias("s1", th.redirectServer.URL+"/long1", "strong-secret")
	th.RedirectWorks("s1", "/long1")

	// Try to add alias request for s2 with wrong secret.
	// Check that redirect fails for s2.
	th.AddAlias("s2", th.redirectServer.URL+"/long1", "wrong-secret", FAILURE)
	th.RedirectIs404("s2")

	// Try to update long url for s1 with wrong secret.
	// Verify the update was not applied by checking redirect.
	th.UpdateOrigURL(
		"s1",                           // alias
		th.redirectServer.URL+"/long1", // oldLong
		th.redirectServer.URL+"/long2", // newLong
		"wrong-secret",                 // secret
		FAILURE,
	)
	th.RedirectWorks("s1", "/long1")

	// Try to update alias for s1 to s2 with wrong secret.
	// Verify the update was not applied by checking redirect.
	th.UpdateAlias("s1", "s2", "wrong-secret", FAILURE)
	th.RedirectWorks("s1", "/long1")
	th.RedirectIs404("s2")

	// Try to remove alias with wrong secret.
	// Verify the update was not applied by checking redirect.
	th.RemoveAlias("s1", "wrong-secret", FAILURE)
	th.RedirectWorks("s1", "/long1")
}

func TestDataUrl(t *testing.T) {
	th := NewTestHelper(t)
	defer th.cleanup()

	// Add and verify data alias. content type = plain text
	th.AddAlias("stext", "data:text/plain;base64,SGVsbG8=", "strong-secret")
	th.httpexpect.GET("/r/stext").Expect().
		Status(http.StatusOK).
		ContentType("text/plain").
		Body().Equal("Hello")

	// Add and verify data alias. content type = html
	th.AddAlias("shtml", "data:text/html;base64,SGVsbG8=", "strong-secret")
	th.httpexpect.GET("/r/shtml").Expect().
		Status(http.StatusOK).
		ContentType("text/html").
		Body().Equal("Hello")
}
