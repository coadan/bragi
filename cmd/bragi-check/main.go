package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/coadan/bragi"
)

func main() {
	profilePath := flag.String("profile", "", "path to a Bragi profile JSON file")
	jsonOutput := flag.Bool("json", false, "write canonical events as JSON Lines")
	flag.Parse()
	if *profilePath == "" || flag.NArg() == 0 {
		fmt.Fprintln(os.Stderr, "usage: bragi-check -profile profile.json [-json] stream.bragi [...]")
		os.Exit(2)
	}

	profileFile, err := os.Open(*profilePath)
	if err != nil {
		fatal(err)
	}
	profile, err := bragi.LoadProfile(profileFile)
	closeErr := profileFile.Close()
	if err != nil {
		fatal(err)
	}
	if closeErr != nil {
		fatal(closeErr)
	}

	failed := false
	for _, path := range flag.Args() {
		if err := check(path, profile, *jsonOutput); err != nil {
			fmt.Fprintf(os.Stderr, "%s: %v\n", path, err)
			failed = true
		}
	}
	if failed {
		os.Exit(1)
	}
}

func check(path string, profile bragi.Profile, jsonOutput bool) error {
	stream, err := os.Open(path)
	if err != nil {
		return err
	}
	defer stream.Close()

	materializer, err := bragi.NewMaterializer(profile)
	if err != nil {
		return err
	}
	decoder := bragi.NewDecoder(bragi.DecoderOptions{MaxLineBytes: profile.Limits.MaxLineBytes, AllowCRLF: true})
	buffer := make([]byte, 4096)
	rejected := false
	for {
		count, readErr := stream.Read(buffer)
		if count > 0 {
			records, issues := decoder.Write(buffer[:count])
			for _, issue := range issues {
				event := materializer.RejectSource(issue)
				rejected = true
				writeEvent(event, jsonOutput)
			}
			for _, record := range records {
				for _, event := range materializer.Apply(record) {
					if event.Kind == "op.rejected" || event.Kind == "commit.rejected" {
						rejected = true
					}
					writeEvent(event, jsonOutput)
				}
			}
		}
		if readErr == io.EOF {
			break
		}
		if readErr != nil {
			return readErr
		}
	}
	for _, issue := range decoder.Finish() {
		event := materializer.RejectSource(issue)
		rejected = true
		writeEvent(event, jsonOutput)
	}
	completeIssues := materializer.ValidateComplete()
	if len(completeIssues) > 0 {
		rejected = true
		if !jsonOutput {
			for _, issue := range completeIssues {
				fmt.Fprintf(os.Stderr, "%s: %s: %s\n", path, issue.Code, issue.Message)
			}
		}
	}
	if rejected {
		return fmt.Errorf("not conformant")
	}
	if !jsonOutput {
		fmt.Printf("%s: ok (%d canonical events)\n", path, len(materializer.Events()))
	}
	return nil
}

func writeEvent(event bragi.Event, enabled bool) {
	if !enabled {
		return
	}
	if err := json.NewEncoder(os.Stdout).Encode(event); err != nil {
		fatal(err)
	}
}

func fatal(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
