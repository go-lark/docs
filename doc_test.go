package docs

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	"testing"
	"time"

	"github.com/go-lark/docs/doctypes"
	"github.com/hilaily/kit/pp"
	"github.com/stretchr/testify/assert"
)

func TestGetContent(t *testing.T) {
	b, c, _, err := docClient().GetContent()
	assert.NoError(t, err)
	resp := string(b)
	assert.NotEmpty(t, resp)
	fmt.Println(resp)
	ioutil.WriteFile("/tmp/doc.json", b, 0777)
	t.Log("c: ", c)
}

func TestGetMeta(t *testing.T) {
	res, err := docClient().GetMeta()
	assert.NoError(t, err)
	assert.Equal(t, baseDomain+"/docs/doccnHA8wfHJtNsVBhH6MkRoh7m", res.URL)
	assert.Equal(t, DocDeleteFlagNormal, res.DeleteFlag)
}

func TestDocDeleteFlag(t *testing.T) {
	type a struct {
		A DocDeleteFlag
		B *DocDeleteFlag
	}

	aa := a{
		DocDeleteFlagDeleted,
		&DocDeleteFlagNormal,
	}
	en, _ := json.Marshal(aa)
	assert.Equal(t, "{\"A\":2,\"B\":0}", string(en))
	res := &a{}
	err := json.Unmarshal(en, res)
	assert.NoError(t, err)
	assert.Equal(t, aa.A, res.A)
	assert.Equal(t, aa.B, res.B)
}

func TestDoc_Create(t *testing.T) {
	b, c, _, err := docClient().GetContent()
	assert.NoError(t, err)
	resp := string(b)
	assert.NotEmpty(t, resp)
	title := c.Title
	title.Elements[0].TextRun.Text = "import test " + time.Now().String()
	doc := getClientNew().OpenFolder(testFolderToken).CreateDoc(title, c.Body)
	assert.Nil(t, doc.Err)
	assert.NotEmpty(t, doc.GetToken())
	//assert.NotEmpty(t, doc.GetURL())
	t.Log("token: ", doc.GetToken())
	doc.ChangeOwner(NewMemberWithEmail(testUserEmail), false, true)
	assert.Nil(t, doc.Err)
}

func TestDoc_AddWholeComment(t *testing.T) {
	_, err := docClient().AddWholeComment("test comment")
	assert.NoError(t, err)
}

func TestQueryEscape(t *testing.T) {
	u := "https%3A%2F%2Fplayer.bilibili.com%2Fplayer.html%3Faid%3D25898700"
	u2 := url.QueryEscape(u)
	t.Log(u2)
}

func TestDoc_CreateDoc(t *testing.T) {
	var blocks []doctypes.IBlocks
	for _, v := range [][]doctypes.IBlocks{
		genParagraph(),
		genHorizonLine(),
		genOtherBlocks(),
	} {
		blocks = append(blocks, v...)
	}

	title := doctypes.NewTitle("test create doc " + time.Now().Format("2006-01-02 15:04:05"))
	body := doctypes.NewBody(blocks...)
	pp.Dump("title: ", title)
	pp.Dump("body: ", body)
	doc := getClientNew().OpenFolder(testFolderToken).CreateDoc(title, body)
	assert.Nil(t, doc.Err)

	doc = doc.ChangeOwner(NewMemberWithEmail(testUserEmail), false, true)
	assert.Nil(t, doc.Err)
}

func genParagraph() []doctypes.IBlocks {
	res := []doctypes.IBlocks{
		doctypes.NewBlockParagraph(doctypes.NewElementTextRun("我是第一行")),
		doctypes.NewBlockParagraph(doctypes.NewElementTextRun("以下都是 blocks")).HeadingLevel(doctypes.ParagraphHeadineLevel1),
		doctypes.NewBlockParagraph(doctypes.NewElementTextRun("我是 paragraph")).HeadingLevel(doctypes.ParagraphHeadineLevel2),
		doctypes.NewBlockParagraph(
			doctypes.NewElementTextRun("以下都是 elements，放到同一行"),
			doctypes.NewElementPerson("ou_3bbe8a09c20e89cce9bff989ed840674"),
			doctypes.NewElementDocsLink(baseDomain+"/docs/doccnMao97bBZWcWfVG3A6NJBTg"),
			doctypes.NewElementReminder(false, 1600507800, false, doctypes.ReminderNotify1dBefore),
			doctypes.NewElementEquation("E_n = - R_H \\left( {\\frac{1}{{n^2 }}} \\right) = \\frac{{ - 2.178 \\times 10^{ - 18} }}{{n^2 }}joule"),
		),
		doctypes.NewBlockParagraphWithTextRun("我是textRun").HeadingLevel(doctypes.ParagraphHeadineLevel2),
		doctypes.NewBlocksList(doctypes.ListNumber, []*doctypes.IndentList{
			{Text: "有序列表第一个", Ident: 1},
			{Text: "有序列表第一个缩进", Ident: 2},
			{Text: "有序列表第二个", Ident: 1},
		}),
		doctypes.NewBlocksList(doctypes.ListBullet, []*doctypes.IndentList{
			{Text: "无序列表第一个", Ident: 1},
			{Text: "无序列表第二个", Ident: 1},
			{Text: "无序列表第二个缩进", Ident: 2},
		}),
		doctypes.NewBlockParagraph(doctypes.NewElementTextRun("居中对齐")).Align(doctypes.ParagraphAlignCenter),
		doctypes.NewBlockParagraph(doctypes.NewElementTextRun("右对齐")).Align(doctypes.ParagraphAlignRight),
		doctypes.NewBlockParagraph(doctypes.NewElementTextRun("引用")).SetQuote(),
		doctypes.NewBlockParagraph(
			doctypes.NewElementTextRun("加粗").SetBold(),
			doctypes.NewElementTextRun("斜体").SetItalic(),
			doctypes.NewElementTextRun("下划线").SetUnderline(),
			doctypes.NewElementTextRun("删除线").SetStrickThrouth(),
			doctypes.NewElementTextRun("Inline Code").SetCodeInline(),
			doctypes.NewElementTextRun("绿颜色").SetColor(&doctypes.ColorRGBA{Red: 98, Green: 210, Blue: 86, Alpha: float64(1)}, nil),
			doctypes.NewElementTextRun("链接").SetLink("https://github.com"),
		),
		doctypes.NewBlockCode([]string{
			"func main()",
			`  fmt.Println("hello world")`,
			"}",
		}),
		doctypes.NewBlockCallout(
			&doctypes.BlockCallout{
				CalloutEmojiID:         "heart",
				CalloutTextColor:       &doctypes.ColorRGBA{Red: 216, Green: 57, Blue: 49},
				CalloutBackgroundColor: &doctypes.ColorRGBA{Red: 254, Green: 255, Blue: 250},
				CalloutBorderColor:     &doctypes.ColorRGBA{Red: 205, Green: 178, Blue: 250},
				Body: doctypes.BlockCalloutBody{
					Blocks: doctypes.NewBlockParagraphWithTextRun("测试高亮").ToBlocks(),
				},
			},
		),
	}
	return res
}

func genHorizonLine() []doctypes.IBlocks {
	return []doctypes.IBlocks{
		doctypes.NewBlockParagraphWithTextRun("我是 horizonLine").HeadingLevel(1),
		doctypes.NewBlockHorizontalLine(),
	}
}

func genOtherBlocks() []doctypes.IBlocks {
	return []doctypes.IBlocks{
		doctypes.NewBlockParagraphWithTextRun("我是数据表").HeadingLevel(1),
		doctypes.NewBlockBitable("", doctypes.BitableViewGrid),
		doctypes.NewBlockBitable("", doctypes.BitableViewKanban),
		doctypes.NewBlockSheet("", 9, 9),
		doctypes.NewBlockChatGroup("oc_637494e772e7ded6df6a6d950921e65a"),
		doctypes.NewBlockEmbeddedPage(doctypes.EmbedBilibili, "https%3A%2F%2Fplayer.bilibili.com%2Fplayer.html%3Faid%3D25898700"),
	}
}

func docClient() *Doc {
	return getClientNew().OpenDoc(testDocToken)
}
