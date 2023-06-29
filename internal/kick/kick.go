package kick

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/url"
	"strings"

	http "github.com/bogdanfinn/fhttp"
)

func (client *Client) getSpecificCookie(domain, cookieName string) string {
	_domain, err := url.Parse(fmt.Sprintf(`https://%s`, domain))
	if err != nil {
		panic(err)
	}

	cookies := client.request.GetCookies(_domain)

	for _, cookie := range cookies {
		if strings.EqualFold(strings.ToLower(cookie.Name), strings.ToLower(cookieName)) {
			return cookie.Value
		}
	}

	return ""
}

func (client *Client) GetCookies() {
	req, err := http.NewRequest("GET", "https://kick.com/", nil)
	if err != nil {
		log.Fatal(err)
		return
	}
	req.Header = http.Header{
		"User-Agent":                {"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:109.0) Gecko/20100101 Firefox/112.0"},
		"Accept":                    {"text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8"},
		"Accept-Language":           {"en-US,en;q=0.5"},
		"DNT":                       {"1"},
		"Upgrade-Insecure-Requests": {"1"},
		"Connection":                {"keep-alive"},
		"Sec-Fetch-Dest":            {"document"},
		"Sec-Fetch-Mode":            {"navigate"},
		"Sec-Fetch-Site":            {"none"},
		"Sec-Fetch-User":            {"?1"},
		"Pragma":                    {"no-cache"},
		"Cache-Control":             {"no-cache"},
	}
	resp, err := client.request.Do(req)
	if err != nil {
		log.Fatal(err)
		return
	}

	if resp.StatusCode != 200 {
		log.Println("error: on requesting cookies..")
		return
	}
	log.Println("successfully requested cookies..")
	client.xsrf = client.getSpecificCookie("kick.com", "XSRF-TOKEN")
}

func (client *Client) RequestTokenProvider() {
	log.Println(client.socketID)
	req, err := http.NewRequest("GET", "https://kick.com/kick-token-provider", nil)
	if err != nil {
		log.Fatal(err)
		return
	}
	req.Header = http.Header{
		"User-Agent":                {"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:109.0) Gecko/20100101 Firefox/112.0"},
		"Accept":                    {"text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8"},
		"Accept-Language":           {"en-US,en;q=0.5"},
		"DNT":                       {"1"},
		"Upgrade-Insecure-Requests": {"1"},
		"Connection":                {"keep-alive"},
		"Sec-Fetch-Dest":            {"document"},
		"Sec-Fetch-Mode":            {"navigate"},
		"Sec-Fetch-Site":            {"none"},
		"Sec-Fetch-User":            {"?1"},
		"Pragma":                    {"no-cache"},
		"Cache-Control":             {"no-cache"},
		"X-Socket-ID":               {client.socketID},
		"Referer":                   {"https://kick.com/"},
		"Authorization":             {"Bearer " + client.xsrf},
		"X-XSRF-TOKEN":              {client.xsrf},
	}

	resp, err := client.request.Do(req)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Println("error: on requesting cookies..")
		return
	}

	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	json.Unmarshal(bodyText, &client.form)
	log.Println("successfully requested token provider..")
}

func (client *Client) SendEmail() {
	req, err := http.NewRequest("POST", "https://kick.com/api/v1/signup/send/email", strings.NewReader(fmt.Sprintf(`{"email":"%s"}`, client.Email)))
	if err != nil {
		log.Fatal(err)
		return
	}
	req.Header = http.Header{
		"User-Agent":      {"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:109.0) Gecko/20100101 Firefox/112.0"},
		"Accept":          {"application/json, text/plain, */*"},
		"Accept-Language": {"en-US"},
		"Content-Type":    {"application/json"},
		"X-Socket-ID":     {client.socketID},
		"Authorization":   {"Bearer " + client.xsrf},
		"X-XSRF-TOKEN":    {client.xsrf},
		"Origin":          {"https://kick.com"},
		"DNT":             {"1"},
		"Connection":      {"keep-alive"},
		"Referer":         {"https:/`/kick.com/"},
		"Sec-Fetch-Dest":  {"empty"},
		"Sec-Fetch-Mode":  {"cors"},
		"Sec-Fetch-Site":  {"same-origin"},
		"Pragma":          {"no-cache"},
		"Cache-Control":   {"no-cache"},
	}

	resp, err := client.request.Do(req)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 204 {
		fmt.Println(resp.StatusCode)
		log.Println("error: on email send..")
		bodyText, err := io.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}
		fmt.Println(string(bodyText))
		return
	}
	log.Println("successfully sent email..")
}

func (client *Client) SendEmailCode(code string) {
	req, err := http.NewRequest("POST", "https://kick.com/api/v1/signup/verify/code", strings.NewReader(fmt.Sprintf(`{"code":"%s","email":"%s"}`, code, client.Email)))
	if err != nil {
		log.Fatal(err)
		return
	}
	req.Header = http.Header{
		"User-Agent":      {"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:109.0) Gecko/20100101 Firefox/112.0"},
		"Accept":          {"application/json, text/plain, */*"},
		"Accept-Language": {"en-US"},
		"Content-Type":    {"application/json"},
		"X-Socket-ID":     {client.socketID},
		"Authorization":   {"Bearer " + client.xsrf},
		"X-XSRF-TOKEN":    {client.xsrf},
		"Origin":          {"https://kick.com"},
		"DNT":             {"1"},
		"Connection":      {"keep-alive"},
		"Referer":         {"https:/`/kick.com/"},
		"Sec-Fetch-Dest":  {"empty"},
		"Sec-Fetch-Mode":  {"cors"},
		"Sec-Fetch-Site":  {"same-origin"},
		"Pragma":          {"no-cache"},
		"Cache-Control":   {"no-cache"},
	}

	resp, err := client.request.Do(req)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 204 {
		fmt.Println(resp.StatusCode)
		log.Println("error: on email send..")
		bodyText, err := io.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}
		fmt.Println(string(bodyText))
		return
	}

	log.Println("successfully sent email code..")
}

func (client *Client) RegisterAccount(username string) (string, error) {
	payload := fmt.Sprintf(
		`{"birthdate":"03/27/1995","username":"%s","email":"%s","cf_captcha_token":"","password":"%s","password_confirmation":"%s","agreed_to_terms":true,"%s":"","_kick_token_valid_from":"%s"}`,
		username,
		client.Email,
		client.Password,
		client.Password,
		client.form.NameFieldName,
		client.form.EncryptedValidFrom,
	)

	req, err := http.NewRequest("POST", "https://kick.com/register", strings.NewReader(payload))
	if err != nil {
		log.Fatal(err)
		return username, err
	}
	req.Header = http.Header{
		"User-Agent":      {"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:109.0) Gecko/20100101 Firefox/112.0"},
		"Accept":          {"application/json, text/plain, */*"},
		"Accept-Language": {"en-US"},
		"Content-Type":    {"application/json"},
		"X-Socket-ID":     {client.socketID},
		"Authorization":   {"Bearer " + client.xsrf},
		"X-XSRF-TOKEN":    {client.xsrf},
		"Origin":          {"https://kick.com"},
		"DNT":             {"1"},
		"Connection":      {"keep-alive"},
		"Referer":         {"https:/`/kick.com/"},
		"Sec-Fetch-Dest":  {"empty"},
		"Sec-Fetch-Mode":  {"cors"},
		"Sec-Fetch-Site":  {"same-origin"},
		"Pragma":          {"no-cache"},
		"Cache-Control":   {"no-cache"},
	}

	resp, err := client.request.Do(req)
	if err != nil {
		log.Fatal(err)
		return username, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 || resp.StatusCode == 201 {
		log.Println("successfully sent register..")
		return username, nil
	}

	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	log.Println("error: on register register..")

	return username, errors.New(string(bodyText))
}
