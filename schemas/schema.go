package schemas

type CallInfo struct {
	// Uid     int    `json:"uid"`
	Channel   string     `json:"channel"`
	Streamers []Streamer `json:"streamers"`
}

type StopRecording struct {
	Uid     int    `json:"uid"`
	Channel string `json:"channel"`
	Rid     string `json:"rid"`
	Sid     string `json:"sid"`
}

type UserCredentials struct {
	Rtc string `json:"rtc"`
	UID int    `json:"uid"`
}

type QueryRecording struct {
	Rid string `json:"rid"`
	Sid string `json:"sid"`
}

type UpdateRecording struct {
	Uid       int        `json:"uid"`
	Channel   string     `json:"channel"`
	Rid       string     `json:"rid"`
	Sid       string     `json:"sid"`
	Streamers []Streamer `json:"streamers"`
}

type Streamer struct {
	Uid      string `json:"uid"`
	ImageURL string `json:"image_url"`
}
