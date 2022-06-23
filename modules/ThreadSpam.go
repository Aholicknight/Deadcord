package modules

import (
	"Deadcord/core"
	"Deadcord/requests"
	"Deadcord/util"
	"math/rand"
	"strconv"
	"sync"
)

func StartMassThreadCreateThreads(channel_id string, thread_name string) {
	var wg sync.WaitGroup
	wg.Add(len(core.RawTokensLoaded))
	for _, token := range core.RawTokensLoaded {
		go func(token string, channel_id string, thread_name string) {
			massThreadWorker(token, channel_id, thread_name)
		}(token, channel_id, thread_name)
	}

	wg.Done()
}

func massThreadWorker(token string, channel_id string, thread_name string) {
	status, status_code, _ := requests.SendDiscordRequest("channels/"+channel_id+"/threads", "POST", token, map[string]interface{}{}, map[string]interface{}{
		"name":                  thread_name + strconv.Itoa(rand.Intn(1000-1)+1),
		"type":                  11,
		"auto_archive_duration": 1440,
	})

	if status {
		switch status_code {
		case 201:
			util.WriteToConsole("Token created thread: "+thread_name+".", 2)
		case 429:
			util.WriteToConsole("Thread request was rate limited.", 1)
		default:
			util.WriteToConsole("Token could not create thread, request failed. Code: "+strconv.Itoa(status_code), 3)
		}
	}
}
