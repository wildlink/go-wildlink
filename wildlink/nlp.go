package wildlink

import (
	"sync"

	"github.com/dghubble/sling"
	"golang.org/x/xerrors"
)

type NLPService struct {
	c     *Client
	sling *sling.Sling
}

func newNLPService(c *Client, sling *sling.Sling) *NLPService {
	return &NLPService{c: c, sling: sling.Path("v2/nlp/")}
}

type AnalyzeParams struct {
	Content string
}

type AnalyzeResult struct {
	ContentPart []string `json:"ContentPart,omitempty"`
	URL         string   `json:"URL,omitempty"`
}

type AnalyzeResultIterator interface {
	Next() bool
	Scan() (*AnalyzeResult, error)
	Err() error
	Close() error
}

type analyzeResultIter struct {
	sync.Mutex
	results *[]AnalyzeResult
	last    bool
	err     error
}

func (c *analyzeResultIter) Close() error {
	c.Lock()
	defer c.Unlock()
	c.results = nil
	c.last = true
	return nil
}

func (c *analyzeResultIter) Err() error {
	return c.err
}

func (c *analyzeResultIter) Scan() (*AnalyzeResult, error) {
	c.Lock()
	defer c.Unlock()

	if c.results != nil && len(*c.results) > 0 {
		item, buffer := (*c.results)[0], (*c.results)[1:]
		c.results = &buffer
		return &item, nil
	}
	return nil, xerrors.New("Nothing left and Next should have stopped you")

}

func (c *analyzeResultIter) Next() bool {
	c.Lock()
	defer c.Unlock()
	if c.results != nil && len(*c.results) > 0 {
		return true
	}
	if c.last {
		return false
	}

	return len(*c.results) > 0
}

func (s *NLPService) ListAnalysis(params *AnalyzeParams) (AnalyzeResultIterator, error) {

	apiError := new(APIError)
	slingReq := s.c.SetAuthHeaders(s.sling.New().Path("analyze"))

	results := new([]AnalyzeResult)

	_, err := slingReq.Post("").BodyJSON(params).Receive(results, apiError)
	if relevantError(err, *apiError) != nil {
		return nil, relevantError(err, *apiError)
	}

	return &analyzeResultIter{results: results}, nil
}
