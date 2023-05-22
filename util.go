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

func debug(val ...any) {
	for _, v := range val {
		switch v.(type) {
		case []byte:
			fmt.Printf("%s\n", v.([]byte))
			return
		case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
			fmt.Printf("%d\n", v.(int))
			return
		case float32, float64:
			fmt.Printf("%f\n", v.(float64))
			return
		case bool:
			fmt.Printf("%t\n", v.(bool))
			return
		default:
			fmt.Printf("%s\n", v)
			return
		}
	}
}

func contains(tgt string, src []byte) bool {
	for i := 0; i < len(src); i++ {
		if i+len(tgt) <= len(src) {
			if tgt == string(src[i:i+len(tgt)]) {
				return true
			}
		}
	}
	return false
}

func containsCount(tgt []string, src string, sep string) (bool, int) {
	b := []string{}
	for _, t := range tgt {
		for i := 0; i < len(src); i++ {
			if (i + len(t)) <= (len(src)) {
				if t == string(src[i:i+len(t)]) {
					b = append(b, t)
				}
			}
		}
	}
	if len(b) > 0 {
		return true, len(b)
	}
	return false, 0
}

func checkHorseName(name *Horse) {
	return
}
