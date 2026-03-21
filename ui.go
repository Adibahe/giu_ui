// Package main presents how to implement large, procedurally-generated table.
package main

import (
	"fmt"
	"image/color"

	g "github.com/AllenDang/giu"
)

func buildRows(messages []message) []*g.TableRowWidget {
	rows := make([]*g.TableRowWidget, len(messages))
	for i := range rows {
		rows[i] = g.TableRow(
			g.Label(fmt.Sprintf("%s", messages[i].Id)),
			g.Label(messages[i].Name),
		)
	}

	if len(rows) > 0 {
		rows[0].BgColor(&(color.RGBA{200, 100, 100, 255}))

	}

	return rows
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
		g.Table().Freeze(0, 1).FastMode(true).Columns(g.TableColumn("Id"), g.TableColumn("Function Name")).Rows(buildRows(messages)...),
	)

	return messages
}
