package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

type Client struct {
	baseURL    string
	publicKey  string
	privateKey string
	httpClient *http.Client
}

var marvelBaseURL = "https://gateway.marvel.com/v1/public/"
var publicKey, privateKey = getKeys()
var httpClient = &http.Client{
		Timeout: 10 * time.Second
}

func getKeys() (publickey string, privatekey string) {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
	pubkey := os.Getenv("MARVEL_PUBLIC_KEY")
	privkey := os.Getenv("MARVEL_PRIVATE_KEY")
	return pubkey, privkey
}
	
func NewClient(publicKey string, privateKey string) *Client {
	client := Client{
		publicKey:  publicKey,
		privateKey: privateKey,
	}
	return &client
}
// TODO create unix ts
// TODO create hash
// TODO fmt url fmt.Sprintf("")
func (c *Client) getmd5Hash(ts int64) string {
	tsHash := strconv.Itoa(int(ts))
	hash := md5.Sum([]byte(tsHash + c.privateKey + c.publicKey))
	return hex.EncodeToString(hash[:])
}

func (c *Client) signURL(url string) string {
	t := time.Now().Unix()
	hash := c.getmd5Hash(t)
	return fmt.Sprintf("%s?ts=%d&apikey=%s&hash=%s", url, t, c.publicKey, hash)
}

func (c *Client) getCharactersLimit(l int) ([]CharacterData, error) {
	url := c.baseURL + fmt.Sprintf("?limit=%d", l)
	url = c.signURL(url)

	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var characterData CharResponseHTTP
	err = json.NewDecoder(resp.Body).Decode(&characterData)
	if err != nil {
		return nil, err
	}

	return characterData.Data.Results, nil
}
