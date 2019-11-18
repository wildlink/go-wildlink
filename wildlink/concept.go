package wildlink

import (
	"sync"

	"github.com/dghubble/sling"
	"golang.org/x/xerrors"
)

type ConceptService struct {
	c     *Client
	sling *sling.Sling
}

func newConceptService(c *Client, sling *sling.Sling) *ConceptService {
	return &ConceptService{c: c, sling: sling.Path("v2/concept")}
}

type ConceptListParams struct {
	Query  string `url:"q,omitempty"`
	Limit  uint64 `url:"limit,omitempty"`
	Cursor string `url:"cursor,omitempty"`
}

type ConceptResults struct {
	Concepts   []*Concept
	PrevCursor string
	NextCursor string
}

type Concept struct {
	ID    string `json:"ID,omitempty"`
	Value string `json:"Value,omitempty"`
	URL   string `json:"URL,omitempty"`
}
type ConceptIterator interface {
	Next() bool
	Scan() (*Concept, error)
	Err() error
	Close() error
}

type conceptIter struct {
	sync.Mutex
	params  *ConceptListParams
	results *ConceptResults
	sling   *sling.Sling
	client  *Client
	last    bool
	err     error
}

func (c *conceptIter) Close() error {
	c.Lock()
	defer c.Unlock()
	c.params = nil
	c.results = nil
	c.last = true
	return nil
}
func (c *conceptIter) Err() error {
	return c.err
}

func (c *conceptIter) Scan() (*Concept, error) {
	c.Lock()
	defer c.Unlock()

	if c.results != nil && len(c.results.Concepts) > 0 {
		item, buffer := c.results.Concepts[0], c.results.Concepts[1:]
		c.results.Concepts = buffer
		return item, nil
	}
	return nil, xerrors.New("Nothing left and Next should have stopped you")

}

func (c *conceptIter) Next() bool {
	c.Lock()
	defer c.Unlock()
	if c.results != nil && len(c.results.Concepts) > 0 {
		return true
	}
	if c.last {
		return false
	}

	results := new(ConceptResults)
	apiError := new(APIError)
	slingReq := c.client.SetAuthHeaders(c.sling.New())
	resp, err := slingReq.Get("").QueryStruct(c.params).Receive(results, apiError)
	if relevantError(err, *apiError) != nil {
		c.err = err
		return false
	}
	c.client.updateTokens(resp)
	c.results = results

	c.params.Cursor = c.results.NextCursor
	c.last = c.params.Cursor == ""
	return len(c.results.Concepts) > 0
}

func (s *ConceptService) List(params *ConceptListParams) (ConceptIterator, error) {
	return &conceptIter{params: params, client: s.c, sling: s.sling}, nil
}
