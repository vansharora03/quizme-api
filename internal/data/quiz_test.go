package data

import (
	"fmt"
	"reflect"
	"testing"
	"vanshadhruvp/quizme-api/internal/validator"
)

func TestValidateQuiz(t *testing.T) {
    bigString := ""
    for i := 0; i <= 100; i++ {
        bigString += "a"
    }

    tests := []struct{
        name string
        inputQuiz Quiz
        expectedErrs map[string]string
    }{
        {"Valid quiz", Quiz{Title: "quiz"}, make(map[string]string)},
        {"Empty quiz title", Quiz{Title: ""}, map[string]string{"title":"must be provided"}},
        {"Big quiz title", Quiz{Title: bigString}, map[string]string{"title":"must be 100 characters or less"}},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            v := validator.New()

            ValidateQuiz(&v, &tt.inputQuiz)

            if !reflect.DeepEqual(v.Errors, tt.expectedErrs) {
                t.Fatalf("INCORRECT VALIDATION: expected %v, got %v", tt.expectedErrs, v.Errors)
            }
        })
    }
}

func TestGetAll(t *testing.T) {
    quizModel := QuizModel{openTestDB(t)}
    result, err := quizModel.GetAll()
    if err != nil {
        t.Fatal(err)
    }

    expectedTitles := []string{"quiz1", "quiz2"}

    expectedLength := 2

    if len(result) != expectedLength {
        t.Fatalf("INCORRECT LENGTH OF ENTRIES: expected %d, got %d", expectedLength, len(result))
    }

    for i, quiz := range result {
        if quiz.Title != expectedTitles[i] {
            t.Fatalf("INCORRECT ENTRY TITLE: expected %q, got %q", expectedTitles[i], quiz.Title)
        }
    }
}

func TestGet(t *testing.T) {
    quizModel := QuizModel{openTestDB(t)}
    tests := []struct{
        name string
        id string
        expectedErr error
        expectedTitle string
    }{{"Valid id", "1", nil, "quiz1"}, 
    {"Valid id", "2", nil, "quiz2"}, 
    {"Invalid id: does not exist", "3", ErrNoRecords, ""},}

    for _, tt := range tests {
        t.Run(tt.name, func (t *testing.T) {
            quiz, err := quizModel.Get(tt.id)

            if tt.expectedErr != nil {
                if err != tt.expectedErr {
                    t.Fatalf("INCORRECT ERROR: expected %q, got %q", tt.expectedErr, err)
                }

                return
            }

            if tt.expectedTitle != quiz.Title {
                t.Fatalf("INCORRECT ENTRY TITLE: expected %q, got %q", tt.expectedTitle, quiz.Title)
            }

        })
    }
}

/**func TestAdd(t *testing.T) {
    quizModel := QuizModel{openTestDB(t)}

    title := "quiz3test";

    quizModel.Add(title);

    stmt := `SELECT title FROM quiz WHERE title = $1`

    row := quizModel.DB.QueryRow(stmt, title)

    var gotTitle string;

    err := row.Scan(&gotTitle)
    if err != nil {
        t.Fatal(err)
    }

    if gotTitle != title {
        t.Fatalf("INCORRECT ENTRY TITLE: expected %q, got %q", title, gotTitle)
    }

}**/

func TestUpdate(t *testing.T) {
    quizmodel := QuizModel{openTestDB(t)}

    newQuiz := Quiz{
        Title: "changedquiz",
        Version: 1,
        ID: 1,
    }

    err := quizmodel.Update(&newQuiz)
    if err != nil {
        t.Fatal(err)
    }

    q, err := quizmodel.Get(fmt.Sprintf("%d", newQuiz.ID))
    if err != nil {
        t.Fatal(err)
    }

    if q.Version != newQuiz.Version {
        t.Fatalf("INCORRECT VERSION: expected %d, got %d", newQuiz.Version, q.Version)
    }

    if q.Title != newQuiz.Title {
        t.Fatalf("INCORRECT VERSION: expected %q, got %q", newQuiz.Title, q.Title)
    }

}







