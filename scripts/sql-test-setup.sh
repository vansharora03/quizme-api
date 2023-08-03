#!/bin/sh

psql -U quizme  -h localhost -d quizmetest -f sql-setup/setup/user_table.sql
psql -U quizme  -h localhost -d quizmetest -f sql-setup/setup/quiz_table.sql
psql -U quizme  -h localhost -d quizmetest -f sql-setup/setup/question_table.sql
psql -U quizme  -h localhost -d quizmetest -f sql-setup/setup/token_table.sql
psql -U quizme  -h localhost -d quizmetest -f sql-setup/setup/score_table.sql
psql -U quizme  -h localhost -d quizmetest -f sql-test-setup/insert_quiz1.sql
psql -U quizme  -h localhost -d quizmetest -f sql-test-setup/insert_quiz2.sql
psql -U quizme  -h localhost -d quizmetest -f sql-test-setup/insert_quiz2_questions.sql


