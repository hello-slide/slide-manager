package slide

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
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

	dateOp := newDateOp()

	slideContent := SlideContent{
		Title:      title,
		Id:         slideId,
		CreateDate: dateOp.getDateJST(),
		ChangeDate: dateOp.getDateJST(),
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

	dateOp := newDateOp()
	slideDetails.ChangeDate = dateOp.getDateJST()

	body, err := json.Marshal(slideDetails)
	if err != nil {
		return nil, err
	}

	slideInfo := state.NewState(s.client, &s.ctx, slideInfoState)
	if err := slideInfo.Set(slideId, body); err != nil {
		return nil, err
	}

	if err := s.changedDateUpdate(true, false, slideId); err != nil {
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
	if err := s.changedDateUpdate(true, true, slideId); err != nil {
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
	if err != nil {
		return nil, err
	}
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

	return []byte(""), nil
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

	if utf8.RuneCount(getData.Value) == 0 {
		return fmt.Errorf("the slide does not exist")
	}
	if err := json.Unmarshal(getData.Value, &slideConfig); err != nil {
		return err
	}

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

	// change slide details.
	_slideData, err := slideInfo.Get(slideId)
	if err != nil {
		return err
	}

	if utf8.RuneCount(_slideData.Value) == 0 {
		return fmt.Errorf("the slide does not exist")
	}

	var slideData SlideData

	if err := json.Unmarshal(_slideData.Value, &slideData); err != nil {
		return err
	}
	slideData.Title = newName

	body, err = json.Marshal(slideData)
	if err != nil {
		return err
	}

	if err := slideInfo.Set(slideId, body); err != nil {
		return err
	}

	return nil
}

// Swap pages
//
// Arguments:
// - slideId: Id of slide,
// - origin: origin index.
// - target: target index.
func (s *SlideManager) SwapPage(slideId string, origin int, target int) error {
	slideData, err := s.GetSlideDetails(slideId)
	if err != nil {
		return err
	}

	if origin >= len(slideData.Pages) || target >= len(slideData.Pages) || origin < 0 || target < 0 {
		return fmt.Errorf("the specified index is out of range")
	}

	buffer := slideData.Pages[origin]
	slideData.Pages[origin] = slideData.Pages[target]
	slideData.Pages[target] = buffer

	dateOp := newDateOp()
	slideData.ChangeDate = dateOp.getDateJST()

	body, err := json.Marshal(slideData)
	if err != nil {
		return err
	}
	slideInfo := state.NewState(s.client, &s.ctx, slideInfoState)
	if err := slideInfo.Set(slideId, []byte(body)); err != nil {
		return err
	}

	if err := s.changedDateUpdate(true, false, slideId); err != nil {
		return nil
	}

	return nil
}

// Delete slide.
//
// Arguments:
// - slideId: Id of slide.
// - storageOp: storage op instance
func (s *SlideManager) Delete(slideId string, storageOp storage.StorageOp) error {
	slideInfo := state.NewState(s.client, &s.ctx, slideInfoState)

	// delete slide config
	getData, err := slideInfo.Get(s.userId)
	if err != nil {
		return err
	}

	var slideConfig SlideConfig

	if utf8.RuneCount(getData.Value) == 0 {
		return fmt.Errorf("the slide does not exist")
	}

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

	// delete slide page info
	if err := slideInfo.Delete(slideId); err != nil {
		return err
	}

	// Delete page data.
	filePath := []string{
		"pages",
		s.userId,
		slideId,
	}
	if err := storageOp.Delete(strings.Join(filePath, "/")); err != nil {
		return err
	}

	return nil

}

// Delete All slide.
//
// Arguments:
// - storageOp: storage op instance
func (s *SlideManager) DeleteAll(storageOp storage.StorageOp) error {
	slideInfo := state.NewState(s.client, &s.ctx, slideInfoState)

	slideData, err := slideInfo.Get(s.userId)
	var slideConfig SlideConfig

	if utf8.RuneCount(slideData.Value) == 0 {
		// If the database with userId as Key does not exist, the slide data is empty.
		return nil
	}

	if err := json.Unmarshal(slideData.Value, &slideConfig); err != nil {
		return err
	}

	if err != nil {
		return err
	}
	for _, pageId := range slideConfig.Slides {
		if err := slideInfo.Delete(pageId.Id); err != nil {
			return err
		}
	}

	if err := slideInfo.Delete(s.userId); err != nil {
		return err
	}

	// Delete page data.
	filePath := []string{
		"pages",
		s.userId,
	}
	if err := storageOp.Delete(strings.Join(filePath, "/")); err != nil {
		return err
	}
	return nil
}

// Delete page.
//
// Arguments:
// - slideId: Id of slide.
// - pageId: Id of page.
// - storageOp: storage op instance
func (s *SlideManager) DeletePage(slideId string, pageId string, storageOp storage.StorageOp) error {
	slideInfo := state.NewState(s.client, &s.ctx, slideInfoState)

	getData, err := slideInfo.Get(slideId)
	if err != nil {
		return err
	}

	var slideData SlideData

	if utf8.RuneCount(getData.Value) == 0 {
		return fmt.Errorf("the slide does not exist")
	}

	if err := json.Unmarshal(getData.Value, &slideData); err != nil {
		return err
	}
	slideData.NumberOfPages--
	dateOp := newDateOp()
	slideData.ChangeDate = dateOp.getDateJST()

	deleteIndex, err := getIndexPage(slideData, pageId)
	if err != nil {
		return err
	}
	newPages := removePage(slideData.Pages, deleteIndex)
	slideData.Pages = newPages

	body, err := json.Marshal(slideData)
	if err != nil {
		return err
	}

	if err := slideInfo.Set(slideId, body); err != nil {
		return err
	}

	if err := s.changedDateUpdate(true, false, slideId); err != nil {
		return err
	}

	filePath := []string{
		"pages",
		s.userId,
		slideId,
		pageId,
	}
	if err := storageOp.Delete(strings.Join(filePath, "/")); err != nil {
		return err
	}
	return nil
}

// Update `change_date`
//
// Arguments:
// - isInfo: If set to true, the user's slide database will be updated.
// - isDetails: If set to true, the detailed database for that slide will be updated.
// - slideId: Id of slide.
func (s *SlideManager) changedDateUpdate(isInfo bool, isDetails bool, slideId string) error {
	dateOp := newDateOp()
	if isInfo {
		slideInfo, err := s.GetInfo()
		if err != nil {
			return err
		}
		targetIndex, err := getIndexSlideConfig(*slideInfo, slideId)
		if err != nil {
			return err
		}
		slideInfo.Slides[targetIndex].ChangeDate = dateOp.getDateJST()

		body, err := json.Marshal(slideInfo)
		if err != nil {
			return err
		}

		_slideInfo := state.NewState(s.client, &s.ctx, slideInfoState)
		if err := _slideInfo.Set(s.userId, body); err != nil {
			return err
		}
	}

	if isDetails {
		slideDetails, err := s.GetSlideDetails(slideId)
		if err != nil {
			return err
		}
		slideDetails.ChangeDate = dateOp.getDateJST()

		body, err := json.Marshal(slideDetails)
		if err != nil {
			return err
		}

		slideInfo := state.NewState(s.client, &s.ctx, slideInfoState)
		if err := slideInfo.Set(slideId, body); err != nil {
			return err
		}
	}
	return nil
}
