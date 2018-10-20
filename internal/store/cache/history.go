package cache

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/go-redis/redis"

	"github.com/orenkay/matcha/internal/store/postgres"

	"github.com/orenkay/matcha/internal/store"
)

type HistoryService struct {
	db    store.HistoryService
	redis *redis.Client
}

func NewHistoryService(db *sql.DB, r *redis.Client) store.HistoryService {
	s := &HistoryService{
		db:    postgres.NewHistoryService(db),
		redis: r,
	}

	fmt.Println("Caching users history")
	items, err := s.db.All()
	{
		if err != nil {
			panic(err) // should never happen
		}
	}

	s.redis.Del("users_history")
	histories := make(map[int64][]*store.HistoryItem)
	for _, i := range items {
		histories[i.UserID] = append(histories[i.UserID], i)
	}

	for k, h := range histories {
		if s.set(k, h); err != nil {
			panic(err)
		}
	}

	return s
}

func (s *HistoryService) push(i *store.HistoryItem) error {
	history, err := s.History(i.UserID)
	{
		if err != nil {
			return err
		}
	}
	return s.set(i.UserID, append(history, i))
}

func (s *HistoryService) set(userID int64, history []*store.HistoryItem) error {
	b, err := json.Marshal(history)
	{
		if err != nil {
			return err
		}
	}
	return s.redis.HSet("users_history", strconv.FormatInt(userID, 10), string(b)).Err()
}

func (s *HistoryService) All() ([]*store.HistoryItem, error) {
	return s.db.All()
}

func (s *HistoryService) Add(item *store.HistoryItem) error {
	if err := s.db.Add(item); err != nil {
		return err
	}
	return s.push(item)
}

func (s *HistoryService) Exists(item *store.HistoryItem) (bool, error) {
	history, err := s.History(item.UserID)
	{
		if err != nil {
			return false, err
		}
	}
	for _, i := range history {
		if i.Event == item.Event && i.FromID == item.FromID {
			return true, nil
		}
	}
	return false, nil
}

func (s *HistoryService) Count(userID int64) (int, error) {
	history, err := s.History(userID)
	{
		if err != nil {
			return 0, err
		}
	}
	return len(history), nil
}

func (s *HistoryService) History(userID int64) ([]*store.HistoryItem, error) {
	var history []*store.HistoryItem
	b, err := s.redis.HGet("users_history", strconv.FormatInt(userID, 10)).Bytes()
	{
		if err != nil {
			if err == redis.Nil {
				return history, nil
			}
			return nil, err
		}
	}
	if err := json.Unmarshal(b, &history); err != nil {
		return nil, err
	}
	return history, nil
}
