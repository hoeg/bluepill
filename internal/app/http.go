package app

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/hoeg/bluepill/internal/app/adaptors"
	api "github.com/hoeg/bluepill/internal/generated"
	middleware "github.com/romulets/oapi-codegen/pkg/gin-middleware"
)

func NewHTTPService(bluepill api.ServerInterface, conf *HTTPConfig) *adaptors.HTTPSService {
	swagger, err := api.GetSwagger()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading swagger spec\n: %s", err)
		os.Exit(1)
	}
	swagger.Servers = nil
	r := gin.Default()
	r.Use(middleware.OapiRequestValidator(swagger))
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"POST", "GET"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return origin == "*"
		},
		MaxAge: 12 * time.Hour,
	}))
	api.RegisterHandlers(r, bluepill)

	log.Println("Server started on port", conf.Port)

	return adaptors.NewHTTPSService(&http.Server{
		Addr:        fmt.Sprintf(":%s", conf.Port),
		Handler:     r,
		ReadTimeout: 5 * time.Second,
	}, conf.CertFile, conf.KeyFile)
}
