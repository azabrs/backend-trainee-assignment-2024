package authorization

import (
	custom_errors "backend-trainee-assignment-2024/errors"
	jwt "backend-trainee-assignment-2024/pkg/JWT"
)


type Authorization struct{
	JWTKey string
}	

func NewAuthorization(JWTKey string) Authorization{
	return Authorization{ JWTKey: JWTKey}
}


func (auth Authorization)IsAdmin(Token string) (bool, error){
	if Token == "" {
		return false, custom_errors.ErrNoTokenProvided
	}
	isAdmin, err := jwt.CheckIsAdminInJWT(Token, auth.JWTKey) 
	if err != nil {
		return false, custom_errors.ErrTokenIsInvalid
	}

	return isAdmin, nil
}