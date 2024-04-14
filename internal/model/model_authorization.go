package model


type AuthorizationData struct{
	Login string 	`form:"login" json:"login" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
	IsAdmin bool 	`form:"is_admin" json:"is_admin"`
}