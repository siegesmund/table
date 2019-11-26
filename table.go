package table

import (
	"encoding/json"
	"fmt"
	"github.com/gocolly/colly"
	"github.com/olekukonko/tablewriter"
	"os"
)

type Table struct {
	Headers []string
	Rows [][]string
}

func (t *Table) Get(url string, i int) error {
	c := colly.NewCollector()
	c.OnHTML("table", func(e *colly.HTMLElement) {
		if e.Index == i {
			t.getHeader(e)
			t.getRows(e)
		}
	})

	err := c.Visit(url)

	if err != nil {
		return err
	}

	c.Wait()

	return nil
}

func (t *Table) getHeader(h *colly.HTMLElement) {
	var headers []string
	h.ForEach("tr th", func(_ int, el *colly.HTMLElement) {
		headers = append(headers, el.Text)
	})
	t.Headers = headers
}

func (t *Table) getRows(rowElements *colly.HTMLElement) {
	var result [][]string
	rowElements.ForEach("tr", func(i int, rowElement *colly.HTMLElement) {
		var row [] string
		rowElement.ForEach("td", func(j int, columnElement *colly.HTMLElement) {
			row = append(row, columnElement.Text)
		})
		if len(row) > 0 {
			result = append(result, row)
		}
	})
	t.Rows = result
}

func (t *Table) Map() []map[string]string {
	var result []map[string]string
	for _, row := range t.Rows {
		var rowMap = make(map[string]string)
		for j, value:= range row {
			name := t.Headers[j]
			rowMap[name] = value
		}
		result = append(result, rowMap)
	}
	return result
}

func (t *Table) JSON() (string, error) {
	m := t.Map()
	b, e := json.Marshal(m)
	return string(b), e
}

func (t *Table) Print() {
	twriter := tablewriter.NewWriter(os.Stdout)
	twriter.SetHeader(t.Headers)
	for _, v := range t.Rows {
		twriter.Append(v)
	}
	twriter.Render()
}

func (t *Table) PrintJSON() error {
	m := t.Map()
	b, e := json.MarshalIndent(m, "", "\t")
	if e != nil {
		return e
	}
	fmt.Println(string(b))
	return nil
}


