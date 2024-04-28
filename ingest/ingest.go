package ingest

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/vilmibm/trunkless/db"
)

const cutupDir = "/home/vilmibm/pg_plaintext/cutup"

// TODO
// - [X] finalize gutenberg ingestion
// - [ ] clean up commands
// - [ ] clean up repo
// - [ ] push and deploy to town with new pg db
// - [ ] gamefaqs extraction
// - [ ] corpus selector
// - [ ] deploy to town
// - [ ] geocities
// - [ ] blog post
// - [ ] launch

func IngestGut() error {
	conn, err := db.Connect()
	if err != nil {
		return err
	}
	defer conn.Close(context.Background())

	dir, err := os.Open(cutupDir)
	if err != nil {
		return fmt.Errorf("could not open %s: %w", cutupDir, err)
	}

	// echo gutenberg | sha1sum | head -c7
	corpusid := "cb20c3e"
	_, err = conn.Exec(context.Background(), "INSERT INTO corpora (id, name) VALUES ($1, $2) ON CONFLICT DO NOTHING", corpusid, "gutenberg")
	if err != nil {
		return fmt.Errorf("failed to create gutenberg corpus: %w", err)
	}

	entries, err := dir.Readdirnames(-1)
	if err != nil {
		return fmt.Errorf("could not read %s: %w", cutupDir, err)
	}

	idx, err := os.Open(path.Join(cutupDir, "_title_index.tsv"))
	if err != nil {
		return fmt.Errorf("failed to open source index: %w", err)
	}

	tx, err := conn.Begin(context.Background())
	if err != nil {
		return fmt.Errorf("could not open transaction: %w", err)
	}

	s := bufio.NewScanner(idx)
	for s.Scan() {
		line := s.Text()
		parts := strings.SplitN(line, "	", 2)
		if len(parts) != 2 {
			return fmt.Errorf("malformed line in sourceMap: %s", line)
		}
		_, err = tx.Exec(context.Background(),
			"INSERT INTO sources (id, corpusid, name) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING",
			parts[0], corpusid, parts[1])
	}

	tx.Commit(context.Background())

	for _, e := range entries {
		if strings.HasPrefix(e, "_") {
			continue
		}
		p := path.Join(cutupDir, e)
		sql := fmt.Sprintf("COPY phrases(sourceid, phrase) FROM '%s'", p)
		_, err = conn.Exec(context.Background(), sql)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to ingest '%s': %s\n", p, err.Error())
		}
	}

	return nil
}
