package equieasy

import (
	"regexp"
	"testing"
)

const (
	// testFile = "data/eqbPDFChartPlus.pdf"
	testFile = "C:/GoProjects/Scratch/Test/collected/BEL050910USA.pdf"
)

func TestGetHorses(t *testing.T) {
	pages, err := GetValidPages(testFile)
	if err != nil {
		t.Errorf("wasnt able to get valid pages from %s\nerror >>> %s\n", testFile, err)
	}
	for i, page := range pages.Pages {
		if horses, err := GetHorses(page); err != nil {
			if len(horses) > 3 {
				t.Logf("Horse 1: %s, Horse 2: %s, Horse 3: %s\n", horses[0].Name, horses[1].Name, horses[2].Name)
			} else {
				t.Errorf("expected to have more horses in a race, got %d on page %d\n%v\n", len(horses), i+1, horses)
			}
			if len(horses) < 1 {
				t.Errorf("expected more than one horse per race, got %d\n", len(horses))
			}
		}
	}
}

func TestGetValidPages(t *testing.T) {
	re := regexp.MustCompile(REGEX_VALID_RACE_PAGE)
	pages, err := GetValidPages(testFile)
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
	pages, err := GetValidPages(testFile)
	if err != nil {
		t.Errorf("wasnt able to get valid pages from %s\nerror >>> %s\n", testFile, err)
	}
	for _, page := range pages.Pages {
		fracs, _, err := Fractionals(page)
		if err != nil {
			t.Errorf("error >>> %s\n", err)
		}
		if len(fracs) < 1 {
			t.Errorf("expected to get more than 0 fracs, got %d\n", len(fracs))
		}
	}
}

func TestApplyFractionals(t *testing.T) {
	pages, err := GetValidPages(testFile)
	if err != nil {
		t.Errorf("wasnt able to get valid pages from %s\nerror >>> %s\n", testFile, err)
	}

	for i, page := range pages.Pages {
		// if i == 5 {
		horses, err := GetHorses(page)
		if err != nil {
			t.Errorf("wasnt able to get horse data for page %d\n", i+1)
		}
		for _, horse := range horses {
			if err := horse.ApplyFractionalData(page); err != nil {
				t.Errorf("wasnt able to asssemble fractional data points for page %d and horse %s\n", i+1, horse.Name)
			}
			if len(horse.Fractionals) < 1 {
				t.Errorf("expected to have fractional data, got len(horse.Fractionals)=%d\n", len(horse.Fractionals))
			}
		}
		// }
	}
}

func TestApplyTrainerData(t *testing.T) {
	pages, err := GetValidPages(testFile)
	if err != nil {
		t.Errorf("wasnt able to get valid pages from %s\nerror >>> %s\n", testFile, err)
	}

	for i, page := range pages.Pages {
		// if i == 5 {
		horses, err := GetHorses(page)
		if err != nil {
			t.Errorf("wasnt able to get horse data for page %d\n", i+1)
		}
		for _, horse := range horses {
			if err := horse.ApplyFractionalData(page); err != nil {
				t.Errorf("wasnt able to asssemble fractional data points for page %d and horse %s\n", i+1, horse.Name)
			}
			if err := horse.ApplyTrainerData(page); err != nil {
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
	pages, err := GetValidPages(testFile)
	if err != nil {
		t.Errorf("wasnt able to get valid pages from %s\nerror >>> %s\n", testFile, err)
	}

	for i, page := range pages.Pages {
		// if i == 5 {
		horses, err := GetHorses(page)
		if err != nil {
			t.Errorf("wasnt able to get horse data for page %d\n", i+1)
		}
		for _, horse := range horses {
			if err := horse.ApplyFractionalData(page); err != nil {
				t.Errorf("wasnt able to asssemble fractional data points for page %d and horse %s\n", i+1, horse.Name)
			}
			if err := horse.ApplyOwnerData(page); err != nil {
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
	pages, err := GetValidPages(testFile)
	if err != nil {
		t.Errorf("wasnt able to get valid pages from %s\nerror >>> %s\n", testFile, err)
	}

	for i, page := range pages.Pages {
		race, err := NewMetadata(page)
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
	}
}

func TestGenericDataFromRegex(t *testing.T) {
	pages, err := GetValidPages(testFile)
	if err != nil {
		t.Errorf("wasnt able to get valid pages from %s\nerror >>> %s\n", testFile, err)
	}

	for i, page := range pages.Pages {
		// if i == 8 {
		meta := RaceMetadata{}
		if horsetype, err := meta.GenericDataFromRegex(page, REGEX_RACE_HORSETYPE); err != nil || horsetype == nil {
			t.Errorf("wasn't able to gather generic data from page %d, %s\nerror >> %s", i+1, horsetype, err)
		}
		if racetype, err := meta.GenericDataFromRegex(page, REGEX_RACE_RACETYPE); err != nil || racetype == nil {
			t.Errorf("wasn't able to gather generic data from page %d, %s\nerror >> %s", i+1, racetype, err)
		}
		if purse, err := meta.GenericDataFromRegex(page, REGEX_RACE_PURSE); err != nil || purse == nil {
			t.Errorf("wasn't able to gather generic data from page %d, %s\nerror >> %s", i+1, purse, err)
		}
		if weather, err := meta.GenericDataFromRegex(page, REGEX_RACE_WEATHER); err != nil || weather == nil {
			t.Errorf("wasn't able to gather generic data from page %d, %s\nerror >> %s", i+1, weather, err)
		}
		if trackcondition, err := meta.GenericDataFromRegex(page, REGEX_RACE_TRACK_CONDITION); err != nil || trackcondition == nil {
			t.Errorf("wasn't able to gather generic data from page %d, %s\nerror >> %s", i+1, trackcondition, err)
		}
		if length, err := meta.GenericDataFromRegex(page, REGEX_RACE_LENGTH); err != nil || length == nil {
			t.Errorf("wasn't able to gather generic data from page %d, %s\nerror >> %s", i+1, length, err)
		}
		if trackrecord, err := meta.GenericDataFromRegex(page, REGEX_RACE_CURRENT_TRACK_RECORD); err != nil || trackrecord == nil {
			t.Errorf("wasn't able to gather generic data from page %d, %s\nerror >> %s", i+1, trackrecord, err)
		}
		if finaltime, err := meta.GenericDataFromRegex(page, REGEX_RACE_FINAL_TIME); err != nil || finaltime == nil {
			t.Errorf("wasn't able to gather generic data from page %d, %s\nerror >> %s", i+1, finaltime, err)
		}
		if fractionaltimes, err := meta.GenericDataFromRegex(page, REGEX_RACE_FRACTIONAL_TIMES); err != nil || fractionaltimes == nil {
			t.Errorf("wasn't able to gather generic data from page %d, %s\nerror >> %s", i+1, fractionaltimes, err)
		}
		if splittimes, err := meta.GenericDataFromRegex(page, REGEX_RACE_SPLIT_TIMES); err != nil || splittimes == nil {
			t.Errorf("wasn't able to gather generic data from page %d, %s\nerror >> %s", i+1, splittimes, err)
		}
		if runup, err := meta.GenericDataFromRegex(page, REGEX_RACE_RUN_UP); err != nil || runup == nil {
			t.Errorf("wasn't able to gather generic data from page %d, %s\nerror >> %s", i+1, runup, err)
		}
		// }
	}
}

func TestNewRacePage(t *testing.T) {
	pages, err := GetValidPages(testFile)
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
