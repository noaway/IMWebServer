package services

import (
	"IMWebServer/models"
	"IMWebServer/parsers"
	"github.com/noaway/heartbeat"
	"io"
	"log"
	"net"
	"os"
)

const BufLength = 128

type Server struct {
	Net           string
	IpPort        string
	Spec          int
	HeartbeatName string
	send          chan []byte
	conn          net.Conn
	read_buf      map[uint16]func(b []byte)
}

func (s *Server) InitServer(c chan []byte) {
	addr, err := net.ResolveTCPAddr(s.Net, s.IpPort)
	if err != nil {
		log.Fatalln(err)
	}

	conn, err := net.DialTCP("tcp", nil, addr)

	if err != nil {
		log.Fatalln(err.Error())
	}

	s.send = c
	s.conn = conn
	s.read_buf = make(map[uint16]func(b []byte), 1024)

	//发送心跳
	var pd parsers.PduHeader
	if s.HeartbeatName == "" || s.Spec == 0 {
		return
	}
	ht, err := heartbeat.NewTast(s.HeartbeatName, s.Spec)
	if err != nil {
		log.Println(err)
		return
	}
	// Run a new mission
	ht.Start(func() error {
		data, _ := pd.HeartbeatPacket()
		s.conn.Write(data)
		return nil
	})

	go s.HandleRece()

	go s.HandleSend()
}

func (s *Server) HandleRece() {
	for {
		data := make([]byte, 0)
		buf := make([]byte, BufLength)
		for {
			n, err := s.conn.Read(buf)
			if err != nil && err != io.EOF {
				checkErr(err)
			}
			data = append(data, buf[:n]...)
			if n != BufLength {
				break
			}
		}
		models.DBChanRece <- data
	}
	s.conn.Close()
}

func (s *Server) HandleSend() {
	for msg := range s.send {
		if _, err := s.conn.Write(msg); err != nil {
			log.Println(err)
		}
	}
	s.conn.Close()
}

type DBServer struct {
	Server
}

type LoginServer struct {
	Server
	MaxConnCnt string
	CurConnCnt string
	HostName   string
}

type RouteServer struct {
	Server
}

func checkErr(err error) {
	if err != nil {
		log.Println(err)
		os.Exit(-1)
	}
}
