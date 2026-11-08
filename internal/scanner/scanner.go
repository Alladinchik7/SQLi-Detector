package scanner

import (
	"context"
	"io"
	"log"
	"net/http"
	"net/url"
	"sync"
	"time"
)

type Task struct {
	target     *url.URL // оригинальный URL
	paramName  string   // в какой параметр инжектим
	paramValue string
	payload    string
}

type Finding struct {
	URL        string    `json:"url"`
	Parameter  string    `json:"parameter"`
	Payload    string    `json:"payload"`
	Vulnerable bool      `json:"vulnerable"`
	Evidence   string    `json:"evidence,omitempty"` // фрагмент тела вокруг сигнатуры, ±80 символов
	StatusCode int       `json:"status_code"`
	Error      string    `json:"error,omitempty"` // сетевая ошибка, если была
	CheckedAt  time.Time `json:"checked_at"`
}

func Run(ctx context.Context, client *http.Client, threads int, cookie string, tasks []Task) []Finding {
	tasksCh := make(chan Task)
	resultsCh := make(chan Finding)

	var wg sync.WaitGroup
	for i := 0; i < threads; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			log.Printf("start worker: %d\n", id)
			for t := range tasksCh {
				resultsCh <- check(ctx, client, cookie, t) // ← вот здесь происходит HTTP
			}
		}(i)
	}

	go func() {
		wg.Wait()
		close(resultsCh)
	}()
	go func() {
		for _, t := range tasks {
			tasksCh <- t
		}
		close(tasksCh)
	}()

	var findings []Finding
	total := len(tasks)
	done := 0
	for f := range resultsCh {
		findings = append(findings, f)
		done++
		if done%10 == 0 || done == total {
			log.Printf("progress: %d/%d tasks", done, total)
		}
	}
	return findings
}

// check — одна задача = один запрос + детекция = один Finding (даже при ошибке)
func check(ctx context.Context, client *http.Client, cookie string, t Task) Finding {
	f := Finding{
		Parameter: t.paramName,
		Payload:   t.payload,
		CheckedAt: time.Now(),
	}

	q := t.target.Query() // Query() возвращает КОПИЮ values — оригинал не трогаем
	q.Set(t.paramName, t.payload)

	u := *t.target // локальная копия URL
	u.RawQuery = q.Encode()
	f.URL = u.String()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		f.Error = err.Error()
		return f
	}
	if cookie != "" {
		req.Header.Set("Cookie", cookie)
	}

	resp, err := client.Do(req)
	if err != nil {
		f.Error = err.Error()
		return f
	}
	defer resp.Body.Close()
	f.StatusCode = resp.StatusCode

	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		f.Error = err.Error()
		return f
	}

	if sig, ok := detectSignature(body); ok {
		f.Vulnerable = true
		f.Evidence = extractEvidence(body, sig)
	}

	return f
}
