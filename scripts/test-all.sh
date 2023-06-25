#!/bin/bash
./scripts/sql-test-setup.sh
grc go test ./... -v -cover count=1
./scripts/sql-test-teardown.sh
