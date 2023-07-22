package scraper

import "time"

const (
	PreSeason  SeasonType = 1
	RegSeason  SeasonType = 2
	PostSeason SeasonType = 3

	customTimeLayout = "2006-01-02T15:04Z"

	ESPNDomain = "https://www.espn.com"
)

type (
	SeasonType int

	Matchup struct {
		AwayTeam, HomeTeam Team
		Date               string
		Time               string
	}
	Team struct {
		TeamName, Abbreviation string
	}

	Matchups []Matchup

	BySeasonType []Week

	CustomTime struct {
		time.Time
	}

	Week struct {
		Text       string     `json:"text"`
		Label      string     `json:"label"`
		StartDate  CustomTime `json:"startDate"`
		EndDate    CustomTime `json:"endDate"`
		SeasonType SeasonType `json:"seasonType"`
		WeekNumber int        `json:"weekNumber"`
		Year       int        `json:"year"`
		URL        string     `json:"url"`
		IsActive   bool       `json:"isActive"`
	}

	Game struct {
		ID          string       `json:"id"`
		Competitors []Competitor `json:"competitors"`
		Date        string       `json:"date"`
		TBD         bool         `json:"tbd"`
		Completed   bool         `json:"completed"`
		Link        string       `json:"link"`
		Teams       []NFLTeam    `json:"teams"`
		IsTie       bool         `json:"isTie"`
		TimeValid   bool         `json:"timeValid"`
	}

	Games map[string][]Game

	Competitor struct {
		ID               string `json:"id"`
		Abbrev           string `json:"abbrev"`
		DisplayName      string `json:"displayName"`
		ShortDisplayName string `json:"shortDisplayName"`
		Logo             string `json:"logo"`
		TeamColor        string `json:"teamColor"`
		AltColor         string `json:"altColor"`
		UID              string `json:"uid"`
		RecordSummary    string `json:"recordSummary"`
		StandingSummary  string `json:"standingSummary"`
		Location         string `json:"location"`
		Links            string `json:"links"`
		Name             string `json:"name"`
		ShortName        string `json:"shortName"`
		IsHome           bool   `json:"isHome"`
	}

	NFLTeam struct {
		ID               string `json:"id"`
		Abbrev           string `json:"abbrev"`
		DisplayName      string `json:"displayName"`
		ShortDisplayName string `json:"shortDisplayName"`
		Logo             string `json:"logo"`
		TeamColor        string `json:"teamColor"`
		AltColor         string `json:"altColor"`
		UID              string `json:"uid"`
		RecordSummary    string `json:"recordSummary"`
		StandingSummary  string `json:"standingSummary"`
		Location         string `json:"location"`
		Links            string `json:"links"`
		Name             string `json:"name"`
		ShortName        string `json:"shortName"`
		IsHome           bool   `json:"isHome"`
	}

	GameInfo struct {
		GameID      string `json:"gid,omitempty"`
		SeasonType  int    `json:"seasonType,omitempty"`
		Status      Status `json:"status,omitempty"`
		StatusState string `json:"statusState,omitempty"`
		Tbd         bool   `json:"tbd,omitempty"`
		Tms         []Tms  `json:"tms,omitempty"`
	}
	Status struct {
		Desc  string `json:"desc,omitempty"`
		Det   string `json:"det,omitempty"`
		ID    string `json:"id,omitempty"`
		State string `json:"state,omitempty"`
	}
	Records struct {
		Type         string `json:"type,omitempty"`
		Summary      string `json:"summary,omitempty"`
		DisplayValue string `json:"displayValue,omitempty"`
	}
	Linescores struct {
		DisplayValue string `json:"displayValue,omitempty"`
	}
	Tms struct {
		Abbrev           string       `json:"abbrev,omitempty"`
		DisplayName      string       `json:"displayName,omitempty"`
		ShortDisplayName string       `json:"shortDisplayName,omitempty"`
		Records          []Records    `json:"records,omitempty"`
		IsHome           bool         `json:"isHome,omitempty"`
		Linescores       []Linescores `json:"linescores,omitempty"`
		Score            string       `json:"score,omitempty"`
		Winner           bool         `json:"winner,omitempty"`
	}
)

func (ct *CustomTime) UnmarshalJSON(b []byte) error {
	s := string(b)
	t, err := time.Parse(customTimeLayout, s[1:len(s)-1])
	if err != nil {
		return err
	}
	ct.Time = t.UTC()
	return nil
}

// Implement the sort.Interface interface for the BySeason type
func (w BySeasonType) Len() int      { return len(w) }
func (w BySeasonType) Swap(i, j int) { w[i], w[j] = w[j], w[i] }
func (w BySeasonType) Less(i, j int) bool {

	if w[i].Year != w[j].Year {
		return w[i].Year < w[j].Year
	}

	if w[i].SeasonType != w[j].SeasonType {
		return w[i].SeasonType < w[j].SeasonType
	}

	return w[i].WeekNumber < w[j].WeekNumber
}
