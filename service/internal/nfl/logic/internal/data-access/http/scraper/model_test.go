package scraper

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBySeasonType_Len(t *testing.T) {
	weeks := BySeasonType{
		{Year: 2021, SeasonType: PreSeason, WeekNumber: 1},
		{Year: 2021, SeasonType: RegSeason, WeekNumber: 2},
		{Year: 2021, SeasonType: PostSeason, WeekNumber: 3},
	}

	assert.Equal(t, 3, weeks.Len(), "Length should be 3")
}

func TestBySeasonType_Swap(t *testing.T) {
	weeks := BySeasonType{
		{Year: 2021, SeasonType: PreSeason, WeekNumber: 2},
		{Year: 2021, SeasonType: RegSeason, WeekNumber: 1},
	}

	weeks.Swap(0, 1)

	assert.Equal(t, RegSeason, weeks[0].SeasonType, "First element should be RegSeason")
	assert.Equal(t, PreSeason, weeks[1].SeasonType, "Second element should be PreSeason")
}

func TestBySeasonType_Less(t *testing.T) {
	weeks := BySeasonType{
		{Year: 2020, SeasonType: RegSeason, WeekNumber: 3},
		{Year: 2021, SeasonType: RegSeason, WeekNumber: 3},
		{Year: 2020, SeasonType: RegSeason, WeekNumber: 3},
		{Year: 2020, SeasonType: PostSeason, WeekNumber: 3},
		{Year: 2020, SeasonType: RegSeason, WeekNumber: 3},
		{Year: 2020, SeasonType: RegSeason, WeekNumber: 4},
	}

	less := weeks.Less(0, 1)
	assert.True(t, less, "First element should be less than the second element")

	less = weeks.Less(2, 3)
	assert.True(t, less, "First element should be less than the second element")

	less = weeks.Less(4, 5)
	assert.True(t, less, "First element should be less than the second element")

}

func TestBySeasonType_Sort(t *testing.T) {
	weeks := BySeasonType{
		{Year: 2021, SeasonType: RegSeason, WeekNumber: 2},
		{Year: 2020, SeasonType: PreSeason, WeekNumber: 1},
		{Year: 2021, SeasonType: PostSeason, WeekNumber: 3},
	}

	sort.Sort(weeks)

	assert.Equal(t, PreSeason, weeks[0].SeasonType, "First element should be PreSeason")
	assert.Equal(t, RegSeason, weeks[1].SeasonType, "Second element should be RegSeason")
	assert.Equal(t, PostSeason, weeks[2].SeasonType, "Third element should be PostSeason")
}
