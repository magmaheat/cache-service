package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
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
	GetDocuments(ctx context.Context, meta entity.Meta, limit int) ([]entity.Meta, error)
	DeleteDocument(ctx context.Context, id string) error
}

type CacheService struct {
	cacheRepo repo.Cache
}

func NewCacheService(cacheRepo repo.Cache) *CacheService {
	return &CacheService{cacheRepo}
}

func (c *CacheService) SaveData(ctx context.Context, meta entity.Meta, jsonField string, file *multipart.FileHeader) error {
	const fn = "service - cache - SaveFile"

	meta.Id = uuid.New().String()

	log.Infof("save document id: %s", meta.Id)

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

	key := fmt.Sprintf("document:%s", meta.Id)

	if err = c.cacheRepo.SaveJson(ctx, key, "meta", string(jsonMeta)); err != nil {
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
		if err = c.cacheRepo.SaveFile(ctx, key, "file", fileData); err != nil {
			return err
		}
	}

	if jsonField != "" {
		err = c.cacheRepo.SaveJson(ctx, key, "json", jsonField)
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

func (c *CacheService) GetDocuments(ctx context.Context, meta entity.Meta, limit int) ([]entity.Meta, error) {
	metaList, err := c.cacheRepo.GetDocuments(ctx)
	if err != nil {
		return nil, err
	}

	_ = metaList

	return nil, nil
}

func (c *CacheService) DeleteDocument(ctx context.Context, id string) error {
	key := fmt.Sprintf("document:%s", id)

	count, err := c.cacheRepo.DeleteDocument(ctx, key)
	if err != nil {
		return err
	}

	if count == 0 {
		return ErrFileNotFound
	}

	return nil
}
