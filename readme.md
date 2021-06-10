1. Copy .env_sample to .env (cp .env_sample .env)
    And add the details in .env file
2. To migrate
    go run ./cmd/web/main.go migrate
3. To seed
    go run ./cmd/web/main.go seed
4. To run the project
    go run ./cmd/web/main.go

5. After seeder is run
login with admin@admin.com/password