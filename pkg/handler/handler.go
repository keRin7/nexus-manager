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

	delete := router.Group("/delete")
	{
		delete.POST("/", h.PostDelete)
	}
	auth := router.Group("/auth")
	{
		auth.GET("/", h.getAuthRoot)
		auth.POST("/", h.postAuthRoot)
	}

	repos := router.Group("/:repo")
	{
		repos.Use(h.authMiddleware)
		repos.GET("/*id", h.getRepos)
	}

	root := router.Group("/")
	{
		root.Use(h.authMiddleware)
		root.GET("/", h.getReposList)
	}

	router.Static("/assets", "./assets")

	return router
}
