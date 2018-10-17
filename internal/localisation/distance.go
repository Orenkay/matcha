package localisation

import (
	"math"

	"github.com/orenkay/matcha/internal/store"
)

func hav(o float64) float64 {
	return math.Pow(math.Sin(o/2), 2)
}

func CoordDistance(o *store.Localisation, d *store.Localisation) float64 {

	olat := o.Lat * math.Pi / 180
	dlat := d.Lat * math.Pi / 180
	olng := o.Lng * math.Pi / 180
	dlng := d.Lng * math.Pi / 180

	r := float64(6378100)

	return 2 * r * math.Asin(math.Sqrt(hav(olat-dlat)+math.Cos(dlat)*math.Cos(olat)*hav(olng-dlng)))
}
