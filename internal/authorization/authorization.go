package authorization


type Authorization struct{
	JWTKey string
}	

func NewAuthorization(JWTKey string) Authorization{
	return Authorization{ JWTKey: JWTKey}
}