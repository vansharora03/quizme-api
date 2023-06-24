package data

import "testing"


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
