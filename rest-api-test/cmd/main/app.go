package main

import (
	"bufio"
	"bytes"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
)

func GetTranslation(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

	/* 	file, err := os.Open("./data/file.htm")
	   	if err != nil {
	   		log.Print("wrong thmthing")
	   	}
	   	defer file.Close()

	   	wr := bytes.Buffer{}
	   	sc := bufio.NewScanner(file)
	   	for sc.Scan() {
	   		wr.WriteString(sc.Text())
	   	} */

	origin := params.ByName("origin")
	res := FindWord(origin, GetWords("file.htm"))

	for i := 0; i < len(res); i++ {
		w.Write([]byte(res[i] + "/n/n/n"))
	}
}

func main() {
	router := httprouter.New()
	router.GET("/api/get-translation/:origin", GetTranslation)

	start(router)
}

func SplitTextToArray(text string) []string {
	bodyText := strings.Split(text, "body")[1]
	splitText := strings.Split(bodyText, "<p>")

	slice := []string{"start"}
	index := 0

	for i := 0; i < len(splitText); i++ {
		if strings.Contains(splitText[i], "style=\"font-weight:bold;\"") {
			index++
			slice = append(slice, "<p>"+splitText[i])
			continue
		}
		slice[index] += "<p>" + splitText[i]
	}
	return slice
}

func GetSplitTranslate(str string) []string {

	return strings.Split(str, "</p>")
}

func GetWords(path string) string {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	wr := bytes.Buffer{}
	sc := bufio.NewScanner(file)
	for sc.Scan() {
		wr.WriteString(sc.Text())
	}
	return wr.String()
}

func FindWord(word string, str string) []string {
	slice := SplitTextToArray(str)
	res := make([]string, 0)
	for i := 0; i < len(slice); i++ {
		split := strings.Split(slice[i], "</p>") //GetSplitTranslate(slice[i])
		if strings.Contains(split[0], word) {
			res = append(res, slice[i]) //strings.ReplaceAll(slice[i], `'`, ""))
		}
	}

	return res
}

func start(router *httprouter.Router) {
	log.Println("starst application")

	listener, err := net.Listen("tcp", ":1234")
	if err != nil {
		panic(err)
	}
	server := &http.Server{
		Handler:      router,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Printf("server is listening port %s\n", "0.0.0.0:1234")
	log.Fatalln(server.Serve(listener))
}
