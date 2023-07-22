package scraper

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"sort"
	"strings"
)

//go:generate mockgen -destination scraper_mock.go -package scraper . ScheduleScraper
var _ ScheduleScraper = &Scraper{}

type (
	ScheduleScraper interface {
		FetchSchedule() (BySeasonType, error)
		FetchGamesForWeeks(weeks []Week) (Games, error)
		FetchGamesForWeek(week Week) (Games, error)
		FetchGameInfo(gameID string) (GameInfo, error)
	}
	Scraper struct {
		httpClient *http.Client
	}
)

func New(httpClient *http.Client) *Scraper {
	return &Scraper{
		httpClient: httpClient,
	}
}

func (s *Scraper) FetchSchedule() (BySeasonType, error) {
	url := ESPNDomain + "/nfl/schedule"

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to create request for %s: %w", url, err)
	}

	res, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("unable to make request for %s: %w", url, err)
	}

	bytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to read bytes from request %s: %w", url, err)
	}

	return findWeeksFromBytes(bytes)
}

func findWeeksFromBytes(bytes []byte) (BySeasonType, error) {
	weekArrayRegExp := `"weeks":\s?(\[.*?\])`
	flatString := strings.Join(strings.Split(string(bytes), "\n"), "")
	regex, err := regexp.Compile(weekArrayRegExp)
	if err != nil {
		return nil, fmt.Errorf("unable to compile regex %s: %w", weekArrayRegExp, err)
	}

	var weeks BySeasonType
	matches := regex.FindAllStringSubmatch(flatString, -1)
	match := matches[(len(matches))-1:][0]
	if len(match) > 1 {
		err = json.Unmarshal([]byte(match[1]), &weeks)
		if err != nil {
			return nil, fmt.Errorf("unable to unmarshal string %s: %w", match[1], err)
		}
	}

	weeks = weeks[len(weeks)-27:]
	sort.Sort(weeks)

	return weeks, nil
}

func (s *Scraper) FetchGamesForWeeks(weeks []Week) (Games, error) {
	games := Games{}

	for _, week := range weeks {
		gamesOfWeek, err := s.FetchGamesForWeek(week)
		if err != nil {
			return nil, fmt.Errorf("unable to fetch games for week %+v: %w", week, err)
		}
		for k, v := range gamesOfWeek {
			games[k] = v
		}
	}

	return games, nil
}

func (s *Scraper) FetchGamesForWeek(week Week) (Games, error) {
	url := ESPNDomain + week.URL

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return Games(nil), fmt.Errorf("unable to create request for %s: %w", url, err)
	}

	res, err := s.httpClient.Do(req)
	if err != nil {
		return Games(nil), fmt.Errorf("unable to make request for %s: %w", url, err)
	}

	bytes, err := io.ReadAll(res.Body)
	if err != nil {
		return Games(nil), fmt.Errorf("unable to read bytes from request %s: %w", url, err)
	}
	res.Body.Close()
	return findGamesFromBytes(bytes)
}
func findGamesFromBytes(bArr []byte) (Games, error) {

	flatString := strings.Join(strings.Split(string(bArr), "\n"), "")
	jsonStr := jsonStringFirstIndex(`"events":`, flatString)

	var games Games
	err := json.Unmarshal([]byte(jsonStr), &games)
	if err != nil {
		return nil, fmt.Errorf("unable to unmarshal string %s: %w", jsonStr, err)
	}

	return games, nil
}
func (s *Scraper) FetchGameInfo(gameID string) (GameInfo, error) {
	url := ESPNDomain + "/nfl/game/_/gameId/" + gameID
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return GameInfo{}, fmt.Errorf("unable to create request for %s: %w", url, err)
	}

	res, err := s.httpClient.Do(req)
	if err != nil {
		return GameInfo{}, fmt.Errorf("unable to make request for %s: %w", url, err)
	}

	bytes, err := io.ReadAll(res.Body)
	if err != nil {
		return GameInfo{}, fmt.Errorf("unable to read bytes from request %s: %w", url, err)
	}
	return findGameInfoFromBytes(bytes)
}

func findGameInfoFromBytes(bArr []byte) (GameInfo, error) {

	flatString := strings.Join(strings.Split(string(bArr), "\n"), "")
	jsonStr := jsonStringLastIndex(`"gmStrp":`, flatString)

	var gameInfo GameInfo
	err := json.Unmarshal([]byte(jsonStr), &gameInfo)
	if err != nil {
		return GameInfo{}, fmt.Errorf("unable to unmarshal string %s: %w", jsonStr, err)
	}

	return gameInfo, nil
}

func jsonStringFirstIndex(word, html string) string {

	str := html[strings.Index(html, word)+len(word):]
	cnt := 0
	strB := strings.Builder{}

	for _, b := range str {
		strB.WriteRune(b)
		switch b {
		case '{':
			cnt++
		case '}':
			cnt--
		}
		if cnt == 0 {
			break
		}
	}

	return strB.String()

}

func jsonStringLastIndex(word, html string) string {

	str := html[strings.LastIndex(html, word)+len(word):]
	cnt := 0
	strB := strings.Builder{}

	for _, b := range str {
		strB.WriteRune(b)
		switch b {
		case '{':
			cnt++
		case '}':
			cnt--
		}
		if cnt == 0 {
			break
		}
	}

	return strB.String()

}
