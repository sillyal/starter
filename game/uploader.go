package game

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"time"
)

func Upload(s string) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, error := writer.CreateFormFile("file", time.Now().Format("20060102150405")+".txt")
	if error != nil {
		panic(fmt.Sprintln("Fail to connect to upload logs to server", error))
	}

	w := bufio.NewWriter(part)
	count, error := w.WriteString(s)
	if error != nil {
		panic(fmt.Sprintln("Fail to write to upload logs to server", error))
	}
	fmt.Printf("Wrote %d bytes\n", count)
	w.Flush()

	err := writer.Close()
	if err != nil {
		panic(fmt.Sprintln("Fail to close multipart writer", err))
	}

	domain := "https://sillyal.com"
	req, err := http.NewRequest("POST", domain+"/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	} else {
		body := &bytes.Buffer{}
		_, err := body.ReadFrom(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		resp.Body.Close()

		fmt.Println(resp.StatusCode)
		fmt.Println(resp.Header)
		fmt.Println(resp.Header["Location"][0])
		fmt.Println(body)
	}
}
