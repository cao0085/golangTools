#!/bin/sh

# 編譯腳本：根據 APP_NAME 環境變數選擇要編譯的程式

echo "========================================="
echo "開始編譯 Go 專案..."
echo "========================================="

case "$APP_NAME" in
  "businessTaxMerger")
    echo "📦 編譯目標: BusinessTaxMerger"
    GOOS=windows GOARCH=amd64 go build -o /app/BusinessTaxMerger.exe /app/apps/businessTaxMerger/main.go
    if [ $? -eq 0 ]; then
      echo "✅ 編譯成功: BusinessTaxMerger.exe"
      ls -lh /app/BusinessTaxMerger.exe
    else
      echo "❌ 編譯失敗"
      exit 1
    fi
    ;;

  "testcase")
    echo "📦 編譯目標: TestCase"
    GOOS=windows GOARCH=amd64 go build -o /app/TestCase.exe /app/apps/testcase/main.go
    if [ $? -eq 0 ]; then
      echo "✅ 編譯成功: TestCase.exe"
      ls -lh /app/TestCase.exe
    else
      echo "❌ 編譯失敗"
      exit 1
    fi
    ;;

  "all")
    echo "📦 編譯目標: 全部"
    echo ""
    echo "正在編譯 BusinessTaxMerger..."
    GOOS=windows GOARCH=amd64 go build -o /app/BusinessTaxMerger.exe /app/apps/businessTaxMerger/main.go

    echo "正在編譯 TestCase..."
    GOOS=windows GOARCH=amd64 go build -o /app/TestCase.exe /app/apps/testcase/main.go

    echo ""
    echo "✅ 全部編譯完成:"
    ls -lh /app/*.exe
    ;;

  *)
    echo "❌ 錯誤: 未知的 APP_NAME='$APP_NAME'"
    echo ""
    echo "可用選項:"
    echo "  - businessTaxMerger  (預設)"
    echo "  - testcase"
    echo "  - all               (編譯全部)"
    echo ""
    echo "使用方式:"
    echo "  docker run --rm -v \${PWD}:/app -e APP_NAME=testcase golang_accountingtools"
    exit 1
    ;;
esac

echo "========================================="
echo "編譯流程結束"
echo "========================================="
