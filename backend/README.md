# BFL System - Blockchain-based Federated Learning Backend

> 연합학습 기반 의료 데이터 가중치 수집 시스템

## 📋 프로젝트 개요

BFL(Blockchain-based Federated Learning) System은 병원 노드들이 로컬에서 학습한 AI 모델 가중치를 중앙 서버로 전송하고, 중앙 서버가 이를 수집·검증·합산하는 분산 연합학습 시스템입니다.

각 병원은 개인정보 보호를 위해 원본 데이터를 외부로 전송하지 않고, 로컬에서 학습된 모델 가중치만을 중앙 서버로 전송합니다. 중앙 서버는 여러 병원으로부터 수신한 가중치를 집계하여 글로벌 모델을 생성합니다.

### 주요 특징

- 🏥 **분산형 학습**: 병원별 로컬 학습 후 가중치만 공유
- 🔐 **무결성 검증**: SHA-256 해시 기반 파일 검증
- 🤖 **자동 집계**: 3개 이상의 가중치 수신 시 자동 합산
- 🔗 **블록체인 연동 준비**: 향후 투명성 및 감사 추적 지원 예정

## 🏗️ 시스템 아키텍처

```
┌─────────────────┐         ┌─────────────────┐         ┌─────────────────┐
│  Hospital Node  │         │  Hospital Node  │         │  Hospital Node  │
│    (HOSP_01)    │         │    (HOSP_02)    │         │    (HOSP_03)    │
└────────┬────────┘         └────────┬────────┘         └────────┬────────┘
         │                           │                           │
         │  POST /submit             │  POST /submit             │
         │  (model weights)          │  (model weights)          │
         └───────────────────────────┼───────────────────────────┘
                                     │
                                     ▼
                          ┌──────────────────────┐
                          │  Central Aggregator  │
                          │                      │
                          │  - File Reception    │
                          │  - Hash Validation   │
                          │  - Storage           │
                          └──────────┬───────────┘
                                     │
                                     │ (3+ files collected)
                                     ▼
                          ┌──────────────────────┐
                          │   Python Script      │
                          │  (aggregate.py)      │
                          │                      │
                          │  AI Weight           │
                          │  Aggregation         │
                          └──────────────────────┘
```

## 📁 프로젝트 구조

```
backend/
├── central-aggregator/          # 중앙 집계 서버
│   ├── main.go                  # 서버 엔트리포인트
│   ├── service/
│   │   └── ai_worker.go        # Python 스크립트 실행 로직
│   ├── utils/
│   │   └── hash.go             # SHA-256 해시 계산
│   ├── storage/                # 수신된 가중치 파일 저장소
│   ├── database/               # DB 관련 코드 (예정)
│   ├── go.mod
│   └── go.sum
│
├── hospital-node/               # 병원 노드
│   ├── main.go                  # 가중치 전송 클라이언트
│   ├── data/                    # 로컬 학습 데이터
│   ├── test.txt                 # 테스트용 샘플 파일
│   └── go.mod
│
└── docker-compose.yml           # Docker 배포 설정
```

## 🛠️ 기술 스택

| Category | Technology |
|----------|-----------|
| **Language** | Go 1.21+ |
| **Web Framework** | Gin |
| **AI Processing** | Python 3.x |
| **Hash Algorithm** | SHA-256 |
| **Deployment** | Docker, Docker Compose |
| **Future Plan** | PostgreSQL, Blockchain Integration |

## 🚀 시작하기

### 사전 요구사항

- Go 1.21 이상
- Python 3.x (AI 집계 스크립트 실행용)
- Docker & Docker Compose (선택사항)

### 설치 및 실행

#### 1. Central Aggregator (중앙 서버) 실행

```bash
cd backend/central-aggregator
go mod tidy
go run main.go
```

서버는 기본적으로 `http://localhost:8080`에서 실행됩니다.

#### 2. Hospital Node (병원 노드) 실행

```bash
cd backend/hospital-node
go mod tidy
go run main.go
```

병원 노드는 `test.txt` 파일을 중앙 서버로 전송합니다.

## 📡 API 명세

### POST /submit

병원 노드가 학습된 모델 가중치를 중앙 서버로 전송하는 엔드포인트입니다.

**Request**

- **Method**: `POST`
- **URL**: `http://localhost:8080/submit`
- **Content-Type**: `multipart/form-data`

**Parameters**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `file` | File | Yes | 모델 가중치 파일 (.bin, .pt, .h5 등) |
| `hospital_id` | String | Yes | 병원 고유 식별자 (예: `HOSP_01`) |
| `round_id` | String | Yes | 학습 라운드 번호 (예: `1`, `2`) |
| `model_version` | String | Yes | 모델 버전 (예: `v1.0`) |

**Response (Success)**

```json
{
  "message": "success",
  "weight_hash": "a3b5c7d9e1f2...",
  "hospital_id": "HOSP_01"
}
```

**Response (Error)**

```json
{
  "error": "파일을 읽을 수 없습니다"
}
```

### 서버 동작 흐름

1. **파일 수신**: 병원으로부터 가중치 파일 및 메타데이터 수신
2. **저장**: `storage/round_{round_id}/{hospital_id}_{model_version}_{filename}` 형식으로 저장
3. **해시 계산**: SHA-256 알고리즘으로 파일 무결성 검증용 해시값 생성
4. **집계 트리거**: 해당 라운드에 3개 이상의 파일이 모이면 자동으로 Python 집계 스크립트 실행
5. **응답 반환**: 병원에게 해시값과 함께 성공 응답 전송

## 🧪 테스트

### Postman을 사용한 테스트

1. Postman에서 새 요청 생성
2. Method를 `POST`로 설정
3. URL: `http://localhost:8080/submit`
4. Body 탭 → `form-data` 선택
5. 다음 필드 추가:
   - `file`: (File 타입) 테스트할 파일 선택
   - `hospital_id`: `HOSP_01`
   - `round_id`: `1`
   - `model_version`: `v1.0`
6. Send 버튼 클릭

### cURL을 사용한 테스트

```bash
curl -X POST http://localhost:8080/submit \
  -F "file=@test.txt" \
  -F "hospital_id=HOSP_01" \
  -F "round_id=1" \
  -F "model_version=v1.0"
```

## 🔐 보안 및 무결성

### SHA-256 해시 검증

모든 수신된 파일은 SHA-256 해시값이 자동으로 계산되어 저장됩니다. 이는 다음과 같은 목적으로 활용됩니다:

- 파일 전송 과정에서의 무결성 검증
- 블록체인 저장을 위한 식별자
- 감사 추적 및 투명성 확보

해시 계산은 `utils/hash.go`의 `CalculateFileHash` 함수에서 수행됩니다.

## 🐳 Docker 배포

```bash
# Docker Compose로 전체 시스템 실행
cd backend
docker-compose up -d

# 로그 확인
docker-compose logs -f

# 종료
docker-compose down
```

## 🗺️ 향후 계획

### Phase 1 (현재)
- ✅ 기본 가중치 수신/저장 기능
- ✅ SHA-256 해시 검증
- ✅ Python 집계 스크립트 실행

### Phase 2 (진행 예정)
- 🔄 PostgreSQL 데이터베이스 연동
  - 메타데이터 저장 (hospital_id, round_id, hash, timestamp 등)
  - 학습 히스토리 관리
- 🔄 블록체인 연동
  - 가중치 해시값 블록체인 기록
  - 투명성 및 감사 추적 강화
- 🔄 인증 및 권한 관리
  - 병원별 API Key 발급
  - JWT 기반 인증

### Phase 3 (계획)
- 📊 대시보드 개발 (모니터링 UI)
- 🔔 알림 시스템 (가중치 수신, 집계 완료 등)
- 📈 성능 최적화 및 스케일링

## 📝 코드 주요 부분 설명

### Central Aggregator 핵심 로직

**파일 수신 및 저장** (`main.go:20-46`)
```go
// POST /submit 핸들러에서 multipart/form-data 처리
hospitalID := c.PostForm("hospital_id")
roundID := c.PostForm("round_id")
modelVersion := c.PostForm("model_version")
fileHeader, err := c.FormFile("file")

// storage/round_{roundID} 디렉토리에 저장
uploadPath := filepath.Join("storage", "round_"+roundID)
dst := filepath.Join(uploadPath, hospitalID+"_"+modelVersion+"_"+fileHeader.Filename)
```

**해시 계산** (`main.go:48-49`, `utils/hash.go:12-24`)
```go
hashString, _ := utils.CalculateFileHash(dst)
// SHA-256 해시값을 hex string으로 반환
```

**AI 집계 트리거** (`main.go:52-60`, `service/ai_worker.go:9-14`)
```go
files, _ := os.ReadDir(uploadPath)
if len(files) >= 3 {
    go func() {
        service.TriggerAggregation(uploadPath)
    }()
}
// Python3 aggregate.py {targetDir} 실행
```

### Hospital Node 핵심 로직

**파일 전송** (`hospital-node/main.go:11-41`)
```go
// multipart/form-data 생성
writer := multipart.NewWriter(body)
writer.CreateFormFile("file", filePath)
writer.WriteField("hospital_id", "HOSP_01")
writer.WriteField("round_id", "1")
writer.WriteField("model_version", "v1.0")

// POST 요청 전송
http.NewRequest("POST", "http://localhost:8080/submit", body)
```

## 🤝 기여하기

이슈 제기 및 Pull Request는 언제나 환영합니다!

## 📄 라이센스

This project is licensed under the MIT License.

## 👥 개발팀

- **Project**: BFL System (Blockchain-based Federated Learning)
- **Repository**: https://github.com/DeepGastro/b-fl-backend

---

**Note**: 이 프로젝트는 의료 데이터의 프라이버시를 보호하면서도 협업 학습을 가능하게 하는 연합학습 시스템의 백엔드 구현입니다.
