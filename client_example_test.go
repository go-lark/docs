package docs

import (
	"fmt"
)

func Example_newClient() {
	client := NewClient("", "")
	folder := client.RootFolder()
	fmt.Println("root folder: ", folder.GetToken())
}

func Example_createSheet() {
	spreadSheets := NewClient("", "").RootFolder().CreateSpreadSheet("sheet title")
	spreadSheets.Share(PermEdit, false, NewMemberWithEmail("aa.com"))
	err := spreadSheets.GetSheetByIndex(1).WriteRows(
		[]string{"name", "age"},
		[][]interface{}{
			{"Ace", 10},
			{"Bob", 11},
		},
	).Err
	if err != nil {
		fmt.Println("err: ", err.Error())
	}
}
