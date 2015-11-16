package cmd

import (
	"log"
	"strings"

	"github.com/codegangsta/cli"
	"golang.org/x/net/context"
)

var cmdRecordsDelete = cli.Command{
	Name:      "delete",
	Usage:     "deletes zone record",
	ArgsUsage: "<zone-id> <record-id> [<record-id> ...]",
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "all",
			Usage: "deletes all zone records",
		},
		cli.StringFlag{
			Name:  "type",
			Usage: "deletes only records of given type (can be comma separated)",
		},
		cli.StringFlag{
			Name:  "ignore",
			Usage: "ignores records of given type (can be comma separated)",
		},
	},
	Action: func(c *cli.Context) {
		if c.Args().First() == "" {
			log.Fatal("Usage error: zone id is required.")
		}

		if !c.Bool("all") {
			if len(c.Args()) < 2 {
				log.Fatal("Usage error: --all flag or at least one record id is required.")
			} else if c.String("type") != "" {
				log.Fatal("Usage error: --type can be only used with --all.")
			} else if c.String("ignore") != "" {
				log.Fatal("Usage error: --type can be only used with --all.")
			}
		}

		var (
			ids    []string
			types  = splitComma(c.String("type"))
			ignore = splitComma(c.String("ignore"))
		)

		if c.Bool("all") {
			records, err := client(c).Records.List(context.Background(), c.Args().First())
			if err != nil {
				log.Fatalf("Error listing records: %v", err)
			}
			for _, record := range records {
				if stringIn(record.Type, ignore) {
					log.Printf("Ignoring record %s (type=%s)", record.ID, record.Type)
					continue
				}
				if len(types) > 0 && !stringIn(record.Type, types) {
					continue
				}
				ids = append(ids, record.ID)
			}
		} else {
			ids = c.Args()[1:]
		}

		for _, id := range ids {
			err := client(c).Records.Delete(context.Background(), c.Args().First(), id)
			if err != nil {
				log.Fatalf("Error deleting %q: %v", id, err)
			}
		}

		log.Println("Done")
	},
}

func splitComma(s string) []string {
	if len(s) == 0 {
		return nil
	}
	return strings.Split(s, ",")
}
