package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Alladinchik/slqi-detector/internal/config"
	"github.com/Alladinchik/slqi-detector/internal/report"
	"github.com/Alladinchik/slqi-detector/internal/scanner"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg := config.Parse()
	if err := cfg.Validate(); err != nil {
		log.Printf("error validate config: %v", err)
		os.Exit(1)
	}
	targets, err := cfg.LoadTargets()
	if err != nil {
		log.Printf("error load targets: %v", err)
		os.Exit(1)
	}

	var tasks []scanner.Task
	for _, t := range targets {
		built, err := scanner.BuildTasks(t, nil)
		if err != nil {
			log.Printf("skip target %q: %v", t, err)
			continue
		}
		tasks = append(tasks, built...)
	}
	if len(tasks) == 0 {
		log.Printf("no tasks to run")
		os.Exit(1)
	}

	client := &http.Client{Timeout: cfg.Timeout}
	start := time.Now()
	log.Printf("scan started: targets=%d tasks=%d workers=%d time_started=%s\n", len(targets), len(tasks), cfg.Threads, start.Format("15:04:05"))
	findings := scanner.Run(ctx, client, cfg.Threads, cfg.Cookie, tasks)
	finished := time.Now()
	log.Printf("scan finished: targets=%d tasks=%d findings=%d time_finished=%s\n", len(targets), len(tasks), len(findings), finished.Format("15:04:05"))

	summary := report.CreateSummary(start, finished, len(targets), len(tasks), findings)
	if err := summary.Write(cfg.OutputPath); err != nil {
		log.Printf("error write report: %v", err)
		os.Exit(1)
	}
}
