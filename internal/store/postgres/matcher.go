package postgres

import (
	"database/sql"

	"github.com/lib/pq"

	"github.com/orenkay/matcha/internal/store"
)

const (
	// TheBigQuery78 hmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmm
	// $1 = userId
	// $2 = Interests array
	// $3 = Gender array
	// $4 = Attraction array
	// $5 = Lat float64
	// $6 = Lng loat64
	TheBigQuery78 = `
	SELECT
		t1.userId

		FROM profiles t1

		LEFT JOIN users_interests t2 ON t1.userId = t2.userId AND (t2.value = any($2))
		LEFT JOIN localisations t3 ON t1.userId = t3.userId
		LEFT JOIN users_likes t4 ON t1.userId = t4.targetId
		LEFT JOIN users_likes t5 ON t5.userId = $1 AND t5.targetId = t1.userId

		WHERE
			PROFILE_CHECK(t1.userId) AND
			-- COORD_DISTANCE($5, t3.lat, $6, t3.lng) <= 30000 AND
			t1.gender = any($3) AND
			t1.attraction = any($4) AND
			t1.userId <> $1 AND
			t5 IS NULL

		GROUP BY
			t1.userId,
			t3.lat,
			t3.lng

		ORDER BY
			COORD_DISTANCE($5, t3.lat, $6, t3.lng),
			COUNT(distinct t2.id) DESC,
			COUNT(distinct t4.id) DESC

	LIMIT 30
	`
)

type MatcherService struct {
	db *sql.DB
}

// func NewMatcherService(db *sql.DB) store.MatcherService {
// 	return &MatcherService{
// 		db: db,
// 	}
// }

func (s *MatcherService) List(p *store.Profile, interests []*store.Interest, loc *store.Localisation) ([]int64, error) {
	var (
		_genders     []string
		_attractions []string
		_interests   []string
		ids          []int64
	)

	// gender
	{
		if p.Attraction == "bi" {
			_genders = []string{"male", "female"}
		}
		if (p.Gender == "male" && p.Attraction == "hetero") || (p.Gender == "female" && p.Attraction == "homo") {
			_genders = []string{"female"}
		}
		if (p.Gender == "male" && p.Attraction == "homo") || (p.Gender == "female" && p.Attraction == "hetero") {
			_genders = []string{"male"}
		}
	}

	// attraction
	{
		if p.Attraction == "bi" {
			_attractions = []string{"hetero", "homo", "bi"}
		} else {
			_attractions = []string{p.Attraction, "bi"}
		}
	}

	for _, i := range interests {
		_interests = append(_interests, i.Value)
	}

	rows, err := s.db.Query(TheBigQuery78, p.UserID, pq.Array(_interests), pq.Array(_genders), pq.Array(_attractions), loc.Lat, loc.Lng)
	{
		if err != nil {
			return nil, err
		}
	}
	defer rows.Close()

	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}
