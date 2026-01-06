package main

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

func main() {
	// 1. 전송할 파일 열기
	filePath := "test.txt"
	file, _ := os.Open(filePath)
	defer file.Close()

	// 2. 데이터를 담을 버퍼 만들기
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// 3. 파일 데이터를 버퍼에 추가
	part, _ := writer.CreateFormFile("file", filePath)
	io.Copy(part, file)

	// 4. 나머지 데이터들을 body에 추가
	writer.WriteField("hospital_id", "HOSP_01")
	writer.WriteField("round_id", "1")
	writer.WriteField("model_version", "v1.0")
	writer.Close()

	// 5. 중앙 서버로 전송
	req, _ := http.NewRequest("POST", "http://localhost:8080/submit", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, _ := client.Do(req)
	defer resp.Body.Close()

	// 6. 결과 확인
	io.Copy(os.Stdout, resp.Body)
}
