package modules

import (
	"Deadcord/core"
	"Deadcord/requests"
	"Deadcord/util"
	"strconv"
	"sync"
)

func StartFriendThreads(user_id string) bool {
	var wg sync.WaitGroup
	wg.Add(len(core.RawTokensLoaded))
	for _, token := range core.RawTokensLoaded {
		go func(server_id string, token string) {
			friendWorker(user_id, token)
		}(user_id, token)
	}

	wg.Done()

	return false
}

func friendWorker(user_id string, token string) {
	status, status_code, _ := requests.SendDiscordRequest("users/@me/relationships/"+user_id, "PUT", token, map[string]interface{}{}, map[string]interface{}{})

	if status {
		switch status_code {
		case 204:
			util.WriteToConsole("Bot successfully sent friend request.", 2)
		case 429:
			util.WriteToConsole("Outgoing friend request was limited.", 1)
		case 404:
			util.WriteToConsole("Could not find user to send friend request.", 1)
		default:
			util.WriteToConsole("Token could not send friend request, request failed. Code: "+strconv.Itoa(status_code), 3)
		}
	}

}
