package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
	"os"
	"project_sem/internal/controllers"
	"project_sem/internal/logger"
	"project_sem/internal/repositories"
	"project_sem/internal/services"
)

var (
	logg *logger.Logger
)

func main() {
	parseFlags()
	initLogger()
	dbParam := fmt.Sprintf("postgres://%s:%s@%s?sslmode=disable", dbUsername, dbPassword, dbConn)
	store, err := repositories.NewDataStore(logg, dbParam)
	if err != nil {
		logg.WithError(err).
			Error("Init DB connection")
		os.Exit(1)
	}

	priceService := services.NewPriceService(logg, store)

	err = http.ListenAndServe(srvAddr, controllers.MetricsRouter(logg, priceService))
	if err != nil {
		logg.WithError(err).
			Errorf("Init http server %s", srvAddr)
		os.Exit(1)
	}
}

func initLogger() {
	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		log.Fatal(err)
	}

	l := &logrus.Logger{
		Out:   os.Stdout,
		Level: level,
		Formatter: &logrus.JSONFormatter{
			FieldMap: logrus.FieldMap{
				logrus.FieldKeyTime:  "@ts",
				logrus.FieldKeyLevel: "@level",
				logrus.FieldKeyMsg:   "@msg",
				logrus.FieldKeyFile:  "@caller",
			},
		},
	}
	logg = logger.NewLogger(l)
}
