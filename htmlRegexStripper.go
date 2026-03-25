package main

import (
	// "regexp"
	// "strings"

	"jaytaylor.com/html2text"
)

func htmlToReadableText(s string) string {
	op := html2text.Options{
		PrettyTables:        true,
		PrettyTablesOptions: html2text.NewPrettyTablesOptions(),
	}
	_, err := html2text.FromString(s, op)
	if err != nil {
		panic(err)
	}
	return s
}

// func htmlToReadableText(s string) string {
// 	// Normalize line endings
// 	s = strings.ReplaceAll(s, "\r\n", "\n")
//
// 	// Convert opening tags with attributes to readable markers
// 	reTagReplacements := []struct {
// 		re  *regexp.Regexp
// 		new string
// 	}{
// 		// Headings
// 		{regexp.MustCompile(`(?i)<h1[^>]*>`), "\n\n"},
// 		{regexp.MustCompile(`(?i)<h2[^>]*>`), "\n\n"},
// 		{regexp.MustCompile(`(?i)<h3[^>]*>`), "\n\n"},
// 		{regexp.MustCompile(`(?i)<h4[^>]*>`), "\n\n"},
// 		{regexp.MustCompile(`(?i)<h5[^>]*>`), "\n\n"},
// 		{regexp.MustCompile(`(?i)<h6[^>]*>`), "\n\n"},
//
// 		// Paragraphs / divs / sections
// 		{regexp.MustCompile(`(?i)<p[^>]*>`), ""},
// 		{regexp.MustCompile(`(?i)<div[^>]*>`), ""},
// 		{regexp.MustCompile(`(?i)<section[^>]*>`), "\n\n"},
//
// 		// Lists
// 		{regexp.MustCompile(`(?i)<ul[^>]*>`), "\n"},
// 		{regexp.MustCompile(`(?i)<ol[^>]*>`), "\n"},
// 		{regexp.MustCompile(`(?i)<li[^>]*>`), "• "},
//
// 		// Code blocks
// 		{regexp.MustCompile(`(?i)<pre[^>]*>`), "\n\n"},
// 		{regexp.MustCompile(`(?i)<code[^>]*>`), ""},
//
// 		// Tables
// 		{regexp.MustCompile(`(?i)<table[^>]*>`), "\n"},
// 		{regexp.MustCompile(`(?i)<thead[^>]*>`), ""},
// 		{regexp.MustCompile(`(?i)<tbody[^>]*>`), ""},
// 		{regexp.MustCompile(`(?i)<tr[^>]*>`), ""},
// 		{regexp.MustCompile(`(?i)<th[^>]*>`), ""},
// 		{regexp.MustCompile(`(?i)<td[^>]*>`), ""},
//
// 		// Links
// 		{regexp.MustCompile(`(?i)<a[^>]*>`), ""},
// 		{regexp.MustCompile(`(?i)<strong[^>]*>`), ""},
// 		{regexp.MustCompile(`(?i)<b[^>]*>`), ""},
// 		{regexp.MustCompile(`(?i)<em[^>]*>`), ""},
// 		{regexp.MustCompile(`(?i)<i[^>]*>`), ""},
// 	}
//
// 	for _, rr := range reTagReplacements {
// 		s = rr.re.ReplaceAllString(s, rr.new)
// 	}
//
// 	// Exact closing tag replacements
// 	replacements := []struct {
// 		old string
// 		new string
// 	}{
// 		// Line breaks
// 		{"<br>", "\n"},
// 		{"<br/>", "\n"},
// 		{"<br />", "\n"},
//
// 		// Paragraphs / block sections
// 		{"</p>", "\n\n"},
// 		{"</div>", "\n"},
// 		{"</section>", "\n\n"},
//
// 		// Lists
// 		{"</ul>", "\n"},
// 		{"</ol>", "\n"},
// 		{"</li>", "\n"},
//
// 		// Headings
// 		{"</h1>", "\n\n"},
// 		{"</h2>", "\n\n"},
// 		{"</h3>", "\n\n"},
// 		{"</h4>", "\n\n"},
// 		{"</h5>", "\n\n"},
// 		{"</h6>", "\n\n"},
//
// 		// Code blocks
// 		{"</pre>", "\n\n"},
// 		{"</code>", ""},
//
// 		// Tables
// 		{"</table>", "\n"},
// 		{"</thead>", ""},
// 		{"</tbody>", ""},
// 		{"</tr>", "\n"},
// 		{"</th>", " : "},
// 		{"</td>", "    "},
//
// 		// Links / emphasis
// 		{"</a>", ""},
// 		{"</strong>", ""},
// 		{"</b>", ""},
// 		{"</em>", ""},
// 		{"</i>", ""},
//
// 		// Rules
// 		{"<hr>", "\n----------------------------------------\n"},
// 		{"<hr/>", "\n----------------------------------------\n"},
// 		{"<hr />", "\n----------------------------------------\n"},
// 	}
//
// 	for _, r := range replacements {
// 		s = strings.ReplaceAll(s, r.old, r.new)
// 	}
//
// 	re := regexp.MustCompile(`(?is)<[^>]+>`)
// 	s = re.ReplaceAllString(s, "")
//
// 	entityReplacements := []struct {
// 		old string
// 		new string
// 	}{
// 		{"&nbsp;", " "},
// 		{"&amp;", "&"},
// 		{"&lt;", "<"},
// 		{"&gt;", ">"},
// 		{"&quot;", `"`},
// 		{"&#39;", "'"},
// 		{"&apos;", "'"},
// 		{"┬á", " "},
// 	}
//
// 	for _, e := range entityReplacements {
// 		s = strings.ReplaceAll(s, e.old, e.new)
// 	}
//
// 	s = regexp.MustCompile(`[ \t]+\n`).ReplaceAllString(s, "\n")
// 	s = regexp.MustCompile(`\n{3,}`).ReplaceAllString(s, "\n\n")
//
// 	return strings.TrimSpace(s)
// }

//lesser version

// func htmlToReadableText(s string) string {
// 	s = strings.ReplaceAll(s, "<br>", "\n")
// 	s = strings.ReplaceAll(s, "<br/>", "\n")
// 	s = strings.ReplaceAll(s, "<br />", "\n")
// 	s = strings.ReplaceAll(s, "</p>", "\n\n")
// 	s = strings.ReplaceAll(s, "</div>", "\n")

// 	re := regexp.MustCompile(`<[^>]*>`)
// 	s = re.ReplaceAllString(s, "")

// 	s = strings.ReplaceAll(s, "&nbsp;", " ")
// 	s = strings.ReplaceAll(s, "&amp;", "&")
// 	s = strings.ReplaceAll(s, "&lt;", "<")
// 	s = strings.ReplaceAll(s, "&gt;", ">")

// 	return strings.TrimSpace(s)
// }
