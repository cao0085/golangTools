# 使用 Linux alpine 映像(更小更快)
FROM golang:1.25.3-alpine

# 設定工作目錄
WORKDIR /app

# 複製 go.mod 和 go.sum
COPY go.mod go.sum ./

# 下載依賴
RUN go mod download

# 複製原始碼
COPY . .

# 設定環境變數（預設編譯 businessTaxMerger）
# 可選值: businessTaxMerger, testcase
ENV APP_NAME=businessTaxMerger

# 編譯腳本
COPY build.sh /build.sh
RUN chmod +x /build.sh

# 執行編譯
CMD ["/build.sh"]
