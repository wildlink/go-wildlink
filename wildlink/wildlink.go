package wildlink

import (
	"net/http"
	"os"
	"sync"

	"github.com/dghubble/sling"
)

const version = "v1.0.0"

var baseURI = "https://api.wfi.re/"

type Client struct {
	sync.Mutex
	sling          *sling.Sling
	appID          uint64
	appKey         string
	device         *Device
	DeviceService  *DeviceService
	ConceptService *ConceptService
}

func init() {
	uri, ok := os.LookupEnv("WILDLINK_URI")
	if ok {
		baseURI = uri
	}
}

func NewClient(httpClient *http.Client) *Client {
	base := sling.New().Client(httpClient).Base(baseURI).Set("User-Agent", "go-wildlink/"+version)

	c := &Client{}
	c.sling = base
	c.DeviceService = newDeviceService(c, base.New())
	c.ConceptService = newConceptService(c, base.New())
	return c
}

func (c *Client) updateTokens(resp *http.Response) {
	c.Lock()
	defer c.Unlock()
	if c.device != nil {
		newToken := resp.Header.Get("X-WF-DeviceToken")
		if newToken != "" {
			c.device.Token = newToken
		}
	}
}

func (c *Client) Connect() error {
	return c.DeviceService.ensure()
}

func (c *Client) SetAppID(id uint64) *Client {
	c.Lock()
	defer c.Unlock()
	c.appID = id
	return c
}

func (c *Client) SetAppKey(key string) *Client {
	c.Lock()
	defer c.Unlock()
	c.appKey = key
	return c
}

func (c *Client) SetDevice(device *Device) *Client {
	c.Lock()
	defer c.Unlock()
	c.device = device
	return c
}

func (c *Client) Device() *Device {
	c.Lock()
	defer c.Unlock()
	return c.device
}
