package telegram

import (
	"bot/internal/storage"
	"context"
	"log"
	"strings"
)

// RndCmd   = "/rnd"
const (
	HelpCmd  = "/help"
	StartCmd = "/start"
	Stat     = "/stat"
)

func (p *Processor) doCmd(ctx context.Context, text string, meta Meta, id int) error {
	text = strings.TrimSpace(text)

	log.Printf("got new command '%s' from '%s", text, meta.UserName)

	switch text {
	case HelpCmd:
		return p.sendHelp(ctx, meta.ChatID)
	case StartCmd:
		return p.sendHello(ctx, meta.ChatID)
	case Stat:
		return p.sendStat(ctx, id, meta.ChatID)
	default:
		return p.sendRandomFilm(ctx, meta.ChatID, text, id)
	}
}

func (p *Processor) sendStat(ctx context.Context, id, chatId int) error {
	v, err := p.storage.Stat(ctx, id)
	if err == nil && v == nil {
		p.cli.SendMassage(ctx, chatId, msgNoHistory)
	} else if err == nil && v != nil {
		p.cli.SendMassage(ctx, chatId, v.ToString())
	}
	return err
}

func (p *Processor) sendRandomFilm(ctx context.Context, chatID int, genre string, id int) error {

	genres := map[string]int{
		"Аниме":           1750,
		"Короткометражка": 15,
		"Биография":       22,
		"Боевик":          3,
		"Вестерн":         13,
		"Военный":         19,
		"Детектив":        17,
		"Детский":         456,
		"Для взрослых":    20,
		"Документальный":  12,
		"Драма":           8,
		"Игра":            27,
		"История":         23,
		"Комедия":         6,
		"Концерт":         1747,
		"Криминал":        16,
		"Мелодрама":       7,
		"Музыка":          21,
		"Мюзикл":          9,
		"Новости":         28,
		"Приключения":     10,
		"Реальное ТВ":     25,
		"Семейный":        11,
		"Спорт":           24,
		"Триллер":         4,
		"Ужасы":           1,
		"Фантастика":      2,
		"Фэнтези":         5,
		"Церемония":       1751,
	}

	if value, ok := genres[genre]; ok {
		var inf *storage.Film
		page := 1
		for ; page < 4; page++ {
			tmpInf, err := p.storage.Search(ctx, value, id, page)
			if err != nil {
				return err
			}
			if tmpInf != nil {
				inf = tmpInf
				break
			}
		}
		if inf == nil {
			return p.cli.SendMassage(ctx, chatID, msgAllMovie)
		}
		if err := p.storage.Save(ctx, inf, id, genre); err != nil {
			return err
		}
		return p.cli.SendMassage(ctx, chatID, inf.ToString())
	}
	return p.cli.SendMassage(ctx, chatID, msgUnknownCommand)
}

func (p *Processor) sendHelp(ctx context.Context, chatID int) error {
	return p.cli.SendMassage(ctx, chatID, msgHelp)
}

func (p *Processor) sendHello(ctx context.Context, chatID int) error {
	return p.cli.SendMassage(ctx, chatID, msgHello)
}
