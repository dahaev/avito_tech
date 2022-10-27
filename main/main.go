package main

import (
	"avito_tech_main/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {

	router := gin.Default()
	router.GET("/clients", GetAllClients)
	//router.GET("/orders")
	router.PATCH("/orders/:order_id", getOrderByID)
	router.POST("/clients", AddClient)
	router.POST("/orders", AddOrder)
	router.GET("/clients/:clientID", GetClientByID)
	router.Run("localhost:8080")

}

func GetAllClients(context *gin.Context) {
	clients := models.GetClients()
	context.JSONP(http.StatusOK, clients)
}

func AddClient(c *gin.Context) {
	var cli models.Client
	if err := c.BindJSON(&cli); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
	} else {
		models.AddClient(cli)
		c.IndentedJSON(http.StatusCreated, cli)
	}

}

func AddOrder(c *gin.Context) {
	var ord models.Order
	if err := c.BindJSON(&ord); err != nil {
		fmt.Println("Ошибка в AddOrder")
		c.AbortWithStatus(http.StatusBadRequest)
	} else {
		fmt.Println(ord)
		models.AddOrder(ord)
		c.IndentedJSON(http.StatusCreated, ord)
	}
}

func getOrderByID(c *gin.Context) {
	order_id := c.Param("order_id")
	order := models.GetOrderById(order_id)
	fmt.Println("Вывод в функции main: ", order)
	if order == nil {
		c.AbortWithStatus(http.StatusNotFound)
	} else {
		order.Completed = !order.Completed
		c.IndentedJSON(http.StatusOK, order)
	}

}

func GetClientByID(c *gin.Context) {
	client_id := c.Param("clientID")
	cli := models.GetClientByID(client_id)
	if cli == nil {
		c.AbortWithStatus(http.StatusNotFound)
	} else {
		c.IndentedJSON(http.StatusOK, cli)
	}
}
