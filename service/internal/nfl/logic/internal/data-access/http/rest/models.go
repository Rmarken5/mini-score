package rest

type Event struct {
	UID    string `json:"uid"`
	Status struct {
		Period       int    `json:"period"`
		DisplayClock string `json:"displayClock"`
		Clock        int    `json:"clock"`
		Type         struct {
			Name        string `json:"name"`
			Description string `json:"description"`
			ID          string `json:"id"`
			State       string `json:"state"`
			Completed   bool   `json:"completed"`
			Detail      string `json:"detail"`
			ShortDetail string `json:"shortDetail"`
		} `json:"type"`
	} `json:"status"`
}

type Content struct {
	SBData struct {
		Events []Event `json:"events"`
	} `json:"sbData"`
}

type ScoreboardResponse struct {
	Content Content `json:"content"`
}
