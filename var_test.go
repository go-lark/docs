package docs

import "os"

var (
	baseDomain = ""

	testAPPID      = ""
	testAPPSecret  = ""
	testAPPID2     = ""
	testAPPSecret2 = ""

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
	baseDomain = os.Getenv("DOCS_BASEDOMAIN")
	// for test bot with website token
	testAPPID = os.Getenv("DOCS_APPID")
	testAPPSecret = os.Getenv("DOCS_APPSECRET")

	// for test bot with backend token
	testAPPID2 = os.Getenv("DOCS_APPID2")
	testAPPSecret2 = os.Getenv("DOCS_APPSECRET2")

	// docs file token
	testSpreadSheetToken = os.Getenv("DOCS_SPREADSHEET_TOKEN")
	testDocToken = os.Getenv("DOCS_DOC_TOKEN")
	testFolderToken = os.Getenv("DOCS_FOLDER_TOKEN")
	testBitableToken = os.Getenv("DOCS_BITABLE_TOKEN")
	testDocxToken = os.Getenv("DOCS_DOCX_TOKEN")

	testBigFile = os.Getenv("DOCS_BIG_FILE")
	testUserEmail = os.Getenv("DOCS_USEREMAIL")

}
