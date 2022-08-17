package handler

import "errors"

type registerPayload struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (p *registerPayload) Validate() error {
	if p.Username == "" {
		return errors.New("please fill username")
	}
	if p.Password == "" {
		return errors.New("please fill password")
	}

	return nil
}
