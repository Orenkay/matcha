docker-compose start matcha
docker-compose exec matcha go run ./cmd/seeder/main.go $1
docker-compose restart matcha