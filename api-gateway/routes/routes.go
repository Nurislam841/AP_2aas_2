package routes

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	RegisterUserRoutes(r)
	RegisterInventoryRoutes(r)
	RegisterOrderRoutes(r)
}
