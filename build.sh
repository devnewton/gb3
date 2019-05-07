#!/bin/sh
go get github.com/blevesearch/bleve/...
( cd gb0 ; go build )
( cd gb2c ; go build )