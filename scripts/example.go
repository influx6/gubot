package scripts

import (
	"github.com/ArthurHlt/gubot/robot"
	"fmt"
	"time"
	"net/http"
	"encoding/json"
	"github.com/gorilla/mux"
)

func init() {
	var conf ExampleScriptConfig
	robot.GetConfig(&conf)
	e := &ExampleScript{
		annoy: make(map[string]bool),
		config: conf,
	}

	e.listen()
	robot.Router().HandleFunc("/gubot/chatsecrets/{channel}", e.handlerChatsecret)
	robot.RegisterScripts([]robot.Script{
		{
			Name: "badger",
			Matcher: "(?i)badger",
			Function: e.badger,
			Type: robot.Tsend,
		},
		{
			Name: "doors",
			Matcher: "(?i)open the (.*) doors",
			Function: e.doors,
			Type: robot.Trespond,
		},
		{
			Name: "lulz",
			Matcher: "(?i)lulz",
			Function: e.lulz,
			Type: robot.Tsend,
		},
		{
			Name: "Ultimate question",
			Description: "Answer to the ultimate question",
			Matcher: "(?i)what is the answer to the ultimate question of life",
			Function: e.ultimateQuestion,
			Type: robot.Tsend,
		},
		{
			Name: "annoy",
			Example: "`annoy me`",
			Matcher: "(?i)^annoy me",
			Function: e.annoyMe,
			Type: robot.Tsend,
		},
		{
			Name: "unannoy",
			Example: "`unannoy me`",
			Matcher: "(?i)^unannoy me",
			Function: e.unannoyMe,
			Type: robot.Tsend,
		},
	})
}

type ExampleScriptConfig struct {
	GubotAnswerToTheUltimateQuestionOfLifeTheUniverseAndEverything string
}
type ExampleScript struct {
	config ExampleScriptConfig
	annoy  map[string]bool
}
type SecretMessage struct {
	Secret string
}

func (e ExampleScript) badger(envelop robot.Envelop, subMatch [][]string) ([]string, error) {
	return []string{"Badgers? BADGERS? WE DON'T NEED NO STINKIN BADGERS"}, nil
}

func (e ExampleScript) doors(envelop robot.Envelop, subMatch [][]string) ([]string, error) {
	doorType := subMatch[0][1]
	if doorType == "pod bay" {
		return []string{"I'm afraid I can't let you do that."}, nil
	}
	return []string{"Opening " + doorType + " doors"}, nil
}

func (e ExampleScript) lulz(envelop robot.Envelop, subMatch [][]string) ([]string, error) {
	return []string{"lol", "rofl", "lmao"}, nil
}

func (e ExampleScript) topic(envelop robot.Envelop, subMatch [][]string) ([]string, error) {
	return []string{envelop.Message + "? that's a paddlin"}, nil
}

func (e *ExampleScript) ultimateQuestion(envelop robot.Envelop, subMatch [][]string) ([]string, error) {
	if e.config.GubotAnswerToTheUltimateQuestionOfLifeTheUniverseAndEverything == "" {
		return []string{"Missing GubotAnswerToTheUltimateQuestionOfLifeTheUniverseAndEverything config parameter: please set and try again"}, nil
	}
	return []string{e.config.GubotAnswerToTheUltimateQuestionOfLifeTheUniverseAndEverything + ", but what is the question?"}, nil
}
func (e *ExampleScript) annoyMe(envelop robot.Envelop, subMatch [][]string) ([]string, error) {
	e.annoy[envelop.ChannelName + envelop.ChannelId] = true
	go func() {
		for {
			if !e.annoy[envelop.ChannelName + envelop.ChannelId] {
				break
			}
			robot.SendMessages(envelop, "AAAAAAAAAAAEEEEEEEEEEEEEEEEEEEEEEEEIIIIIIIIHHHHHHHHHH")
			time.Sleep(2 * time.Second)
		}
	}()
	return []string{"Hey, want to hear the most annoying sound in the world?"}, nil
}
func (e *ExampleScript) unannoyMe(envelop robot.Envelop, subMatch [][]string) ([]string, error) {
	message := "Not annoying you right now, am I?"
	if e.annoy[envelop.ChannelName + envelop.ChannelId] {
		e.annoy[envelop.ChannelName + envelop.ChannelId] = false
		message = "GUYS, GUYS, GUYS!"
	}
	return []string{message}, nil
}
func (e ExampleScript) handlerChatsecret(w http.ResponseWriter, req *http.Request) {
	var secretMessage SecretMessage
	err := json.NewDecoder(req.Body).Decode(&secretMessage)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		b, _ := json.MarshalIndent(struct {
			Code    int
			Message string
		}{http.StatusBadRequest, err.Error()}, "", "\t")
		w.Write(b)
		return
	}
	vars := mux.Vars(req)
	err = robot.SendMessages(robot.Envelop{
		ChannelName: vars["channel"],
	}, "I have a secret: " + secretMessage.Secret)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		b, _ := json.MarshalIndent(struct {
			Code    int
			Message string
		}{http.StatusInternalServerError, err.Error()}, "", "\t")
		w.Write(b)
		return
	}
	w.WriteHeader(http.StatusOK)

}
func (e *ExampleScript) listen() {
	go func() {
		for event := range robot.On(robot.EVENT_ROBOT_CHANNEL_ENTER) {
			gubotEvent := robot.ToGubotEvent(event)
			err := robot.RespondMessages(gubotEvent.Envelop, "Hi", "Target Acquired", "Firing", "Hello friend.", "Gotcha", "I see you")
			if err != nil {
				fmt.Println(err)
			}
		}
	}()
	go func() {
		for event := range robot.On(robot.EVENT_ROBOT_CHANNEL_LEAVE) {
			gubotEvent := robot.ToGubotEvent(event)
			err := robot.RespondMessages(gubotEvent.Envelop, "Are you still there?", "Target lost", "Searching")
			if err != nil {
				fmt.Println(err)
			}
		}
	}()
	go func() {
		for event := range robot.On(robot.EVENT_ROBOT_USER_ONLINE) {
			gubotEvent := robot.ToGubotEvent(event)
			robot.SendMessages(gubotEvent.Envelop, "Hello again", "It's been a while")
		}
	}()

}