package ingest

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/vilmibm/trunkless/db"
)

// TODO
// - [X] finalize gutenberg ingestion
// - [ ] clean up commands
// 	- [X] get down to just ingest/cutup/serve
//  - [ ] add arguments for generalizing
// - [X] clean up repo
// - [ ] push and deploy to town with new pg db
// - [ ] gamefaqs extraction
// - [ ] corpus selector
// - [ ] deploy to town
// - [ ] geocities
// - [ ] blog post
// - [ ] launch

type IngestOpts struct {
	Conn     *pgx.Conn
	Corpus   string
	CutupDir string
}

func Ingest(o IngestOpts) error {
	conn := o.Conn
	cutupDir := o.CutupDir

	dir, err := os.Open(o.CutupDir)
	if err != nil {
		return fmt.Errorf("could not open %s: %w", cutupDir, err)
	}
	defer dir.Close()

	entries, err := dir.Readdirnames(-1)
	if err != nil {
		return fmt.Errorf("could not read %s: %w", cutupDir, err)
	}

	idx, err := os.Open(path.Join(cutupDir, "_title_index.tsv"))
	if err != nil {
		return fmt.Errorf("failed to open source index: %w", err)
	}
	defer idx.Close()

	corpusid := db.StrToID(o.Corpus)
	tablename := fmt.Sprintf("phrases_%s", corpusid)
	_, err = conn.Exec(context.Background(),
		fmt.Sprintf(`CREATE TABLE %s (
			id SERIAL PRIMARY KEY,
			sourceid char(7) NOT NULL,
			phrase TEXT,
			FOREIGN KEY (sourceid) REFERENCES sources(id)
		)`, tablename))
	if err != nil {
		return fmt.Errorf("could not create table '%s': %w", tablename, err)
	}

	_, err = conn.Exec(context.Background(),
		"INSERT INTO corpora (id, name) VALUES ($1, $2) ON CONFLICT DO NOTHING",
		corpusid, o.Corpus)
	if err != nil {
		return fmt.Errorf("failed to create '%s' corpus: %w", o.Corpus, err)
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
		sql := fmt.Sprintf("COPY %s(sourceid, phrase) FROM '%s'", tablename, p)
		_, err = conn.Exec(context.Background(), sql)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to ingest '%s': %s\n", p, err.Error())
		}
	}
	return nil
}
