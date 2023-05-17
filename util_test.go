package equieasy

import (
	"testing"

	"github.com/dslipak/pdf"
)

func TestIsValidRacePage(t *testing.T) {
	r, err := pdf.Open(testFile)
	if err != nil {
		t.Errorf("unable to open test file: %s", testFile)
	}

	listofPages := []pdf.Page{r.Page(1), r.Page(2), r.Page(3), r.Page(4)}

	if testFile == "data/eqbPDFChartPlus.pdf" {
		for i, page := range listofPages {
			if isValidRacePage(page) {
				switch i {
				case 0, 1, 3:
					continue
				case 2:
					t.Errorf("expected this page not to be valid, page=%d", i+1)
				default:
					t.Errorf("expected this page not to be valid, page=%d", i+1)
				}
			} else {
				switch i {
				case 0, 1, 3:
					t.Errorf("expected this page to be valid, page=%d", i+1)
				case 2:
					continue
				default:
					t.Errorf("expected this page to be valid, page=%d", i+1)
				}
			}
		}
	}
}

func TestRoundFloat(t *testing.T) {
	listofValues := []float64{3.3333333333, 4.55678, 2.1, 3.456, 0.1 + 0.2}
	expectedValues := []float64{3.333, 4.557, 2.100, 3.456, 0.3}
	for i, val := range listofValues {
		rounded := roundFloat(val, 3)
		if rounded != expectedValues[i] {
			t.Errorf("failed to round the float64 value | expected >> %f got >> %f", expectedValues[i], rounded)
		} else {
			continue
		}
	}
}
