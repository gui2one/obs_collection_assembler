package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"example.com/obs_collection_assembler/utils"
	"github.com/fatih/color"
	"github.com/sqweek/dialog"
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

func filter_paths(paths []string) []string {
	var good_paths []string
	for _, path := range paths {
		if checkFileExists(path) {
			good_paths = append(good_paths, path)
		}
	}
	return good_paths
}

func replacePaths(base_str string, paths []string, root_path string) string {
	var final_str = base_str
	for _, path := range paths {
		base_name := filepath.Base(path)
		final_str = strings.Replace(final_str, path, root_path+"/"+base_name, -1)
	}

	return final_str
}

func main() {
	filename, err := dialog.File().Filter("JSON file", "json").Load()

	if err != nil {
		log.Println("Aborting ...")
		return
	}

	// var json_path = "C:/gui2one/OBS_scene_collections/TEST_SCENE.json"
	var json_path = filename
	root_dir := filepath.Dir(json_path)
	collection_name := strings.Split(filepath.Base(json_path), ".")[0]
	fmt.Println("root dir is     : " + root_dir)
	fmt.Println("Collection name : " + collection_name)

	collection_dir := root_dir + "/" + collection_name
	os.Mkdir(collection_dir, os.ModeAppend)

	var dat = utils.ReadJsonToStringMap(json_path)

	file_paths := extractFilePaths(dat)

	good_paths := filter_paths(file_paths)

	fmt.Printf("Good paths are : \n %v\n\n", good_paths)
	fmt.Println("Found File Paths:")
	for _, path := range good_paths {
		fmt.Println(path)
		file_exists := checkFileExists(path)

		if file_exists {

			fmt.Println("Copying file " + path)
			copy(path, collection_dir+"/"+filepath.Base(path))
		}
	}

	data, err := utils.ReadJsonToString(json_path)
	if err != nil {
		log.Printf("Problem reading json file ->  %s", json_path)
	}

	str_with_replaced_paths := replacePaths(data, good_paths, collection_name)

	converted_json_path := filepath.Join(root_dir, collection_name+"_converted.json")
	os.WriteFile(converted_json_path, []byte(str_with_replaced_paths), os.ModeAppend)
	color.Green("Wrote %s to disk\n", converted_json_path)

}
