# QQMusicDecoder-go

[https://github.com/WXRIW/QQMusicDecoder](https://github.com/WXRIW/QQMusicDecoder) 的golang实现

可通过 Id 直接通过网络获取歌词 `Helper.GetLyrics(string, string)`；  
可通过 Mid 直接通过网络获取歌词 `Helper.GetLyricsByMid(string, string)`；  
也可以直接解密 QRC 歌词 `Decrypter.DecryptLyrics(string)`。  

### 特别感谢
[WXRIW/QQMusicDecoder](https://github.com/WXRIW/QQMusicDecoder)
[fred913/goqrcdec](https://github.com/fred913/goqrcdec)