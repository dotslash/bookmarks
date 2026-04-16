# Developer Setup Guide

This guide covers everything you need to know to set up the `bookmarks` project locally for development.

## Prerequisites
1. **Go Environment**: You need [Go](https://go.dev/doc/install) installed on your system to compile and run the application.
2. **GCC / C Compiler**: A C compiler (like `gcc`) must be installed on your system. This is required because the `github.com/mattn/go-sqlite3` driver relies on `CGO` to compile native SQLite bindings.
3. **SQLite3**: Either the `sqlite3` CLI installed, or it's simply bundled natively during Go build (no direct dependency needed since the DB file exists in test data).

## Initial Setup

For ease of use, you can use the newly provided `Makefile` to get started.

### 1. Initialize the Database
Before running the application, you need an empty database file containing the correct schema.
Run the following to set up a developer database (`dev.db`):
```bash
make setup
```
*(This command simply copies `internal/testdata/test_db` to `dev.db` in your root folder. Alternatively, you can use `./scripts/setup_db.sh dev.db` but that requires the `sqlite3` tool to be installed).*

### 2. Build the Application
To fetch dependencies and compile the server:
```bash
make build
```

### 3. Run Locally
To spin up the web server locally (running on `http://localhost:8085` using `dev.db`):
```bash
make run
```
You can now access the application at http://localhost:8085.

### 4. Running Tests
To run the project tests:
```bash
make test
```

## Admin Features Setup
If you want to test the admin features (like adding/editing/deleting bookmarks and making aliases that begin with `_` hidden), you must define the `bm_secret` key:

1. Open `dev.db` using your SQLite viewer (or `sqlite3 dev.db`).
2. Insert a secret:
```sql
INSERT INTO config ("key", "value") VALUES ("bm_secret", "YOUR_SECRET_KEY");
```
3. Whenever issuing requests that modify data, you will now need to pass this secret along in the request.
