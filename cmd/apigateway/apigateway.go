package main

import (
	"log"
	"net/http"
	"promo/internal/apigateway/app"
	apigateway "promo/internal/apigateway/swagger"

	"github.com/gin-gonic/gin"
)

func main() {
	server := app.NewApigatewayServer("userservice:1337")

	r := gin.Default()

	apigateway.RegisterHandlers(r, server)

	s := &http.Server{
		Handler: r,
		Addr:    "0.0.0.0:8080",
	}

	log.Fatal(s.ListenAndServe())
}
