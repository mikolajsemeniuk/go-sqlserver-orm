package sailing

import (
	"encoding/json"
	"errors"
	"net/mail"
)

type SignUpRequest struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

func (r *SignUpRequest) UnmarshalJSON(data []byte) error {
	type Input SignUpRequest
	var input Input

	if err := json.Unmarshal(data, &input); err != nil {
		return err
	}

	if _, err := mail.ParseAddress(input.Email); err != nil {
		return err
	}

	if len(input.Name) < 2 || len(input.Name) > 255 {
		return errors.New("name has to be in range of 2 to 255")
	}

	if len(input.Password) < 2 || len(input.Password) > 255 {
		return errors.New("password has to be in range of 2 to 255")
	}

	*r = SignUpRequest(input)
	return nil
}

type SignInRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (r *SignInRequest) UnmarshalJSON(data []byte) error {
	type Input SignInRequest
	var input Input

	if err := json.Unmarshal(data, &input); err != nil {
		return err
	}

	if _, err := mail.ParseAddress(input.Email); err != nil {
		return err
	}

	if len(input.Password) < 2 || len(input.Password) > 255 {
		return errors.New("password has to be in range of 2 to 255")
	}

	*r = SignInRequest(input)
	return nil
}
