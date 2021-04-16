package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/keRin7/nexus-manager/nexusmanager"
)

type Handler struct {
	nexusmanager *nexusmanager.NexusManager
}

func NewHandler(nexusmanager *nexusmanager.NexusManager) *Handler {
	return &Handler{
		nexusmanager: nexusmanager,
	}
}

func (h *Handler) InitRoutes() *gin.Engine {

	router := gin.New()

	//auth := router.Group("/auth")
	//{
	//	lists.GET("/auth", h.getAuthRoot)
	//		auth.POST("/sign-up", h.signUp)
	//		auth.POST("/sign-in", h.signIn)
	//}
	api := router.Group("/")
	{
		lists := api.Group("coolrocket")
		{
			lists.Use(h.authMiddleware)
			//	lists.POST("/", h.createList)
			lists.GET("/:id", h.getRoot)
			//	lists.GET("/:id", h.getListById)
			//	lists.PUT("/:id", h.updateList)
			//	lists.DELETE("/:id", h.deleteList)
		}
		delete := api.Group("delete")
		{
			//	lists.POST("/", h.createList)
			delete.POST("/", h.PostDelete)
			//	lists.GET("/:id", h.getListById)
			//	lists.PUT("/:id", h.updateList)
			//	lists.DELETE("/:id", h.deleteList)
		}
		auth := api.Group("auth")
		{
			auth.GET("/", h.getAuthRoot)
			auth.POST("/", h.postAuthRoot)
		}

		//	lists.POST("/", h.createList)

		//	lists.GET("/:id", h.getListById)
		//	lists.PUT("/:id", h.updateList)
		//	lists.DELETE("/:id", h.deleteList)

	}

	router.Static("/assets", "./assets")

	return router
}
