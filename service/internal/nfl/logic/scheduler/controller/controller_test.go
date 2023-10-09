package controller

import (
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	repository2 "github.com/rmarken5/mini-score/service/internal/nfl/logic/internal/data-access/db/repository"
	"github.com/rmarken5/mini-score/service/internal/nfl/logic/internal/data-access/http/scraper"
	"github.com/rmarken5/mini-score/service/internal/nfl/logic/scheduler/controller/internal/mocks"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

//go:generate mockgen -destination ./internal/mocks/repository_mock.go -package mocks -source=../../../../nfl/data-access/db/repository/repository.go Repository
//go:generate mockgen -destination ./internal/mocks/schedule_scraper_mock.go -package mocks  -source=../../../../nfl/data-access/http/scraper/scraper.go ScheduleScraper
//go:generate mockgen -destination ./internal/mocks/schedule_requestor_mock.go -package mocks  -source=../../../../nfl/data-access/http/rest/requester.go Requester

func TestLogic_KeepScheduleSynchronized(t *testing.T) {

	var (
		repoGameOne = repository2.Game{
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
		mockScraper   func(ctrl *gomock.Controller) *mocks.MockScheduleScraper
		mockRepo      func(ctrl *gomock.Controller) *mocks.MockRepository
		mockRequester func(ctrl *gomock.Controller) *mocks.MockRequester
	}{
		"should insert game when game doesn't exist": {
			mockScraper: func(ctrl *gomock.Controller) *mocks.MockScheduleScraper {
				mockScraper := mocks.NewMockScheduleScraper(ctrl)
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
			mockRepo: func(ctrl *gomock.Controller) *mocks.MockRepository {
				mockRepo := mocks.NewMockRepository(ctrl)
				mockRepo.EXPECT().GetGame(gomock.Any()).Return(repository2.Game{}, repository2.ErrNoGame).AnyTimes()
				mockRepo.EXPECT().InsertGame(gomock.Any()).Return(nil).AnyTimes()
				mockRepo.EXPECT().GetTeamByAbv(gomock.Any()).Return(&repository2.Team{ID: uuid.New()}, nil).AnyTimes()
				return mockRepo
			},
		},
		"should continue when insert fails": {
			mockScraper: func(ctrl *gomock.Controller) *mocks.MockScheduleScraper {
				mockScraper := mocks.NewMockScheduleScraper(ctrl)
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
			mockRepo: func(ctrl *gomock.Controller) *mocks.MockRepository {
				mockRepo := mocks.NewMockRepository(ctrl)
				mockRepo.EXPECT().GetGame(gomock.Any()).Return(repository2.Game{}, repository2.ErrNoGame).AnyTimes()
				mockRepo.EXPECT().InsertGame(gomock.Any()).Return(errors.New("error")).AnyTimes()
				mockRepo.EXPECT().GetTeamByAbv(gomock.Any()).Return(&repository2.Team{ID: uuid.New()}, nil).AnyTimes()
				return mockRepo
			},
		},
		"should continue when update when game exists": {
			mockScraper: func(ctrl *gomock.Controller) *mocks.MockScheduleScraper {
				mockScraper := mocks.NewMockScheduleScraper(ctrl)
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
			mockRepo: func(ctrl *gomock.Controller) *mocks.MockRepository {
				mockRepo := mocks.NewMockRepository(ctrl)
				mockRepo.EXPECT().GetGame(gomock.Any()).Return(repoGameOne, nil).AnyTimes()
				mockRepo.EXPECT().UpdateGameTime(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
				mockRepo.EXPECT().UpdateQuarterScore(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
				mockRepo.EXPECT().GetTeamByAbv(gomock.Any()).Return(&repository2.Team{ID: uuid.New()}, nil).AnyTimes()
				return mockRepo
			},
		},
		"should continue when update fails": {
			mockScraper: func(ctrl *gomock.Controller) *mocks.MockScheduleScraper {
				mockScraper := mocks.NewMockScheduleScraper(ctrl)
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
			mockRepo: func(ctrl *gomock.Controller) *mocks.MockRepository {
				mockRepo := mocks.NewMockRepository(ctrl)
				mockRepo.EXPECT().GetGame(gomock.Any()).Return(repoGameOne, nil).AnyTimes()
				mockRepo.EXPECT().UpdateGameTime(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
				mockRepo.EXPECT().UpdateQuarterScore(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("error")).AnyTimes()
				mockRepo.EXPECT().GetTeamByAbv(gomock.Any()).Return(&repository2.Team{ID: uuid.New()}, nil).AnyTimes()
				return mockRepo
			},
		},
		"should continue when getGame fails": {
			mockScraper: func(ctrl *gomock.Controller) *mocks.MockScheduleScraper {
				mockScraper := mocks.NewMockScheduleScraper(ctrl)
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
			mockRepo: func(ctrl *gomock.Controller) *mocks.MockRepository {
				mockRepo := mocks.NewMockRepository(ctrl)
				mockRepo.EXPECT().GetGame(gomock.Any()).Return(repository2.Game{}, errors.New("error")).AnyTimes()
				return mockRepo
			},
		},
		"should continue when FetchGamesForWeeks fails": {
			mockScraper: func(ctrl *gomock.Controller) *mocks.MockScheduleScraper {
				mockScraper := mocks.NewMockScheduleScraper(ctrl)
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
			mockScraper: func(ctrl *gomock.Controller) *mocks.MockScheduleScraper {
				mockScraper := mocks.NewMockScheduleScraper(ctrl)
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
			var mockRepository *mocks.MockRepository
			if tc.mockRepo != nil {
				mockRepository = tc.mockRepo(ctrl)
			}

			var mockRequester *mocks.MockRequester
			if tc.mockRequester != nil {
				mockRequester = tc.mockRequester(ctrl)
			}

			var mockScheduleScraper *mocks.MockScheduleScraper
			if tc.mockScraper != nil {
				mockScheduleScraper = tc.mockScraper(ctrl)
			}

			l := NewLogic(zerolog.New(os.Stdout), mockScheduleScraper, mockRepository, mockRequester)

			go func() {
				l.KeepScheduleSynchronized(exitChannel, time.Second)
			}()

			exitChannel <- true

		})
	}
}

func TestLogic_UpdateGame(t *testing.T) {

	testCases := map[string]struct {
		mockRepo     func(ctrl *gomock.Controller) *mocks.MockRepository
		haveGameInfo scraper.GameInfo
		haveLogic    func(mockRepository *mocks.MockRepository) *Logic
	}{
		"should insert scores for quarters when they don't exist": {
			mockRepo: func(ctrl *gomock.Controller) *mocks.MockRepository {
				mockRepository := mocks.NewMockRepository(ctrl)
				mockRepository.EXPECT().GetQuarterScoreBy(gomock.Any(), gomock.Any(), gomock.Any()).Return(repository2.GameQuarterScore{}, repository2.ErrNoQuarterScore).Times(2)
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
			haveLogic: func(mockRepository *mocks.MockRepository) *Logic {
				l := &Logic{
					logger: zerolog.Nop(), repo: mockRepository,
					gameTeamQuarterCache: make(map[string]int),
				}
				return l
			},
		},
		"should update scores for quarters when they exist": {
			mockRepo: func(ctrl *gomock.Controller) *mocks.MockRepository {
				mockRepository := mocks.NewMockRepository(ctrl)
				mockRepository.EXPECT().GetQuarterScoreBy(gomock.Any(), gomock.Any(), gomock.Any()).Return(repository2.GameQuarterScore{}, nil).Times(2)
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
			haveLogic: func(mockRepository *mocks.MockRepository) *Logic {
				l := &Logic{
					logger: zerolog.Nop(), repo: mockRepository,
					gameTeamQuarterCache: make(map[string]int),
				}
				return l
			},
		},
		"should skip processing when cache is current": {

			mockRepo: func(ctrl *gomock.Controller) *mocks.MockRepository {
				mockRepository := mocks.NewMockRepository(ctrl)
				mockRepository.EXPECT().GetQuarterScoreBy(gomock.Any(), gomock.Any(), gomock.Any()).Return(repository2.GameQuarterScore{}, nil).Times(0)
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
			haveLogic: func(mockRepository *mocks.MockRepository) *Logic {
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
			var mockRepo *mocks.MockRepository
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
		mockRepo      func(ctrl *gomock.Controller) *mocks.MockRepository
		expectedError error
	}{
		"should get games from between dates from repo": {
			mockRepo: func(ctrl *gomock.Controller) *mocks.MockRepository {
				mockRepo := mocks.NewMockRepository(ctrl)
				mockRepo.EXPECT().GetGames(gomock.Any(), gomock.Any()).Return([]repository2.Game{}, nil)
				return mockRepo
			},
		},
		"should log and return error from repo": {
			mockRepo: func(ctrl *gomock.Controller) *mocks.MockRepository {
				mockRepo := mocks.NewMockRepository(ctrl)
				mockRepo.EXPECT().GetGames(gomock.Any(), gomock.Any()).Return([]repository2.Game(nil), myErr)
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

			var mockRepo *mocks.MockRepository
			if tc.mockRepo != nil {
				mockRepo = tc.mockRepo(ctrl)
			}

			logic := NewLogic(zerolog.New(os.Stdout), nil, mockRepo, nil)

			_, err := logic.GetGamesBetweenDates(time.Now(), time.Now())
			assert.ErrorIs(t, err, tc.expectedError)
		})
	}
}

func TestLogic_GetGameInfo(t *testing.T) {
	var myErr = errors.New("error")
	testCases := map[string]struct {
		mockScraper   func(ctrl *gomock.Controller) *mocks.MockScheduleScraper
		expectedError error
	}{
		"should get gameInfo from scraper": {
			mockScraper: func(ctrl *gomock.Controller) *mocks.MockScheduleScraper {
				mockScheduleScraper := mocks.NewMockScheduleScraper(ctrl)
				mockScheduleScraper.EXPECT().FetchGameInfo(gomock.Any()).Return(scraper.GameInfo{}, nil)
				return mockScheduleScraper
			},
		},
		"should log and return error from scraper": {
			mockScraper: func(ctrl *gomock.Controller) *mocks.MockScheduleScraper {
				mockScheduleScraper := mocks.NewMockScheduleScraper(ctrl)
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

			var mockScraper *mocks.MockScheduleScraper
			if tc.mockScraper != nil {
				mockScraper = tc.mockScraper(ctrl)
			}

			logic := NewLogic(zerolog.New(os.Stdout), mockScraper, nil, nil)

			_, err := logic.GetGameInfo(uuid.NewString())
			assert.ErrorIs(t, err, tc.expectedError)
		})
	}
}
