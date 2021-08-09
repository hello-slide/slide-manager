package slide

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
	"unicode/utf8"

	"github.com/dapr/go-sdk/client"
	"github.com/hello-slide/slide-manager/state"
	"github.com/hello-slide/slide-manager/storage"
	"github.com/hello-slide/slide-manager/token"
)

type SlideManager struct {
	ctx    context.Context
	userId string
	client *client.Client
}

func NewSlideManager(ctx context.Context, daprClient *client.Client, userId string) *SlideManager {

	return &SlideManager{
		ctx:    ctx,
		userId: userId,
		client: daprClient,
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
	slideId, err := token.CreateId(title)
	if err != nil {
		return "", err
	}
	now := time.Now()
	nowUTC := now.UTC()
	jst := time.FixedZone("Asia/Tokyo", 9*60*60)
	nowJST := nowUTC.In(jst)

	slideContent := SlideContent{
		Title:      title,
		Id:         slideId,
		CreateDate: nowJST.Format("20060102150405"),
		ChangeDate: nowJST.Format("20060102150405"),
	}

	slideConfig, err := s.GetInfo()
	if err != nil {
		return "", err
	}

	slideConfig.NumberOfSlides++
	slideConfig.Slides = append(slideConfig.Slides, slideContent)

	body, err := json.Marshal(slideConfig)
	if err != nil {
		return "", err
	}

	slideInfo := state.NewState(s.client, &s.ctx, slideInfoState)
	if err := slideInfo.Set(s.userId, body); err != nil {
		return "", err
	}

	return slideId, nil
}

// Create Page
//
// Arguments:
// - slideId: Id of slide.
// - pageType: page type.
func (s *SlideManager) CreatePage(slideId string, pageType string) (*PageData, error) {
	slideDetails, err := s.GetSlideDetails(slideId)
	if err != nil {
		return nil, err
	}

	pageId, err := token.CreateId(slideId)
	if err != nil {
		return nil, err
	}

	pageDate := &PageData{
		PageId: pageId,
		Type:   pageType,
	}

	slideDetails.NumberOfPages++
	slideDetails.Pages = append(slideDetails.Pages, *pageDate)

	body, err := json.Marshal(slideDetails)
	if err != nil {
		return nil, err
	}

	slideInfo := state.NewState(s.client, &s.ctx, slideInfoState)
	if err := slideInfo.Set(slideId, body); err != nil {
		return nil, err
	}

	return pageDate, nil
}

// Write page data.
//
// Arguments:
// - data: page data.
// - slideId: Id of slide.
// - pageId: Id of page.
// - storageOp: storage op instance
func (s *SlideManager) SetPage(data []byte, slideId string, pageId string, storageOp storage.StorageOp) error {
	dirs := []string{
		"pages",
		s.userId,
		slideId,
	}
	if err := storageOp.WriteFile(dirs, pageId, data); err != nil {
		return err
	}
	return nil
}

// Get Slides infomation of user.
func (s *SlideManager) GetInfo() (*SlideConfig, error) {

	slideInfo := state.NewState(s.client, &s.ctx, slideInfoState)
	getData, err := slideInfo.Get(s.userId)
	if err != nil {
		return nil, err
	}

	if utf8.RuneCount(getData.Value) != 0 {
		var slideConfig SlideConfig

		if err := json.Unmarshal(getData.Value, &slideConfig); err != nil {
			return nil, err
		}
		return &slideConfig, nil
	}
	// Not exist
	return &SlideConfig{
		NumberOfSlides: 0,
		Slides:         []SlideContent{},
	}, nil
}

// Get slide detail data.
//
// Arguments:
// - slideId: Id of slide.
func (s *SlideManager) GetSlideDetails(slideId string) (*SlideData, error) {
	slideInfo := state.NewState(s.client, &s.ctx, slideInfoState)
	getSlideData, err := slideInfo.Get(slideId)
	if err != nil {
		return nil, err
	}

	if utf8.RuneCount(getSlideData.Value) != 0 {
		var slideData SlideData

		if err := json.Unmarshal(getSlideData.Value, &slideData); err != nil {
			return nil, err
		}
		return &slideData, nil
	}
	// Not exist
	// Create Slide Data

	slidesConfig, err := s.GetInfo()
	targetIndex, err := getIndexSlideConfig(*slidesConfig, slideId)
	if err != nil {
		return nil, err
	}
	slideConfig := slidesConfig.Slides[targetIndex]

	newSlideInfo := &SlideData{
		NumberOfPages: 0,
		Pages:         []PageData{},
		SlideContent:  slideConfig,
	}
	body, err := json.Marshal(newSlideInfo)
	if err != nil {
		return nil, err
	}

	if err := slideInfo.Set(slideId, body); err != nil {
		return nil, err
	}

	return newSlideInfo, nil
}

// Get page data.
//
// Arguments:
// - slideId: Id of slide.
// - pageId: Id of page.
// - storageOp: storage op instance
func (s *SlideManager) GetPage(slideId string, pageId string, storageOp storage.StorageOp) ([]byte, error) {
	dirs := []string{
		"pages",
		s.userId,
		slideId,
	}
	isExist, err := storageOp.FileExist(dirs, pageId)
	if err != nil {
		return nil, err
	}

	if isExist {
		return storageOp.ReadFile(dirs, pageId)
	}
	return nil, fmt.Errorf("Page data is not exist.")
}

// Rename slide
//
// Arguments:
// - slideId: slide id.
// - newName: new name(title)
func (s *SlideManager) Rename(slideId string, newName string) error {
	slideInfo := state.NewState(s.client, &s.ctx, slideInfoState)
	getData, err := slideInfo.Get(s.userId)
	if err != nil {
		return err
	}

	var slideConfig SlideConfig

	if utf8.RuneCount(getData.Value) != 0 {
		if err := json.Unmarshal(getData.Value, &slideConfig); err != nil {
			return err
		}
		slideConfig.NumberOfSlides--

		targetIndex, err := getIndexSlideConfig(slideConfig, slideId)
		if err != nil {
			return err
		}
		slideConfig.Slides[targetIndex].Title = newName

		body, err := json.Marshal(slideConfig)
		if err != nil {
			return err
		}

		if err := slideInfo.Set(s.userId, body); err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("The slide does not exist.")
}

// Delete slide.
//
// Arguments:
// - slideId: Id of slide.
func (s *SlideManager) Delete(slideId string) error {
	slideInfo := state.NewState(s.client, &s.ctx, slideInfoState)
	getData, err := slideInfo.Get(s.userId)
	if err != nil {
		return err
	}

	var slideConfig SlideConfig

	if utf8.RuneCount(getData.Value) != 0 {
		if err := json.Unmarshal(getData.Value, &slideConfig); err != nil {
			return err
		}
		slideConfig.NumberOfSlides--

		deleteIndex, err := getIndexSlideConfig(slideConfig, slideId)
		if err != nil {
			return err
		}
		newSlides := removeSlides(slideConfig.Slides, deleteIndex)
		slideConfig.Slides = newSlides

		body, err := json.Marshal(slideConfig)
		if err != nil {
			return err
		}

		if err := slideInfo.Set(s.userId, body); err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("The slide does not exist.")

}

// Delete All slide.
func (s *SlideManager) DeleteAll() error {
	slideInfo := state.NewState(s.client, &s.ctx, slideInfoState)
	if err := slideInfo.Delete(s.userId); err != nil {
		return err
	}
	return nil
}
