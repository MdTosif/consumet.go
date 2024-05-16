package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/go-resty/resty/v2"
)

// Key and IV should be kept secure and unique.
var (
	keyHex = "37911490979715163134003223491201"
	key, _ = hex.DecodeString(keyHex)
	ivHex  = "3134003223491201" // IV should be unique for each encryption operation.
	iv, _  = hex.DecodeString(ivHex)
)

// Pad the data to the block size using PKCS#7 padding
func padData(data []byte) []byte {
	padding := aes.BlockSize - (len(data) % aes.BlockSize)
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padText...)
}

func encrypt(data string) (string, error) {
	// Convert string to byte slice
	dataBytes := []byte(data)

	key := []byte("37911490979715163134003223491201")
	iv := []byte("3134003223491201")

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// The IV should be the same length as the block size of the cipher
	if len(iv) != aes.BlockSize {
		return "", fmt.Errorf("IV length must be %d bytes", aes.BlockSize)
	}

	// Pad the data to the block size
	dataBytes = padData(dataBytes)

	// Create a new AES cipher block mode for CBC encryption
	mode := cipher.NewCBCEncrypter(block, iv)

	// Encrypt the data
	encryptedData := make([]byte, len(dataBytes))
	mode.CryptBlocks(encryptedData, dataBytes)

	// Encode encrypted data to base64 string
	encodedString := base64.StdEncoding.EncodeToString(encryptedData)

	return encodedString, nil
}

func decrypt(data string, keyNum int) (string, error) {
	// Decode base64-encoded string to bytes
	encryptedData, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return "", err
	}
	key := []byte("37911490979715163134003223491201")
	if keyNum > 1 {
		key = []byte("54674138327930866480207815084989")
	}
	iv := []byte("3134003223491201")

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	if len(encryptedData) < aes.BlockSize {
		return "", fmt.Errorf("encrypted data is too short")
	}

	// The IV should be the same length as the block size of the cipher
	if len(iv) != aes.BlockSize {
		return "", fmt.Errorf("IV length must be %d bytes", aes.BlockSize)
	}

	// Create a new AES cipher block mode for CBC encryption
	mode := cipher.NewCBCDecrypter(block, iv)

	// Decrypt the data
	decryptedData := make([]byte, len(encryptedData))
	mode.CryptBlocks(decryptedData, encryptedData)

	// Remove PKCS#7 padding
	padding := decryptedData[len(decryptedData)-1]
	if padding > aes.BlockSize || padding == 0 {
		return "", fmt.Errorf("invalid padding")
	}
	decryptedData = decryptedData[:len(decryptedData)-int(padding)]

	return string(decryptedData), nil
}

func getGogocdnLink(epidoseUrl string) string {
	// URL for fetching data
	baseURL := epidoseUrl //"https://anitaku.so/fairy-tail-2014-episode-1"

	resp, err := http.Get(baseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	// Parse HTML response
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Extract video URL
	gogocdn := doc.Find("#load_anime > div > div > iframe").AttrOr("src", "")
	return strings.TrimSpace(gogocdn)
}

func getDecryptedData(episodeUrl string) string {
	// Create a Resty Client
	client := resty.New()

	gogocdn := getGogocdnLink(episodeUrl)//"https://anitaku.so/sousou-no-frieren-no-mahou-episode-12"

	// println(gogocdn)

	resp2, err := http.Get((gogocdn))
	if err != nil {
		log.Fatal(err)
	}
	defer resp2.Body.Close()

	doc2, err := goquery.NewDocumentFromReader(resp2.Body)
	if err != nil {
		log.Fatal(err)
	}

	gogoUrl, _ := url.Parse(gogocdn)
	videoId := gogoUrl.Query().Get("id")

	// Extract and decrypt script value
	s := doc2.Find("script[data-name='episode']")
	scriptValue := s.AttrOr("data-value", "")

	// Decrypt the data
	decryptedData, err := decrypt(scriptValue, 1)
	if err != nil {
		log.Fatal("Decryption error:", err)
	}

	e, _ := encrypt(videoId)
	// d,_ := decrypt(e)

	// fmt.Println("Encrypted Data:", scriptValue, )
	// fmt.Println("Decrypted Data:", decryptedData)
	// fmt.Println("Video ID:", d)

	newUrl, _ := url.ParseQuery(decryptedData)
	newUrl.Set("id", e)
	newUrl.Set("alias", videoId)

	// println("url query: ", newUrl.Encode())

	vidUrl := "https://" + gogoUrl.Host + "/encrypt-ajax.php?" + newUrl.Encode()
	println(vidUrl)

	http.Get(vidUrl)

	resp, err := client.R().
		SetHeader("X-Requested-With", "XMLHttpRequest").
		Get(vidUrl)

	if err != nil {
		println("Error:", err.Error())
	}

	var respAjax struct {
		Data string `json:"data"`
	}

	err = json.Unmarshal(resp.Body(), &respAjax)

	epp, _ := decrypt(respAjax.Data, 2)
	if err != nil {
		println("Error:", err.Error())
	}
	return epp
}

func main() {
	
	epp := getDecryptedData("https://anitaku.so/sousou-no-frieren-no-mahou-episode-12")

	var src struct {
		Source []struct {
			File string `json:"File"`
		} `json:"source"`
	}
	err := json.Unmarshal([]byte(epp), &src)
	if err != nil {
		println("Error:", err.Error())
	}

	println(src.Source[0].File)

	fmt.Println("xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx")

}
