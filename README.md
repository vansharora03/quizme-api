# QuizMe Backend
RESTful API for QuizMe.
## How to use
### Global:
- GET -> `/v1/healthcheck` : displays some information about the server environment
### On quiz:
- GET -> `/v1/quiz` : displays all quizzes in a list (no questions; sorting and filtering soon)
- GET -> `/v1/quiz/:id` : displays information about a specific quiz (yes questions)
- POST -> `/v1/quiz` : add a quiz to database, all you need to supply is json object with a `title string` field
- POST -> `/v1/quiz/:id/score` : will return your score back to you as an integer in [0,100], user must supply a json object `{ "answers" : int[] }`, each integer representing the index of the chosen answer.
- POST -> `/v1/quiz/:id/question` : post a question to the specified quiz, must supply json object with: `prompt string, choices string[], correct_index integer`
## To run the server do the following:
1. Set up psql with a `quizme` database
2. Set up a user `quizme` with a password and save this password to env variable `PGPASSWORD`
3. This command will run the server
   - `make run`
## To run tests do the following:
1. Set up psql with a `quizmetest` database
2. Set up a user `quizme` (can just use the user used to run the server)
3. Have your password for that user in env variable `PGPASSWORD` (again, can just use the user password used to run the server)
4. Tests will run with one of these script commands
   1. `make run test` if you have grc installed (Linux) OR
   2. `make run test-no-color` if you do not have grc installed
   
  
