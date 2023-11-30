package docs

func newDocx(token string, client *Client) *Docx {
	return &Docx{
		baseClient: client,
		tokenIns: tokenIns{
			token: token,
		},
	}
}

// Docx represent a doc file
type Docx struct {
	Err error
	tokenIns
	baseClient *Client
}

func (d *Docx) GetMeta() (*MetaRespMetas, error) {
	return d.baseClient.GetMeta(d.token, FileTypeWiki, "")
}

func (d *Docx) Statistics() (*FileStatistics, error) {
	return newFile(d.baseClient).statistics(d.token, FileTypeDocx)
}
