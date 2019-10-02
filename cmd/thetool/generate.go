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
	"gopkg.in/yaml.v2"
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
	text, _ := difflib.GetUnifiedDiffString(diff)
	_, _ = fmt.Fprint(os.Stderr, text)
	return fmt.Errorf("the dist folder is not up to date, please run make generate command and commit the result")
}

func compareAndWrite(fl string, data []byte) error {
	if *compare {
		c, err := ioutil.ReadFile(fl)
		if err != nil {
			return fmt.Errorf("the target file is not exist: %w", err)
		}
		return compareFiles(c, data)
	}

	return ioutil.WriteFile(fl, data, 0600)

}

func generate(cmd *command, fl *File) error {
	path := filepath.Join(*dist, strings.ToLower(fl.Name))

	j, err := json.MarshalIndent(fl, "", "  ")
	if err != nil {
		return fmt.Errorf("converting to json failed: %w", err)
	}

	if err := compareAndWrite(path+".json", j); err != nil {
		return err
	}

	y, err := yaml.Marshal(fl)
	if err != nil {
		return fmt.Errorf("converting to yaml failed: %w", err)
	}

	if err := compareAndWrite(path+".yml", y); err != nil {
		return err
	}

	return nil
}

func init() {
	cur, _ := os.Getwd()
	cmd := newCommand("generate", "generate the dist folder", generate)
	dist = cmd.Flags.String("dist", cur, "The dist folder")
	compare = cmd.Flags.Bool("compare", false, "Compare the current dist with the base file")
	registerCommand(cmd)
}
