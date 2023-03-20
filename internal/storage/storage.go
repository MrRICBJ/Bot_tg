package storage

import (
	"fmt"
	"golang.org/x/net/context"
	"time"
)

type Storage interface {
	IsExists(ctx context.Context, filmId, userID, genre int) (bool, error)
	Search(ctx context.Context, genre, id, page int) (*Film, error)
	Save(ctx context.Context, film *Film, userID int, genre string) error
	Stat(ctx context.Context, userID int) (*History, error)
}

type Film struct {
	Year int     `bson:"year"`
	Rate float64 `bson:"rate"`
	Name string  `bson:"name"`
	Url  string  `bson:"url"`
	Id   int     `bson:"id"`
}

type History struct {
	Genre string
	Movie string
	Time  time.Time
	Stat  []Stat
}

type Stat struct {
	Genre string
	Count int
}

func (h *History) ToString() string {
	head := fmt.Sprintf("Your first request🥳:\n 👉 %s at %d:%d:%d %d.%d.%d\n --->\"%s\"<--- 👈\n",
		h.Genre, h.Time.Hour(), h.Time.Minute(), h.Time.Second(), h.Time.Day(), h.Time.Month(), h.Time.Year(), h.Movie)
	stat := "\nQuery Statistics:\n"
	for i := 0; i < len(h.Stat); i++ {
		stat = fmt.Sprintf("%s %s - %d\n", stat, h.Stat[i].Genre, h.Stat[i].Count)
	}

	return head + stat
}

func (p *Film) ToString() string {
	return fmt.Sprintf("Название: %s\nГод: %d\nРейтинг: %.1f\nСсылка: %s\n", p.Name, p.Year, p.Rate, p.Url)
}
