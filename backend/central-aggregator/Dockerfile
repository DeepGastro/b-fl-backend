# 1단계: Go 빌드
FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o main .

# 2단계: 실행 환경 (Python + Torch 포함)
FROM python:3.9-slim
WORKDIR /app

# Go 실행 파일 복사
COPY --from=builder /app/main .
COPY aggregate.py .
COPY requirements.txt .

# 🛠️ PyTorch CPU 버전 및 필수 라이브러리 설치
# --no-cache-dir 옵션으로 도커 이미지 용량을 최소화합니다.
RUN pip install --no-cache-dir -r requirements.txt

# 가중치 저장용 폴더
RUN mkdir -p storage

EXPOSE 8080
CMD ["./main"]