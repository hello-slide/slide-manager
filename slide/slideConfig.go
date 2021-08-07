package slide

// Detailed information for each slide.
type SlideData struct {
	NumberOfPages int      `json:"number_of_pages"`
	Indexes       []string `json:"indexes"`
	SlideContent
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
