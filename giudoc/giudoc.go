package giudoc

import (
	"strings"

	g "github.com/AllenDang/giu"
	"golang.org/x/net/html"
)

type Param struct {
	Name string
	Type string
	Desc string
}

type Requirement struct {
	Key   string
	Value string
}

type Link struct {
	Text string
	Href string
}

type Doc struct {
	Description  []Node
	Syntax       string
	Params       []Param
	ReturnType   string
	ReturnDesc   string
	Remarks      [][]Node
	Requirements []Requirement
	SeeAlso      []Link
}

type Node struct {
	Text string
	Link *Link
}

//////////////////////////////////////////////////////////
// PUBLIC ENTRY
//////////////////////////////////////////////////////////

func FromHTML(input string) g.Widget {
	doc := parseHTML(input)
	return render(doc)
}

//////////////////////////////////////////////////////////
// PARSER
//////////////////////////////////////////////////////////

func parseHTML(input string) Doc {
	root, _ := html.Parse(strings.NewReader(input))

	var d Doc
	var current string
	var lastParam *Param

	var walk func(*html.Node)
	walk = func(n *html.Node) {

		if n.Type == html.ElementNode {

			switch n.Data {

			case "h2":
				current = strings.ToLower(extractText(n))

			case "p":
				nodes := extractRichText(n)

				switch current {
				case "":
					if len(d.Description) == 0 {
						d.Description = nodes
					}

				case "parameters":
					text := flatten(nodes)

					if strings.HasPrefix(text, "Type:") {
						if lastParam != nil {
							lastParam.Type = strings.TrimPrefix(text, "Type: ")
						}
					} else if lastParam != nil {
						lastParam.Desc += text + " "
					}

				case "return value", "return-value":
					text := flatten(nodes)
					if strings.HasPrefix(text, "Type:") {
						d.ReturnType = strings.TrimPrefix(text, "Type: ")
					} else {
						d.ReturnDesc += text + " "
					}

				case "remarks":
					d.Remarks = append(d.Remarks, nodes)

				case "see also":
					// links handled separately
				}

			case "code":
				if current == "syntax" {
					d.Syntax = extractText(n)
				}
				if current == "parameters" {
					name := extractText(n)
					p := Param{Name: name}
					d.Params = append(d.Params, p)
					lastParam = &d.Params[len(d.Params)-1]
				}


			case "a":
				if current == "see also" {
					link := Link{
						Text: extractText(n),
						Href: getAttr(n, "href"),
					}
					d.SeeAlso = append(d.SeeAlso, link)
				}

			case "tr":
				if current == "requirements" {
					cells := extractCells(n)
					if len(cells) == 2 {
						d.Requirements = append(d.Requirements,
							Requirement{cells[0], cells[1]})
					}
				}
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			walk(c)
		}
	}

	walk(root)
	return d
}

//////////////////////////////////////////////////////////
// RENDERER
//////////////////////////////////////////////////////////

func render(d Doc) g.Widget {
	return g.Layout{

		renderRichText(d.Description),

		section("Syntax", codeBlock(d.Syntax)),

		paramSection(d.Params),

		returnSection(d.ReturnType, d.ReturnDesc),

		remarksSection(d.Remarks),

		requirementsTable(d.Requirements),

		linkSection(d.SeeAlso),
	}
}

//////////////////////////////////////////////////////////
// RICH TEXT (THIS IS THE IMPORTANT PART)
//////////////////////////////////////////////////////////

func renderRichText(nodes []Node) g.Widget {
	var widgets []g.Widget

	for _, n := range nodes {
		if n.Link != nil {
			link := n.Link

			// clickable text
			widgets = append(widgets,
				g.Button(link.Text).OnClick(func() {
					g.OpenURL(link.Href)
				}),
			)
		} else {
			widgets = append(widgets, g.Label(n.Text))
		}
	}

	return g.Row(widgets...)
}

//////////////////////////////////////////////////////////
// SECTIONS
//////////////////////////////////////////////////////////

func section(title string, w ...g.Widget) g.Widget {
	return g.Layout{
		g.Separator(),
		g.Label(title),
		g.Layout(w),
	}
}

func codeBlock(code string) g.Widget {
	return g.InputTextMultiline(&code).
		Size(-1, 120).
		Flags(g.InputTextFlagsReadOnly)
}

//////////////////////////////////////////////////////////
// PARAMETERS
//////////////////////////////////////////////////////////

func paramSection(params []Param) g.Widget {
	if len(params) == 0 {
		return g.Layout{}
	}

	var w []g.Widget

	for _, p := range params {
		w = append(w,
			g.Label(p.Name),
			g.Label("Type: "+p.Type),
			g.Label(p.Desc),
			g.Separator(),
		)
	}

	return section("Parameters", w...)
}

//////////////////////////////////////////////////////////
// RETURN
//////////////////////////////////////////////////////////

func returnSection(t, desc string) g.Widget {
	if t == "" && desc == "" {
		return g.Layout{}
	}

	return section("Return Value",
		g.Label("Type: "+t),
		g.Label(desc),
	)
}

//////////////////////////////////////////////////////////
// REMARKS
//////////////////////////////////////////////////////////

func remarksSection(r [][]Node) g.Widget {
	if len(r) == 0 {
		return g.Layout{}
	}

	var w []g.Widget

	for _, nodes := range r {
		w = append(w, renderRichText(nodes))
	}

	return section("Remarks", w...)
}

//////////////////////////////////////////////////////////
// REQUIREMENTS
//////////////////////////////////////////////////////////

func requirementsTable(reqs []Requirement) g.Widget {
	if len(reqs) == 0 {
		return g.Layout{}
	}

	var rows []*g.TableRowWidget

	for _, r := range reqs {
		rows = append(rows,
			g.TableRow(
				g.Label(r.Key),
				g.Label(r.Value),
			),
		)
	}

	return section("Requirements",
		g.Table().
			ID("req").
			Columns(
				g.TableColumn("Requirement"),
				g.TableColumn("Value"),
			).
			Rows(rows...),
	)
}

//////////////////////////////////////////////////////////
// LINKS (CLICKABLE)
//////////////////////////////////////////////////////////

func linkSection(links []Link) g.Widget {
	if len(links) == 0 {
		return g.Layout{}
	}

	var w []g.Widget

	for _, l := range links {
		link := l
		w = append(w,
			g.Button(link.Text).OnClick(func() {
				g.OpenURL(link.Href)
			}),
		)
	}

	return section("See Also", w...)
}

//////////////////////////////////////////////////////////
// HELPERS
//////////////////////////////////////////////////////////

func extractText(n *html.Node) string {
	var b strings.Builder

	var walk func(*html.Node)
	walk = func(node *html.Node) {
		if node.Type == html.TextNode {
			b.WriteString(node.Data)
		}
		for c := node.FirstChild; c != nil; c = c.NextSibling {
			walk(c)
		}
	}

	walk(n)
	return strings.TrimSpace(b.String())
}

func extractRichText(n *html.Node) []Node {
	var nodes []Node

	for c := n.FirstChild; c != nil; c = c.NextSibling {

		if c.Type == html.TextNode {
			nodes = append(nodes, Node{Text: c.Data})
		}

		if c.Type == html.ElementNode && c.Data == "a" {
			nodes = append(nodes, Node{
				Link: &Link{
					Text: extractText(c),
					Href: getAttr(c, "href"),
				},
			})
		}
	}

	return nodes
}

func extractCells(n *html.Node) []string {
	var cells []string

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode && (c.Data == "td" || c.Data == "th") {
			cells = append(cells, extractText(c))
		}
	}

	return cells
}

func flatten(nodes []Node) string {
	var b strings.Builder
	for _, n := range nodes {
		b.WriteString(n.Text)
	}
	return strings.TrimSpace(b.String())
}

func getAttr(n *html.Node, key string) string {
	for _, a := range n.Attr {
		if a.Key == key {
			return a.Val
		}
	}
	return ""
}
