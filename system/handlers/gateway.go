package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"question/system"
	"strconv"
)

func (svc *Service) GetGateway(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	gateway, err := svc.DBSvc.GetGateway(int64(id))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gateway)
}

type CreateGatewayParams struct {
	Name        string     `json:"name"`
	IpAddresses system.Ips `json:"ip_addresses"`
}

func (svc *Service) CreateGateway(c *gin.Context) {
	var params CreateGatewayParams

	err := c.BindJSON(&params)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if len(params.IpAddresses) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "ip addresses required",
			"params":  "ip_addresses",
		})
		return
	}

	if len(params.Name) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "gateway name required",
			"params":  "name",
		})
		return
	}

	gsParams := map[system.GatewayQueryParam]interface{}{system.GatewayQueryParamName: params.Name}
	gs, err := svc.QueryGateway(gsParams)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(gs) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "gateway with same name exists",
			"params":  "name",
		})
		return
	}

	ipMap := make(map[system.Ip]struct{})
	for _, ip := range params.IpAddresses {
		if _, ok := ipMap[ip]; ok {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "gateway ip addresses should not be duplicated",
				"params":  "ip_addresses",
			})

		} else {
			ipMap[ip] = struct{}{}
		}
	}

	gateway := &system.Gateway{
		Name:        params.Name,
		IpAddresses: params.IpAddresses,
	}
	err = svc.DBSvc.CreateGateway(gateway)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gateway)
}
