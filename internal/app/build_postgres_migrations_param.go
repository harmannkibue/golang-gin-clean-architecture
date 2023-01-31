package app

import (
	"github.com/harmannkibue/golang_gin_clean_architecture/config"
)

// BuildPostgresUrl const conString = "postgres://YourUserName:YourPassword@YourHostname:5432/YourDatabaseName?sslmode=disabled";
func BuildPostgresUrl(cfg *config.Config) string {

	// keys to be passed to postgres  for initialization
	keys := []string{"//", ":", "@", ":", "/"}
	// Load string with environment variables for postgres
	values := []string{cfg.VaUser, cfg.VaPass, cfg.VaHost, cfg.VaPort, cfg.VaDb}
	returnString := "postgres:"

	for i, ky := range keys {
		returnString += ky + values[i]
	}

	return returnString
}
