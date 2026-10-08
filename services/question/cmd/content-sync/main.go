// Command content-sync validates authored content and loads it into the database.
//
//	content-sync --validate   validate only (no writes); checks skills against the DB when DATABASE_URL is set
//	content-sync              validate, then idempotently upsert worlds, nodes, and lessons
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/prepio/prepio/services/question/internal/lesson"
	"github.com/prepio/prepio/services/question/internal/store"
	"github.com/prepio/prepio/shared/postgres"
)

func main() {
	root := flag.String("content", envOrDefault("CONTENT_DIR", "content"), "path to the content directory")
	validateOnly := flag.Bool("validate", false, "validate content and exit without writing")
	flag.Parse()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	content, err := lesson.Load(*root)
	if err != nil {
		log.Fatalf("content-sync: load: %v", err)
	}

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" && !*validateOnly {
		dsn = "postgres://prepio:prepio@localhost:5432/prepio?sslmode=disable"
	}

	var syncStore *store.ContentSyncStore
	catalog := lesson.Catalog{}
	if dsn != "" {
		pool, err := postgres.New(ctx, dsn)
		if err != nil {
			log.Fatalf("content-sync: postgres: %v", err)
		}
		defer pool.Close()
		syncStore = store.NewContentSyncStore(pool)
		catalog, err = syncStore.Catalog(ctx)
		if err != nil {
			log.Fatalf("content-sync: %v", err)
		}
	} else {
		log.Println("content-sync: DATABASE_URL not set, skipping skill and topic catalog checks")
	}

	if issues := lesson.Validate(content, catalog); len(issues) > 0 {
		fmt.Fprintf(os.Stderr, "content-sync: %d validation problem(s):\n  - %s\n", len(issues), strings.Join(issues, "\n  - "))
		os.Exit(1)
	}
	log.Printf("content-sync: %d world(s), %d lesson(s) valid", len(content.Worlds), len(content.Lessons))
	if *validateOnly {
		return
	}

	report, err := syncStore.Sync(ctx, content)
	if err != nil {
		log.Fatalf("content-sync: sync: %v", err)
	}
	log.Printf("content-sync: worlds=%d nodes=%d lessons created=%d new_version=%d unchanged=%d deprecated(lessons=%d worlds=%d nodes=%d)",
		report.Worlds, report.Nodes, report.LessonsCreated, report.LessonsNewVersion, report.LessonsUnchanged,
		report.LessonsDeprecated, report.WorldsDeprecated, report.NodesDeprecated)
}

func envOrDefault(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
