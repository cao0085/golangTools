# AccountingTools

會計小工具集

## 環境需求

- Go 版本: `go1.25.3 windows/amd64`
- 主要套件: `github.com/xuri/excelize/v2`

### 使用 Docker 自動執行輸出

```bash
# 建立 Docker Image
docker build -t golang_accountingtools .

# 編譯 EXE 檔案指令 (預設輸出 businessTaxMerger）
docker run --rm -v ${PWD}:/app golang_accountingtools

# 指定輸出
docker run --rm -v ${PWD}:/app -e APP_NAME=testcase golang_accountingtools

# 編譯全部 /apps
docker run --rm -v ${PWD}:/app -e APP_NAME=all golang_accountingtools
```

### 進入容器內部執行

```bash
# 進入容器內部+掛載目錄
docker run --rm -it -v ${PWD}:/app -w /app golang:1.25.3-alpine sh

# 編譯 BusinessTaxMerger
GOOS=windows GOARCH=amd64 go build -o BusinessTaxMerger.exe ./apps/businessTaxMerger/main.go

# 編譯 TestCase
GOOS=windows GOARCH=amd64 go build -o TestCase.exe ./apps/testcase/main.go
```

### 本機執行

```bash
go run ./apps/businessTaxMerger/main.go
go build -o BusinessTaxMerger.exe ./apps/businessTaxMerger/main.go
```

### 專案結構

``` md
AccountingTools_Golang/
├── apps/
│   ├── businessTaxMerger/
│   │   └── main.go              # 營業稅批次處理工具
│   └── testcase/
│       └── main.go              # 測試用最小實現
├── Dockerfile                   # Docker 建置設定
├── build.sh                     # 編譯腳本
├── go.mod                       # Go module 定義檔
├── go.sum                       # 依賴版本鎖定檔
└── README.md                    # 本檔案
```

### 注意事項

1. Windows 編譯需確保終端機支援 UTF-8 編碼
2. 拖曳檔案到終端機時，路徑會自動加上引號，程式已處理此情況
3. 編譯的 exe 檔案包含所有依賴，可直接在其他 Windows 機器上執行
4. 如需減少執行檔大小，可使用 `-ldflags="-s -w"` 參數編譯
