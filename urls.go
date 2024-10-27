package main

import (
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func URLDefaultHandler(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	purgeExpiredTags(db, c.MustGet("expiration").(string))

	c.Redirect(http.StatusMovedPermanently, c.MustGet("defaultSite").(string))
}

func URLRequestHandler(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	purgeExpiredTags(db, c.MustGet("expiration").(string))

	url := URL{}
	db.Where("tag = ?", c.Param("tag")).First(&url)

	if url.ID == 0 {
		c.Redirect(http.StatusMovedPermanently, c.MustGet("defaultSite").(string))
	} else {
		c.Redirect(http.StatusMovedPermanently, url.URL)
	}
}

func URLCreateHandler(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	purgeExpiredTags(db, c.MustGet("expiration").(string))

	type Header struct {
		Authorization string `header:"Authorization"`
	}

	h := Header{}
	if err := c.ShouldBindHeader(&h); err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	var token string
	authorization := strings.Split(h.Authorization, "Token ")
	if len(authorization) == 2 {
		token = authorization[1]
	} else {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	if token != c.MustGet("token").(string) {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	type request struct {
		URL string `json:"url" binding:"required"`
	}

	var r request
	if err := c.BindJSON(&r); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	url := URL{}
	db.Where("url = ?", r.URL).First(&url)
	if url.ID != 0 {
		c.JSON(http.StatusCreated, gin.H{
			"id":         url.ID,
			"created_at": url.CreatedAt,
			"tag":        url.Tag,
			"url":        url.URL,
		})
		return
	}

	var tag string
	attempts := 0
	for {
		tag = randomTag()

		url := URL{}
		db.Where("tag = ?", tag).First(&url)

		if url.ID == 0 {
			break
		} else {
			attempts++
			if attempts >= 10 {
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}
		}
	}

	newUrl := URL{
		Tag: tag,
		URL: r.URL,
	}
	db.Create(&newUrl)

	c.JSON(http.StatusCreated, gin.H{
		"id":         newUrl.ID,
		"created_at": newUrl.CreatedAt,
		"tag":        newUrl.Tag,
		"url":        newUrl.URL,
	})
}

func randomTag() string {
	runes := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_")
	tag := make([]rune, 4)

	for i := range tag {
		tag[i] = runes[rand.Intn(len(runes))]
	}

	return string(tag)
}

func purgeExpiredTags(db *gorm.DB, expiration string) {
	d, _ := time.ParseDuration(expiration)

	db.Where("created_at <= ?", time.Now().UTC().Add(d)).Delete(&URL{})
}
