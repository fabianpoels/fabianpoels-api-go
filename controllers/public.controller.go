package controllers

import (
	"net/http"

	"github.com/fabianpoels/fabianpoels-api-go/collections"
	"github.com/fabianpoels/fabianpoels-api-go/db"
	"github.com/fabianpoels/fabianpoels-api-go/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type PublicController struct{}

func (ctrl PublicController) Ascents(c *gin.Context) {
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

	c.JSON(http.StatusOK, publicAscents)
}
