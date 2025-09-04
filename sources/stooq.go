package sources

import (
	"context"
	"encoding/csv"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/aurelsandu/StreamLang/stream"
)

// StooqQuoteSource: ia snapshot prețuri din stooq.com în format CSV (polling).
// symbols: exemplu []{"AAPL.US","MSFT.US"}; dacă nu pui .US, îl adăugăm noi.
func StooqQuoteSource(ctx context.Context, symbols []string, every time.Duration) stream.Stream[Quote] {
	if every <= 0 {
		every = 3 * time.Second
	}
	// normalizăm sufixul .US dacă lipsește
	norm := make([]string, 0, len(symbols))
	for _, s := range symbols {
		s = strings.TrimSpace(s)
		if s == "" {
			continue
		}
		if !strings.Contains(s, ".") {
			s += ".US"
		}
		norm = append(norm, strings.ToLower(s))
	}

	ch := make(chan stream.Item[Quote])
	go func() {
		defer close(ch)
		client := &http.Client{Timeout: 10 * time.Second}
		ticker := time.NewTicker(every)
		defer ticker.Stop()

		emit := func(q Quote) bool {
			select {
			case <-ctx.Done():
				return false
			case ch <- stream.Item[Quote]{Value: q}:
				return true
			}
		}
		emitErr := func(err error) bool {
			select {
			case <-ctx.Done():
				return false
			case ch <- stream.Item[Quote]{Err: err}:
				return true
			}
		}

		fetch := func() {
			if len(norm) == 0 {
				emitErr(fmt.Errorf("stooq: no symbols"))
				return
			}
			// ex: https://stooq.com/q/l/?s=aapl.us,msft.us&f=sd2t2ohlcv&h&e=csv
			u, _ := url.Parse("https://stooq.com/q/l/")
			q := u.Query()
			q.Set("s", strings.Join(norm, ","))
			q.Set("f", "sd2t2ohlcv") // include symbol, date, time, open, high, low, close, volume
			q.Set("h", "")
			q.Set("e", "csv")
			u.RawQuery = q.Encode()

			req, _ := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
			req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; StreamLang/1.0)")
			resp, err := client.Do(req)
			if err != nil {
				emitErr(fmt.Errorf("stooq request: %w", err))
				return
			}
			defer resp.Body.Close()
			if resp.StatusCode != 200 {
				emitErr(fmt.Errorf("stooq status: %s", resp.Status))
				return
			}

			r := csv.NewReader(resp.Body)
			// header: Symbol,Date,Time,Open,High,Low,Close,Volume
			records, err := r.ReadAll()
			if err != nil {
				emitErr(fmt.Errorf("stooq csv: %w", err))
				return
			}
			if len(records) <= 1 {
				emitErr(fmt.Errorf("stooq: empty payload"))
				return
			}
			now := time.Now()
			for i := 1; i < len(records); i++ {
				row := records[i]
				if len(row) < 7 {
					continue
				}
				sym := strings.ToUpper(strings.TrimSuffix(row[0], ".US"))
				closeStr := row[6]
				var price float64
				fmt.Sscanf(closeStr, "%f", &price)
				if price <= 0 {
					continue
				}
				if !emit(Quote{Symbol: sym, Price: price, Time: now}) {
					return
				}
			}
		}

		fetch() // imediat
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				fetch()
			}
		}
	}()
	return stream.Stream[Quote]{Ctx: ctx, Ch: ch}
}