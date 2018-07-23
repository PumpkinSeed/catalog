package catalog

import "errors"

var (
	ErrUndefinedService      = errors.New("undefined service")
	ErrServiceRequestInvalid = errors.New("service request must contain at least an ID")
)
