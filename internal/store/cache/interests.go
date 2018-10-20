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

type InterestService struct {
	db    store.InterestService
	redis *redis.Client
}

func NewInterestService(db *sql.DB, r *redis.Client) store.InterestService {
	s := &InterestService{
		db:    postgres.NewInterestService(db),
		redis: r,
	}

	fmt.Println("Caching users interests")
	interests, err := s.db.Interests()
	{
		if err != nil {
			panic(err) // Should not happen
		}
	}

	s.redis.Del("users_interests")

	interestsMap := make(map[int64][]*store.Interest)
	for _, i := range interests {
		if err := s.push(i); err != nil {
			panic(err) // should not happen
		}
	}

	for _, i := range interests {
		interestsMap[i.UserID] = append(interestsMap[i.UserID], i)
	}

	for k, v := range interestsMap {
		if s.set(k, v); err != nil {
			panic(err)
		}
	}

	return s
}

func (s *InterestService) get(userID int64) ([]*store.Interest, error) {
	var interests []*store.Interest
	b, err := s.redis.HGet("users_interests", strconv.FormatInt(userID, 10)).Bytes()
	{
		if err != nil {
			if err == redis.Nil {
				return interests, nil
			}
			return nil, err
		}
	}
	if err := json.Unmarshal(b, &interests); err != nil {
		return nil, err
	}
	return interests, nil
}

func (s *InterestService) set(userID int64, interests []*store.Interest) error {
	b, err := json.Marshal(interests)
	{
		if err != nil {
			return err
		}
	}
	return s.redis.HSet("users_interests", strconv.FormatInt(userID, 10), string(b)).Err()
}

func (s *InterestService) push(i *store.Interest) error {
	interests, err := s.get(i.UserID)
	{
		if err != nil {
			return err
		}
	}
	return s.set(i.UserID, append(interests, i))
}

func (s *InterestService) Add(interest *store.Interest) error {
	if err := s.db.Add(interest); err != nil {
		return err
	}
	return s.push(interest)
}

func (s *InterestService) Remove(userID int64, slug string) error {
	if err := s.db.Remove(userID, slug); err != nil {
		return err
	}
	interests, err := s.get(userID)
	{
		if err != nil {
			return err
		}
	}
	index := -1
	{
		for i, interest := range interests {
			if interest.Value == slug {
				index = i
				break
			}
		}
		if index == -1 {
			return nil
		}
	}

	interests[index] = interests[len(interests)-1]
	interests[len(interests)-1] = nil
	interests = interests[:len(interests)-1]
	return s.set(userID, interests)
}

func (s *InterestService) Interests() ([]*store.Interest, error) {
	return s.db.Interests()
}

func (s *InterestService) Interest(userID int64, slug string) (*store.Interest, error) {
	interests, err := s.get(userID)
	{
		if err != nil {
			return nil, err
		}
	}
	for _, i := range interests {
		if i.Value == slug {
			return i, nil
		}
	}
	return nil, nil
}

func (s *InterestService) AllByUser(userID int64) ([]*store.Interest, error) {
	return s.get(userID)
}

func (s *InterestService) AllBySlug(slug string, limit int) ([]*store.Interest, error) {
	return s.db.AllBySlug(slug, limit)
}
