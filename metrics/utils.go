package metrics

import (
	"bufio"
	"io/ioutil"
	"log"
	"os"
	"regexp"
)

func getFilesInDir(root string) (files []string) {
	filer, err := ioutil.ReadDir(root)
	if err != nil {
		log.Println(err)
		return files
	}

	for _, f := range filer {
		files = append(files, root+"/"+f.Name())
	}
	return files
}

func matchFilesLine(files []string, regx string) map[string]string {
	var lines = make(map[string]string)
	for _, f := range files {
		txty := matchFileLine(f, regx)
		if txty != "" {
			lines[f] = txty
		}
	}
	return lines
}

func matchFileLine(f string, regx string) (lines string) {
	file, err := os.Open(f)
	if err != nil {
		log.Printf("failed opening file: %s\n", err)
		return
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		txty := scanner.Text()
		matched, _ := regexp.MatchString(regx, txty)
		if matched == true {
			lines = txty
			break
		}
	}

	file.Close()
	return
}

func parseStringFloat(value interface{}) float64 {
	switch value.(type) {
	case float64:
		return value.(float64)
	default:
		return 0
	}
}
