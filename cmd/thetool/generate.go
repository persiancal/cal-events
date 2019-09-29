package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/pmezard/go-difflib/difflib"
)

var (
	compare *bool
	dist    *string
)

func compareFiles(src, dst []byte) error {
	if bytes.Compare(src, dst) == 0 {
		return nil
	}

	diff := difflib.UnifiedDiff{
		A:        difflib.SplitLines(string(src)),
		B:        difflib.SplitLines(string(dst)),
		FromFile: "new",
		ToFile:   "old",
	}
	text, err := difflib.GetUnifiedDiffString(diff)
	_, _ = fmt.Fprint(os.Stderr, text)
	return err
}

func generate(cmd *command, fl *File) error {
	j, err := json.MarshalIndent(fl, "", "  ")
	if err != nil {
		return fmt.Errorf("converting to json failed: %w", err)
	}

	path := filepath.Join(*dist, strings.ToLower(fl.Name)+".json")
	if *compare {
		c, err := ioutil.ReadFile(path)
		if err != nil {
			return fmt.Errorf("the target file is not exist: %w", err)
		}
		return compareFiles(c, j)
	}

	return ioutil.WriteFile(path, j, 0600)
}

func init() {
	cur, _ := os.Getwd()
	cmd := newCommand("generate", "generate the dist folder", generate)
	dist = cmd.Flags.String("dist", cur, "The dist folder")
	compare = cmd.Flags.Bool("compare", false, "Compare the current dist with the base file")
	registerCommand(cmd)
}
