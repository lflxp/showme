package pkg

import (
	"encoding/json"
	"fmt"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"

	//"github.com/Cistern/sflow"
	"net"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

type Collected struct {
	DeviceName      string //设备名称
	SnapShotLen     int32
	SnapShotLenUint uint32
	Promiscuous     bool //是否开启混杂模式
	Timeout         time.Duration
	Udpbool         bool   //是否开启udp sample and netflow传输
	Host            string //udp 发送客户端及端口 127.0.0.1:8888
	CounterHost     string //udp counter 传输
	EsPath          string // elasticsearch address path
	IsEs            bool   // 是否传送到es
	Index           string // es索引名称
}

func (this *Collected) SendUdp(result string, counter bool) {
	if counter {
		conn, err := net.Dial("udp", this.CounterHost)
		defer conn.Close()
		if err != nil {
			log.Error(err.Error())
		}
		conn.Write([]byte(result))
	} else {
		conn, err := net.Dial("udp", this.Host)
		defer conn.Close()
		if err != nil {
			log.Error(err.Error())
		}
		conn.Write([]byte(result))
	}
}

func (this *Collected) CheckInfo(ppp []byte) {
	p := gopacket.NewPacket(ppp, layers.LayerTypeEthernet, gopacket.Default)
	if p.ErrorLayer() != nil {
		log.Error("failed :", p.ErrorLayer().Error())
	}
	eth := p.Layer(layers.LayerTypeEthernet)
	if eth != nil {
		log.Debug("Ethernet layer detected.")
		ethernetPacket, _ := eth.(*layers.Ethernet)
		log.Debug("Source MAC:", ethernetPacket.SrcMAC)
		log.Debug("Destionation MAC:", ethernetPacket.DstMAC)
		// Ethernet type is typically IPv4 but could be ARP or other
		log.Debug("Ethernet type: ", ethernetPacket.EthernetType)
		log.Debugln()
	}

	Dq := p.Layer(layers.LayerTypeDot1Q)
	if Dq != nil {
		log.Debug("LayerTypeDot1Q layer detected.")
		dd, _ := Dq.(*layers.Dot1Q)
		log.Debug("DropEligible:", dd.DropEligible)
		log.Debug("Priority:", dd.Priority)
		// Ethernet type is typically IPv4 but could be ARP or other
		log.Debug("Dot1Q type: ", dd.Type)
		log.Debug("VLANIdentifier: ", dd.VLANIdentifier)
		log.Debugln()
	}

	icmq := p.Layer(layers.LayerTypeICMPv4)
	if icmq != nil {
		log.Debug("LayerTypeICMPv4 layer detected.")
		ic, _ := icmq.(*layers.ICMPv4)
		log.Debug("Checksum:", ic.Checksum)
		log.Debug("Id:", ic.Id)
		// Ethernet type is typically IPv4 but could be ARP or other
		log.Debug("Seq: ", ic.Seq)
		log.Debug("TypeCode: ", ic.TypeCode.String())
		log.Debugln()
	}

	// Let's see if the packet is IP ∂(even though the ether type told us)
	ipLayer := p.Layer(layers.LayerTypeIPv4)
	if ipLayer != nil {
		log.Debug("IPv4 layer detected.")
		ip, _ := ipLayer.(*layers.IPv4)

		// IP layer variables:
		// Version (Either 4 or 6)
		// IHL (IP Header Length in 32-bit words)
		// TOS, Length, Id, Flags, FragOffset, TTL, Protocol (TCP?),
		// Checksum, SrcIP, DstIP
		log.Debugf("From %s to %s\n", ip.SrcIP, ip.DstIP)
		log.Debug("Protocol: ", ip.Protocol)
		log.Debugln()
	}

	// Let's see if the packet is TCP
	udpLayer := p.Layer(layers.LayerTypeUDP)
	if udpLayer != nil {
		log.Debug("UDP layer detected.")
		udp, _ := udpLayer.(*layers.UDP)

		// TCP layer variables:
		// SrcPort, DstPort, Seq, Ack, DataOffset, Window, Checksum, Urgent
		// Bool flags: FIN, SYN, RST, PSH, ACK, URG, ECE, CWR, NS
		log.Debugf("From port %d to %d\n", udp.SrcPort, udp.DstPort)
		log.Debugln("Checksum number: ", udp.Checksum)
		log.Debugln()

		pp := gopacket.NewPacket(udp.Payload, layers.LayerTypeSFlow, gopacket.Default)
		if pp.ErrorLayer() == nil {
			log.Errorln("UDP SFLOW detected")
		}
	}
	sflowlayer := p.Layer(layers.LayerTypeSFlow)
	if sflowlayer != nil {
		log.Debugln("SFLOW layer detected")
	}

	// When iterating through packet.Layers() above,
	// if it lists Payload layer then that is the same as
	// this applicationLayer. applicationLayer contains the payload
	applicationLayer := p.ApplicationLayer()
	if applicationLayer != nil {
		log.Debugln("Application layer/Payload found.")
		log.Debugf("%d %s\n", len(applicationLayer.Payload()), applicationLayer.Payload())
		// Search for a string inside the payload
		if strings.Contains(string(applicationLayer.Payload()), "HTTP") {
			log.Debugln("HTTP found!")
		}
		log.Debugln(applicationLayer.Payload())
		log.Debugln()
	}

	// Check for errors
	if err := p.ErrorLayer(); err != nil {
		log.Debugln("Error decoding some part of the packet:", err)
	}

	for _, x := range p.Layers() {
		log.Infoln(x.LayerType().String())
	}

	log.Debugln()
}

func (this *Collected) ListenSFlowSample(protocol, port string) {
	//Open Device
	handle, err := pcap.OpenLive(this.DeviceName, this.SnapShotLen, this.Promiscuous, this.Timeout)
	if err != nil {
		log.Error(err.Error())
		return
	}
	defer handle.Close()

	//Set filter
	var filter string = fmt.Sprintf("%s and port %s", protocol, port)
	err = handle.SetBPFFilter(filter)
	if err != nil {
		log.Error(err)
	}
	log.Info(fmt.Sprintf("Only capturing %s port %s packets.", protocol, port))
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	for packet := range packetSource.Packets() {
		// this.CheckInfo(packet.Data())
		Origin := NewData()
		p := gopacket.NewPacket(packet.Data(), layers.LayerTypeEthernet, gopacket.Default)
		if p.ErrorLayer() != nil {
			log.Errorln("failed LayerTypeEthernet:", p.ErrorLayer().Error())
		}

		Origin.Init(p)

		udpLayer := p.Layer(layers.LayerTypeUDP)
		if udpLayer != nil {
			// fmt.Println("UDP layer detected.")
			udp, _ := udpLayer.(*layers.UDP)
			pp := gopacket.NewPacket(udp.Payload, layers.LayerTypeSFlow, gopacket.Default)
			if pp.ErrorLayer() != nil {
				log.Errorln("failed LayerTypeUDP:", p.ErrorLayer().Error())
			}

			if got, ok := pp.ApplicationLayer().(*layers.SFlowDatagram); ok {
				go func(datas []layers.SFlowFlowSample, got *layers.SFlowDatagram) {
					for _, y := range datas {
						//beego.Critical(len(y.Records),y.RecordCount)
						tmp := NewFlowSamples()
						//tmp.InitOriginData(p)
						tmp.Data = Origin
						tmp.InitFlowSampleData(y)

						// fmt.Println(p.Dump())
						// b, err := json.Marshal(tmp)
						// if err != nil {
						// 	fmt.Println(err.Error())
						// }
						if this.Udpbool {
							b, err := json.Marshal(tmp)
							if err != nil {
								log.Error(err.Error())
							}
							this.SendUdp(string(b), false)
						} else if this.IsEs {
							result, err := ParseSflowV5ToEs(tmp, nil)
							if err != nil {
								log.Errorln("err", err.Error())
							} else {
								DataChannel <- result
								// err = parse.CreateEs(result, "doc", fmt.Sprintf("ABC%d", time.Now().UnixNano()))
								// if err != nil {
								// 	log.Errorln("send error:", err.Error(), result)
								// }
							}
						} else {
							b, err := json.Marshal(tmp)
							if err != nil {
								log.Errorln(err.Error())
							}
							log.Debugln(string(b))
						}
					}
				}(got.FlowSamples, got)
			}
		}
	}
}

func (this *Collected) ListenSflowCounter(protocol, port string) {
	//Open Device
	handle, err := pcap.OpenLive(this.DeviceName, this.SnapShotLen, this.Promiscuous, this.Timeout)
	if err != nil {
		log.Error(err)
		return
	}
	defer handle.Close()

	//Set filter
	var filter string = fmt.Sprintf("%s and port %s", protocol, port)
	err = handle.SetBPFFilter(filter)
	if err != nil {
		log.Error(err)
	}
	log.Info(fmt.Sprintf("Only capturing %s port %s packets.", protocol, port))
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	for packet := range packetSource.Packets() {
		//fmt.Println(packet.Dump())
		p := gopacket.NewPacket(packet.Data(), layers.LayerTypeEthernet, gopacket.Default)
		if p.ErrorLayer() != nil {
			log.Errorln("failed LayerTypeEthernet:", p.ErrorLayer().Error())
		}

		udpLayer := packet.Layer(layers.LayerTypeUDP)
		if udpLayer != nil {
			//fmt.Println("UDP layer detected.")
			udp, _ := udpLayer.(*layers.UDP)

			pp := gopacket.NewPacket(udp.Payload, layers.LayerTypeSFlow, gopacket.Default)
			if pp.ErrorLayer() != nil {
				log.Errorln("failed LayerTypeUDP:", p.ErrorLayer().Error())
			}
			if got, ok := pp.ApplicationLayer().(*layers.SFlowDatagram); ok {
				go func(datas []layers.SFlowCounterSample) {
					if len(datas) > 0 {
						//log.Error(udp.Payload)
						tmp := NewCounterFlow()
						tmp.InitOriginData(p)
						for _, y := range datas {
							//beego.Critical(len(y.Records),y.RecordCount)
							tmp.InitCounterSample(y)
						}

						if this.Udpbool {
							b, err := json.Marshal(tmp)
							if err != nil {
								log.Errorln(err.Error())
							}
							this.SendUdp(string(b), true)
						} else if this.IsEs {
							result, err := ParseSflowV5ToEs(nil, tmp)
							if err != nil {
								log.Errorln("err", err.Error())
							} else {
								DataChannel <- result
								// err = parse.CreateEs(result, "doc", fmt.Sprintf("ABC%d", time.Now().UnixNano()))
								// if err != nil {
								// 	log.Errorln("send error:", err.Error(), result)
								// }
							}
						} else {
							b, err := json.Marshal(tmp)
							if err != nil {
								log.Errorln(err.Error())
							}
							log.Debugln(string(b))
						}
					}
				}(got.CounterSamples)
			}
		}

		//sflow := packet.Layer(layers.LayerTypeSFlow)
		//if sflow != nil {
		//	fmt.Println("SFLOW layer detected")
		//}
	}
}

func (this *Collected) ListenSflowAll(protocol, port string) {
	//Open Device
	handle, err := pcap.OpenLive(this.DeviceName, this.SnapShotLen, this.Promiscuous, this.Timeout)
	if err != nil {
		log.Error(err)
		return
	}
	defer handle.Close()

	//Set filter
	var filter string = fmt.Sprintf("%s and port %s", protocol, port)
	err = handle.SetBPFFilter(filter)
	if err != nil {
		log.Error(err)
	}
	log.Info(fmt.Sprintf("Only capturing %s port %s packets.", protocol, port))
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	for packet := range packetSource.Packets() {
		Origin := NewData()
		p := gopacket.NewPacket(packet.Data(), layers.LayerTypeEthernet, gopacket.Default)
		if p.ErrorLayer() != nil {
			log.Errorln("failed LayerTypeEthernet:", p.ErrorLayer().Error())
		}
		//fmt.Println(p.Dump())
		Origin.Init(p)
		udpLayer := p.Layer(layers.LayerTypeUDP)
		if udpLayer != nil {
			//fmt.Println("UDP layer detected.")
			udp, _ := udpLayer.(*layers.UDP)
			pp := gopacket.NewPacket(udp.Payload, layers.LayerTypeSFlow, gopacket.Default)
			if pp.ErrorLayer() != nil {
				log.Errorln("failed LayerTypeUDP:", p.ErrorLayer().Error())
				// go func(data []byte) {
				// 	s := &layers.SFlowDatagram{}
				// 	s.DecodeSampleFromBytes(data, gopacket.NilDecodeFeedback)
				// 	for _, y := range s.FlowSamples {
				// 		tmp := NewFlowSamples()
				// 		//tmp.InitOriginData(p)
				// 		tmp.Data = Origin
				// 		tmp.InitFlowSampleData(y)
				// 		//for _, yy := range y.Records {
				// 		//	if g1, ok1 := yy.(layers.SFlowRawPacketFlowRecord); ok1 {
				// 		//		tmp.ParseLayers(g1.Header)
				// 		//		b, err := json.Marshal(tmp)
				// 		//		if err != nil {
				// 		//			fmt.Println(err.Error())
				// 		//		}
				// 		//		if this.Udpbool {
				// 		//			this.SendUdp(string(b), false)
				// 		//		} else {
				// 		//			fmt.Println(string(b))
				// 		//		}
				// 		//	}
				// 		//}
				// 		b, err := json.Marshal(tmp)
				// 		if err != nil {
				// 			fmt.Println(err.Error())
				// 		}
				// 		if this.Udpbool {
				// 			this.SendUdp(string(b), false)
				// 		} else {
				// 			fmt.Println(string(b))
				// 		}
				// 	}

				// 	sc := &layers.SFlowDatagram{}
				// 	sc.DecodeCounterFromBytes(data, gopacket.NilDecodeFeedback)
				// 	if len(s.CounterSamples) != 0 {
				// 		//log.Error("Error out of bounds ")
				// 		tmp := NewCounterFlow()
				// 		//tmp.InitOriginData(p)
				// 		tmp.InitCounterSampleStruct(sc)

				// 		b, err := json.Marshal(tmp)
				// 		if err != nil {
				// 			fmt.Println(err.Error())
				// 		}
				// 		if this.Udpbool {
				// 			this.SendUdp(string(b), true)
				// 		} else {
				// 			fmt.Println(string(b))
				// 		}
				// 	}
				// }(udp.Payload)
			}

			if got, ok := pp.ApplicationLayer().(*layers.SFlowDatagram); ok {
				go func(Sample []layers.SFlowFlowSample, Counter []layers.SFlowCounterSample) {
					if len(Sample) > 0 {
						for _, y := range Sample {
							//beego.Critical(len(y.Records),y.RecordCount)
							tmp := NewFlowSamples()
							//tmp.InitOriginData(p)
							tmp.Data = Origin
							tmp.InitFlowSampleData(y)

							if this.Udpbool {
								b, err := json.Marshal(tmp)
								if err != nil {
									log.Errorln(err.Error())
								}
								this.SendUdp(string(b), false)
							} else if this.IsEs {
								result, err := ParseSflowV5ToEs(tmp, nil)
								if err != nil {
									log.Errorln("err", err.Error())
								} else {
									DataChannel <- result
									// err = parse.CreateEs(result, "doc", fmt.Sprintf("ABC%d", time.Now().UnixNano()))
									// if err != nil {
									// 	log.Errorln("send error:", err.Error(), result)
									// }
								}
							} else {
								b, err := json.Marshal(tmp)
								if err != nil {
									log.Errorln(err.Error())
								}
								log.Debugln(string(b))
							}
						}
					}

					if len(Counter) > 0 {
						tmp := NewCounterFlow()
						tmp.InitOriginData(p)
						for _, y := range Counter {
							//beego.Critical(len(y.Records),y.RecordCount)
							tmp.InitCounterSample(y)
						}

						if this.Udpbool {
							b, err := json.Marshal(tmp)
							if err != nil {
								log.Errorln(err.Error())
							}
							this.SendUdp(string(b), true)
						} else if this.IsEs {
							result, err := ParseSflowV5ToEs(nil, tmp)
							if err != nil {
								log.Errorln("err", err.Error())
							} else {
								DataChannel <- result
								// err = parse.CreateEs(result, "doc", fmt.Sprintf("ABC%d", time.Now().UnixNano()))
								// if err != nil {
								// 	log.Errorln("send error:", err.Error(), result)
								// }
							}
						} else {
							b, err := json.Marshal(tmp)
							if err != nil {
								log.Errorln(err.Error())
							}
							log.Debugln(string(b))
						}
					}
				}(got.FlowSamples, got.CounterSamples)
			}
		}
	}
}

func (this *Collected) ListenNetFlowV5(protocol, port string) {
	//Open Device
	handle, err := pcap.OpenLive(this.DeviceName, this.SnapShotLen, this.Promiscuous, this.Timeout)
	if err != nil {
		log.Error(err)
		panic(err)
	}
	defer handle.Close()

	//Set filter
	var filter string = fmt.Sprintf("%s and port %s", protocol, port)
	err = handle.SetBPFFilter(filter)
	if err != nil {
		log.Error(err)
	}
	log.Info(fmt.Sprintf("Only capturing %s port %s packets.", protocol, port))
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	for packet := range packetSource.Packets() {
		//log.Info("1")
		go func(packet gopacket.Packet) {
			//log.Info("2")
			//log.Error("############开始解析#############")
			udpLayer := packet.Layer(layers.LayerTypeUDP)
			if udpLayer != nil {
				//fmt.Println("UDP layer detected.")
				udp, _ := udpLayer.(*layers.UDP)

				tmp := NetFlowV5{}

				for _, x := range tmp.PayLoadToNetFlowV5(udp.Payload, packet.NetworkLayer().NetworkFlow().Src().String()) {
					this.SendUdp(x, false)
				}
				//log.Error(len(data))
				//fmt.Println(data)
			}
		}(packet)
	}
}
