package cache

import (
	"encoding/json"
	"strconv"

	"github.com/go-redis/redis"

	"github.com/orenkay/matcha/internal/store"
)

type NotificationService struct {
	redis *redis.Client
	is    store.IgnoreService
}

func NewNotificationService(r *redis.Client, is store.IgnoreService) *NotificationService {
	return &NotificationService{
		redis: r,
		is:    is,
	}
}

func (s *NotificationService) Set(userID int64, notifs []*store.Notification) error {
	b, err := json.Marshal(notifs)
	{
		if err != nil {
			return err
		}
	}
	return s.redis.HSet("notifications", strconv.FormatInt(userID, 10), string(b)).Err()
}

func (s *NotificationService) Get(userID int64) ([]*store.Notification, error) {
	var notifications []*store.Notification
	b, err := s.redis.HGet("notifications", strconv.FormatInt(userID, 10)).Bytes()
	{
		if err != nil {
			if err == redis.Nil {
				return notifications, nil
			}
			return nil, err
		}
	}
	if err := json.Unmarshal(b, &notifications); err != nil {
		return nil, err
	}
	return notifications, nil
}

func (s *NotificationService) Push(to int64, from int64, evt string, data interface{}) error {
	ignored, err := s.is.Ignored(to, from)
	{
		if err != nil || ignored {
			return err
		}
	}

	notifications, err := s.Get(to)
	{
		if err != nil {
			return err
		}
	}
	return s.Set(to, append(notifications, &store.Notification{
		Event: evt,
		Data:  data,
	}))
}

func (s *NotificationService) Notifications(userID int64) []*store.Notification {
	notifs, err := s.Get(userID)
	{
		if err == nil {
			s.redis.HDel("notifications", strconv.FormatInt(userID, 10))
		}
	}
	return notifs
}
