package errors

type ValidationError struct {
	Field    string      `json:"field,omitempty"`
	Tag      string      `json:"tag,omitempty"`
	Param    interface{} `json:"param,omitempty"`
	Value    interface{} `json:"value,omitempty"`
	ErrorStr string      `json:"message,omitempty"`
}

func (t *ValidationError) Error() string {
	return t.ErrorStr
}
