#!/bin/bash
./sql-test-setup.sh
go test ./... -v -cover -count=1
./sql-test-teardown.sh
