package shape

type DataShape struct {
	Widget struct {
		Debug  string `json:"debug"`
		Window struct {
			Title  string `json:"title"`
			Name   string `json:"name"`
			Width  int    `json:"width"`
			Height int    `json:"height"`
		} `json:"window"`
		Image struct {
			Src       string `json:"src"`
			Name      string `json:"name"`
			HOffset   int    `json:"hOffset"`
			VOffset   int    `json:"vOffset"`
			Alignment string `json:"alignment"`
		} `json:"image"`
		Text struct {
			Data      string `json:"data"`
			Size      int    `json:"size"`
			Style     string `json:"style"`
			Name      string `json:"name"`
			HOffset   int    `json:"hOffset"`
			VOffset   int    `json:"vOffset"`
			Alignment string `json:"alignment"`
			OnMouseUp string `json:"onMouseUp"`
		} `json:"text"`
	} `json:"widget"`
}
