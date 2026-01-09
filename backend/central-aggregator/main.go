package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"central-aggregator/handlers"
	"central-aggregator/service"
	"central-aggregator/utils"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// gin engine instance 생성
	r := gin.Default()

	// --- CORS 설정 시작 ---
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"}, // 프론트엔드 주소 허용
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	// --- CORS 설정 끝 ---

	// 가중치와 블록체인에 보낼 내용 받고, 형식 만들기
	r.POST("/submit", func(c *gin.Context) {
		// form 데이터에서 병원 id, 라운드, 버전 읽어오기
		hospitalID := c.PostForm("hospital_id")
		roundID := c.PostForm("round_id")
		modelVersion := c.PostForm("model_version")

		// 파일 받기
		fileHeader, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "파일을 읽을 수 없습니다"})
			return
		}

		// 받은 파일 로컬 디스크에 저장할 폴더 생성(폴더가 없을 경우 생성하도록)
		uploadPath := filepath.Join("storage", "round_"+roundID)
		if _, err := os.Stat(uploadPath); os.IsNotExist(err) {
			os.MkdirAll(uploadPath, 0755)
		}

		// 저장될 파일명 결정(ex: storage/round_1/HOPS_01_v1.0_20210922123456.bin)
		dst := filepath.Join(uploadPath, hospitalID+"_"+modelVersion+"_"+fileHeader.Filename)

		// 서버 하드디스크에 파일 저장
		if err := c.SaveUploadedFile(fileHeader, dst); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "파일 저장 오류"})
			return
		}

		// 해시값 생성
		hashString, _ := utils.CalculateFileHash(dst)

		// TODO: 이걸 몇번에 한번씩 가중치 합산? 그리고 ai파트에게 파일 경로를 주면 가중치 파일이 한개가 생성돼서 나오는건가?
		// ai 프로그램 실행(3번 호출당 한번씩 가중치 합산)
		files, _ := os.ReadDir(uploadPath)
		if len(files) >= 3 {
			go func() {
				if err := service.TriggerAggregation(uploadPath); err != nil {
					fmt.Println("AI 가중치 합산 오류", err)
				}
			}()
		}

		// 결과 출력
		fmt.Printf("--- 가중치 수신 및 저장 성공 ---\n")
		fmt.Printf("병원 ID: %s\n", hospitalID)
		fmt.Printf("라운드: %s\n", roundID)
		fmt.Printf("모델 버전: %s\n", modelVersion)
		fmt.Printf("파일명: %s\n", fileHeader.Filename)
		fmt.Printf("해시값: %s\n", hashString)
		fmt.Printf("------------------------\n")

		// 가중치 수신을 정상적으로 받았는지 응답 반환(병원에게)
		c.JSON(http.StatusOK, gin.H{
			"message":     "success",
			"weight_hash": hashString,
			"hospital_id": hospitalID,
		})
	})

	r.GET("/status", handlers.GetStatus)

	// 3. 서버 실행
	r.Run(":8080")
}
