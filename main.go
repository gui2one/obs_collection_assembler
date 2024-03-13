package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func extractFilePaths(data interface{}) []string {
	var paths []string

	switch v := data.(type) {
	case string:
		// Use a regular expression to identify potential file paths
		// Adjust the regular expression based on your specific needs
		// re := regexp.MustCompile(`(?:[a-zA-Z]:\\|/)?[\w\d./\\-]+`)
		re := regexp.MustCompile(`(?:[a-zA-Z]:)?(?:/[\w\d.-]+)+`)
		matches := re.FindAllString(v, -1)
		paths = append(paths, matches...)
	case map[string]interface{}:
		for _, value := range v {
			paths = append(paths, extractFilePaths(value)...)
		}
	case []interface{}:
		for _, item := range v {
			paths = append(paths, extractFilePaths(item)...)
		}
	}

	return paths
}

func readJSON(json_path string) map[string]interface{} {
	var data, err = os.ReadFile(json_path)
	if err != nil {
		log.Fatal(err)
	} else {
		var json_data map[string]interface{}
		err := json.Unmarshal(data, &json_data)
		if err != nil {
			log.Fatal(err)
		}
		return json_data
	}

	return nil
}

func checkFileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return !errors.Is(err, os.ErrNotExist)
}
func copy(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()

	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}

func main() {
	fmt.Println(len(os.Args), os.Args)

	if len(os.Args) < 2 {
		log.Fatal("No JSON file path provided, I need an argument !")
	}

	// var json_path = "C:/gui2one/OBS_scene_collections/TEST_SCENE.json"
	var json_path = os.Args[1]
	root_dir := filepath.Dir(json_path)
	collection_name := strings.Split(filepath.Base(json_path), ".")[0]
	fmt.Println("root dir is     : " + root_dir)
	fmt.Println("COllection name : " + collection_name)

	collection_dir := root_dir + "/" + collection_name
	os.Mkdir(collection_dir, os.ModeAppend)
	var dat = readJSON(json_path)

	filePaths := extractFilePaths(dat)

	fmt.Println("Found File Paths:")
	for _, path := range filePaths {
		fmt.Println(path)
		file_exists := checkFileExists(path)

		if file_exists {

			fmt.Println("Copying file " + path)
			copy(path, collection_dir+"/"+filepath.Base(path))
		}
	}

}
