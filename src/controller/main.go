package main

import(
  "net"
	"fmt"
	"os"
	"time"
	"strings"
)

func main(){
	laddr,err := net.ResolveUDPAddr("udp4",":12346");if err != nil {return}
	raddr,err := net.ResolveUDPAddr("udp4","127.0.0.1:12345");if err != nil {return}
	conn,err := net.DialUDP("udp4",laddr,raddr);if err != nil {return}
	buf := []byte(strings.Join(os.Args[1:]," "));
	n,err := conn.Write(buf);if (err != nil)||(n != len(buf)){fmt.Printf("ERROR: %d - %s\n",n,err);return}
	fmt.Printf("OK: %d\n",n);
	conn.SetReadDeadline(time.Now().Add(time.Second * 5));
	recbuf := make([]byte,4*1024);
	n,addr,err := conn.ReadFromUDP(recbuf);
	fmt.Printf("%v : %v : %v : %s", n,addr,err,string(recbuf[0:n]));
}	;
