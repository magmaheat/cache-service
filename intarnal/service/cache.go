package service

import (
	"context"
	"encoding/json"
	"github.com/magmaheat/cache-service/intarnal/entity"
	"github.com/magmaheat/cache-service/intarnal/repo"
	log "github.com/sirupsen/logrus"
	"io"
	"mime/multipart"
	"time"
)

type Cache interface {
	SaveData(ctx context.Context, meta entity.Meta, jsonField string, file *multipart.FileHeader) error
	GetDocument(ctx context.Context, id string) (*entity.Document, error)
}

type CacheService struct {
	cacheRepo repo.Cache
}

func NewCacheService(cacheRepo repo.Cache) *CacheService {
	return &CacheService{cacheRepo}
}

func (c *CacheService) SaveData(ctx context.Context, meta entity.Meta, jsonField string, file *multipart.FileHeader) error {
	const fn = "service - cache - SaveFile"

	exists, err := c.cacheRepo.CheckId(ctx, meta.Id)
	if err != nil {
		return err
	}

	if exists {
		return ErrFileAlreadyExists
	}

	meta.Created = time.Now()

	jsonMeta, err := json.Marshal(meta)
	if err != nil {
		log.Errorf("%s - json.Marshal: %v", fn, err)
		return err
	}

	if err = c.cacheRepo.SaveJson(ctx, meta.Id, "meta", string(jsonMeta)); err != nil {
		return err
	}

	if file != nil {
		src, err := file.Open()
		if err != nil {
			log.Errorf("%s - file.Open: %v", fn, err)
			return err
		}
		defer src.Close()

		fileData, err := io.ReadAll(src)
		if err != nil {
			log.Errorf("%s - io.ReadAll: %v", fn, err)
		}
		if err = c.cacheRepo.SaveFile(ctx, meta.Id, "file", fileData); err != nil {
			return err
		}
	}

	if jsonField != "" {
		err = c.cacheRepo.SaveJson(ctx, meta.Id, "json", jsonField)
	}

	return nil
}

func (c *CacheService) GetDocument(ctx context.Context, id string) (*entity.Document, error) {
	allFields, err := c.cacheRepo.GetDocument(ctx, id)
	if err != nil {
		return &entity.Document{}, err
	}

	if len(allFields) == 0 {
		return &entity.Document{}, ErrFileNotFound
	}

	var meta entity.Meta
	strMeta := allFields["meta"]
	err = json.Unmarshal([]byte(strMeta), &meta)
	if err != nil {
		log.Errorf("service - cache - GetDocument - Unmarshal: %v", err)
		return &entity.Document{}, err
	}

	body := []byte(allFields["file"])
	jsonFiled := allFields["json"]
	document := entity.NewDocument(body, meta.Mime, meta.Name, jsonFiled)

	return document, nil
}
