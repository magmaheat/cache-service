package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/magmaheat/cache-service/internal/entity"
	"github.com/magmaheat/cache-service/internal/repo"
	log "github.com/sirupsen/logrus"
	"io"
	"mime/multipart"
	"slices"
	"sort"
	"time"
)

const defaultCountFiles = 10

type Cache interface {
	SaveData(ctx context.Context, meta entity.Meta, jsonField string, file *multipart.FileHeader) error
	GetDocument(ctx context.Context, id string) (*entity.Document, error)
	GetDocuments(ctx context.Context, input entity.SearchDocuments) (entity.MetaSlice, error)
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

	log.Debugf("save meta document id: %s", meta.Id)

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

	log.Debugf("save file document id: %s", meta.Id)

	if jsonField != "" {
		err = c.cacheRepo.SaveJson(ctx, key, "json", jsonField)
		if err != nil {
			return err
		}

		log.Debugf("save json document id: %s", meta.Id)

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

	log.Debugf("counts field in result: %d", len(allFields))

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

func (c *CacheService) GetDocuments(ctx context.Context, input entity.SearchDocuments) (entity.MetaSlice, error) {
	metaList, err := c.cacheRepo.GetDocuments(ctx)
	if err != nil {
		return nil, err
	}

	log.Debugf("count meta files for filter search: %d", len(metaList))

	var result entity.MetaSlice

	for _, meta := range metaList {
		switch {
		case !slices.Contains(meta.Grant, input.Login):
			continue
		case input.File != nil && *input.File != meta.File:
			continue
		case input.Public != nil && *input.Public != meta.Public:
			continue
		case input.Mime != "" && input.Mime != meta.Mime:
			continue
		case input.Name != "" && input.Name != meta.Name:
			continue
		default:
			result = append(result, meta)
		}
	}

	log.Debugf("found example: %d", len(result))

	if len(result) == 0 {
		return nil, ErrFileNotFound
	}

	sort.Sort(result)

	if input.Limit != 0 && input.Limit <= len(result) {
		return result[:input.Limit], nil
	}

	return result[:defaultCountFiles], nil
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
