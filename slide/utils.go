package slide

import "fmt"

// Pop element from list.
func removeSlides(s []SlideContent, i int) []SlideContent {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

// Pop element from list.
func removePage(s []PageData, i int) []PageData {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

// Returns the index of the corresponding slide ID.
//
// Arguments:
// - slideConfig: SlideConfig
// - targetId: target slide id.
//
// Returns:
// - int: Index of the corresponding targetId.
func getIndexSlideConfig(slideConfig SlideConfig, targetId string) (int, error) {
	var targetIndex int
	var isExist bool = false
	for index, data := range slideConfig.Slides {
		if data.Id == targetId {
			targetIndex = index
			isExist = true
			break
		}
	}
	if !isExist {
		return 0, fmt.Errorf("The specified slide ID does not exist.")
	}
	return targetIndex, nil
}

// Returns the index of the corresponding page ID.
//
// Arguments:
// - slideData: SlideData
// - targetId: target slide id.
//
// Returns:
// - int: Index of the corresponding targetId.
func getIndexPage(slideData SlideData, targetId string) (int, error) {
	var targetIndex int
	var isExist bool = false
	for index, data := range slideData.Pages {
		if data.PageId == targetId {
			targetIndex = index
			isExist = true
			break
		}
	}
	if !isExist {
		return 0, fmt.Errorf("The specified slide ID does not exist.")
	}
	return targetIndex, nil
}
