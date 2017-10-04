package FormatHTML

import (
	"github.com/l1k3ab0t/GoJugg/lib/GameEngine"
	"strconv"
	"html/template"
)

func FormatTeamLIst(teams []GameEngine.Team) template.HTML {
	formated:=template.HTML("<table>")
	formated = formated+"<tr> <th> Team ID </th> <th> Team Name </th> </tr>"
	for _,v:=range teams{
		formated = formated+ template.HTML("\n		<tr> <td>" + strconv.Itoa(v.ID) + "</td>" + "<td>" + v.Name + "</td> </tr> " )
	}
	formated = formated+"\n	 </table>"
	return formated
}

func FormatBracket(games []GameEngine.Game) template.HTML {
	formated:=template.HTML("<table>")
	formated = formated+"<tr> <th> Team 1 ID </th> <th> Team 1 </th> <th></th> <th> Team 2 ID </th> <th> Team 2 </th> <th>Result </th> </tr>"
	for _,v:=range games{
		formated = formated+ template.HTML("\n		<tr> <td>" + strconv.Itoa(v.Opponent1.ID) + "</td>" + "<td>" + v.Opponent1.Name + "</td> <td>	vs		</td> " )
		formated = formated+ template.HTML(" <td>" + strconv.Itoa(v.Opponent2.ID) + "</td>" + "<td>" + v.Opponent2.Name + "</td> <td>"+strconv.Itoa(v.Result.Team1Juggs)+" : "+strconv.Itoa(v.Result.Team2Juggs) + "</td> </tr> " )
	}
	formated = formated+"\n	 </table>"
	return formated
}