package v1

import "github.com/gin-gonic/gin"

type ProductsHandler interface {
	Register(router *gin.RouterGroup)
}

type Router struct {
	Products ProductsHandler
}

func (router Router) RegisterV1(group *gin.RouterGroup) {
	router.Products.Register(group.Group("/products"))
}
