package core

// TaxRecord 營業稅 - 營業人進銷項資料檔
type TaxRecord struct {
	// ===== 基本資訊 =====

	// FormatCode 格式代號 X(002) - 位置 1-2
	FormatCode string

	// DeclarantTaxId 申報營業人統編號 X(009) - 位置 3-11
	DeclarantTaxId string

	// SequenceNumber 流水號 X(007) - 位置 12-18
	SequenceNumber string

	// ===== 資料所屬年月 =====

	// DataYear 資料所屬年度 9(003) - 位置 19-21
	DataYear string

	// DataMonth 資料所屬月份 9(002) - 位置 22-23
	DataMonth string

	// ===== 買受人/營業/銷售人統編號 =====

	// BuyerTaxId 買受人統一編號 X(008) - 位置 24-31 (欄位共用)
	BuyerTaxId string

	// BusinessNumber 發票訖號 9(008) - 位置 24-31 (欄位共用)
	BusinessNumber string

	// SellerTaxId 銷售人統一編號 X(008) - 位置 32-39
	SellerTaxId string

	// ===== 統一發票 =====

	// InvoicePrefix 發票字軌 X(002) - 位置 40-41
	InvoicePrefix string

	// InvoiceStartNumber 發票(起)號碼 9(008) - 位置 42-49
	InvoiceStartNumber string

	// ===== 總額張數 =====

	// TotalSheets 彙總張數 9(004) - 位置 32-35
	TotalSheets string

	// Blank1 空白 X(004) - 位置 36-39 (欄位共用)
	Blank1 string

	// OtherVoucherNumber 其他憑證號碼 X(010) - 位置 40-49
	OtherVoucherNumber string

	// UtilitySequenceNumber 公用事業裁其流水號 X(010) - 位置 40-49
	UtilitySequenceNumber string

	// Blank2 空白 X(004) - 位置 32-35
	Blank2 string

	// CustomsTaxPaymentNumber 海關代徵營業稅繳納證號碼 X(014) - 位置 36-49
	CustomsTaxPaymentNumber string

	// ===== 金額 =====

	// SalesAmount 銷售金額 9(012) - 位置 50-61 (欄位共用)
	SalesAmount string

	// TaxBase 營業稅稅基 9(012) - 位置 50-61 (欄位共用)
	TaxBase string

	// ===== 課稅別 =====

	// TaxType 課稅別 X(001) - 位置 62-62
	TaxType string

	// TaxAmount 營業稅額 9(010) - 位置 63-72
	TaxAmount string

	// ===== 扣抵代號 =====

	// DeductionCode 扣抵代號 X(001) - 位置 73-73
	DeductionCode string

	// Blank3 空白 X(005) - 位置 74-78
	Blank3 string

	// SpecialTaxRate 特種稅額類稅率 X(001) - 位置 79-79
	SpecialTaxRate string

	// AggregationMark 彙加性記 X(001) - 位置 80-80 (欄位共用)
	AggregationMark string

	// AllocationMark 分攤性記 X(001) - 位置 80-80 (欄位共用)
	AllocationMark string

	// CustomsClearanceMark 通關方式註記 X(001) - 位置 81-81
	CustomsClearanceMark string

	// ===== 原始資料 =====

	// RawData 原始行資料
	RawData string

	// LineNumber 行號（在原始檔案中的位置）
	LineNumber int

	// SourceFileName 來源檔案名稱
	SourceFileName string
}

// TxtFileInfo TXT 檔案資訊
type TxtFileInfo struct {
	FilePath  string
	FileName  string
	LineCount int
}

// NewTxtFileInfo 建立 TxtFileInfo
func NewTxtFileInfo(filePath string, lineCount int) *TxtFileInfo {
	return &TxtFileInfo{
		FilePath:  filePath,
		FileName:  getFileName(filePath),
		LineCount: lineCount,
	}
}

// getFileName 從完整路徑取得檔案名稱
func getFileName(filePath string) string {
	for i := len(filePath) - 1; i >= 0; i-- {
		if filePath[i] == '\\' || filePath[i] == '/' {
			return filePath[i+1:]
		}
	}
	return filePath
}
