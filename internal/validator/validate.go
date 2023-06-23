package validator

import "vanshadhruvp/quizme-api/internal/data"

func ValidateQuiz(quiz data.Quiz) error {
	if quiz.Title == "" {
		return NewValidationError("Title is required")
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
