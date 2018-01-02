package Objects

import (
	"os"
	"log"
	"bufio"
	"strings"
	"strconv"
	"github.com/l1k3ab0t/GoJugg/lib/FormatHTML"
)

type Settings struct {
	webcfg     bool
	name       string
	list       string
	fields     int
	gameMode   int
	groupCount  int
	consoleLog bool
	port       int
	aTCount    int
}

type line struct {
	Linenumber int
	Content    string
}

func readFile(path string) []line {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	i := 0
	var str []line
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if scanner.Text() != "" && strings.Index(scanner.Text(), "#") != 0 {
			i++
			str = append(str, line{i, scanner.Text()})
		}
	}
	file.Close()
	return str
}

func splitConfig(config string) []string {
	splits := strings.Split(config, "=")
	return splits

}

func readTeamList(f string) []Team {
	teamList := readFile(f)
	var t []Team
	for i := range teamList {
		log.Println(strconv.Itoa(teamList[i].Linenumber) + " " + teamList[i].Content)
		t = append(t, Team{FormatHTML.FormatTeamName(teamList[i].Content), teamList[i].Linenumber, teamList[i].Linenumber, 1, nil, 0})
	}
	return t
}

func (s Settings) LoadCFG(configPath string) (Settings, []Team) {
	cfg := readFile(configPath)
	for _, v := range cfg {
		if splitConfig(v.Content)[0] == "EnableWebConfig" {
			if splitConfig(v.Content)[1] == "TRUE" {
				s.webcfg = true
			} else {
				s.webcfg = false
			}
		}
		if splitConfig(v.Content)[0] == "TournamentName" {
			s.name = splitConfig(v.Content)[1]
		}
		if splitConfig(v.Content)[0] == "List" {
			s.list = splitConfig(v.Content)[1]
		}
		if splitConfig(v.Content)[0] == "GameMode" {

			i, err := strconv.Atoi(splitConfig(v.Content)[1])
			if err == nil {
				s.gameMode = i
			} else {
				log.Println("Wrong Value set in Gamemode")
			}
		}
		if splitConfig(v.Content)[0] == "GroupCount" {

			i, err := strconv.Atoi(splitConfig(v.Content)[1])
			if err == nil {
				s.groupCount = i - 1
			} else {
				log.Println("Wrong Value set in GroupCount")
			}
		}
		if splitConfig(v.Content)[0] == "AdvancingFromGroups" {
			i, err := strconv.Atoi(splitConfig(v.Content)[1])
			if err != nil {
				log.Fatalf("error setting Advancing Teams Count: %v", err)
			} else {
				s.aTCount = i
			}
		}
		if splitConfig(v.Content)[0] == "ConsoleLog" {
			if splitConfig(v.Content)[1] == "TRUE" {
				s.consoleLog = true
			} else {
				s.consoleLog = false
			}
		}
		if splitConfig(v.Content)[0] == "Port" {
			i, err := strconv.Atoi(splitConfig(v.Content)[1])
			if err != nil {
				log.Fatalf("error setting Port: %v", err)
			} else {
				s.port = i
			}

		}

	}

	log.Println(s)
	return s, readTeamList(s.list)
}

