package modules

import (
	"Deadcord/core"
	"Deadcord/requests"
	"encoding/json"
	"errors"
	"log"
	"strconv"
	"strings"
)

func InServer(server_id string, token string) bool {
	status, status_code, _ := requests.SendDiscordRequest("guilds/"+server_id, "GET", token, map[string]interface{}{}, map[string]interface{}{})

	if status && status_code == 200 {
		return true
	} else {
		return false
	}
}

func GetMessages(channel_id string, amount int, token string) ([]byte, error) {
	status, status_code, messages_json := requests.SendDiscordRequest("channels/"+channel_id+"/messages?limit="+strconv.Itoa(amount), "GET", token, map[string]interface{}{}, map[string]interface{}{})

	if status && status_code == 200 {
		return messages_json, nil
	} else {
		return nil, errors.New("get messages request failed, code not ok")
	}
}

func GetGuildIdAndChannelIdFromInvite(invite string) (string, string, error) {
	invite_parts := strings.Split(invite, "/")
	invite_code := ""

	if len(invite_parts) > 1 {
		if invite_parts[3] == "invite" {
			invite_code = invite_parts[4]
		} else {
			invite_code = invite_parts[3]
		}
	} else {
		invite_code = invite_parts[0]
	}

	status, status_code, invite_json := requests.SendDiscordRequest("invites/"+invite_code, "GET", core.RawTokensLoaded[0], map[string]interface{}{}, map[string]interface{}{})

	if status && status_code == 200 {
		var invite core.Invite
		if err := json.Unmarshal(invite_json, &invite); err != nil {
			log.Fatal(err)
		}

		return invite.Guild.ID, invite.Channel.ID, nil
	} else {
		return "", "", errors.New("get guild from invite request failed, code not ok")
	}
}

func GetChannels(server_id string) (int, []string) {
	var channels []string

	channel_status, channel_status_code, channel_json := requests.SendDiscordRequest("guilds/"+server_id+"/channels", "GET", core.RawTokensLoaded[0], map[string]interface{}{}, map[string]interface{}{})

	var result core.GuildChannels
	if err := json.Unmarshal(channel_json, &result); err != nil {
		log.Fatal(err)
	}

	if channel_status && channel_status_code == 200 {
		for _, channel := range result {
			if channel.Type == 0 {
				channels = append(channels, channel.Name+":"+channel.ID)
			}
		}
	}

	return channel_status_code, channels
}

func GuildHasCommunityVerification(server_id string) bool {
	status, status_code, guild_json := requests.SendDiscordRequest("guilds/"+server_id, "GET", core.RawTokensLoaded[0], map[string]interface{}{}, map[string]interface{}{})

	if status && status_code == 200 {
		var guild core.Guild
		if err := json.Unmarshal(guild_json, &guild); err != nil {
			log.Fatal(err)
		}

		for _, feature := range guild.Features {
			if feature == "MEMBER_VERIFICATION_GATE_ENABLED" {
				return true
			}
		}
	}

	return false
}
