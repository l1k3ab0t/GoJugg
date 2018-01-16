package Saveing

import (
	"os"
	"github.com/l1k3ab0t/GoJugg/lib/Objects"
	"encoding/xml"
	"io/ioutil"
	"log"
)
type Save struct{
	Games []Groups	`xml:"Group"`
	TeamList []Objects.Team
	S Objects.Settings	`xml:"Settings"`
}

type Groups struct {
	Id        int      `xml:"id,attr"`
	GroupGames []Round	`xml:"Round"`
}

type Round struct {
	Id        int      `xml:"id,attr"`
	GroupGames []Objects.Game	`xml:"Game"`
}


func Createave (filename string,t Objects.Tournament ) {
	var S Save
	S.TeamList=t.TeamList()
	for i,v:= range t.Games(){
		x:=[]Round{}
		for i2,v2:= range v{
			x=append(x,Round{i2,v2})
		}
		S.Games=append(S.Games, Groups{i,x})
	}
	output, err := xml.MarshalIndent(S, "  ", "    ")
	if err != nil {
		log.Printf("error: %v\n", err)
	}
	os.Stdout.Write([]byte(xml.Header))

	os.Stdout.Write(output)
	err = ioutil.WriteFile(filename, output, 0644)
	if err != nil {
		panic(err)
	}
}