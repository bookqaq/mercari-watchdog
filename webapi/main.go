package webapi

import (
	"crypto/tls"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Boot() {
	router := gin.Default()
	router.Use(corsHandler())
	getAllRouters(router)

	srv := &http.Server{
		Addr:    ":6456",
		Handler: router,
		TLSConfig: &tls.Config{
			MinVersion: tls.VersionTLS12,
		},
	}

	if err := srv.ListenAndServeTLS("./resource/ssl/fullchain.pem", "./resource/ssl/privkey.pem"); err != nil && err != http.ErrServerClosed {
		log.Fatalf("listen: %s\n", err)
	}
}
