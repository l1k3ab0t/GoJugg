package Objects

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
)

type Settings struct {
	webcfg     bool
	name       string
	list       string
	fields     int
	gameMode   int
	groupCount int
	roundCount int
	consoleLog bool
	port       int
	aTCount    int
	customFieldNames bool
}

type line struct {
	Linenumber int
	Content    string
}

func (s Settings) Webcfg() bool {
	return s.webcfg
}

func (s Settings) Name() string {
	return s.name
}

func (s *Settings) SetName(Name string) {
	s.name = Name
}

func (s Settings) List() string {
	return s.list
}

func (s Settings) Fields() int {
	return s.fields
}

func (s *Settings) SetFields(Fields int) {
	s.fields = Fields
}

func (s Settings) GameMode() int {
	return s.gameMode
}

func (s *Settings) SetGameMode(GameMode int) {
	s.gameMode = GameMode
}

func (s Settings) GroupCount() int {
	return s.groupCount
}

func (s *Settings) SetGroupCount(GCount int) {
	s.groupCount = GCount
}

func (s Settings) RoundCount() int {
	return s.roundCount
}

func (s *Settings) SetRoundCount(RCount int) {
	s.roundCount = RCount
}

func (s Settings) ConsoleLog() bool {
	return s.consoleLog
}

func (s Settings) Port() int {
	return s.port
}

func (s Settings) ATCount() int {
	return s.aTCount
}

func (s Settings) CFieldNames() bool {
	return s.customFieldNames
}

func (s *Settings) SetCFieldNames(customFieldNames bool) {
	s.customFieldNames = customFieldNames
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

func ReadTeamList(f string) []Team {
	teamList := readFile(f)
	var t []Team
	for i := range teamList {
		log.Println(strconv.Itoa(teamList[i].Linenumber) + " " + teamList[i].Content)
		t = append(t, Team{formatTeamName(teamList[i].Content), teamList[i].Linenumber, teamList[i].Linenumber, 1, nil, 0})
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
		if splitConfig(v.Content)[0] == "RoundCount" {

			i, err := strconv.Atoi(splitConfig(v.Content)[1])
			if err == nil {
				s.roundCount = i - 1
			} else {
				log.Println("Wrong Value set in RoundCount")
			}
		}

		if splitConfig(v.Content)[0] == "Fields" {

			i, err := strconv.Atoi(splitConfig(v.Content)[1])
			if err == nil {
				s.fields = i - 1
			} else {
				log.Println("Wrong Value set in Fields")
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

		if splitConfig(v.Content)[0] == "CustomFieldNames" {
			if splitConfig(v.Content)[1] == "TRUE" {
				s.customFieldNames = true
			} else {
				s.customFieldNames = false
			}
		}


	}

	log.Println(s)
	return s, ReadTeamList(s.list)
}

func formatTeamName(name string) string {
	var rname string
	for _, char := range name {
		if char == 32 { //32 == " "
			rname = rname + "-"
		}else if char == 228{
			rname = rname + "ae"
		}else if char == 196{
			rname = rname + "Ae"
		}else if char == 246{
			rname = rname + "oe"
		}else if char == 214{
			rname = rname + "Oe"
		}else if char == 252{
			rname = rname + "ue"
		}else if char == 220{
			rname = rname + "Ue"
		} else {
			rname = rname + string(char)
		}
	}
	return rname
}