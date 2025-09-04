package sources

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/aurelsandu/StreamLang/stream"
)

type Quote struct {
	Symbol string
	Price  float64
	Time   time.Time
}

func YahooQuoteSource(ctx context.Context, symbols []string, every time.Duration) stream.Stream[Quote] {
	if every <= 0 {
		every = 3 * time.Second
	}
	ch := make(chan stream.Item[Quote])

	go func() {
		defer close(ch)

		client := &http.Client{Timeout: 8 * time.Second}
		baseQuote := "https://query1.finance.yahoo.com/v7/finance/quote"
		baseChart := "https://query1.finance.yahoo.com/v8/finance/chart"

		reqHeaders := func(req *http.Request) {
			req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) StreamLang/1.0")
			req.Header.Set("Accept", "application/json,text/plain,*/*")
			req.Header.Set("Accept-Language", "en-US,en;q=0.9")
			req.Header.Set("Connection", "keep-alive")
		}

		ticker := time.NewTicker(every)
		defer ticker.Stop()

		emitErr := func(err error) bool {
			select {
			case <-ctx.Done():
				return false
			case ch <- stream.Item[Quote]{Err: err}:
				return true
			}
		}
		emit := func(q Quote) bool {
			select {
			case <-ctx.Done():
				return false
			case ch <- stream.Item[Quote]{Value: q}:
				return true
			}
		}

		fetchQuoteBatch := func() (handled bool, err error) {
			u, _ := url.Parse(baseQuote)
			qp := u.Query()
			qp.Set("symbols", strings.Join(symbols, ","))
			u.RawQuery = qp.Encode()

			req, _ := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
			reqHeaders(req)
			resp, err := client.Do(req)
			if err != nil {
				return false, fmt.Errorf("yahoo quote request: %w", err)
			}
			defer resp.Body.Close()
			if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
				// lasă caller-ul să facă fallback pe chart
				return false, fmt.Errorf("unauthorized: %s", resp.Status)
			}
			if resp.StatusCode != 200 {
				return false, fmt.Errorf("yahoo quote status: %s", resp.Status)
			}

			var payload struct {
				QuoteResponse struct {
					Result []struct {
						Symbol       string  `json:"symbol"`
						RegularPrice float64 `json:"regularMarketPrice"`
						RegularTime  int64   `json:"regularMarketTime"`
					} `json:"result"`
					Error any `json:"error"`
				} `json:"quoteResponse"`
			}
			if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
				return false, fmt.Errorf("yahoo quote decode: %w", err)
			}

			for _, r := range payload.QuoteResponse.Result {
				ts := time.Unix(r.RegularTime, 0)
				if r.RegularTime == 0 {
					ts = time.Now()
				}
				if !emit(Quote{Symbol: r.Symbol, Price: r.RegularPrice, Time: ts}) {
					return true, nil
				}
			}
			return true, nil
		}

		fetchChartOne := func(sym string) error {
			u := fmt.Sprintf("%s/%s", baseChart, url.PathEscape(sym))
			req, _ := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
			reqHeaders(req)

			q := req.URL.Query()
			q.Set("interval", "1m")
			q.Set("range", "1d")
			req.URL.RawQuery = q.Encode()

			resp, err := client.Do(req)
			if err != nil {
				return fmt.Errorf("yahoo chart %s request: %w", sym, err)
			}
			defer resp.Body.Close()
			if resp.StatusCode != 200 {
				return fmt.Errorf("yahoo chart %s status: %s", sym, resp.Status)
			}

			var payload struct {
				Chart struct {
					Result []struct {
						Meta struct {
							Symbol string  `json:"symbol"`
							// uneori ai price și în Meta
							RegularMarketPrice float64 `json:"regularMarketPrice"`
						} `json:"meta"`
						Timestamp []int64 `json:"timestamp"`
						Indicators struct {
							Quote []struct {
								Close []float64 `json:"close"`
							} `json:"quote"`
						} `json:"indicators"`
					} `json:"result"`
					Error any `json:"error"`
				} `json:"chart"`
			}
			if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
				return fmt.Errorf("yahoo chart %s decode: %w", sym, err)
			}
			if len(payload.Chart.Result) == 0 {
				return errors.New("chart: empty result")
			}
			res := payload.Chart.Result[0]
			// încearcă close-ul cel mai recent nenul
			var price float64
			var ts time.Time
			if len(res.Indicators.Quote) > 0 && len(res.Indicators.Quote[0].Close) > 0 {
				closes := res.Indicators.Quote[0].Close
				tstamps := res.Timestamp
				for i := len(closes) - 1; i >= 0; i-- {
					if closes[i] > 0 {
						price = closes[i]
						if i < len(tstamps) && tstamps[i] > 0 {
							ts = time.Unix(tstamps[i], 0)
						} else {
							ts = time.Now()
						}
						break
					}
				}
			}
			if price == 0 && res.Meta.RegularMarketPrice > 0 {
				price = res.Meta.RegularMarketPrice
				ts = time.Now()
			}
			if price == 0 {
				return fmt.Errorf("chart %s: no price", sym)
			}
			if !emit(Quote{Symbol: sym, Price: price, Time: ts}) {
				return context.Canceled
			}
			return nil
		}

		doFetch := func() {
			handled, err := fetchQuoteBatch()
			if err != nil {
				// 401/403 sau altă eroare — o raportăm și încercăm fallback pe chart
				emitErr(fmt.Errorf("yahoo quote: %v", err))
			}
			if handled {
				return
			}
			// fallback pe chart per simbol
			for _, s := range symbols {
				if err := fetchChartOne(s); err != nil && !errors.Is(err, context.Canceled) {
					emitErr(err)
				}
			}
		}

		// primul fetch imediat
		doFetch()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				doFetch()
			}
		}
	}()

	return stream.Stream[Quote]{Ctx: ctx, Ch: ch}
}