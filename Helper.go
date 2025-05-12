package QQMusicDecoder

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/beevik/etree"
	"github.com/valyala/fastjson"
)

var (
	UserAgent = "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/63.0.3239.132 Safari/537.36"
	Cookie    = "os=pc;osver=Microsoft-Windows-10-Professional-build-16299.125-64bit;appver=2.0.3.131777;channel=netease;__remember_me=true"

	VerbatimXmlMappingDict = map[string]string{
		"content":     "orig",  // 原文
		"contentts":   "ts",    // 译文
		"contentroma": "roma",  // 罗马音
		"Lyric_1":     "lyric", // 解压后的内容
	}
)

func GetLyricsByMid(mid string) (*QqLyricsResponse, error) {
	song := GetSong(mid)
	if song == nil || song.Data == nil || len(song.Data) == 0 {
		return nil, fmt.Errorf("歌曲不存在")
	}
	return GetLyrics(fmt.Sprintf("%d", song.Data[0].Id))
}

func GetSong(mid string) *SongResponse {
	callBack := "getOneSongInfoCallback"

	body := map[string]string{
		"songmid":       mid,
		"tpl":           "yqq_song_detail",
		"format":        "jsonp",
		"callback":      callBack,
		"g_tk":          "5381",
		"jsonpCallback": callBack,
		"loginUin":      "0",
		"hostUin":       "0",
		"outCharset":    "utf8",
		"notice":        "0",
		"platform":      "yqq",
		"needNewCode":   "0",
	}

	//https://c.y.qq.com/v8/fcg-bin/fcg_play_single_song.fcg
	/*var json = await PostAsync("https://c.y.qq.com/v8/fcg-bin/fcg_play_single_song.fcg", data);

	try
	{
	    return JsonSerializer.Deserialize<SongResponse>(json);
	}
	catch
	{
	    return null;
	}*/

	str, err := Post("https://c.y.qq.com/v8/fcg-bin/fcg_play_single_song.fcg", body)
	if err != nil {
		return nil
	}

	//去除 getOneSongInfoCallback( 字符串
	str = str[len(callBack)+1 : len(str)-1]
	// 去除最后的 ) 字符串

	var json = fastjson.MustParse(str)

	var response SongResponse

	response.Code = json.GetInt("code")

	arr := json.GetArray("data")
	if len(arr) == 0 {
		return &response
	}
	var song Song
	song.Id = arr[0].GetInt("id")
	response.Data = append(response.Data, song)
	return &response
}

func Post(url_ string, param map[string]string) (string, error) {
	// 1. 编码表单数据
	formData := url.Values{}
	for k, v := range param {
		formData.Add(k, v)
	}

	// 2. 创建请求对象
	req, err := http.NewRequest("POST", url_, strings.NewReader(formData.Encode()))
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %w", err)
	}

	// 3. 设置请求头
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Referer", "https://c.y.qq.com/")
	req.Header.Set("User-Agent", UserAgent)
	req.Header.Set("Cookie", Cookie)

	// 5. 使用自定义客户端发送
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	// 6. 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应失败: %w", err)
	}

	return string(body), nil
}

func GetLyrics(id string) (*QqLyricsResponse, error) {
	resp, err := Post("https://c.y.qq.com/qqmusic/fcgi-bin/lyric_download.fcg", map[string]string{
		"version":     "15",
		"miniversion": "82",
		"lrctype":     "4",
		"musicid":     id,
	})
	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}
	//resp = resp.Replace("<!--", "").Replace("-->", "");
	// 去除响应中的注释
	resp = strings.Replace(resp, "<!--", "", -1)
	resp = strings.Replace(resp, "-->", "", -1)

	dict := map[string]*etree.Element{}

	RecursionFindElement(&CreateXmlDom(resp).Element, VerbatimXmlMappingDict, dict)
	resout := &QqLyricsResponse{
		Lyrics: "",
		Trans:  "",
	}
	for key, pair := range dict {
		//log.Println(key, pair.Text())
		text := pair.Text()
		if text == "" {
			continue
		}
		var decompressText string

		if decompressText, err = DecryptLyrics(HexStringToByteArray(text)); err != nil {
			continue
		}

		s := ""
		if strings.Contains(decompressText, "<?xml") {
			doc := CreateXmlDom(decompressText)
			subDict := map[string]*etree.Element{}
			RecursionFindElement(&doc.Element, VerbatimXmlMappingDict, subDict)
			if subDict["lyric"] != nil {
				//s = subDict["lyric"]
				// 获取subDict["lyric"] 的LyricContent=属性
				s = subDict["lyric"].SelectAttrValue("LyricContent", "")
			}
		} else {
			s = decompressText
		}

		if s != "" {
			switch key {
			case "orig":
				resout.Lyrics = s // LyricUtils.DealVerbatimLyric(s, SearchSourceEnum.QQ_MUSIC)
				break
			case "ts":
				resout.Trans = s // LyricUtils.DealVerbatimLyric(s, SearchSourceEnum.QQ_MUSIC)
				break
			}
		}
	}
	if resout.Lyrics == "" && resout.Trans == "" {
		return nil, fmt.Errorf("歌词解析失败")
	}
	return resout, nil
}
