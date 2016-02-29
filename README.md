bookmarks
=========

The code for [bm.yesteapea.com](http://bm.yesteapea.com). I use the site to bookmark websites with custom redirect URLs. The website is written completely in Go (earlier ~~PHP~~). What makes this more usable to me is this library called [editable grid](https://github.com/webismymind/editablegrid).

I use sqlite (the poor man's DB!) to persist data. The code expects 'foo.db' (yes,Im lazy) in the 'src' directory. Find the db schema below
```
CREATE TABLE "aliases" (
	`orig`	 TEXT UNIQUE,
	`alias`	 TEXT UNIQUE,
	`rec_id` INTEGER PRIMARY KEY AUTOINCREMENT
)
CREATE TABLE `config` (
	`key`	TEXT,
	`value`	TEXT,
	PRIMARY KEY(key)
)
# In config table there needs to be a record with key='bm_secret' and value='{YOUR_SECRET_KEY}' 
```
*Note* : Even though the website is publicly accessible, content can be modified only be me (One has to enter the secret key to edit content).


To start the server run `go run src/*.go http://localhost:8085`

