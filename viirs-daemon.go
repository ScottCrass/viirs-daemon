package main

import (
    "bufio"
    "fmt"
    "net"
    "os"
    "strings"
    "database/sql"
   _ "github.com/lib/pq"
)

const (
  host     = "*****"
  port     = 5432
  user     = "*****"
  dbname   = "*****"
  password   = "*****"
)

type Client struct {
    socket net.Conn
    data   chan []byte
}

var DB *sql.DB

// ConnectDB opens a connection to the database
func ConnectDB() *sql.DB {
  psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
  db, err := sql.Open("postgres", psqlInfo)

  if err != nil {
    panic(err)
  }

 return db
}


func (client *Client) receive() {
    for {
        
    message := make([]byte, 4096)
        if len(message) == 0{
    fmt.Println("Empty message")
    }
    length, err := client.socket.Read(message)
        if err != nil {
            client.socket.Close()
            fmt.Println(err)
            panic(err)
        }
        if length > 0 {
	    fmt.Println(string(message))
            insertRecord(string(message))
        }
        if length == 0 {
           fmt.Println("Panic: Empty record. ")
           panic(err)
        }
    }
}

func insertRecord(s string) {
    records := strings.Fields(s)
    sqlStatement := `INSERT INTO okce_only (co2filtered, sensor_temp, co2raw)  VALUES ($1, $2, $3)`
     res, err := DB.Exec(sqlStatement, records[0], records[1], records[2])
     fmt.Println(sqlStatement)
     fmt.Println(records)
      if err != nil {
        panic(err)
     }
     if res != nil {
         fmt.Println(res)
     }
}

func startClientMode() {
    fmt.Println("Opening Database Connection...")
    DB = ConnectDB()
    fmt.Println("Starting client...")
    connection, error := net.Dial("tcp", "127.0.0.1:1234")
    if error != nil {
        fmt.Println(error)
    }
    client := &Client{socket: connection}
    go client.receive()
    for {
        reader := bufio.NewReader(os.Stdin)
        message, _ := reader.ReadString('\n')
        connection.Write([]byte(strings.TrimRight(message, "\n")))
    }
}

func main() {

        startClientMode()
 
}

