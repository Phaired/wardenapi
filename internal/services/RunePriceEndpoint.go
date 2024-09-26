package services

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"os"
	"time"
	"wardenapi/internal/models"
	"wardenapi/internal/utils"
)

func MapRunePriceRoutes(router *gin.Engine) {
	router.POST("/runeprice", func(c *gin.Context) {
		utils.Filter(c, InsertRunePrice)
	})
	router.GET("/runeprice", func(c *gin.Context) {
		utils.Filter(c, GetLatestRunePrice)
	})
	router.GET("/runeprice/history", func(c *gin.Context) {
		utils.Filter(c, GetRunePriceHistory)
	})
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

func GetRunePriceHistory(c *gin.Context) {
	server := c.Query("server")
	runeName := c.Query("rune_name")
	startDate, _ := time.Parse(time.RFC3339, c.Query("start_date"))
	endDate, _ := time.Parse(time.RFC3339, c.Query("end_date"))

	conn := utils.GetConnection()
	defer conn.Close(context.Background())
	sql := `
		SELECT rune_name, server, runeprice.date, price
		FROM runeprice
		WHERE server = $1
		AND rune_name = $2
		AND CASE WHEN $3 = timestamp '0001-01-01 00:00:00' THEN TRUE ELSE runeprice.date >= $3 END
		AND CASE WHEN $4 = timestamp '0001-01-01 00:00:00' THEN TRUE ELSE runeprice.date <= $4 END
		ORDER BY date DESC
	`
	rows := utils.DoRequest(conn, sql, server, runeName, startDate, endDate)
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
