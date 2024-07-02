package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

func webGame() {
	log.SetFlags(0)
	http.HandleFunc("/new", newHandler)
	http.HandleFunc("/", home)
	http.ListenAndServe(*addr, nil)
	log.Fatal(http.ListenAndServe(*addr, nil))

}
func home(w http.ResponseWriter, r *http.Request) {
	homeTemplate.Execute(w, "ws://"+r.Host+"/new")
}

type HangmanWeb struct {
	conn *websocket.Conn
}

func (h *HangmanWeb) GetDisplayConn() interface{} {
	return h.conn
}

func (h *HangmanWeb) RenderGame(placeholder []string, chances int, entries map[string]bool) {

	keys := []string{}
	for key, _ := range entries {
		keys = append(keys, key)
	}
	c := h.GetDisplayConn().(*websocket.Conn)
	str := fmt.Sprintf("%v Chances left:%d Guesses:%v", placeholder, chances, keys)
	c.WriteMessage(websocket.TextMessage, []byte(str))
}
func (h *HangmanWeb) getInput() string {
	c := h.GetDisplayConn().(*websocket.Conn)
	_, message, err := c.ReadMessage()
	if err != nil {
		log.Println("read:", err)
		return ""
	}

	return string(message)
}

var upgrader = websocket.Upgrader{} // use default options
func newHandler(w http.ResponseWriter, r *http.Request) {

	hangman := HangmanWeb{}

	var word string = WordConstant
	if !*development {
		word = getWord()
	}

	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}

	hangman.conn = c

	defer c.Close()
	if play(&hangman, word) {
		err = c.WriteMessage(websocket.TextMessage, []byte("You win"))
	} else {
		err = c.WriteMessage(websocket.TextMessage, []byte("You hanged up"))
		err = c.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("word was: %v", word)))
	}

	if err != nil {
		log.Println("write:", err)
		return
	}

}

var homeTemplate = template.Must(template.New("").Parse(`
<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<script>  
window.addEventListener("load", function(evt) {

    var output = document.getElementById("output");
    var input = document.getElementById("input");
    var ws;

    var print = function(message) {
        var d = document.createElement("div");
        d.textContent = message;
        output.appendChild(d);
        output.scroll(0, output.scrollHeight);
    };

    document.getElementById("open").onclick = function(evt) {
        if (ws) {
            return false;
        }
        ws = new WebSocket("{{.}}");
        ws.onopen = function(evt) {
            print("OPEN");
        }
        ws.onclose = function(evt) {
            print("CLOSE");
            ws = null;
        }
        ws.onmessage = function(evt) {
            print("RESPONSE: " + evt.data);
        }
        ws.onerror = function(evt) {
            print("ERROR: " + evt.data);
        }
        return false;
    };

    document.getElementById("send").onclick = function(evt) {
        if (!ws) {
            return false;
        }
        print("SEND: " + input.value);
        ws.send(input.value);
        return false;
    };

    document.getElementById("close").onclick = function(evt) {
        if (!ws) {
            return false;
        }
        ws.close();
        return false;
    };

});
</script>
</head>
<body>
<table>
<tr><td valign="top" width="50%">
<form>
<button id="open">Start New Game</button>
<p><input id="input" type="text" value="Hello world!">
<button id="send">Send</button>
</form>
</td><td valign="top" width="50%">
<div id="output" style="max-height: 70vh;overflow-y: scroll;"></div>
</td></tr></table>
</body>
</html>
`))
