package server

import (
	"backend-trainee-assignment-2024/config"
	"backend-trainee-assignment-2024/internal/authorization"
	"backend-trainee-assignment-2024/internal/bannerusecase"
	"backend-trainee-assignment-2024/internal/cache"
	"time"

	"github.com/gin-gonic/gin"
)


type Server struct{
	data config.ServerConfig
	buc bannerusecase.BannerUseCase
	authorization authorization.Authorization
	cache *cache.Cache
}




func NewServer(s config.ServerConfig, buc bannerusecase.BannerUseCase, authorization authorization.Authorization, c *cache.Cache) Server{
	return Server{
		data : s,
		buc : buc,
		authorization : authorization,
		cache: c,
	}
}

func (s Server)Start(){
	go func(){
		for{
			timeout := time.After(time.Minute * 5)
			<- timeout
			s.cache.UpdateCache()
		}
	}()
	serv := gin.Default()
	serv.POST("/register", s.Register)
	serv.GET("/user_banner", s.GetUserBanner)
	serv.GET("/banner", s.GetFiltredBanners)
	serv.POST("/banner", s.CreateBanner)
	serv.PATCH("/banner/:id", s.UpdateBanner)
	serv.DELETE("/banner/:id", s.DeleteBanner)
	serv.Run(":8080")
}