package main

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"time"
	"io"
	"github.com/gin-gonic/gin"

	"github.com/aurelsandu/StreamLang/operators"
//	"github.com/aurelsandu/StreamLang/sinks"
	"github.com/aurelsandu/StreamLang/stream"
)

// randomSensor produce valori întregi aleatorii la interval fix (simulează un senzor).
func randomSensor(ctx context.Context, min, max int, every time.Duration) stream.Stream[int] {
	ch := make(chan stream.Item[int])
	go func() {
		defer close(ch)
		ticker := time.NewTicker(every)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				val := rand.Intn(max-min+1) + min
				ch <- stream.Item[int]{Value: val}
			}
		}
	}()
	return stream.Stream[int]{Ctx: ctx, Ch: ch}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	router := gin.Default()

	// mic landing page ca să testezi rapid SSE în browser
	router.GET("/", func(c *gin.Context) {
		c.Header("Content-Type", "text/html; charset=utf-8")
		c.String(http.StatusOK, `<!doctype html>
<html>
<head><meta charset="utf-8"><title>StreamLang SSE Demo</title></head>
<body>
<h3>StreamLang × Gin (SSE)</h3>
<pre id="log"></pre>
<script>
  const log = document.getElementById("log");
  const es = new EventSource("/stream");
  es.addEventListener("data", (e) => {
    const d = JSON.parse(e.data);
    log.textContent += "[" + new Date(d.ts).toLocaleTimeString() + "] " + d.value + "\n";
  });
  es.addEventListener("error", (e) => {
    log.textContent += "[error] " + e.data + "\n";
  });
</script>
</body>
</html>`)
	})

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Endpointul SSE: combinăm doi senzori și debouncem la 500ms.
	router.GET("/stream", func(c *gin.Context) {
		// context legat de request; când clientul închide conexiunea, se oprește fluxul
		ctx := c.Request.Context()

		// doi senzori: temperatură și umiditate
		temp := randomSensor(ctx, 15, 30, 200*time.Millisecond) // °C
		hum := randomSensor(ctx, 40, 70, 350*time.Millisecond)  // %

		// CombineLatest(temp, hum) => string "climate: T°C / H%"
		pipe := operators.CombineLatest2Op[int, int, string](hum, func(t, h int) string {
			return fmt.Sprintf("climate: %d°C / %d%%", t, h)
		})
		combined := stream.PipeOp(temp, pipe)

		// Debounce 500ms (să nu inundăm clientul)
		final := stream.PipeOp(combined, operators.ThrottleLatestOp[string](500*time.Millisecond))
	//	final := stream.PipeOp(combined, operators.DebounceOp[string](500*time.Millisecond))
		
		// Setăm SSE headers
		c.Header("Content-Type", "text/event-stream")
		c.Header("Cache-Control", "no-cache")
		c.Header("Connection", "keep-alive")

		// Emitere SSE: folosim c.Stream pentru a scrie evenimentele pe măsură ce vin
		c.Stream(func(w io.Writer) bool {
			select {
			case <-ctx.Done():
				return false
			case it, ok := <-final.Ch:
				if !ok {
					return false
				}
				if it.Err != nil {
					c.SSEvent("error", it.Err.Error())
					return true
				}
				c.SSEvent("data", gin.H{
					"ts":    time.Now().UnixMilli(),
					"value": it.Value,
				})
				return true
			}
		})
	})

	// Pornire server
	if err := router.Run(":8080"); err != nil {
		panic(err)
	}
}