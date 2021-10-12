package config

import (
	"os"
	"testing"
)

func TestLoadYaml(t *testing.T) {
	Load()
	c := Get()

	if c.Logging.Level != "debug" {
		t.Fatal("Config must have been loaded from yaml")
	}
}

func TestLoadEnv(t *testing.T) {
	if err := os.Setenv("EMVI_WIKI_LOGLEVEL", "debug"); err != nil {
		t.Fatal(err)
	}

	if err := os.Setenv("EMVI_WIKI_MIGRATE_DB_HOST", "migrate_db"); err != nil {
		t.Fatal(err)
	}

	loadConfigFromEnv()
	c := Get()

	if c.Logging.Level != "debug" {
		t.Fatalf("Config must have been loaded from env, but was: %v", c.Logging.Level)
	}

	if c.Migrate.Database.Host != "migrate_db" {
		t.Fatalf("Config database must have been loaded from env, but was: %v", c.Migrate.Database.Host)
	}
}

func TestGetEnv(t *testing.T) {
	if v := getEnv("NAME", "default"); v != "default" {
		t.Fatalf("Default value must have been returned, but was: %v", v)
	}

	if err := os.Setenv("EMVI_WIKI_TEST", "test"); err != nil {
		t.Fatal(err)
	}

	if v := getEnv("TEST", "default"); v != "test" {
		t.Fatalf("Value must have been returned, but was: %v", v)
	}
}

func TestGetEnvInt(t *testing.T) {
	if v := getEnvInt("NAME", 42); v != 42 {
		t.Fatalf("Default value must have been returned, but was: %v", v)
	}

	if err := os.Setenv("EMVI_WIKI_TEST", "33"); err != nil {
		t.Fatal(err)
	}

	if v := getEnvInt("TEST", 42); v != 33 {
		t.Fatalf("Value must have been returned, but was: %v", v)
	}
}

func TestGetEnvBool(t *testing.T) {
	if v := getEnvBool("NAME", true); !v {
		t.Fatalf("Default value must have been returned, but was: %v", v)
	}

	if err := os.Setenv("EMVI_WIKI_TEST", "true"); err != nil {
		t.Fatal(err)
	}

	if v := getEnvBool("TEST", false); !v {
		t.Fatalf("Value must have been returned, but was: %v", v)
	}
}
