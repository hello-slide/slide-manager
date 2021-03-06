package slide

// Detailed information for each slide.
type SlideData struct {
	NumberOfPages int        `json:"number_of_pages"`
	Pages         []PageData `json:"pages"`
	SlideContent
}

type PageData struct {
	PageId string `json:"page_id"`
	Type   string `json:"type"`
}

// Information for each slide.
type SlideContent struct {
	Title      string `json:"title"`
	Id         string `json:"id"`
	CreateDate string `json:"create_date"`
	ChangeDate string `json:"change_date"`
}

// Describe the slide information possessed by the user.
type SlideConfig struct {
	NumberOfSlides int            `json:"number_of_slides"`
	Slides         []SlideContent `json:"slides"`
}
