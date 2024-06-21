package cmd

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"log"
	"net/http"
	"os"
	"os/signal"
	"simpleAPI/config"
	"simpleAPI/db"
	"simpleAPI/handler"
	"syscall"
	"time"
)

var serveCommand = &cobra.Command{
	Use:   "serve",
	Short: "Start the HTTP server and other modules",
	Run:   serve,
}

func serve(cmd *cobra.Command, args []string) {
	Conf := config.Init(configPath)
	dtb, err := db.NewMySQL(&Conf)
	defer dtb.Close()
	if err != nil {
		log.Panic(err)
	}
	fetchHalls := handler.FetchHalls(dtb)
	ReserveHall := handler.ReserveHall(dtb)
	engine := gin.New()
	v1 := engine.Group("/api/v1/hall-list")
	{
		v1.GET("/all-halls", fetchHalls)
	}
	V1Reserve := engine.Group("/api/v1/reserve")
	{
		V1Reserve.GET("/fetch", fetchHalls)
		V1Reserve.POST("/apply", ReserveHall)

	}
	srv := &http.Server{
		Addr:    ":3000",
		Handler: engine,
	}

	go func() {
		log.Printf("listen and serve %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
			log.Printf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("server forced to shutdown:", err)
	}

	log.Println("bye")
}
