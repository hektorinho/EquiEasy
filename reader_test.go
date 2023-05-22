package equieasy

import (
	"regexp"
	"testing"

	pdf "github.com/dslipak/pdf"
)

const (
	// testFile = "data/eqbPDFChartPlus.pdf"

	// testFile = "C:/GoProjects/Scratch/Test/collected/BEL050910USA.pdf"
	// testFile = "C:/GoProjects/Scratch/Test/collected/BEL062202USA.pdf"
	// testFile = "C:/GoProjects/Scratch/Test/collected/BEL062409USA.pdf"
	// testFile = "C:/GoProjects/Scratch/Test/collected/BEL071807USA.pdf"
	// testFile = "C:/GoProjects/Scratch/Test/collected/BEL072094USA.pdf"
	// testFile = "C:/GoProjects/Scratch/Test/collected/BEL092996USA.pdf"
	// testFile = "C:/GoProjects/Scratch/Test/collected/BEL101812USA.pdf"
	// testFile = "C:/GoProjects/Scratch/Test/collected/BEL102012USA.pdf"
	// testFile = "C:/GoProjects/Scratch/DownloadFromEquibase/dl/AQU/AQU010115USA.pdf"
	// testFile = "C:/GoProjects/Scratch/DownloadFromEquibase/dl/AQU/AQU013191USA.pdf"
	// testFile = "C:/GoProjects/Scratch/DownloadFromEquibase/dl/AQU/AQU010292USA.pdf"
	testFile = "C:/GoProjects/Scratch/DownloadFromEquibase/dl/AQU/AQU010107USA.pdf"

	//
	// testFile = "C:/GoProjects/Scratch/DownloadFromEquibase/dl/AQU/AQU011608USA.pdf"

	// testFile = "C:/GoProjects/Scratch/Test/collected/BEL063022USA.pdf2"
	// testFile = "C:/GoProjects/Scratch/Test/collected/BEL070317USA.pdf2"
	// testFile = "C:/GoProjects/Scratch/Test/collected/BEL070806USA.pdf" // race 3 only 2 horses no betting...

	// testFile = "C:/GoProjects/Scratch/Test/collected/BEL051001USA2.pdf"
	// testFile = "C:/GoProjects/Scratch/Test/collected/BEL051118USA.pdf"
	// testFile = "C:/GoProjects/Scratch/Test/collected/BEL051514USA.pdf2"
	// testFile = "C:/GoProjects/Scratch/Test/collected/BEL060304USA.pdf"
)

// func TestContainsCount(t *testing.T) {
// 	str := []string{"---", "*"}
// 	rowdata := "5 High Chaparall * * * 4 56.00 fucking failed"
// 	ok, cnt := containsCount(str, rowdata, " ")
// 	// fmt.Println(ok, cnt)
// }

func TestHorses(t *testing.T) {
	r, err := pdf.Open(testFile)
	if err != nil {
		t.Errorf("err opening pdf >>> %s", err)
	}
	pages, err := GetValidPages(testFile, r)
	if err != nil {
		t.Errorf("wasnt able to get valid pages from %s\nerror >>> %s\n", testFile, err)
	}
	for i, page := range pages.Pages {
		// if i == 1 {
		if horses, err := Horses(page); err != nil {
			if len(horses) > 3 {
				t.Logf("Horse 1: %s, Horse 2: %s, Horse 3: %s\n", horses[0].Name, horses[1].Name, horses[2].Name)
			} else {
				t.Errorf("expected to have more horses in a race, got %d on page %d\n%v\n", len(horses), i+1, horses)
			}
			if len(horses) < 1 {
				t.Errorf("expected more than one horse per race, got %d\n", len(horses))
			}
		}
		// }
	}
}

func TestGetValidPages(t *testing.T) {
	re := regexp.MustCompile(REGEX_VALID_RACE_PAGE)
	r, err := pdf.Open(testFile)
	if err != nil {
		t.Errorf("err opening pdf >>> %s", err)
	}
	pages, err := GetValidPages(testFile, r)
	if err != nil {
		t.Errorf("wasnt able to get valid pages from %s\nerror >>> %s\n", testFile, err)
	}
	counter := 0
	for _, page := range pages.Pages {
		rows, err := page.GetTextByRow()
		if err != nil {
			t.Errorf("wasnt able to get valid pages from %s\nerror >>> %s\n", testFile, err)
		}
		rowdata := []byte{}
		for _, row := range rows {
			if row.Position == 760 {
				for _, word := range row.Content {
					rowdata = append(rowdata, word.S...)
					rowdata = append(rowdata, " "...)
				}

				if re.Match(rowdata) {
					counter++
				}
			}
		}
	}
	if len(pages.Pages) != counter {
		t.Errorf("wasnt able to get valid pages from %s\nGot len(pages)=%d, counted to %d valid pages\n", testFile, len(pages.Pages), counter)
	}
}

func TestFractionals(t *testing.T) {
	r, err := pdf.Open(testFile)
	if err != nil {
		t.Errorf("err opening pdf >>> %s", err)
	}
	pages, err := GetValidPages(testFile, r)
	if err != nil {
		t.Errorf("wasnt able to get valid pages from %s\nerror >>> %s\n", testFile, err)
	}
	for i, page := range pages.Pages {
		// if i == 3 {
		fracs, _, err := fractionals(page)
		if err != nil {
			t.Errorf("page >> %d || error >>> %s\n", i+1, err)
		}
		if len(fracs) < 1 {
			t.Errorf("expected to get more than 0 fracs, got %d\n", len(fracs))
		}
		// }
	}
}

func TestApplyFractionals(t *testing.T) {
	r, err := pdf.Open(testFile)
	if err != nil {
		t.Errorf("err opening pdf >>> %s", err)
	}
	pages, err := GetValidPages(testFile, r)
	if err != nil {
		t.Errorf("wasnt able to get valid pages from %s\nerror >>> %s\n", testFile, err)
	}

	for i, page := range pages.Pages {
		// fmt.Println(i)
		if i == 6 {
			horses, err := Horses(page)
			if err != nil {
				t.Errorf("wasnt able to get horse data for page %d\n", i+1)
			}
			for _, horse := range horses {
				checkHorseName(horse)
				// if err := horse.applyFractionalData(page); err != nil {
				// 	t.Errorf("wasnt able to asssemble fractional data points for page %d and horse %s\n", i+1, horse.Name)
				// }
				// doesn't always come with full fractional data
				// if len(horse.Fractionals) < 1 {
				// 	t.Errorf("expected to have fractional data, got len(horse.Fractionals)=%d\n", len(horse.Fractionals))
				// }
			}
		}
	}
}

func TestApplyTrainerData(t *testing.T) {
	r, err := pdf.Open(testFile)
	if err != nil {
		t.Errorf("err opening pdf >>> %s", err)
	}
	pages, err := GetValidPages(testFile, r)
	if err != nil {
		t.Errorf("wasnt able to get valid pages from %s\nerror >>> %s\n", testFile, err)
	}

	for i, page := range pages.Pages {
		// if i == 2 {
		horses, err := Horses(page)
		if err != nil {
			t.Errorf("wasnt able to get horse data for page %d\n", i+1)
		}
		for _, horse := range horses {
			// if err := horse.applyFractionalData(page); err != nil {
			// 	t.Errorf("wasnt able to asssemble fractional data points for page %d and horse %s\n", i+1, horse.Name)
			// }
			if err := horse.applyTrainerData(page); err != nil {
				t.Errorf("wasnt able to asssemble trainer for page %d and horse %s\n", i+1, horse.Name)
			}
			if len(horse.Trainer) < 1 {
				t.Errorf("expected to have trainer data, got trainer=%s\n for horse=%s", horse.Trainer, horse.Name)
			}
		}
		// }
	}
}

func TestApplyOwnerData(t *testing.T) {
	r, err := pdf.Open(testFile)
	if err != nil {
		t.Errorf("err opening pdf >>> %s", err)
	}
	pages, err := GetValidPages(testFile, r)
	if err != nil {
		t.Errorf("wasnt able to get valid pages from %s\nerror >>> %s\n", testFile, err)
	}

	for i, page := range pages.Pages {
		// if i == 1 {
		horses, err := Horses(page)
		if err != nil {
			t.Errorf("wasnt able to get horse data for page %d\n", i+1)
		}
		for _, horse := range horses {
			// if err := horse.applyFractionalData(page); err != nil {
			// 	t.Errorf("wasnt able to asssemble fractional data points for page %d and horse %s\n", i+1, horse.Name)
			// }
			if err := horse.applyOwnerData(page); err != nil {
				t.Errorf("wasnt able to asssemble trainer for page %d and horse %s\n", i+1, horse.Name)
			}
			if len(horse.Owners) < 1 {
				t.Errorf("page %d | expected to have owner data, got owner=%s\n from horse=%s", i, horse.Owners, horse.Name)
			}
		}
		// }
	}
}

func TestTrackAndDate(t *testing.T) {
	r, err := pdf.Open(testFile)
	if err != nil {
		t.Errorf("err opening pdf >>> %s", err)
	}
	pages, err := GetValidPages(testFile, r)
	if err != nil {
		t.Errorf("wasnt able to get valid pages from %s\nerror >>> %s\n", testFile, err)
	}

	for i, page := range pages.Pages {
		// if i == 5 {
		race, err := Metadata(page)
		if err != nil {
			t.Errorf("wasn't able to gather track and date data from page %d\nerror >> %s\n", i+1, err)
		}
		if race.Date.IsZero() {
			t.Errorf("wasn't able to gather date data from page, got date=%s\n", race.Date)
		}
		if race.Track == "" {
			t.Errorf("wasn't able to gather track data from page, got track=%s\n", race.Track)
		}
		if race.Number == 0 {
			t.Errorf("wasn't gather the race number. got number=%d", race.Number)
		}
		// }
	}
}

func TestGenericDataFromRegex(t *testing.T) {
	r, err := pdf.Open(testFile)
	if err != nil {
		t.Errorf("err opening pdf >>> %s", err)
	}
	pages, err := GetValidPages(testFile, r)
	if err != nil {
		t.Errorf("wasnt able to get valid pages from %s\nerror >>> %s\n", testFile, err)
	}

	for i, page := range pages.Pages {
		// if i == 8 {
		meta := RaceMetadata{}
		if horsetype, err := meta.genericDataFromRegex(page, REGEX_RACE_HORSETYPE); err != nil || horsetype == nil {
			t.Errorf("wasn't able to gather horse type data from page %d, %s\nerror >> %s", i+1, horsetype, err)
		}
		if racetype, err := meta.genericDataFromRegex(page, REGEX_RACE_RACETYPE); err != nil || racetype == nil {
			t.Errorf("wasn't able to gather RACETYPE data from page %d, %s\nerror >> %s", i+1, racetype, err)
		}
		if purse, err := meta.genericDataFromRegex(page, REGEX_RACE_PURSE); err != nil || purse == nil {
			t.Errorf("wasn't able to gather PURSE data from page %d, %s\nerror >> %s", i+1, purse, err)
		}
		///  Weather data not always included
		// if weather, err := meta.genericDataFromRegex(page, REGEX_RACE_WEATHER); err != nil || weather == nil {
		// 	t.Errorf("wasn't able to gather WEATHER data from page %d, %s\nerror >> %s", i+1, weather, err)
		// }
		if trackcondition, err := meta.genericDataFromRegex(page, REGEX_RACE_TRACK_CONDITION); err != nil || trackcondition == nil {
			t.Errorf("wasn't able to gather TRACK_CONDITION data from page %d, %s\nerror >> %s", i+1, trackcondition, err)
		}
		if length, err := meta.genericDataFromRegex(page, REGEX_RACE_LENGTH); err != nil || length == nil {
			t.Errorf("wasn't able to gather RACE_LENGTH data from page %d, %s\nerror >> %s", i+1, length, err)
		}
		if trackrecord, err := meta.genericDataFromRegex(page, REGEX_RACE_CURRENT_TRACK_RECORD); err != nil || trackrecord == nil {
			t.Errorf("wasn't able to gather CURRENT_TRACK_RECORD data from page %d, %s\nerror >> %s", i+1, trackrecord, err)
		}
		if finaltime, err := meta.genericDataFromRegex(page, REGEX_RACE_FINAL_TIME); err != nil || finaltime == nil {
			t.Errorf("wasn't able to gather RACE_FINAL_TIME data from page %d, %s\nerror >> %s", i+1, finaltime, err)
		}

		// Fractional times can sometimes be omitted...
		// if fractionaltimes, err := meta.genericDataFromRegex(page, REGEX_RACE_FRACTIONAL_TIMES); err != nil || fractionaltimes == nil {
		// 	t.Errorf("wasn't able to gather FRACTIONAL_TIMES data from page %d, %s\nerror >> %s", i+1, fractionaltimes, err)
		// }

		// Split times can sometimes be nil if clock malfunction during race ...
		// if splittimes, err := meta.genericDataFromRegex(page, REGEX_RACE_SPLIT_TIMES); err != nil || splittimes == nil {
		// 	t.Errorf("wasn't able to gather SPLIT_TIMES data from page %d, %s\nerror >> %s", i+1, splittimes, err)
		// }
		if runup, err := meta.genericDataFromRegex(page, REGEX_RACE_RUN_UP); err != nil || runup == nil {
			t.Errorf("wasn't able to gather RUN_UP data from page %d, %s\nerror >> %s", i+1, runup, err)
		}
		// }
	}
}

func TestNewRacePage(t *testing.T) {
	r, err := pdf.Open(testFile)
	if err != nil {
		t.Errorf("err opening pdf >>> %s", err)
	}
	pages, err := GetValidPages(testFile, r)
	if err != nil {
		t.Errorf("wasnt able to get valid pages from %s\nerror >>> %s\n", testFile, err)
	}

	for i, page := range pages.Pages {
		_, err := NewRacePage(page)
		if err != nil {
			t.Errorf("wasnt able to construct a race page from page=%d\nerror >>> %s\n", i+1, err)
		}
	}
}
