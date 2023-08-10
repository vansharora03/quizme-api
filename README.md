# QuizMe Backend
RESTful API for QuizMe.
## How to use
## Authentication -> Becoming a registered user
- Sign up using POST on `/v1/user`
- Get temporary log in token from `/v1/user`
- Use this token in this format in the header of every request requiring authentication `Authorization: Bearer TOKEN`
- When this token expires, you must get a new one from the login route.
### Global:
- GET -> `/v1/healthcheck` : displays some information about the server environment
### On quiz:
- GET -> `/v1/quiz` : displays all quizzes in a list (no questions; sorting and filtering soon)
- GET -> `/v1/quiz/:id` : displays information about a specific quiz (yes questions)
- POST -> `/v1/quiz` : add a quiz to database, all you need to supply is json object with a `title string` field, must be registered user.
- PUT -> `/v1/quiz/:id` : update a quiz in a database, must supply a json object with `title string` and a valid `version int`, must be registered user AND owner of quiz.
- POST -> `/v1/quiz/:id/score` : will return back a score object (look at internal/data/score.go for more info) from input: `{ "answers" : int[] }`, each integer representing the index of the chosen answer. Must be registered user.
- GET -> `/v1/quiz/:id/score`: retrieves all of user's scores on a quiz. Must be registered user.
- POST -> `/v1/quiz/:id/question` : post a question to the specified quiz, must supply json object with: `prompt string, choices string[], correct_index integer`, must be registered user AND owner of quiz.
- PUT -> `/v1/quiz/:id/question/:questionID` : update a question, must supply json object with: `prompt, choices, correct_index, version`, must be registered user AND owner of quiz.
### On user:
- POST -> `/v1/user`: registers a user with input: `{"username":string, "email":string, "password":string}`
- POST -> `/v1/user/login`: logs user in and supplies auth token with expiry date, input: `{"email":string, "password":string}`

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
   
  
