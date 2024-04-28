package cmd

import (
	"github.com/spf13/cobra"
	"github.com/vilmibm/trunkless/db"
	"github.com/vilmibm/trunkless/ingest"
)

func init() {
	// TODO option for cutupDir

	ingestCmd.Flags().StringP("cutupdir", "d", "", "directory to files produced by cutup cmd")
	ingestCmd.MarkFlagRequired("cutupdir")
	rootCmd.AddCommand(ingestCmd)
}

var ingestCmd = &cobra.Command{
	Use:   "ingest corpusname",
	Short: "ingest already cut-up corpora from disk into database",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cutupDir := cmd.Flags().Lookup("cutupdir").Value.String()
		corpus := args[0]

		conn, err := db.Connect()
		if err != nil {
			return err
		}

		opts := ingest.IngestOpts{
			Conn:     conn,
			CutupDir: cutupDir,
			Corpus:   corpus,
		}

		return ingest.Ingest(opts)
	},
}

// thoughts
//
// having multitenancy in the db makes phrase selection harder. i need to determine the ID offsets for each corpus's phrase list.
// currently waiting on an explain analyze for:
// explain analyze select min(p.id),max(p.id) from phrases p join sources s on s.id = p.sourceid and s.corpusid='cb20c3e';
// planning time 12ms
// exec time 91s
// trying again with inner join which was fast but not noticeably; the explain looks the same (which makes sense--no rows with null allowed are involved).

// if i stick with this i can expect several minutes(!) of startup time to the server; however, since i'm generating ID lookups outside of sql, my lookup should still be O(1).
// some options:
// - change everything so every corpus is in its own table:
//   ${corpus}_phrases: id, sourceid, text
//   corpora: id, name
//   sources: id, corpusid, name
// - cache the result of the min/max id analysis. i could do this to disk or in the db...i would probably do it in the db:
//   id_ranges: corpusid, minid, maxid

// thinking about this more, as i add corpora the phrases table is going to
// grow into the billions (assuming other sources are similar in scale to
// gutenberg). turns out postgresql has table partitioning but idk if that will
// help me since the ID space will be shared.

// having a table per corpus's phrases will also make tearing down corpora easier -- otherwise i have to regen the entire phrases table to remove gaps in ID space.

// so it's settled; I'm going to retool for table-per-corpus.
