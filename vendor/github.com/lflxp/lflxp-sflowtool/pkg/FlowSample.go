package pkg

import (
	"encoding/binary"
	"errors"
	"fmt"
	"net"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

// SFlowBaseFlowRecord holds the fields common to all records
// of type SFlowFlowRecordType
type SFlowBaseFlowRecord struct {
	EnterpriseID   string
	Format         string
	FlowDataLength uint32
}

type SFlowRawPacketFlowRecord struct {
	SFlowBaseFlowRecord SFlowBaseFlowRecord
	HeaderProtocol      string
	FrameLength         uint32
	PayloadRemoved      uint32
	HeaderLength        uint32
	Header              Header
}

//flow流详细信息
type Header struct {
	FlowRecords   uint32 //flow流数据量
	Packets       int    //包个数
	Bytes         uint32 //字节大小
	RateBytes     uint32 //自动采样率计算
	SrcMac        string
	DstMac        string
	SrcIP         string
	DstIP         string
	Ipv4_version  uint8
	Ipv4_ihl      uint8
	Ipv4_tos      uint8
	Ipv4_ttl      uint8
	Ipv4_protocol string
	SrcPort       string //如果是icmp的就只把数据写入这个
	DstPort       string
}

// SFlowExtendedSwitchFlowRecord give additional information
// about the sampled packet if it's available. It's mainly
// useful for getting at the incoming and outgoing VLANs
// An agent may or may not provide this information.
type SFlowExtendedSwitchFlowRecord struct {
	SFlowBaseFlowRecord  SFlowBaseFlowRecord
	IncomingVLAN         uint32
	IncomingVLANPriority uint32
	OutgoingVLAN         uint32
	OutgoingVLANPriority uint32
}

// SFlowExtendedRouterFlowRecord gives additional information
// about the layer 3 routing information used to forward
// the packet
type SFlowExtendedRouterFlowRecord struct {
	SFlowBaseFlowRecord    SFlowBaseFlowRecord
	NextHop                net.IP
	NextHopSourceMask      uint32
	NextHopDestinationMask uint32
}

// **************************************************
//  Packet Ethernet Data Record
// **************************************************

//  0                      15                      31
//  +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//  |                     Tag                       |
//  +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//  |                    Length                     |
//  +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//  |                  Length Bytes                 |
//  +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//  |  Src Mac  | Dst Mac   |
//  +--+--+--+--+--+--+--+--+
type SFlowEthernetFrameRecord struct {
	//为2代表是Ethernet Frame Data字段
	Format uint32
	//总的字节数（不包含tag和length字段）
	Length uint32
	//源mac地址8字节
	SrcMac []byte
	//目的mac地址8字节
	DstMac []byte
	Type   uint32
}

func decodeSFlowEthernetFrameRecord(data *[]byte) (SFlowEthernetFrameRecord, error) {
	sef := SFlowEthernetFrameRecord{}

	*data, sef.Format = (*data)[4:], binary.BigEndian.Uint32((*data)[:4])
	*data, sef.Length = (*data)[4:], binary.BigEndian.Uint32((*data)[:4])
	*data, sef.SrcMac = (*data)[6:], (*data)[:6]
	*data, sef.DstMac = (*data)[6:], (*data)[:6]
	*data, sef.Type = (*data)[4:], binary.BigEndian.Uint32((*data)[:4])

	return sef, nil
}

// SFlowExtendedGatewayFlowRecord describes information treasured by
// nework engineers everywhere: AS path information listing which
// BGP peer sent the packet, and various other BGP related info.
// This information is vital because it gives a picture of how much
// traffic is being sent from / received by various BGP peers.

// Extended gateway records have the following structure:

//  0                      15                      31
//  +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//  |      20 bit Interprise (0)     |12 bit format |
//  +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//  |                  record length                |
//  +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//  |   IP version of next hop router (1=v4|2=v6)   |
//  +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//  /     Next Hop address (v4=4byte|v6=16byte)     /
//  /                                               /
//  +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//  |                       AS                      |
//  +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//  |                  Source AS                    |
//  +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//  |                    Peer AS                    |
//  +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//  |                  AS Path Count                |
//  +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//  /                AS Path / Sequence             /
//  /                                               /
//  +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//  /                   Communities                 /
//  /                                               /
//  +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//  |                    Local Pref                 |
//  +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+

// AS Path / Sequence:

//  0                      15                      31
//  +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//  |     AS Source Type (Path=1 / Sequence=2)      |
//  +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//  |              Path / Sequence length           |
//  +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//  /              Path / Sequence Members          /
//  /                                               /
//  +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+

// Communities:

//  0                      15                      31
//  +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//  |                communitiy length              |
//  +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//  /              communitiy Members               /
//  /                                               /
//  +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+

type SFlowExtendedGatewayFlowRecord struct {
	SFlowBaseFlowRecord SFlowBaseFlowRecord
	NextHop             net.IP
	AS                  uint32
	SourceAS            uint32
	PeerAS              uint32
	ASPathCount         uint32
	ASPath              []layers.SFlowASDestination
	Communities         []uint32
	LocalPref           uint32
}

type SFlowExtendedUserFlow struct {
	SFlowBaseFlowRecord SFlowBaseFlowRecord
	SourceCharSet       string
	SourceUserID        string
	DestinationCharSet  string
	DestinationUserID   string
}

type Data struct {
	Datagram        Datagram
	DatagramVersion uint32
	AgentAddress    net.IP
	SubAgentID      uint32
	SequenceNumber  uint32
	AgentUptime     uint32
	SampleCount     uint32
}

func NewData() *Data {
	return &Data{}
}

func (this *Data) Init(p gopacket.Packet) error {
	if p.ErrorLayer() != nil {
		return errors.New(fmt.Sprintf("failed : %s", p.ErrorLayer().Error()))
	}
	eth := p.Layer(layers.LayerTypeEthernet)
	if eth != nil {
		ethernetPacket, _ := eth.(*layers.Ethernet)
		this.Datagram.SrcMac = ethernetPacket.SrcMAC.String()
		this.Datagram.DstMac = ethernetPacket.DstMAC.String()
	} else {
		return errors.New("LayerTypeEthernet error")
	}

	// Let's see if the packet is IP ∂(even though the ether type told us)
	ipLayer := p.Layer(layers.LayerTypeIPv4)
	if ipLayer != nil {
		ip, _ := ipLayer.(*layers.IPv4)

		// IP layer variables:
		// Version (Either 4 or 6)
		// IHL (IP Header Length in 32-bit words)
		// TOS, Length, Id, Flags, FragOffset, TTL, Protocol (TCP?),
		// Checksum, SrcIP, DstIP
		this.Datagram.SrcIP = ip.SrcIP.String()
		this.Datagram.DstIP = ip.DstIP.String()
	} else {
		return errors.New("LayerTypeIPv4 error")
	}

	udpLayer := p.Layer(layers.LayerTypeUDP)
	if udpLayer != nil {
		udp, _ := udpLayer.(*layers.UDP)
		this.Datagram.SrcPort = udp.SrcPort.String()
		this.Datagram.DstPort = udp.DstPort.String()

		pp := gopacket.NewPacket(udp.Payload, layers.LayerTypeSFlow, gopacket.Default)
		if pp.ErrorLayer() != nil {
			//fmt.Println(pp.Data())
			this.decodeDataFromBytes(pp.Data())
		}
		if got, ok := pp.ApplicationLayer().(*layers.SFlowDatagram); ok {
			this.DatagramVersion = got.DatagramVersion
			this.AgentAddress = got.AgentAddress
			this.SubAgentID = got.SubAgentID
			this.SequenceNumber = got.SequenceNumber
			this.AgentUptime = got.AgentUptime
			this.SampleCount = got.SampleCount
		}
	} else {
		return errors.New("LayerTypeUDP error")
	}
	return nil
}

func (this *Data) decodeDataFromBytes(data []byte) error {
	var agentAddressType layers.SFlowIPType

	data, this.DatagramVersion = data[4:], binary.BigEndian.Uint32(data[:4])
	data, agentAddressType = data[4:], layers.SFlowIPType(binary.BigEndian.Uint32(data[:4]))
	data, this.AgentAddress = data[agentAddressType.Length():], data[:agentAddressType.Length()]
	data, this.SubAgentID = data[4:], binary.BigEndian.Uint32(data[:4])
	data, this.SequenceNumber = data[4:], binary.BigEndian.Uint32(data[:4])
	data, this.AgentUptime = data[4:], binary.BigEndian.Uint32(data[:4])
	data, this.SampleCount = data[4:], binary.BigEndian.Uint32(data[:4])
	return nil
}

//原始报文信息即交换机物理设备信息
type Datagram struct {
	SrcMac  string
	DstMac  string
	SrcIP   string
	DstIP   string
	SrcPort string
	DstPort string
}

type FlowSamples struct {
	Data                           *Data
	EnterpriseID                   string
	Format                         string
	SampleLength                   uint32
	SequenceNumber                 uint32
	SourceIDClass                  string
	SourceIDIndex                  string
	SamplingRate                   uint32
	SamplePool                     uint32
	Dropped                        uint32
	InputInterfaceFormat           uint32
	InputInterface                 uint32
	OutputInterfaceFormat          uint32
	OutputInterface                uint32
	RecordCount                    uint32
	SFlowRawPacketFlowRecord       SFlowRawPacketFlowRecord
	SFlowExtendedSwitchFlowRecord  SFlowExtendedSwitchFlowRecord
	SFlowExtendedRouterFlowRecord  SFlowExtendedRouterFlowRecord
	SFlowExtendedGatewayFlowRecord SFlowExtendedGatewayFlowRecord
	SFlowExtendedUserFlow          SFlowExtendedUserFlow
	//SFlowEthernetFrameRecord SFlowEthernetFrameRecord
}

func NewFlowSamples() *FlowSamples {
	return &FlowSamples{}
}

func (this *FlowSamples) SendUdp(result, CounterHost, Host string, counter bool) {
	if counter {
		conn, err := net.Dial("udp", CounterHost)
		defer conn.Close()
		if err != nil {
			panic(err)
		}
		conn.Write([]byte(result))
	} else {
		conn, err := net.Dial("udp", Host)
		defer conn.Close()
		if err != nil {
			panic(err)
		}
		conn.Write([]byte(result))
	}
}

func (this *FlowSamples) InitOriginData(p gopacket.Packet) error {
	if p.ErrorLayer() != nil {
		return errors.New(fmt.Sprintf("failed : %s", p.ErrorLayer().Error()))
	}
	eth := p.Layer(layers.LayerTypeEthernet)
	if eth != nil {
		ethernetPacket, _ := eth.(*layers.Ethernet)
		this.Data.Datagram.SrcMac = ethernetPacket.SrcMAC.String()
		this.Data.Datagram.DstMac = ethernetPacket.DstMAC.String()
	} else {
		return errors.New("LayerTypeEthernet error")
	}

	// Let's see if the packet is IP ∂(even though the ether type told us)
	ipLayer := p.Layer(layers.LayerTypeIPv4)
	if ipLayer != nil {
		ip, _ := ipLayer.(*layers.IPv4)

		// IP layer variables:
		// Version (Either 4 or 6)
		// IHL (IP Header Length in 32-bit words)
		// TOS, Length, Id, Flags, FragOffset, TTL, Protocol (TCP?),
		// Checksum, SrcIP, DstIP
		this.Data.Datagram.SrcIP = ip.SrcIP.String()
		this.Data.Datagram.DstIP = ip.DstIP.String()
	} else {
		return errors.New("LayerTypeIPv4 error")
	}

	udpLayer := p.Layer(layers.LayerTypeUDP)
	if udpLayer != nil {
		udp, _ := udpLayer.(*layers.UDP)
		this.Data.Datagram.SrcPort = udp.SrcPort.String()
		this.Data.Datagram.DstPort = udp.DstPort.String()

		pp := gopacket.NewPacket(udp.Payload, layers.LayerTypeSFlow, gopacket.Default)
		if got, ok := pp.ApplicationLayer().(*layers.SFlowDatagram); ok {
			this.Data.DatagramVersion = got.DatagramVersion
			this.Data.AgentAddress = got.AgentAddress
			this.Data.SubAgentID = got.SubAgentID
			this.Data.SequenceNumber = got.SequenceNumber
			this.Data.AgentUptime = got.AgentUptime
			this.Data.SampleCount = got.SampleCount
		}
	} else {
		return errors.New("LayerTypeUDP error")
	}
	return nil
}

func (this *FlowSamples) InitFlowSampleData(p layers.SFlowFlowSample) error {
	this.EnterpriseID = p.EnterpriseID.String()
	this.Format = p.Format.String()
	this.SampleLength = p.SampleLength
	this.SequenceNumber = p.SequenceNumber
	this.SourceIDClass = p.SourceIDClass.String()
	this.SourceIDIndex = fmt.Sprintf("%v", p.SourceIDIndex)
	this.SamplingRate = p.SamplingRate
	this.SamplePool = p.SamplePool
	this.Dropped = p.Dropped
	this.InputInterfaceFormat = p.InputInterfaceFormat
	this.InputInterface = p.InputInterface
	this.OutputInterfaceFormat = p.OutputInterfaceFormat
	this.OutputInterface = p.OutputInterface
	this.RecordCount = p.RecordCount

	for _, yy := range p.Records {
		if g1, ok1 := yy.(layers.SFlowRawPacketFlowRecord); ok1 {
			this.SFlowRawPacketFlowRecord.SFlowBaseFlowRecord.EnterpriseID = g1.EnterpriseID.String()
			this.SFlowRawPacketFlowRecord.SFlowBaseFlowRecord.Format = g1.Format.String()
			this.SFlowRawPacketFlowRecord.SFlowBaseFlowRecord.FlowDataLength = g1.FlowDataLength
			this.SFlowRawPacketFlowRecord.HeaderProtocol = g1.HeaderProtocol.String()
			this.SFlowRawPacketFlowRecord.FrameLength = g1.FrameLength
			this.SFlowRawPacketFlowRecord.PayloadRemoved = g1.PayloadRemoved
			this.SFlowRawPacketFlowRecord.HeaderLength = g1.HeaderLength
			this.SFlowRawPacketFlowRecord.Header.FlowRecords = g1.FlowDataLength
			this.SFlowRawPacketFlowRecord.Header.Bytes = g1.FrameLength
			this.SFlowRawPacketFlowRecord.Header.RateBytes = g1.FrameLength * this.SamplingRate
			this.SFlowRawPacketFlowRecord.Header.Packets = 1
			this.ParseLayers(g1.Header)
		} else if g2, ok2 := yy.(layers.SFlowExtendedSwitchFlowRecord); ok2 {
			this.SFlowExtendedSwitchFlowRecord.SFlowBaseFlowRecord.EnterpriseID = g2.EnterpriseID.String()
			this.SFlowExtendedSwitchFlowRecord.SFlowBaseFlowRecord.FlowDataLength = g2.FlowDataLength
			this.SFlowExtendedSwitchFlowRecord.SFlowBaseFlowRecord.Format = g2.Format.String()
			this.SFlowExtendedSwitchFlowRecord.IncomingVLAN = g2.IncomingVLAN
			this.SFlowExtendedSwitchFlowRecord.IncomingVLANPriority = g2.IncomingVLANPriority
			this.SFlowExtendedSwitchFlowRecord.OutgoingVLAN = g2.OutgoingVLAN
			this.SFlowExtendedSwitchFlowRecord.OutgoingVLANPriority = g2.OutgoingVLANPriority
		} else if g3, ok3 := yy.(layers.SFlowExtendedRouterFlowRecord); ok3 {
			this.SFlowExtendedRouterFlowRecord.SFlowBaseFlowRecord.EnterpriseID = g3.EnterpriseID.String()
			this.SFlowExtendedRouterFlowRecord.SFlowBaseFlowRecord.FlowDataLength = g3.FlowDataLength
			this.SFlowExtendedRouterFlowRecord.SFlowBaseFlowRecord.Format = g3.Format.String()
			this.SFlowExtendedRouterFlowRecord.NextHop = g3.NextHop
			this.SFlowExtendedRouterFlowRecord.NextHopSourceMask = g3.NextHopSourceMask
			this.SFlowExtendedRouterFlowRecord.NextHopDestinationMask = g3.NextHopDestinationMask
		} else if g4, ok4 := yy.(layers.SFlowExtendedGatewayFlowRecord); ok4 {
			this.SFlowExtendedGatewayFlowRecord.SFlowBaseFlowRecord.EnterpriseID = g4.EnterpriseID.String()
			this.SFlowExtendedGatewayFlowRecord.SFlowBaseFlowRecord.FlowDataLength = g4.FlowDataLength
			this.SFlowExtendedGatewayFlowRecord.SFlowBaseFlowRecord.Format = g4.Format.String()
			this.SFlowExtendedGatewayFlowRecord.NextHop = g4.NextHop
			this.SFlowExtendedGatewayFlowRecord.AS = g4.AS
			this.SFlowExtendedGatewayFlowRecord.SourceAS = g4.SourceAS
			this.SFlowExtendedGatewayFlowRecord.PeerAS = g4.PeerAS
			this.SFlowExtendedGatewayFlowRecord.ASPathCount = g4.ASPathCount
			this.SFlowExtendedGatewayFlowRecord.ASPath = g4.ASPath
			this.SFlowExtendedGatewayFlowRecord.Communities = g4.Communities
			this.SFlowExtendedGatewayFlowRecord.LocalPref = g4.LocalPref
		} else if g5, ok5 := yy.(layers.SFlowExtendedUserFlow); ok5 {
			this.SFlowExtendedUserFlow.SFlowBaseFlowRecord.EnterpriseID = g5.EnterpriseID.String()
			this.SFlowExtendedUserFlow.SFlowBaseFlowRecord.Format = g5.Format.String()
			this.SFlowExtendedUserFlow.SFlowBaseFlowRecord.FlowDataLength = g5.FlowDataLength
			this.SFlowExtendedUserFlow.SourceCharSet = fmt.Sprintf("%v", g5.SourceCharSet)
			this.SFlowExtendedUserFlow.SourceUserID = g5.SourceUserID
			this.SFlowExtendedUserFlow.DestinationCharSet = fmt.Sprintf("%v", g5.DestinationCharSet)
			this.SFlowExtendedUserFlow.DestinationUserID = g5.DestinationUserID
		}
		//else if g6,ok6 := yy.(layers.SFlowEthernetFrameRecord); ok6 {
		//	this.SFlowEthernetFrameRecord.Format = g6.Format
		//	this.SFlowEthernetFrameRecord.Length = g6.Length
		//	this.SFlowEthernetFrameRecord.SrcMac = net.HardwareAddr{g6.SrcMac}
		//	this.SFlowEthernetFrameRecord.DstMac = net.HardwareAddr{g6.DstMac}
		//	this.SFlowEthernetFrameRecord.Type = g6.Type
		//}
	}
	return nil
}

func (this *FlowSamples) ParseLayers(p gopacket.Packet) error {
	if p.ErrorLayer() != nil {
		return errors.New(fmt.Sprintf("failed : %s", p.ErrorLayer().Error()))
	}
	eth := p.Layer(layers.LayerTypeEthernet)
	if eth != nil {
		ethernetPacket, _ := eth.(*layers.Ethernet)
		this.SFlowRawPacketFlowRecord.Header.SrcMac = ethernetPacket.SrcMAC.String()
		this.SFlowRawPacketFlowRecord.Header.DstMac = ethernetPacket.DstMAC.String()
	} else {
		return errors.New("LayerTypeEthernet error")
	}

	// Let's see if the packet is IP ∂(even though the ether type told us)
	ipLayer := p.Layer(layers.LayerTypeIPv4)
	if ipLayer != nil {
		ip, _ := ipLayer.(*layers.IPv4)

		// IP layer variables:
		// Version (Either 4 or 6)
		// IHL (IP Header Length in 32-bit words)
		// TOS, Length, Id, Flags, FragOffset, TTL, Protocol (TCP?),
		// Checksum, SrcIP, DstIP
		this.SFlowRawPacketFlowRecord.Header.SrcIP = ip.SrcIP.String()
		this.SFlowRawPacketFlowRecord.Header.DstIP = ip.DstIP.String()
		this.SFlowRawPacketFlowRecord.Header.Ipv4_version = ip.Version
		this.SFlowRawPacketFlowRecord.Header.Ipv4_ihl = ip.IHL
		this.SFlowRawPacketFlowRecord.Header.Ipv4_tos = ip.TOS
		this.SFlowRawPacketFlowRecord.Header.Ipv4_ttl = ip.TTL
		this.SFlowRawPacketFlowRecord.Header.Ipv4_protocol = ip.Protocol.String()
	} else {
		return errors.New("LayerTypeIPv4 error")
	}

	//TCP
	tcplayer := p.Layer(layers.LayerTypeTCP)
	if tcplayer != nil {
		tcp, _ := tcplayer.(*layers.TCP)
		this.SFlowRawPacketFlowRecord.Header.SrcPort = tcp.SrcPort.String()
		this.SFlowRawPacketFlowRecord.Header.DstPort = tcp.DstPort.String()
	}

	//UDP
	udpLayer := p.Layer(layers.LayerTypeUDP)
	if udpLayer != nil {
		udp, _ := udpLayer.(*layers.UDP)
		this.SFlowRawPacketFlowRecord.Header.SrcPort = udp.SrcPort.String()
		this.SFlowRawPacketFlowRecord.Header.DstPort = udp.DstPort.String()
	} else {
		return errors.New("LayerTypeUDP error")
	}

	//ICMP
	//TCP
	icmplayer := p.Layer(layers.LayerTypeICMPv4)
	if icmplayer != nil {
		icmp, _ := icmplayer.(*layers.ICMPv4)
		this.SFlowRawPacketFlowRecord.Header.SrcPort = icmp.TypeCode.String()
	}

	return nil
}
