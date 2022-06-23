package modules

import (
	"Deadcord/core"
	"Deadcord/requests"
	"Deadcord/util"
	"encoding/json"
	"log"
	"strings"
)

var verify_keywords = []string{
	"verify",
	"verification",
	"rules",
	"human",
	"access",
}

func validVerifyKeyword(check_string string) bool {
	for _, verify_keyword := range verify_keywords {
		if strings.Contains(check_string, verify_keyword) {
			return true
		}
	}

	return false
}

func StartAutoVerifyThreads(server_id string) int {
	return reactVerify(server_id)
}

func reactVerify(server_id string) int {
	var channel_data core.GuildChannels
	status, status_code, channel_json := requests.SendDiscordRequest("guilds/"+server_id+"/channels", "GET", core.RawTokensLoaded[0], map[string]interface{}{}, map[string]interface{}{})

	if status && status_code == 200 {
		if err := json.Unmarshal(channel_json, &channel_data); err != nil {
			log.Fatal(err)
		}

		for _, channel_object := range channel_data {
			if validVerifyKeyword(channel_object.Name) {
				util.WriteToConsole("Found possible verify channel. Attempting to verify...", 0)

				scraped_messages, _ := GetMessages(channel_object.ID, 50, core.RawTokensLoaded[0])

				if len(scraped_messages) > 0 {
					var messages core.Message
					if err := json.Unmarshal(scraped_messages, &messages); err != nil {
						log.Fatal(err)
					}

					for _, message := range messages {
						if validVerifyKeyword(message.Content) {
							for _, reaction := range message.Reactions {
								StartReactThreads(message.ChannelID, message.ID, reaction.Emoji.Name, false)
							}
						}
					}

				} else {
					return 1
				}
			}
		}

	} else {
		return 2
	}

	return 0
}
