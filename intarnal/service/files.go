package service

import "github.com/magmaheat/cache-service/intarnal/repo"

type Files interface {
}

type FilesService struct {
	filesRepo repo.Files
}

func NewFilesService(filesRepo repo.Files) *FilesService {
	return &FilesService{filesRepo}
}
