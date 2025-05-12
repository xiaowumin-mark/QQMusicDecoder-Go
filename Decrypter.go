package QQMusicDecoder

import (
	"bytes"
	"compress/flate"
	"encoding/hex"
	"io"

	"github.com/fred913/goqrcdec"
	//"github.com/fred913/goqrcdec"
)

var (
	QQKey = []byte("!@#)(*$%123ZXC!@!@#)(NHL")
)

func DecryptLyrics(encrypted []byte) (string, error) {
	//encryptedTextByte := HexStringToByteArray(encrypted)
	//data := make([]byte, len(encryptedTextByte))
	//schedule := make([][][]byte, 3)
	//for i := 0; i < 3; i++ {
	//	schedule[i] = make([][]byte, 16)
	//	for j := 0; j < 16; j++ {
	//		schedule[i][j] = make([]byte, 6)
	//	}
	//}
	//TripleDESKeySetup(QQKey, schedule, DECRYPT)
	//for i := 0; i < len(encryptedTextByte); i += 8 {
	//	temp := make([]byte, 8)
	//	TripleDESCrypt(encryptedTextByte[i:], temp, schedule)
	//	for j := 0; j < 8; j++ {
	//		data[i+j] = temp[j]
	//	}
	//}
	//
	//var unzip = SharpZipLibDecompress(data)
	////var result = Encoding.UTF8.GetString(unzip);
	//return string(unzip), nil
	res, err := goqrcdec.DecodeQRC(encrypted)
	if err != nil {
		return "", err
	}
	return string(res), nil
}

func HexStringToByteArray(hexString string) []byte {
	length := len(hexString)
	byts := make([]byte, length/2)
	for i := 0; i < length; i += 2 {
		b, _ := hex.DecodeString(hexString[i : i+2])
		byts[i/2] = b[0]
	}
	return byts
}

func SharpZipLibDecompress(data []byte) []byte {
	compressed := bytes.NewReader(data)
	decompressed := new(bytes.Buffer)

	// 创建 flate.Reader（相当于 C# 的 InflaterInputStream）
	reader := flate.NewReader(compressed)
	defer reader.Close()

	// 拷贝解压后的数据到 buffer
	_, err := io.Copy(decompressed, reader)
	if err != nil {
		return nil
	}

	return decompressed.Bytes()
}
