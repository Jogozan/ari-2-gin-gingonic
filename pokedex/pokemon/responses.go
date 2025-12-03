package pokemon

import "github.com/gin-gonic/gin"

type APIResponse struct {
	Data  interface{} `json:"data,omitempty"`
	Error []string    `json:"error,omitempty"`
}

func RespondOK(c *gin.Context, data interface{}) {
	c.JSON(200, APIResponse{Data: data})
}

func RespondCreated(c *gin.Context, data interface{}) {
	// RespondCreated returns an HTTP 201 with the created resource wrapped
	// into the standard APIResponse envelope. Use for successful POST creations.
	c.JSON(201, APIResponse{Data: data})
}

func RespondError(c *gin.Context, status int, errors []string) {
	// RespondError returns an error response using the APIResponse envelope.
	// `status` allows choosing the HTTP status code (400/404/429/...)
	// and `errors` is a list of human-friendly error messages.
	c.JSON(status, APIResponse{Error: errors})
}
