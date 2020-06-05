package gin_di_router

import "github.com/gin-gonic/gin"

func New(group *gin.RouterGroup) *GinDIRouter {
	return &GinDIRouter{
		group: group,
	}
}
