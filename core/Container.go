package core

import (
	"math/rand"
)

var (
	ActionFlag       = 0
	DeadcordVersion  = 1.4
	RawTokensLoaded  []string
	BuiltTokenStruct map[int]map[string]string

	RawProxiesLoaded []string
	BuiltProxyStruct map[int]map[string]string
)

func SetTokens(raw_token_list []string, built_token_list map[int]map[string]string) int {
	RawTokensLoaded = raw_token_list
	BuiltTokenStruct = built_token_list

	return len(RawTokensLoaded)
}

func SetProxies(raw_proxy_list []string, built_proxy_list map[int]map[string]string) int {
	RawProxiesLoaded = raw_proxy_list
	BuiltProxyStruct = built_proxy_list

	return len(RawProxiesLoaded)
}

func RandomToken() string {
	random_token := RawTokensLoaded[rand.Intn(len(RawTokensLoaded))]
	return random_token
}
