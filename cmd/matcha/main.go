package main

import (
	"database/sql"

	"github.com/go-redis/redis"

	"github.com/orenkay/matcha/internal/server"
	"github.com/orenkay/matcha/internal/store"
	"github.com/orenkay/matcha/internal/store/cache"
	"github.com/orenkay/matcha/internal/store/postgres"

	_ "github.com/lib/pq"
)

func main() {
	db, err := sql.Open("postgres", "host=postgres user=matcha password=matcha dbname=matcha sslmode=disable")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	r := redis.NewClient(&redis.Options{
		Addr: "redis:6379",
	})
	if err := r.Ping().Err(); err != nil {
		panic(err)
	}
	defer r.Close()

	ignoreService := postgres.NewIgnoreService(db)
	historyService := cache.NewHistoryService(db, r)
	notificationService := cache.NewNotificationService(r, ignoreService)

	s := &store.Store{
		UserService:         postgres.NewUserService(db),
		ValidationService:   postgres.NewValidationService(db),
		PresenceService:     postgres.NewPresenceService(db),
		ReportService:       postgres.NewReportService(db),
		MatchService:        postgres.NewMatchService(db, notificationService),
		MessageService:      postgres.NewMessageService(db, notificationService),
		PicturesService:     postgres.NewPicturesService(db),
		ProfileService:      cache.NewProfileService(db, r),
		LikesService:        cache.NewLikesService(db, r, notificationService),
		InterestService:     cache.NewInterestService(db, r),
		LocalisationService: cache.NewLocalisationService(db, r),
		AuthTokenService:    cache.NewAuthTokenService(r),
		IgnoreService:       ignoreService,
		HistoryService:      historyService,
		NotificationService: notificationService,
	}

	server := server.New(s)
	server.Run("0.0.0.0:3000")

}
