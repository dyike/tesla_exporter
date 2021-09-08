package tesla

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

func Login(email, password string) string {
	// Step1 generate Verifier and Chanllenge
	verifier, challenge := genVerifierAndChallenge()
	// Step2 Get Login html
	inputMap, cookie := getLoginHtml(challenge)
	// Step3 Login with email password
	code := loginWithEmailPassword(email, password, challenge, inputMap, cookie)

	// Step4 Get bear access token by code above
	bearToken := getBearAccessToken(code, verifier)
	// Step5 Exchange the bearer access token to real access token
	accessToken := exchangeBearToRealAccessToken((bearToken["access_token"]).(string))
	return accessToken
}

// 1
func genVerifierAndChallenge() (verifier string, challenge string) {
	verifier = genVerifier()
	challenge = genChallenge()
	return
}

// 2
func getLoginHtml(challenge string) (map[string]string, string) {
	params := url.Values{}
	params.Add("locale", "zh-CN")
	params.Add("client_id", "ownerapi")
	params.Add("code_challenge", challenge)
	params.Add("code_challenge_method", "S256")
	params.Add("redirect_uri", "https://auth.tesla.com/void/callback")
	params.Add("response_type", "code")
	params.Add("scope", "openid email offline_access")
	params.Add("state", genState())
	params.Add("prompt", "login")

	reqUrl, _ := url.Parse("https://auth-global.tesla.com/oauth2/v3/authorize")
	reqUrl.RawQuery = params.Encode()
	urlPath := reqUrl.String()
	resp, err := http.Get(urlPath)
	if err != nil {

	}
	defer resp.Body.Close()

	respHeaders := resp.Header.Get("set-cookie")
	body, _ := ioutil.ReadAll(resp.Body)

	nameList := regMatch(string(body), "input type", "name")
	valueList := regMatch(string(body), "input type", "value")
	inputMap := make(map[string]string)
	for i, name := range nameList {
		inputMap[name] = valueList[i]
	}
	inputMap["_phase"] = "authenticate"
	inputMap["_process"] = "1"
	inputMap["privacy_consent"] = "1"

	return inputMap, respHeaders
}

// 3
func loginWithEmailPassword(email, password, challenge string, input map[string]string, cookie string) string {
	params := url.Values{}
	params.Add("locale", "zh-CN")
	params.Add("client_id", "ownerapi")
	params.Add("code_challenge", challenge)
	params.Add("code_challenge_method", "S256")
	params.Add("redirect_uri", "https://auth.tesla.com/void/callback")
	params.Add("response_type", "code")
	params.Add("scope", "openid email offline_access")
	params.Add("state", genState())
	params.Add("prompt", "login")

	reqUrl, _ := url.Parse("https://auth.tesla.cn/oauth2/v3/authorize")
	reqUrl.RawQuery = params.Encode()
	urlPath := reqUrl.String()

	urlValues := url.Values{}
	urlValues.Set("identity", email)
	urlValues.Set("credential", password)
	for k, v := range input {
		urlValues.Set(k, v)
	}

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	req, _ := http.NewRequest("POST", urlPath, strings.NewReader(urlValues.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Cookie", cookie)
	resp, err := client.Do(req)
	if err != nil {

	}
	defer resp.Body.Close()

	location, _ := resp.Location()
	code := location.Query().Get("code")
	return code
}

// 4
func getBearAccessToken(code, verifier string) map[string]interface{} {
	client := &http.Client{}
	data := make(map[string]interface{})
	data["grant_type"] = "authorization_code"
	data["client_id"] = "ownerapi"
	data["code"] = code
	data["code_verifier"] = verifier
	data["redirect_uri"] = "https://auth.tesla.com/void/callback"
	bytesData, _ := json.Marshal(data)
	req, _ := http.NewRequest("POST", "https://auth.tesla.cn/oauth2/v3/token", bytes.NewReader(bytesData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Linux; Android 10; Pixel 3a Build/QQ1A.200205.002; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/80.0.3987.162")
	resp, err := client.Do(req)
	if err != nil {
	}
	defer resp.Body.Close()
	res, _ := ioutil.ReadAll(resp.Body)
	m := make(map[string]interface{})
	json.Unmarshal(res, &m)
	return m
}

// 5
func exchangeBearToRealAccessToken(bearToken string) string {
	client := &http.Client{}
	params := make(map[string]string)
	params["grant_type"] = "urn:ietf:params:oauth:grant-type:jwt-bearer"
	params["client_id"] = "81527cff06843c8634fdc09e8ac0abefb46ac849f38fe1e431c2ef2106796384"
	bytesData, _ := json.Marshal(params)
	req, _ := http.NewRequest("POST", "https://owner-api.teslamotors.com/oauth/token", bytes.NewReader(bytesData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("authorization", "Bearer "+bearToken)

	resp, err := client.Do(req)
	if err != nil {

	}
	defer resp.Body.Close()

	res, _ := ioutil.ReadAll(resp.Body)
	return string(res)
}
