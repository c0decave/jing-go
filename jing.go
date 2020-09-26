package main

import (
	"bufio"
	"fmt"

	//"io/ioutil"
	"encoding/json"
	"log"
	"os"

	"github.com/go-ping/ping"
	"github.com/pborman/getopt/v2"
)

var (
	target      = "127.0.0.1"
	target_list = ""
	sleep       = 100
	privilege   = false
	count       = 3
	timeout     = 3000000000
	json_out    = false
	target_str  = ""
)

func init() {

	getopt.FlagLong(&target, "target", 't', "the target host to ping")
	getopt.FlagLong(&target_list, "target_list", 'T', "define targetlist with hosts to ping")
	getopt.FlagLong(&sleep, "sleep", 's', "sleep between packets")
	getopt.FlagLong(&count, "count", 'c', "how many packets to send")
	getopt.FlagLong(&privilege, "privilege", 'r', "admin/root privileges available")
	getopt.FlagLong(&json_out, "json-output", 'J', "prefer json output")
	//	getopt.FlagLong(&timeout,"")
}

type ping_args struct {
	Target      string
	Target_list string
	Timeout     int
	Count       int
	Privilege   bool
}

type ping_result struct {
	Target string
	Sent   int
	Recv   int
	Alive  bool
}

func pingme(target string, count int, timeout int) {

	var alive bool
	pinger, err := ping.NewPinger(target)

	if err != nil {
		panic(err)

	}
	// fix this with time import and milliseconds
	pinger.Timeout = 3000000000
	pinger.Count = 3

	pinger.OnFinish = func(stats *ping.Statistics) {
		if stats.PacketsRecv > 0 {
			alive = true
		} else {
			alive = false
		}

		if json_out == false {
			fmt.Printf("%s (%s) %d/%d (Sent/Recv) Alive: %t\n", stats.Addr, stats.IPAddr, stats.PacketsSent, stats.PacketsRecv, alive)
		} else {

			p_result := ping_result{

				Target: target,
				Sent:   stats.PacketsSent,
				Recv:   stats.PacketsRecv,
				Alive:  alive,
			}
			json_result, err := json.Marshal(p_result)
			if err != nil {
				log.Fatalf("Unable to encode")
			}
			fmt.Println(string(json_result))
		}
	}

	//pinger.SetPrivileged(true)
	err = pinger.Run()
	if err != nil {

		fmt.Printf("Failed to ping target host :%s", err)
	}

}

func main() {

	getopt.Parse()
	p_args := ping_args{
		Target:      target,
		Target_list: target_list,
		Timeout:     timeout,
		Count:       count,
		Privilege:   privilege}

	var json_data []byte

	json_data, err := json.Marshal(p_args)
	if err != nil {
		log.Println(err)
	}
	if json_out == true {
		fmt.Println(string(json_data))
	}

	if target_list != "" {

		buf, err := os.Open(target_list)
		if err != nil {
			log.Fatal(err)
		}

		snl := bufio.NewScanner(buf)
		for snl.Scan() {
			pingme(snl.Text(), 3, 3000000000)

		}
	} else {
		pingme(target, 3, 3000000000)
	}
	if json_out == false {
		fmt.Printf("Done\n")
	}
}
