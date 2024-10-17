package cli

import (
	"fmt"
	"github.com/artarts36/singlecli/color"
	"strings"
)

type Output interface {
	PrintMarkdownTable(headers []string, rows [][]string)
	PrintColoredBlock(color color.ConsoleColor, text string)
}

type output struct {
}

func (output) PrintMarkdownTable(headers []string, rows [][]string) {
	if len(rows) == 0 {
		return
	}

	colOffsets := make([]int, len(headers))

	for i, header := range headers {
		colOffsets[i] = len(header)
	}

	for _, row := range rows {
		for colID, col := range row {
			if len(col) > colOffsets[colID] {
				colOffsets[colID] = len(col)
			}
		}
	}

	headerString := make([]string, 0, len(headers)*3-2)
	separatorString := make([]string, 0, len(headers)*2-1)

	for i, header := range headers {
		headerString = append(headerString, header)

		if i < len(headers)-1 {
			headerString = append(headerString, strings.Repeat(
				" ",
				colOffsets[i]-len(header)+2,
			), "| ")
		}

		// fill separator line

		separatorString = append(separatorString, strings.Repeat("-", colOffsets[i]+2))

		if i < len(headers)-1 {
			separatorString = append(separatorString, "|")
		}
	}

	fmt.Println(strings.Join(headerString, ""))
	fmt.Println(strings.Join(separatorString, ""))

	for _, row := range rows {
		rowString := make([]string, 0, len(row)*3-1)

		for colID, col := range row {
			rowString = append(rowString, col)

			if colID < len(row)-1 {
				rowString = append(
					rowString,
					strings.Repeat(" ", colOffsets[colID]-len(col)+2),
					"| ",
				)
			}
		}

		fmt.Println(strings.Join(rowString, ""))
	}
}

func (output) PrintColoredBlock(col color.ConsoleColor, text string) {
	fmt.Printf("\x1b[4%dm %s \x1b[0m\n", col, text)
}
