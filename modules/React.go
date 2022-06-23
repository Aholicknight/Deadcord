package modules

import (
	"Deadcord/core"
	"Deadcord/requests"
	"Deadcord/util"
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/enescakir/emoji"
)

func StartReactThreads(channel_id string, message_id string, emoji string, suffix bool) {
	var wg sync.WaitGroup
	wg.Add(len(core.RawTokensLoaded))
	for _, token := range core.RawTokensLoaded {
		go func(token string, channel_id string, message_id string, emoji string, suffix bool) {
			reactWorker(token, channel_id, message_id, emoji, suffix)
		}(token, channel_id, message_id, emoji, suffix)
	}

	wg.Done()
}

func reactWorker(token string, channel_id string, message_id string, emoji_string string, suffix bool) {
	react_emoji := ""

	if suffix {
		react_emoji = strings.TrimSuffix(emoji.Parse(":"+emoji_string+":"), " ")
		fmt.Println(react_emoji)
	} else {
		react_emoji = emoji_string
	}

	if !strings.Contains(react_emoji, ":") {
		status, status_code, _ := requests.SendDiscordRequest("channels/"+channel_id+"/messages/"+message_id+"/reactions/"+react_emoji+"/@me", "PUT", token, map[string]interface{}{}, map[string]interface{}{})

		if status {
			switch status_code {
			case 204:
				util.WriteToConsole("Bot reacted with: [ "+react_emoji+" ].", 2)
			case 429:
				util.WriteToConsole("Reaction request was rate limited.", 1)
			default:
				util.WriteToConsole("Bot could not react, request failed. Code:  "+strconv.Itoa(status_code), 3)
			}
		}

	} else {
		util.WriteToConsole("Do not use ':' in the emoji name.", 1)
	}
}
