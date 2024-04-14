package main

import (
	"backend-trainee-assignment-2024/config"
	"backend-trainee-assignment-2024/internal/authorization"
	"backend-trainee-assignment-2024/internal/bannerusecase"
	"backend-trainee-assignment-2024/internal/server"
	postgres "backend-trainee-assignment-2024/internal/storage"
	"log"
)

func main(){
	viperInst, err := config.LoadConfig()
	if err != nil{
		log.Fatal("Error in LoadConfig. Error:", err)
	}
	conf, err := config.ParseConfig(viperInst)
	if err != nil{
		log.Fatal("Error in ParseConfig. Error:", err)
	}
	pg, err := postgres.InitDb(conf.Postgres)
	if err != nil{
		log.Fatal(err)
	}
	defer pg.Db.Close()
	authorization := authorization.NewAuthorization(conf.JWTKey)
	buc := bannerusecase.New(pg)
	s := server.NewServer(conf.Server, buc, authorization)
	s.Start()

}