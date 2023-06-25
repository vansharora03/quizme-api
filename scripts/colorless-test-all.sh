#!/bin/bash
./scripts/sql-test-setup.sh
go test ./... -v -cover -count=1
./scripts/sql-test-teardown.sh
