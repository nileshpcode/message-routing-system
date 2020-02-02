package handlers

import (
	context2 "context"
	"github.com/gin-gonic/gin"
	"net/http"
	"question/system"
	"strconv"
)

func (svc *Service) GetRoute(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	route, err := svc.DBSvc.GetGateway(int64(id))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, route)
}

type CreateRouteParams struct {
	Prefix    string `json:"prefix"`
	GatewayId int64  `json:"gateway_id"`
}

func (svc *Service) CreateRoute(c *gin.Context) {
	var params CreateRouteParams

	err := c.BindJSON(&params)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if len(params.Prefix) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "prefix required",
			"params":  "prefix",
		})
		return
	}

	routeParams := map[system.RouteQueryParam]interface{}{system.RouteQueryParamPrefix: params.Prefix}
	gs, _, err := svc.QueryRoute(routeParams)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(gs) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "duplicate prefix",
			"params":  "prefix",
		})
		return
	}

	route := &system.Route{
		Prefix: params.Prefix,
		Gateway: &system.Gateway{
			ID: params.GatewayId,
		},
	}

	context := context2.Background()
	tx, err := svc.Dbo.DB.BeginTx(context, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = svc.DBSvc.CreateRoute(route)
	if err != nil {
		_ = tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	route.Gateway, err = svc.DBSvc.GetGateway(route.Gateway.ID)
	if err != nil {
		_ = tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// commit everything
	_ = tx.Commit()

	c.JSON(http.StatusCreated, route)
}

func (svc *Service) SearchRoute(c *gin.Context) {
	number := c.Param("number")

	var prefixMap map[string]*system.Route
	var err error

	_, prefixMap, err = svc.QueryRoute(map[system.RouteQueryParam]interface{}{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if prefixMap == nil {
		c.JSON(http.StatusOK, prefixMap)
		return
	}

	var bestMatchRoute system.Route

	// for each sequential substring check if prefix present
	for index := 0; index < len(number); index++ {
		prefix := number[:index+1]
		if r, ok := prefixMap[prefix]; ok {
			bestMatchRoute = *r
			continue
		}
	}
	if bestMatchRoute.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"route": "not found"})
		return
	}

	c.JSON(http.StatusOK, bestMatchRoute)
}
