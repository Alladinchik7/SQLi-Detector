package config

import (
	"bufio"
	"errors"
	"flag"
	"os"
	"strings"
	"time"
)

// internal/config
type Config struct {
	TargetURL   string        // -u
	TargetsFile string        // -f
	Threads     int           // -t, default 10
	OutputPath  string        // -o, default "" (stdout)
	Cookie      string        // -cookie, сырое значение заголовка Cookie
	Timeout     time.Duration // -timeout, default 10s
}

var (
	errorAbsentTarget   = errors.New("укажите цель: -u <url> ИЛИ -f <file>")
	errorWrongTarget    = errors.New("-u и -f одновременно указывать нельзя, выберите один источник")
	errorInvalidThreads = errors.New("-t должен быть >= 1")
	errorInvalidTimeout = errors.New("-timeout должен быть > 0")
)

func Parse() *Config {
	cfg := &Config{Threads: 10, Timeout: 10 * time.Second}
	flag.StringVar(&cfg.TargetURL, "u", "", "один URL для проверки")
	flag.StringVar(&cfg.TargetsFile, "f", "", "файл со списком URL")
	flag.IntVar(&cfg.Threads, "t", 10, "число воркеров")
	flag.StringVar(&cfg.OutputPath, "o", "", "путь к JSON-отчёту (по умолчанию stdout)")
	flag.StringVar(&cfg.Cookie, "cookie", "", "сырое значение заголовка Cookie")
	flag.DurationVar(&cfg.Timeout, "timeout", 10*time.Second, "таймаут запроса")
	flag.Parse()

	return cfg
}

func (c *Config) Validate() error {
	switch {
	case c.TargetURL == "" && c.TargetsFile == "":
		return errorAbsentTarget
	case c.TargetURL != "" && c.TargetsFile != "":
		return errorWrongTarget
	case c.Threads < 1:
		return errorInvalidThreads
	case c.Timeout <= 0:
		return errorInvalidTimeout
	}

	return nil
}

func (c *Config) LoadTargets() ([]string, error) {
	var targets []string
	if c.TargetsFile != "" {
		file, err := os.Open(c.TargetsFile)
		if err != nil {
			return nil, err
		}
		defer file.Close()

		sc := bufio.NewScanner(file)
		for sc.Scan() {
			line := strings.TrimSpace(sc.Text())

			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}

			targets = append(targets, line)
		}

		if sc.Err() != nil {
			return nil, sc.Err()
		}
	} else {
		targets = append(targets, c.TargetURL)
	}

	return targets, nil
}
