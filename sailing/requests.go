package sailing

import (
	"encoding/json"
	"errors"
	"net/mail"
)

type RegisterRequest struct {
	Email    string `json:"email" example:"mike@mock.com"`
	Name     string `json:"name" example:"mike"`
	Password string `json:"password" example:"P@ssw0rd"`
}

func (r *RegisterRequest) UnmarshalJSON(data []byte) error {
	type Input RegisterRequest
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

	*r = RegisterRequest(input)
	return nil
}

type LoginRequest struct {
	Email    string `json:"email" example:"mike@mock.com"`
	Password string `json:"password" example:"P@ssw0rd"`
}

func (r *LoginRequest) UnmarshalJSON(data []byte) error {
	type Input LoginRequest
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

	*r = LoginRequest(input)
	return nil
}
