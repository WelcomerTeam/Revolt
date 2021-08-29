package revolt

import (
	"bytes"
	"context"
	"image/color"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"strconv"
	"sync"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/savsgio/gotils"
	"github.com/vmihailenco/msgpack"
	"nhooyr.io/websocket"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

const RevoltWS = "wss://ws.revolt.chat?format=json"
const RevoltHTTPBase = "https://api.revolt.chat"
const AutumnHTTPBase = "https://autumn.revolt.chat"
const usingMsgpack = false

func init() {
	if usingMsgpack {
		panic("stop using msgpack")
	}
}

type RevoltBot struct {
	ctx context.Context

	Token string

	usersMu sync.RWMutex
	Users   map[string]*User

	guildsMu sync.RWMutex
	Guilds   map[string]*Guild

	channelsMu sync.RWMutex
	Channels   map[string]*Channel

	membersMu sync.RWMutex
	Members   map[string]*GuildMember

	wsConn *websocket.Conn
}

func NewRevoltBot(token string) (rb *RevoltBot) {
	ctx := context.WithValue(context.Background(), "revolt", "revolt-bot")

	rb = &RevoltBot{
		ctx:   ctx,
		Token: token,

		Users:    make(map[string]*User),
		Guilds:   make(map[string]*Guild),
		Channels: make(map[string]*Channel),
		Members:  make(map[string]*GuildMember),
	}

	return rb
}

func (rb *RevoltBot) Post(path string, data interface{}) (resp *http.Response, err error) {
	buf := bytes.NewBuffer(nil)
	json.NewEncoder(buf).Encode(data)

	req, err := http.NewRequest("POST", RevoltHTTPBase+path, buf)
	if err != nil {
		return nil, err
	}

	req.Header.Set("x-bot-token", rb.Token)
	return http.DefaultClient.Do(req)
}

func (rb *RevoltBot) Get(path string) (resp *http.Response, err error) {
	req, err := http.NewRequest("GET", RevoltHTTPBase+path, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("x-bot-token", rb.Token)
	return http.DefaultClient.Do(req)
}

func (rb *RevoltBot) UploadFile(fileName string, fileContent []byte) (autumnID string, err error) {
	b := new(bytes.Buffer)
	w := multipart.NewWriter(b)
	part, _ := w.CreateFormFile("file", fileName)
	part.Write(fileContent)
	w.Close()

	req, err := http.NewRequest("POST", AutumnHTTPBase+"/attachments", b)
	if err != nil {
		return "", err
	}

	req.Header.Set("x-bot-token", rb.Token)
	req.Header.Set("Content-Type", w.FormDataContentType())
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}

	dat, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	autumnID = json.Get(dat, "id").ToString()

	return autumnID, nil
}

func (rb *RevoltBot) FetchUser(userID string) (user *User, err error) {
	resp, err := rb.Get("/users/" + userID)
	if err != nil {
		return nil, err
	}

	res, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(res, &user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (rb *RevoltBot) SendMessage(channelID string, messageRequest *MessageRequest) (message *Message, err error) {
	if messageRequest.Nonce == "" {
		messageRequest.Nonce = strconv.FormatInt(time.Now().Unix(), 10)
		println("Dont forget to add a nonce")
	}

	resp, err := rb.Post("/channels/"+channelID+"/messages", messageRequest)
	if err != nil {
		return nil, err
	}

	res, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	println(gotils.B2S(res))

	err = json.Unmarshal(res, &message)
	if err != nil {
		return nil, err
	}

	return message, nil
}

func (rb *RevoltBot) Start() (err error) {
	conn, _, err := websocket.Dial(rb.ctx, RevoltWS, nil)
	if err != nil {
		return err
	}

	rb.wsConn = conn

	println("CONNECTED TO " + RevoltWS)

	go rb.Heartbeat()

	rb.SendEvent(Authenticate{
		SentBase: SentBase{"Authenticate"},
		Token:    &rb.Token,
	})

	for {
		_, buf, err := rb.wsConn.Read(rb.ctx)
		if err != nil {
			println(err)
			return err
		}

		mType := json.Get(buf, "type").ToString()

		go rb.OnDispatch(mType, buf)
	}
}

func (rb *RevoltBot) Heartbeat() {
	t := time.NewTicker(time.Second * 20)

	for {
		select {
		case <-t.C:
			rb.SendEvent(Ping{
				SentBase: SentBase{"Ping"},
				Time:     int(time.Now().Unix()),
			})
		case <-rb.ctx.Done():
			return
		}
	}
}

func (rb *RevoltBot) SendEvent(data interface{}) (err error) {
	var val []byte

	if usingMsgpack {
		val, err = msgpack.Marshal(data)
	} else {
		val, err = json.Marshal(data)
	}

	if err != nil {
		println(err.Error())

		return err
	}

	println("<-", gotils.B2S(val))

	return rb.wsConn.Write(rb.ctx, websocket.MessageText, val)
}

func (rb *RevoltBot) OnDispatch(messageType string, data []byte) (err error) {
	println("-> ", gotils.B2S(data))

	switch messageType {
	case "Authenticated":
		o := Authenticated{}
		err = json.Unmarshal(data, &o)
		if err != nil {
			return err
		}

		rb.OnAuthenticated(o)
	case "Pong":
		o := Pong{}
		err = json.Unmarshal(data, &o)
		if err != nil {
			return err
		}

		rb.OnPong(o)
	case "Ready":
		o := Ready{}
		err = json.Unmarshal(data, &o)
		if err != nil {
			return err
		}

		rb.OnReady(o)
	case "Message":
		o := MessageCreate{}
		err = json.Unmarshal(data, &o)
		if err != nil {
			return err
		}

		rb.OnMessageCreate(o)
	case "MessageUpdate":
		o := MessageUpdate{}
		err = json.Unmarshal(data, &o)
		if err != nil {
			return err
		}

		rb.OnMessageUpdate(o)
	case "MessageDelete":
		o := MessageDelete{}
		err = json.Unmarshal(data, &o)
		if err != nil {
			return err
		}

		rb.OnMessageDelete(o)
	case "ChannelCreate":
		o := ChannelCreate{}
		err = json.Unmarshal(data, &o)
		if err != nil {
			return err
		}

		rb.OnChannelCreate(o)
	case "ChannelUpdate":
		o := ChannelUpdate{}
		err = json.Unmarshal(data, &o)
		if err != nil {
			return err
		}

		rb.OnChannelUpdate(o)
	case "ChannelDelete":
		o := ChannelDelete{}
		err = json.Unmarshal(data, &o)
		if err != nil {
			return err
		}

		rb.OnChannelDelete(o)
	case "ChannelGroupJoin":
		o := ChannelGroupJoin{}
		err = json.Unmarshal(data, &o)
		if err != nil {
			return err
		}

		rb.OnChannelGroupJoin(o)
	case "ChannelGroupLeave":
		o := ChannelGroupLeave{}
		err = json.Unmarshal(data, &o)
		if err != nil {
			return err
		}

		rb.OnChannelGroupLeave(o)
	case "ChannelStartTyping":
		o := ChannelStartTyping{}
		err = json.Unmarshal(data, &o)
		if err != nil {
			return err
		}

		rb.OnChannelStartTyping(o)
	case "ChannelStopTyping":
		o := ChannelStopTyping{}
		err = json.Unmarshal(data, &o)
		if err != nil {
			return err
		}

		rb.OnChannelStopTyping(o)
	case "ChannelAck":
		o := ChannelAck{}
		err = json.Unmarshal(data, &o)
		if err != nil {
			return err
		}

		rb.OnChannelAck(o)
	case "ServerUpdate":
		o := ServerUpdate{}
		err = json.Unmarshal(data, &o)
		if err != nil {
			return err
		}

		rb.OnServerUpdate(o)
	case "ServerDelete":
		o := ServerDelete{}
		err = json.Unmarshal(data, &o)
		if err != nil {
			return err
		}

		rb.OnServerDelete(o)
	case "ServerMemberUpdate":
		o := ServerMemberUpdate{}
		err = json.Unmarshal(data, &o)
		if err != nil {
			return err
		}

		rb.OnServerMemberUpdate(o)
	case "ServerMemberJoin":
		o := ServerMemberJoin{}
		err = json.Unmarshal(data, &o)
		if err != nil {
			return err
		}

		rb.OnServerMemberJoin(o)
	case "ServerMemberLeave":
		o := ServerMemberLeave{}
		err = json.Unmarshal(data, &o)
		if err != nil {
			return err
		}

		rb.OnServerMemberLeave(o)
	case "ServerRoleUpdate":
		o := ServerRoleUpdate{}
		err = json.Unmarshal(data, &o)
		if err != nil {
			return err
		}

		rb.OnServerRoleUpdate(o)
	case "ServerRoleDelete":
		o := ServerRoleDelete{}
		err = json.Unmarshal(data, &o)
		if err != nil {
			return err
		}

		rb.OnServerRoleDelete(o)
	case "UserUpdate":
		o := UserUpdate{}
		err = json.Unmarshal(data, &o)
		if err != nil {
			return err
		}

		rb.OnUserUpdate(o)
	case "UserRelationship":
		o := UserRelationship{}
		err = json.Unmarshal(data, &o)
		if err != nil {
			return err
		}

		rb.OnUserRelationship(o)
	default:
		println(messageType + " not implemented")
	}

	return
}

func (rb *RevoltBot) OnAuthenticated(o Authenticated) {}
func (rb *RevoltBot) OnPong(o Pong)                   {}
func (rb *RevoltBot) OnReady(o Ready) {
	rb.channelsMu.Lock()
	for _, c := range o.Channels {
		rb.Channels[c.ID] = c
	}
	rb.channelsMu.Unlock()

	rb.guildsMu.Lock()
	for _, g := range o.Guilds {
		rb.Guilds[g.ID] = g
	}
	rb.guildsMu.Unlock()

	rb.usersMu.Lock()
	for _, u := range o.Users {
		rb.Users[u.ID] = u
	}
	rb.usersMu.Unlock()
}
func (rb *RevoltBot) OnMessageCreate(o MessageCreate) {
	if v, ok := o.Message.RawContent.(string); ok {
		o.Message.ContentType = "message"
		o.Message.Content = v
	} else if v, ok := o.Message.RawContent.(*MessageContent); ok {
		o.Message.ContentType = v.Type
		o.Message.By = *v.By
		o.Message.Content = *v.Content
		o.Message.TargetID = *v.ID
		o.Message.Name = *v.Name
	}

	if o.Message.Content == "/pog" {
		_, err := rb.SendMessage(o.Message.ChannelID, &MessageRequest{
			Content: "pog",
		})
		if err != nil {
			println(err.Error())
		}
	}

	if o.Message.Content == "/rock" {
		f, _ := ioutil.ReadFile("C:/users/blane/desktop/realrock.png")
		autumnID, err := rb.UploadFile("rock.png", f)
		if err != nil {
			println(err.Error())
		}

		println(autumnID)

		_, err = rb.SendMessage(o.Message.ChannelID, &MessageRequest{
			Content:     "heres a rock",
			Attachments: []string{autumnID},
			Nonce:       strconv.FormatInt(time.Now().Unix(), 10),
		})
		if err != nil {
			println(err.Error())
		}
	}

	// println(o.Message.Author + " said '" + o.Message.Content + "' in channel " + o.Message.ChannelID)
}
func (rb *RevoltBot) OnMessageUpdate(o MessageUpdate)           {}
func (rb *RevoltBot) OnMessageDelete(o MessageDelete)           {}
func (rb *RevoltBot) OnChannelCreate(o ChannelCreate)           {}
func (rb *RevoltBot) OnChannelUpdate(o ChannelUpdate)           {}
func (rb *RevoltBot) OnChannelDelete(o ChannelDelete)           {}
func (rb *RevoltBot) OnChannelGroupJoin(o ChannelGroupJoin)     {}
func (rb *RevoltBot) OnChannelGroupLeave(o ChannelGroupLeave)   {}
func (rb *RevoltBot) OnChannelStartTyping(o ChannelStartTyping) {}
func (rb *RevoltBot) OnChannelStopTyping(o ChannelStopTyping)   {}
func (rb *RevoltBot) OnChannelAck(o ChannelAck)                 {}
func (rb *RevoltBot) OnServerUpdate(o ServerUpdate)             {}
func (rb *RevoltBot) OnServerDelete(o ServerDelete)             {}
func (rb *RevoltBot) OnServerMemberUpdate(o ServerMemberUpdate) {}
func (rb *RevoltBot) OnServerMemberJoin(o ServerMemberJoin) {

	var b bytes.Buffer

	user, err := rb.FetchUser(o.UserID)
	if err != nil {
		println(err.Error())

		return
	}

	json.NewEncoder(&b).Encode(ImageCreateArguments{
		FilesizeLimit: 10000000,
		Options: ImageOpts{
			Text:                "Welcome " + user.Username,
			ImageURL:            "https://autumn.revolt.chat/avatars/" + user.Avatar.ID + "?format=png&max_side=256",
			Background:          "revolt",
			Font:                "Raleway-Bold",
			BorderColour:        color.RGBA{253, 68, 83, 0},
			BorderWidth:         16,
			TextAlignmentX:      1,
			TextAlignmentY:      1,
			ProfileBorderColour: color.RGBA{17, 24, 34, 0},
			TextStroke:          8,
			TextStrokeColour:    color.RGBA{255, 255, 255, 0},
			TextColour:          color.RGBA{253, 68, 83, 0},
		},
	})

	resp, _ := http.Post("http://localhost:4200/images", "application/json", &b)
	body, _ := ioutil.ReadAll(resp.Body)

	autumnID, _ := rb.UploadFile("welcome.png", body)

	rb.guildsMu.RLock()
	g := rb.Guilds[o.GuildID]
	rb.guildsMu.RUnlock()

	msg, err := rb.SendMessage(g.SystemMessages.UserJoined, &MessageRequest{
		Attachments: []string{autumnID},
	})

	println(msg, msg.ID)

	if err != nil {
		println(err.Error())
	}
}
func (rb *RevoltBot) OnServerMemberLeave(o ServerMemberLeave) {}
func (rb *RevoltBot) OnServerRoleUpdate(o ServerRoleUpdate)   {}
func (rb *RevoltBot) OnServerRoleDelete(o ServerRoleDelete)   {}
func (rb *RevoltBot) OnUserUpdate(o UserUpdate)               {}
func (rb *RevoltBot) OnUserRelationship(o UserRelationship)   {}
