package modules

import (
	"Deadcord/core"
	"Deadcord/requests"
	"Deadcord/util"
	"strconv"
	"sync"
)

func StartNickThreads(server_id string, nickname string) {
	var wg sync.WaitGroup
	wg.Add(len(core.RawTokensLoaded))
	for _, token := range core.RawTokensLoaded {
		go func(token string, server_id string, nickname string) {
			nickWorker(token, server_id, nickname)
		}(token, server_id, nickname)
	}

	wg.Done()
}

func nickWorker(token string, server_id string, nickname string) {

	nickname_string := ""

	switch nickname {
	case "reset":
		nickname_string = " "
	default:
		nickname_string = nickname
	}

	status, status_code, _ := requests.SendDiscordRequest("guilds/"+server_id+"/members/@me", "PATCH", token, map[string]interface{}{}, map[string]interface{}{
		"nick": nickname_string,
	})

	if status {
		switch status_code {
		case 200:
			util.WriteToConsole("Token chnaged nickname to: "+nickname_string+".", 2)
		case 429:
			util.WriteToConsole("Change nickname request rate limited.", 1)
		default:
			util.WriteToConsole("Token could not change nickname, request failed. Code: "+strconv.Itoa(status_code), 3)
		}
	}

}
