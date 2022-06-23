package requests

import (
	"Deadcord/core"
	"Deadcord/util"
	"bytes"
	"encoding/hex"
	"encoding/json"
	"io"
	"log"
	"math/rand"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func ReadyRequestCors(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}

func GetCfData() (string, string) {
	resp, err := http.Get("https://discord.com")

	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		body_bytes, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}

		body_data := string(body_bytes)

		reg_a := regexp.MustCompile("r:'[^']*'")
		r := strings.Replace(strings.Replace(reg_a.FindStringSubmatch(body_data)[0], "r:'", "", 1), "'", "", 1)

		reg_b := regexp.MustCompile("m:'[^']*'")
		m := strings.Replace(strings.Replace(reg_b.FindStringSubmatch(body_data)[0], "m:'", "", 1), "'", "", 1)

		return r, m

	} else {
		return "", ""
	}
}

func CfbmCookieValue(r string, m string) string {
	var hexes []string

	for i := 0; i < 2; i++ {
		token := make([]byte, 16)
		if _, err := rand.Read(token); err != nil {
			log.Fatal(err)
		}

		if len(hexes) != 2 {
			hex := string(hex.EncodeToString(token))
			hexes = append(hexes, hex)
		}
	}

	cfbm_payload := core.CfbmPayload{
		M: m,
		Results: []interface{}{
			hexes[0],
			hexes[1],
		},
		Timing: rand.Intn(120-40) + 40,
		Fp: core.Fp{
			ID: 3,
			E: core.E{
				R:  []int{1920, 1080},
				Ar: []int{1040, 1920},
				Pr: 1,
				Cd: 24,
				Wb: false,
				Wp: false,
				Wn: false,
				Ch: true,
				Ws: false,
				Wd: true,
			},
		},
	}

	cookie_client := http.Client{}
	cookie_json, err := json.Marshal(cfbm_payload)

	if err != nil {
		log.Fatal(err)
	}

	req, err := http.NewRequest("POST", "https://discord.com/cdn-cgi/bm/cv/result?req_id="+r, bytes.NewBuffer(cookie_json))

	if err != nil {
		log.Fatal(err)
	}

	resp, err := cookie_client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	return resp.Cookies()[0].Value
}

func GetDiscordCookies() string {
	r, m := GetCfData()
	CfbmCookieValue(r, m)

	resp, err := http.Get("https://discord.com")

	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	return "__cf_bm=" + CfbmCookieValue(r, m) + "__dcfduid=" + resp.Cookies()[0].Value + "; __sdcfduid=" + resp.Cookies()[1].Value + "; locale=en-GB;"

}

func GetDiscordFingerprint() string {
	resp, err := http.Get("https://discord.com/api/v9/experiments")

	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		log.Fatal(err)
	}

	var science core.Science
	if err := json.Unmarshal(body, &science); err != nil {
		log.Fatal(err)
	}

	return science.Fingerprint

}

func JsonResponse(code int, message string, data map[string]interface{}) []byte {
	response := make(map[string]interface{})
	response["code"] = code
	response["message"] = message
	response["data"] = data

	switch code {
	case 200:
		util.WriteToConsole(response["message"].(string), 2)
	case 400:
		util.WriteToConsole(response["message"].(string), 1)
	case 500:
		util.WriteToConsole(response["message"].(string), 3)
	default:
		util.WriteToConsole(response["message"].(string), 0)
	}

	raw_json_response, err := json.Marshal(response)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}

	return raw_json_response
}

func GetNonce() int64 {
	nonce_raw := strconv.FormatInt((time.Now().UTC().UnixNano()/1000000)-1420070400000, 2) + "0000000000000000000000"
	nonce, _ := strconv.ParseInt(nonce_raw, 2, 64)
	return nonce
}

func ErrorResponse(message string) []byte {
	return []byte(JsonResponse(500, message, map[string]interface{}{}))
}

func AllParametersError() []byte {
	return []byte(JsonResponse(400, "All parameters must be provided.", map[string]interface{}{}))
}

func SimpleGet(url string) (int, []byte) {
	resp, err := http.Get(url)

	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		log.Fatal(err)
	}

	return resp.StatusCode, body
}
