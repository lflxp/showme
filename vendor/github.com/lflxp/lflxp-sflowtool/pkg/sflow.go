package pkg

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/elastic/beats/libbeat/common"
)

func ParseSflowV5ToEs(sample *FlowSamples, counter *SFlowCounterSample) (string, error) {
	if sample != nil {
		result := common.MapStr{}
		// result["@timestamp"] = time.Now().Format("2006-01-02T15:04:05.000Z")
		result["@timestamp"] = fmt.Sprintf("%s+0800", time.Now().Format("2006-01-02T15:04:05"))
		result["type"] = "sflow"
		result["transport"] = "udp"
		result["src"] = sample.Data.Datagram.SrcMac
		result["dst"] = sample.Data.Datagram.DstMac
		result["status"] = "ok"
		result["notes"] = "sample"

		sflowEvent := common.MapStr{}
		result["sflow"] = sflowEvent

		//sflow agent info
		result["Datagram"] = common.MapStr{
			"IPLength": 4,
			"SrcIP":    sample.Data.Datagram.SrcIP,
			"DstIP":    sample.Data.Datagram.DstIP,
			"SrcPort":  sample.Data.Datagram.SrcPort,
			"DstPort":  sample.Data.Datagram.DstPort,
		}

		//SFlow info
		// result["type"] = "sample"
		result["DatagramVersion"] = sample.Data.DatagramVersion
		result["AgentAddress"] = sample.Data.AgentAddress
		result["SubAgentID"] = sample.Data.SubAgentID
		result["SequenceNumber"] = sample.Data.SequenceNumber
		result["AgentUptime"] = sample.Data.AgentUptime
		result["SampleCount"] = sample.Data.SampleCount
		result["EnterpriseID"] = sample.EnterpriseID
		result["Format"] = sample.Format
		result["SampleLength"] = sample.SampleLength
		result["SequenceNumber"] = sample.SequenceNumber
		result["SourceIDClass"] = sample.SourceIDClass
		result["SourceIDIndex"] = sample.SourceIDIndex
		result["SamplingRate"] = sample.SamplingRate
		result["SamplePool"] = sample.SamplePool
		result["Dropped"] = sample.Dropped
		result["InputInterfaceFormat"] = sample.InputInterfaceFormat
		result["InputInterface"] = sample.InputInterface
		result["OutputInterfaceFormat"] = sample.OutputInterfaceFormat
		result["OutputInterface"] = sample.OutputInterface
		result["RecordCount"] = sample.RecordCount

		//SFlowRawPacketFlowRecord
		result["SFlowRawPacketFlowRecord"] = common.MapStr{
			"SFlowBaseFlowRecord": common.MapStr{
				"EnterpriseID":   sample.SFlowRawPacketFlowRecord.SFlowBaseFlowRecord.EnterpriseID,
				"Format":         sample.SFlowRawPacketFlowRecord.SFlowBaseFlowRecord.Format,
				"FlowDataLength": sample.SFlowRawPacketFlowRecord.SFlowBaseFlowRecord.FlowDataLength,
			},
			"HeaderProtocol": sample.SFlowRawPacketFlowRecord.HeaderProtocol,
			"FrameLength":    sample.SFlowRawPacketFlowRecord.FrameLength,
			"PayloadRemoved": sample.SFlowRawPacketFlowRecord.PayloadRemoved,
			"HeaderLength":   sample.SFlowRawPacketFlowRecord.HeaderLength,
			"Header": common.MapStr{
				"FlowRecords":   sample.SFlowRawPacketFlowRecord.Header.FlowRecords,
				"Packets":       sample.SFlowRawPacketFlowRecord.Header.Packets,
				"Bytes":         sample.SFlowRawPacketFlowRecord.Header.Bytes,
				"RateBytes":     sample.SFlowRawPacketFlowRecord.Header.RateBytes,
				"SrcMac":        sample.SFlowRawPacketFlowRecord.Header.SrcMac,
				"DstMac":        sample.SFlowRawPacketFlowRecord.Header.DstMac,
				"SrcIP":         sample.SFlowRawPacketFlowRecord.Header.SrcIP,
				"DstIP":         sample.SFlowRawPacketFlowRecord.Header.DstIP,
				"Ipv4_version":  sample.SFlowRawPacketFlowRecord.Header.Ipv4_version,
				"Ipv4_ihl":      sample.SFlowRawPacketFlowRecord.Header.Ipv4_ihl,
				"Ipv4_tos":      sample.SFlowRawPacketFlowRecord.Header.Ipv4_tos,
				"Ipv4_ttl":      sample.SFlowRawPacketFlowRecord.Header.Ipv4_ttl,
				"Ipv4_protocol": sample.SFlowRawPacketFlowRecord.Header.Ipv4_protocol,
				"SrcPort":       sample.SFlowRawPacketFlowRecord.Header.SrcPort,
				"DstPort":       sample.SFlowRawPacketFlowRecord.Header.DstPort,
			},
		}

		//SFlowExtendedSwitchFlowRecord
		result["SFlowExtendedSwitchFlowRecord"] = common.MapStr{
			"SFlowBaseFlowRecord": common.MapStr{
				"EnterpriseID":   sample.SFlowExtendedSwitchFlowRecord.SFlowBaseFlowRecord.EnterpriseID,
				"Format":         sample.SFlowExtendedSwitchFlowRecord.SFlowBaseFlowRecord.Format,
				"FlowDataLength": sample.SFlowExtendedSwitchFlowRecord.SFlowBaseFlowRecord.FlowDataLength,
			},
			"IncomingVLAN":         sample.SFlowExtendedSwitchFlowRecord.IncomingVLAN,
			"IncomingVLANPriority": sample.SFlowExtendedSwitchFlowRecord.IncomingVLANPriority,
			"OutgoingVLAN":         sample.SFlowExtendedSwitchFlowRecord.OutgoingVLAN,
			"OutgoingVLANPriority": sample.SFlowExtendedSwitchFlowRecord.OutgoingVLANPriority,
		}

		//SFlowExtendedRouterFlowRecord
		result["SFlowExtendedRouterFlowRecord"] = common.MapStr{
			"SFlowBaseFlowRecord": common.MapStr{
				"EnterpriseID":   sample.SFlowExtendedRouterFlowRecord.SFlowBaseFlowRecord.EnterpriseID,
				"Format":         sample.SFlowExtendedRouterFlowRecord.SFlowBaseFlowRecord.Format,
				"FlowDataLength": sample.SFlowExtendedRouterFlowRecord.SFlowBaseFlowRecord.FlowDataLength,
			},
			"NextHop":                sample.SFlowExtendedRouterFlowRecord.NextHop,
			"NextHopSourceMask":      sample.SFlowExtendedRouterFlowRecord.NextHopSourceMask,
			"NextHopDestinationMask": sample.SFlowExtendedRouterFlowRecord.NextHopDestinationMask,
		}

		//SFlowExtendedGatewayFlowRecord
		result["SFlowExtendedGatewayFlowRecord"] = common.MapStr{
			"SFlowBaseFlowRecord": common.MapStr{
				"EnterpriseID":   sample.SFlowExtendedGatewayFlowRecord.SFlowBaseFlowRecord.EnterpriseID,
				"Format":         sample.SFlowExtendedGatewayFlowRecord.SFlowBaseFlowRecord.Format,
				"FlowDataLength": sample.SFlowExtendedGatewayFlowRecord.SFlowBaseFlowRecord.FlowDataLength,
			},
			"NextHop":     sample.SFlowExtendedGatewayFlowRecord.NextHop,
			"AS":          sample.SFlowExtendedGatewayFlowRecord.AS,
			"SourceAS":    sample.SFlowExtendedGatewayFlowRecord.SourceAS,
			"PeerAS":      sample.SFlowExtendedGatewayFlowRecord.PeerAS,
			"ASPathCount": sample.SFlowExtendedGatewayFlowRecord.ASPathCount,
			"ASPath":      sample.SFlowExtendedGatewayFlowRecord.ASPath,
			"Communities": sample.SFlowExtendedGatewayFlowRecord.Communities,
			"LocalPref":   sample.SFlowExtendedGatewayFlowRecord.LocalPref,
		}

		//SFlowExtendedUserFlow
		result["SFlowExtendedUserFlow"] = common.MapStr{
			"SFlowBaseFlowRecord": common.MapStr{
				"EnterpriseID":   sample.SFlowExtendedUserFlow.SFlowBaseFlowRecord.EnterpriseID,
				"Format":         sample.SFlowExtendedUserFlow.SFlowBaseFlowRecord.Format,
				"FlowDataLength": sample.SFlowExtendedUserFlow.SFlowBaseFlowRecord.FlowDataLength,
			},
			"SourceCharSet":      sample.SFlowExtendedUserFlow.SourceCharSet,
			"SourceUserID":       sample.SFlowExtendedUserFlow.SourceUserID,
			"DestinationCharSet": sample.SFlowExtendedUserFlow.DestinationCharSet,
			"DestinationUserID":  sample.SFlowExtendedUserFlow.DestinationUserID,
		}
		data, err := json.Marshal(result)
		if err != nil {
			return "", err
		}
		return string(data), nil
	}

	if counter != nil {
		fields := common.MapStr{}
		fields["@timestamp"] = fmt.Sprintf("%s+0800", time.Now().Format("2006-01-02T15:04:05"))
		fields["type"] = "sflow"
		fields["transport"] = "udp"
		fields["src"] = counter.Data.Datagram.SrcMac
		fields["dst"] = counter.Data.Datagram.DstMac
		fields["status"] = common.ERROR_STATUS
		fields["notes"] = "counter"

		sflowEvent := common.MapStr{}
		fields["sflow"] = sflowEvent

		//sflow agent info
		fields["Data"] = common.MapStr{
			"IPLength": 4,
			"SrcIP":    counter.Data.Datagram.SrcIP,
			"DstIP":    counter.Data.Datagram.DstIP,
			"SrcPort":  counter.Data.Datagram.SrcPort,
			"DstPort":  counter.Data.Datagram.DstPort,
		}

		//SFlow info
		// fields["type"] = "counter"
		fields["DatagramVersion"] = counter.Data.DatagramVersion
		fields["AgentAddress"] = counter.Data.AgentAddress
		fields["SubAgentID"] = counter.Data.SubAgentID
		fields["SequenceNumber"] = counter.Data.SequenceNumber
		fields["AgentUptime"] = counter.Data.AgentUptime
		fields["SampleCount"] = counter.Data.SampleCount
		fields["EnterpriseID"] = counter.EnterpriseID
		fields["Format"] = counter.Format
		fields["SampleLength"] = counter.SampleLength
		fields["SequenceNumber"] = counter.SequenceNumber
		fields["SourceIDClass"] = counter.SourceIDClass
		fields["SourceIDIndex"] = counter.SourceIDIndex
		fields["RecordCount"] = counter.RecordCount

		//SFlowGenericInterfaceCounters
		fields["SFlowGenericInterfaceCounters"] = common.MapStr{
			"SFlowBaseFlowRecord": common.MapStr{
				"EnterpriseID":   counter.SFlowGenericInterfaceCounters.SFlowBaseCounterRecord.EnterpriseID,
				"Format":         counter.SFlowGenericInterfaceCounters.SFlowBaseCounterRecord.Format,
				"FlowDataLength": counter.SFlowGenericInterfaceCounters.SFlowBaseCounterRecord.FlowDataLength,
			},
			"IfIndex":            counter.SFlowGenericInterfaceCounters.IfIndex,
			"IfType":             counter.SFlowGenericInterfaceCounters.IfType,
			"IfSpeed":            counter.SFlowGenericInterfaceCounters.IfSpeed,
			"IfDirection":        counter.SFlowGenericInterfaceCounters.IfDirection,
			"IfStatus":           counter.SFlowGenericInterfaceCounters.IfStatus,
			"IfInOctets":         counter.SFlowGenericInterfaceCounters.IfInOctets,
			"IfInUcastPkts":      counter.SFlowGenericInterfaceCounters.IfInUcastPkts,
			"IfInMulticastPkts":  counter.SFlowGenericInterfaceCounters.IfInMulticastPkts,
			"IfInBroadcastPkts":  counter.SFlowGenericInterfaceCounters.IfInBroadcastPkts,
			"IfInDiscards":       counter.SFlowGenericInterfaceCounters.IfInDiscards,
			"IfInErrors":         counter.SFlowGenericInterfaceCounters.IfInErrors,
			"IfInUnknownProtos":  counter.SFlowGenericInterfaceCounters.IfInUnknownProtos,
			"IfOutOctets":        counter.SFlowGenericInterfaceCounters.IfOutOctets,
			"IfOutUcastPkts":     counter.SFlowGenericInterfaceCounters.IfOutUcastPkts,
			"IfOutMulticastPkts": counter.SFlowGenericInterfaceCounters.IfOutMulticastPkts,
			"IfOutBroadcastPkts": counter.SFlowGenericInterfaceCounters.IfOutBroadcastPkts,
			"IfOutDiscards":      counter.SFlowGenericInterfaceCounters.IfOutDiscards,
			"IfOutErrors":        counter.SFlowGenericInterfaceCounters.IfOutErrors,
			"IfPromiscuousMode":  counter.SFlowGenericInterfaceCounters.IfPromiscuousMode,
		}

		//SFlowEthernetCounters
		fields["SFlowEthernetCounters"] = common.MapStr{
			"SFlowBaseFlowRecord": common.MapStr{
				"EnterpriseID":   counter.SFlowEthernetCounters.SFlowBaseCounterRecord.EnterpriseID,
				"Format":         counter.SFlowEthernetCounters.SFlowBaseCounterRecord.Format,
				"FlowDataLength": counter.SFlowEthernetCounters.SFlowBaseCounterRecord.FlowDataLength,
			},
			"AlignmentErrors":           counter.SFlowEthernetCounters.AlignmentErrors,
			"FCSErrors":                 counter.SFlowEthernetCounters.FCSErrors,
			"SingleCollisionFrames":     counter.SFlowEthernetCounters.SingleCollisionFrames,
			"MultipleCollisionFrames":   counter.SFlowEthernetCounters.MultipleCollisionFrames,
			"SQETestErrors":             counter.SFlowEthernetCounters.SQETestErrors,
			"DeferredTransmissions":     counter.SFlowEthernetCounters.DeferredTransmissions,
			"LateCollisions":            counter.SFlowEthernetCounters.LateCollisions,
			"ExcessiveCollisions":       counter.SFlowEthernetCounters.ExcessiveCollisions,
			"InternalMacTransmitErrors": counter.SFlowEthernetCounters.InternalMacTransmitErrors,
			"CarrierSenseErrors":        counter.SFlowEthernetCounters.CarrierSenseErrors,
			"FrameTooLongs":             counter.SFlowEthernetCounters.FrameTooLongs,
			"InternalMacReceiveErrors":  counter.SFlowEthernetCounters.InternalMacReceiveErrors,
			"SymbolErrors":              counter.SFlowEthernetCounters.SymbolErrors,
		}

		//SFlowProcessorCounters
		fields["SFlowProcessorCounters"] = common.MapStr{
			"SFlowBaseFlowRecord": common.MapStr{
				"EnterpriseID":   counter.SFlowProcessorCounters.SFlowBaseCounterRecord.EnterpriseID,
				"Format":         counter.SFlowProcessorCounters.SFlowBaseCounterRecord.Format,
				"FlowDataLength": counter.SFlowProcessorCounters.SFlowBaseCounterRecord.FlowDataLength,
			},
			"FiveSecCpu":  counter.SFlowProcessorCounters.FiveSecCpu,
			"OneMinCpu":   counter.SFlowProcessorCounters.OneMinCpu,
			"FiveMinCpu":  counter.SFlowProcessorCounters.FiveMinCpu,
			"TotalMemory": counter.SFlowProcessorCounters.TotalMemory,
			"FreeMemory":  counter.SFlowProcessorCounters.FreeMemory,
		}
		data, err := json.Marshal(fields)
		if err != nil {
			return "", err
		}
		return string(data), nil
	}
	return "", errors.New("none found")
}
