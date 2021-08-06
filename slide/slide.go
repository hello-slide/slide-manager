package slide

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"cloud.google.com/go/storage"
	storageOp "github.com/hello-slide/slide-manager/storage"
	"github.com/hello-slide/slide-manager/token"
)

type SlideManager struct {
	ctx       context.Context
	storageOp *storageOp.StorageOp
	userId    string
}

func NewSlideManager(ctx context.Context, client storage.Client, bucketName string, userId string) *SlideManager {
	storageOp := storageOp.NewStorageOp(ctx, client, bucketName)

	return &SlideManager{
		ctx:       ctx,
		storageOp: storageOp,
		userId:    userId,
	}
}

// Create slide
//
// Arguments:
// - title: Slide title.
//
// Return:
// - id string: Slide id
func (s *SlideManager) Create(title string) (string, error) {
	slideInfoPath := []string{s.userId}
	fileName := "slide_config.json"

	slideId, err := token.CreateSlideId(title)
	if err != nil {
		return "", err
	}
	t := time.Now()

	slideContent := SlideContent{
		Title:      title,
		Id:         slideId,
		CreateDate: t.Format("20060102150405"),
		ChangeDate: t.Format("20060102150405"),
	}
	var slideConfig SlideConfig

	isExist, err := s.storageOp.FileExist(slideInfoPath, fileName)
	if err != nil {
		return "", err
	}

	if isExist {
		slideConfigByte, err := s.storageOp.ReadFile(slideInfoPath, fileName)
		if err != nil {
			return "", err
		}
		if err := json.Unmarshal(slideConfigByte, &slideConfig); err != nil {
			return "", err
		}
		slideConfig.NumberOfSlides++
		slideConfig.Slides = append(slideConfig.Slides, slideContent)
	} else {
		slideConfig = SlideConfig{
			UserId:         s.userId,
			NumberOfSlides: 1,
			Slides: []SlideContent{
				slideContent,
			},
		}
	}

	body, err := json.Marshal(slideConfig)
	if err != nil {
		return "", err
	}

	if err := s.storageOp.WriteFile(slideInfoPath, fileName, body); err != nil {
		return "", err
	}

	return slideId, nil
}

// Get Slides infomation of user.
func (s *SlideManager) GetInfo() (*SlideConfig, error) {
	slideInfoPath := []string{s.userId}
	fileName := "slide_config.json"

	isExist, err := s.storageOp.FileExist(slideInfoPath, fileName)
	if err != nil {
		return nil, err
	}
	if isExist {
		var slideConfig SlideConfig
		slideConfigByte, err := s.storageOp.ReadFile(slideInfoPath, fileName)
		if err != nil {
			return nil, err
		}

		if err := json.Unmarshal(slideConfigByte, &slideConfig); err != nil {
			return nil, err
		}
		return &slideConfig, nil
	}
	// Not exist
	return &SlideConfig{
		UserId:         s.userId,
		NumberOfSlides: 0,
		Slides:         nil,
	}, nil
}

// Delete slide.
//
// Arguments:
// - slideId: Id of slide.
func (s *SlideManager) Delete(slideId string) error {
	slideInfoPath := []string{s.userId}
	fileName := "slide_config.json"

	var slideConfig SlideConfig

	isExist, err := s.storageOp.FileExist(slideInfoPath, fileName)
	if err != nil {
		return err
	}

	if isExist {
		slideConfigByte, err := s.storageOp.ReadFile(slideInfoPath, fileName)
		if err != nil {
			return err
		}
		if err := json.Unmarshal(slideConfigByte, &slideConfig); err != nil {
			return err
		}
		slideConfig.NumberOfSlides--

		var deleteIndex int
		for index, data := range slideConfig.Slides {
			if data.Id == slideId {
				deleteIndex = index
				break
			}
		}
		newSlides := removeSlides(slideConfig.Slides, deleteIndex)
		slideConfig.Slides = newSlides

		body, err := json.Marshal(slideConfig)
		if err != nil {
			return err
		}

		if err := s.storageOp.WriteFile(slideInfoPath, fileName, body); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("The slide does not exist.")
	}
	return nil
}

// Pop element from list.
func removeSlides(s []SlideContent, i int) []SlideContent {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}
