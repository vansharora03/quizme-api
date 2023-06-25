#!/bin/sh

psql -U quizme  -h localhost -d quizme -f sql-setup/setup/quiz_table.sql
psql -U quizme  -h localhost -d quizme -f sql-setup/setup/question_table.sql
