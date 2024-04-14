package server

import (
	"backend-trainee-assignment-2024/config"
	"backend-trainee-assignment-2024/internal/authorization"
	"backend-trainee-assignment-2024/internal/bannerusecase"

	"github.com/gin-gonic/gin"
)


type Server struct{
	data config.ServerConfig
	buc bannerusecase.BannerUseCase
	authorization authorization.Authorization
}




func NewServer(s config.ServerConfig, buc bannerusecase.BannerUseCase, authorization authorization.Authorization) Server{
	return Server{
		data : s,
		buc : buc,
		authorization : authorization,
	}
}

func (s Server)Start(){
	serv := gin.Default()
	serv.POST("/register", s.Register)
	serv.GET("/user_banner", s.GetUserBanner)
	serv.GET("/banner", s.GetFiltredBanners)
	serv.POST("/banner", s.CreateBanner)
	serv.PATCH("/banner/:id", s.UpdateBanner)
	serv.DELETE("/banner/:id", s.DeleteBanner)
	serv.Run(":8080")
}