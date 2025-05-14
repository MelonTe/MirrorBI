package utils

import (
	"errors"
	"io"
	"mime/multipart"
	"os"
)

// 上传临时文件，返回文件在项目的相对路径
func SaveFileToLocal(file *multipart.FileHeader) (string, error) {
	// 1.获取文件名，校验文件大小
	if file.Size > 20*1024*1024 {
		return "", errors.New("文件大于20MB")
	}
	fileName := file.Filename
	//2.拼接目标地址
	dstPath := "./temp/" + fileName
	//3.保存文件到目标地址
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()
	dst, err := os.Create(dstPath)
	if err != nil {
		return "", err
	}
	defer dst.Close()
	_, err = io.Copy(dst, src)
	if err != nil {
		return "", err
	}
	return dstPath, nil
}

// 删除目标地址的文件
func DeleteFile(filePath string) error {
	err := os.Remove(filePath)
	if err != nil {
		return err
	}
	return nil
}
