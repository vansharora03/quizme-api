package validator

import "vanshadhruvp/quizme-api/internal/data"

func ValidateQuiz(quiz data.Quiz) error {
	if quiz.Title == "" {
		return NewValidationError("Title is required")
	}

	return nil
}

func ValidateQuestion(question data.Question) error {
	if question.Prompt == "" {
		return NewValidationError("Prompt is required")
	} else if question.Choices == nil {
		return NewValidationError("Choices are required")
	} else if question.CorrectIndex < 0 || question.CorrectIndex >= int32(len(question.Choices)) {
		return NewValidationError("Correct index is required")
	}
	return nil

}

type ValidationError struct {
	ErrorMessage string
}

func NewValidationError(message string) *ValidationError {
	return &ValidationError{
		ErrorMessage: message,
	}
}

func (v *ValidationError) Error() string {
	return v.ErrorMessage
}
