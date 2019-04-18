package utils

import (
	"bufio"
	"bytes"
	"fmt"
	"strings"
	"time"

	"github.com/lflxp/showme/utils/decoder/http"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

type Data struct {
	SrcMac     string            `json:"srcmac"`
	DstMac     string            `json:"dstmac"`
	SrcIp      string            `json:"srcip"`
	DstIp      string            `json:"Dstip"`
	Protocol   string            `json:"protocol"`
	SrcPort    string            `json:"srcport"`
	DstPort    string            `json:"dstport"`
	Type       string            `json:"type"`
	Method     string            `json:"method"`
	Url        string            `json:"url"`
	Version    string            `json:"version"`
	Headers    map[string]string `json:"headers"`
	StatusCode int               `json:"statusCode"`
	StatusMsg  string            `json:"statusMsg"`
	Time       string            `json:"time"`
	Response   string            `json:"response"`
}

type Device struct {
	Device       string
	Snapshot_len int32
	Promiscuous  bool
	Timeout      time.Duration
	Handle       *pcap.Handle
}

func (this *Device) init() error {
	handle, err := pcap.OpenLive(this.Device, this.Snapshot_len, this.Promiscuous, this.Timeout)
	if err != nil {
		return err
	}
	this.Handle = handle
	return nil
}

func NewDevice(deviceName string) (*Device, error) {
	info := &Device{
		Device:       deviceName,
		Snapshot_len: 1024,
		Promiscuous:  false,
		Timeout:      30 * time.Second,
	}
	err := info.init()
	return info, err
}

func WatchDogEasy(name string) {
	// fmt.Println("watchDog", name)
	handle, err := NewDevice(name)
	if err != nil {
		fmt.Println("w eerr", err.Error())
	}
	defer handle.Handle.Close()

	// Use the handle as a packet source to process all packets
	packetSource := gopacket.NewPacketSource(handle.Handle, handle.Handle.LinkType())
	for p := range packetSource.Packets() {
		// Process packet here
		// fmt.Println(p)
		fmt.Println(p.String())
	}
}

func WatchDogString(datainfo chan interface{}, name string) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Safe quit WatchDogString")
		}
	}()
	// fmt.Println("watchDog", name)
	handle, err := NewDevice(name)
	if err != nil {
		fmt.Println("w eerr", err.Error())
	}
	defer handle.Handle.Close()

	// Use the handle as a packet source to process all packets
	packetSource := gopacket.NewPacketSource(handle.Handle, handle.Handle.LinkType())
	for p := range packetSource.Packets() {
		// Process packet here
		// fmt.Println(p)
		// fmt.Println(p.String())
		datainfo <- p.String()
	}
}

func WatchDog(datainfo chan interface{}, name string) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Safe quit WatchDog", r)
		}
	}()
	// fmt.Println("watchDog", name)
	handle, err := NewDevice(name)
	if err != nil {
		fmt.Println("w eerr", err.Error())
	}
	defer handle.Handle.Close()
	// Use the handle as a packet source to process all packets
	packetSource := gopacket.NewPacketSource(handle.Handle, handle.Handle.LinkType())
	for p := range packetSource.Packets() {
		// Process packet here
		// fmt.Println(p)
		// fmt.Println(p.String())
		info := &Data{}

		eth := p.Layer(layers.LayerTypeEthernet)
		if eth != nil {
			// fmt.Println("Ethernet layer detected.")
			ethernetPacket, _ := eth.(*layers.Ethernet)
			// fmt.Println("Source MAC:", ethernetPacket.SrcMAC)
			// fmt.Println("Destionation MAC:", ethernetPacket.DstMAC)
			// // Ethernet type is typically IPv4 but could be ARP or other
			// fmt.Println("Ethernet type: ", ethernetPacket.EthernetType)
			// fmt.Println()
			info.SrcMac = ethernetPacket.SrcMAC.String()
			info.DstMac = ethernetPacket.DstMAC.String()
		}

		// Dq := p.Layer(layers.LayerTypeDot1Q)
		// if Dq != nil {
		// 	fmt.Println("LayerTypeDot1Q layer detected.")
		// 	dd, _ := Dq.(*layers.Dot1Q)
		// 	fmt.Println("DropEligible:", dd.DropEligible)
		// 	fmt.Println("Priority:", dd.Priority)
		// 	// Ethernet type is typically IPv4 but could be ARP or other
		// 	fmt.Println("Dot1Q type: ", dd.Type)
		// 	fmt.Println("VLANIdentifier: ", dd.VLANIdentifier)
		// 	fmt.Println()
		// }

		// icmq := p.Layer(layers.LayerTypeICMPv4)
		// if icmq != nil {
		// 	fmt.Println("LayerTypeICMPv4 layer detected.")
		// 	ic, _ := icmq.(*layers.ICMPv4)
		// 	fmt.Println("Checksum:", ic.Checksum)
		// 	fmt.Println("Id:", ic.Id)
		// 	// Ethernet type is typically IPv4 but could be ARP or other
		// 	fmt.Println("Seq: ", ic.Seq)
		// 	fmt.Println("TypeCode: ", ic.TypeCode.String())
		// 	fmt.Println()
		// }

		// // Let's see if the packet is IP ∂(even though the ether type told us)
		ipLayer := p.Layer(layers.LayerTypeIPv4)
		if ipLayer != nil {
			// fmt.Println("IPv4 layer detected.")
			ip, _ := ipLayer.(*layers.IPv4)

			// IP layer variables:
			// Version (Either 4 or 6)
			// IHL (IP Header Length in 32-bit words)
			// TOS, Length, Id, Flags, FragOffset, TTL, Protocol (TCP?),
			// Checksum, SrcIP, DstIP
			// fmt.Printf("From %s to %s\n", ip.SrcIP, ip.DstIP)
			// fmt.Println("Protocol: ", ip.Protocol)
			// fmt.Println()
			info.SrcIp = ip.SrcIP.String()
			info.DstIp = ip.DstIP.String()
			info.Protocol = ip.Protocol.String()
		}

		tcpLayer := p.Layer(layers.LayerTypeTCP)
		if tcpLayer != nil {
			tcp, _ := tcpLayer.(*layers.TCP)

			info.SrcPort = tcp.SrcPort.String()
			info.DstPort = tcp.DstPort.String()
		}

		// Let's see if the packet is TCP
		udpLayer := p.Layer(layers.LayerTypeUDP)
		if udpLayer != nil {
			// fmt.Println("UDP layer detected.")
			udp, _ := udpLayer.(*layers.UDP)

			// TCP layer variables:
			// SrcPort, DstPort, Seq, Ack, DataOffset, Window, Checksum, Urgent
			// Bool flags: FIN, SYN, RST, PSH, ACK, URG, ECE, CWR, NS
			// fmt.Printf("From port %d to %d\n", udp.SrcPort, udp.DstPort)
			// fmt.Println("Checksum number: ", udp.Checksum)
			// fmt.Println()

			// pp := gopacket.NewPacket(udp.Payload, layers.LayerTypeSFlow, gopacket.Default)
			// if pp.ErrorLayer() == nil {
			// 	fmt.Println("UDP SFLOW detected")
			// }
			info.SrcPort = udp.SrcPort.String()
			info.DstPort = udp.DstPort.String()
		}
		// sflowlayer := p.Layer(layers.LayerTypeSFlow)
		// if sflowlayer != nil {
		// 	fmt.Println("SFLOW layer detected")
		// }

		// When iterating through packet.Layers() above,
		// if it lists Payload layer then that is the same as
		// this applicationLayer. applicationLayer contains the payload
		applicationLayer := p.ApplicationLayer()
		if applicationLayer != nil {
			// fmt.Println("Application layer/Payload found.")
			// fmt.Printf("%d %s\n", len(applicationLayer.Payload()), applicationLayer.Payload())
			// Search for a string inside the payload
			if strings.Contains(string(applicationLayer.Payload()), "HTTP") {
				// fmt.Println("HTTP found!")
				// fmt.Println(applicationLayer.LayerType().String())
				// fmt.Println(p.Dump())
				decoder := http.Decoder{}
				decoder.SetFilter("")
				decoder.Buf = bufio.NewReader(bytes.NewReader(applicationLayer.Payload()))
				_data, err := decoder.DecodeHttp()
				if err != nil {
					fmt.Println(err)
				}

				// req := _data.(*http.HttpReq)
				// if req != nil {
				// 	fmt.Println("Request", req)
				// }

				// resp := _data.(*http.HttpResp)
				// if resp != nil {
				// 	fmt.Println("Response", resp)
				// }
				info.Time = time.Now().Format("2006-01-02 15:04:05")
				switch _data.(type) {
				case *http.HttpReq:
					// fmt.Println("Request", _data.(*http.HttpReq), string(_data.(*http.HttpReq).RawBody()))
					tmp := _data.(*http.HttpReq)
					info.Type = "request"
					info.Method = tmp.Method
					info.Url = tmp.Url
					info.Version = tmp.Version
					info.Headers = tmp.Headers
					info.Response = string(tmp.RawBody())
				case *http.HttpResp:
					// fmt.Println("Response", _data.(*http.HttpResp), string(_data.(*http.HttpResp).RawBody()))
					tmp := _data.(*http.HttpResp)
					info.Type = "response"
					info.Version = tmp.Version
					info.StatusCode = tmp.StatusCode
					info.StatusMsg = tmp.StatusMsg
					info.Headers = tmp.Headers
					info.Response = string(tmp.RawBody())
				}

			}
			// fmt.Println(applicationLayer.Payload())
			// fmt.Println()
		}

		datainfo <- info
		// _, ok := <-datainfo
		// if ok {
		// 	datainfo <- info
		// } else {
		// 	return
		// }

		// Check for errors
		// if err := p.ErrorLayer(); err != nil {
		// 	fmt.Println("Error decoding some part of the packet:", err)
		// }

		// for _, x := range p.Layers() {
		// 	fmt.Println(x.LayerType().String())
		// }

		// fmt.Println()
	}
}

func WatchDogData(limit int, name string) []*Data {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in f", r)
		}
	}()
	handle, err := NewDevice(name)
	if err != nil {
		fmt.Println("w eerr", err.Error())
	}
	defer handle.Handle.Close()

	num := 0
	rs := []*Data{}

	// Use the handle as a packet source to process all packets
	packetSource := gopacket.NewPacketSource(handle.Handle, handle.Handle.LinkType())
	for p := range packetSource.Packets() {
		if num > limit {
			break
			return rs
		}
		info := &Data{}

		eth := p.Layer(layers.LayerTypeEthernet)
		if eth != nil {
			ethernetPacket, _ := eth.(*layers.Ethernet)
			info.SrcMac = ethernetPacket.SrcMAC.String()
			info.DstMac = ethernetPacket.DstMAC.String()
		}

		// // Let's see if the packet is IP ∂(even though the ether type told us)
		ipLayer := p.Layer(layers.LayerTypeIPv4)
		if ipLayer != nil {
			// fmt.Println("IPv4 layer detected.")
			ip, _ := ipLayer.(*layers.IPv4)

			info.SrcIp = ip.SrcIP.String()
			info.DstIp = ip.DstIP.String()
			info.Protocol = ip.Protocol.String()
		}

		tcpLayer := p.Layer(layers.LayerTypeTCP)
		if tcpLayer != nil {
			tcp, _ := tcpLayer.(*layers.TCP)

			info.SrcPort = tcp.SrcPort.String()
			info.DstPort = tcp.DstPort.String()
		}

		// Let's see if the packet is TCP
		udpLayer := p.Layer(layers.LayerTypeUDP)
		if udpLayer != nil {
			// fmt.Println("UDP layer detected.")
			udp, _ := udpLayer.(*layers.UDP)

			info.SrcPort = udp.SrcPort.String()
			info.DstPort = udp.DstPort.String()
		}
		// sflowlayer := p.Layer(layers.LayerTypeSFlow)
		// if sflowlayer != nil {
		// 	fmt.Println("SFLOW layer detected")
		// }

		// When iterating through packet.Layers() above,
		// if it lists Payload layer then that is the same as
		// this applicationLayer. applicationLayer contains the payload
		applicationLayer := p.ApplicationLayer()
		if applicationLayer != nil {
			// fmt.Println("Application layer/Payload found.")
			// fmt.Printf("%d %s\n", len(applicationLayer.Payload()), applicationLayer.Payload())
			// Search for a string inside the payload
			if strings.Contains(string(applicationLayer.Payload()), "HTTP") {
				// fmt.Println("HTTP found!")
				// fmt.Println(applicationLayer.LayerType().String())
				// fmt.Println(p.Dump())
				decoder := http.Decoder{}
				decoder.SetFilter("")
				decoder.Buf = bufio.NewReader(bytes.NewReader(applicationLayer.Payload()))
				_data, err := decoder.DecodeHttp()
				if err != nil {
					fmt.Println(err)
				}

				// req := _data.(*http.HttpReq)
				// if req != nil {
				// 	fmt.Println("Request", req)
				// }

				// resp := _data.(*http.HttpResp)
				// if resp != nil {
				// 	fmt.Println("Response", resp)
				// }
				info.Time = time.Now().Format("2006-01-02 15:04:05")
				switch _data.(type) {
				case *http.HttpReq:
					// fmt.Println("Request", _data.(*http.HttpReq), string(_data.(*http.HttpReq).RawBody()))
					tmp := _data.(*http.HttpReq)
					info.Type = "request"
					info.Method = tmp.Method
					info.Url = tmp.Url
					info.Version = tmp.Version
					info.Headers = tmp.Headers
					info.Response = string(tmp.RawBody())
				case *http.HttpResp:
					// fmt.Println("Response", _data.(*http.HttpResp), string(_data.(*http.HttpResp).RawBody()))
					tmp := _data.(*http.HttpResp)
					info.Type = "response"
					info.Version = tmp.Version
					info.StatusCode = tmp.StatusCode
					info.StatusMsg = tmp.StatusMsg
					info.Headers = tmp.Headers
					info.Response = string(tmp.RawBody())
				}
			}
		}
		num++
		rs = append(rs, info)
	}
	return rs
}
