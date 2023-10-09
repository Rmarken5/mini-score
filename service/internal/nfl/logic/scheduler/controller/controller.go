package controller

import (
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/rmarken5/mini-score/service/internal/nfl/logic/internal/data-access/db/repository"
	"github.com/rmarken5/mini-score/service/internal/nfl/logic/internal/data-access/http/rest"
	"github.com/rmarken5/mini-score/service/internal/nfl/logic/internal/data-access/http/scraper"
	"github.com/rs/zerolog"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

type (
	Controller interface {
		KeepScheduleSynchronized(loopExiter <-chan bool, iterationInterval time.Duration)
		GetGamesBetweenDates(start time.Time, end time.Time) ([]repository.Game, error)
		GetGameInfo(gameID string) (scraper.GameInfo, error)
		UpdateGame(gameInfo scraper.GameInfo)
		FinalizeGame(gameInfo scraper.GameInfo)
	}

	Logic struct {
		logger               zerolog.Logger
		scrapper             scraper.ScheduleScraper
		repo                 repository.Repository
		requester            rest.Requester
		gameTeamQuarterCache map[string]int
		scoreCacheLock       sync.RWMutex
		clockCache           map[string]string
		clockCacheLock       sync.RWMutex
	}
)

func NewLogic(logger zerolog.Logger, db *sqlx.DB) *Logic {

	logger = logger.With().Str("service", "Logic").Logger()

	httpClient := &http.Client{}
	s := scraper.New(httpClient)
	restRequester := rest.NewRequester(logger, httpClient)

	return &Logic{
		logger:               logger,
		scrapper:             s,
		repo:                 repository.NewRepository(logger, db),
		requester:            restRequester,
		gameTeamQuarterCache: make(map[string]int),
		clockCache:           make(map[string]string),
	}
}
func (l *Logic) KeepScheduleSynchronized(loopExiter <-chan bool, iterationInterval time.Duration) {
	logger := l.logger.With().Str("method", "KeepScheduleSynchronized").Logger()
	logger.Info().Msgf("Starting KeepScheduleSynchronized")

	logger.Info().Msgf("calling sync")
	if err := l.syncSchedule(); err != nil {
		logger.Info().Err(err).Msgf("error syncing schedule")
	}
	for {
		select {
		case <-loopExiter:
			return
		case <-time.After(iterationInterval):
			logger.Info().Msgf("calling sync")
			if err := l.syncSchedule(); err != nil {
				logger.Info().Err(err).Msgf("error syncing schedule")
			}
		}
	}
}

func (l *Logic) syncSchedule() error {
	weeks, err := l.scrapper.FetchSchedule()
	if err != nil {
		return err
	}

	gameMap, err := l.scrapper.FetchGamesForWeeks(weeks)
	if err != nil {
		return err
	}
	for _, games := range gameMap {
		l.processGameMap(games)
	}
	return nil
}
func (l *Logic) processGameMap(games []scraper.Game) {
	logger := l.logger.With().Str("method", "processGameMap").Logger()
	for _, game := range games {
		if err := l.processGame(game); err != nil {
			logger.Info().Err(err).Fields(game).Msg("processing game")
		}
	}
}

func (l *Logic) processGame(game scraper.Game) error {
	repoGame, err := l.repo.GetGame(game.ID)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrNoGame):
			err := l.insertGame(game)
			if err != nil {
				return err
			}
			for _, team := range game.Teams {
				for i := 0; i < 4; i++ {
					quarterScore := repository.GameQuarterScore{
						GameID:  game.ID,
						TeamID:  team.Abbrev,
						Quarter: strconv.Itoa(i + 1),
						Score:   0,
					}
					err := l.repo.InsertQuarterScore(quarterScore)
					if err != nil {
						return err
					}
				}
			}
			return nil
		}
		return err
	}
	return l.updateGameTime(game, repoGame)
}
func (l *Logic) updateGameTime(game scraper.Game, repoGame repository.Game) error {
	layoutStr := "2006-01-02T15:04Z"
	gameTime, err := time.Parse(layoutStr, game.Date)
	if err != nil {
		return err
	}
	if !repoGame.GameTime.Equal(gameTime) {
		err := l.repo.UpdateGameTime(game.ID, gameTime)
		if err != nil {
			return err
		}
	}
	return nil
}

func (l *Logic) insertGame(game scraper.Game) error {
	logger := l.logger.With().Str("method", "insertGame").Logger()

	repoGame, err := l.fromScraperGameToRepo(game)
	if err != nil {
		logger.Info().Err(err).Msgf("unable to convert from scraper game to repo game")
		return err
	}
	err = l.repo.InsertGame(repoGame)
	if err != nil {
		logger.Info().Err(err).Msgf("unable to insert game")
		return err
	}
	return nil
}

func (l *Logic) fromScraperGameToRepo(game scraper.Game) (repository.Game, error) {

	layoutStr := "2006-01-02T15:04Z"
	gameTime, err := time.Parse(layoutStr, game.Date)
	if err != nil {
		return repository.Game{}, err
	}
	compOne := game.Competitors[0]
	compTwo := game.Competitors[1]

	teamOne, err := l.repo.GetTeamByAbv(compOne.Abbrev)
	if err != nil {
		return repository.Game{}, err
	}

	teamTwo, err := l.repo.GetTeamByAbv(compTwo.Abbrev)
	if err != nil {
		return repository.Game{}, err
	}
	g := repository.Game{
		ID:        game.ID,
		GameTime:  gameTime,
		Quarter:   "",
		GameClock: "",
	}

	g.AwayTeam = teamOne.ID.String()
	g.HomeTeam = teamTwo.ID.String()
	if compOne.IsHome {
		g.HomeTeam = teamOne.ID.String()
		g.AwayTeam = teamTwo.ID.String()
	}

	return g, nil
}

func (l *Logic) GetGamesBetweenDates(start time.Time, end time.Time) ([]repository.Game, error) {
	logger := l.logger.With().Str("method", "GetGamesBetweenDates").Logger()

	logger.Info().Msgf("getting games between: %s - %s", start, end)
	games, err := l.repo.GetGames(start, &end)
	if err != nil {
		logger.Error().Err(err).Msgf("error getting games")
		return nil, err
	}
	return games, err
}

func (l *Logic) GetGameInfo(gameID string) (scraper.GameInfo, error) {
	logger := l.logger.With().Str("method", "GetGameInfo").Logger()

	logger.Info().Msgf("Getting game info for: %s", gameID)
	return l.scrapper.FetchGameInfo(gameID)
}

func (l *Logic) UpdateGame(info scraper.GameInfo) {
	logger := l.logger.With().Str("method", "UpdateGame").Logger()
	go func() {
		if err := l.updateGameQuarterScore(info); err != nil {
			logger.Error().Err(err).Msgf("while trying to update game")
		}
	}()
	go func() {
		if err := l.updateGameClock(info.GameID); err != nil {
			logger.Error().Err(err).Msgf("while trying to update game clock")
		}
	}()
}

func (l *Logic) updateGameQuarterScore(gameInfo scraper.GameInfo) error {
	logger := l.logger.With().Str("method", "updateGameQuarterScore").Logger()
	l.logger.Info().Msgf("Updating game info for %s", gameInfo.GameID)

	gameInfoDAO, err := dtoFromGameInfo(gameInfo)
	if err != nil {
		return err
	}

	for _, team := range gameInfoDAO.TeamScores {
		for i, score := range team.Score {
			quarter := strconv.Itoa(i + 1)
			var scoreNum int
			scoreNum, err := strconv.Atoi(score)
			if err != nil {
				scoreNum = 0
			}

			cacheKey := gameInfoDAO.GameID + team.TeamAbbreviation + quarter
			if l.isTeamQuarterCacheCurrent(cacheKey, scoreNum) {
				logger.Info().Str("cache key", cacheKey).Msg("cache is current - skipping update")
				continue
			}
			logger.Info().Str("cache key", cacheKey).Msg("cache is not current - performing update")

			_, err = l.repo.GetQuarterScoreBy(gameInfoDAO.GameID, team.TeamAbbreviation, quarter)
			if err != nil {
				switch {
				case errors.Is(err, repository.ErrNoQuarterScore):
					logger.Info().Str("cache key", cacheKey).Msg("no quarter for team - performing insert")

					quarterScore := repository.GameQuarterScore{
						GameID:  gameInfoDAO.GameID,
						TeamID:  team.TeamAbbreviation,
						Quarter: quarter,
						Score:   scoreNum,
					}
					err := l.repo.InsertQuarterScore(quarterScore)
					if err != nil {
						logger.Error().Err(err).Msgf("while trying to insert game quarter score: %+v", quarterScore)
						continue
					}
				default:
					logger.Error().Err(err).Msg("while trying to get game quarter score.")
				}
				continue
			}
			logger.Info().Str("cache key", cacheKey).Msg("quarter for team exists - performing update")
			err = l.repo.UpdateQuarterScore(scoreNum, gameInfoDAO.GameID, team.TeamAbbreviation, quarter)
			if err != nil {
				logger.Error().Err(err).Msg("while trying to update game quarter score.")
				continue
			}
			l.updateTeamQuarterCache(cacheKey, scoreNum)
		}
	}
	return nil
}

func dtoFromGameInfo(gameInfo scraper.GameInfo) (repository.GameScoreDTO, error) {
	var firstTeamLineScore []string
	var secondTeamLineScore []string

	if len(gameInfo.Tms) < 2 {
		return repository.GameScoreDTO{}, fmt.Errorf("not enough teams in info to process: %+v", gameInfo.Tms)
	}

	for _, score := range gameInfo.Tms[0].Linescores {
		firstTeamLineScore = append(firstTeamLineScore, score.DisplayValue)
	}

	for _, score := range gameInfo.Tms[1].Linescores {
		secondTeamLineScore = append(secondTeamLineScore, score.DisplayValue)
	}

	gameScore := repository.GameScoreDTO{
		GameID: gameInfo.GameID,
		TeamScores: []repository.ScoreDTO{
			{
				TeamAbbreviation: gameInfo.Tms[0].Abbrev,
				Score:            firstTeamLineScore,
			},
			{
				TeamAbbreviation: gameInfo.Tms[1].Abbrev,
				Score:            secondTeamLineScore,
			},
		},
	}

	return gameScore, nil
}

func (l *Logic) isTeamQuarterCacheCurrent(key string, score int) bool {
	l.scoreCacheLock.RLock()
	defer l.scoreCacheLock.RUnlock()
	if cacheScore, ok := l.gameTeamQuarterCache[key]; ok {
		return cacheScore == score
	}
	return false
}

func (l *Logic) updateTeamQuarterCache(key string, score int) {
	l.scoreCacheLock.Lock()
	defer l.scoreCacheLock.Unlock()
	l.gameTeamQuarterCache[key] = score
}

func (l *Logic) updateGameClock(gameID string) error {
	logger := l.logger.With().Str("method", "updateGameClock").Logger()

	clock, quarter, err := l.getClockAndPeriodForGameID(gameID)
	if err != nil {
		logger.Error().Err(err).Msgf("while getting game clock for gameID: %s", gameID)
		return err
	}

	cachedClock := l.clockCacheByGameID(gameID)
	if cachedClock == clock {
		logger.Info().Msgf("Cached clock :%s matches fetched clock :%s, skipping update", cachedClock, clock)
		return nil
	}

	l.setClockCache(gameID, clock)

	err = l.repo.UpdateQuarterGameClock(gameID, quarter, clock)
	if err != nil {
		logger.Error().Err(err).Msgf("while updating game clock/quarter for gameID: %s", gameID)
		return err
	}
	return nil
}

func (l *Logic) setClockCache(gameID, clock string) {
	l.clockCacheLock.Lock()
	defer l.clockCacheLock.Unlock()
	l.clockCache[gameID] = clock
}

func (l *Logic) clockCacheByGameID(gameID string) string {
	l.clockCacheLock.RLock()
	defer l.clockCacheLock.RUnlock()
	return l.clockCache[gameID]
}

func (l *Logic) getClockAndPeriodForGameID(gameID string) (string, string, error) {
	logger := l.logger.With().Str("method", "GetScheduleForGameID").Logger()
	resp, err := l.requester.GetScoreboard()
	if err != nil {
		logger.Error().Err(err).Msgf("While getting scoreboard for gameID: %s", gameID)
		return "", "", err
	}
	gameTime, period, err := l.getTimeAndPeriodFromScoreboard(gameID, resp)
	if err != nil {
		logger.Error().Err(err).Msgf("unable to get time for gameID: %s", gameID)
		return "", "", err
	}
	return gameTime, period, nil
}

func (l *Logic) getTimeAndPeriodFromScoreboard(gameID string, response rest.ScoreboardResponse) (string, string, error) {
	var ev rest.Event
	for _, event := range response.Content.SBData.Events {
		if isGameIDMatch(gameID, event.UID) {
			ev = event
			break
		}
	}
	if ev.UID == "" {
		return "", "", fmt.Errorf("no match for game in scoreboard")
	}
	return ev.Status.DisplayClock, strconv.Itoa(ev.Status.Period), nil

}

func isGameIDMatch(gameID, gameUID string) bool {
	uid := gameUID[strings.LastIndex(gameUID, ":")+1:]
	return gameID == uid
}

func (l *Logic) FinalizeGame(gameInfo scraper.GameInfo) {
	logger := l.logger.With().Str("method", "FinalizeGame").Logger()
	defer l.clearGameClockCache(gameInfo.GameID)
	defer l.clearGameCache(gameInfo.GameID)
	go func(gameInfo scraper.GameInfo) {
		if err := l.updateGameQuarterScore(gameInfo); err != nil {
			logger.Error().Err(err).Msgf("while trying to update game")
		}
	}(gameInfo)
	go func(gameID, quarter string, gameClock string) {
		if err := l.repo.UpdateQuarterGameClock(gameID, quarter, gameClock); err != nil {
			logger.Error().Err(err).Msgf("while trying to update game clock")

		}
	}(gameInfo.GameID, "F", "Final")
}

func (l *Logic) clearGameClockCache(gameID string) {
	l.clockCacheLock.Lock()
	defer l.clockCacheLock.Unlock()
	delete(l.clockCache, gameID)
}

func (l *Logic) clearGameCache(gameID string) {
	l.scoreCacheLock.Lock()
	defer l.scoreCacheLock.Unlock()
	delete(l.gameTeamQuarterCache, gameID)
}
