package service

import (
	"errors"
	"fmt"
	"mime/multipart"
	"os"

	"mhelper_server/global"

	"mhelper_server/pkg/upload"
)

type FileInfo struct {
	Name      string
	AccessUrl string
}

func (svc *Service) UploadFile(fileType upload.FileType, file multipart.File, fileHeader *multipart.FileHeader) (*FileInfo, error) {
	fileName := upload.GetFileName(fileHeader.Filename)

	// 校验文件后缀名
	if !upload.CheckContainExt(fileType, fileName) {
		return nil, errors.New("file suffix is not supported")
	}
	// 校验文件大小
	if upload.CheckMaxSize(fileType, file) {
		return nil, errors.New("exceeded maximum file limit")
	}

	uploadSavePath := upload.GetSavePath()
	// 检查存储文件目录是否穿件
	if upload.CheckSavePath(uploadSavePath) {
		// 未创建情形 创建之
		if err := upload.CreateSavePath(uploadSavePath, os.ModePerm); err != nil {
			return nil, errors.New("failed to create save directory")
		}
	}
	// 检查存储文件目录权限
	if upload.CheckPermission(uploadSavePath) {
		return nil, errors.New("insufficient file permissions")
	}

	dst := uploadSavePath + "/" + fileName
	if err := upload.SaveFile(fileHeader, dst); err != nil {
		return nil, err
	}

	accessUrl := global.AppSetting.UploadServerUrl + "/" + fileName
	accessUrl = fmt.Sprintf(accessUrl, global.ServerSetting.HttpPort)
	return &FileInfo{Name: fileName, AccessUrl: accessUrl}, nil
}
