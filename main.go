package main

import (
	"encoding/binary"
	"fmt"
	"math"
	"net"
	"os"
	"strings"
	"time"
)

const proto, addr = "udp", "127.0.0.1:6454"

const sleepRate = 28 * time.Millisecond

const dmxMax = 255

var reds []int
var greens []int
var blues []int

func main() {

	// set base dmx data to all 0s
	dmxdata := []byte("\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00")
	var dmxdataP = &dmxdata

	//rgb offsets
	greens = append(greens, 1, 4, 7, 10, 13, 16, 19, 22, 25, 28)
	reds = append(reds, 2, 5, 8, 11, 14, 17, 20, 23, 26, 29)
	blues = append(blues, 3, 6, 9, 12, 15, 18, 21, 24, 27, 30)
	L1 := []int{4, 5, 6, 28, 29, 30}
	L2 := []int{10, 11, 12, 26, 27, 25}
	L3 := []int{13, 14, 15, 22, 23, 24}
	L4 := []int{16, 17, 18, 19, 20, 21}

	// other init const
	var dmxPerc float32 = 1.00
	dmxPercP := &dmxPerc
	var dmxRedPerc float32 = 0.75
	dmxRedPercP := &dmxRedPerc
	var dmxGreenPerc float32 = 0.25
	dmxGreebPercP := &dmxGreenPerc
	var dmxBluePerc float32 = 0.00
	dmxBluePercP := &dmxBluePerc
	toggleFollow := false

	// setup udp listener, close when done
	udpAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:2347")
	check(err)
	ua, err := net.ListenUDP("udp", udpAddr)
	check(err)
	defer ua.Close()

	//start stream of artnet packets
	go sendArtPacket(dmxdata)

	fmt.Printf("service started")

	for {

		buffer := make([]byte, 1024)
		n, _, err := ua.ReadFromUDP(buffer)
		check(err)

		/* get id */
		id, val := parseOSC(buffer[:n])
		trimmed := strings.Trim(id, "\x00")

		idMap := make(map[string]bool)
		idMap[trimmed] = true

		// im sorry
		if idMap["brightness"] {
			dmxPerc = val
		} else if idMap["toggleFollow"] {
			if toggleFollow {
				toggleFollow = false
				go resetBrightness(dmxdataP)
			} else {
				toggleFollow = true
			}
		} else if idMap["follow"] {
			if toggleFollow {
				go follow(dmxdataP, dmxMax, dmxPercP, val, dmxRedPercP, dmxGreebPercP, dmxBluePercP)
			}
		} else if idMap["red"] {
			go blip(dmxdataP, reds, 1.00, dmxPercP)
		} else if idMap["red25"] {
			go blip(dmxdataP, reds, .45, dmxPercP)
		} else if idMap["green"] {
			go blip(dmxdataP, greens, 1.00, dmxPercP)
		} else if idMap["green25"] {
			go blip(dmxdataP, greens, .45, dmxPercP)
		} else if idMap["blue"] {
			go blip(dmxdataP, blues, 1.00, dmxPercP)
		} else if idMap["blue25"] {
			go blip(dmxdataP, blues, .45, dmxPercP)
		} else if idMap["scatter"] {
			go scatter(dmxdataP, dmxMax, dmxPercP)
		} else if idMap["redWheel"] {
			dmxRedPerc = val
		} else if idMap["greenWheel"] {
			dmxGreenPerc = val
		} else if idMap["blueWheel"] {
			dmxBluePerc = val
		} else if idMap["L1"] {
			go blipWithColorWheels(dmxdataP, L1, 1.00, dmxPercP, dmxRedPercP, dmxGreebPercP, dmxBluePercP)
		} else if idMap["L2"] {
			go blipWithColorWheels(dmxdataP, L2, 1.00, dmxPercP, dmxRedPercP, dmxGreebPercP, dmxBluePercP)
		} else if idMap["L3"] {
			go blipWithColorWheels(dmxdataP, L3, 1.00, dmxPercP, dmxRedPercP, dmxGreebPercP, dmxBluePercP)
		} else if idMap["L4"] {
			go blipWithColorWheels(dmxdataP, L4, 1.00, dmxPercP, dmxRedPercP, dmxGreebPercP, dmxBluePercP)
		}
	}

}

func follow(dmxdataP *[]byte, dmxMax float32, dmxPercP *float32, val float32, dmxRedPercP *float32, dmxGreenPercP *float32, dmxBluePercP *float32) {
	val = val * 10
	for _, id := range reds {
		brightness := val * (dmxMax * *dmxRedPercP)
		(*dmxdataP)[id-1] = byte(brightness)
	}
	for _, id := range greens {
		brightness := val * (dmxMax * *dmxGreenPercP)
		(*dmxdataP)[id-1] = byte(brightness)
	}
	for _, id := range blues {
		brightness := val * (dmxMax * *dmxBluePercP)
		(*dmxdataP)[id-1] = byte(brightness)
	}
}

func resetBrightness(dmxdataP *[]byte) {
	for i := 1; i <= 30; i++ {
		(*dmxdataP)[i] = byte(0)
	}
}

func scatter(dmxdataP *[]byte, val int, dmxPercP *float32) {
	for i := 1; i <= 15; i++ {
		time.Sleep(15 * time.Millisecond)
		brightness := dmxMax * *dmxPercP
		(*dmxdataP)[i] = byte(brightness)
		(*dmxdataP)[30-i] = byte(brightness)
	}
	time.Sleep(50 * time.Millisecond)
	for i := 1; i <= 30; i++ {
		(*dmxdataP)[i] = byte(0)
	}

}

func blip(dmxdataP *[]byte, ids []int, volume float32, dmxPercP *float32) {
	for i := 255 / 10; i >= 1; i-- {
		time.Sleep(10 * time.Millisecond)
		brightness := float32(i*10) * *dmxPercP * volume
		for _, id := range ids {
			(*dmxdataP)[id-1] = byte(brightness)
		}
	}
	for _, id := range ids {
		(*dmxdataP)[id-1] = byte(0)
	}
}

func blipWithColorWheels(dmxdataP *[]byte, ids []int, volume float32, dmxPercP *float32, dmxRedPercP *float32, dmxGreenPercP *float32, dmxBluePercP *float32) {

	for i := 255 / 10; i >= 1; i-- {
		time.Sleep(10 * time.Millisecond)
		stdVal := float32(i * 10)
		redVal := stdVal * (dmxMax * *dmxRedPercP) * volume
		greenVal := stdVal * (dmxMax * *dmxGreenPercP) * volume
		blueVal := stdVal * (dmxMax * *dmxBluePercP) * volume

		// need to put this into a struct for faster
		// checking rather than itr through all these
		for _, red := range reds {
			for _, id := range ids {
				if red == id {
					(*dmxdataP)[id-1] = byte(redVal)
				}
			}
		}
		for _, green := range greens {
			for _, id := range ids {
				if green == id {
					(*dmxdataP)[id-1] = byte(greenVal)
				}
			}
		}
		for _, blue := range blues {
			for _, id := range ids {
				if blue == id {
					(*dmxdataP)[id-1] = byte(blueVal)
				}
			}
		}

	}

	for i := 1; i <= 30; i++ {
		(*dmxdataP)[i-1] = byte(0)
	}
}

func parseOSC(packet []byte) (id string, value float32) {

	// grab id before data
	s := string(packet[1:])
	i := strings.Split(s, ",f")[0]

	// grab last 4 packets that have the data
	bits := binary.BigEndian.Uint32(packet[len(string(packet))-4:])
	float := math.Float32frombits(bits)

	return i, float
}

func sendArtPacket(data []byte) {

	header := []byte("Art-Net\x00\x00P\x00\x0e\x00\x00\x00\x00\x02\x00")

	// i := 1
	// for _, v := range header {
	// 	println(i, v)
	// 	i++
	// }

	udpServer, err := net.ResolveUDPAddr(proto, addr)
	check(err)

	conn, err := net.DialUDP(proto, nil, udpServer)
	check(err)
	defer conn.Close()

	for {
		conn.Write(append(header[:], data[:]...))
		time.Sleep(sleepRate)
	}
}

func check(err error) {
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
