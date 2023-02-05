package env

import (
	"flag"
	"os"
	"strconv"
)

var (
	PORT          string = "3000"
	ENABLE_BIGLY  bool   = false
	ENABLE_GAMBLE bool   = false
)

var (
	Token        string
	DATABASE_URL string
	SUID         string
	GID          string
)

func Parse() {
	flag.StringVar(&Token, "t", lookupEnvOrString("DISCORD_TOKEN", Token), "Bot Token")
	flag.StringVar(&PORT, "port", lookupEnvOrString("PORT", PORT), "Health Check Endpoint")
	flag.BoolVar(&ENABLE_BIGLY, "bigly", lookupEnv("ENABLE_BIGLY"), "Feature Flag to Enable Bigly Slash Command")
	flag.BoolVar(&ENABLE_GAMBLE, "ENABLE_GAMBLE", lookupEnv("ENABLE_GAMBLE"), "Feature Flag to Enable Lets Gamble Slash Command")
	flag.StringVar(&DATABASE_URL, "db", lookupEnvOrString("DATABASE_URL", DATABASE_URL), "Database Url")
	flag.StringVar(&SUID, "su", lookupEnvOrString("SU", SUID), "Superuser Discord Id")
	flag.StringVar(&GID, "g", lookupEnvOrString("G", GID), "Guild Discord Id")
	flag.Parse()
}

func lookupEnvOrString(key string, defaultVal string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}

	return defaultVal
}

func lookupEnv(key string) bool {
	if val, ok := os.LookupEnv(key); ok {
		b, err := strconv.ParseBool(val)

		if err != nil {
			return false
		}

		return b
	}

	return false
}
