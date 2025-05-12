package main

import (
	"fmt"

	"github.com/xiaowumin-mark/QQMusicDecoder-Go"
)

func main() {
	// 使用MID获取歌词
	fmt.Println(QQMusicDecoder.GetLyricsByMid("0016hWfr285oV5"))
	// 通过ID获取歌词
	fmt.Println(QQMusicDecoder.GetLyrics("732381"))
}
