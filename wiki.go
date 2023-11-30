package docs

func newWiki(token string, client *Client) *Wiki {
	return &Wiki{
		baseClient: client,
		tokenIns: tokenIns{
			token: token,
		},
	}
}

// Docx represent a doc file
type Wiki struct {
	Err error
	tokenIns
	baseClient *Client
}

func (w *Wiki) GetMeta() (*MetaRespMetas, error) {
	return w.baseClient.GetMeta(w.token, FileTypeWiki, "")
}

func (w *Wiki) Statistics() (*FileStatistics, error) {
	return newFile(w.baseClient).statistics(w.token, FileTypeWiki)
}
