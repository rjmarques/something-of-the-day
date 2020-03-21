#!/bin/bash
export HOST="${DB_HOST:-localhost}"
export POSTGRES_URL=postgres://postgres:mysecretpassword@$HOST:5432/somethingoftheday?sslmode=disable
go test -count=1 -v ./datastore/...
