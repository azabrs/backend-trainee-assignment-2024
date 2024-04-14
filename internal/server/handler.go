package server

import (
	custom_errors "backend-trainee-assignment-2024/errors"
	"backend-trainee-assignment-2024/internal/model"
	jwt "backend-trainee-assignment-2024/pkg/JWT"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func (s *Server)GetUserBanner(c *gin.Context){
	var param model.GetUserBannerParam

	if err := c.ShouldBindQuery(&param); err != nil{
		log.Println("Error:", err)
		c.String(http.StatusBadRequest, fmt.Sprintf("Error: %s", err))
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	param.Token = c.GetHeader("token")
	if param.Token == ""{
		log.Println("Unable get token")
		c.String(http.StatusBadRequest, "Unable get token")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	isAdmin, err := s.authorization.IsAdmin(param.Token)
	if err != nil{
		log.Println(err)
		c.String(http.StatusUnauthorized, fmt.Sprintf("Error: %s", err))
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}


	data, err := s.buc.GetBanner(param, isAdmin)
	if err != nil{
		if err == custom_errors.ErrBannerNotFound{
			log.Println(err)
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		log.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, data.Content)
}


func (s *Server)GetFiltredBanners(c *gin.Context){
	var param model.GetUserFiltredBannerParam
	param.TagID = -1
	if tag_id, exist := c.GetQuery("tag_id"); exist && tag_id != ""{
		tag_id_int, err := strconv.Atoi(tag_id)
		if err != nil{
			log.Printf("Inccorect tag ID: %v", err)
			c.String(http.StatusBadRequest, fmt.Sprintf("Inccorect tag ID: %v", err))
			return
		}
		param.TagID = int64(tag_id_int)
	}

	param.FeatureID = -1
	if feature_id, exist := c.GetQuery("feature_id"); exist{
		feature_id_int, err := strconv.Atoi(feature_id)
		if err != nil{
			log.Printf("Inccorect feature ID: %v", err)
			c.String(http.StatusBadRequest, fmt.Sprintf("Inccorect feature ID: %v", err))
			return
		}
		param.FeatureID = int64(feature_id_int)
	}

	if err := c.ShouldBindQuery(&param); err != nil{
		log.Println("Error:", err)
		c.String(http.StatusBadRequest, fmt.Sprintf("Error: %s", err))
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	param.Token = c.GetHeader("token")
	if param.Token == ""{
		log.Println("Unable get token")
		c.String(http.StatusBadRequest, "Unable get token")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	isAdmin, err := s.authorization.IsAdmin(param.Token)
	if err != nil{
		log.Println(err)
		c.String(http.StatusUnauthorized, fmt.Sprintf("Error: %s", err))
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	if !isAdmin{
		c.AbortWithStatus(http.StatusForbidden)
		return
	}
	banners, err := s.buc.GetFiltredBanners(param)
	if err != nil{
		log.Println(err)
		c.String(http.StatusInternalServerError, fmt.Sprintf("Error: %s", err))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, banners)

}

func (s *Server)CreateBanner(c *gin.Context){
	var bodydata model.RequestFiltredBodyBanners
	c.ShouldBindBodyWith(&bodydata, binding.JSON)
	
	Token := c.GetHeader("token")
	if Token == ""{
		log.Println("Unable get token")
		c.String(http.StatusBadRequest, "Unable get token")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	isAdmin, err := s.authorization.IsAdmin(Token)
	if err != nil{
		log.Println(err)
		c.String(http.StatusUnauthorized, fmt.Sprintf("Error: %s", err))
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	if !isAdmin{
		c.AbortWithStatus(http.StatusForbidden)
		return
	}


	if err := s.buc.CreateBanner(bodydata); err != nil{
		log.Println(err)
		c.String(http.StatusInternalServerError, fmt.Sprintf("Error: %s", err))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.String(http.StatusCreated  , "")
	c.AbortWithStatus(http.StatusCreated)
}


func (s *Server)UpdateBanner(c *gin.Context){
	temp := c.Param("id")
	var bodydata model.RequestFiltredBodyBanners
	c.ShouldBindBodyWith(&bodydata, binding.JSON)
	banner_id, err := strconv.Atoi(temp)
	if err != nil{
		log.Println(err)
		c.String(http.StatusBadRequest, fmt.Sprintf("Error: %s", err))
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	bodydata.BannerId = banner_id

	Token := c.GetHeader("token")
	if Token == ""{
		log.Println("Unable get token")
		c.String(http.StatusBadRequest, "Unable get token")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	isAdmin, err := s.authorization.IsAdmin(Token)
	if err != nil{
		log.Println(err)
		c.String(http.StatusUnauthorized, fmt.Sprintf("Error: %s", err))
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	if !isAdmin{
		c.AbortWithStatus(http.StatusForbidden)
		return
	}


	if err := s.buc.UpdateBanner(bodydata); err != nil{
		if errors.Is(err, custom_errors.ErrBannerNotFound){
			log.Println(err)
			c.String(http.StatusNotFound, fmt.Sprintf("Error: %s", err))
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		log.Println(err)
		c.String(http.StatusInternalServerError, fmt.Sprintf("Error: %s", err))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.String(http.StatusNoContent, "")
	c.AbortWithStatus(http.StatusNoContent)
}

func (s *Server)DeleteBanner(c *gin.Context){
	temp := c.Param("id")
	banner_id, err := strconv.Atoi(temp)
	if err != nil{
		log.Println(err)
		c.String(http.StatusBadRequest, fmt.Sprintf("Error: %s", err))
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	Token := c.GetHeader("token")
	if Token == ""{
		log.Println("Unable get token")
		c.String(http.StatusBadRequest, "Unable get token")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	isAdmin, err := s.authorization.IsAdmin(Token)
	if err != nil{
		log.Println(err)
		c.String(http.StatusUnauthorized, fmt.Sprintf("Error: %s", err))
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	if !isAdmin{
		c.AbortWithStatus(http.StatusForbidden)
		return
	}


	if err := s.buc.DeleteBanner(banner_id); err != nil{
		if errors.Is(err, custom_errors.ErrBannerNotFound){
			log.Println(err)
			c.String(http.StatusNotFound, fmt.Sprintf("Error: %s", err))
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		log.Println(err)
		c.String(http.StatusInternalServerError, fmt.Sprintf("Error: %s", err))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.String(http.StatusNoContent, "")
	c.AbortWithStatus(http.StatusNoContent)
}

func(s * Server)Register(c *gin.Context){
	var AuthorizationData model.AuthorizationData
	if err := c.ShouldBindJSON(&AuthorizationData); err != nil{
		log.Println(err)
		c.String(http.StatusBadRequest, fmt.Sprintf("Error: %s", err))
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	var isAdmin bool
	if c.GetHeader("Admin") == "true"{
		isAdmin = true
	} else if c.GetHeader("admin") != "false" && c.GetHeader("admin") != "" {
		log.Println("incorrect admin field")
		c.String(http.StatusBadRequest, fmt.Sprintf("Error: %s", "incorrect admin field"))
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if AuthorizationData.Password == "" ||  AuthorizationData.Login == "" {
		log.Println("incorrect password or login field")
		c.String(http.StatusBadRequest, fmt.Sprintf("Error: %s", "incorrect password or login field"))
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	AuthorizationData.IsAdmin = isAdmin
	if err := s.buc.Register(AuthorizationData); err != nil{
		if errors.Is(err, custom_errors.ErrAlreadyRegistered){
			log.Println(err)
			c.String(http.StatusConflict, fmt.Sprintf("Error: %s", err))
			c.AbortWithStatus(http.StatusConflict)
			return

		}
		log.Println(err)
		c.String(http.StatusInternalServerError, fmt.Sprintf("Error: %s", err))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	tokenStr, err := jwt.CreateJWT(isAdmin, []byte(s.authorization.JWTKey), time.Now().Add(24*time.Hour)) 
	if err != nil {
		log.Println(err)
		c.String(http.StatusInternalServerError, fmt.Sprintf("Error: %s", err))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	buf := map[string]string{"Token" : tokenStr}
	c.String(http.StatusInternalServerError, fmt.Sprintf("Error: %s", err))
	c.JSON(http.StatusCreated, buf)
	c.AbortWithStatus(http.StatusCreated)
}