package main

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

func main() {
	// 1. 설정 정보
	gcpIP := "34.158.207.77"
	url := fmt.Sprintf("http://%s:8080/submit", gcpIP)

	filePath := "model.pth" // 보낼 가중치 파일
	hospitalID := "HOSPITAL_SEOUL_01"
	roundID := "1"
	version := "v1.0"

	// 2. 파일 열기
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("파일이 없습니다. 테스트용 model.pth를 만들어주세요.")
		return
	}
	defer file.Close()

	// 3. Multipart 바디 생성
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// 필드 추가 (병원ID, 라운드, 버전)
	_ = writer.WriteField("hospital_id", hospitalID)
	_ = writer.WriteField("round_id", roundID)
	_ = writer.WriteField("model_version", version)

	// 파일 데이터 추가
	part, _ := writer.CreateFormFile("file", filePath)
	io.Copy(part, file)
	writer.Close()

	// 4. 전송 요청
	fmt.Printf("🚀 %s 서버로 가중치 전송 시작...\n", gcpIP)
	resp, err := http.Post(url, writer.FormDataContentType(), body)
	if err != nil {
		fmt.Printf("❌ 전송 실패: %v\n", err)
		return
	}
	defer resp.Body.Close()

	// 5. 결과 확인
	respBody, _ := io.ReadAll(resp.Body)
	fmt.Printf("✅ 서버 응답 [%s]: %s\n", resp.Status, string(respBody))
}
