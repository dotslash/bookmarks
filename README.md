bookmarks
=========

The code for [yesteapea.com/bm](http://yesteapea.com/bm) .I use the site to bookmark websites with custom redirect URLs. The website is written completely in Go (earlier ~~PHP~~). What makes this more usable to me is this library called [editable grid](https://github.com/webismymind/editablegrid).

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
*Note* : Even though the website is publicly accessible, content can be modified only be me (One has to enter the secret key to edit content). There is one more useful functionality which will hide any bookmark with alias that starts with `_` unless the secret is typed.


To start the server run `go run src/*.go http://localhost:8085`  
The logs would go to `~/log/bm-info.log`, `~/log/bm-error.log`

TODO
===
- Tidy up the logs. (The present log is mix-mash of native log library and logrus)
- The startup time is very slow on ec2 (my tiny t2.micro). Investigate why.

**Credits** : This is my first task/mini-project on Go and the code is heavily borrowed from [thenewstack.io/make-a-restful-json-api-go](http://thenewstack.io/make-a-restful-json-api-go/)
