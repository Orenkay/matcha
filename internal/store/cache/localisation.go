package cache

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/go-redis/redis"

	"github.com/orenkay/matcha/internal/localisation"
	"github.com/orenkay/matcha/internal/store/postgres"

	"github.com/orenkay/matcha/internal/store"
)

type LocalisationService struct {
	redis *redis.Client
	db    store.LocalisationService
}

func NewLocalisationService(db *sql.DB, r *redis.Client) store.LocalisationService {
	s := &LocalisationService{
		db:    postgres.NewLocalisationService(db),
		redis: r,
	}

	fmt.Println("Caching users localisations")
	locs, err := s.db.Localisations()
	{
		if err != nil {
			panic(err) // should never happen
		}
	}

	s.redis.Del("localisations")
	for _, e := range locs {
		if err := s.Set(e); err != nil {
			panic(err) // should never happen
		}
	}
	return s
}

func (s *LocalisationService) Set(loc *store.Localisation) error {
	b, err := json.Marshal(loc)
	{
		if err != nil {
			return err
		}
	}
	return s.redis.HSet("localisations", strconv.FormatInt(loc.UserID, 10), string(b)).Err()
}

func (s *LocalisationService) Add(loc *store.Localisation) error {
	if err := s.db.Add(loc); err != nil {
		return err
	}
	return s.Set(loc)
}

func (s *LocalisationService) Localisations() ([]*store.Localisation, error) {
	var localisations []*store.Localisation

	m, err := s.redis.HGetAll("localisations").Result()
	{
		if err != nil {
			return nil, err
		}
	}
	for _, s := range m {
		loc := &store.Localisation{}
		if err := json.Unmarshal([]byte(s), loc); err != nil {
			return nil, err
		}
		localisations = append(localisations, loc)
	}
	return localisations, nil
}

func (s *LocalisationService) Update(loc *store.Localisation) error {
	if err := s.db.Update(loc); err != nil {
		return err
	}
	s.Set(loc)
	return nil
}

func (s *LocalisationService) DistanceBetween(userID, targetID int64) (float64, error) {
	loc1, err := s.Localisation(userID)
	{
		if err != nil {
			return 0, err
		}
	}
	loc2, err := s.Localisation(targetID)
	{
		if err != nil {
			return 0, err
		}
	}
	if loc1 == nil || loc2 == nil {
		return -1, nil
	}
	return localisation.CoordDistance(loc1, loc2), nil
}

func (s *LocalisationService) Localisation(userID int64) (*store.Localisation, error) {
	b, err := s.redis.HGet("localisations", strconv.FormatInt(userID, 10)).Bytes()
	{
		if err != nil {
			if err == redis.Nil {
				return nil, nil
			}
			return nil, err
		}
	}
	loc := &store.Localisation{}
	if err := json.Unmarshal(b, &loc); err != nil {
		return nil, err
	}
	return loc, nil
}
