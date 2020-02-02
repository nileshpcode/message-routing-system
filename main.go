package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
	"os"
	"question/system/handlers"
	"question/system/sqlite"
)

func main() {
	err := os.Setenv("CGO_ENABLED", "1")
	if err != nil {
		log.Fatal(err)
	}

	db, err := sqlite.Open("./test.db")
	if err != nil {
		log.Fatal(err)
	}

	svc := sqlite.DBSvc{
		Dbo: db,
	}

	err = svc.CreateTestSchema()
	if err != nil {
		log.Fatal(err)
	}

	router := gin.Default()
	router.HandleMethodNotAllowed = true
	router.NoMethod(func(c *gin.Context) {
		c.JSON(http.StatusMethodNotAllowed, gin.H{"result": false, "error": "Method Not Allowed"})
		return
	})
	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"result": false, "error": "Endpoint Not Found"})
		return
	})

	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	router.Use(gin.Recovery())

	handlerSvc := &handlers.Service{DBSvc: svc}
	router = AppendRoutes(router, handlerSvc)
	if err := router.Run(":8000"); err != nil {
		log.Fatal(err)
	}
}

func AppendRoutes(r *gin.Engine, svc *handlers.Service) (engine *gin.Engine) {

	gatewayRouter := r.Group("gateway")
	{
		gatewayRouter.POST("/", svc.CreateGateway)
		gatewayRouter.GET("/:id/", svc.GetGateway)
	}

	routeRouter := r.Group("route")
	{
		routeRouter.POST("/", svc.CreateRoute)
		routeRouter.GET("/:id/", svc.GetRoute)
	}

	searchRouteRouter := r.Group("/search/route")
	{
		searchRouteRouter.GET("/:number/", svc.SearchRoute)
	}

	return r
}
