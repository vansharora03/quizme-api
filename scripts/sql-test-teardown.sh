#!/bin/bash

psql -U quizme  -h localhost -d quizmetest -f sql-setup/teardown/score_table.sql
psql -U quizme  -h localhost -d quizmetest -f sql-setup/teardown/question_table.sql
psql -U quizme  -h localhost -d quizmetest -f sql-setup/teardown/quiz_table.sql
psql -U quizme  -h localhost -d quizmetest -f sql-setup/teardown/user_table.sql
psql -U quizme  -h localhost -d quizmetest -f sql-setup/teardown/token_table.sql
psql -U quizme  -h localhost -d quizmetest -f sql-setup/teardown/user_table.sql
