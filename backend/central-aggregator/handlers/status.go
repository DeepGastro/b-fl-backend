package handlers

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

// GetStatus는 특정 라운드의 파일 제출 현황을 반환합니다.
func GetStatus(c *gin.Context) {
	roundID := c.Query("round_id") // URL 파라미터에서 round_id 추출
	if roundID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "round_id가 필요합니다."})
		return
	}

	uploadPath := filepath.Join("storage", "round_"+roundID)

	// 폴더 읽기
	files, err := os.ReadDir(uploadPath)
	if err != nil {
		// 폴더가 없으면 아직 아무도 제출 안 한 것
		c.JSON(http.StatusOK, gin.H{
			"round_id": roundID,
			"files":    []string{},
			"count":    0,
		})
		return
	}

	var fileList []string
	for _, file := range files {
		if !file.IsDir() {
			fileList = append(fileList, file.Name())
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"round_id": roundID,
		"files":    fileList,
		"count":    len(fileList),
	})
}
