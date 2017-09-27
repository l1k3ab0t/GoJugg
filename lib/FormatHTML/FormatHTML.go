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
