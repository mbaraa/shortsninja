package main

import (
	"fmt"
	"github.com/baraa-almasri/shortsninja/globals"
	"github.com/baraa-almasri/shortsninja/handlers"
	"github.com/baraa-almasri/shortsninja/initialsetup"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

func printAvailableTextFiles() {
	op, err := exec.Command("ls", ".").Output()
	if err != nil {
		panic(err)
	}

	files0 := strings.Split(string(op), "\n")
	var files []string

	for _, v := range files0 {
		if strings.Contains(v, ".txt") {
			files = append(files, v)
		}
	}

	for _, v := range files {
		println(v)
	}
}

// still testing :)
func main() {
	switch os.Args[1] {
	case "setup":
		println("available text files")
		printAvailableTextFiles()
		print("enter text file to be used as a short urls source, \nif you want to generate a new file type -1: ")
		var fileName string
		fmt.Scanf("%s", &fileName)

		var f *os.File
		if fileName == "-1" {
			f, _ = os.Create(fileName)
			var amount, numChars int
			print("enter number of wanted chars in each url: ")
			fmt.Scanf("%d", &numChars)

			print("enter amount of wanted urls: ")
			fmt.Scanf("%d", &amount)

			initialsetup.GenerateShortURLsFile(numChars, amount, f)
		} else {
			f, _ = os.OpenFile(fileName, os.O_WRONLY, 0755)
		}

		initialsetup.GenerateMuxFileUsingURLsFile(f)

		break
	case "start":
		f, _ := os.Open("./urls.txt")
		_, _, _ = globals.LoadGlobals(f)
		mux := handlers.GetMux()
		mux.HandleFunc("/shorten/", handlers.AddURL)
		log.Fatal(http.ListenAndServe(":8080", mux))

		break
	default:
		println("hmm, that's wrong ain't it?")
	}
}
