package models

type Poll struct {
	ID      int      `json:"id"`
	Topic   string   `json:"topic"`
	Options []Option `json:"options"`
}

type Option struct {
	ID    int    `json:"id"`
	Text  string `json:"text"`
	Votes int    `json:"votes"`
}

type CreatePollRequest struct {
	Topic   string   `json:"topic"`
	Options []string `json:"options"`
}
