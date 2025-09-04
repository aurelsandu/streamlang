package main

import (
//	"context"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/aurelsandu/StreamLang/operators"
	"github.com/aurelsandu/StreamLang/sources"
	"github.com/aurelsandu/StreamLang/stream"
)

func main() {
	router := gin.Default()

// SSE: /stocks?symbols=AAPL,MSFT,GOOG&every=2s&provider=stooq|yahoo|mock
	router.GET("/stocks", func(c *gin.Context) {
		symbols := c.Query("symbols")
		if symbols == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "query param 'symbols' is required, e.g. AAPL,MSFT"})
			return
		}
		everyStr := c.DefaultQuery("every", "3s")
		every, err := time.ParseDuration(everyStr)
		if err != nil || every <= 0 {
			every = 3 * time.Second
		}
		provider := c.DefaultQuery("provider", "stooq") // implicit: stooq (mai stabil)

		ctx := c.Request.Context()
		var src stream.Stream[sources.Quote]

		switch provider {
		case "yahoo":
			src = sources.YahooQuoteSource(ctx, splitCSV(symbols), every)
		case "mock":
			src = sources.MockQuoteSource(ctx, splitCSV(symbols), every)
		default: // "stooq"
			src = sources.StooqQuoteSource(ctx, splitCSV(symbols), every)
		}

		// throttling pentru UI
		out := stream.PipeOp(src, operators.ThrottleLatestOp[sources.Quote](500*time.Millisecond))

		// SSE headers
		c.Header("Content-Type", "text/event-stream")
		c.Header("Cache-Control", "no-cache")
		c.Header("Connection", "keep-alive")

		c.Stream(func(w io.Writer) bool {
			select {
				case <-ctx.Done():
					return false
				case it, ok := <-out.Ch:
					if !ok {
						return false
					}
					if it.Err != nil {
						c.SSEvent("error", gin.H{"error": it.Err.Error()})
						return true
					}
					c.SSEvent("quote", gin.H{
						"ts":     it.Value.Time.UnixMilli(),
						"symbol": it.Value.Symbol,
						"price":  it.Value.Price,
					})
					return true
				}
		})
	})

	// Landing simplu de test
	router.GET("/", func(c *gin.Context) {
		c.Header("Content-Type", "text/html; charset=utf-8")
		c.String(http.StatusOK, `<!doctype html>
<html>
<head><meta charset="utf-8"><title>Stocks SSE</title></head>
<body>
<h3>StreamLang × Gin — Stocks (Yahoo polling)</h3>
<p>Exemplu: <a href="/stocks?symbols=AAPL,MSFT,GOOG&every=2s">/stocks?symbols=AAPL,MSFT,GOOG&every=2s</a></p>
<pre id="log"></pre>
<script>
const log = document.getElementById("log");
const es = new EventSource("/stocks?symbols=AAPL,MSFT,GOOG&every=2s");
es.addEventListener("quote", (e) => {
  const d = JSON.parse(e.data);
  log.textContent += d.symbol + "  " + d.price.toFixed(2) + "  @" + new Date(d.ts).toLocaleTimeString() + "\n";
});
es.addEventListener("error", (e) => {
  try { const d = JSON.parse(e.data); log.textContent += "[error] " + d.error + "\n"; }
  catch { log.textContent += "[error]\n"; }
});
</script>
</body>
</html>`)
	})


	if err := router.Run(":8080"); err != nil {
		panic(err)
	}
}

// splitCSV: mic utilitar pentru a separa "AAPL,MSFT,GOOG" în []string
func splitCSV(s string) []string {
	var out []string
	cur := ""
	for _, r := range s {
		if r == ',' || r == ';' || r == ' ' {
			if cur != "" {
				out = append(out, cur)
				cur = ""
			}
			continue
		}
		cur += string(r)
	}
	if cur != "" {
		out = append(out, cur)
	}
	return out
}