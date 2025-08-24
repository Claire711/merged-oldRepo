package service

import (
	"gorm.io/gorm"
	"server/global"
	"server/model/appTypes"
	"server/model/database"
	"server/model/other"
	"server/model/request"
	"server/utils"
)

type AdvertisementService struct {
}

func (advertisementService *AdvertisementService) AdvertisementInfo() (ads []database.Advertisement, total int64, err error) {
	err = global.DB.Model(&database.Advertisement{}).Count(&total).Find(&ads).Error
	if err != nil {
		return nil, 0, err
	}
	return ads, total, nil
}

func (advertisementService *AdvertisementService) AdvertisementCreate(req request.AdvertisementCreate) error {
	advertisementToCreate := database.Advertisement{
		AdImage: req.AdImage,
		Title:   req.Title,
		Link:    req.Link,
		Content: req.Content,
	}
	//事务回滚，修改广告封面图片类型
	return global.DB.Transaction(func(tx *gorm.DB) error {
		if err := utils.ChangeImagesCategory(tx, []string{advertisementToCreate.AdImage}, appTypes.AdImage); err != nil {
			return err
		}
		return tx.Create(&advertisementToCreate).Error
	})
}

func (advertisementService *AdvertisementService) AdvertisementDelete(req request.AdvertisementDelete) error {
	//判断是否为空
	if len(req.IDs) == 0 {
		return nil
	}
	return global.DB.Transaction(func(tx *gorm.DB) error {
		for _, id := range req.IDs {
			var advertisementToDelete database.Advertisement
			//id是主键，Take() 和 First() 查找结果无区别
			if err := tx.Take(&advertisementToDelete, id).Error; err != nil {
				return err
			}
			if err := utils.InitImagesCategory(tx, []string{advertisementToDelete.AdImage}); err != nil {
				return err
			}
			if err := tx.Delete(&advertisementToDelete).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (advertisementService *AdvertisementService) AdvertisementUpdate(req request.AdvertisementUpdate) error {
	updates := struct {
		Link    string `json:"link"`
		Title   string `json:"title"`
		Content string `json:"content"`
	}{
		Title:   req.Title,
		Link:    req.Link,
		Content: req.Content,
	}
	return global.DB.Take(&database.Advertisement{}, req.ID).Updates(updates).Error
}

func (advertisementService *AdvertisementService) AdvertisementList(info request.AdvertisementList) (interface{}, int64, error) {
	db := global.DB
	if info.Content != nil {
		db = db.Where("content LIKE ?", "%"+*info.Content+"%")
	}
	if info.Title != nil {
		db = db.Where("title LIKE ?", "%"+*info.Title+"%")
	}
	options := other.MySQLOption{Where: db,
		PageInfo: info.PageInfo}
	return utils.MySQLPagination(&database.Advertisement{}, options)
}
