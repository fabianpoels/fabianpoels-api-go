package middleware

import (
	"log"
	"net/http"
	"strings"

	"github.com/fabianpoels/fabianpoels-api-go/collections"
	"github.com/fabianpoels/fabianpoels-api-go/db"
	"github.com/fabianpoels/fabianpoels-api-go/models"
	"github.com/fabianpoels/fabianpoels-api-go/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson"
)

type authHeader struct {
	IDToken string `header:"Authorization"`
}

func ValidateJwt() gin.HandlerFunc {
	// load non-request based stuff
	var mongoClient = db.GetDbClient()

	return func(c *gin.Context) {
		h := authHeader{}
		err := c.ShouldBindHeader(&h)

		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"Message": "No token provided"})
			return
		}

		idTokenHeader := strings.Split(h.IDToken, "Bearer ")

		if len(idTokenHeader) < 2 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"Message": "No token provided"})
			return
		}

		// validate the JWT and read the user _id
		tokenString := idTokenHeader[1]
		token, err := utils.ParseJwt(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"Message": "Token invalid"})
			return
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"Message": "Token invalid"})
			return
		}
		email := claims["email"]

		// look up the user
		var user models.User
		err = collections.GetUserCollection(mongoClient).FindOne(c, bson.D{{Key: "email", Value: email}}).Decode(&user)
		if err != nil {
			log.Println(err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"Message": "user not found"})
			return
		}

		c.Set("user", user)
		c.Next()
	}
}

func GetUserFromContext(c *gin.Context) (models.User, bool) {
	userInterface, exists := c.Get("user")
	if !exists {
		return models.User{}, false
	}
	user, ok := userInterface.(models.User)
	return user, ok
}
