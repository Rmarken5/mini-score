package controller

import (
	"errors"
	"github.com/google/uuid"
	"github.com/rmarken5/mini-score/service/internal/nfl/logic/internal/data-access/db/repository"
	"github.com/rmarken5/mini-score/service/internal/nfl/logic/internal/data-access/http/rest"
	"github.com/rmarken5/mini-score/service/internal/nfl/logic/internal/data-access/http/scraper"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"os"
	"sync"
	"testing"
	"time"
)

func TestLogic_KeepScheduleSynchronized(t *testing.T) {

	var (
		repoGameOne = repository.Game{
			ID:        uuid.NewString(),
			GameTime:  time.Now(),
			Quarter:   "",
			GameClock: "",
			AwayTeam:  uuid.NewString(),
			HomeTeam:  uuid.NewString(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			DeletedAt: nil,
		}
	)

	testCases := map[string]struct {
		mockScraper   func(ctrl *gomock.Controller) *scraper.MockScheduleScraper
		mockRepo      func(ctrl *gomock.Controller) *repository.MockRepository
		mockRequester func(ctrl *gomock.Controller) *rest.MockRequester
	}{
		"should insert game when game doesn't exist": {
			mockScraper: func(ctrl *gomock.Controller) *scraper.MockScheduleScraper {
				mockScraper := scraper.NewMockScheduleScraper(ctrl)
				mockScraper.EXPECT().FetchSchedule().Return([]scraper.Week{{
					Text:  uuid.NewString(),
					Label: uuid.NewString(),
					StartDate: scraper.CustomTime{
						Time: time.Now(),
					},
					EndDate: scraper.CustomTime{
						Time: time.Now(),
					},
					SeasonType: 1,
					WeekNumber: 1,
					Year:       2023,
					URL:        uuid.NewString(),
					IsActive:   true,
				},
					{
						Text:  uuid.NewString(),
						Label: uuid.NewString(),
						StartDate: scraper.CustomTime{
							Time: time.Now(),
						},
						EndDate: scraper.CustomTime{
							Time: time.Now(),
						},
						SeasonType: 1,
						WeekNumber: 1,
						Year:       2023,
						URL:        uuid.NewString(),
						IsActive:   true,
					},
					{
						Text:  uuid.NewString(),
						Label: uuid.NewString(),
						StartDate: scraper.CustomTime{
							time.Now(),
						},
						EndDate: scraper.CustomTime{
							Time: time.Now(),
						},
						SeasonType: 1,
						WeekNumber: 1,
						Year:       2023,
						URL:        uuid.NewString(),
						IsActive:   true,
					},
					{
						Text:  uuid.NewString(),
						Label: uuid.NewString(),
						StartDate: scraper.CustomTime{
							Time: time.Now(),
						},
						EndDate: scraper.CustomTime{
							Time: time.Now(),
						},
						SeasonType: 1,
						WeekNumber: 1,
						Year:       2023,
						URL:        uuid.NewString(),
						IsActive:   true,
					},
					{
						Text:  uuid.NewString(),
						Label: uuid.NewString(),
						StartDate: scraper.CustomTime{
							Time: time.Now(),
						},
						EndDate: scraper.CustomTime{
							Time: time.Now(),
						},
						SeasonType: 1,
						WeekNumber: 1,
						Year:       2023,
						URL:        uuid.NewString(),
						IsActive:   true,
					},
				}, nil).AnyTimes()
				mockScraper.EXPECT().FetchGamesForWeeks(gomock.Any()).Return(scraper.Games{
					"20230911": {
						scraper.Game{
							ID:          "12345",
							Competitors: []scraper.Competitor{{Abbrev: "SF"}, {Abbrev: "PIT"}},
							Date:        "2023-09-10T13:00Z",
							TBD:         false,
							Completed:   false,
							Link:        "",
							Teams:       nil,
							IsTie:       false,
							TimeValid:   false,
						},
					},
				}, nil).AnyTimes()

				return mockScraper
			},
			mockRepo: func(ctrl *gomock.Controller) *repository.MockRepository {
				mockRepo := repository.NewMockRepository(ctrl)
				mockRepo.EXPECT().GetGame(gomock.Any()).Return(repository.Game{}, repository.ErrNoGame).AnyTimes()
				mockRepo.EXPECT().InsertGame(gomock.Any()).Return(nil).AnyTimes()
				mockRepo.EXPECT().GetTeamByAbv(gomock.Any()).Return(&repository.Team{ID: uuid.New()}, nil).AnyTimes()
				return mockRepo
			},
		},
		"should continue when insert fails": {
			mockScraper: func(ctrl *gomock.Controller) *scraper.MockScheduleScraper {
				mockScraper := scraper.NewMockScheduleScraper(ctrl)
				mockScraper.EXPECT().FetchSchedule().Return([]scraper.Week{{
					Text:  uuid.NewString(),
					Label: uuid.NewString(),
					StartDate: scraper.CustomTime{
						Time: time.Now(),
					},
					EndDate: scraper.CustomTime{
						Time: time.Now(),
					},
					SeasonType: 1,
					WeekNumber: 1,
					Year:       2023,
					URL:        uuid.NewString(),
					IsActive:   true,
				},
					{
						Text:  uuid.NewString(),
						Label: uuid.NewString(),
						StartDate: scraper.CustomTime{
							Time: time.Now(),
						},
						EndDate: scraper.CustomTime{
							Time: time.Now(),
						},
						SeasonType: 1,
						WeekNumber: 1,
						Year:       2023,
						URL:        uuid.NewString(),
						IsActive:   true,
					},
					{
						Text:  uuid.NewString(),
						Label: uuid.NewString(),
						StartDate: scraper.CustomTime{
							time.Now(),
						},
						EndDate: scraper.CustomTime{
							Time: time.Now(),
						},
						SeasonType: 1,
						WeekNumber: 1,
						Year:       2023,
						URL:        uuid.NewString(),
						IsActive:   true,
					},
					{
						Text:  uuid.NewString(),
						Label: uuid.NewString(),
						StartDate: scraper.CustomTime{
							Time: time.Now(),
						},
						EndDate: scraper.CustomTime{
							Time: time.Now(),
						},
						SeasonType: 1,
						WeekNumber: 1,
						Year:       2023,
						URL:        uuid.NewString(),
						IsActive:   true,
					},
					{
						Text:  uuid.NewString(),
						Label: uuid.NewString(),
						StartDate: scraper.CustomTime{
							Time: time.Now(),
						},
						EndDate: scraper.CustomTime{
							Time: time.Now(),
						},
						SeasonType: 1,
						WeekNumber: 1,
						Year:       2023,
						URL:        uuid.NewString(),
						IsActive:   true,
					},
				}, nil).AnyTimes()
				mockScraper.EXPECT().FetchGamesForWeeks(gomock.Any()).Return(scraper.Games{
					"20230911": {
						scraper.Game{
							ID:          "12345",
							Competitors: []scraper.Competitor{{Abbrev: "SF"}, {Abbrev: "PIT"}},
							Date:        "2023-09-10T13:00Z",
							TBD:         false,
							Completed:   false,
							Link:        "",
							Teams:       nil,
							IsTie:       false,
							TimeValid:   false,
						},
					},
				}, nil).AnyTimes()

				return mockScraper
			},
			mockRepo: func(ctrl *gomock.Controller) *repository.MockRepository {
				mockRepo := repository.NewMockRepository(ctrl)
				mockRepo.EXPECT().GetGame(gomock.Any()).Return(repository.Game{}, repository.ErrNoGame).AnyTimes()
				mockRepo.EXPECT().InsertGame(gomock.Any()).Return(errors.New("error")).AnyTimes()
				mockRepo.EXPECT().GetTeamByAbv(gomock.Any()).Return(&repository.Team{ID: uuid.New()}, nil).AnyTimes()
				return mockRepo
			},
		},
		"should continue when update when game exists": {
			mockScraper: func(ctrl *gomock.Controller) *scraper.MockScheduleScraper {
				mockScraper := scraper.NewMockScheduleScraper(ctrl)
				mockScraper.EXPECT().FetchSchedule().Return([]scraper.Week{{
					Text:  uuid.NewString(),
					Label: uuid.NewString(),
					StartDate: scraper.CustomTime{
						Time: time.Now(),
					},
					EndDate: scraper.CustomTime{
						Time: time.Now(),
					},
					SeasonType: 1,
					WeekNumber: 1,
					Year:       2023,
					URL:        uuid.NewString(),
					IsActive:   true,
				},
					{
						Text:  uuid.NewString(),
						Label: uuid.NewString(),
						StartDate: scraper.CustomTime{
							Time: time.Now(),
						},
						EndDate: scraper.CustomTime{
							Time: time.Now(),
						},
						SeasonType: 1,
						WeekNumber: 1,
						Year:       2023,
						URL:        uuid.NewString(),
						IsActive:   true,
					},
					{
						Text:  uuid.NewString(),
						Label: uuid.NewString(),
						StartDate: scraper.CustomTime{
							time.Now(),
						},
						EndDate: scraper.CustomTime{
							Time: time.Now(),
						},
						SeasonType: 1,
						WeekNumber: 1,
						Year:       2023,
						URL:        uuid.NewString(),
						IsActive:   true,
					},
					{
						Text:  uuid.NewString(),
						Label: uuid.NewString(),
						StartDate: scraper.CustomTime{
							Time: time.Now(),
						},
						EndDate: scraper.CustomTime{
							Time: time.Now(),
						},
						SeasonType: 1,
						WeekNumber: 1,
						Year:       2023,
						URL:        uuid.NewString(),
						IsActive:   true,
					},
					{
						Text:  uuid.NewString(),
						Label: uuid.NewString(),
						StartDate: scraper.CustomTime{
							Time: time.Now(),
						},
						EndDate: scraper.CustomTime{
							Time: time.Now(),
						},
						SeasonType: 1,
						WeekNumber: 1,
						Year:       2023,
						URL:        uuid.NewString(),
						IsActive:   true,
					},
				}, nil).AnyTimes()
				mockScraper.EXPECT().FetchGamesForWeeks(gomock.Any()).Return(scraper.Games{
					"20230911": {
						scraper.Game{
							ID:          "12345",
							Competitors: []scraper.Competitor{{Abbrev: "SF"}, {Abbrev: "PIT"}},
							Date:        "2023-09-10T13:00Z",
							TBD:         false,
							Completed:   false,
							Link:        "",
							Teams:       nil,
							IsTie:       false,
							TimeValid:   false,
						},
					},
				}, nil).AnyTimes()

				return mockScraper
			},
			mockRepo: func(ctrl *gomock.Controller) *repository.MockRepository {
				mockRepo := repository.NewMockRepository(ctrl)
				mockRepo.EXPECT().GetGame(gomock.Any()).Return(repoGameOne, nil).AnyTimes()
				mockRepo.EXPECT().UpdateGameTime(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
				mockRepo.EXPECT().UpdateQuarterScore(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
				mockRepo.EXPECT().GetTeamByAbv(gomock.Any()).Return(&repository.Team{ID: uuid.New()}, nil).AnyTimes()
				return mockRepo
			},
		},
		"should continue when update fails": {
			mockScraper: func(ctrl *gomock.Controller) *scraper.MockScheduleScraper {
				mockScraper := scraper.NewMockScheduleScraper(ctrl)
				mockScraper.EXPECT().FetchSchedule().Return([]scraper.Week{{
					Text:  uuid.NewString(),
					Label: uuid.NewString(),
					StartDate: scraper.CustomTime{
						Time: time.Now(),
					},
					EndDate: scraper.CustomTime{
						Time: time.Now(),
					},
					SeasonType: 1,
					WeekNumber: 1,
					Year:       2023,
					URL:        uuid.NewString(),
					IsActive:   true,
				},
					{
						Text:  uuid.NewString(),
						Label: uuid.NewString(),
						StartDate: scraper.CustomTime{
							Time: time.Now(),
						},
						EndDate: scraper.CustomTime{
							Time: time.Now(),
						},
						SeasonType: 1,
						WeekNumber: 1,
						Year:       2023,
						URL:        uuid.NewString(),
						IsActive:   true,
					},
					{
						Text:  uuid.NewString(),
						Label: uuid.NewString(),
						StartDate: scraper.CustomTime{
							time.Now(),
						},
						EndDate: scraper.CustomTime{
							Time: time.Now(),
						},
						SeasonType: 1,
						WeekNumber: 1,
						Year:       2023,
						URL:        uuid.NewString(),
						IsActive:   true,
					},
					{
						Text:  uuid.NewString(),
						Label: uuid.NewString(),
						StartDate: scraper.CustomTime{
							Time: time.Now(),
						},
						EndDate: scraper.CustomTime{
							Time: time.Now(),
						},
						SeasonType: 1,
						WeekNumber: 1,
						Year:       2023,
						URL:        uuid.NewString(),
						IsActive:   true,
					},
					{
						Text:  uuid.NewString(),
						Label: uuid.NewString(),
						StartDate: scraper.CustomTime{
							Time: time.Now(),
						},
						EndDate: scraper.CustomTime{
							Time: time.Now(),
						},
						SeasonType: 1,
						WeekNumber: 1,
						Year:       2023,
						URL:        uuid.NewString(),
						IsActive:   true,
					},
				}, nil).AnyTimes()
				mockScraper.EXPECT().FetchGamesForWeeks(gomock.Any()).Return(scraper.Games{
					"20230911": {
						scraper.Game{
							ID:          "12345",
							Competitors: []scraper.Competitor{{Abbrev: "SF"}, {Abbrev: "PIT"}},
							Date:        "2023-09-10T13:00Z",
							TBD:         false,
							Completed:   false,
							Link:        "",
							Teams:       nil,
							IsTie:       false,
							TimeValid:   false,
						},
					},
				}, nil).AnyTimes()

				return mockScraper
			},
			mockRepo: func(ctrl *gomock.Controller) *repository.MockRepository {
				mockRepo := repository.NewMockRepository(ctrl)
				mockRepo.EXPECT().GetGame(gomock.Any()).Return(repoGameOne, nil).AnyTimes()
				mockRepo.EXPECT().UpdateGameTime(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
				mockRepo.EXPECT().UpdateQuarterScore(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("error")).AnyTimes()
				mockRepo.EXPECT().GetTeamByAbv(gomock.Any()).Return(&repository.Team{ID: uuid.New()}, nil).AnyTimes()
				return mockRepo
			},
		},
		"should continue when getGame fails": {
			mockScraper: func(ctrl *gomock.Controller) *scraper.MockScheduleScraper {
				mockScraper := scraper.NewMockScheduleScraper(ctrl)
				mockScraper.EXPECT().FetchSchedule().Return([]scraper.Week{{
					Text:  uuid.NewString(),
					Label: uuid.NewString(),
					StartDate: scraper.CustomTime{
						Time: time.Now(),
					},
					EndDate: scraper.CustomTime{
						Time: time.Now(),
					},
					SeasonType: 1,
					WeekNumber: 1,
					Year:       2023,
					URL:        uuid.NewString(),
					IsActive:   true,
				},
					{
						Text:  uuid.NewString(),
						Label: uuid.NewString(),
						StartDate: scraper.CustomTime{
							Time: time.Now(),
						},
						EndDate: scraper.CustomTime{
							Time: time.Now(),
						},
						SeasonType: 1,
						WeekNumber: 1,
						Year:       2023,
						URL:        uuid.NewString(),
						IsActive:   true,
					},
					{
						Text:  uuid.NewString(),
						Label: uuid.NewString(),
						StartDate: scraper.CustomTime{
							time.Now(),
						},
						EndDate: scraper.CustomTime{
							Time: time.Now(),
						},
						SeasonType: 1,
						WeekNumber: 1,
						Year:       2023,
						URL:        uuid.NewString(),
						IsActive:   true,
					},
					{
						Text:  uuid.NewString(),
						Label: uuid.NewString(),
						StartDate: scraper.CustomTime{
							Time: time.Now(),
						},
						EndDate: scraper.CustomTime{
							Time: time.Now(),
						},
						SeasonType: 1,
						WeekNumber: 1,
						Year:       2023,
						URL:        uuid.NewString(),
						IsActive:   true,
					},
					{
						Text:  uuid.NewString(),
						Label: uuid.NewString(),
						StartDate: scraper.CustomTime{
							Time: time.Now(),
						},
						EndDate: scraper.CustomTime{
							Time: time.Now(),
						},
						SeasonType: 1,
						WeekNumber: 1,
						Year:       2023,
						URL:        uuid.NewString(),
						IsActive:   true,
					},
				}, nil).AnyTimes()
				mockScraper.EXPECT().FetchGamesForWeeks(gomock.Any()).Return(scraper.Games{
					"20230911": {
						scraper.Game{
							ID:          "12345",
							Competitors: []scraper.Competitor{{Abbrev: "SF"}, {Abbrev: "PIT"}},
							Date:        "2023-09-10T13:00Z",
							TBD:         false,
							Completed:   false,
							Link:        "",
							Teams:       nil,
							IsTie:       false,
							TimeValid:   false,
						},
					},
				}, nil).AnyTimes()

				return mockScraper
			},
			mockRepo: func(ctrl *gomock.Controller) *repository.MockRepository {
				mockRepo := repository.NewMockRepository(ctrl)
				mockRepo.EXPECT().GetGame(gomock.Any()).Return(repository.Game{}, errors.New("error")).AnyTimes()
				return mockRepo
			},
		},
		"should continue when FetchGamesForWeeks fails": {
			mockScraper: func(ctrl *gomock.Controller) *scraper.MockScheduleScraper {
				mockScraper := scraper.NewMockScheduleScraper(ctrl)
				mockScraper.EXPECT().FetchSchedule().Return([]scraper.Week{{
					Text:  uuid.NewString(),
					Label: uuid.NewString(),
					StartDate: scraper.CustomTime{
						Time: time.Now(),
					},
					EndDate: scraper.CustomTime{
						Time: time.Now(),
					},
					SeasonType: 1,
					WeekNumber: 1,
					Year:       2023,
					URL:        uuid.NewString(),
					IsActive:   true,
				},
					{
						Text:  uuid.NewString(),
						Label: uuid.NewString(),
						StartDate: scraper.CustomTime{
							Time: time.Now(),
						},
						EndDate: scraper.CustomTime{
							Time: time.Now(),
						},
						SeasonType: 1,
						WeekNumber: 1,
						Year:       2023,
						URL:        uuid.NewString(),
						IsActive:   true,
					},
					{
						Text:  uuid.NewString(),
						Label: uuid.NewString(),
						StartDate: scraper.CustomTime{
							time.Now(),
						},
						EndDate: scraper.CustomTime{
							Time: time.Now(),
						},
						SeasonType: 1,
						WeekNumber: 1,
						Year:       2023,
						URL:        uuid.NewString(),
						IsActive:   true,
					},
					{
						Text:  uuid.NewString(),
						Label: uuid.NewString(),
						StartDate: scraper.CustomTime{
							Time: time.Now(),
						},
						EndDate: scraper.CustomTime{
							Time: time.Now(),
						},
						SeasonType: 1,
						WeekNumber: 1,
						Year:       2023,
						URL:        uuid.NewString(),
						IsActive:   true,
					},
					{
						Text:  uuid.NewString(),
						Label: uuid.NewString(),
						StartDate: scraper.CustomTime{
							Time: time.Now(),
						},
						EndDate: scraper.CustomTime{
							Time: time.Now(),
						},
						SeasonType: 1,
						WeekNumber: 1,
						Year:       2023,
						URL:        uuid.NewString(),
						IsActive:   true,
					},
				}, nil).AnyTimes()
				mockScraper.EXPECT().FetchGamesForWeeks(gomock.Any()).Return(make(scraper.Games), errors.New("error")).AnyTimes()

				return mockScraper
			},
		},
		"should continue when FetchSchedule fails": {
			mockScraper: func(ctrl *gomock.Controller) *scraper.MockScheduleScraper {
				mockScraper := scraper.NewMockScheduleScraper(ctrl)
				mockScraper.EXPECT().FetchSchedule().Return([]scraper.Week(nil), errors.New("error")).AnyTimes()

				return mockScraper
			},
		},
	}

	for name, tc := range testCases {
		name := name
		tc := tc
		t.Run(name, func(t *testing.T) {
			exitChannel := make(chan bool)
			ctrl := gomock.NewController(t)
			var mockRepository *repository.MockRepository
			if tc.mockRepo != nil {
				mockRepository = tc.mockRepo(ctrl)
			}

			var mockRequester *rest.MockRequester
			if tc.mockRequester != nil {
				mockRequester = tc.mockRequester(ctrl)
			}

			var mockScheduleScraper *scraper.MockScheduleScraper
			if tc.mockScraper != nil {
				mockScheduleScraper = tc.mockScraper(ctrl)
			}

			l := Logic{
				logger:               zerolog.New(os.Stdout),
				scrapper:             mockScheduleScraper,
				repo:                 mockRepository,
				requester:            mockRequester,
				gameTeamQuarterCache: nil,
				scoreCacheLock:       sync.RWMutex{},
				clockCache:           nil,
				clockCacheLock:       sync.RWMutex{},
			}

			go func() {
				l.KeepScheduleSynchronized(exitChannel, time.Second)
			}()

			exitChannel <- true

		})
	}
}

func TestLogic_UpdateGame(t *testing.T) {

	testCases := map[string]struct {
		mockRepo     func(ctrl *gomock.Controller) *repository.MockRepository
		haveGameInfo scraper.GameInfo
		haveLogic    func(mockRepository *repository.MockRepository) *Logic
	}{
		"should insert scores for quarters when they don't exist": {
			mockRepo: func(ctrl *gomock.Controller) *repository.MockRepository {
				mockRepository := repository.NewMockRepository(ctrl)
				mockRepository.EXPECT().GetQuarterScoreBy(gomock.Any(), gomock.Any(), gomock.Any()).Return(repository.GameQuarterScore{}, repository.ErrNoQuarterScore).Times(2)
				mockRepository.EXPECT().InsertQuarterScore(gomock.Any()).Return(nil).Times(2)

				return mockRepository
			},
			haveGameInfo: scraper.GameInfo{
				GameID:     "123",
				SeasonType: 0,
				Status: scraper.Status{
					Desc:  "",
					Det:   "",
					ID:    "",
					State: "",
				},
				StatusState: "",
				Tbd:         false,
				Tms: []scraper.Tms{
					{
						Abbrev:           "ABC",
						DisplayName:      "",
						ShortDisplayName: "",
						Records:          nil,
						IsHome:           false,
						Linescores: []scraper.Linescores{
							{DisplayValue: "7"},
						},
						Score:  "",
						Winner: false,
					},
					{
						Abbrev:           "DEF",
						DisplayName:      "",
						ShortDisplayName: "",
						Records:          nil,
						IsHome:           false,
						Linescores: []scraper.Linescores{
							{DisplayValue: "10"},
						},
						Score:  "",
						Winner: false,
					},
				},
			},
			haveLogic: func(mockRepository *repository.MockRepository) *Logic {
				l := &Logic{
					logger: zerolog.Nop(), repo: mockRepository,
					gameTeamQuarterCache: make(map[string]int),
				}
				return l
			},
		},
		"should update scores for quarters when they exist": {
			mockRepo: func(ctrl *gomock.Controller) *repository.MockRepository {
				mockRepository := repository.NewMockRepository(ctrl)
				mockRepository.EXPECT().GetQuarterScoreBy(gomock.Any(), gomock.Any(), gomock.Any()).Return(repository.GameQuarterScore{}, nil).Times(2)
				mockRepository.EXPECT().UpdateQuarterScore(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).Times(2)

				return mockRepository
			},
			haveGameInfo: scraper.GameInfo{
				GameID:     "123",
				SeasonType: 0,
				Status: scraper.Status{
					Desc:  "",
					Det:   "",
					ID:    "",
					State: "",
				},
				StatusState: "",
				Tbd:         false,
				Tms: []scraper.Tms{
					{
						Abbrev:           "ABC",
						DisplayName:      "",
						ShortDisplayName: "",
						Records:          nil,
						IsHome:           false,
						Linescores: []scraper.Linescores{
							{DisplayValue: "7"},
						},
						Score:  "",
						Winner: false,
					},
					{
						Abbrev:           "DEF",
						DisplayName:      "",
						ShortDisplayName: "",
						Records:          nil,
						IsHome:           false,
						Linescores: []scraper.Linescores{
							{DisplayValue: "10"},
						},
						Score:  "",
						Winner: false,
					},
				},
			},
			haveLogic: func(mockRepository *repository.MockRepository) *Logic {
				l := &Logic{
					logger: zerolog.Nop(), repo: mockRepository,
					gameTeamQuarterCache: make(map[string]int),
				}
				return l
			},
		},
		"should skip processing when cache is current": {

			mockRepo: func(ctrl *gomock.Controller) *repository.MockRepository {
				mockRepository := repository.NewMockRepository(ctrl)
				mockRepository.EXPECT().GetQuarterScoreBy(gomock.Any(), gomock.Any(), gomock.Any()).Return(repository.GameQuarterScore{}, nil).Times(0)
				mockRepository.EXPECT().UpdateQuarterScore(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).Times(0)

				return mockRepository
			},
			haveGameInfo: scraper.GameInfo{
				GameID:     "123",
				SeasonType: 0,
				Status: scraper.Status{
					Desc:  "",
					Det:   "",
					ID:    "",
					State: "",
				},
				StatusState: "",
				Tbd:         false,
				Tms: []scraper.Tms{
					{
						Abbrev:           "ABC",
						DisplayName:      "",
						ShortDisplayName: "",
						Records:          nil,
						IsHome:           false,
						Linescores: []scraper.Linescores{
							{DisplayValue: "7"},
						},
						Score:  "",
						Winner: false,
					},
					{
						Abbrev:           "DEF",
						DisplayName:      "",
						ShortDisplayName: "",
						Records:          nil,
						IsHome:           false,
						Linescores: []scraper.Linescores{
							{DisplayValue: "10"},
						},
						Score:  "",
						Winner: false,
					},
				},
			},
			haveLogic: func(mockRepository *repository.MockRepository) *Logic {
				gameTeamQuarterCache := make(map[string]int)
				gameTeamQuarterCache["123ABC1"] = 7
				gameTeamQuarterCache["123DEF1"] = 10
				l := &Logic{
					logger: zerolog.Nop(), repo: mockRepository,
					gameTeamQuarterCache: gameTeamQuarterCache,
				}
				return l
			},
		},
	}

	for name, tc := range testCases {
		name := name
		tc := tc
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			var mockRepo *repository.MockRepository
			if tc.mockRepo != nil {
				mockRepo = tc.mockRepo(ctrl)
			}
			var logic *Logic
			if tc.haveLogic != nil {
				logic = tc.haveLogic(mockRepo)
			}
			logic.updateGameQuarterScore(tc.haveGameInfo)
		})
	}

}

func TestLogic_GetGamesBetweenDates(t *testing.T) {
	var myErr = errors.New("error")
	testCases := map[string]struct {
		mockRepo      func(ctrl *gomock.Controller) *repository.MockRepository
		expectedError error
	}{
		"should get games from between dates from repo": {
			mockRepo: func(ctrl *gomock.Controller) *repository.MockRepository {
				mockRepo := repository.NewMockRepository(ctrl)
				mockRepo.EXPECT().GetGames(gomock.Any(), gomock.Any()).Return([]repository.Game{}, nil)
				return mockRepo
			},
		},
		"should log and return error from repo": {
			mockRepo: func(ctrl *gomock.Controller) *repository.MockRepository {
				mockRepo := repository.NewMockRepository(ctrl)
				mockRepo.EXPECT().GetGames(gomock.Any(), gomock.Any()).Return([]repository.Game(nil), myErr)
				return mockRepo
			},
			expectedError: myErr,
		},
	}
	for name, tc := range testCases {
		name := name
		tc := tc
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			var mockRepo *repository.MockRepository
			if tc.mockRepo != nil {
				mockRepo = tc.mockRepo(ctrl)
			}

			l := Logic{
				logger:               zerolog.New(os.Stdout),
				scrapper:             nil,
				repo:                 mockRepo,
				requester:            nil,
				gameTeamQuarterCache: nil,
				scoreCacheLock:       sync.RWMutex{},
				clockCache:           nil,
				clockCacheLock:       sync.RWMutex{},
			}

			_, err := l.GetGamesBetweenDates(time.Now(), time.Now())
			assert.ErrorIs(t, err, tc.expectedError)
		})
	}
}

func TestLogic_GetGameInfo(t *testing.T) {
	var myErr = errors.New("error")
	testCases := map[string]struct {
		mockScraper   func(ctrl *gomock.Controller) *scraper.MockScheduleScraper
		expectedError error
	}{
		"should get gameInfo from scraper": {
			mockScraper: func(ctrl *gomock.Controller) *scraper.MockScheduleScraper {
				mockScheduleScraper := scraper.NewMockScheduleScraper(ctrl)
				mockScheduleScraper.EXPECT().FetchGameInfo(gomock.Any()).Return(scraper.GameInfo{}, nil)
				return mockScheduleScraper
			},
		},
		"should log and return error from scraper": {
			mockScraper: func(ctrl *gomock.Controller) *scraper.MockScheduleScraper {
				mockScheduleScraper := scraper.NewMockScheduleScraper(ctrl)
				mockScheduleScraper.EXPECT().FetchGameInfo(gomock.Any()).Return(scraper.GameInfo{}, myErr)
				return mockScheduleScraper
			},
			expectedError: myErr,
		},
	}
	for name, tc := range testCases {
		name := name
		tc := tc
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			var mockScraper *scraper.MockScheduleScraper
			if tc.mockScraper != nil {
				mockScraper = tc.mockScraper(ctrl)
			}

			l := Logic{
				logger:               zerolog.New(os.Stdout),
				scrapper:             mockScraper,
				repo:                 nil,
				requester:            nil,
				gameTeamQuarterCache: nil,
				scoreCacheLock:       sync.RWMutex{},
				clockCache:           nil,
				clockCacheLock:       sync.RWMutex{},
			}

			_, err := l.GetGameInfo(uuid.NewString())
			assert.ErrorIs(t, err, tc.expectedError)
		})
	}
}
