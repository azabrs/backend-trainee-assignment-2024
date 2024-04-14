package postgres

import (
	"backend-trainee-assignment-2024/config"
	custom_errors "backend-trainee-assignment-2024/errors"
	"backend-trainee-assignment-2024/internal/model"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"

	_ "github.com/lib/pq"
)


type Postgres struct{
	Db *sql.DB
}

func InitDb(c config.PostgresConfig) (Postgres, error){
	str := fmt.Sprintf("dbname=%s user=%s password=%s port=%s sslmode = disable", c.DbName, c.User, c.Password, c.Port)
	db, err := sql.Open("postgres", str)
	if err != nil{
		return Postgres{},  fmt.Errorf("cant connect to db, %v", err)
	}
	return Postgres{Db: db}, nil
}

func(pg *Postgres)SelectBanner(featureID, tagID int64, bannerdata *model.RequestBodyBanner) (error){
	query := `SELECT banners.id, banners.feature_id, banner_tags.tag_id, banners_data.content, banners.is_active
	FROM banners_data
	INNER JOIN banners ON banners.data_id = banners_data.id
	INNER JOIN banner_tags ON banners.id = banner_tags.banner_id 
	WHERE banners.feature_id = $1 AND banner_tags.tag_id = $2`
	if err := pg.Db.QueryRow(query, featureID, tagID).Scan(&bannerdata.BannerID, &bannerdata.FeatureId, &bannerdata.TagIds, &bannerdata.Content, &bannerdata.IsActive); err != nil{
		if errors.Is(err, sql.ErrNoRows){
			return custom_errors.ErrBannerNotFound
		}
		return err
	}
	return nil
}

//.Scan(&bannerdata.FeatureId, &banner_id, &bannerdata.Content, &bannerdata.IsActive)
func(pg *Postgres)SelectFiltredBanners(featureID, tagID, limit, offset int64) ([]model.RequestFiltredBodyBanners, error){
	query := `SELECT banners.feature_id, banners.id, banners_data.content, banners.is_active
	FROM banners_data
	INNER JOIN banners ON banners.data_id = banners_data.id
	INNER JOIN banner_tags ON banners.id = banner_tags.banner_id 
	WHERE ($1 = -1 or banners.feature_id = $1) AND ($2 = -1 OR banner_tags.tag_id = $2)
	GROUP BY banners.id, banners.feature_id, banners_data.content, banners.is_active
`
	query_tag := `SELECT banner_tags.tag_id FROM banner_tags
				INNER JOIN banners ON banners.id = banner_tags.banner_id 
				WHERE banners.id = $1
				`
	var banners []model.RequestFiltredBodyBanners
	rows, err := pg.Db.Query(query, featureID, tagID)
	if err != nil{
		return []model.RequestFiltredBodyBanners{}, err
	}
	for rows.Next(){
		var banner model.RequestFiltredBodyBanners
		rows.Scan(&banner.FeatureId, &banner.BannerId, &banner.Content, &banner.IsActive)
		tags, err := pg.Db.Query(query_tag, banner.BannerId)
		if err != nil{
			return []model.RequestFiltredBodyBanners{}, err
		}
		for tags.Next(){
			var tag int
			tags.Scan(&tag)
			banner.TagIds = append(banner.TagIds, tag)
		}
		banners = append(banners, banner)
	}
	return banners, nil
}

func (pg *Postgres) DeleteBanner(id int) error{
	query := `
			DELETE FROM banners
			WHERE id = $1;
			`
	res, err := pg.Db.Exec(query, id)
	if c, err := res.RowsAffected(); c == 0 || err != nil{
		return custom_errors.ErrBannerNotFound
	}
	return err
}

func (pg *Postgres) CreateBanner(banner model.RequestFiltredBodyBanners) error{
	exist, err := pg.IsBannerExist(banner)
	if err != nil{
		return err
	} else if exist{
		return fmt.Errorf("banner already exist")
	}
	query := `SELECT id FROM banners_data
			ORDER BY id desc 
			LIMIT 1`
	var data_id int
	pg.Db.QueryRow(query).Scan(&data_id)

	query = `INSERT INTO banners_data (id, content)
	VALUES($1, $2)`
	if _, err = pg.Db.Exec(query, data_id + 1, banner.Content); err != nil{
		return err
	}

	query = `
	INSERT INTO banners (feature_id, data_id, is_active)
	VALUES ($1, $2, $3)
	RETURNING id;
`
	if _, err = pg.Db.Exec(query, banner.FeatureId, data_id + 1, banner.IsActive); err != nil{
		return err
	}

	query = `SELECT id FROM banners
	ORDER BY id desc 
	LIMIT 1`
	var banner_id int
	pg.Db.QueryRow(query).Scan(&banner_id)

	query = `INSERT INTO banner_tags (banner_id, tag_id)
			VALUES($1, $2)`
	for _, tag := range(banner.TagIds){
		if _, err = pg.Db.Exec(query, banner_id, tag); err != nil{
			return err
		}
	}

	return nil
}

func (pg *Postgres) UpdateBanner(newbanner model.RequestFiltredBodyBanners) error{
		for _, tag := range newbanner.TagIds {
		var banner model.RequestBodyBanner
		err := pg.SelectBanner(int64(tag), int64(newbanner.FeatureId), &banner)
		if err != nil {
			if strings.Contains(err.Error(), "no rows in result set") {
				continue
			}
		}
		if int(banner.BannerID) == newbanner.BannerId {
			continue
		}
		return fmt.Errorf("banner with new feature and tags already exist")
	}

	var dataIdOldStr string
	query := `
					SELECT data_id
					FROM banners
					WHERE id = $1
				`
	err := pg.Db.QueryRow(query, int(newbanner.BannerId)).Scan(&dataIdOldStr)
	if err != nil {
		return custom_errors.ErrBannerNotFound
	}
	dataIdOld, err := strconv.Atoi(dataIdOldStr)
	if err != nil {
		return err
	}
	query = `
	UPDATE banners
	SET feature_id = $2, is_active = $3
	WHERE id = $1;
`
	_, err = pg.Db.Exec(query, newbanner.BannerId, newbanner.FeatureId, newbanner.IsActive)
	if err != nil {
		return err
	}

	query = `
	DELETE FROM banner_tags
	WHERE banner_id = $1;
	`
	_, err = pg.Db.Exec(query, newbanner.BannerId)
	if err != nil {
	return err
	}

	for _, tag := range newbanner.TagIds {
		query = `
			INSERT INTO banner_tags (banner_id, tag_id)
			VALUES ($1, $2);
		`
		_, err = pg.Db.Exec(query, newbanner.BannerId, tag)
		if err != nil {
			return err
		}
	}

	query = `
	UPDATE banners_data
	SET content = $2
	WHERE id = $1;
	`
	_, err = pg.Db.Exec(query, dataIdOld, newbanner.Content)
	if err != nil {
		return err
	}
	return nil
}

func (pg *Postgres) IsBannerExist(banner model.RequestFiltredBodyBanners) (bool, error){

	for _, tag := range(banner.TagIds){
		var m model.RequestBodyBanner
		if err := pg.SelectBanner(int64(banner.FeatureId), int64(tag), &m); err != nil{
			if strings.Contains(err.Error(), "no rows in result set") {
				continue
			}
			return false, err
		}
		return true, nil
	}
	return false, nil
}

func (pg *Postgres) Register(login, hashPassword string, isAdmin bool) error{
	query := `INSERT INTO users(login, password_hash, is_admin) VALUES($1, $2, $3)`
	
	res, err := pg.Db.Exec(query, login, hashPassword, isAdmin)
	if err != nil{
		return err
	}
	count, err := res.RowsAffected()
	if err != nil{
		return err
	} 
	if count == 0{
		return custom_errors.ErrAlreadyRegistered
	}
	return nil
}