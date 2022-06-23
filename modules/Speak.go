package modules

import (
	"Deadcord/core"
	"Deadcord/requests"
	"Deadcord/util"
	"strconv"
	"strings"
	"sync"
)

func StartSpeakThreads(server_id string, message string) bool {
	var wg sync.WaitGroup
	wg.Add(len(core.RawTokensLoaded))

	var speak_channels []string
	channel_status_code, raw_channel_list := GetChannels(server_id)

	for _, channel := range raw_channel_list {
		channel_chunks := strings.Split(channel, ":")
		speak_channels = append(speak_channels, channel_chunks[1])
	}

	if channel_status_code != 200 {
		return false
	}

	for _, channel := range speak_channels {
		go func(channel string, message string) {
			speakWorker(channel, message)
		}(channel, message)
	}

	wg.Done()

	return false
}

func speakWorker(channel_id string, message string) {
	use_token := core.RandomToken()

	status, status_code, _ := requests.SendDiscordRequest("channels/"+channel_id+"/messages", "POST", use_token, map[string]interface{}{}, map[string]interface{}{
		"content": message,
		"nonce":   requests.GetNonce(),
		"tts":     false,
	})

	if status {
		switch status_code {
		case 200:
			util.WriteToConsole("Bot succesfully sent message.", 2)
		case 403:
			util.WriteToConsole("Bot could not send message, no access.", 1)
		default:
			util.WriteToConsole("Token could not send message, request failed. Code: "+strconv.Itoa(status_code), 3)
		}

	}

}
