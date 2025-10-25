package controllers

import (
	"net/http"
	"time"

	"github.com/fabianpoels/fabianpoels-api-go/cache"
	"github.com/fabianpoels/fabianpoels-api-go/collections"
	"github.com/fabianpoels/fabianpoels-api-go/db"
	"github.com/fabianpoels/fabianpoels-api-go/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type AdminController struct {
}

func (ctrl AdminController) Ascents(c *gin.Context) {
	mongoClient := db.GetDbClient()

	cursor, err := collections.GetAscentCollection(mongoClient).Find(c, bson.M{}, options.Find().SetSort(bson.M{"number": -1}))
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	defer cursor.Close(c)

	var ascents []models.Ascent
	if err = cursor.All(c, &ascents); err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, ascents)
}

func (ctrl AdminController) AddAscent(c *gin.Context) {
	var ascent models.Ascent
	err := c.BindJSON(&ascent)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err})
		return
	}

	mongoClient := db.GetDbClient()
	ascentCollection := collections.GetAscentCollection(mongoClient)

	ascent.CreatedAt = time.Now()
	ascent.UpdatedAt = time.Now()
	opts := options.FindOne().SetSort(bson.D{{"number", -1}})
	var lastAscent models.Ascent
	err = ascentCollection.FindOne(c, bson.D{}, opts).Decode(&lastAscent)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			ascent.Number = 1
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to query last ascent"})
			return
		}
	} else {
		ascent.Number = lastAscent.Number + 1
	}

	_, err = ascentCollection.InsertOne(c, ascent)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	cacheService := cache.Service{C: c}
	cacheService.Del(AscentsCacheKey)

	c.JSON(http.StatusCreated, ascent)
}
