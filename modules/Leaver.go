package modules

import (
	"Deadcord/core"
	"Deadcord/requests"
	"Deadcord/util"
	"strconv"
	"sync"
)

func StartLeaveGuildThreads(server_id string) bool {
	var wg sync.WaitGroup
	wg.Add(len(core.RawTokensLoaded))
	for _, token := range core.RawTokensLoaded {
		go func(server_id string, token string) {
			leaveWorker(server_id, token)
		}(server_id, token)
	}

	wg.Done()

	return false
}

func leaveWorker(server_id string, token string) {
	status, status_code, _ := requests.SendDiscordRequest("users/@me/guilds/"+server_id, "DELETE", token, map[string]interface{}{}, map[string]interface{}{
		"lurking": false,
	})

	if status {
		if status_code == 204 {
			util.WriteToConsole("Bot successfully left guild.", 2)
		} else {
			util.WriteToConsole("Token could not leave guild, request failed. Code: "+strconv.Itoa(status_code), 3)
		}
	}

}
