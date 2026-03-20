package main

import (
	_ "embed"
	"fmt"

	g "github.com/AllenDang/giu"
)

//go:embed style.css
var cssStyle []byte

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
	tableRows := make([]*g.TableRowWidget, 0, len(messages))
	for _, msg := range messages {
		tableRows = append(tableRows,
			g.TableRow(
				g.Label(msg.Id),
				g.Label(msg.Name),
			),
		)
	}

	rightPanelWidgets := []g.Widget{
		g.Label("Function Calls"),
		g.Separator(),
	}

	if len(messages) == 0 {
		rightPanelWidgets = append(rightPanelWidgets, g.Label("No function calls yet"))
	} else {
		for i, msg := range messages {
			rightPanelWidgets = append(
				rightPanelWidgets,
				g.Label(fmt.Sprintf("%d. %s()", i+1, msg.Name)),
			)
		}
	}

	g.SingleWindow().Layout(
		g.Label("Incoming Messages"),
		g.Separator(),

		g.Row(
			g.Child().Size(450, 500).Layout(
				g.Label("Messages (ID / Name)"),
				g.Separator(),
				g.Table().
					Columns(
						g.TableColumn("ID"),
						g.TableColumn("Name"),
					).
					Rows(tableRows...),
			),

			// Right side: function call list
			g.Child().Size(300, 500).Layout(
				rightPanelWidgets...,
			),
		),
	)

	return messages
}
