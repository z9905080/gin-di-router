# gin-di-router
gin tool of dependency injection with router

----

# Install

`go get -u github.com/z9905080/gin-di-router`

# Example

```go
    package main
    
    import (
        ginDIRouter "github.com/z9905080/gin-di-router"
    	"github.com/gin-gonic/gin"
    	"log"
    )
    
    
    type ExampleControllerST struct {
    }
    
    func (exC *ExampleControllerST) TestThisFunction() (ginDIRouter.APIType, []gin.HandlerFunc) {
    
    	apiType := ginDIRouter.Get
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
    
    func main() {
    	c := gin.Default()
    	diRouter := ginDIRouter.New(c.Group("test"))
    	diRouter.Register(new(ExampleControllerST))
    	c.Run(":8082")
    }

```