package modules

import (
	"Deadcord/requests"
	"net/http"
)

func StartWebhookDelete(webhook string) bool {
	webhook_headers := http.Header{"Content-type": []string{"application/json"}}

	status, status_code, _ := requests.RequestTemplate(http.Client{}, "DELETE", webhook, webhook_headers, map[string]interface{}{})

	if status && status_code == 204 {
		return true
	}

	return true
}
