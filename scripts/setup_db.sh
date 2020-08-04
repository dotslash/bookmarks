#!/bin/bash
#!/bin/sh
sqlite3 src/foo.db <<EOF
CREATE TABLE "aliases" (
    orig TEXT,
    alias TEXT UNIQUE,
    rec_id       INTEGER PRIMARY KEY AUTOINCREMENT
);
CREATE INDEX aliases_orig_index ON "aliases" (orig);
CREATE TABLE config (
	key	TEXT,
	value	TEXT,
	PRIMARY KEY(key)
)
EOF
