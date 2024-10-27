package main

import (
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type ShortestConfig struct {
	DB          *gorm.DB
	DefaultSite string
	Expiration  string
	Token       string
}

func shortestMiddleware(config *ShortestConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("db", config.DB)
		c.Set("defaultSite", config.DefaultSite)
		c.Set("expiration", config.Expiration)
		c.Set("token", config.Token)

		c.Next()
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	db, err := gorm.Open(postgres.Open(os.Getenv("DATABASE_DSN")), &gorm.Config{})
	if err != nil {
		log.Fatalf("unable to connect to database: %v", err)
	}

	defaultSite := os.Getenv("DEFAULT_SITE")
	if len(defaultSite) == 0 {
		log.Fatal("DEFAULT_SITE not set")
	}

	expiration := os.Getenv("EXPIRATION")
	if len(expiration) == 0 {
		log.Fatal("EXPIRATION not set")
	}

	if _, err := time.ParseDuration(expiration); err != nil {
		log.Fatalf("invalid EXPIRATION: %v", err)
	}

	token := os.Getenv("TOKEN")
	if len(token) == 0 {
		log.Fatalf("TOKEN not set")
	}

	db.AutoMigrate(&URL{})

	r := gin.Default()
	r.Use(shortestMiddleware(&ShortestConfig{
		DB:          db,
		DefaultSite: defaultSite,
		Token:       token,
		Expiration:  expiration,
	}))

	r.GET("/:tag", URLRequestHandler)

	r.GET("/", URLDefaultHandler)
	r.POST("/", URLCreateHandler)

	r.Run(":8000")
}
