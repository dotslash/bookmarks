[![Build Status](https://travis-ci.com/dotslash/bookmarks.svg?branch=master)](https://travis-ci.com/dotslash/bookmarks)

Bookmarks
=========

This is the code for [bm.suram.in](http://bm.suram.in). I use this to bookmark websites with custom redirect URLs. 

The application is written in Go and uses [editable grid](https://github.com/webismymind/editablegrid) to list/search/update the bookmarked URLs. The app uses sqlite (the poor man's DB!) to persist data. It expects a `foo.db` in `bookmarks/bookmarks` directory. 

To get an empty sqlite file with the correct schema, use this online utility - [https://sqliteonline.com/#fiddle-5d6c0626e3699dmuk01a04iq](https://sqliteonline.com/#fiddle-5d6c0626e3699dmuk01a04iq)
```sql
CREATE TABLE "aliases" (
    `orig`   TEXT,
    `alias`  TEXT UNIQUE,
    `rec_id` INTEGER PRIMARY KEY AUTOINCREMENT
);
CREATE INDEX aliases_orig_index ON "aliases" (orig);

CREATE TABLE `config` (
	`key`	TEXT,
	`value`	TEXT,
	PRIMARY KEY(key)
)

```
## Admin features
The application has 2 "admin" features. To enable these features there needs to be a record with key set to `bm_secret` and value set to `{YOUR_SECRET_KEY}` in the `config` table. 
1. After this, content can be modified only if the correct secret is passed in the request.
2. Any bookmark with alias that starts with `_` will be hidden unless the secret is passed in the request.

## Installation
Clone the repository to $GOPATH/src/github.com/dotslash/bookmarks and the following to start the server.
```sh
go build
./bookmarks http://localhost:8085 8085
```
The logs will go to `~/log/bm-info.log`, `~/log/bm-error.log`

Check `scripts/supervisor_aws.conf` to see how I install the server.

Chrome Extension
================

The application also comes with a compatible chrome extension. Check out its [Readme](https://github.com/dotslash/bookmarks/tree/master/chrome_plugin) 

