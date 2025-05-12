package QQMusicDecoder

type QqLyricsResponse struct {
	Lyrics string `json:"lyrics"`
	Trans  string `json:"trans"`
}

type SongResponse struct {
	Code int    `json:"code"`
	Data []Song `json:"data"`
}

type Song struct {
	Id int `json:"id"`
}
