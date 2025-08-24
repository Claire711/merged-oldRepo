package service

import (
	"gorm.io/gorm"
	"mime/multipart"
	"server/global"
	"server/model/appTypes"
	"server/model/database"
	"server/model/other"
	"server/model/request"
	"server/utils"

	"server/utils/upload"
)

type ImageService struct {
}

func (imageService *ImageService) ImageUpload(file *multipart.FileHeader) (string, error) {
	OSS := upload.NewOSS()
	url, filename, err := OSS.UploadImage(file)
	if err != nil {
		return "", err
	}
	return url, global.DB.Create(&database.Image{
		URL:      url,
		Name:     filename,
		Storage:  global.Config.System.Storage(),
		Category: appTypes.Null,
	}).Error
}
func (imageService *ImageService) ImageDelete(req request.ImageDelete) error {
	if len(req.IDs) == 0 {
		return nil
	}
	var images []database.Image
	if err := global.DB.Find(&images, req.IDs).Error; err != nil {
		return err
	}
	// 遍历 images 切片中的每个图片
	for _, image := range images {
		// 对每个图片执行事务操作
		if err := global.DB.Transaction(func(tx *gorm.DB) error {
			// 事务内的具体操作：1. 初始化存储客户端 2. 删除数据库记录 3. 删除存储文件
			oss := upload.NewOssWithStorage(image.Storage)
			if err := global.DB.Delete(&image).Error; err != nil {
				return err
			}
			return oss.DeleteImage(image.Name)
		}); err != nil {
			return err
		}
	}
	return nil
}
func (imageService *ImageService) ImageList(info request.ImageList) (interface{}, int64, error) {
	db := global.DB
	if info.Name != nil {
		db = db.Where("name LIKE ?", "%"+*info.Name+"%")
	}
	if info.Category != nil {
		category := appTypes.ToCategory(*info.Category)
		db = db.Where("category = ?", category)
	}

	if info.Storage != nil {
		storage := appTypes.ToStorage(*info.Storage)
		db = db.Where("storage = ?", storage)
	}
	option := other.MySQLOption{
		PageInfo: info.PageInfo,
		Where:    db,
	}
	return utils.MySQLPagination(&database.Image{}, option)
}
