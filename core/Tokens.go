package core

import (
	"Deadcord/util"
	"bufio"
	"fmt"
	"log"
	"os"

	ua "github.com/mileusna/useragent"
)

func LoadTokens() (bool, []string, map[int]map[string]string) {
	if _, err := os.Stat("./tokens.txt"); os.IsNotExist(err) {
		util.WriteToConsole("No token file found. Deadcord will automaticlly create one. Restart Deadcord.", 3)
		util.Sleep(5)
		token_file, err := os.Create("./tokens.txt")

		if err != nil {
			log.Fatal(err)
		}

		_, err = token_file.WriteString("! Put your tokens in this file, each token on a new line, with no quotation marks.")

		if err != nil {
			log.Fatal(err)
		}

		token_file.Close()

		return false, nil, nil
	} else {
		tokens_loaded, err := parseTokenFile("./tokens.txt")

		if err != nil {
			log.Fatal(err)
		}

		if len(tokens_loaded) > 0 {
			util.WriteToConsole("Creating unique profiles for tokens.", 0)

			label_num := 1
			var token_struct = map[int]map[string]string{}
			for _, token := range tokens_loaded {
				token_struct[label_num] = map[string]string{}
				random_agent := util.RandomUserAgent()
				parse_agent := ua.Parse(random_agent)

				token_struct[label_num]["browser"] = parse_agent.Name
				token_struct[label_num]["token"] = token
				token_struct[label_num]["agent"] = random_agent
				token_struct[label_num]["os"] = parse_agent.OS
				token_struct[label_num]["browser_version"] = parse_agent.Version
				token_struct[label_num]["os_version"] = parse_agent.OSVersion

				label_num += 1
			}

			return true, tokens_loaded, token_struct
		} else {
			util.WriteToConsole("No tokens could be loaded from ./tokens.txt.", 3)
			util.Sleep(5)
			return false, nil, nil
		}

	}
}

func ResetTokenServiceWithManualTokens(token_list []string) int {
	token_remove := os.Remove("./tokens.txt")

	if token_remove != nil {
		log.Fatal(token_remove)
	}

	WriteLines(token_list, "./tokens.txt")

	status, raw_tokens, built_tokens := LoadTokens()

	if status {
		return SetTokens(raw_tokens, built_tokens)
	} else {
		return 0
	}
}

func parseTokenFile(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func GetTokenInfo(token string) (bool, map[string]string) {
	for _, check_token_struct := range BuiltTokenStruct {
		for _, data_value := range check_token_struct {
			if data_value == token {
				return true, check_token_struct
			}
		}
	}

	return false, nil
}

func WriteLines(lines []string, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	for _, line := range lines {
		fmt.Fprintln(w, line)
	}
	return w.Flush()
}
