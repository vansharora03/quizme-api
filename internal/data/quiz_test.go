package data

import "testing"

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
