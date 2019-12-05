package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"
)

type Content struct {
	Filename      string `json:"filename"`
	Fileextension string `json:"fileextension"`
	Header        string `json:"header"`
}

type ArrayofWord struct {
	Key   string
	Value int
}

type ArrayofChar struct {
	Key   string
	Value int
}

var sortWord []ArrayofWord

var sortChar []ArrayofChar

func main() {

	fmt.Println("Starting ...")

	http.HandleFunc("/upload", getData)

	if err := http.ListenAndServe(":9090", nil); err != nil {
		log.Fatal(err)
	}
}

func getData(w http.ResponseWriter, r *http.Request) {

	filePath := getFilePath(w, r)

	text := parseFile(filePath)

	words := parseTextToWords(text)

	wordsCounter := findnumberofwords(words)

	charsCounter := findnumberofchars(text)

	fmt.Printf("Find words : %d \n", len(wordsCounter))
	fmt.Print("Top 10 words : \n")

	for k, v := range wordsCounter {
		sortWord = append(sortWord, ArrayofWord{k, v})
	}

	sort.Slice(sortWord, func(i, j int) bool {
		return sortWord[i].Value > sortWord[j].Value
	})

	i := 0

	for _, kv := range sortWord {

		i++

		if i < 11 {
			fmt.Printf("%d. Word : %s, Amount : %d\n", i, kv.Key, kv.Value)
		}

	}

	fmt.Println("----")

	fmt.Printf("Find chars : %d \n", len(text))
	fmt.Print("Top 10 chars : \n")

	for k, v := range charsCounter {
		sortChar = append(sortChar, ArrayofChar{k, v})
	}

	sort.Slice(sortChar, func(i, j int) bool {
		return sortChar[i].Value > sortChar[j].Value
	})

	j := 0

	for _, kv := range sortChar {

		j++

		if j < 11 {

			fmt.Printf("%d. Char : %s, Amount : %d\n", j, kv.Key, kv.Value)
		}

	}

}

func getFilePath(w http.ResponseWriter, r *http.Request) string {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		// fmt.Fprintf(w, "invalid_http_method")
		return "invalid_http_method"
	}

	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "Parseform() err: %v", err)
	}

	file, handler, err := r.FormFile("upload")
	if err != nil {
		fmt.Println("form file err: ", err)
		// return ("form file err: ", err)
	}

	defer file.Close()

	fmt.Fprintf(w, "%v", handler.Header)

	filePath := "./files/" + handler.Filename

	f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println("open file err: ", err)
		// return "open file err: ", err
	}

	defer f.Close()

	io.Copy(f, file)

	return filePath
}

func parseFile(filepath string) string {

	file, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	b, err := ioutil.ReadAll(file)

	if err != nil {
		fmt.Print(err)
	}

	return string(b)
}

func parseTextToWords(text string) []string {

	words := strings.Fields(text)

	return words
}

func findnumberofwords(list []string) map[string]int {

	duplicateFrequency := make(map[string]int)

	for _, item := range list {
		_, exist := duplicateFrequency[item]

		if exist {
			duplicateFrequency[item]++
		} else {
			duplicateFrequency[item] = 1
		}
	}

	return duplicateFrequency

}

func findnumberofchars(text string) map[string]int {

	duplicateFrequency := make(map[string]int)

	for _, char := range text {
		_, exist := duplicateFrequency[string(char)]

		if exist {
			duplicateFrequency[string(char)]++
		} else {
			duplicateFrequency[string(char)] = 1
		}
	}

	return duplicateFrequency

}
