package equieasy

import (
	"math"
	"regexp"

	pdf "github.com/dslipak/pdf"
)

// Checks if page contains valid race data or is a continuation from previous race.
func IsValidRacePage(page pdf.Page) bool {
	re := regexp.MustCompile(REGEX_VALID_RACE_PAGE)
	rows, err := page.GetTextByRow()
	if err != nil {
		return false
	}

	rowdata := []byte{}
	for _, row := range rows {
		if row.Position == 760 {
			for _, word := range row.Content {
				rowdata = append(rowdata, word.S...)
				rowdata = append(rowdata, " "...)
			}
		}
		if re.Match(rowdata) {
			return true
		}
	}
	return false
}

// Used to round a float to a certain precision
func roundFloat(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}
