package main

import (
	library_app "github.com/Nkez/lib-app.git"
	"github.com/Nkez/lib-app.git/pkg/handler"
	"github.com/Nkez/lib-app.git/pkg/repository"
	services "github.com/Nkez/lib-app.git/pkg/services"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
)

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))

	if err := ConfigInit(); err != nil {
		logrus.Fatalf("error instaling configs: %s", err.Error())
	}

	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("error loading env variables: %s", err.Error())
	}

	db, err := repository.NewPostgresDB(repository.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
		Password: os.Getenv("DB_PASSWORD"),
	})
	if err != nil {
		logrus.Fatalf("failed to initialize db : %s", err.Error())
	}
	repos := repository.NewRepository(db)
	services := services.NewService(repos)
	handlers := handler.NewHandler(services)

	go services.SendEmail()

	if err := ConfigInit(); err != nil {
		logrus.Fatalf("error init config %s", err.Error())
	}
	srv := new(library_app.Server)
	if err := srv.Run(viper.GetString("port"), handlers.InitRouter()); err != nil {
		logrus.Fatalf("error occured while running http server: %s", err.Error())
	}

	logrus.Print("TodoApp Started")

}

func ConfigInit() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
