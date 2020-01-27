package controllers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"
)

type ViewData struct {
	Title       string
	Description string
}

type ArrayofWord struct {
	Key   string
	Value int
}

type ArrayofChar struct {
	Key   string
	Value int
}

type Content struct {
	WordsNumber int
	CharsNumber int
	Words       []ArrayofWord
	Char        []ArrayofChar
}

var sortWord []ArrayofWord

var sortChar []ArrayofChar

func GetuploadForm(w http.ResponseWriter, r *http.Request) {

	data := ViewData{
		Title:       "Uploading file",
		Description: "Put your file",
	}
	t, _ := template.ParseFiles("templates/upload.html")

	t.Execute(w, data)

}

func GetStatistic(w http.ResponseWriter, r *http.Request) {
	GetData(w, r)
}

func GetData(w http.ResponseWriter, r *http.Request) []byte {

	filePath := GetFilePath(w, r)

	text, err := ParseFile(filePath)

	if err != nil {
		log.Println(err)
	}

	words := ParseTextToWords(text)

	wordsCounter := FindNumberOfWords(words)

	charsCounter := FindNumberOfChars(text)

	//fmt.Printf("Find words : %d \n", len(wordsCounter))
	//fmt.Print("Top 10 words : \n")

	for k, v := range wordsCounter {
		sortWord = append(sortWord, ArrayofWord{k, v})
	}

	sort.Slice(sortWord, func(i, j int) bool {
		return sortWord[i].Value > sortWord[j].Value
	})

	//i := 0
	//
	//for _, kv := range sortWord {
	//
	//	i++
	//
	//	if i < 11 {
	//		fmt.Printf("%d. Word : %s, Amount : %d\n", i, kv.Key, kv.Value)
	//	}
	//
	//}
	//
	//fmt.Println("----")
	//
	//fmt.Printf("Find chars : %d \n", len(text))
	//fmt.Print("Top 10 chars : \n")

	for k, v := range charsCounter {
		sortChar = append(sortChar, ArrayofChar{k, v})
	}

	sort.Slice(sortChar, func(i, j int) bool {
		return sortChar[i].Value > sortChar[j].Value
	})

	//j := 0
	//
	//for _, kv := range sortChar {
	//
	//	j++
	//
	//	if j < 11 {
	//
	//		fmt.Printf("%d. Char : %s, Amount : %d\n", j, kv.Key, kv.Value)
	//	}
	//
	//}

	data := Content{
		WordsNumber: len(wordsCounter),
		CharsNumber: len(text),
		Words:       sortWord[:10],
		Char:        sortChar[:10],
	}

	//var jsonData []byte

	jsonData, err := json.Marshal(data)

	if err != nil {
		fmt.Println(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)

	return jsonData
}

func GetFilePath(w http.ResponseWriter, r *http.Request) string {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return "invalid_http_method"
	}

	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "Parseform() err: %v", err)
	}

	file, handler, err := r.FormFile("upload")
	if err != nil {
		fmt.Println("form file err: ", err)
	}

	defer file.Close()

	filePath := "./files/" + handler.Filename

	f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println("open file err: ", err)
	}

	defer f.Close()

	io.Copy(f, file)

	return filePath
}

func ParseFile(filepath string) (string, error) {

	file, err := os.Open(filepath)
	if err != nil {
		return "", err
	}

	defer file.Close()

	read, err := ioutil.ReadAll(file)

	if err != nil {
		return "", err
	}

	return string(read), nil
}

func ParseTextToWords(text string) []string {

	words := strings.Fields(text)

	return words
}

func FindNumberOfWords(list []string) map[string]int {

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

func FindNumberOfChars(text string) map[string]int {

	duplicateFrequency := make(map[string]int)

	parseString := strings.Replace(text, " ", "", -1)

	for _, char := range parseString {
		_, exist := duplicateFrequency[string(char)]

		if exist {
			duplicateFrequency[string(char)]++
		} else {
			duplicateFrequency[string(char)] = 1
		}
	}

	return duplicateFrequency

}
