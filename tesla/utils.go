package tesla

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"regexp"
)

func PostRequest(client *http.Client, token *Token, url string, payload []byte) ([]byte, error) {
	req, err := http.NewRequest("POST", BaseURL+url, bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}
	req.Header.Add("User-Agent", UserAgent)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	if token != nil {
		req.Header.Add("Authorization", "Bearer "+token.AccessToken)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}
	return body, nil
}

func GetRequest(client *http.Client, token *Token, url string) ([]byte, error) {
	req, err := http.NewRequest("GET", BaseURL+url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("User-Agent", UserAgent)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	if token != nil {
		req.Header.Add("Authorization", "Bearer "+token.AccessToken)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s", http.StatusText(resp.StatusCode))
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func genVerifier() (verifier string) {
	var p [87]byte
	if _, err := io.ReadFull(rand.Reader, p[:]); err != nil {
		return ""
	}
	verifier = base64.RawURLEncoding.EncodeToString(p[:])
	return
}

func genState() (state string) {
	var p [21]byte
	if _, err := io.ReadFull(rand.Reader, p[:]); err != nil {
		return ""
	}
	state = base64.RawURLEncoding.EncodeToString(p[:])
	return
}

func genChallenge() (challenge string) {
	b := sha256.Sum256([]byte(challenge))
	challenge = base64.RawURLEncoding.EncodeToString(b[:])
	return
}

func regMatch(source, element, attr string) []string {
	reg := "<" + element + "[^<>]*?\\s" + attr + "=['\"]?(.*?)['\"]?(\\s.*?)?>"
	re := regexp.MustCompile(reg)
	match := re.FindAllStringSubmatch(source, -1)
	var res = make([]string, len(match))
	for i, item := range match {
		res[i] = item[1]
	}
	return res
}
