package main

import (
	"context"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/aurelsandu/StreamLang/operators"
	"github.com/aurelsandu/StreamLang/stream"
)

// Interval: emite 1,2,3,... la fiecare 'every'
func Interval(ctx context.Context, every time.Duration) stream.Stream[int] {
	ch := make(chan stream.Item[int])
	go func() {
		defer close(ch)
		t := time.NewTicker(every)
		defer t.Stop()
		i := 0
		for {
			select {
			case <-ctx.Done():
				return
			case <-t.C:
				i++
				ch <- stream.Item[int]{Value: i}
			}
		}
	}()
	return stream.Stream[int]{Ctx: ctx, Ch: ch}
}

func main() {
	r := gin.Default()

	// mic landing page pentru test
	r.GET("/", func(c *gin.Context) {
		c.Header("Content-Type", "text/html; charset=utf-8")
		c.String(http.StatusOK, `<!doctype html>
<html>
<head><meta charset="utf-8"><title>SSE Int Stream</title></head>
<body>
<h3>StreamLang × SSE — int stream</h3>
<p>Endpoint: <code>/ints?every=200ms</code></p>
<pre id="log"></pre>
<script>
const log = document.getElementById("log");
const es  = new EventSource("/ints?every=200ms");
es.addEventListener("tick", (e) => {
  const d = JSON.parse(e.data);
  log.textContent += d.value + "\n";
});
es.addEventListener("error", (e) => {
  try { const d = JSON.parse(e.data); log.textContent += "[error] " + d.error + "\n"; }
  catch { log.textContent += "[error]\n"; }
});
</script>
</body>
</html>`)
	})

	// SSE: /ints?every=200ms
	r.GET("/ints", func(c *gin.Context) {
		everyStr := c.DefaultQuery("every", "200ms")
		every, err := time.ParseDuration(everyStr)
		if err != nil || every <= 0 {
			every = 200 * time.Millisecond
		}
		ctx := c.Request.Context()

		src := Interval(ctx, every) // Stream[int]

		// opțional: limitează rata spre UI (1/s)
		pipe := operators.ThrottleLatestOp[int](1 * time.Second)
		out := stream.PipeOp(src, pipe)

		// SSE headers
		c.Header("Content-Type", "text/event-stream")
		c.Header("Cache-Control", "no-cache")
		c.Header("Connection", "keep-alive")

		// scriem evenimentele pe măsură ce vin
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
				c.SSEvent("tick", gin.H{"value": it.Value})
				return true
			}
		})
	})

	if err := r.Run(":8080"); err != nil {
		panic(err)
	}
}
