package wildlink

import (
	"runtime"

	"github.com/dghubble/sling"
)

type DeviceService struct {
	c     *Client
	sling *sling.Sling
}

func newDeviceService(c *Client, sling *sling.Sling) *DeviceService {
	return &DeviceService{c: c, sling: sling.Path("v2/device")}
}

type deviceCreateParams struct {
	OS             string
	InstallChannel string
	DeviceKey      string
}

type Device struct {
	ID    uint64 `json:"DeviceID,omitempty"`
	Key   string `json:"DeviceKey,omitempty"`
	Token string `json:"DeviceToken,omitempty"`
}

func (s *DeviceService) ensure() error {
	device := new(Device)
	if s.c.Device() != nil {
		device = s.c.Device()
		if device.Token != "" {
			return nil
		}
	}
	apiError := new(APIError)
	slingReq := s.c.SetAuthHeaders(s.sling.New())

	cp := &deviceCreateParams{OS: runtime.GOOS}
	cp.DeviceKey = device.Key

	_, err := slingReq.Post("").BodyJSON(cp).Receive(device, apiError)
	if relevantError(err, *apiError) == nil {
		// Set the device
		s.c.SetDevice(device)
	}
	return relevantError(err, *apiError)
}
