package cache

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/go-redis/redis"
	"github.com/orenkay/matcha/internal/store"
	"github.com/orenkay/matcha/internal/store/postgres"
)

type ProfileService struct {
	db    store.ProfileService
	redis *redis.Client
}

func NewProfileService(db *sql.DB, redis *redis.Client) store.ProfileService {
	s := &ProfileService{
		db:    postgres.NewProfileService(db),
		redis: redis,
	}
	fmt.Println("Caching users profiles")
	profiles, err := s.db.Profiles()
	{
		if err != nil {
			panic(err) // should never happen
		}
	}

	s.redis.Del("profiles")
	for _, p := range profiles {
		if err := s.Set(p.UserID, p); err != nil {
			panic(err) // should never happen
		}
	}
	return s
}

func (s *ProfileService) Get(userID int64) (*store.Profile, error) {
	b, err := s.redis.HGet("profiles", strconv.FormatInt(userID, 10)).Bytes()
	{
		if err != nil {
			if err == redis.Nil {
				return nil, nil
			}
			return nil, err
		}
	}
	p := &store.Profile{}
	if err := json.Unmarshal(b, p); err != nil {
		return nil, err
	}
	return p, nil
}

func (s *ProfileService) Set(userID int64, profile *store.Profile) error {
	b, err := json.Marshal(profile)
	{
		if err != nil {
			return err
		}
	}
	return s.redis.HSet("profiles", strconv.FormatInt(userID, 10), string(b)).Err()
}

func (s *ProfileService) Del(userID int64) {
	s.redis.HDel("profiles", strconv.FormatInt(userID, 10))
}

func (s *ProfileService) Profiles() ([]*store.Profile, error) {
	var profiles []*store.Profile

	m, err := s.redis.HGetAll("profiles").Result()
	{
		if err != nil && err != redis.Nil {
			return nil, err
		}
	}

	// We cache here
	if err == redis.Nil {
		profiles, err := s.db.Profiles()
		{
			if err != nil {
				return nil, err
			}
		}
		for _, p := range profiles {
			go s.Set(p.UserID, p)
		}
	}

	for _, s := range m {
		p := &store.Profile{}
		if err := json.Unmarshal([]byte(s), p); err != nil {
			return nil, err
		}
		profiles = append(profiles, p)
	}

	return profiles, nil
}

func (s *ProfileService) Profile(userID int64) (*store.Profile, error) {
	return s.Get(userID)
}

func (s *ProfileService) Add(profile *store.Profile) error {
	if err := s.db.Add(profile); err != nil {
		return err
	}
	return s.Set(profile.UserID, profile)
}

func (s *ProfileService) Delete(userID int64) error {
	if err := s.db.Delete(userID); err != nil {
		return err
	}
	s.Del(userID)
	return nil
}

func (s *ProfileService) Update(profile *store.Profile) error {
	if err := s.db.Update(profile); err != nil {
		return err
	}
	return s.Set(profile.UserID, profile)
}
