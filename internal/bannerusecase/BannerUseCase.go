package bannerusecase

import (
	"backend-trainee-assignment-2024/internal/model"
	postgres "backend-trainee-assignment-2024/internal/storage"
	"fmt"
	"golang.org/x/crypto/bcrypt"
)


type BannerUseCase struct{
	psql postgres.Postgres
}

func New(psql postgres.Postgres) BannerUseCase{
	return BannerUseCase{psql : psql}
}


func (buc BannerUseCase) GetBanner(param model.GetUserBannerParam, isAdmin bool) (model.RequestBodyBanner, error){
	var bannerdata model.RequestBodyBanner
	if err := buc.psql.SelectBanner(param.FeatureID, param.TagID, &bannerdata); err != nil{
		return model.RequestBodyBanner{}, err
	}
	if !bannerdata.IsActive && !isAdmin{
		return model.RequestBodyBanner{}, fmt.Errorf("banner innactive")
	}
	return bannerdata, nil
}


func (buc BannerUseCase) GetFiltredBanners(param model.GetUserFiltredBannerParam) ([]model.RequestFiltredBodyBanners, error){
	var banners []model.RequestFiltredBodyBanners
	banners, err := buc.psql.SelectFiltredBanners(param.FeatureID, param.TagID, param.Limit, param.Ofset)
	if err != nil{
		return []model.RequestFiltredBodyBanners{}, err
	}

	
	if len(banners) == 0{
		return []model.RequestFiltredBodyBanners{}, fmt.Errorf("banners was not found")
	}
	return banners, nil
}


func (buc BannerUseCase) DeleteBanner(id int) error{
	err := buc.psql.DeleteBanner(id)
	return err
}

func (buc BannerUseCase) CreateBanner(banner model.RequestFiltredBodyBanners) error{
	if err := buc.psql.CreateBanner(banner); err != nil{
		return err
	}
	return nil
}


func (buc BannerUseCase) UpdateBanner(banner model.RequestFiltredBodyBanners) error{
	if err := buc.psql.UpdateBanner(banner); err != nil{
		return err
	}
	return nil
}

func (buc BannerUseCase) Register(user model.AuthorizationData) error{
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.MinCost)
	if err != nil {
		return err
	}
	if err := buc.psql.Register(user.Login, string(passwordHash), user.IsAdmin); err != nil{
		return err
	}
	return nil
}