package equieasy

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	pdf "github.com/dslipak/pdf"
)

const (
	REGEX_VALID_RACE_PAGE      = `.*\- [A-Za-z]+ [0-9]+\, [0-9]{4} \- Race [0-9]+`
	REGEX_VALID_CANCELLED      = `Cancelled.*\-.*`
	REGEX_GET_HORSES           = `(?P<datetrack>([\-]{3}\s+)|([0-9]{1}[0-9]*(Jan|Feb|Mar|Apr|May|Jun|Jul|Aug|Sep|Oct|Nov|Dec)[0-9]{1}[0-9]* [A-Z]{2}[A-Z]*))(?P<pgm>(\s)*[0-9ABCX]+) (?P<horsename>[A-Za-z0-9\s\'\"\!\.\,\-\_\*\$\?]+[A-Z\(\)\w]*) \((?P<jockey>[A-Za-z0-9\,\.\s\'\-]+)\) (?P<wgt>[0-9]{3})( |.*)(?P<me>[A-Za-z\-\s]+) (?P<postposition>[0-9]{1}|[0-9]{2}) .* (?P<odds>[0-9]+\.[0-9\*]+) (?P<comment>[A-Za-z0-9\(\)\{\}\[\]\/\,\.\s\-\&\:\;\'\"\|]+)`
	REGEX_LAST_DATE_TRACK      = `(?P<date>[0-9]{1}[0-9]*(Jan|Feb|Mar|Apr|May|Jun|Jul|Aug|Sep|Oct|Nov|Dec)[0-9]{1}[0-9]*) (?P<track>[A-Z]{2}[A-Z]*)`
	REGEX_FRACTIONALS          = `(Pg m Horse Name) (?P<fracs>(Start|1\/4|3\/8|1\/2) .* Str Fin)`
	REGEX_TOP_VALUES           = `([0-9\/\s])+|(Head)+[\s]+|(Nose)+[\s]+|(Neck)+[\s]+`
	REGEX_HORSE_POSITIONAL     = `[0-9]{1}[0-9ABX]* (?P<horsename>[A-Za-z0-9\s\'\"\!\.\,\-\_\*\$\?]+[A-Z\(\)\w]*) [0-9]{1}[0-9]* [0-9]{1}[0-9]* [0-9]{1}[0-9]*.*`
	REGEX_HORSE_TABLE_W_START  = `[0-9]{1}[0-9ABX]* (?P<horsename>[A-Za-z0-9\s\'\"\!\.\,\-\_\*\$\?]+[A-Z\(\)\w]*) (?P<start>[0-9]{1}[0-9]*) (?P<firstfrac>[0-9]{1}[0-9]*)`
	REGEX_HORSE_TABLE_WO_START = `[0-9]{1}[0-9ABX]* (?P<horsename>[A-Za-z0-9\s\'\"\!\.\,\-\_\*\$\?]+[A-Z\(\)\w]*) (?P<firstfrac>[0-9]{1}[0-9]*)`
	REGEX_TRAINERS             = `Trainers:(?P<trainers>.*)`
	REGEX_TRAINERS_STOP        = `Owners:.*`
	REGEX_OWNERS               = `Owners:(?P<owners>.*)`
	REGEX_OWNERS_STOP          = `Footnotes.*`
	REGEX_IND_TRAINER_OWNER    = `(?P<number>[0-9]{1}[0-9abxABX]*)(?P<sep>(\s|\-)(\s|\-)*(\s)*)(?P<name>[A-Za-z0-9\,\.\s\'\-\(\)]+)`

	REGEX_RACE_TRACK                = `[ \*\[\]\/\\\\?\<\>]`
	REGEX_RACE_NUMBER               = `Race (?P<number>[0-9]+)`
	REGEX_RACE_HORSETYPE            = `(MAIDEN|CLAIMING|STARTER|ALLOWANCE|STAKES).* - (?P<value>[A-Za-z0-9]+)`
	REGEX_RACE_RACETYPE             = `(?P<value>(MAIDEN|CLAIMING|STARTER|ALLOWANCE|STAKES).*) - ([A-Za-z0-9]+)`
	REGEX_RACE_PURSE                = `(Purse: )(?P<value>.*)`
	REGEX_RACE_WEATHER              = `(Weather: )(?P<value>.*) Track: .*`
	REGEX_RACE_TRACK_CONDITION      = `(Weather: )(.*) (Track: )(?P<value>.*)`
	REGEX_RACE_LENGTH               = `(?P<value>.* (On The Inner turf|On The Outer turf|On The Dirt|On The Turf|On The Downhill turf|On The Downhill Turf))(.*)`
	REGEX_RACE_CURRENT_TRACK_RECORD = `.*(Current Track Record:|Track Record:) (?P<value>.*)`
	REGEX_RACE_FINAL_TIME           = `.*Final Time: (?P<value>.*)`
	REGEX_RACE_FRACTIONAL_TIMES     = `Fractional Times: (?P<value>.*) Final Time:.*`
	REGEX_RACE_SPLIT_TIMES          = `Split Times: (?P<value>.*)`
	REGEX_RACE_RUN_UP               = `Run-Up: (?P<value>.*)`

	OFFSET_TOP_X = 3.892
	OFFSET_TOP_Y = 3.351
)

type RacePage struct {
	Metadata *RaceMetadata
	Horses   []*Horse
}

type RaceMetadata struct {
	RaceHash           string
	Track              string
	Type               string
	Number             int
	HorseType          string
	Date               time.Time
	Purse              string
	Weather            string
	TrackCondition     string
	RaceLength         string
	CurrentTrackRecord string
	FinalTime          string
	FractionalTimes    string
	SplitTimes         string
	RunUp              string
}

type Horse struct {
	RaceHash      string
	Name          string
	Number        string
	PostPosition  string
	Weight        string
	MedEquipment  string
	Jockey        string
	Trainer       string
	Owners        string
	Odds          string
	Comments      string
	LastRaced     time.Time
	LastTrack     string
	Fractionals   []Fractional
	withTopOffset bool
}

type Fractional struct {
	Distance string
	Lengths  string
	Position string
}

type ValidPages struct {
	Filename string
	Pages    []pdf.Page
}

// Returns all valid race pages in an Equibase Document.
func GetValidPages(f string) (ValidPages, error) {
	racepages := []pdf.Page{}
	r, err := pdf.Open(f)
	if err != nil {
		return ValidPages{}, err
	}

	for i := 1; i <= r.NumPage(); i++ {
		page := r.Page(i)
		if isValidRacePage(page) {
			racepages = append(racepages, page)
		}
	}
	return ValidPages{Filename: f, Pages: racepages}, nil
}

// Generates a new RacePage including all horses and metadata about the race.
func NewRacePage(page pdf.Page) (*RacePage, error) {
	meta, err := Metadata(page)
	if err != nil {
		return nil, err
	}
	horses, err := Horses(page)
	if err != nil {
		return nil, err
	}
	return &RacePage{Metadata: meta, Horses: horses}, nil
}

// Create a new metadata struct from a page
func Metadata(page pdf.Page) (*RaceMetadata, error) {
	re := regexp.MustCompile(REGEX_RACE_TRACK)
	renum := regexp.MustCompile(REGEX_RACE_NUMBER)
	r := &RaceMetadata{}
	var racehash string

	rows, err := page.GetTextByRow()
	if err != nil {
		return nil, err
	}

	for i, row := range rows {
		rowdata := []byte{}
		if row.Position == 760 {
			for _, word := range row.Content {
				rowdata = append(rowdata, word.S...)
				rowdata = append(rowdata, " "...)
			}
			if i == 0 {
				racehash = fmt.Sprintf("%x", md5.Sum(rowdata))
			}

			data := bytes.Split(rowdata, []byte(" - "))

			if len(data) != 3 {
				return nil, fmt.Errorf("wasn't able to gather track and date data from race")
			}

			r.Track = string(re.ReplaceAll(data[0], []byte("")))
			date, err := time.Parse("January 2, 2006", string(data[1]))
			if err != nil {
				return nil, err
			}
			r.Date = date
			r.Number, err = strconv.Atoi(string(bytes.TrimSpace(renum.ReplaceAll(data[2], []byte("$number")))))
			if err != nil {
				return nil, err
			}
		}
	}
	horsetype, err := r.genericDataFromRegex(page, REGEX_RACE_HORSETYPE)
	if err != nil {
		return nil, err
	}
	racetype, err := r.genericDataFromRegex(page, REGEX_RACE_RACETYPE)
	if err != nil {
		return nil, err
	}
	purse, err := r.genericDataFromRegex(page, REGEX_RACE_PURSE)
	if err != nil {
		return nil, err
	}
	weather, err := r.genericDataFromRegex(page, REGEX_RACE_WEATHER)
	if err != nil {
		return nil, err
	}
	trackcondition, err := r.genericDataFromRegex(page, REGEX_RACE_TRACK_CONDITION)
	if err != nil {
		return nil, err
	}
	length, err := r.genericDataFromRegex(page, REGEX_RACE_LENGTH)
	if err != nil {
		return nil, err
	}
	trackrecord, err := r.genericDataFromRegex(page, REGEX_RACE_CURRENT_TRACK_RECORD)
	if err != nil {
		return nil, err
	}
	finaltime, err := r.genericDataFromRegex(page, REGEX_RACE_FINAL_TIME)
	if err != nil {
		return nil, err
	}
	fractionaltimes, err := r.genericDataFromRegex(page, REGEX_RACE_FRACTIONAL_TIMES)
	if err != nil {
		return nil, err
	}
	splittimes, err := r.genericDataFromRegex(page, REGEX_RACE_SPLIT_TIMES)
	if err != nil {
		return nil, err
	}
	runup, err := r.genericDataFromRegex(page, REGEX_RACE_RUN_UP)
	if err != nil {
		return nil, err
	}
	r.RaceHash = string(racehash)
	r.HorseType = string(horsetype)
	r.Type = string(racetype)
	r.Purse = string(purse)
	r.Weather = string(weather)
	r.TrackCondition = string(trackcondition)
	r.RaceLength = string(length)
	r.CurrentTrackRecord = string(trackrecord)
	r.FinalTime = string(finaltime)
	r.FractionalTimes = string(fractionaltimes)
	r.SplitTimes = string(splittimes)
	r.RunUp = string(runup)
	return r, nil
}

// Get a list of all Horses racing on a race page.
func Horses(page pdf.Page) ([]*Horse, error) {
	horses := []*Horse{}
	re := regexp.MustCompile(REGEX_GET_HORSES)
	re2 := regexp.MustCompile(REGEX_LAST_DATE_TRACK)
	racehash := ""

	rows, err := page.GetTextByRow()
	// fmt.Println(rows, " err >>>> ", err)
	if err != nil {
		return nil, err
	}
	for i, row := range rows {
		rowdata := []byte{}
		for _, word := range row.Content {
			rowdata = append(rowdata, word.S...)
			rowdata = append(rowdata, " "...)
		}
		if i == 0 {
			racehash = fmt.Sprintf("%x", md5.Sum(rowdata))
		}

		if re.Match(rowdata) {
			dt := string(re.ReplaceAll(rowdata, []byte("$datetrack")))
			date := time.Time{}
			if dt[:3] != "---" {
				date, err = time.Parse("2Jan06", string(re2.ReplaceAllString(dt, "$date")))
				if err != nil {
					return nil, err
				}
			}
			horse := &Horse{
				RaceHash:     racehash,
				Name:         strings.TrimSpace(string(re.ReplaceAll(rowdata, []byte("$horsename")))),
				Number:       strings.TrimSpace(string(re.ReplaceAll(rowdata, []byte("$pgm")))),
				PostPosition: strings.TrimSpace(string(re.ReplaceAll(rowdata, []byte("$postposition")))),
				Weight:       strings.TrimSpace(string(re.ReplaceAll(rowdata, []byte("$weight")))),
				MedEquipment: strings.TrimSpace(string(re.ReplaceAll(rowdata, []byte("$me")))),
				Jockey:       strings.TrimSpace(string(re.ReplaceAll(rowdata, []byte("$jockey")))),
				Odds:         strings.TrimSpace(string(re.ReplaceAll(rowdata, []byte("$odds")))),
				Comments:     strings.TrimSpace(string(re.ReplaceAll(rowdata, []byte("$comment")))),
				LastRaced:    date,
				LastTrack:    strings.TrimSpace(string(re2.ReplaceAllString(dt, "$track"))),
			}
			horses = append(horses, horse)
		}
	}
	return horses, nil
}

// TODO:
// Scans a wider margin for DOH (Distance to Other Horses) values rather than just the expected spot.
func SetTopOffsetOn() {}

// Returns a list of the positioned fractionals for a race.
func fractionals(page pdf.Page) ([]string, int64, error) {
	re := regexp.MustCompile(REGEX_FRACTIONALS)
	rows, err := page.GetTextByRow()
	if err != nil {
		return nil, -99, err
	}

	for _, row := range rows {
		rowdata := []byte{}
		for _, word := range row.Content {
			rowdata = append(rowdata, word.S...)
			rowdata = append(rowdata, " "...)
		}
		if re.Match(rowdata) {
			fractionals := strings.Split(string(re.ReplaceAll(bytes.Trim(rowdata, " "), []byte("$fracs"))), " ")
			return fractionals, row.Position, nil
		}
	}
	return nil, -99, fmt.Errorf("failed to gather any fractional data from page")
}

func (h *Horse) applyFractionalData(page pdf.Page) error {
	frs, y, err := fractionals(page)
	if err != nil {
		return err
	}

	rows, err := page.GetTextByRow()
	if err != nil {
		return err
	}

	re := regexp.MustCompile(REGEX_HORSE_TABLE_W_START)
	re2 := regexp.MustCompile(REGEX_HORSE_TABLE_WO_START)

	fractionals := []Fractional{}
	listoftop := [][]float64{}

	for _, row := range rows {
		if row.Position < y {
			rowdata := []byte{}
			prevX := 9999.0
			for _, word := range row.Content {
				if word.X > prevX {
					rowdata = append(rowdata, "| "...)
				}
				rowdata = append(rowdata, word.S...)
				rowdata = append(rowdata, " "...)
				prevX = word.X
			}

			var horseString string
			if len(h.Name) > 7 {
				horseString = fmt.Sprintf("%s %s", h.Number, h.Name[:5])
			} else {
				horseString = fmt.Sprintf("%s %s", h.Number, h.Name)
			}

			if len(rowdata) >= len(horseString) {
				if horseString == string(rowdata[:len(horseString)]) {
					values := bytes.Split(rowdata, []byte(" | "))
					for _, val := range values {
						switch frs[0] {
						case "Start":
							if re.Match(val) {
								start := Fractional{Position: strings.Trim(string(re.ReplaceAll(val, []byte("$start"))), " ")}
								firstfrac := Fractional{Position: strings.Trim(string(re.ReplaceAll(val, []byte("$firstfrac"))), " ")}
								fractionals = append(fractionals, start, firstfrac)
							} else {
								frac := Fractional{Position: strings.Trim(string(val), " ")}
								fractionals = append(fractionals, frac)
							}
						default:
							if re2.Match(val) {
								firstfrac := Fractional{Position: strings.Trim(string(re.ReplaceAll(val, []byte("$firstfrac"))), " ")}
								fractionals = append(fractionals, firstfrac)
							} else {
								frac := Fractional{Position: strings.Trim(string(val), " ")}
								fractionals = append(fractionals, frac)
							}

						}
					}

					prevX = 0
					for _, word := range row.Content {
						if word.X > prevX {
							listoftop = append(listoftop, []float64{roundFloat(word.X+OFFSET_TOP_X, 3), roundFloat(word.Y+OFFSET_TOP_Y, 3)})
						}
						prevX = word.X
					}
				}
			}
		}
	}

	listoftop[0][0] = roundFloat(listoftop[1][0]-(listoftop[2][0]-listoftop[1][0]), 3)

	offset := len(fractionals) - len(listoftop)

	// TODO:
	// When function in place, set default to 0.0
	offset_from_word := 4.0
	if h.withTopOffset {
		offset_from_word = 4.0
	}

	for i, item := range listoftop {
		for _, row := range rows {
			prevX := 0.0
			for _, word := range row.Content {
				if (word.X >= item[0]-offset_from_word) && (word.X <= item[0]+offset_from_word) && (word.Y == item[1]) {
					if prevX == word.X {
						fractionals[i+offset].Lengths += " " + word.S
					}
					if prevX != word.X {
						fractionals[i+offset].Lengths = word.S
					}
				}
				prevX = word.X
			}
		}
	}

	for i, fr := range frs {
		fractionals[i].Distance = fr
		h.Fractionals = append(h.Fractionals, fractionals[i])
	}
	return nil
}

// Scans the document for Trainers and applies them to the horse
func (h *Horse) applyTrainerData(page pdf.Page) error {
	re := regexp.MustCompile(REGEX_TRAINERS)
	retr := regexp.MustCompile(REGEX_IND_TRAINER_OWNER)
	restop := regexp.MustCompile(REGEX_TRAINERS_STOP)
	trainerdatavalid := false

	rows, err := page.GetTextByRow()
	if err != nil {
		return err
	}

	trainerdata := []byte{}
	for _, row := range rows {
		rowdata := []byte{}
		for _, word := range row.Content {
			rowdata = append(rowdata, word.S...)
			rowdata = append(rowdata, " "...)
		}

		if re.Match(rowdata) {
			trainerdatavalid = true
		}

		if restop.Match(rowdata) {
			trainerdatavalid = false
		}

		if trainerdatavalid {
			trainerdata = append(trainerdata, rowdata...)
		}
	}
	trainerdata = re.ReplaceAll(trainerdata, []byte("$trainers"))
	trainers := bytes.Split(trainerdata, []byte(";"))
	for _, tr := range trainers {
		tr = bytes.Trim(tr, " ")
		tnumb := bytes.ToUpper(retr.ReplaceAll(tr, []byte("$number")))
		tname := retr.ReplaceAll(tr, []byte("$name"))

		if h.Number == string(tnumb) {
			h.Trainer = string(tname)
		}
	}
	return nil
}

// Scans the document for Owners and applies them to the horse
func (h *Horse) applyOwnerData(page pdf.Page) error {
	re := regexp.MustCompile(REGEX_OWNERS)
	retr := regexp.MustCompile(REGEX_IND_TRAINER_OWNER)
	restop := regexp.MustCompile(REGEX_OWNERS_STOP)
	ownerdatavalid := false

	rows, err := page.GetTextByRow()
	if err != nil {
		return err
	}

	ownerdata := []byte{}
	for _, row := range rows {
		rowdata := []byte{}
		for _, word := range row.Content {
			rowdata = append(rowdata, word.S...)
			rowdata = append(rowdata, " "...)
		}

		if re.Match(rowdata) {
			ownerdatavalid = true
		}

		if restop.Match(rowdata) {
			ownerdatavalid = false
		}

		if ownerdatavalid {
			ownerdata = append(ownerdata, rowdata...)
		}
	}
	ownerdata = re.ReplaceAll(ownerdata, []byte("$owners"))

	owners := bytes.Split(ownerdata, []byte(";"))
	for _, or := range owners {
		or = bytes.Trim(or, " ")
		onumb := bytes.ToUpper(retr.ReplaceAll(or, []byte("$number")))
		oname := retr.ReplaceAll(or, []byte("$name"))

		if h.Number == string(onumb) {
			h.Owners = string(oname)
		}
	}
	return nil
}

// Gather generic data from page from a regex expression.
func (r *RaceMetadata) genericDataFromRegex(page pdf.Page, regex string) ([]byte, error) {
	re := regexp.MustCompile(regex)

	rows, err := page.GetTextByRow()
	if err != nil {
		return nil, err
	}

	for _, row := range rows {
		rowdata := []byte{}
		for _, word := range row.Content {
			rowdata = append(rowdata, word.S...)
			rowdata = append(rowdata, " "...)
		}

		if re.Match(rowdata) {
			data := re.ReplaceAll(rowdata, []byte("$value"))
			return data, nil
		}
	}
	if regex == REGEX_RACE_SPLIT_TIMES {
		return nil, nil
	}
	return nil, fmt.Errorf("wasn't able to grab the generic data requested\nregex >> %s", regex)
}
