package gin_di_router

import (
	"github.com/gin-gonic/gin"
	"log"
	"reflect"
	"regexp"
	"strings"
	"sync"
)

// APIType gin api type
type APIType int

const (
	Default APIType = iota // Default
	Get
	Delete
	Put
	Patch
	Post
	Options
	Any
	Head
)

// GinDIRouter Dependency Router of gin
type GinDIRouter struct {
	group *gin.RouterGroup
}

// Group Get Gin Group (getter)
func (diRouter *GinDIRouter) Group() *gin.RouterGroup {
	return diRouter.group
}

// SetGroup set new gin Group (setter)
func (diRouter *GinDIRouter) SetGroup(group *gin.RouterGroup) {
	diRouter.group = group
}

type service struct {
	name   string                 // name of service
	rcvr   reflect.Value          // receiver of methods for the service
	typ    reflect.Type           // type of the receiver
	method map[string]*methodType // registered methods
}

type methodType struct {
	sync.Mutex // protects counters
	method     reflect.Method
}

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

func ToSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

func (diRouter *GinDIRouter) Register(receiver interface{}) {
	service := new(service)
	service.typ = reflect.TypeOf(receiver)
	service.rcvr = reflect.ValueOf(receiver)

	// Install the methods
	service.method = suitableMethods(service.typ, true)

	for methodName, methodData := range service.method {
		result := methodData.method.Func.Call([]reflect.Value{service.rcvr})
		callAPIType := result[0].Interface().(APIType)
		callMethod := result[1].Interface().([]gin.HandlerFunc)
		newMethodName := ToSnakeCase(methodName)
		switch callAPIType {
		case Get:
			{
				diRouter.group.GET(newMethodName, callMethod...)
			}
		case Delete:
			{
				diRouter.group.DELETE(newMethodName, callMethod...)
			}
		case Put:
			{
				diRouter.group.PUT(newMethodName, callMethod...)
			}
		case Post:
			{
				diRouter.group.POST(newMethodName, callMethod...)
			}
		case Patch:
			{
				diRouter.group.PATCH(newMethodName, callMethod...)
			}
		case Options:
			{
				diRouter.group.OPTIONS(newMethodName, callMethod...)
			}
		case Any:
			{
				diRouter.group.Any(newMethodName, callMethod...)
			}
		case Head:
			{
				diRouter.group.HEAD(newMethodName, callMethod...)
			}
		case Default:
			{
				log.Println("no math api type to set route")
			}

		}
	}

}

// RegisterWithGroup Register with gin Group.
func (diRouter *GinDIRouter) RegisterWithGroup(receiver interface{}, group *gin.RouterGroup) {
	service := new(service)
	service.typ = reflect.TypeOf(receiver)
	service.rcvr = reflect.ValueOf(receiver)

	// Install the methods
	service.method = suitableMethods(service.typ, true)

	for methodName, methodData := range service.method {
		result := methodData.method.Func.Call([]reflect.Value{service.rcvr})
		callAPIType := result[0].Interface().(APIType)
		callMethod := result[1].Interface().([]gin.HandlerFunc)
		if len(result) >= 3 {
			callRouterPath := result[2].Interface().(string)
			methodName = callRouterPath
		}
		newMethodName := ToSnakeCase(methodName)
		switch callAPIType {
		case Get:
			{
				group.GET(newMethodName, callMethod...)
			}
		case Delete:
			{
				group.DELETE(newMethodName, callMethod...)
			}
		case Put:
			{
				group.PUT(newMethodName, callMethod...)
			}
		case Post:
			{
				group.POST(newMethodName, callMethod...)
			}
		case Patch:
			{
				group.PATCH(newMethodName, callMethod...)
			}
		case Options:
			{
				group.OPTIONS(newMethodName, callMethod...)
			}
		case Any:
			{
				group.Any(newMethodName, callMethod...)
			}
		case Head:
			{
				group.HEAD(newMethodName, callMethod...)
			}
		case Default:
			{
				log.Println("no math api type to set route")
			}

		}
	}
}

// suitableMethods returns suitable api methods of type, it will report
// error using log if reportErr is true.
func suitableMethods(typ reflect.Type, reportErr bool) map[string]*methodType {
	methods := make(map[string]*methodType)
	for m := 0; m < typ.NumMethod(); m++ {
		method := typ.Method(m)
		mType := method.Type
		mName := method.Name
		// Method must be exported.
		if method.PkgPath != "" {
			continue
		}
		// Method needs four ins: receiver
		if mType.NumIn() != 1 {
			if reportErr {
				log.Println("method ", mName, " has wrong number of ins:", mType.NumIn())
			}
			continue
		}

		methods[mName] = &methodType{method: method}
	}
	return methods
}
