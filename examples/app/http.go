package main

import (
	"context"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/aurelsandu/StreamLang/operators"
	"github.com/aurelsandu/StreamLang/sinks"
	"github.com/aurelsandu/StreamLang/sources"
	"github.com/aurelsandu/StreamLang/stream"
)

func main() {
	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// GET /sum?nums=1,2,3    sau    /sum?n=5   (=> 1..5)
	r.GET("/sum", func(c *gin.Context) {
		nums, err := readNumbers(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
		defer cancel()

		src := sources.FromSlice(ctx, nums)
		fold := operators.FoldOp[int, int](0, func(acc, x int) int { return acc + x })
		out := stream.PipeOp(src, fold)

		vals, err := sinks.ToSlice(out)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		sum := 0
		if len(vals) > 0 {
			sum = vals[0]
		}
		c.JSON(http.StatusOK, gin.H{
			"input": nums,
			"sum":   sum,
		})
	})

	// GET /scan?nums=1,2,3    sau    /scan?n=5  (=> 1..5)
	r.GET("/scan", func(c *gin.Context) {
		nums, err := readNumbers(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
		defer cancel()

		src := sources.FromSlice(ctx, nums)
		scan := operators.ScanOp[int, int](0, func(acc, x int) int { return acc + x })
		out := stream.PipeOp(src, scan)

		vals, err := sinks.ToSlice(out)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"input":        nums,
			"running_sum":  vals, // ex: [1,3,6,10]
			"final_sum":    lastOrZero(vals),
			"count_inputs": len(nums),
		})
	})

	// landing mic pentru test rapid
	r.GET("/", func(c *gin.Context) {
		c.Header("Content-Type", "text/html; charset=utf-8")
		c.String(http.StatusOK, `<!doctype html>
<html><head><meta charset="utf-8"><title>StreamLang HTTP demo</title></head>
<body>
<h3>StreamLang × Gin — HTTP demo</h3>
<ul>
  <li><a href="/sum?nums=1,2,3,4,5">/sum?nums=1,2,3,4,5</a></li>
  <li><a href="/sum?n=5">/sum?n=5</a> (generează 1..5)</li>
  <li><a href="/scan?nums=1,2,3,4,5">/scan?nums=1,2,3,4,5</a></li>
  <li><a href="/scan?n=5">/scan?n=5</a></li>
</ul>
</body></html>`)
	})

	if err := r.Run(":8080"); err != nil {
		panic(err)
	}
}

func readNumbers(c *gin.Context) ([]int, error) {
	if raw := strings.TrimSpace(c.Query("nums")); raw != "" {
		parts := strings.Split(raw, ",")
		out := make([]int, 0, len(parts))
		for _, p := range parts {
			p = strings.TrimSpace(p)
			if p == "" {
				continue
			}
			n, err := strconv.Atoi(p)
			if err != nil {
				return nil, err
			}
			out = append(out, n)
		}
		return out, nil
	}
	// fallback: ?n=K  => [1..K]
	if nk := strings.TrimSpace(c.DefaultQuery("n", "")); nk != "" {
		k, err := strconv.Atoi(nk)
		if err != nil || k < 0 {
			if err == nil {
				err = strconv.ErrSyntax
			}
			return nil, err
		}
		out := make([]int, k)
		for i := 1; i <= k; i++ {
			out[i-1] = i
		}
		return out, nil
	}
	// implicit: gol
	return []int{}, nil
}

func lastOrZero(xs []int) int {
	if len(xs) == 0 {
		return 0
	}
	return xs[len(xs)-1]
}