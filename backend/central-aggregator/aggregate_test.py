import sys
import os

if __name__ == "__main__":
    target_dir = sys.argv[1] # Go에서 넘겨준 uploadPath
    print(f"--- 파이썬 합산기 작동 시작 ---")
    print(f"대상 폴더: {target_dir}")
    files = os.listdir(target_dir)
    print(f"발견된 파일들: {files}")
    print(f"--- 합산 완료 (가짜) ---")