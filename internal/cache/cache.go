package cache

import (
	custom_errors "backend-trainee-assignment-2024/errors"
	"backend-trainee-assignment-2024/internal/model"
	postgres "backend-trainee-assignment-2024/internal/storage"
	"sync"
)

type FeatutureTag struct{
	Feature int64
	Tag int64
}

type Cache struct{
	featureMutex sync.RWMutex
	tagsMutex sync.RWMutex
	tagsFeature sync.RWMutex
	featureData map[int64][]model.RequestFiltredBodyBanners
	tagsData	map[int64][]model.RequestFiltredBodyBanners
	featureTagData map[FeatutureTag]model.RequestFiltredBodyBanners
	stor postgres.Postgres
}


func NewCache(stor postgres.Postgres) Cache{
	return Cache{
		featureData	 : make(map[int64][]model.RequestFiltredBodyBanners),
		tagsData	 : make(map[int64][]model.RequestFiltredBodyBanners),	
		featureTagData : make(map[FeatutureTag]model.RequestFiltredBodyBanners),
		stor: stor,
	}

}

func (c *Cache)UpdateCache() error{
	banners, err := c.stor.GetAllBanners()
	if err != nil{
		return err
	}
	featureData	 := make(map[int64][]model.RequestFiltredBodyBanners)
	tagsData	 := make(map[int64][]model.RequestFiltredBodyBanners)
	featureTagData := make(map[FeatutureTag]model.RequestFiltredBodyBanners)
	for _, banner := range(banners){
		featureData[int64(banner.FeatureId)] = append(featureData[int64(banner.FeatureId)], banner)
		for _, tag := range(banner.TagIds){
			tagsData[int64(tag)] = append(featureData[int64(tag)], banner)
			key := FeatutureTag{
				Feature: int64(banner.FeatureId),
				Tag: int64(tag),
			}
			featureTagData[key] = banner
		}

	}
	var wg sync.WaitGroup

	wg.Add(3)
	go func(mutex *sync.RWMutex){
		defer wg.Done()
		mutex.Lock()
		defer mutex.Unlock()
		c.featureData = featureData
	}(&c.featureMutex)

	go func(mutex *sync.RWMutex){
		defer wg.Done()
		mutex.Lock()
		defer mutex.Unlock()
		c.tagsData = tagsData
	}(&c.featureMutex)

	go func(mutex *sync.RWMutex){
		defer wg.Done()
		mutex.Lock()
		defer mutex.Unlock()
		c.featureTagData = featureTagData
	}(&c.featureMutex)


	return nil
}

func(c *Cache)BannerOnTag(tag int64) ([]model.RequestFiltredBodyBanners, error){
	c.tagsMutex.RLock()
	defer c.tagsMutex.RUnlock()
	res, ok := c.tagsData[tag]
	if !ok{
		return []model.RequestFiltredBodyBanners{}, custom_errors.ErrBannerNotFound
	}
	return res, nil
}

func(c *Cache)BannerOnFeature(feature int64) ([]model.RequestFiltredBodyBanners, error){
	c.featureMutex.RLock()
	defer c.featureMutex.RUnlock()
	res, ok := c.featureData[feature]
	if !ok{
		return []model.RequestFiltredBodyBanners{}, custom_errors.ErrBannerNotFound
	}
	return res, nil
}

func(c *Cache)BannerOnFeatureTags(feature, tag int64) (model.RequestFiltredBodyBanners, error){
	c.featureMutex.RLock()
	defer c.featureMutex.RUnlock()
	key := FeatutureTag{
		Feature: feature,
		Tag: tag,
	}
	res, ok := c.featureTagData[key]
	if !ok{
		return model.RequestFiltredBodyBanners{}, custom_errors.ErrBannerNotFound
	}
	return res, nil
}