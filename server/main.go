package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"

	//"io"
	"fmt"
        "strings"
)

var upgrader = websocket.Upgrader{}

func spam(feedback_channel chan bool, connection *websocket.Conn){
   log.Println("Spamming the client")
   for i := 0; i < 10; i ++ {
      connection.WriteMessage(websocket.TextMessage,[]byte("Hello!"))
      time.Sleep(time.Millisecond * 200)
   }
   log.Println("Spammer is finished")
   close(feedback_channel)
}

func dispatch(message string){
   log.Println(fmt.Sprintf("Handling %s", message))
}

func read_from_connection(con *websocket.Conn, output chan string){
   log.Println("Reading from connection")
   defer close(output)
   for {
      message_type, message, err := con.ReadMessage()
      if err != nil {
         log.Println(fmt.Sprintf("Error while reading from websocket: %v", err))
         break
      }
      if message_type != websocket.TextMessage {
         log.Println("Got unexpected message type bytes")
         break
      }
      message_string := strings.TrimSpace(string(message))
      if message_string == "quit" {
         log.Println("Received quit")
         break
      }
      log.Println(fmt.Sprintf("Got %s", message_string))
      output <- message_string 
   }
}

func main(){
   http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request){
      log.Println("New connection")
      con, err := upgrader.Upgrade(w,req,nil)
      if err != nil {
         log.Println("upgrade:",err)
         return
      }
      defer con.Close()
      connection_output := make(chan string)
      go read_from_connection(con,connection_output)
      for message := range connection_output {
         log.Println(fmt.Sprintf("Handling message %s", message))
      }
      log.Println("Goodbye!")
   })
   panic(http.ListenAndServe(":8080", nil))
}
