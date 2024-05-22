package gogoanime

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github/mdtosif/go-anime/provider/interfacetype"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/go-resty/resty/v2"
)

var client = resty.New()

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

func getGogocdnLink(epidoseUrl string) (string, error) {
	// URL for fetching data
	baseURL := epidoseUrl //"https://anitaku.so/fairy-tail-2014-episode-1"

	resp, err := http.Get(baseURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Parse HTML response
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", err
	}

	// Extract video URL
	gogocdn := doc.Find("#load_anime > div > div > iframe").AttrOr("src", "")
	return strings.TrimSpace(gogocdn), nil
}

func getDecryptedData(episodeUrl string) (string, error) {

	gogocdn, err := getGogocdnLink(episodeUrl)

	if err != nil {
		return "", err
	}

	resp2, err := http.Get((gogocdn))
	if err != nil {
		return "", err
	}
	defer resp2.Body.Close()

	doc2, err := goquery.NewDocumentFromReader(resp2.Body)
	if err != nil {
		return "", err
	}

	gogoUrl, _ := url.Parse(gogocdn)
	videoId := gogoUrl.Query().Get("id")

	// Extract and decrypt script value
	s := doc2.Find("script[data-name='episode']")
	scriptValue := s.AttrOr("data-value", "")

	// Decrypt the data
	decryptedData, err := decrypt(scriptValue, 1)
	if err != nil {
		return "", err
	}

	e, _ := encrypt(videoId)

	newUrl, _ := url.ParseQuery(decryptedData)
	newUrl.Set("id", e)
	newUrl.Set("alias", videoId)

	vidUrl := "https://" + gogoUrl.Host + "/encrypt-ajax.php?" + newUrl.Encode()

	http.Get(vidUrl)

	resp, err := client.R().
		SetHeader("X-Requested-With", "XMLHttpRequest").
		Get(vidUrl)

	if err != nil {
		return "", err
	}

	var respAjax struct {
		Data string `json:"data"`
	}

	err = json.Unmarshal(resp.Body(), &respAjax)

	epp, _ := decrypt(respAjax.Data, 2)
	if err != nil {
		return "", err
	}
	return epp, nil
}



func getAnimeSrcs(episodeId string) ([]interfacetype.AnimeEpisodeSource, error) {
	var episodesSrc []interfacetype.AnimeEpisodeSource
	epp, err := getDecryptedData("https://anitaku.so/" + episodeId)
	if err != nil {
		return nil, err
	}

	var src struct {
		Source []struct {
			File string `json:"File"`
		} `json:"source"`
	}

	err = json.Unmarshal([]byte(epp), &src)
	if err != nil {
		return nil, err
	}

	fileUrl := src.Source[0].File
	lastIndexOfSlash := strings.LastIndex(fileUrl, "/")

	fileHost := fileUrl[:lastIndexOfSlash]

	m3u8File, _ := client.R().Get(fileUrl)

	// Define the regular expression pattern to extract stream information
	pattern := `#EXT-X-STREAM-INF:PROGRAM-ID=\d+,BANDWIDTH=\d+,RESOLUTION=(\d+x\d+),NAME="([^"]+)"`
	re := regexp.MustCompile(pattern)

	// Find all matches in the EXT M3U content
	matches := re.FindAllStringSubmatch(m3u8File.String(), -1)

	// Extract m3u8 links along with names and resolutions
	for _, match := range matches {
		resolution := match[1]
		name := match[2]

		// Extract m3u8 link
		link := ""
		lines := strings.Split(m3u8File.String(), "\n")
		for idx, line := range lines {
			if strings.HasPrefix(line, "#EXT-X-STREAM-INF") && strings.Contains(line, fmt.Sprintf("RESOLUTION=%s", resolution)) && strings.Contains(line, fmt.Sprintf("NAME=\"%s\"", name)) {
				link = lines[idx+1]
				break
			}
		}

		// Print extracted information
		episodesSrc = append(episodesSrc, interfacetype.AnimeEpisodeSource{
			Name:    name,
			Quality: resolution,
			Url:     fileHost + "/" + link,
		})
	}
	return episodesSrc, nil
}

// eps, err := GetAnimeSrcs("sousou-no-frieren-no-mahou-episode-12")
// if err != nil{
// 	println(err.Error(), err)
// }
// for _, v := range eps {
// 	println("name: ", v.Name, "Url: ", v.Url)
// }
