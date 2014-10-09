package main

import (
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
	//"os"
	"code.google.com/p/go.net/websocket"
	// "io"
	// "io/ioutil"
	"encoding/json"
	"net/http"
)

// Echo the data received on the WebSocket.
/*func EchoServer(ws *websocket.Conn) {
	io.Copy(ws, ws)
}*/

type TextPayload struct {
	Kind  string `json:"kind"`
	Value string `json:"value"`
}

type GPPayload struct {
	Kind  string `json:"kind"`
	Which uint8  `json:"which"`
	Value int16  `json:"value"`
}

var gamepadConn *websocket.Conn

var gp GPPayload

func sendMsg(kind string, which uint8, value int16) {
	if gamepadConn == nil {
		return
	}
	gp.Kind = kind
	gp.Which = which
	gp.Value = value
	out_, err := json.Marshal(gp)
	if err != nil {
		return
	}
	//fmt.Printf("|-> %s\n", out_)
	websocket.Message.Send(gamepadConn, string(out_))
}

// Echo the data received on the WebSocket.
func GamepadServer(ws *websocket.Conn) {
	fmt.Printf("gamepad handler called...\n")

	var tp TextPayload
	var in []byte
	var out []byte
	var err error

	gamepadConn = ws

	for {
		websocket.Message.Receive(ws, &in)
		//fmt.Printf("|<- %s\n", in)

		err = json.Unmarshal(in, &tp)
		if err != nil {
			fmt.Errorf("%s\n", err)
			continue
		}

		tp.Value += " ... back at ya!"

		out, err = json.Marshal(tp)
		if err != nil {
			fmt.Errorf("%s\n", err)
			continue
		}
		//fmt.Printf("|-> %s\n", out)
		websocket.Message.Send(ws, string(out))
	}
}

func WSServer() {
	// http.Handle("/echo", websocket.Handler(EchoServer))
	http.Handle("/gamepad", websocket.Handler(GamepadServer))

	fmt.Printf("WS SERVER running on 127.0.0.1:12345...\n")

	err := http.ListenAndServe(":12345", nil)

	/*serverAddr = server.Listener.Addr().String()
	  log.Print("Test WebSocket server listening on ", serverAddr)*/

	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}

// var winTitle string = "Go-SDL2 Events"
// var winWidth, winHeight int = 100, 100

func main() {
	// fmt.Println("1")
	go WSServer()
	// fmt.Println("2")

	// var window *sdl.Window
	// var renderer *sdl.Renderer
	var event sdl.Event
	var running bool

	/*window = sdl.CreateWindow(winTitle, sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		winWidth, winHeight, sdl.WINDOW_SHOWN)
	if window == nil {
		fmt.Fprintf(os.Stderr, "Failed to create window: %s\n", sdl.GetError())
		os.Exit(1)
	}

	renderer = sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if renderer == nil {
		fmt.Fprintf(os.Stderr, "Failed to create renderer: %s\n", sdl.GetError())
		os.Exit(2)
	}*/

	sdl.InitSubSystem(sdl.INIT_JOYSTICK)

	numJoys := sdl.NumJoysticks()
	fmt.Printf("# joysticks: %d\n", numJoys)

	if numJoys > 0 {
		joyName := sdl.JoystickNameForIndex(0)
		fmt.Printf("name: %s\n\n", joyName)
		joy := sdl.JoystickOpen(0)
		fmt.Printf("# axis:    %d\n# balls:   %d\n# buttons: %d\n# hats:    %d\n\n", joy.NumAxes(), joy.NumBalls(), joy.NumButtons(), joy.NumHats())
	}

	sdl.JoystickEventState(sdl.ENABLE)

	/*
		SDL_JoyAxisEvent	Occurs when an axis changes
		SDL_JoyBallEvent	Occurs when a joystick trackball's position changes
		SDL_JoyHatEvent	    Occurs when a hat's position changes
		SDL_JoyButtonEvent	Occurs when a button is pressed or released
	*/

	running = true
	for running {
		for event = sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch t := event.(type) {
			case *sdl.QuitEvent:
				running = false

			/*case *sdl.MouseMotionEvent:
				fmt.Printf("[%d ms] MouseMotion\ttype:%d\tid:%d\tx:%d\ty:%d\txrel:%d\tyrel:%d\n",
					t.Timestamp, t.Type, t.Which, t.X, t.Y, t.XRel, t.YRel)
			case *sdl.MouseButtonEvent:
				fmt.Printf("[%d ms] MouseButton\ttype:%d\tid:%d\tx:%d\ty:%d\tbutton:%d\tstate:%d\n",
					t.Timestamp, t.Type, t.Which, t.X, t.Y, t.Button, t.State)
			case *sdl.MouseWheelEvent:
				fmt.Printf("[%d ms] MouseWheel\ttype:%d\tid:%d\tx:%d\ty:%d\n",
					t.Timestamp, t.Type, t.Which, t.X, t.Y)*/

			/*case *sdl.KeyUpEvent:
			fmt.Printf("KEY    sym: %c | modifiers: %d | state: %d | repeat: %d\n", t.Keysym.Sym, t.Keysym.Mod, t.State, t.Repeat)*/

			case *sdl.JoyAxisEvent:
				// fmt.Printf("AXIS   which: %d | axis:   %d | value: %d\n", t.Which, t.Axis, t.Value)
				sendMsg("axis", t.Axis, t.Value)
			case *sdl.JoyButtonEvent:
				// fmt.Printf("BUTTON which: %d | button: %d | state: %d\n", t.Which, t.Button, t.State)
				sendMsg("button", t.Button, int16(t.State))
			case *sdl.JoyHatEvent:
				// fmt.Printf("HAT    which: %d | hat:    %d | value: %d\n", t.Which, t.Hat, t.Value)
				sendMsg("hat", t.Hat, int16(t.Value))
				/*case *sdl.JoyBallEvent:
				fmt.Printf("BALL   which: %d | ball:   %d | relPos: ( %d , %d )\n", t.Which, t.Ball, t.XRel, t.YRel)*/
			}
		}
	}

	/*renderer.Destroy()
	window.Destroy()*/
}
