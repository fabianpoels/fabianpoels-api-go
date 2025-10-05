package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/fabianpoels/fabianpoels-api-go/cache"
	"github.com/fabianpoels/fabianpoels-api-go/collections"
	"github.com/fabianpoels/fabianpoels-api-go/db"
	"github.com/fabianpoels/fabianpoels-api-go/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type PublicController struct{}

const cacheKey = "ascents:public"

func (ctrl PublicController) Ascents(c *gin.Context) {
	// try to load from cache
	cacheAscents, err := publicAscentsFromCache(c)
	if err == nil {
		c.JSON(http.StatusOK, cacheAscents)
		return
	}

	// load from DB
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

	publicAscents := make([]models.PublicAscent, len(ascents))
	for i, ascent := range ascents {
		publicAscents[i] = models.SerializeAscent(ascent)
	}

	storePublicAscentsInCache(c, publicAscents)

	c.JSON(http.StatusOK, publicAscents)
}

func publicAscentsFromCache(c *gin.Context) (ascents []models.PublicAscent, err error) {
	cacheService := cache.Service{C: c}
	cacheString, err := cacheService.Get(cacheKey)
	if err != nil {
		return ascents, err
	}

	err = json.Unmarshal([]byte(cacheString), &ascents)
	if err != nil {
		return ascents, err
	}

	return ascents, nil
}

func storePublicAscentsInCache(c *gin.Context, ascents []models.PublicAscent) {
	cacheService := cache.Service{C: c}
	cacheString, _ := json.Marshal(ascents)
	cacheService.Set(cacheKey, cacheString, 0)
}
