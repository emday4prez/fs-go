package service

import "log"

type FileService struct {
}

func NewFileService() *FileService {
	log.Println("new FileService created")
	return &FileService{}
}
