package main

import (
	"database/sql"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/icrowley/fake"
	"github.com/orenkay/matcha/internal/store"

	_ "github.com/lib/pq"
	"github.com/orenkay/matcha/internal/store/postgres"
)

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

type FakeNotificationsService struct{}

func (s *FakeNotificationsService) Push(to int64, from int64, evt string, data interface{}) error {
	return nil
}

func (s *FakeNotificationsService) Notifications(userID int64) []*store.Notification {
	return nil
}

func MaxTimestamp() int {
	curr := time.Now()
	max := time.Date(curr.Year()-18, curr.Month(), curr.Day(), curr.Hour(), curr.Minute(), curr.Second(), curr.Nanosecond(), curr.Location())
	return int(max.Unix())
}

func Seed(s *store.Store, interestsTags []string) {
	user := &store.User{
		Username: "fake_account",
		Email:    "fake_account",
		Password: "fake_account",
	}
	checkError(s.UserService.Add(user))

	profile := &store.Profile{
		UserID:     user.ID,
		Attraction: []string{"hetero", "bi", "homo"}[rand.Intn(3)],
		Gender:     []string{"male", "female"}[rand.Intn(2)],
		Birthdate:  int64(rand.Intn(MaxTimestamp())),
		LastName:   fake.LastName(),
		Bio:        fake.SentencesN(2),
	}
	if profile.Gender == "male" {
		profile.FirstName = fake.MaleFirstName()
	} else {
		profile.FirstName = fake.FemaleFirstName()
	}
	checkError(s.ProfileService.Add(profile))

	picture := &store.Picture{
		UserID: user.ID,
		IsPP:   true,
	}

	if profile.Gender == "male" {
		picture.Path = fmt.Sprintf("https://randomuser.me/api/portraits/men/%d.jpg", rand.Intn(100))
	} else {
		picture.Path = fmt.Sprintf("https://randomuser.me/api/portraits/women/%d.jpg", rand.Intn(100))
	}

	checkError(s.PicturesService.Add(picture))

	for i := 0; i < rand.Intn(20); i++ {
		interest := &store.Interest{
			UserID: user.ID,
			Value:  interestsTags[rand.Intn(len(interestsTags))],
		}
		checkError(s.InterestService.Add(interest))
	}

	loc := &store.Localisation{
		UserID:  user.ID,
		Lat:     42 + rand.Float64()*(52-42),
		Lng:     -4 + rand.Float64()*(8-(-4)),
		Address: "FAKE_ADDRESS",
	}
	checkError(s.LocalisationService.Add(loc))
	checkError(s.PresenceService.Add(user.ID))

	visites := rand.Intn(100)
	for i := 0; i < visites; i++ {
		checkError(
			s.HistoryService.Add(
				&store.HistoryItem{
					UserID: user.ID,
					FromID: int64(i),
					Event:  "visite",
				}),
		)
	}
	for i := 0; i < rand.Intn(visites+1); i++ {
		checkError(s.LikesService.Add(int64(i), user.ID))
	}
}

func main() {
	db, err := sql.Open("postgres", "host=postgres user=matcha password=matcha dbname=matcha sslmode=disable")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	count := 100

	if len(os.Args) == 2 {
		c, err := strconv.Atoi(os.Args[1])
		{
			if err == nil {
				count = c
			}
		}
	}

	s := &store.Store{
		UserService:         postgres.NewUserService(db),
		ProfileService:      postgres.NewProfileService(db),
		PicturesService:     postgres.NewPicturesService(db),
		InterestService:     postgres.NewInterestService(db),
		LocalisationService: postgres.NewLocalisationService(db),
		PresenceService:     postgres.NewPresenceService(db),
		HistoryService:      postgres.NewHistoryService(db),
		LikesService:        postgres.NewLikesService(db, &FakeNotificationsService{}),
	}

	var wg sync.WaitGroup
	var interestsTags []string
	var dones uint64
	semaphore := make(chan struct{}, 90)

	for i := 0; i < 50; i++ {
		interestsTags = append(interestsTags, fake.Word())
	}

	for i := 0; i < count; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			semaphore <- struct{}{}
			defer func() {
				<-semaphore
			}()

			atomic.AddUint64(&dones, 1)
			fmt.Printf("\r\033[1K%d / %d", dones, count)
			Seed(s, interestsTags)
		}()
	}

	wg.Wait()
	fmt.Println("\nDatabase seeded !")
}
