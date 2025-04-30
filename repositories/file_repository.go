package repositories

import (
	"jellyfin_uploader/models"
	"jellyfin_uploader/util"
)

func CreateFile(file *models.File) error {
	return util.DB.Create(file).Error
}
func UpdateFile(file *models.File) error {
	return util.DB.Save(file).Error
}
