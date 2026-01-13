package main

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"time"
)

func main() {
	// 1. 설정 정보
	gcpIP := "34.158.207.77"
	url := fmt.Sprintf("http://%s:8080/submit", gcpIP)

	filePath := "model.pth"
	hospitalID := "HOSPITAL_SEOUL_01" // 테스트할 때 01, 02, 03으로 바꿔서 해보세요!
	roundID := "1"
	version := "v1.0"

	// --- [STEP 1] 로컬 가중치 파일 업로드 ---
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("파일이 없습니다. 테스트용 model.pth를 만들어주세요.")
		return
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	_ = writer.WriteField("hospital_id", hospitalID)
	_ = writer.WriteField("round_id", roundID)
	_ = writer.WriteField("model_version", version)
	part, _ := writer.CreateFormFile("file", filePath)
	io.Copy(part, file)
	writer.Close()
	file.Close() // 업로드 후 파일 닫기

	fmt.Printf("🚀 %s 서버로 가중치 전송 시작...\n", gcpIP)
	resp, err := http.Post(url, writer.FormDataContentType(), body)
	if err != nil {
		fmt.Printf("❌ 전송 실패: %v\n", err)
		return
	}

	respBody, _ := io.ReadAll(resp.Body)
	fmt.Printf("✅ 서버 응답 [%s]: %s\n", resp.Status, string(respBody))
	resp.Body.Close()

	// --- [STEP 2] 합산 결과 기다리기 (Polling) ---
	// 업로드가 성공했으니 이제 서버가 합산을 끝낼 때까지 기다립니다.
	downloadGlobalModel(gcpIP, roundID)
}

func downloadGlobalModel(gcpIP string, roundID string) {
	url := fmt.Sprintf("http://%s:8080/download/%s", gcpIP, roundID)
	fmt.Printf("\n⏳ 라운드 %s 합산 결과 대기 중...\n", roundID)

	for {
		resp, err := http.Get(url)
		if err != nil {
			fmt.Printf("❌ 서버 연결 실패: %v\n", err)
			time.Sleep(5 * time.Second)
			continue
		}

		if resp.StatusCode == http.StatusOK {
			out, _ := os.Create("global_model_received.pth")
			io.Copy(out, resp.Body)
			out.Close()
			resp.Body.Close()

			fmt.Println("\n✅ 글로벌 모델 다운로드 완료! 다음 학습을 시작합니다.")
			break
		}

		// 아직 파일이 없으면(404 등) 5초 대기 후 재시도
		resp.Body.Close()
		fmt.Print(".") // 대기 중임을 표시
		time.Sleep(5 * time.Second)
	}
}
