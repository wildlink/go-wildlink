package wildlink

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/dghubble/sling"
)

func (c *Client) SetAuthHeaders(s *sling.Sling) *sling.Sling {

	timeValue := time.Now().UTC().Format("2006-01-02T15:04:05.000Z")
	var senderToken, deviceToken string
	if c.Device() != nil {
		deviceToken = c.Device().Token
	}

	sigString := timeValue + "\n" +
		deviceToken + "\n" +
		senderToken + "\n"

	appSignature := computeHexHmac256(sigString, c.appKey)
	return s.Set("X-WF-DateTime", timeValue).Set("Authorization", fmt.Sprintf("WFAV1 %v:%v:%v:%v", c.appID, appSignature, deviceToken, senderToken))

}

func computeHexHmac256(message string, secret string) string {
	key := []byte(secret)
	h := hmac.New(sha256.New, key)
	_, _ = h.Write([]byte(message))
	return hex.EncodeToString(h.Sum(nil))
}
