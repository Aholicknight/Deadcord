package modules

import (
	"Deadcord/core"
	"Deadcord/requests"
	"Deadcord/util"
	"encoding/json"
	"log"
	"strconv"
	"strings"
	"sync"
)

var (
	JoinResults    int  = 0
	AttemptedJoins int  = 0
	HaltJoin       bool = false
)

func StartJoinGuildThreads(invite string, server_id string, channel_id string, delay int) int {
	HaltJoin = false
	JoinResults = 0
	AttemptedJoins = 0

	if JoinResults != 0 {
		JoinResults = 0
	}

	var wg sync.WaitGroup
	wg.Add(len(core.RawTokensLoaded))

	join_channel := make(chan int)

	for _, token := range core.RawTokensLoaded {
		util.Sleep(delay)
		go func(token string, invite string, server_id string, channel_id string, join_channel chan int) {
			joinWorker(token, invite, server_id, channel_id, join_channel)
		}(token, invite, server_id, channel_id, join_channel)
	}

	join_channel_results := <-join_channel

	close(join_channel)
	wg.Done()

	return join_channel_results
}

func joinWorker(token string, invite string, server_id string, channel_id string, join_results chan int) {
	invite_clean := strings.ReplaceAll(invite, "https://", "")
	invite_parts := strings.Split(invite_clean, "/")
	invite_code := ""

	switch len(invite_parts) {
	case 1:
		invite_code = invite_parts[0]
	case 2:
		invite_code = invite_parts[1]
	case 3:
		invite_code = invite_parts[2]
	case 4:
		invite_code = invite_parts[4]
	}

	status, status_code, join_json := requests.SendDiscordRequest("invites/"+invite_code+"?inputValue="+invite+"&with_counts=true&with_expiration=true", "POST", token, map[string]interface{}{
		"location":              "Join Guild",
		"location_guild_id":     server_id,
		"location_channel_id":   channel_id,
		"location_channel_type": 0,
	}, map[string]interface{}{})

	if status {
		if !HaltJoin {
			switch status_code {
			case 200:
				var guild_data core.GuildJoin
				if err := json.Unmarshal(join_json, &guild_data); err != nil {
					log.Fatal(err)
				}

				if GuildHasCommunityVerification(server_id) {
					bypassMembershipScreen(server_id, token, invite_code)
				}

				JoinResults++
				util.WriteToConsole("Bot successfully joined guild.", 2)

			case 400:
				util.WriteToConsole("Could not join guild, captcha blocked request.", 1)
				HaltJoin = true
			case 404:
				util.WriteToConsole("Guild not found, or invite invalid.", 1)
				HaltJoin = true
			case 429:
				util.WriteToConsole("IP ratelimited or Cloudflare banned.", 1)
			default:
				util.WriteToConsole("Bot could not join guild, request failed. Code: "+strconv.Itoa(status_code), 3)
			}
		}
	}

	AttemptedJoins++

	if JoinResults == len(core.RawTokensLoaded) || AttemptedJoins == len(core.RawTokensLoaded) {
		join_results <- JoinResults
		return
	}
}

func bypassMembershipScreen(server_id string, token string, invite_code string) {
	get_status, get_status_code, get_member_json := requests.SendDiscordRequest("guilds/"+server_id+"/member-verification?with_guild=false&invite_code="+invite_code, "GET", token, map[string]interface{}{}, map[string]interface{}{})

	if get_status && get_status_code == 200 {
		util.Sleep(1)
		var membership core.MembershipScreening
		if err := json.Unmarshal(get_member_json, &membership); err != nil {
			log.Fatal(err)
		}

		for i := 0; i < len(membership.FormFields); i++ {
			membership.FormFields[i].Required = true
		}

		put_status, put_status_code, put_member_json := requests.SendDiscordRequest("guilds/"+server_id+"/requests/@me", "PUT", token, map[string]interface{}{}, map[string]interface{}{
			"version":     membership.Version,
			"form_fields": membership.FormFields,
			"description": membership.Description,
		})

		if put_status {
			switch put_status_code {
			case 201:
				var membership_status core.MembershipStatus
				if err := json.Unmarshal(put_member_json, &membership_status); err != nil {
					log.Fatal(err)
				}

				if membership_status.ApplicationStatus == "APPROVED" {
					util.WriteToConsole("Bot bypassed member screening.", 2)
				} else {
					util.WriteToConsole("Bot was not approved during member screening.", 1)
				}
			case 410:
				util.WriteToConsole("Bot already completed member screening.", 1)
			default:
				util.WriteToConsole("Bot could not bypass member screening. Code: "+strconv.Itoa(put_status_code), 3)
			}
		}
	}
}
