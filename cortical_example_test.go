package cortical_test

import (
	"context"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/owulveryck/cortical"
	"log"
	"net/http"
)

// echo is a dummy type that reads a message, and send back an "ack"
type echo struct{}

func new() *echo {
	return &echo{}
}

// NewCortex is filling the  ...
func (e *echo) NewCortex(ctx context.Context) (cortical.GetInfoFromCortexFunc, cortical.SendInfoToCortex) {
	c := make(chan []byte)
	return func(ctx context.Context) chan []byte {
			return c
		}, func(ctx context.Context, b *[]byte) {
			c <- *b
		}
}

func Example() {
	brain := &cortical.Cortical{
		Upgrader: websocket.Upgrader{},
		Cortexes: []cortical.Cortex{new()},
	}
	http.HandleFunc("/ws", brain.ServeWS)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, index)
	})

	log.Fatal(http.ListenAndServe(":8080", nil))

}

var index = `
<html>
    <head>
	<title>demo</title>
	<script type="text/javascript">
	window.onload = function () {
	    var conn;
	    var log = document.getElementById("log");
	    function appendLog(item) {
		var doScroll = log.scrollTop > log.scrollHeight - log.clientHeight - 1;
		log.appendChild(item);
		if (doScroll) {
		    log.scrollTop = log.scrollHeight - log.clientHeight;
		}
	    }
	    setInterval(function () {
		if (!conn) {
		    return false;
		}
		conn.send("ping");
		return false;
	    },1000);
	    if (window["WebSocket"]) {
		conn = new WebSocket("ws://" + document.location.host + "/ws");
		conn.onclose = function (evt) {
		    var item = document.createElement("div");
		    item.innerHTML = "<b>Connection closed.</b>";
		    appendLog(item);
		};
		conn.onmessage = function (evt) {
		    var messages = evt.data.split('\n');
		    for (var i = 0; i < messages.length; i++) {
			var item = document.createElement("div");
			item.innerText = messages[i];
			appendLog(item);
		    }
		};
	    } else {
		var item = document.createElement("div");
		item.innerHTML = "<b>Your browser does not support WebSockets.</b>";
		appendLog(item);
	    }
	};
	</script>
    </head>
    <body>
	<p id="log"> </p>
	<script>

	</script>
    </body>
</html>
`
