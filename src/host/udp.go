package main

import(
	"net"
	"os"
	"fmt"
)

type TMessage struct{
	proto uint8
	length uint32
	data []byte
}

func processBlock(conn *net.UDPConn, from net.Addr, data []byte){
	var Result string;
	fmt.Printf("[%v:%v] %v\n",conn,from,data);
	r,e := Execute(string(data));
	if e != nil {
		Result = fmt.Sprintf("ERROR: %s",e.Error())
	}else{
		Result = r
	}
	fmt.Printf("%s\n",Result);
	conn.WriteTo([]byte(Result),from);
}

var ContinueUDPFlag bool = true;

func udp(){
	addr,err := net.ResolveUDPAddr("udp4",":12345");if err != nil {return}
	conn,err := net.ListenUDP("udp4",addr);if err != nil {return}
	var buf []byte = make([]byte,4*1024);
	for ContinueUDPFlag {
		n,addr,err := conn.ReadFromUDP(buf);
		if err != nil {
			os.Stdout.Write([]byte("ERROR: " + err.Error() + "\n"));
		}else{
			os.Stdout.Write([]byte(fmt.Sprintf("(%v)[%d]: %s \n",addr,n,buf[0:n])));
			go processBlock(conn,addr,buf[0:n]);
		}
	}
}

/*
Assets

Known networks are "tcp", "tcp4" (IPv4-only), "tcp6" (IPv6-only), "udp", "udp4" (IPv4-only), "udp6" (IPv6-only), "ip", "ip4" (IPv4-only), "ip6" (IPv6-only), "unix", "unixgram" and "unixpacket". 

*/