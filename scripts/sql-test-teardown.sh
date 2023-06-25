#!/bin/bash

psql -U quizme  -h localhost -d quizmetest -f sql-setup/teardown/question_table.sql
psql -U quizme  -h localhost -d quizmetest -f sql-setup/teardown/quiz_table.sql
