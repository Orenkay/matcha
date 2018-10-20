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

type LikesService struct {
	db    store.LikesService
	redis *redis.Client
}

func NewLikesService(db *sql.DB, r *redis.Client, ns store.NotificationService) store.LikesService {
	s := &LikesService{
		db:    postgres.NewLikesService(db, ns),
		redis: r,
	}

	fmt.Println("Caching users likes")
	likes, err := s.db.Likes()
	{
		if err != nil {
			panic(err)
		}
	}

	s.redis.Del("users_likes")
	s.redis.Del("users_likes_count")

	likesMap := make(map[int64][]*store.Like)
	for _, l := range likes {
		likesMap[l.UserID] = append(likesMap[l.UserID], l)
		if s.redis.HIncrBy("users_likes_count", strconv.FormatInt(l.TargetID, 10), 1).Err(); err != nil {
			panic(err)
		}
	}
	for k, l := range likesMap {
		if s.set(k, l); err != nil {
			panic(err)
		}
	}
	return s
}

func (s *LikesService) get(userID int64) ([]*store.Like, error) {
	var likes []*store.Like
	b, err := s.redis.HGet("users_likes", strconv.FormatInt(userID, 10)).Bytes()
	{
		if err != nil {
			if err == redis.Nil {
				return likes, nil
			}
			return nil, err
		}
	}
	if err := json.Unmarshal(b, &likes); err != nil {
		return nil, err
	}
	return likes, nil
}

func (s *LikesService) set(userID int64, likes []*store.Like) error {
	b, err := json.Marshal(likes)
	{
		if err != nil {
			return err
		}
	}
	return s.redis.HSet("users_likes", strconv.FormatInt(userID, 10), string(b)).Err()
}

func (s *LikesService) push(userID int64, l *store.Like) error {
	likes, err := s.get(userID)
	{
		if err != nil {
			return err
		}
	}
	err = s.redis.HIncrBy("users_likes_count", strconv.FormatInt(l.TargetID, 10), 1).Err()
	{
		if err != nil {
			return err
		}
	}
	return s.set(userID, append(likes, l))
}

func (s *LikesService) remove(userID int64, targetID int64) error {
	likes, err := s.get(userID)
	{
		if err != nil {
			return err
		}
	}
	index := -1
	{
		for i, l := range likes {
			if l.TargetID == targetID {
				index = i
				break
			}
		}
		if index == -1 {
			return nil
		}
	}

	err = s.redis.HIncrBy("users_likes_count", strconv.FormatInt(targetID, 10), -1).Err()
	{
		if err != nil {
			return err
		}
	}

	likes[index] = likes[len(likes)-1]
	likes[len(likes)-1] = nil
	likes = likes[:len(likes)-1]
	return s.set(userID, likes)
}

func (s *LikesService) Likes() ([]*store.Like, error) {
	return s.db.Likes()
}

func (s *LikesService) Like(userID, targetID int64) (bool, error) {
	return s.db.Like(userID, targetID)
}

func (s *LikesService) Add(userID, id int64) error {
	if err := s.db.Add(userID, id); err != nil {
		return err
	}
	return s.push(userID, &store.Like{UserID: userID, TargetID: id})
}

func (s *LikesService) Count(userID int64) (int, error) {
	count, err := s.redis.HGet("users_likes_count", strconv.FormatInt(userID, 10)).Int()
	{
		if err != nil && err == redis.Nil {
			return 0, nil
		}
	}
	return count, err
}

func (s *LikesService) Remove(userID, id int64) error {
	if err := s.db.Remove(userID, id); err != nil {
		return err
	}
	return s.remove(userID, id)
}
