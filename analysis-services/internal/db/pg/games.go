package pg

import (
	"api.cloud.io/internal/models"
	"context"
	"encoding/json"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"time"
)

func initGames(jsondir string) error {
	cnt, err := getCount(context.Background())
	if err != nil {
		return err
	}
	if cnt == 0 {
		err = addGames(jsondir)
	}
	return err
}
func init(){
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339})
}
func addGames(jsondir string) error {
	log.Info().Msgf("Try to add init games")
	bt := []byte(getJsondata())
	games := make([]models.Games, 0)
	err := json.Unmarshal(bt, &games)
	if err != nil {
		log.Error().Msgf("Faild to load init games %s",err)
		return err
	}
	log.Info().Msgf("found %d games",len(games))
	for _, game := range games {
		err = AddGame(game)
		if err != nil {
			return err
		}
	}
	return nil
}
func AddGame(game models.Games) error {
	const q = `INSERT INTO games (
		rank, name, platform, year, genre, publisher, na_sale, eu_sale, jp_sale, other_sale, global_sale)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`
	_, err := TheDB.Exec(
		q, game.Rank, game.Name, game.Platform, game.Year, game.Genre, game.Publisher, game.NASale, game.EUSale, game.JPSale, game.OtherSale, game.GlobalSale)
	if err != nil {
		return err
	}
	return nil
}

func getCount(ctx context.Context) (int, error) {
	var err error
	q := `
		SELECT
			count(*)
		FROM games
	`
	rows, err := Query(ctx, q)
	if err != nil {
		return -1, err
	}
	defer rows.Close()
	count := 0
	if rows.Next() {
		err = rows.Scan(&count)
	}
	return count, err
}

func GetGamesByRank(ctx context.Context, limit, offset int) ([]models.Games, error) {
	var err error
	q := `
		SELECT
			rank,
			name,
			platform,
			year,
			genre,
			publisher,
			na_sale,
			eu_sale,
			jp_sale,
			other_sale,
			global_sale
		FROM games
		ORDER BY rank
		LIMIT $1 OFFSET $2
	`
	rows, err := Query(ctx, q, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	games := make([]models.Games, 0)
	for rows.Next() {
		var game models.Games
		err = rows.Scan(&game.Rank, &game.Name, &game.Platform, &game.Year, &game.Genre, &game.Publisher, &game.NASale, &game.EUSale, &game.JPSale, &game.OtherSale, &game.GlobalSale)
		if err != nil {
			return nil, err
		}
		games = append(games, game)
	}
	return games, err
}

func GetGamesBySell(ctx context.Context, limit, offset int) ([]models.Games, error) {
	var err error
	q := `
		SELECT
			rank,
			name,
			platform,
			year,
			genre,
			publisher,
			na_sale,
			eu_sale,
			jp_sale,
			other_sale,
			global_sale
		FROM games
		WHERE eu_sale > na_sale
		ORDER BY rank
		LIMIT $1 OFFSET $2
	`
	rows, err := Query(ctx, q, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	games := make([]models.Games, 0)
	for rows.Next() {
		var game models.Games
		err = rows.Scan(&game.Rank, &game.Name, &game.Platform, &game.Year, &game.Genre, &game.Publisher, &game.NASale, &game.EUSale, &game.JPSale, &game.OtherSale, &game.GlobalSale)
		if err != nil {
			return nil, err
		}
		games = append(games, game)
	}
	return games, err
}

func GetPlatforms(ctx context.Context) ([]string, error) {
	var err error
	q := `
		SELECT
			platform
		FROM games
		GROUP BY platform
	`
	rows, err := Query(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	platforms := make([]string, 0)
	for rows.Next() {
		var platform string
		err = rows.Scan(&platform)
		if err != nil {
			return nil, err
		}
		platforms = append(platforms, platform)
	}
	return platforms, err
}

func GetGamesByPlatform(ctx context.Context, limit int, platform string) ([]models.Games, error) {
	var err error
	q := `
		SELECT
			rank,
			name,
			platform,
			year,
			genre,
			publisher,
			na_sale,
			eu_sale,
			jp_sale,
			other_sale,
			global_sale
		FROM games
		WHERE platform=$1
		ORDER BY rank
		LIMIT $2
	`
	rows, err := Query(ctx, q, platform, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	games := make([]models.Games, 0)
	for rows.Next() {
		var game models.Games
		err = rows.Scan(&game.Rank, &game.Name, &game.Platform, &game.Year, &game.Genre, &game.Publisher, &game.NASale, &game.EUSale, &game.JPSale, &game.OtherSale, &game.GlobalSale)
		if err != nil {
			return nil, err
		}
		games = append(games, game)
	}
	return games, err
}

func GetYears(ctx context.Context) ([]int, error) {
	var err error
	q := `
		SELECT
			year
		FROM games
		GROUP BY year
	`
	rows, err := Query(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	years := make([]int, 0)
	for rows.Next() {
		var year int
		err = rows.Scan(&year)
		if err != nil {
			return nil, err
		}
		years = append(years, year)
	}
	return years, err
}

func GetGamesByYear(ctx context.Context, limit int, year int) ([]models.Games, error) {
	var err error
	q := `
		SELECT
			rank,
			name,
			platform,
			year,
			genre,
			publisher,
			na_sale,
			eu_sale,
			jp_sale,
			other_sale,
			global_sale
		FROM games
		WHERE year=$1
		ORDER BY rank
		LIMIT $2
	`
	rows, err := Query(ctx, q, year, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	games := make([]models.Games, 0)
	for rows.Next() {
		var game models.Games
		err = rows.Scan(&game.Rank, &game.Name, &game.Platform, &game.Year, &game.Genre, &game.Publisher, &game.NASale, &game.EUSale, &game.JPSale, &game.OtherSale, &game.GlobalSale)
		if err != nil {
			return nil, err
		}
		games = append(games, game)
	}
	return games, err
}

func GetGenres(ctx context.Context) ([]string, error) {
	var err error
	q := `
		SELECT
			genre
		FROM games
		GROUP BY genre
	`
	rows, err := Query(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	genres := make([]string, 0)
	for rows.Next() {
		var genre string
		err = rows.Scan(&genre)
		if err != nil {
			return nil, err
		}
		genres = append(genres, genre)
	}
	return genres, err
}

func GetGamesByGenre(ctx context.Context, limit int, genre string) ([]models.Games, error) {
	var err error
	q := `
		SELECT
			rank,
			name,
			platform,
			year,
			genre,
			publisher,
			na_sale,
			eu_sale,
			jp_sale,
			other_sale,
			global_sale
		FROM games
		WHERE genre=$1
		ORDER BY rank
		LIMIT $2
	`
	rows, err := Query(ctx, q, genre, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	games := make([]models.Games, 0)
	for rows.Next() {
		var game models.Games
		err = rows.Scan(&game.Rank, &game.Name, &game.Platform, &game.Year, &game.Genre, &game.Publisher, &game.NASale, &game.EUSale, &game.JPSale, &game.OtherSale, &game.GlobalSale)
		if err != nil {
			return nil, err
		}
		games = append(games, game)
	}
	return games, err
}

func GetGamesByName(ctx context.Context, name string) ([]models.Games, error) {
	var err error
	q := `
		SELECT
			rank,
			name,
			platform,
			year,
			genre,
			publisher,
			na_sale,
			eu_sale,
			jp_sale,
			other_sale,
			global_sale
		FROM games
		WHERE name like $1
		ORDER BY rank
	`
	rows, err := Query(ctx, q, "%"+name+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	games := make([]models.Games, 0)
	for rows.Next() {
		var game models.Games
		err = rows.Scan(&game.Rank, &game.Name, &game.Platform, &game.Year, &game.Genre, &game.Publisher, &game.NASale, &game.EUSale, &game.JPSale, &game.OtherSale, &game.GlobalSale)
		if err != nil {
			return nil, err
		}
		games = append(games, game)
	}
	return games, err
}

func GetTopSellForYearByPlatform(ctx context.Context, platform string, year int) ([]models.Games, error) {
	var err error
	q := `
		SELECT
			rank,
			name,
			platform,
			year,
			genre,
			publisher,
			na_sale,
			eu_sale,
			jp_sale,
			other_sale,
			global_sale
		FROM games
		WHERE platform = $1 AND year=$2
		ORDER BY rank
		LIMIT 5
	`
	rows, err := Query(ctx, q, platform, year)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	games := make([]models.Games, 0)
	for rows.Next() {
		var game models.Games
		err = rows.Scan(&game.Rank, &game.Name, &game.Platform, &game.Year, &game.Genre, &game.Publisher, &game.NASale, &game.EUSale, &game.JPSale, &game.OtherSale, &game.GlobalSale)
		if err != nil {
			return nil, err
		}
		games = append(games, game)
	}
	return games, err
}

func GetTotalSellByGenre(ctx context.Context, start int, end int) (map[string]float64, error) {
	var err error
	q := `
		SELECT
			genre,
			sum(global_sale)
		FROM games
		WHERE  year>=$1 AND year<=$2
		GROUP BY genre
		LIMIT 5
	`
	rows, err := Query(ctx, q, start, end)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := make(map[string]float64)
	for rows.Next() {
		var genre string
		var totalSell float64
		err = rows.Scan(&genre, &totalSell)
		if err != nil {
			return nil, err
		}
		res[genre] = totalSell
	}
	return res, err
}

func GetTotalSellByPublisher(ctx context.Context, start int, end int, publisher string) (map[int]float64, error) {
	var err error
	q := `
		SELECT
			year,
			sum(global_sale)
		FROM games
		WHERE  year>=$1 AND year<=$2 AND publisher=$3
		GROUP BY year
	`
	rows, err := Query(ctx, q, start, end, publisher)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := make(map[int]float64)
	for rows.Next() {
		var year int
		var totalSell float64
		err = rows.Scan(&year, &totalSell)
		if err != nil {
			return nil, err
		}
		res[year] = totalSell
	}
	return res, err
}

func GetTotalSellByYear(ctx context.Context, start int, end int) (map[int]float64, error) {
	var err error
	q := `
		SELECT
			year,
			sum(global_sale)
		FROM games
		WHERE  year>=$1 AND year<=$2
		GROUP BY year
	`
	rows, err := Query(ctx, q, start, end)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := make(map[int]float64)
	for rows.Next() {
		var year int
		var totalSell float64
		err = rows.Scan(&year, &totalSell)
		if err != nil {
			return nil, err
		}
		res[year] = totalSell
	}
	return res, err
}

func GetTotalSellByName(ctx context.Context, name string) (float64, float64, float64, float64, float64, error) {
	var err error
	q := `
		SELECT
			sum(na_sale),
			sum(eu_sale),
			sum(jp_sale),
			sum(other_sale),
			sum(global_sale)
		FROM games
		WHERE  name=$1
		GROUP BY name
	`
	rows, err := Query(ctx, q, name)
	if err != nil {
		return 0, 0, 0, 0, 0, err
	}
	defer rows.Close()
	var naSell, euSell, jpSell, otherSell, globalSell float64
	if rows.Next() {
		err = rows.Scan(&naSell, &euSell, &jpSell, &otherSell, &globalSell)
		if err != nil {
			return 0, 0, 0, 0, 0, err
		}
	}
	return naSell, euSell, jpSell, otherSell, globalSell, err
}
