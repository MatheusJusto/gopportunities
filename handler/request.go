package handler

import "fmt"

type UpdateOpeningRequest struct {
	Role     string  `json:"role"`
	Company  string  `json:"company"`
	Location string  `json:"location"`
	Remote   *bool   `json:"remote"`
	Link     string  `json:"link"`
	Salary   float64 `json:"salary"`
}

func (r *UpdateOpeningRequest) Validate() error {
	//if any field is provided, validation is truthy
	if r.Role != "" || r.Company != "" || r.Location != "" || r.Remote != nil || r.Link != "" || r.Salary <= 0 {
		return nil
	}
	//if none of fields were provided, return falsy
	return fmt.Errorf("at least one valid field must be provided")
}
