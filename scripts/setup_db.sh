#!/bin/bash
#!/bin/sh
# Arg1: Where to save the db file.

if [ $1 ]; then
  if [ -e $1 ]; then
    echo "$1 already exists"
    exit 1
  fi
  echo "Saving sqlite file to $1"
else
  echo "Pass the sqlite file destination"
  exit 1
fi


sqlite3 $1 <<EOF
CREATE TABLE "aliases" (
    orig TEXT,
    alias TEXT UNIQUE,
    rec_id       INTEGER PRIMARY KEY AUTOINCREMENT
);
CREATE INDEX aliases_orig_index ON "aliases" (orig);
CREATE TABLE config (
  key TEXT,
  value TEXT,
  PRIMARY KEY(key)
)
EOF
