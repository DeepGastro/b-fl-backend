package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	// 1. gin engine instance 생성
	r := gin.Default()

	// 2. 가중치와 블록체인에 보낼 내용 받고, 형식 만들기
	r.POST("/submit", func(c *gin.Context) {
		// A. form 데이터에서 병원 id, 라운드, 버전 읽어오기
		hospitalID := c.PostForm("hospital_id")
		roundID := c.PostForm("round_id")
		modelVersion := c.PostForm("model_version")

		// B. 파일 읽기
		file, header, err := c.Request.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"errer": "파일을 읽을 수 없습니다"})
			return
		}
		defer file.Close()

		// C. 해시 계산하기
		hash := sha256.New()
		if _, err := io.Copy(hash, file); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "해시 계산 오류"})
			return
		}
		hashString := hex.EncodeToString(hash.Sum(nil))

		// D. 결과 출력
		fmt.Printf("--- 신규 가중치 수신 ---\n")
		fmt.Printf("병원 ID: %s\n", hospitalID)
		fmt.Printf("라운드: %s\n", roundID)
		fmt.Printf("모델 버전: %s\n", modelVersion)
		fmt.Printf("파일명: %s\n", header.Filename)
		fmt.Printf("해시값: %s\n", hashString)
		fmt.Printf("------------------------\n")

		// E. 가중치 수신을 정상적으로 받았는지 응답 반환(병원에게)
		c.JSON(http.StatusOK, gin.H{
			"message":     "success",
			"weight_hash": hashString,
			"hospital_id": hospitalID,
		})
	})

	// 3. 서버 실행
	r.Run(":8080")
}
