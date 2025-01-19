package main

import (
	"flag"
	"io/fs"
	"jsonbytes"
	"log"
	"os"
	"path/filepath"
)

var dirFlag = flag.String("dir", ".", "the absolute path of the directory to search for JSON files")
var logNotJson = flag.Bool("lognotjson", false, "enables logging that files aren't json")
var skipIsJson = flag.Bool("skipisjson", false, "disables logging that files are json")

func main() {
	flag.Parse()
	foundInvalidJson := false
	err := filepath.WalkDir(*dirFlag, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			log.Printf("😭 Could not walk %s, err: %s\n", path, err.Error())
		} else if d == nil || d.IsDir() {
			return nil
		}
		fileContent, err := os.ReadFile(path)
		if err != nil {
			log.Printf("😭 Could not open %s, err: %s\n", path, err.Error())
			return nil
		}
		err = jsonbytes.IsJson(fileContent)
		if err != nil && *logNotJson {
			foundInvalidJson = true
			log.Printf("❌ %s is not JSON.\n", path)
		} else if err == nil && !*skipIsJson {
			log.Printf("✅ %s is JSON!\n", path)
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	if foundInvalidJson {
		log.Println("Some of the files checked weren't valid JSON! 😭")
	} else {
		log.Println("Every file checked was valid JSON! 🥳")
	}
}
