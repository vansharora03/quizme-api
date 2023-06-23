#!/bin/bash
./sql-test-setup.sh
grc go test ./... -v -cover
./sql-test-teardown.sh
