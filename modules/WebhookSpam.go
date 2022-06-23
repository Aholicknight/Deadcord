package modules

import (
	"Deadcord/core"
	"Deadcord/requests"
	"Deadcord/util"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

func StartWebhookSpam(webhook string, username string, message string) {
	go func(webhook string) {
		for {

			if core.ActionFlag == 1 {
				return
			}

			webhook_headers := http.Header{"Content-type": []string{"application/json"}}

			status, status_code, webhook_json := requests.RequestTemplate(http.Client{}, "POST", webhook, webhook_headers, map[string]interface{}{
				"username": username,
				"content":  message,
			})

			if status {
				switch status_code {
				case 429:
					var rate_limit core.RateLimit
					if err := json.Unmarshal(webhook_json, &rate_limit); err != nil {
						log.Fatal(err)
					}

					pause_time_int := int(rate_limit.RetryAfter) / 1000

					pause_time_string := strconv.Itoa(pause_time_int)

					if pause_time_int > 0 {
						util.WriteToConsole("Webhook rate limited, pausing for: "+pause_time_string+" seconds.", 1)
					}

					util.Sleep(int(rate_limit.RetryAfter) / 1000)
				case 404:
					util.WriteToConsole("Webhook no longer exists.", 3)
					return
				case 401:
					util.WriteToConsole("An error occured while trying to spam webhook.", 3)
					return
				}
			}
		}
	}(webhook)
}
