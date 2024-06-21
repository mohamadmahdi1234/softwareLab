package handler

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"simpleAPI/db"
)

func FetchHalls(db *db.MySQL) gin.HandlerFunc {
	return func(c *gin.Context) {
		halls, err := db.GetHalls(c)
		if err != nil {
			log.Println(err)
			_ = c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		c.JSON(http.StatusOK, halls)
		return
	}
}
