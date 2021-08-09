package slide

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
	"unicode/utf8"

	"github.com/dapr/go-sdk/client"
	"github.com/hello-slide/slide-manager/state"
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
	slideId, err := token.CreateSlideId(title)
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
	slideInfo := state.NewState(s.client, &s.ctx, slideInfoState)
	getData, err := slideInfo.Get(s.userId)
	if err != nil {
		return "", err
	}

	var slideConfig SlideConfig

	if utf8.RuneCount(getData.Value) != 0 {
		if err := json.Unmarshal(getData.Value, &slideConfig); err != nil {
			return "", err
		}
		slideConfig.NumberOfSlides++
		slideConfig.Slides = append(slideConfig.Slides, slideContent)
	} else {
		slideConfig = SlideConfig{
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

	if err := slideInfo.Set(s.userId, body); err != nil {
		return "", err
	}

	return slideId, nil
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

		if err := slideInfo.Set(s.userId, body); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("The slide does not exist.")
	}
	return nil
}

// Delete All slide.
func (s *SlideManager) DeleteAll() error {
	slideInfo := state.NewState(s.client, &s.ctx, slideInfoState)
	if err := slideInfo.Delete(s.userId); err != nil {
		return err
	}
	return nil
}

// Pop element from list.
func removeSlides(s []SlideContent, i int) []SlideContent {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}
