#!/bin/bash
./sql-test-setup.sh
go test ./... -v -cover
./sql-test-teardown.sh
