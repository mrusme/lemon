package pushover

import (
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"image/png"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/gorilla/websocket"
	"github.com/mrusme/lemon/inbox"
)

type PushoverStreamReturn int

const (
	PushoverOk PushoverStreamReturn = iota
	PushoverError
	PushoverReconnectRequest
	PushoverReauthenticateRequest
	PushoverEndRequest
)

const PushoverIconsURLFmt = "https://api.pushover.net/icons/%s.png"

type PushoverMessage struct {
	ID      int    `json:"id"`
	IDstr   string `json:"id_str"`
	UMID    int    `json:"umid"`
	UMIDstr string `json:"umid_str"`
	AID     int    `json:"aid"`
	AIDstr  string `json:"aid_str"`
	App     string `json:"app"`

	Title          string `json:"title"`
	Message        string `json:"message"`
	Icon           string `json:"icon"`
	Date           int    `json:"date"`
	QueuedDate     int    `json:"queued_date"`
	DispatchedDate int    `json:"dispatched_date"`
	Priority       int    `json:"priority"`
	Sound          string `json:"sound"`
	URL            string `json:"url"`
	URLTitle       string `json:"url_title"`
	Acked          int    `json:"acked"`
	Receipt        string `json:"receipt"`
	HTML           int    `json:"html"`
}

type PushoverMessagesResponse struct {
	Status   int               `json:"status"`
	Request  string            `json:"request"`
	Messages []PushoverMessage `json:"messages"`
}

type Pushover struct {
	deviceId string
	secret   string
	ibx      chan inbox.Message
}

func New(ibx chan inbox.Message, deviceId string, secret string) (*Pushover, error) {
	var po = new(Pushover)

	po.deviceId = deviceId
	po.secret = secret
	po.ibx = ibx

	return po, nil
}

func (po *Pushover) Stream() {
	for {
		status, err := po.stream()
		switch status {
		case PushoverOk:
			log.Println("pushover terminated normally, quitting")
			os.Exit(0)
		case PushoverError:
			log.Printf("pushover error: %s", err)
			os.Exit(int(status))
		case PushoverReconnectRequest:
			log.Println("pushover requested reconnect")
			continue
		case PushoverReauthenticateRequest:
			log.Println("pushover requested re-auth, quitting")
			os.Exit(int(status))
		case PushoverEndRequest:
			log.Println("pushover requested end, quitting")
			os.Exit(int(status))
		}
	}
}

func (po *Pushover) stream() (PushoverStreamReturn, error) {
	u := url.URL{
		Scheme: "wss",
		Host:   "client.pushover.net",
		Path:   "/push",
	}

	log.Printf("connecting to %s ...", u.String())

	ws, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return PushoverError, err
	}
	defer ws.Close()

	ws.WriteMessage(
		websocket.TextMessage,
		[]byte(fmt.Sprintf(
			"login:%s:%s\n",
			po.deviceId,
			po.secret,
		)),
	)

	for {
		_, push, err := ws.ReadMessage()
		if err != nil {
			return PushoverError, err
		}
		pushStr := string(push)
		switch pushStr {
		case "#":
			// Keep-alive packet, no response needed.
			continue
		case "!":
			// A new message has arrived; you should perform a sync.
			msgs, err := po.getMessages()
			if err != nil {
				log.Println(err)
			}

			for _, msg := range msgs {
				fmt.Println(msg)
				icon, iconPath, err := po.getIcon(msg.Icon)
				if err != nil {
					log.Printf("Error: %s\n", err)
					continue
				}
				// TODO: if msg.HTML == 1 { convert HTML to text }
				ibxMsg := inbox.Message{
					Icon:     icon,
					IconPath: iconPath,
					Title:    msg.Title,
					Text:     msg.Message,
					URL:      msg.URL,
				}
				po.ibx <- ibxMsg
			}

			if err := po.deleteMessages(msgs); err != nil {
				log.Println(err)
			}
			// Reload request; you should drop your connection and re-connect.
			return PushoverReconnectRequest, nil
		case "E":
			// Error; a permanent problem occured and you should not automatically
			// re-connect. Prompt the user to login again or re-enable the device.
			return PushoverReauthenticateRequest, nil
		case "A":
			// Error; the device logged in from another session and this session is
			// being closed. Do not automatically re-connect.
			return PushoverEndRequest, nil
		}
	}
}

func (po *Pushover) getMessages() ([]PushoverMessage, error) {
	var err error
	u := "https://api.pushover.net/1/messages"

	pushoverClient := http.Client{
		Timeout: time.Second * 30,
	}

	req, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "lemon")

	q := req.URL.Query()
	q.Add("device_id", po.deviceId)
	q.Add("secret", po.secret)
	req.URL.RawQuery = q.Encode()

	res, err := pushoverClient.Do(req)
	if err != nil {
		return nil, err
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	messagesResponse := PushoverMessagesResponse{}
	err = json.Unmarshal(body, &messagesResponse)
	if err != nil {
		return nil, err
	}

	if messagesResponse.Status != 1 {
		return nil, errors.New(fmt.Sprintf("Status was %d\n", messagesResponse.Status))
	}

	return messagesResponse.Messages, nil
}

func (po *Pushover) deleteMessages(msgs []PushoverMessage) error {
	data := url.Values{
		"secret":  {po.secret},
		"message": {msgs[len(msgs)-1].IDstr},
	}

	resp, err := http.PostForm(
		"https://api.pushover.net/1/devices/"+
			po.deviceId+"/update_highest_message.json",
		data,
	)
	if err != nil {
		return err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return errors.New("Status code not 200")
	}

	return nil
}

func (po *Pushover) getIcon(iconName string) (image.Image, string, error) {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return nil, "", err
	}

	if err := os.MkdirAll(filepath.Join(cacheDir, "lemon"), 0755); err != nil {
		return nil, "", err
	}

	filePath := filepath.Join(cacheDir, "lemon", iconName+".png")
	_, err = os.Stat(filePath)
	if os.IsNotExist(err) {
		if err := po.downloadIcon(iconName, filePath); err != nil {
			return nil, "", err
		}
	}

	file, err := os.Open(filePath)
	if err != nil {
		return nil, "", err
	}
	defer file.Close()

	img, err := png.Decode(file)
	if err != nil {
		return nil, "", err
	}

	return img, filePath, nil
}

func (po *Pushover) downloadIcon(iconName string, filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	response, err := http.Get(fmt.Sprintf(PushoverIconsURLFmt, iconName))
	if err != nil {
		return err
	}
	defer response.Body.Close()

	_, err = io.Copy(file, response.Body)
	if err != nil {
		return err
	}

	return nil
}
