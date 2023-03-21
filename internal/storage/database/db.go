package db

import (
	"bot/internal/config"
	"bot/internal/lib/e"
	"bot/internal/storage"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"io/ioutil"
	"log"
	"net/http"
)

type Storage struct {
	Config config.Config
	Db     *sqlx.DB
}

type Movie struct {
	KinopoiskId  int    `json:"kinopoiskId"`
	ImdbId       string `json:"imdbId"`
	NameRu       string `json:"nameRu"`
	NameEn       string `json:"nameEn"`
	NameOriginal string `json:"nameOriginal"`
	Countries    []struct {
		Country string `json:"country"`
	} `json:"countries"`
	Genres []struct {
		Genre string `json:"genre"`
	} `json:"genres"`
	RatingKinopoisk  float64 `json:"ratingKinopoisk"`
	RatingImdb       float64 `json:"ratingImdb"`
	Year             int     `json:"year"`
	Type             string  `json:"type"`
	PosterUrl        string  `json:"posterUrl"`
	PosterUrlPreview string  `json:"posterUrlPreview"`
}

type MoviesResponse struct {
	Total      int     `json:"total"`
	TotalPages int     `json:"totalPages"`
	Items      []Movie `json:"items"`
}

func New(ctx context.Context, config config.Config) Storage {
	db, err := dbConnect(ctx, config.Post)
	if err != nil {
		log.Fatal(e.Wrap("Failed to connect to database", err))
	}
	return Storage{
		Config: config,
		Db:     db,
	}

}

func dbConnect(ctx context.Context, cfg config.PostgresConfig) (*sqlx.DB, error) {
	fmt.Println(cfg)
	db, err := sqlx.ConnectContext(ctx, "postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.DBName, cfg.Password, cfg.SSLMode))
	if err != nil {
		return nil, err
	}

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func (s Storage) Search(ctx context.Context, genre, id, page int) (*storage.Film, error) {
	url := fmt.Sprintf("https://kinopoiskapiunofficial.tech/api/v2.2/films?genres=%d&page=%d", genre, page)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, e.Wrap("can't create a new request", err)
	}

	req.Header.Set("X-API-KEY", s.Config.ApiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, e.Wrap("request submission error", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, e.Wrap("incorrect response status", err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, e.Wrap("read request body error", err)
	}

	var searchRes MoviesResponse
	if err := json.Unmarshal(body, &searchRes); err != nil {
		return nil, e.Wrap("translation error in json", err)
	}

	fmt.Println("afhu- ", len(searchRes.Items))
	for _, film := range searchRes.Items {
		if ok, err := s.IsExists(ctx, film.KinopoiskId, id, genre); ok == false && err == nil {
			if film.NameRu == "" {
				film.NameRu = film.NameEn
				if film.NameRu == "" {
					film.NameRu = film.NameOriginal
				}
			}
			return &storage.Film{
				Id:   film.KinopoiskId,
				Rate: film.RatingKinopoisk,
				Name: film.NameRu,
				Year: film.Year,
				Url:  fmt.Sprintf("https://www.kinopoisk.ru/film/%d/", film.KinopoiskId),
			}, nil
		} else if err != nil {
			return nil, err
		}
	}
	return nil, nil
}

func (s Storage) Save(ctx context.Context, film *storage.Film, userID int, genre string) error {
	q := `INSERT INTO movies (user_id, genre, movie_id, movie_name, time) VALUES ($1, $2, $3, $4, CURRENT_TIMESTAMP)`
	_, err := s.Db.ExecContext(ctx, q, userID, genre, film.Id, film.Name)
	if err != nil {
		return e.Wrap("can't save response", err)
	}
	return nil
}

func (s Storage) IsExists(ctx context.Context, filmId, userID, genre int) (bool, error) {
	var count int
	q := `SELECT COUNT(*) FROM movies WHERE user_id = $1 AND movie_id = $2`
	err := s.Db.QueryRowContext(ctx, q, userID, filmId).Scan(&count)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, e.Wrap("failed to get count", err)
	}

	return count > 0, nil
}

func (s Storage) Stat(ctx context.Context, userID int) (*storage.History, error) {
	hist := storage.History{}
	q := `SELECT genre, movie_name, time FROM movies WHERE user_id = $1 ORDER BY time ASC LIMIT 1`
	row := s.Db.QueryRowContext(ctx, q, userID)
	err := row.Scan(&hist.Genre, &hist.Movie, &hist.Time)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, e.Wrap("failed to get count", err)
	}
	v, err := s.groupStat(ctx)
	if err != nil {
		return nil, err
	}
	hist.Stat = v
	return &hist, nil
}

func (s Storage) groupStat(ctx context.Context) ([]storage.Stat, error) {
	listStat := make([]storage.Stat, 0)
	q := `SELECT genre, COUNT(*) as count FROM movies GROUP BY genre`
	row, err := s.Db.QueryContext(ctx, q)
	if err != nil {
		if err == sql.ErrNoRows {
			// обработка ситуации, когда запрос не найден
			return nil, nil
		}
		return nil, e.Wrap("database connection error", err)
	}
	defer row.Close()
	for row.Next() {
		tmp := storage.Stat{}
		err = row.Scan(&tmp.Genre, &tmp.Count)
		if err != nil {
			return nil, e.Wrap("database connection error", err)
		}
		listStat = append(listStat, tmp)
	}
	return listStat, nil
}
