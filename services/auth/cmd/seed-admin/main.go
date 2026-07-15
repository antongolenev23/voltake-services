package main

import (
	"context"
	"log"
	"os"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/antongolenev23/voltake-services/services/auth/internal/config"
	"github.com/antongolenev23/voltake-services/services/auth/internal/repository/postgres"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cfg := config.MustLoadSeed()

	pool, err := postgres.NewPgxpool(ctx, &cfg.Repository)

	if err != nil {
		log.Fatalf("failed to connect db: %v", err)
	}

	defer pool.Close()

	repo := postgres.New(pool)

	email := os.Getenv("ADMIN_EMAIL")
	password := os.Getenv("ADMIN_PASSWORD")

	if email == "" || password == "" {
		log.Fatal("ADMIN_EMAIL and ADMIN_PASSWORD required")
	}

	exists, err := repo.ExistsByEmail(ctx, email)

	if err != nil {
		log.Fatalf("failed check admin: %v", err)
	}

	if exists {
		log.Println("admin already exists")
		return
	}

	hash, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)

	if err != nil {
		log.Fatalf("failed hash password: %v", err)
	}

	_, err = repo.SaveAdmin(ctx, email, hash)

	if err != nil {
		log.Fatalf("failed create admin: %v", err)
	}

	log.Println("admin created")
}
