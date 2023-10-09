package scraper

import (
	"embed"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"testing"
)

//go:embed test-data
var testHTML embed.FS

func Test_FindWeeksFromBytes(t *testing.T) {

	data, err := testHTML.Open("test-data/test_schedule_data.html")
	require.NoError(t, err)
	defer data.Close()

	bytes, err := io.ReadAll(data)

	evnts, err := findWeeksFromBytes(bytes)
	require.NoError(t, err)

	assert.Greater(t, len(evnts), 1)
}

func Test_FindGamesFromBytes(t *testing.T) {
	data, err := testHTML.Open("test-data/test_game_date.html")
	require.NoError(t, err)
	defer data.Close()

	bytes, err := io.ReadAll(data)

	evnts, err := findGamesFromBytes(bytes)
	require.NoError(t, err)

	assert.Len(t, evnts, 1)
}

func Test_FindGameInfoFromBytes(t *testing.T) {
	data, err := testHTML.Open("test-data/test_game_info.html")
	require.NoError(t, err)
	defer data.Close()

	bytes, err := io.ReadAll(data)

	gameInfo, err := findGameInfoFromBytes(bytes)
	require.NoError(t, err)

	assert.Equal(t, gameInfo.Tms[0].Abbrev, "CLE")
	assert.Equal(t, gameInfo.Tms[1].Abbrev, "NYJ")

}
