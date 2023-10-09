package scheduler

import (
	"github.com/rmarken5/mini-score/service/internal/nfl/logic/internal/data-access/db/repository"
	"github.com/rmarken5/mini-score/service/internal/nfl/logic/internal/general"
	"github.com/rmarken5/mini-score/service/internal/nfl/logic/scheduler/controller"
	"github.com/rs/zerolog"
	"strconv"
	"strings"
	"sync"
	"time"
)

type (
	Scheduler struct {
		logger     zerolog.Logger
		controller controller.Controller
		games      map[string]repository.Game
		lock       sync.RWMutex
	}
)

func New(logger zerolog.Logger, ctrl controller.Controller) *Scheduler {
	return &Scheduler{
		logger:     logger.With().Str("service", "scheduler").Logger(),
		controller: ctrl,
		games:      make(map[string]repository.Game),
	}
}

func (s *Scheduler) Run() {
	gameChannel := make(chan repository.Game)
	end := make(chan bool)
	go s.controller.KeepScheduleSynchronized(nil, time.Hour*24)
	go s.SynchronizeCurrentWeek(gameChannel)
	go s.RunScheduler(gameChannel)
	<-end
}

// RunScheduler runs forever.
// It starts a goroutine to fetch game data for each game that should be started.
func (s *Scheduler) RunScheduler(gameChan <-chan repository.Game) {
	logger := s.logger.With().Str("method", "RunScheduler").Logger()
	logger.Info().Msgf("starting scheduler")
	for {
		select {
		case game := <-gameChan:
			logger.Debug().Fields(game).Msgf("getting game from channel")
			if !s.IsGameInList(game.ID) {
				logger.Debug().Fields(game).Msgf("game not lin list")
				s.AddGame(game)
				go s.GetGameInfo(game)
			}
		}
	}
}

func (s *Scheduler) AddGame(game repository.Game) {
	s.lock.Lock()
	defer s.lock.Unlock()
	if g, ok := s.games[game.ID]; !ok || !g.GameTime.Equal(game.GameTime) {
		s.games[game.ID] = game
	}
}

func (s *Scheduler) RemoveGame(gameID string) {
	s.lock.Lock()
	defer s.lock.Unlock()
	delete(s.games, gameID)
}

func (s *Scheduler) IsGameInList(gameID string) bool {
	s.lock.RLock()
	defer s.lock.RUnlock()
	_, isInList := s.games[gameID]
	return isInList
}

func (s *Scheduler) SynchronizeCurrentWeek(gameChan chan<- repository.Game) error {
	logger := s.logger.With().Str("method", "SynchronizeCurrentWeek").Logger()
	for {
		now := time.Now()
		startTime := general.StartTime(now)
		endTime := general.EndTime(now)

		logger.Info().Msgf("getting games for the week between %s - %s", startTime, endTime)

		games, err := s.controller.GetGamesBetweenDates(startTime, endTime)
		if err != nil {
			logger.Info().Err(err).Msgf("error synchronizing current week: %v")
		}
		for _, game := range games {
			logger.Debug().Fields(game).Msgf("adding game to channel")
			gameChan <- game
		}
		tomorrow := time.Now().UTC().AddDate(0, 0, 1)
		midNight := time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(), 0, 0, 0, 0, time.UTC)
		duration := midNight.Sub(time.Now().UTC())
		logger.Info().Msgf("updating weekly schedule in %v", duration)
		time.Sleep(duration)
	}
}

func (s *Scheduler) GetGameInfo(game repository.Game) {
	logger := s.logger.With().Str("method", "GetGameInfo").Logger()

	for {
		logger.Info().Msgf("getting game info for event: %v - %s vs. %s", game, game.AwayTeam, game.HomeTeam)

		info, err := s.controller.GetGameInfo(game.ID)
		if err != nil {
			logger.Error().Err(err).Msgf("error getting game info for event: %v - %s vs. %s", game, game.AwayTeam, game.HomeTeam)
		}
		switch info.Status.Desc {
		case "Final":
			s.controller.FinalizeGame(info)

			logger.Info().Msgf("game %s - %s vs. %s ended. Exiting get game info", game.ID, game.AwayTeam, game.HomeTeam)
			s.RemoveGame(game.ID)
			return
		case "Scheduled":
			sleepDuration, err := calculateSleepTimeInSeconds(logger, info.Status.Det)
			if err != nil {
				logger.Info().Err(err).Msgf("error calculating sleep time")
				sleepDuration = time.Second
			}
			logger.Info().Msgf("Sleeping game: %s - %s vs. %s for %s", game.ID, game.AwayTeam, game.HomeTeam, sleepDuration)
			time.Sleep(sleepDuration)
		default:
			s.controller.UpdateGame(info)
		}
		time.Sleep(time.Second)
	}
}

func calculateSleepTimeInSeconds(logger zerolog.Logger, dateTime string) (time.Duration, error) {
	layout := "1/2/2006 - 3:04 PM MST"
	dateParts := strings.Split(dateTime, " - ")
	dateTime = dateParts[0] + "/" + strconv.Itoa(time.Now().Year()) + " - " + dateParts[1]
	// Parse the input string using the layout
	t, err := time.Parse(layout, dateTime)
	if err != nil {
		return 0, err
	}

	logger.Info().Str("time", t.String()).Msg("time from game")
	return t.Sub(time.Now().UTC().Add(-4 * time.Hour)), nil
}
