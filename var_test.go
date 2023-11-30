package docs

import "os"

var (
	feishuDomain = ""
	tenantDomain = ""

	testAPPID     = ""
	testAPPSecret = ""

	testSpreadSheetToken = ""
	testDocToken         = ""
	testFolderToken      = ""
	testBitableToken     = ""
	testDocxToken        = ""

	testBigFile   = ""
	testUserEmail = ""
)

// some variable for test
func init() {
	feishuDomain = os.Getenv("DOCS_FEISHU_DOMAIN")
	tenantDomain = os.Getenv("DOCS_BASEDOMAIN")

	testAPPID = os.Getenv("DOCS_APPID")
	testAPPSecret = os.Getenv("DOCS_APPSECRET")

	// docs file token
	testSpreadSheetToken = os.Getenv("DOCS_SPREADSHEET_TOKEN")
	testDocToken = os.Getenv("DOCS_DOC_TOKEN")
	testFolderToken = os.Getenv("DOCS_FOLDER_TOKEN")
	testBitableToken = os.Getenv("DOCS_SPREADSHEET_TOKEN")
	testDocxToken = os.Getenv("DOCS_DOCX_TOKEN")

	testBigFile = os.Getenv("DOCS_BIG_FILE")
	testUserEmail = os.Getenv("DOCS_USEREMAIL")

}
