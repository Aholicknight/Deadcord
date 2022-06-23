package main

import (
	"Deadcord/core"
	"Deadcord/modules"
	"Deadcord/requests"
	"Deadcord/util"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	// Remove this import to build for Linux based systems.
	"golang.org/x/sys/windows"
)

var (
	DeadcordVersionNumberString = fmt.Sprintf("%g", core.DeadcordVersion)
)

func ping_tokens(w http.ResponseWriter, r *http.Request) {
	requests.ReadyRequestCors(w)

	util.WriteToConsole("Pinging tokens...", 0)

	var alive []string
	var locked []string
	var limited []string
	var invalid []string
	var cloudflare []string

	token_results := modules.StartPingTokens()

	for _, token_ping_result := range token_results {
		token_ping_split := strings.Split(token_ping_result, ":")
		token := token_ping_split[1]

		switch token_ping_split[0] {
		case "0":
			alive = append(alive, token)
		case "1":
			invalid = append(invalid, token)
		case "2":
			locked = append(locked, token)
		case "3":
			cloudflare = append(cloudflare, token)
		case "4":
			limited = append(limited, token)
		}
	}

	if len(locked) > 0 || len(invalid) > 0 {
		alive_token_list := append(alive, limited...)
		dead_token_list := append(locked, invalid...)
		core.WriteLines(dead_token_list, "dead-tokens.txt")
		core.ResetTokenServiceWithManualTokens(alive_token_list)
	}

	result_string := fmt.Sprintf(util.Green+"%d tokens alive.\n"+util.Yellow+"%d tokens invalid.\n"+util.Red+"%d tokens locked.\n"+util.Blue+"%d tokens rate-limited.\n"+util.Cyan+"%d tokens Cloudflare banned.", len(alive), len(invalid), len(locked), len(limited), len(cloudflare))
	fmt.Println(result_string)

	w.Write(requests.JsonResponse(200, "All tokens pinged: "+strconv.Itoa(len(alive))+" alive tokens.", map[string]interface{}{}))
}

func start_spam(w http.ResponseWriter, r *http.Request) {
	requests.ReadyRequestCors(w)

	util.WriteToConsole("Starting Spam...", 0)

	core.ActionFlag = 0

	r.ParseForm()
	server_id := r.Form.Get("server_id")
	channels := r.Form.Get("channels")
	messages := r.Form.Get("messages")
	spam_mode := r.Form.Get("mode")
	spam_tts := r.Form.Get("tts")

	if util.AllParameters([]string{server_id, channels, messages, spam_mode, spam_tts}) {
		spam_mode_num, err := strconv.Atoi(spam_mode)

		if err != nil {
			w.Write(requests.ErrorResponse("Invalid spam mode parameter type."))
			return
		}

		spam_tts_bool, err := strconv.ParseBool(spam_tts)

		if err != nil {
			w.Write(requests.ErrorResponse("Invalid TTS parameter type."))
			return
		}

		messages := strings.Split(messages, "\n")

		start_spam_routines := modules.StartSpamThreads(server_id, channels, messages, spam_mode_num, spam_tts_bool)

		switch start_spam_routines {
		case 1:
			w.Write(requests.ErrorResponse("Could not start spam, message content hit the character limit, or something went wrong."))
			return
		case 2:
			w.Write(requests.ErrorResponse("Could not start spam, no open channels found."))
			return
		}

	} else {
		w.Write(requests.AllParametersError())
	}

}

func stop_all(w http.ResponseWriter, r *http.Request) {
	requests.ReadyRequestCors(w)

	util.WriteToConsole("Stopping all running actions...", 0)

	core.ActionFlag = 1
	w.Write(requests.JsonResponse(200, "Attempted to stop running actions.", map[string]interface{}{}))
}

func start_typing_spam(w http.ResponseWriter, r *http.Request) {
	requests.ReadyRequestCors(w)

	util.WriteToConsole("Starting typing spam...", 0)

	core.ActionFlag = 0

	r.ParseForm()
	channel_id := r.Form.Get("channel_id")

	if util.AllParameters([]string{channel_id}) {
		modules.StartTypingSpamThreads(channel_id)

		w.Write(requests.JsonResponse(200, "Attempted to start typing spam.", map[string]interface{}{}))
	} else {
		w.Write(requests.AllParametersError())
	}
}

func react(w http.ResponseWriter, r *http.Request) {
	requests.ReadyRequestCors(w)

	util.WriteToConsole("Starting mass reacting...", 0)

	r.ParseForm()
	channel_id := r.Form.Get("channel_id")
	message_id := r.Form.Get("message_id")
	emoji := r.Form.Get("emoji")

	if util.AllParameters([]string{channel_id, message_id, emoji}) {
		modules.StartReactThreads(channel_id, message_id, emoji, true)

		w.Write(requests.JsonResponse(200, "Bots attempted to react.", map[string]interface{}{}))
	} else {
		w.Write(requests.AllParametersError())
	}
}

func change_nick(w http.ResponseWriter, r *http.Request) {
	requests.ReadyRequestCors(w)

	util.WriteToConsole("Starting mass nickname...", 0)

	r.ParseForm()
	server_id := r.Form.Get("server_id")
	nickname := r.Form.Get("nickname")

	if util.AllParameters([]string{server_id, nickname}) {
		modules.StartNickThreads(server_id, nickname)

		w.Write(requests.JsonResponse(200, "Bots attempted to change their nickname.", map[string]interface{}{}))
	} else {
		w.Write(requests.AllParametersError())
	}
}

func join_guild(w http.ResponseWriter, r *http.Request) {
	requests.ReadyRequestCors(w)

	util.WriteToConsole("Attempting to join guild...", 0)

	r.ParseForm()
	guild_invite := r.Form.Get("invite")
	join_delay := r.Form.Get("delay")

	if util.AllParameters([]string{guild_invite, join_delay}) {
		join_result_number := 0

		delay, err := strconv.Atoi(join_delay)
		if err != nil {
			w.Write(requests.ErrorResponse("Invalid delay parameter type."))
			return
		}

		server_id, channel_id, err := modules.GetGuildIdAndChannelIdFromInvite(guild_invite)

		if err != nil {
			w.Write(requests.ErrorResponse("Unable to get guild ID from invite. Invite invalid."))
			return
		}

		join_result_number = modules.StartJoinGuildThreads(guild_invite, server_id, channel_id, delay)

		if join_result_number > 0 {

			if len(server_id) > 0 {
				util.WriteToConsole("Attempting to auto-verify bots.", 2)
				status := modules.StartAutoVerifyThreads(server_id)

				switch status {
				case 1:
					w.Write(requests.ErrorResponse("No verification messages found."))
				case 2:
					w.Write(requests.ErrorResponse("Automatic verification request failed. Code not ok."))
				}
			}

		} else {
			w.Write(requests.JsonResponse(500, "Tokens could not join guild.", map[string]interface{}{}))
		}

		w.Write(requests.JsonResponse(200, strconv.Itoa(join_result_number)+"/"+strconv.Itoa(len(core.RawTokensLoaded))+" tokens joined guild.", map[string]interface{}{}))

	} else {
		w.Write(requests.AllParametersError())
	}
}

func leave_guild(w http.ResponseWriter, r *http.Request) {
	requests.ReadyRequestCors(w)

	util.WriteToConsole("Attempting to leave guild...", 0)

	r.ParseForm()
	server_id := r.Form.Get("server_id")

	if util.AllParameters([]string{server_id}) {
		modules.StartLeaveGuildThreads(server_id)

		w.Write(requests.JsonResponse(200, "Bots attempted to leave the target guild.", map[string]interface{}{}))
	} else {
		w.Write(requests.AllParametersError())
	}
}

func send_friend_requests(w http.ResponseWriter, r *http.Request) {
	requests.ReadyRequestCors(w)

	util.WriteToConsole("Attempting to send friend requests...", 0)

	r.ParseForm()
	user_id := r.Form.Get("user_id")

	if util.AllParameters([]string{user_id}) {
		modules.StartFriendThreads(user_id)

		w.Write(requests.JsonResponse(200, "Bots attempted to send target friend requests.", map[string]interface{}{}))
	} else {
		w.Write(requests.AllParametersError())
	}
}

func remove_friend_requests(w http.ResponseWriter, r *http.Request) {
	requests.ReadyRequestCors(w)

	util.WriteToConsole("Attempting to remove friend requests...", 0)

	r.ParseForm()
	user_id := r.Form.Get("user_id")

	if util.AllParameters([]string{user_id}) {
		modules.StartRemoveFriendThreads(user_id)

		w.Write(requests.JsonResponse(200, "Bots attempted to remove friend requests from target.", map[string]interface{}{}))
	} else {
		w.Write(requests.AllParametersError())
	}
}

func speak(w http.ResponseWriter, r *http.Request) {
	requests.ReadyRequestCors(w)

	util.WriteToConsole("Starting mass speak...", 0)

	r.ParseForm()
	server_id := r.Form.Get("server_id")
	message := r.Form.Get("message")

	if util.AllParameters([]string{server_id, message}) {
		modules.StartSpeakThreads(server_id, message)

		w.Write(requests.JsonResponse(200, "Bots attempted to send messages in all open channels.", map[string]interface{}{}))
	} else {
		w.Write(requests.AllParametersError())
	}
}

func start_webhook_spam(w http.ResponseWriter, r *http.Request) {
	requests.ReadyRequestCors(w)

	util.WriteToConsole("Starting webhook spam...", 0)

	core.ActionFlag = 0

	r.ParseForm()
	webhook := r.Form.Get("webhook")
	username := r.Form.Get("username")
	message := r.Form.Get("message")

	if util.AllParameters([]string{webhook, username, message}) {

		modules.StartWebhookSpam(webhook, username, message)

		w.Write(requests.JsonResponse(200, "Attempting to start webhook spam.", map[string]interface{}{}))
	} else {
		w.Write(requests.AllParametersError())
	}
}

func delete_webhook(w http.ResponseWriter, r *http.Request) {
	requests.ReadyRequestCors(w)

	util.WriteToConsole("Attempting to delete webhook...", 0)

	r.ParseForm()
	webhook := r.Form.Get("webhook")

	if util.AllParameters([]string{webhook}) {

		if modules.StartWebhookDelete(webhook) {
			w.Write(requests.JsonResponse(200, "Webhook has been deleted.", map[string]interface{}{}))
		} else {
			w.Write(requests.ErrorResponse("Could not delete webhook."))
		}

	} else {
		w.Write(requests.AllParametersError())
	}
}

func disguise_tokens(w http.ResponseWriter, r *http.Request) {
	requests.ReadyRequestCors(w)

	util.WriteToConsole("Attempting to disguise tokens...", 0)

	modules.StartDisguiseThreads()

	w.Write(requests.JsonResponse(200, "Bots attempted to disguise.", map[string]interface{}{}))
}

func start_thread_spam(w http.ResponseWriter, r *http.Request) {
	requests.ReadyRequestCors(w)

	util.WriteToConsole("Attempting to start thread spam...", 0)

	core.ActionFlag = 0

	r.ParseForm()
	channel_id := r.Form.Get("channel_id")
	thread_name := r.Form.Get("thread_name")

	if util.AllParameters([]string{channel_id, thread_name}) {
		modules.StartMassThreadCreateThreads(channel_id, thread_name)
	} else {
		w.Write(requests.AllParametersError())
	}
}

func fetch_channels(w http.ResponseWriter, r *http.Request) {
	requests.ReadyRequestCors(w)

	util.WriteToConsole("Attempting to fetch channels...", 0)

	r.ParseForm()
	server_id := r.Form.Get("server_id")

	channel_status_code, found_channels := modules.GetChannels(server_id)

	if channel_status_code != 200 {
		w.Write(requests.ErrorResponse("An error occured when attempting to fetch guild channels. Code: " + strconv.Itoa(channel_status_code)))
		return
	}

	if len(found_channels) > 0 {
		w.Write(requests.JsonResponse(200, "Successfully fetched guild channels.", map[string]interface{}{
			"channels": found_channels,
		}))
	} else {
		w.Write(requests.ErrorResponse("No open channels available."))
	}
}

var deadcord_banner string = `
    ██████╗ ███████╗ █████╗ ██████╗  ██████╗ ████████╗ ██████╗ ██████╗   ┏━━━━━━━━━━━━━━━━━━ Info ━━━━━━━━━━━━━━━━┓
    ██╔══██╗██╔════╝██╔══██╗██╔══██╗██╔════╝██████████╗██╔══██╗██╔══██╗   ` + util.Purple + `@ Package:` + util.ColorReset + ` Deadcord-Engine
    ██║  ██║█████╗  ███████║██║  ██║██║     ██║ ██  ██║██████╔╝██║  ██║   ` + util.Purple + `@ Tokens:` + util.ColorReset + ` %d tokens loaded.
    ██║  ██║██╔══╝  ██╔══██║██║  ██║██║     ████  ████║██╔══██╗██║  ██║   ` + util.Purple + `@ Warning:` + util.Red + ` Use at your own risk.` + util.ColorReset + ` 
    ██████╔╝███████╗██║  ██║██████╔╝╚██████╗╚████████╔╝██║  ██║██████╔╝   ` + util.Purple + `@ Author:` + util.ColorReset + ` https://github.com/GalaxzyDev` + util.Purple + `
    ╚═════╝ ╚══════╝╚═╝  ╚═╝╚═════╝  ╚═════╝ █═█═█═█═╝ ╚═╝  ╚═╝╚═════╝   ┗━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┛

	        ` + util.White + ` The best Discord raid tool, powerful, trusted, and purposeful. ` + util.Blue + `Using the power of Golang.` + util.Purple + `
────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────

   ` + util.White + ` You need to download our Better Discord (https://betterdiscord.app/) plugin to interact with Deadcord. You can
     download our plugin by joining our Discord server. We offer support and assistance for setting up our plugin. ` + util.Purple + `

────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────

`

func main() {

	// START Remove this block to run on Linux systems.
	stdout := windows.Handle(os.Stdout.Fd())
	var originalMode uint32

	windows.GetConsoleMode(stdout, &originalMode)
	windows.SetConsoleMode(stdout, originalMode|windows.ENABLE_VIRTUAL_TERMINAL_PROCESSING)
	// END Remove this block to run on Linux systems.

	util.WriteToConsole("Initializing output logger.", 0)
	core.InitLogger()

	token_status, raw_tokens, built_tokens := core.LoadTokens()
	proxy_status, raw_proxies, built_proxies := core.LoadProxies()

	if token_status {

		returned_token_amount := core.SetTokens(raw_tokens, built_tokens)

		if proxy_status {
			returned_proxy_amount := core.SetProxies(raw_proxies, built_proxies)
			util.WriteToConsole(strconv.Itoa(returned_proxy_amount)+" active proxies loaded.", 2)
		}

		banner_template := fmt.Sprintf(strings.ReplaceAll(deadcord_banner, "█", util.White+"█"+util.Purple), returned_token_amount)
		fmt.Println(banner_template)

		if len(raw_tokens) > 350 {
			util.WriteToConsole("Your token file exceeds the safely tested 350 token amount. Using this amount of tokens may lead to unexpected side-effects. Deadcord and the developers are not responsible for any loss of tokens. Continue at your own risk.", 1)
		}

		util.WriteToConsole("Starting Deadcord version: "+DeadcordVersionNumberString, 0)

		util.WriteToConsole(util.GetQuote(), 0)

		main_router := mux.NewRouter()

		api_router := main_router.PathPrefix("/deadcord/").Subrouter()
		api_router.HandleFunc("/ping-tokens", ping_tokens).Methods("GET")
		api_router.HandleFunc("/start-spam", start_spam).Methods("POST")
		api_router.HandleFunc("/stop-all", stop_all).Methods("GET")
		api_router.HandleFunc("/start-typing-spam", start_typing_spam).Methods("POST")
		api_router.HandleFunc("/join-guild", join_guild).Methods("POST")
		api_router.HandleFunc("/leave-guild", leave_guild).Methods("POST")
		api_router.HandleFunc("/react", react).Methods("POST")
		api_router.HandleFunc("/nick", change_nick).Methods("POST")
		api_router.HandleFunc("/disguise", disguise_tokens).Methods("GET")
		api_router.HandleFunc("/friend", send_friend_requests).Methods("POST")
		api_router.HandleFunc("/remove-friend", remove_friend_requests).Methods("POST")
		api_router.HandleFunc("/speak", speak).Methods("POST")
		api_router.HandleFunc("/start-webhook-spam", start_webhook_spam).Methods("POST")
		api_router.HandleFunc("/start-thread-spam", start_thread_spam).Methods("POST")
		api_router.HandleFunc("/delete-webhook", delete_webhook).Methods("POST")
		api_router.HandleFunc("/fetch-channels", fetch_channels).Methods("POST")

		open_port, err := net.Listen("tcp", ":0")
		if err != nil {
			log.Fatal(err)
		}

		use_open_port := strings.Split(open_port.Addr().String(), ":")[3]

		util.WriteToConsole("Deadcord is ready and running on port: "+use_open_port+".", 2)

		log.Fatal(http.Serve(open_port, main_router))
	}

}
