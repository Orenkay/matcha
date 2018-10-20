package matcher

import (
	"encoding/json"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"

	"github.com/go-chi/render"
	"github.com/orenkay/matcha/internal/api"
	"github.com/orenkay/matcha/internal/store"
)

type MatchedProfile struct {
	UserID     int64
	Interests  []*store.Interest
	Birthdate  int64
	Popularity float64
	Distance   float64
	Visited    bool
	Loc        *store.Localisation
}

func LookingFor(p *store.Profile) []string {
	if (p.Gender == "male" && p.Attraction == "hetero") || (p.Gender == "female" && p.Attraction == "homo") {
		return []string{"female"}
	}
	if (p.Gender == "male" && p.Attraction == "homo") || (p.Gender == "female" && p.Attraction == "hetero") {
		return []string{"male"}
	}
	return []string{"female", "male"}
}

func IsMatchable(p1 *store.Profile, p2 *store.Profile) bool {
	var cond1, cond2 bool

	for _, g := range LookingFor(p1) {
		if p2.Gender == g {
			cond1 = true
			break
		}
	}
	for _, g := range LookingFor(p2) {
		if p1.Gender == g {
			cond2 = true
			break
		}
	}

	return cond1 && cond2
}

type FilterHelper func(*MatchedProfile) bool

func MatchableProfiles(s *store.Store, up *store.Profile, h FilterHelper) ([]*MatchedProfile, error) {
	profiles, err := s.ProfileService.Profiles()
	{
		if err != nil {
			return nil, err
		}
	}
	ret := make([]*MatchedProfile, 0, len(profiles))
	for _, p := range profiles {
		if p.UserID == up.UserID {
			continue
		}

		if !IsMatchable(up, p) {
			continue
		}

		interests, err := s.InterestService.AllByUser(p.UserID)
		{
			if err != nil || len(interests) == 0 {
				continue
			}
		}

		loc, err := s.LocalisationService.Localisation(p.UserID)
		{
			if err != nil || loc == nil {
				continue
			}
		}

		pp, err := s.PicturesService.PP(p.UserID)
		{
			if err != nil || pp == nil {
				continue
			}
		}

		distance, err := s.LocalisationService.DistanceBetween(up.UserID, p.UserID)
		{
			if err != nil || distance == -1 {
				continue
			}
		}

		visited, err := s.HistoryService.Exists(&store.HistoryItem{
			FromID: up.UserID,
			UserID: p.UserID,
			Event:  "visite",
		})
		{
			if err != nil {
				continue
			}
		}

		visites, err := s.HistoryService.Count(p.UserID)
		{
			if err != nil {
				continue
			}
		}
		likes, err := s.LikesService.Count(p.UserID)
		{
			if err != nil {
				continue
			}
		}
		var popularity float64
		if visites > 0 {
			popularity = float64(likes) / float64(visites)
		}

		mp := &MatchedProfile{
			UserID:     p.UserID,
			Interests:  interests,
			Birthdate:  p.Birthdate,
			Distance:   distance,
			Visited:    visited,
			Popularity: popularity,
			Loc:        loc,
		}

		if h(mp) {
			ret = append(ret, mp)
		}
	}
	return ret, nil
}

type MergeHelper func(*MatchedProfile, *MatchedProfile) bool

func MatchedMerge(h MergeHelper, a []*MatchedProfile, b []*MatchedProfile) []*MatchedProfile {
	s := make([]*MatchedProfile, 0, len(a)+len(b))
	for len(a) > 0 || len(b) > 0 {
		if len(a) == 0 {
			return append(s, b...)
		}
		if len(b) == 0 {
			return append(s, a...)
		}
		if h(a[0], b[0]) {
			s = append(s, a[0])
			a = a[1:]
		} else {
			s = append(s, b[0])
			b = b[1:]
		}
	}
	return s
}

func MatchedMergeSort(s []*MatchedProfile, h MergeHelper) []*MatchedProfile {
	if len(s) <= 1 {
		return s
	}
	n := len(s) / 2
	return MatchedMerge(h, MatchedMergeSort(s[:n], h), MatchedMergeSort(s[n:], h))
}

func interestsScore(interests1 []*store.Interest, interests2 []*store.Interest) float64 {

	dict := make(map[string]interface{})
	var n float64

	for _, e := range interests1 {
		dict[e.Value] = nil
	}

	for _, e := range interests2 {
		if _, ok := dict[e.Value]; ok {
			n++
		}
	}

	return n / float64(len(interests1))
}

type MatcherRequest struct {
	Filters struct {
		Age struct {
			Max int `json:"max"`
			Min int `json:"min"`
		} `json:"age"`

		Distance struct {
			Max float64 `json:"max"`
			Min float64 `json:"min"`
		} `json:"distance"`

		Popularity struct {
			Max float64 `json:"max"`
			Min float64 `json:"min"`
		} `json:"popularity"`
	}
	Sort struct {
		By   string
		Desc bool
	} `json:"sort"`
	IsSuggestion bool
}

func (data *MatcherRequest) Bind(r *http.Request) error {

	if filtersJSON := r.URL.Query().Get("filters"); filtersJSON != "" {
		if err := json.Unmarshal([]byte(filtersJSON), &data.Filters); err != nil {
			return err
		}
		data.Filters.Distance.Max *= 1000
		data.Filters.Distance.Min *= 1000
	}

	if sortJSON := r.URL.Query().Get("sort"); sortJSON != "" {
		if err := json.Unmarshal([]byte(sortJSON), &data.Sort); err != nil {
			return err
		}
	}

	var err error
	data.IsSuggestion, err = strconv.ParseBool(r.URL.Query().Get("suggestion"))
	{
		if err != nil {
			return err
		}
	}
	return nil
}

func Matcher(s *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value("user").(*store.User)

		offset, err := strconv.Atoi(chi.URLParam(r, "offset"))
		{
			if err != nil {
				render.Render(w, r, api.ErrInternal(err))
				return
			}
		}

		data := &MatcherRequest{}
		if err := data.Bind(r); err != nil {
			render.Render(w, r, api.ErrInvalidRequest(err))
			return
		}

		up, err := s.ProfileService.Profile(user.ID)
		{
			if err != nil {
				render.Render(w, r, api.ErrInternal(err))
				return
			}
		}
		interests, err := s.InterestService.AllByUser(user.ID)
		{
			if err != nil {
				render.Render(w, r, api.ErrInternal(err))
				return
			}
		}

		matches, err := MatchableProfiles(s, up, func(mp *MatchedProfile) bool {
			// age filters

			age := int(time.Since(time.Unix(mp.Birthdate, 0)).Hours() / 24 / 365)
			if data.Filters.Age.Max > 0 && age > data.Filters.Age.Max {
				return false
			}

			if data.Filters.Age.Min > 0 && age < data.Filters.Age.Min {
				return false
			}

			// distance filters
			if data.Filters.Distance.Max > 0 && data.Filters.Distance.Max < mp.Distance {
				return false
			}

			if data.Filters.Distance.Min > 0 && data.Filters.Distance.Min > mp.Distance {
				return false
			}

			// If we are in suggestions mode we filters already visited profiles
			if data.IsSuggestion && mp.Visited == true {
				return false
			}

			// populairty filters
			if data.Filters.Popularity.Max > 0 && data.Filters.Popularity.Max < mp.Popularity {
				return false
			}

			if data.Filters.Popularity.Min > 0 && data.Filters.Popularity.Min > mp.Popularity {
				return false
			}

			return true
		})
		{
			if err != nil {
				render.Render(w, r, api.ErrInternal(err))
				return
			}
		}

		matches = MatchedMergeSort(matches, func(a *MatchedProfile, b *MatchedProfile) bool {

			var result bool

			switch data.Sort.By {
			case "distance":
				result = a.Distance <= b.Distance
				break

			case "age":
				result = a.Birthdate >= b.Birthdate
				break

			case "popularity":
				result = a.Popularity <= b.Popularity
				break

			default:
				aScore := -a.Distance
				bScore := -b.Distance

				aScore -= aScore * (interestsScore(interests, a.Interests) / 10)
				bScore -= bScore * (interestsScore(interests, b.Interests) / 10)

				aScore -= aScore * (a.Popularity / 20)
				bScore -= bScore * (b.Popularity / 20)

				result = aScore >= bScore
				break
			}

			if data.Sort.Desc {
				return !result
			}
			return result
		})

		var ids []int64
		count := len(matches) - offset
		if count > 0 {
			limit := int(math.Min(float64(count), 30))
			for _, p := range matches[offset : offset+limit] {
				ids = append(ids, p.UserID)
			}
		}
		render.Render(w, r, api.DefaultResponse(http.StatusOK, ids))
	}
}
