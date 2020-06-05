package gin_di_router

import (
	"github.com/gin-gonic/gin"
	"log"
	"testing"
)

type ExampleControllerST struct {
}

func (exC *ExampleControllerST) TestThisFunction() (APIType, []gin.HandlerFunc) {

	apiType := Get
	middlewareFuncList := []gin.HandlerFunc{
		func(c *gin.Context) {
			log.Println("past this middleware")
		},
	}
	controlFunc := func(c *gin.Context) {
		log.Println("call get_user_data")
	}
	return apiType, append(middlewareFuncList, controlFunc)
}

func TestGinDIRouter_Register(t *testing.T) {
	c := gin.Default()
	diRouter := &GinDIRouter{
		group: c.Group("test"),
	}
	diRouter.Register(new(ExampleControllerST))

	c.Run(":8082")
}
