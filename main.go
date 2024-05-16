package main

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// Key and IV should be kept secure and unique.
var (
	keyHex = "37911490979715163134003223491201"
	key, _ = hex.DecodeString(keyHex)
	ivHex  = "3134003223491201" // IV should be unique for each encryption operation.
	iv, _  = hex.DecodeString(ivHex)
)

func encrypt(data []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// Encrypt the data using AES-GCM with the provided IV.
	return aesGCM.Seal(nil, iv, data, nil), nil
}

func decrypt(encryptedData []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// Decrypt the data using AES-GCM with the provided IV.
	return aesGCM.Open(nil, iv, encryptedData, nil)
}

func main() {
	// URL for fetching data
	baseURL := "https://anitaku.so/fairy-tail-2014-episode-1"

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
	fmt.Println("gogocdn:", strings.TrimSpace(gogocdn))

	resp2, err := http.Get(strings.TrimSpace(gogocdn))
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

	// Decode hex-encoded string to bytes
	encryptedData, err := base64.StdEncoding.DecodeString(scriptValue)
	if err != nil {
		log.Fatal(err)
	}

	// Decrypt the data
	decryptedData, err := decrypt(encryptedData)
	if err != nil {
		log.Fatal("Decryption error:", err)
	}

	fmt.Println("Encrypted Data:", scriptValue)
	fmt.Println("Decrypted Data:", string(decryptedData))
	fmt.Println("Video ID:", videoId)
}
