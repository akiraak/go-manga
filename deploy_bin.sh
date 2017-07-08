#!/bin/sh
export GOOS=linux
export GOARCH=amd64
BIN_PATH="~/go/src/github.com/akiraak/go-manga"

go build server.go
scp server gmanganow:$BIN_PATH/server_new
echo Built server

go build update_books.go
scp update_books gmanganow:$BIN_PATH/update_books_new
echo Built update_books

go build mailfile.go
scp mailfile gmanganow:$BIN_PATH/mailfile_new
echo Deployed mailfile
