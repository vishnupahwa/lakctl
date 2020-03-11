#!/usr/bin/env bash

go build ./cmd/lakctl/
go test ./...
rm lakctl