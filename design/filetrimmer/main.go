package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/spf13/viper"
)

var (
	baseFolder string
	origin     string
	target     string
)

func initV() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.SetConfigType("yaml")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Fatal error config file: %s \n", err)
	}
	baseFolder = viper.GetString("base_folder")
	baseFolder = filepath.Clean(baseFolder)
	origin = viper.GetString("origin")
	target = viper.GetString("target")
	if baseFolder == "" {
		log.Fatal("base_folder is required")
	}
	log.Printf("baseFolder: %s\n", baseFolder)
	log.Printf("origin: %s\n", origin)
	log.Printf("target: %s\n", target)
}

func main() {
	initV()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		trimFiles(baseFolder, &wg)
	}()

	wg.Wait()
}

func trimFiles(currFolder string, wg *sync.WaitGroup) {
	// find files
	log.Printf("find files in %s\n", currFolder)

	if _, err := os.Stat(currFolder); os.IsNotExist(err) {
		log.Fatalf("Folder %s does not exist", currFolder)
	}

	// read files
	files, err := os.ReadDir(currFolder)
	if err != nil {
		log.Fatalf("Cannot read folder %s", currFolder)
	}

	for _, file := range files {
		if file.IsDir() {
			wg.Add(1)
			go trimFiles(currFolder+"/"+file.Name(), wg)
		} else if strings.Contains(file.Name(), origin) {
			oldFilePath := currFolder + "/" + file.Name()
			newName := strings.Replace(file.Name(), origin, target, 1)
			newFilePath := currFolder + "/" + newName
			err := os.Rename(oldFilePath, newFilePath)
			if err != nil {
				log.Fatalf("Cannot rename file %s to %s", oldFilePath, newFilePath)
			}
			log.Printf("Renamed %s to %s\n", oldFilePath, newFilePath)
		}
	}

	wg.Done()
}
