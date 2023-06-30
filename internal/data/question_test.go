package data

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/lib/pq"
)


func TestGetAllByQuizID(t *testing.T) {

    questionModel := QuestionModel{openTestDB(t)}

    questions, err := questionModel.GetAllByQuizID("2")
    if err != nil {
        t.Fatal(err)
    }

    if len(questions) != 2 {
        t.Fatalf("INCORRECT ENTRY LENGTH: expected %d, got %d", 2, len(questions))
    }


}


func TestAddQuestion(t *testing.T) {
    questionModel := QuestionModel{openTestDB(t)}

    tests := []struct{
        name string
        prompt string
        choices []string
        correctIndex int32
        quizID int64
        expectedError error
    }{{"Existent quizID", "prompt123", []string{"A", "Z", "Q"}, 2, 1, nil},
      {"Nonexistent quizID", "prompt123", []string{"A", "Z", "Q"}, 2, 3, ErrNoRecords},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            question := Question{
                Prompt: tt.prompt,
                Choices: tt.choices,
                CorrectIndex: tt.correctIndex,
                QuizID: tt.quizID,
            }

            err := questionModel.AddQuestion(&question)
            if err != nil {
                if err != tt.expectedError {
                    t.Fatalf("INCORRECT ERROR: expected %q, got %q", tt.expectedError, err)
                }
                return
            }

            // Query should return the same question 
            stmt := `SELECT id, created_at, version, prompt, choices, correct_index, quiz_id
            FROM question WHERE id = $1`

            ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
            defer cancel()

            row := questionModel.DB.QueryRowContext(ctx, stmt, question.ID)
            dupeQuestion := Question{}

            err = row.Scan(&dupeQuestion.ID,
                &dupeQuestion.CreatedAt,
                &dupeQuestion.Version,
                &dupeQuestion.Prompt,
                pq.Array(&dupeQuestion.Choices),
                &dupeQuestion.CorrectIndex,
                &dupeQuestion.QuizID)

            if err != nil {
                t.Fatal(err)
            }

            if !reflect.DeepEqual(question, dupeQuestion) {
                t.Fatalf("INCORRECT ENTRY: expected %v, got %v", question, dupeQuestion)
            }

        })

    }




}


func TestUpdateQuestion(t *testing.T) {
    questionModel := QuestionModel{openTestDB(t)}
    
    q := Question{
        Prompt: "updated",
        ID: 1,
        Version: 1,
        Choices: []string{"up", "date", "d"},
        CorrectIndex: 2,
    }

    err := questionModel.Update(&q)
    if err != nil {
        t.Fatal(err)
    }

    if q.Version != 2 {
        t.Fatalf("INCORRECT VERSION: expected %d, got %d", 2, q.Version)
    }

    stmt := `SELECT prompt, choices, correct_index
    FROM question
    WHERE id = $1`

    ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
    defer cancel()

    var gotQ Question

    err = questionModel.DB.QueryRowContext(ctx, stmt, q.ID).Scan(
        &gotQ.Prompt, pq.Array(&gotQ.Choices), &gotQ.CorrectIndex)

    if err != nil {
        t.Fatal(err)
    }

    if gotQ.Prompt != q.Prompt {
        t.Fatalf("INCORRECT ENTRY PROMPT: expected %q, got %q", q.Prompt, gotQ.Prompt)
    }

    if !reflect.DeepEqual(q.Choices, gotQ.Choices) {
        t.Fatalf("INCORRECT ENTRY PROMPT: expected %v, got %v", q.Choices, gotQ.Choices)
    }

    if gotQ.CorrectIndex != q.CorrectIndex {
        t.Fatalf("INCORRECT ENTRY PROMPT: expected %d, got %d", q.CorrectIndex, gotQ.CorrectIndex)
    }
}








