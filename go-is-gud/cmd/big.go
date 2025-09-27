package main

import (
	"context"
	"fmt"
	"guigou/bot-is-gud/db"
	"guigou/bot-is-gud/env"
	"os"
	"time"

	"github.com/guptarohit/asciigraph"
	"go.uber.org/zap"
)

func main() {
	env.Parse()
	if len(os.Args) < 2 {
		fmt.Println("Usage: big <uid>")
		os.Exit(1)
	}

	uid := os.Args[1]
	logger := zap.NewExample().Sugar()
	defer logger.Sync()
	ctx := context.Background()
	conn := db.New(logger, ctx)
	if conn == nil {
		fmt.Fprintln(os.Stderr, "Database connection failed")
		os.Exit(1)
	}

	rows, err := conn.Query(ctx, `SELECT start, expire FROM presences WHERE uid = $1 AND start > NOW() - INTERVAL '30 days'`, uid)
	if err != nil {
		fmt.Println("Error querying presence data:", err)
		os.Exit(1)
	}
	defer rows.Close()

	usage := make(map[string]float64)
	for rows.Next() {
		var start, expire time.Time
		err := rows.Scan(&start, &expire)
		if err != nil {
			fmt.Println("Error scanning row:", err)
			continue
		}
		day := start.Format("2006-01-02")
		duration := expire.Sub(start).Seconds()
		usage[day] += duration
	}

	data := make([]float64, 30)
	today := time.Now()
	for i := 29; i >= 0; i-- {
		day := today.AddDate(0, 0, -i).Format("2006-01-02")
		seconds := usage[day]
		data[29-i] = seconds / 3600.0 // hours
	}

	graph := asciigraph.Plot(data,
		asciigraph.Height(12),
		asciigraph.Caption(fmt.Sprintf("Daily Usage (hours) for %s", uid)),
	)
	fmt.Println(graph)
}
