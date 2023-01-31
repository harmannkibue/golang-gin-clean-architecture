package v1

import "github.com/gin-gonic/gin"

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
