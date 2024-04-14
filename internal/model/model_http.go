package model

type BannerData struct{
	Title	string 			`json:"title" json:"title" binding "required"`
	TextContent	string 		`json:"text" json:"text" binding "required"`
	UrlContent	string 		`json:"url" json:"url" binding "required"`
}

type RequestBodyBanner struct {
	BannerID  int64	 `json:"banner_ids"`
	TagIds    int64  `json:"tag_ids"`
	FeatureId int    `json:"feature_id"`
	Content   BannerData `json:"content"`
	IsActive  bool   `json:"is_active"`
}

type RequestFiltredBodyBanners struct {
	BannerId	int  `json:"banner_id"`
	TagIds    []int  `json:"tag_ids"`
	FeatureId int    `json:"feature_id"`
	Content   BannerData `json:"content"`
	IsActive  bool   `json:"is_active"`
}

type GetUserBannerParam struct{
	Token           string 
	TagID           int64  `form:"tag_id" binding:"required"`
	FeatureID       int64  `form:"feature_id" binding:"required"`
	UseLastRevision bool   `form:"use_last_revision" binding:"required"`
}

type GetUserFiltredBannerParam struct{
	Token           string 
	TagID           int64  
	FeatureID       int64  
	UseLastRevision bool   `form:"use_last_revision" binding:"required"`
	Limit 			int64	`form "limit" binding:"required"`
	Ofset			int64	`form "offset" binding:"required"`
}
