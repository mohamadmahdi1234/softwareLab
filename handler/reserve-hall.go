package handler

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"simpleAPI/db"
	"simpleAPI/db/entity"
	"time"
)

func Resp(status int, message string) map[string]interface{} {
	return map[string]interface{}{"status": status, "data": map[string]string{"message": message}}
}

func ReserveHall(db *db.MySQL) gin.HandlerFunc {
	return func(c *gin.Context) {
		var reserveReq entity.ReserveReq
		if err := c.ShouldBindJSON(&reserveReq); err != nil {
			c.JSON(http.StatusUnprocessableEntity, Resp(http.StatusUnprocessableEntity, err.Error()))
			return
		}
		date, err := time.Parse("2006-01-02", reserveReq.Date)
		if err != nil {
			c.JSON(http.StatusBadRequest, Resp(http.StatusUnprocessableEntity, err.Error()))
			return
		}

		err = db.ReserveHall(c, reserveReq, date)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, Resp(http.StatusInternalServerError, err.Error()))
			return
		}
		c.JSON(http.StatusOK, Resp(http.StatusOK, "success"))
		return
	}
}
