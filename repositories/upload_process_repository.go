package repositories

import (
	"jellyfin_uploader/models"
	"jellyfin_uploader/util"
)

func ListUploadProcesses() ([]models.UploadProcess, error) {
	var uploadProcesses []models.UploadProcess
	err := util.DB.Model(&models.UploadProcess{}).Preload("Files").Find(&uploadProcesses).Error
	if err != nil {
		return nil, err
	}
	return uploadProcesses, nil
}

func CreateUploadProcess(uploadProcess *models.UploadProcess) error {
	return util.DB.Create(uploadProcess).Error
}

func DeleteUploadProcess(id string) error {
	return util.DB.Delete(&models.UploadProcess{}, id).Error
}
