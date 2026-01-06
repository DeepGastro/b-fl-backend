package service

// ai 프로그램을 실행시키는 프로그램입니다!
import (
	"os"
	"os/exec"
)

func TriggerAggregation(targetDir string) error {
	cmd := exec.Command("python3", "aggregate.py", targetDir)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
