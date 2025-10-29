package handlers

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"strings"

	"golang.org/x/crypto/bcrypt"

	"cabb/internal/store"
	"cabb/pkg/config"
)

// BootstrapAdminUser ensures there is at least one admin account when credentials are provided via env.
func BootstrapAdminUser(ctx context.Context, db *store.DB, cfg config.Config) error {
	if db == nil || db.SQL == nil {
		return nil
	}
	email := strings.TrimSpace(strings.ToLower(cfg.AdminBootstrapEmail))
	password := strings.TrimSpace(cfg.AdminBootstrapPassword)
	if email == "" || password == "" {
		return nil
	}
	displayName := strings.TrimSpace(cfg.AdminBootstrapName)
	if displayName == "" {
		displayName = "Plane Admin"
	}
	user, err := db.GetAdminUserByEmail(ctx, email)
	if err == nil {
		if !user.Active {
			if err := db.UpdateAdminUser(ctx, user.ID, displayName, user.Role, true); err != nil {
				return err
			}
			log.Printf("reactivated bootstrap admin %s", email)
		}
		return nil
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return err
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	if _, err := db.CreateAdminUser(ctx, email, displayName, string(hash), "admin", true); err != nil {
		return err
	}
	log.Printf("bootstrapped admin user %s", email)
	return nil
}
