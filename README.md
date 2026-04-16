# Bookmarks

A self-hosted bookmark manager with short URL redirects. Save websites with custom aliases and access them via short URLs like `https://bm.suram.in/r/gh`.

Built with Go and SQLite. See it live at [bm.suram.in](http://bm.suram.in).

## Features

- **Short URL redirects** — bookmark any URL with a custom alias
- **Inline editing** — click any cell to edit directly in the table
- **Search & filter** — live search across all bookmarks
- **Admin controls** — secret-based auth for adding, editing, and deleting
- **Hidden aliases** — aliases starting with `_` are hidden unless authenticated
- **Template aliases** — regex-based alias patterns for dynamic redirects
- **Chrome extension** — quickly add bookmarks from the browser ([readme](chrome_plugin/))

## Quick Start

### Prerequisites

- [Go](https://go.dev/doc/install) (or use [GVM](https://github.com/moovweb/gvm) to manage versions)
- GCC (required for `go-sqlite3` CGO bindings)
- Make

### Setup & Run

```sh
make setup   # creates dev.db from the test schema
make build   # fetches deps and compiles the binary
make run     # runs setup + build, then starts the server
```

The server starts at **http://localhost:8085**.

### Other Commands

```sh
make test    # run tests
make clean   # remove binary and dev.db
```

## Database Schema

```sql
CREATE TABLE "aliases" (
    `orig`   TEXT,
    `alias`  TEXT UNIQUE,
    `rec_id` INTEGER PRIMARY KEY AUTOINCREMENT
);
CREATE INDEX aliases_orig_index ON "aliases" (orig);

CREATE TABLE `config` (
    `key`   TEXT,
    `value` TEXT,
    PRIMARY KEY(key)
);
```

Use `internal/testdata/test_db` for an empty sqlite file with this schema, or run `scripts/setup_db.sh <path>` if you have `sqlite3` installed.

## Admin Features

To enable admin features, insert a secret into the `config` table:

```sql
INSERT INTO config ("key", "value") VALUES ("bm_secret", "YOUR_SECRET_KEY");
```

Once set:
1. Modifications (add/edit/delete) require the secret in the request.
2. Aliases starting with `_` are hidden unless the secret is provided.

## Project Structure

```
├── main.go                  # entry point
├── internal/
│   ├── db.go                # sqlite storage layer
│   ├── handlers.go          # HTTP request handlers
│   ├── router.go            # route definitions
│   ├── structures.go        # data types and response builders
│   ├── logger.go            # logging setup
│   ├── server_test.go       # integration tests
│   └── static/              # frontend assets (served at /)
│       ├── index.html
│       ├── css/style.css
│       └── js/app.js
├── scripts/
│   ├── setup_db.sh          # create a new db with schema
│   └── supervisor_aws.conf  # example supervisor config
├── chrome_plugin/            # chrome extension source
├── Makefile                  # dev workflow commands
├── DEV_SETUP.md             # detailed developer setup guide
└── ISSUES.md                # tracked issues and feature requests
```

## License

See [LICENCE](LICENCE).
