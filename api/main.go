package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/splinter0/api/security"
	"github.com/splinter0/api/views"
)

const (
	CRT string = "server.crt"
	KEY string = "server.key"
)

func main() {
	r := gin.Default()

	// Routes
	r.POST("/login", views.Login)
	// Middleware
	r.Use(security.AuthMiddleware())
	// Authorization required
	r.GET("/", views.Index)

	//security.CreateAdmin("splinter", "wow")

	// Start HTTPS
	err := http.ListenAndServeTLS(":443", CRT, KEY, r)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
