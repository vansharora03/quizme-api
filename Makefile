test:
	./scripts/test-all.sh

test-no-color:
	./scripts/colorless-test-all.sh

run:
	go run ./... -dsn=postgres://quizme:${PGPASSWORD}@localhost/quizme

open-psql:
	./scripts/psql-open.sh

sql-up:
	./scripts/sql-setup.sh

sql-down:
	./scripts/sql-teardown.sh
 
