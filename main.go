package main

import (
	"flag"
	"fmt"
	"log"
	"path/filepath"
	"strings"
	"sync"

	"github.com/alireza-bonab/elm-processor/eml"
	"github.com/alireza-bonab/elm-processor/file"
	"github.com/alireza-bonab/elm-processor/project"
)

func validateFileContent(body string, keywords []string) bool {
	// if body contains any of the keywords return true else return false
	for _, keyword := range keywords {
		if strings.Contains(body, keyword) {
			return true
		}
	}
	return false
}

func processFile(fileChan <-chan string, project project.Project) <-chan string {
	out := make(chan string)
	go func() {
		defer close(out)
		for sourceFilePath := range fileChan {
			destDir := filepath.Join(project.DestDir, project.Name)
			// parse file
			mail, err := eml.ParseFile(sourceFilePath)
			if err != nil {
				log.Fatal(err)
				continue
			}
			// validate file content
			if !validateFileContent(mail.Body, project.Keywords) {
				continue
			}
			// get unique file name
			destFilePath, err := file.GetUniqueFileName(sourceFilePath, destDir)
			if err != nil {
				log.Fatal(err)
				continue
			}

			// copy file
			err = file.CopyFile(sourceFilePath, destFilePath)
			if err != nil {
				log.Fatal(err)
				continue
			}

			out <- fmt.Sprintf("%s - %s", project.Name, destFilePath)
		}
	}()

	return out
}

func processProject(project project.Project) []<-chan string {
	filesChan := file.WalkDirChan(project.SourceDir, ".eml")

	projectDestPath := filepath.Join(project.DestDir, project.Name)

	// recreate dest dir
	err := file.RecreateDir(projectDestPath)
	if err != nil {
		log.Fatal(err)
	}

	// make an array of processors
	var processors []<-chan string

	for i := 0; i < 5; i++ {
		processors = append(processors, processFile(filesChan, project))
	}

	return processors
}

func mergeProcessors(processors []<-chan string) <-chan string {
	out := make(chan string)
	var wg sync.WaitGroup
	wg.Add(len(processors))

	for _, processor := range processors {
		go func(p <-chan string) {
			for file := range p {
				out <- file
			}
			wg.Done()
		}(processor)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

func main() {
	log.SetFlags(0)
	configPath := flag.String("config", "", "")
	flag.Parse()

	if *configPath == "" {
		log.Fatal("config file path is required")
	}

	projects, err := project.GetAllProjects(*configPath)
	if err != nil {
		log.Fatal(err)
	}

	var projectProcessors []<-chan string

	for _, project := range projects {
		projectProcessors = append(projectProcessors, processProject(project)...)
	}

	mergedProcessors := mergeProcessors(projectProcessors)

	for file := range mergedProcessors {
		log.Println(file)
	}

	log.Println("done")

}
