package main

import (
	"fmt"
	"io"
	"strings"
	"text/tabwriter"
)

type response struct {
	columns []string
	result  [][]interface{}
}

func tabWriter(w io.Writer) *tabwriter.Writer {
	return tabwriter.NewWriter(w, 0, 0, 1, ' ', tabwriter.AlignRight|tabwriter.Debug)
}

func printResult(w *tabwriter.Writer, res *response) {
	header := ""
	spaceRow := ""
	columnLength := len(res.columns)

	for i := range res.columns {
		tab := " \t "
		if columnLength == i+1 {
			tab = " "
		}
		header += res.columns[i] + tab
		spaceRow += " " + tab
	}

	fmt.Fprintln(w, header)
	fmt.Fprintln(w, spaceRow)
	for _, row := range res.result {
		rowStr := ""
		for range row {
			rowStr += "%v \t "
		}

		rowStr = strings.TrimSuffix(rowStr, "\t ")
		fmt.Fprintln(w, fmt.Sprintf(rowStr, row...))
	}

	w.Flush()
}
