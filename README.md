Bookmarks
=========

This is the code for [bm.suram.in](http://bm.suram.in). I use this to bookmark websites with custom redirect URLs. 

The application is written in Go and uses [editable grid](https://github.com/webismymind/editablegrid) to list/search/update the bookmarked URLs. The app uses sqlite (the poor man's DB!) to persist data. It expects a `foo.db` in the `src` directory. Find the db schema below
```sql
CREATE TABLE "aliases" (
    `orig` TEXT,
    `alias`TEXT UNIQUE,
    `rec_id`       INTEGER PRIMARY KEY AUTOINCREMENT
);
CREATE INDEX aliases_orig_index ON "aliases" (orig);

CREATE TABLE `config` (
	`key`	TEXT,
	`value`	TEXT,
	PRIMARY KEY(key)
)
# In config table there needs to be a record with key='bm_secret' and value='{YOUR_SECRET_KEY}' 
```
*Note* : Even though the website is publicly accessible, content can be modified only be me (One has to enter the secret key to edit content).
There is one more useful functionality which will hide any bookmark with alias that starts with `_` unless the secret is typed.


To start the server do the following from `src` directory.
```
go build -o bookmarks.bin *.go
./bookmarks.bin http://localhost:8085 8085
```
The logs will go to `~/log/bm-info.log`, `~/log/bm-error.log`


Chrome Extension
================

The application also comes with a compatible chrome extension. Check out its [Readme](https://github.com/dotslash/bookmarks/tree/master/chrome_plugin) 

