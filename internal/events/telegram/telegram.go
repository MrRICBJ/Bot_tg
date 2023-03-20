package telegram

import (
	telegram2 "bot/internal/clients/telegram"
	"bot/internal/events"
	"bot/internal/lib/e"
	"bot/internal/storage"
	"context"
	"errors"
)

type Processor struct {
	cli     *telegram2.Client
	offset  int
	storage storage.Storage
}

type Meta struct {
	ChatID   int
	UserName string
}

var (
	ErrUnknownEventType = errors.New("unknown event type")
	ErrUnknownMetaType  = errors.New("unknown meta type")
)

func New(client *telegram2.Client, storage storage.Storage) *Processor {
	return &Processor{
		cli:     client,
		storage: storage,
	}
}

func (p *Processor) Fetch(ctx context.Context, limit int) ([]events.Event, error) {
	updates, err := p.cli.Updates(ctx, p.offset, limit)
	if err != nil {
		return nil, e.Wrap("can't get events", err)
	}

	if len(updates) == 0 {
		return nil, nil
	}

	res := make([]events.Event, 0, len(updates))

	for _, u := range updates {
		res = append(res, event(u))
	}

	p.offset = updates[len(updates)-1].ID + 1

	return res, nil
}

func (p *Processor) Process(ctx context.Context, event events.Event) error {
	switch event.Type {
	case events.Message:
		return p.processMessage(ctx, event)
	default:
		return e.Wrap("can't process message", ErrUnknownEventType)
	}
}

func (p *Processor) processMessage(ctx context.Context, event events.Event) error {
	meta, err := meta(event)
	if err != nil {
		return e.Wrap("can't process message", err)
	}

	if err := p.doCmd(ctx, event.Text, meta, event.UserId); err != nil {
		return e.Wrap("can't process message", err)
	}

	return nil
}

func meta(event events.Event) (Meta, error) {
	res, ok := event.Meta.(Meta)
	if !ok {
		return Meta{}, e.Wrap("can't get meta", ErrUnknownMetaType)
	}
	return res, nil
}

func event(u telegram2.Update) events.Event {
	updType := fetchType(u)
	res := events.Event{
		Type:   updType,
		Text:   fetchText(u),
		UserId: u.Message.From.UserId,
	}

	if updType == events.Message {
		res.Meta = Meta{
			ChatID:   u.Message.Chat.ID,
			UserName: u.Message.From.Username,
		}
	}

	return res
}

func fetchType(u telegram2.Update) events.Type {
	if u.Message == nil {
		return events.Unknown
	}
	return events.Message
}

func fetchText(u telegram2.Update) string {
	if u.Message == nil {
		return ""
	}
	return u.Message.Text
}
