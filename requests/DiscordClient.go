package requests

import (
	"Deadcord/core"
	"bytes"
	"crypto/tls"
	b64 "encoding/base64"
	"encoding/json"
	"io"
	"log"
	"math/rand"
	"net/http"
)

var (
	CookieString      = GetDiscordCookies()
	GlobalFingerprint = GetDiscordFingerprint()
	BaseURLs          = []string{
		"https://discord.com/api/v9/",
		"https://ptb.discord.com/api/v9/",
		"https://canary.discord.com/api/v9/",
	}
)

func SendDiscordRequest(endpoint string, method string, token string, custom_context map[string]interface{}, data map[string]interface{}) (bool, int, []byte) {

	discord_client := http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				MinVersion: tls.VersionTLS13,
				CipherSuites: []uint16{
					tls.TLS_AES_128_GCM_SHA256,
					tls.TLS_AES_256_GCM_SHA384,
					tls.TLS_CHACHA20_POLY1305_SHA256,
				},
				InsecureSkipVerify: true,
				CurvePreferences:   []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
			},
		},
	}

	discord_client_build := "122222"

	_, token_data := core.GetTokenInfo(token)

	x_super_props := map[string]string{
		"os":                       token_data["os"],
		"browser":                  token_data["browser"],
		"device":                   "",
		"system_locale":            "en-US",
		"browser_user_agent":       token_data["agent"],
		"browser_version":          token_data["browser_version"],
		"os_version":               token_data["os_version"],
		"referrer":                 "",
		"referring_domain":         "",
		"referrer_current":         "",
		"referring_domain_current": "",
		"release_channel":          "stable",
		"client_build_number":      discord_client_build,
		"client_event_source":      "None",
	}

	x_track_props := map[string]string{
		"os":                       token_data["os"],
		"browser":                  token_data["browser"],
		"device":                   "",
		"system_locale":            "en-US",
		"browser_user_agent":       token_data["agent"],
		"browser_version":          token_data["browser_version"],
		"os_version":               token_data["os_version"],
		"referrer":                 "",
		"referring_domain":         "",
		"referrer_current":         "",
		"referring_domain_current": "",
		"release_channel":          "stable",
		"client_build_number":      discord_client_build,
		"client_event_source":      "None",
	}

	x_super_props_json, err := json.Marshal(x_super_props)

	if err != nil {
		log.Fatal(err)
	}

	x_track_props_json, err := json.Marshal(x_track_props)

	if err != nil {
		log.Fatal(err)
	}

	x_super_props_refined := b64.StdEncoding.EncodeToString(x_super_props_json)
	x_track_props_refined := b64.StdEncoding.EncodeToString(x_track_props_json)

	discord_headers := http.Header{
		"Host":               []string{"discord.com"},
		"Accept":             []string{"*/*"},
		"Accept-language":    []string{"en-GB"},
		"Authorization":      []string{token},
		"Alt-Used":           []string{"discord.com"},
		"Content-type":       []string{"application/json"},
		"Cookie":             []string{CookieString},
		"DNT":                []string{"1"},
		"Origin":             []string{"https://discord.com"},
		"Referer":            []string{"https://discord.com/channels/@me"},
		"Sec-fetch-dest":     []string{"empty"},
		"Sec-fetch-mode":     []string{"cors"},
		"Sec-fetch-site":     []string{"same-origin"},
		"Sec-ch-ua":          []string{"Not A;Brand';v='99', 'Chromium';v='96', 'Google Chrome';v='96'"},
		"Sec-ch-ua-mobile":   []string{"0"},
		"Sec-ch-ua-platform": []string{token_data["os"]},
		"TE":                 []string{"Trailers"},
		"User-Agent":         []string{token_data["agent"]},
		"X-Debug-options":    []string{"bugReporterEnabled"},
		"X-Discord-Locale":   []string{"en-US"},
		"X-Track":            []string{x_track_props_refined},
		"X-Fingerprint":      []string{GlobalFingerprint},
		"X-Super-Properties": []string{x_super_props_refined},
	}

	if len(custom_context) > 0 {
		x_context_props_json, err := json.Marshal(custom_context)

		if err != nil {
			log.Fatal(err)
		}

		x_context_props_refined := b64.StdEncoding.EncodeToString(x_context_props_json)

		discord_headers["X-Context-Properties"] = []string{x_context_props_refined}
	}

	used_base_url := BaseURLs[rand.Intn(len(BaseURLs))]

	switch method {
	case "GET":
		status, status_code, body := GetRequestTemplate(discord_client, used_base_url+endpoint, discord_headers)
		return status, status_code, body
	case "POST":
		status, status_code, body := RequestTemplate(discord_client, "POST", used_base_url+endpoint, discord_headers, data)
		return status, status_code, body
	case "PUT":
		status, status_code, body := RequestTemplate(discord_client, "PUT", used_base_url+endpoint, discord_headers, data)
		return status, status_code, body
	case "PATCH":
		status, status_code, body := RequestTemplate(discord_client, "PATCH", used_base_url+endpoint, discord_headers, data)
		return status, status_code, body
	case "DELETE":
		status, status_code, body := RequestTemplate(discord_client, "DELETE", used_base_url+endpoint, discord_headers, data)
		return status, status_code, body
	}

	return false, 0, nil
}

func RequestTemplate(client http.Client, request_type string, url string, headers http.Header, json_payload map[string]interface{}) (bool, int, []byte) {

	patch_json, err := json.Marshal(json_payload)

	if err != nil {
		log.Fatal(err)
	}

	req, err := http.NewRequest(request_type, url, bytes.NewBuffer(patch_json))

	if err != nil {
		log.Fatal(err)
	}

	req.Header = headers

	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)

	if err != nil {
		log.Fatal(err)
	}

	return true, res.StatusCode, []byte(body)
}

func GetRequestTemplate(client http.Client, url string, headers http.Header) (bool, int, []byte) {

	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		log.Fatal(err)
	}

	req.Header = headers

	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)

	if err != nil {
		log.Fatal(err)
	}

	return true, res.StatusCode, []byte(body)
}
