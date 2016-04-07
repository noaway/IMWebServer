package main

import (
	"IMWebServer/models"
	"IMWebServer/parsers"
	// "IMWebServer/models/IM_Server"
	"IMWebServer/services"
	// "github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"runtime"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	go func() {
		ds := services.DBServer{}
		ds.HeartbeatName = "db_server"
		ds.Spec = 5
		ds.IpPort = models.Config("db_proxy_server", "IpPort1")
		ds.Net = models.Config("db_proxy_server", "Net")
		ds.InitServer(models.DBChanSend)
		parsers.ProtoInit()
	}()

	// go func(c chan []byte) {
	// 	ds := services.LoginServer{}
	// 	ds.IpPort = models.Config("login_server", "IpPort1")
	// 	ds.Net = models.Config("login_server", "Net")
	// 	ds.InitServer(c)
	// }(models.LoginChan)

	// go func(c chan []byte) {
	// 	ds := services.RouteServer{}
	// 	ds.IpPort = models.Config("route_server", "IpPort1")
	// 	ds.Net = models.Config("route_server", "Net")
	// 	ds.InitServer(c)
	// }(models.RouteChan)

	// phead := &parsers.PduHeader{
	// 	Service_id: 0x0007,
	// 	Command_id: 0x0709,
	// }
	// data, _ := phead.RenderByte(&IM_Server.IMMsgServInfo{
	// 	Ip1:        proto.String(models.Config("login_server", "ListenIP1")),
	// 	Ip2:        proto.String(models.Config("login_server", "ListenIP2")),
	// 	Port:       proto.Uint32(0),
	// 	WebimPort:  proto.Uint32(uint32(models.ConfigInt("login_server", "ListenPort1"))),
	// 	MaxConnCnt: proto.Uint32(uint32(models.ConfigInt("default", "MaxConnCnt"))),
	// 	CurConnCnt: proto.Uint32(uint32(models.ConfigInt("default", "MaxConnCnt"))),
	// 	HostName:   proto.String(models.Config("default", "HostName")),
	// })

	// models.LoginChan <- data

	// go func() {
	// 	for msg := range models.DBChan {
	// 		log.Println(string(msg))
	// 	}
	// }()

	ip_port := models.Config("default", "IpPort")
	log.Printf("Listen : %s \n", ip_port)

	http.HandleFunc("/ws/echo", echo)
	http.HandleFunc("/ws/index", index)
	log.Fatalln(http.ListenAndServe(ip_port, nil))
}

func echo(w http.ResponseWriter, r *http.Request) {
	var upgrader = websocket.Upgrader{}
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	log.Println("new connection")
	go parsers.Exec(c)
}

func index(w http.ResponseWriter, r *http.Request) {
	html := `
            <script type="text/javascript">
         var sock = null;
         var wsuri = "ws://127.0.0.1:8080/ws/echo";

         window.onload = function() {

            console.log("onload");

            sock = new WebSocket(wsuri);

            sock.onopen = function() {
               console.log("connected to " + wsuri);
            }

            sock.onclose = function(e) {
               console.log("connection closed (" + e.code + ")");
            }

            sock.onmessage = function(e) {
               var node=document.createElement("li");
               var textnode=document.createTextNode(e.data);
               node.appendChild(textnode);
               document.getElementById("msg").appendChild(node);
            }
         };

         function send() {
            var msg = document.getElementById('message').value;
            sock.send(msg);
         };
      </script>
      <h1>WebSocket Echo Test</h1>
      <ul id="msg"></ul>
		<button onclick="send();">Send Message</button>
         Message: <textarea id="message" values="" rows="3" cols="20">
	`
	w.Write([]byte(html))
}
