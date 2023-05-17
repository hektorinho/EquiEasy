package equieasy

import (
	"fmt"
	"math"
	"regexp"

	pdf "github.com/dslipak/pdf"
)

// Checks if page contains valid race data or is a continuation from previous race.
func isValidRacePage(page pdf.Page) bool {
	re := regexp.MustCompile(REGEX_VALID_RACE_PAGE)
	reCancelled := regexp.MustCompile(REGEX_VALID_CANCELLED)
	rows, err := page.GetTextByRow()
	if err != nil {
		fmt.Println("failed to receive rows back...")
		return false
	}

	firstrowdata := []byte{}
	secondrowdata := []byte{}
	for _, row := range rows {
		if row.Position == 760 {
			for _, word := range row.Content {
				firstrowdata = append(firstrowdata, word.S...)
				firstrowdata = append(firstrowdata, " "...)
			}
		}
		if row.Position == 751 {
			for _, word := range row.Content {
				secondrowdata = append(secondrowdata, word.S...)
				secondrowdata = append(secondrowdata, " "...)
			}
		}
	}
	if reCancelled.Match(secondrowdata) {
		return false
	}
	if re.Match(firstrowdata) {
		return true
	}
	return false
}

// Used to round a float to a certain precision
func roundFloat(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}
