package dto

import "regexp"

type SignUpDTO struct {
	Email                string `json:"email"`
	Password             string `json:"password"`
	PasswordConfirmation string `json:"password_confirmation"`
}

func (s *SignUpDTO) IsValidPassword() bool {
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(s.Password)
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(s.Password)
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(s.Password)
	hasSpecial := regexp.MustCompile(`[^a-zA-Z0-9]`).MatchString(s.Password)
	hasMinLength := len(s.Password) >= 8

	return hasUpper && hasLower && hasNumber && hasSpecial && hasMinLength
}
