package services

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"os"
	"wardenapi/internal/models"
	"wardenapi/internal/utils"
)

func MapRunePriceRoutes(router *gin.Engine) {
	router.POST("/runeprice", InsertRunePrice)
	router.GET("/runeprice", GetLatestRunePrice)
}

func InsertRunePrice(c *gin.Context) {
	var runes []models.RunePrice
	c.BindJSON(&runes)

	for _, runePrice := range runes {
		conn := utils.GetConnection()
		defer conn.Close(context.Background())

		sql := `
			INSERT INTO public.runeprice (rune_name, server, date, price)
			VALUES ($1, $2, $3, $4);
		`

		utils.DoRequest(conn, sql, runePrice.Rune_name, runePrice.Server, runePrice.Date, runePrice.Price)
	}

	c.JSON(200, gin.H{
		"message": "done",
	})
}

func GetLatestRunePrice(c *gin.Context) {
	server := c.Query("server")

	conn := utils.GetConnection()
	defer conn.Close(context.Background())
	sql := `
		SELECT runeprice.rune_name, runeprice.server, latest.date, price
		FROM runeprice
         JOIN (SELECT rune_name, server, MAX(date) AS date
               FROM runeprice
               GROUP BY rune_name, server)
    	 AS latest ON runeprice.rune_name = latest.rune_name AND runeprice.date = latest.date AND
                 runeprice.server = latest.server
		WHERE runeprice.server = $1
	`
	rows := utils.DoRequest(conn, sql, server)
	var results []models.RunePrice
	for rows.Next() {
		runePrice := models.RunePrice{}
		err := rows.Scan(&runePrice.Rune_name, &runePrice.Server, &runePrice.Date, &runePrice.Price)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Query failed: %v\n", err)
		} else {
			results = append(results, runePrice)
		}
	}

	c.JSON(200, gin.H{
		"results": results,
	})
}
