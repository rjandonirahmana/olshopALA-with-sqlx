package handler

type Response struct {
	Meta Meta        `json:"meta"`
	Data interface{} `json:"data"`
}

type Meta struct {
	Message string
	Code    int
	Status  string
}

func APIResponse(message string, code int, status string, data interface{}) Response {
	meta := Meta{
		Message: message,
		Code:    code,
		Status:  status,
	}

	NewRensponse := Response{
		Meta: meta,
		Data: data,
	}

	return NewRensponse
}
