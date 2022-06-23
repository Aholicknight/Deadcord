package modules

import (
	"Deadcord/core"
	"Deadcord/requests"
	"Deadcord/util"
	b64 "encoding/base64"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"sync"
)

var image_urls = [2]string{
	"https://picsum.photos/512/512",
	"https://cataas.com/cat?width=512&height=512",
}

func StartDisguiseThreads() {
	var wg sync.WaitGroup
	wg.Add(len(core.RawTokensLoaded))

	for _, token := range core.RawTokensLoaded {
		go func(token string) {
			disguiseWorker(token)
		}(token)
	}
}

func disguiseWorker(token string) {

	random_image_api := image_urls[rand.Intn(len(image_urls))]

	resp, err := http.Get(random_image_api)

	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatal(err)
	}

	image := "data:image/png;base64," + b64.StdEncoding.EncodeToString(body)

	status, status_code, _ := requests.SendDiscordRequest("users/@me", "PATCH", token, map[string]interface{}{}, map[string]interface{}{
		"avatar": image,
	})

	if status {
		switch status_code {
		case 200:
			util.WriteToConsole("Successfully changed token profile picture.", 2)
		case 429:
			util.WriteToConsole("IP ratelimited or Cloudflare banned.", 1)
		default:
			util.WriteToConsole("Token could not chnage profile picture, request failed. Code: "+strconv.Itoa(status_code), 3)
		}
	}
}
