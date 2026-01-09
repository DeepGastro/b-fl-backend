package service

// ai 프로그램을 실행시키는 프로그램입니다!
import (
	"os"
	"os/exec"
)

// 여기서 python3 인지 python 인지 운영체제에 따라 잘 결정해야할듯?
func TriggerAggregation(targetDir string) error {
	cmd := exec.Command("python", "aggregate.py", targetDir)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
