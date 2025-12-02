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
	//TODO
	c.JSON(201, APIResponse{Data: data})
}

func RespondError(c *gin.Context, status int, errors []string) {
	//TODO
	c.JSON(status, APIResponse{Error: errors})
}
