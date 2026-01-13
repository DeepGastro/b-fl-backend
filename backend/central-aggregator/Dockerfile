# --- 1단계: Go 빌드 ---
FROM golang:1.21-bullseye AS builder
WORKDIR /app

# 캐시 효율을 위해 모듈 파일 먼저 복사
COPY go.mod go.sum ./
RUN go mod download

# 전체 소스 복사 (handlers, service, utils 등 포함)
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# --- 2단계: 실행 환경 (Python 3.9 + PyTorch) ---
FROM python:3.9-slim
WORKDIR /app

# 시스템 라이브러리 설치 (Go 바이너리 실행 호환성 확보)
RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

# 빌드된 Go 실행 파일 가져오기
COPY --from=builder /app/main .

# 파이썬 코드 및 모델 데이터 전체 복사
# (run_server.py, server.py, models/ 폴더 등이 모두 들어갑니다)
COPY . .

# AI 의존성 설치 (requirements.txt에 torch 필수!)
RUN pip install --no-cache-dir -r requirements.txt

# 권한 설정 및 저장소 생성
RUN mkdir -p storage && chmod 777 storage
RUN chmod +x ./main

EXPOSE 8080

# 백엔드 서버 실행
CMD ["./main"]