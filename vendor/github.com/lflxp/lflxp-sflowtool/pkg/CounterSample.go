package pkg

import (
	"errors"
	"fmt"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

// Counter samples report information about various counter
// objects. Typically these are items like IfInOctets, or
// CPU / Memory stats, etc. SFlow will report these at regular
// intervals as configured on the agent. If one were sufficiently
// industrious, this could be used to replace the typical
// SNMP polling used for such things.
type SFlowCounterSample struct {
	Data                          Data
	EnterpriseID                  string
	Format                        string
	SampleLength                  uint32
	SequenceNumber                uint32
	SourceIDClass                 string
	SourceIDIndex                 string
	RecordCount                   uint32
	SFlowGenericInterfaceCounters SFlowGenericInterfaceCounters
	SFlowEthernetCounters         SFlowEthernetCounters
	SFlowProcessorCounters        SFlowProcessorCounters
}

// **************************************************
//  Counter Record
// **************************************************

//  0                      15                      31
//  +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//  |      20 bit Interprise (0)     |12 bit format |
//  +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//  |                  counter length               |
//  +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//  /                   counter data                /
//  /                                               /
//  +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+

type SFlowBaseCounterRecord struct {
	EnterpriseID   string
	Format         string
	FlowDataLength uint32
}

// **************************************************
//  Counter Record
// **************************************************

//  0                      15                      31
//  +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//  |      20 bit Interprise (0)     |12 bit format |
//  +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//  |                  counter length               |
//  +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//  |                    IfIndex                    |
//  +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//  |                    IfType                     |
//  +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//  |                   IfSpeed                     |
//  +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//  |                   IfDirection                 |
//  +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//  |                    IfStatus                   |
//  +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//  |                   IFInOctets                  |
//  |                                               |
//  +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//  |                   IfInUcastPkts               |
//  +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//  |                  IfInMulticastPkts            |
//  +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//  |                  IfInBroadcastPkts            |
//  +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//  |                    IfInDiscards               |
//  +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//  |                    InInErrors                 |
//  +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//  |                  IfInUnknownProtos            |
//  +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//  |                   IfOutOctets                 |
//  |                                               |
//  +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//  |                   IfOutUcastPkts              |
//  +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//  |                  IfOutMulticastPkts           |
//  +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//  |                  IfOutBroadcastPkts           |
//  +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//  |                   IfOutDiscards               |
//  +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//  |                    IfOUtErrors                |
//  +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//  |                 IfPromiscouousMode            |
//  +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+

type SFlowGenericInterfaceCounters struct {
	SFlowBaseCounterRecord SFlowBaseCounterRecord
	IfIndex                uint32
	IfType                 uint32
	IfSpeed                uint64
	IfDirection            uint32
	IfStatus               uint32
	IfInOctets             uint64
	IfInUcastPkts          uint32
	IfInMulticastPkts      uint32
	IfInBroadcastPkts      uint32
	IfInDiscards           uint32
	IfInErrors             uint32
	IfInUnknownProtos      uint32
	IfOutOctets            uint64
	IfOutUcastPkts         uint32
	IfOutMulticastPkts     uint32
	IfOutBroadcastPkts     uint32
	IfOutDiscards          uint32
	IfOutErrors            uint32
	IfPromiscuousMode      uint32
}

// **************************************************
//  Counter Record
// **************************************************

//  0                      15                      31
//  +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//  |      20 bit Interprise (0)     |12 bit format |
//  +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//  |                  counter length               |
//  +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//  /                   counter data                /
//  /                                               /
//  +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+

type SFlowEthernetCounters struct {
	SFlowBaseCounterRecord    SFlowBaseCounterRecord
	AlignmentErrors           uint32
	FCSErrors                 uint32
	SingleCollisionFrames     uint32
	MultipleCollisionFrames   uint32
	SQETestErrors             uint32
	DeferredTransmissions     uint32
	LateCollisions            uint32
	ExcessiveCollisions       uint32
	InternalMacTransmitErrors uint32
	CarrierSenseErrors        uint32
	FrameTooLongs             uint32
	InternalMacReceiveErrors  uint32
	SymbolErrors              uint32
}

// **************************************************
//  Processor Counter Record
// **************************************************
//  0                      15                      31
//  +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//  |      20 bit Interprise (0)     |12 bit format |
//  +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//  |                  counter length               |
//  +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//  |                    FiveSecCpu                 |
//  +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//  |                    OneMinCpu                  |
//  +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//  |                    GiveMinCpu                 |
//  +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//  |                   TotalMemory                 |
//  +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//  |                    FreeMemory                 |
//  +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+

type SFlowProcessorCounters struct {
	SFlowBaseCounterRecord SFlowBaseCounterRecord
	FiveSecCpu             uint32 // 5 second average CPU utilization
	OneMinCpu              uint32 // 1 minute average CPU utilization
	FiveMinCpu             uint32 // 5 minute average CPU utilization
	TotalMemory            uint64 // total memory (in bytes)
	FreeMemory             uint64 // free memory (in bytes)
}

func NewCounterFlow() *SFlowCounterSample {
	return &SFlowCounterSample{}
}

func (this *SFlowCounterSample) InitOriginData(p gopacket.Packet) error {
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

func (this *SFlowCounterSample) InitCounterSample(p layers.SFlowCounterSample) error {
	this.EnterpriseID = p.EnterpriseID.String()
	this.Format = p.Format.String()
	this.SampleLength = p.SampleLength
	this.SequenceNumber = p.SequenceNumber
	this.SourceIDClass = p.SourceIDClass.String()
	this.SourceIDIndex = fmt.Sprintf("%v", p.SourceIDIndex)
	this.RecordCount = p.RecordCount
	for _, yy := range p.Records {
		if g1, ok1 := yy.(layers.SFlowGenericInterfaceCounters); ok1 {
			this.SFlowGenericInterfaceCounters.SFlowBaseCounterRecord.EnterpriseID = g1.EnterpriseID.String()
			this.SFlowGenericInterfaceCounters.SFlowBaseCounterRecord.Format = g1.Format.String()
			this.SFlowGenericInterfaceCounters.SFlowBaseCounterRecord.FlowDataLength = g1.FlowDataLength
			this.SFlowGenericInterfaceCounters.IfIndex = g1.IfIndex
			this.SFlowGenericInterfaceCounters.IfType = g1.IfType
			this.SFlowGenericInterfaceCounters.IfSpeed = g1.IfSpeed
			this.SFlowGenericInterfaceCounters.IfDirection = g1.IfDirection
			this.SFlowGenericInterfaceCounters.IfStatus = g1.IfStatus
			this.SFlowGenericInterfaceCounters.IfInOctets = g1.IfInOctets
			this.SFlowGenericInterfaceCounters.IfInUcastPkts = g1.IfInUcastPkts
			this.SFlowGenericInterfaceCounters.IfInMulticastPkts = g1.IfInMulticastPkts
			this.SFlowGenericInterfaceCounters.IfInBroadcastPkts = g1.IfInBroadcastPkts
			this.SFlowGenericInterfaceCounters.IfInDiscards = g1.IfInDiscards
			this.SFlowGenericInterfaceCounters.IfInErrors = g1.IfInErrors
			this.SFlowGenericInterfaceCounters.IfInUnknownProtos = g1.IfInUnknownProtos
			this.SFlowGenericInterfaceCounters.IfOutOctets = g1.IfOutOctets
			this.SFlowGenericInterfaceCounters.IfOutUcastPkts = g1.IfOutUcastPkts
			this.SFlowGenericInterfaceCounters.IfOutMulticastPkts = g1.IfOutMulticastPkts
			this.SFlowGenericInterfaceCounters.IfOutBroadcastPkts = g1.IfOutBroadcastPkts
			this.SFlowGenericInterfaceCounters.IfOutDiscards = g1.IfOutDiscards
			this.SFlowGenericInterfaceCounters.IfOutErrors = g1.IfOutErrors
			this.SFlowGenericInterfaceCounters.IfPromiscuousMode = g1.IfPromiscuousMode
		} else if g2, ok2 := yy.(layers.SFlowEthernetCounters); ok2 {
			this.SFlowEthernetCounters.SFlowBaseCounterRecord.EnterpriseID = g2.EnterpriseID.String()
			this.SFlowEthernetCounters.SFlowBaseCounterRecord.Format = g2.Format.String()
			this.SFlowEthernetCounters.SFlowBaseCounterRecord.FlowDataLength = g2.FlowDataLength
			this.SFlowEthernetCounters.AlignmentErrors = g2.AlignmentErrors
			this.SFlowEthernetCounters.FCSErrors = g2.FCSErrors
			this.SFlowEthernetCounters.SingleCollisionFrames = g2.SingleCollisionFrames
			this.SFlowEthernetCounters.MultipleCollisionFrames = g2.MultipleCollisionFrames
			this.SFlowEthernetCounters.SQETestErrors = g2.SQETestErrors
			this.SFlowEthernetCounters.DeferredTransmissions = g2.DeferredTransmissions
			this.SFlowEthernetCounters.LateCollisions = g2.LateCollisions
			this.SFlowEthernetCounters.ExcessiveCollisions = g2.ExcessiveCollisions
			this.SFlowEthernetCounters.InternalMacTransmitErrors = g2.InternalMacTransmitErrors
			this.SFlowEthernetCounters.CarrierSenseErrors = g2.CarrierSenseErrors
			this.SFlowEthernetCounters.FrameTooLongs = g2.FrameTooLongs
			this.SFlowEthernetCounters.InternalMacReceiveErrors = g2.InternalMacReceiveErrors
			this.SFlowEthernetCounters.SymbolErrors = g2.SymbolErrors
		} else if g3, ok3 := yy.(layers.SFlowProcessorCounters); ok3 {
			this.SFlowProcessorCounters.SFlowBaseCounterRecord.EnterpriseID = g3.EnterpriseID.String()
			this.SFlowProcessorCounters.SFlowBaseCounterRecord.Format = g3.Format.String()
			this.SFlowProcessorCounters.SFlowBaseCounterRecord.FlowDataLength = g3.FlowDataLength
			this.SFlowProcessorCounters.FiveSecCpu = g3.FiveSecCpu
			this.SFlowProcessorCounters.OneMinCpu = g3.OneMinCpu
			this.SFlowProcessorCounters.FiveMinCpu = g3.FiveMinCpu
			this.SFlowProcessorCounters.TotalMemory = g3.TotalMemory
			this.SFlowProcessorCounters.FreeMemory = g3.FreeMemory
		} else {
			return errors.New("nothing deteced")
		}
	}
	return nil
}

func (this *SFlowCounterSample) InitCounterSampleStruct(p *layers.SFlowDatagram) error {
	for _, c := range p.CounterSamples {
		this.EnterpriseID = c.EnterpriseID.String()
		this.Format = c.Format.String()
		this.SampleLength = c.SampleLength
		this.SequenceNumber = c.SequenceNumber
		this.SourceIDClass = c.SourceIDClass.String()
		this.SourceIDIndex = fmt.Sprintf("%v", c.SourceIDIndex)
		this.RecordCount = c.RecordCount
		for _, yy := range c.Records {
			if g1, ok1 := yy.(layers.SFlowGenericInterfaceCounters); ok1 {
				this.SFlowGenericInterfaceCounters.SFlowBaseCounterRecord.EnterpriseID = g1.EnterpriseID.String()
				this.SFlowGenericInterfaceCounters.SFlowBaseCounterRecord.Format = g1.Format.String()
				this.SFlowGenericInterfaceCounters.SFlowBaseCounterRecord.FlowDataLength = g1.FlowDataLength
				this.SFlowGenericInterfaceCounters.IfIndex = g1.IfIndex
				this.SFlowGenericInterfaceCounters.IfType = g1.IfType
				this.SFlowGenericInterfaceCounters.IfSpeed = g1.IfSpeed
				this.SFlowGenericInterfaceCounters.IfDirection = g1.IfDirection
				this.SFlowGenericInterfaceCounters.IfStatus = g1.IfStatus
				this.SFlowGenericInterfaceCounters.IfInOctets = g1.IfInOctets
				this.SFlowGenericInterfaceCounters.IfInUcastPkts = g1.IfInUcastPkts
				this.SFlowGenericInterfaceCounters.IfInMulticastPkts = g1.IfInMulticastPkts
				this.SFlowGenericInterfaceCounters.IfInBroadcastPkts = g1.IfInBroadcastPkts
				this.SFlowGenericInterfaceCounters.IfInDiscards = g1.IfInDiscards
				this.SFlowGenericInterfaceCounters.IfInErrors = g1.IfInErrors
				this.SFlowGenericInterfaceCounters.IfInUnknownProtos = g1.IfInUnknownProtos
				this.SFlowGenericInterfaceCounters.IfOutOctets = g1.IfOutOctets
				this.SFlowGenericInterfaceCounters.IfOutUcastPkts = g1.IfOutUcastPkts
				this.SFlowGenericInterfaceCounters.IfOutMulticastPkts = g1.IfOutMulticastPkts
				this.SFlowGenericInterfaceCounters.IfOutBroadcastPkts = g1.IfOutBroadcastPkts
				this.SFlowGenericInterfaceCounters.IfOutDiscards = g1.IfOutDiscards
				this.SFlowGenericInterfaceCounters.IfOutErrors = g1.IfOutErrors
				this.SFlowGenericInterfaceCounters.IfPromiscuousMode = g1.IfPromiscuousMode
			} else if g2, ok2 := yy.(layers.SFlowEthernetCounters); ok2 {
				this.SFlowEthernetCounters.SFlowBaseCounterRecord.EnterpriseID = g2.EnterpriseID.String()
				this.SFlowEthernetCounters.SFlowBaseCounterRecord.Format = g2.Format.String()
				this.SFlowEthernetCounters.SFlowBaseCounterRecord.FlowDataLength = g2.FlowDataLength
				this.SFlowEthernetCounters.AlignmentErrors = g2.AlignmentErrors
				this.SFlowEthernetCounters.FCSErrors = g2.FCSErrors
				this.SFlowEthernetCounters.SingleCollisionFrames = g2.SingleCollisionFrames
				this.SFlowEthernetCounters.MultipleCollisionFrames = g2.MultipleCollisionFrames
				this.SFlowEthernetCounters.SQETestErrors = g2.SQETestErrors
				this.SFlowEthernetCounters.DeferredTransmissions = g2.DeferredTransmissions
				this.SFlowEthernetCounters.LateCollisions = g2.LateCollisions
				this.SFlowEthernetCounters.ExcessiveCollisions = g2.ExcessiveCollisions
				this.SFlowEthernetCounters.InternalMacTransmitErrors = g2.InternalMacTransmitErrors
				this.SFlowEthernetCounters.CarrierSenseErrors = g2.CarrierSenseErrors
				this.SFlowEthernetCounters.FrameTooLongs = g2.FrameTooLongs
				this.SFlowEthernetCounters.InternalMacReceiveErrors = g2.InternalMacReceiveErrors
				this.SFlowEthernetCounters.SymbolErrors = g2.SymbolErrors
			} else if g3, ok3 := yy.(layers.SFlowProcessorCounters); ok3 {
				this.SFlowProcessorCounters.SFlowBaseCounterRecord.EnterpriseID = g3.EnterpriseID.String()
				this.SFlowProcessorCounters.SFlowBaseCounterRecord.Format = g3.Format.String()
				this.SFlowProcessorCounters.SFlowBaseCounterRecord.FlowDataLength = g3.FlowDataLength
				this.SFlowProcessorCounters.FiveSecCpu = g3.FiveSecCpu
				this.SFlowProcessorCounters.OneMinCpu = g3.OneMinCpu
				this.SFlowProcessorCounters.FiveMinCpu = g3.FiveMinCpu
				this.SFlowProcessorCounters.TotalMemory = g3.TotalMemory
				this.SFlowProcessorCounters.FreeMemory = g3.FreeMemory
			} else {
				return errors.New("nothing deteced")
			}
		}
	}

	return nil
}
