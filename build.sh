#!/bin/sh
go build ./cmd/gb0
go build ./cmd/gb2c
upx --lzma --best gb0
upx --lzma --best gb2c
tar czf gc2.tar.gz gc2
