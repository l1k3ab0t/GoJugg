package Objects

import (
	"encoding/xml"
	"io/ioutil"
	"log"
	"os"
)

type Save struct {
	Games    []Groups `xml:"Group"`
	TeamList []Team
	S        SaveSettings `xml:"Settings"`
}

type Groups struct {
	Id         int     `xml:"id,attr"`
	GroupGames []Round `xml:"Round"`
}

type Round struct {
	Id         int    `xml:"id,attr"`
	GroupGames []Game `xml:"Game"`
}

type SaveSettings struct {
	Webcfg           bool   `xml:"WegCFG"`
	Name             string `xml:"TName"`
	List             string `xml:"TeamListPath"`
	Fields           int    `xml:"NofFields"`
	GameMode         int    `xml:"GameMode"`
	GroupCount       int    `xml:"GroupCount"`
	RoundCount       int    `xml:"RoundCount"`
	ConsoleLog       bool   `xml:"ConsoleLog"`
	Port             int    `xml:"Port"`
	ATCount          int
	CustomFieldNames bool `xml:"CustomTeamNames"`
}

func CreateSave(filename string, t Tournament) {
	if filename == "" {
		filename = "Save.xml"
	}
	var S Save
	S.TeamList = t.TeamList()
	for i, v := range t.Games() {
		x := []Round{}
		for i2, v2 := range v {
			x = append(x, Round{i2, v2})
		}
		S.Games = append(S.Games, Groups{i, x})
	}
	S.S.RoundCount = t.S.RoundCount()
	S.S.GroupCount = t.S.GroupCount()
	S.S.List = t.S.List()
	S.S.Fields = t.S.Fields()
	S.S.GameMode = t.S.GameMode()
	S.S.Name = t.S.Name()
	S.S.ATCount = t.S.ATCount()
	S.S.ConsoleLog = t.S.ConsoleLog()
	S.S.CustomFieldNames = t.S.webcfg
	S.S.Port = t.S.Port()
	S.S.Webcfg = t.S.Webcfg()
	output, err := xml.MarshalIndent(S, "  ", "    ")
	if err != nil {
		log.Printf("error: %v\n", err)
	}
	os.Stdout.Write([]byte(xml.Header))
	os.Stdout.Write(output)
	err = ioutil.WriteFile(filename, output, 0644)
	if err != nil {
		panic(err)
	} else {
		log.Println("Succsessful saved, ", filename)
	}
}
