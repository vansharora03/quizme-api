#!/bin/bash
./sql-test-setup.sh
grc go test ./... -v -cover count=1
./sql-test-teardown.sh
