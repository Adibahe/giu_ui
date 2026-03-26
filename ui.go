package main

import (
	"fmt"
	"image/color"
	"strconv"
	"strings"

	"github.com/AllenDang/giu"
	g "github.com/AllenDang/giu"
)

var sashPos1 float32 = 600
var sashPos2 float32 = 500
var selectedFunction int = -1
var details string = "(⌐■_■) click on function calls"
var text = ""

func buildRows(messages []message) []*g.TableRowWidget {
	rows := make([]*g.TableRowWidget, len(messages))
	if len(messages) == 0 {
		return []*g.TableRowWidget{
			g.TableRow(g.Label("	( ╥﹏╥) 	Nothing yet!!")),
		}
	} else {
		for i := range rows {
			if messages[i].Name == "" {
				messages[i].Name = getName(messages[i].Id)
			}
			rows[i] = g.TableRow(
				g.Label(strconv.Itoa(i+1)),
				g.Selectable(messages[i].Id).
					Selected(selectedFunction == i).
					OnClick(func() {
						selectedFunction = i
						go getDesc(messages[i].Id, &details)
						fmt.Println("Selected:", messages[i].Name)
					}),
				g.Selectable(messages[i].Name).
					Selected(selectedFunction == i).
					OnClick(func() {
						selectedFunction = i
						go getDesc(messages[i].Id, &details)
						fmt.Println("Selected:", messages[i].Name)
					}),
			)
			if selectedFunction == i {
				rows[i].BgColor(&color.RGBA{80, 120, 200, 100})
			}

			if matchSearch(messages[i].Name, text) {
				rows[i].BgColor(&color.RGBA{255, 220, 100, 120})
			}

		}
	}

	return rows
}

func matchSearch(name string, text string) bool {
	if text == "" {
		return false
	}
	return strings.Contains(strings.ToLower(name), strings.ToLower(text))
}

func loop(msgchan chan message, messages []message) []message {

	for {
		select {
		case msg := <-msgchan:
			messages = append(messages, msg)
		default:
			goto render
		}
	}

render:

	g.SingleWindow().Layout(
		g.SplitLayout(
			g.DirectionVertical,
			&sashPos2,
			g.Layout{
				g.Style().
					SetFontSize(15).
					To(
						g.Label("Function List"),
					),
				g.Separator(),
				// g.SplitLayout(
				// 	g.DirectionHorizontal,
				// 	&sashPos1,
				// 	g.Layout{
				// 		g.Table().Size(-1, -1).Freeze(0, 1).FastMode(true).Columns(g.TableColumn("S/N."), g.TableColumn("Id"), g.TableColumn("Function Name")).Rows(buildRows(messages)...),
				// 	},
				// 	g.Layout{
				// 		g.Separator(),
				// 		g.Row(g.Label("Search"), g.InputText(&text).Size(-1).AutoComplete(keys).Size(300)),
				// 	},
				// ),
				g.Layout{
					g.Table().Flags(g.TableFlagsSizingFixedFit).Size(-1, 900).Freeze(0, 1).FastMode(true).Columns(g.TableColumn("S/N."), g.TableColumn("Id").InnerWidthOrWeight(40), g.TableColumn("Function Name")).Rows(buildRows(messages)...),
				},
				g.Layout{
					g.Separator(),
					g.Row(g.Label("Search"), g.InputText(&text).Size(-1)),
				},
			},
			g.Layout{
				g.Style().
					SetFontSize(15).
					To(
						g.Label("Description"),
					),
				g.Separator(),
				g.Style().SetFontSize(15).SetFont(monoFont).
					To(
						g.InputTextMultiline(&details).Flags(giu.InputTextFlagsReadOnly).Size(-1, -1),
					),
			},
		),
	)

	return messages
}
