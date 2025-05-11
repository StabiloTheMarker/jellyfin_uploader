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

func GetFile(id string) (*models.File, error) {
	var file models.File
	err := util.DB.Model(file).First(&file, id).Error
	if err != nil {
		return nil, err
	}
	return &file, nil
}

func DeleteFile(id string) error {
	return util.DB.Delete(&models.File{}, id).Error
}
