package pkg

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"

	log "github.com/sirupsen/logrus"
)

type HeaderV5 struct {
	Version          uint16 `json:"version" codec:"version"`
	FlowRecords      uint16 `json:"flow_records" codec:"flow_records"`
	Uptime           uint32 `json:"uptime" codec:"uptime"`
	UnixSec          uint32 `json:"unix_sec" codec:"unix_sec"`
	UnixNsec         uint32 `json:"unix_nsec" codec:"unix_nsec"`
	FlowSeqNum       uint32 `json:"flow_seq_num" codec:"flow_seq_num"`
	EngineType       uint8  `json:"engine_type" codec:"engine_type"`
	EngineId         uint8  `json:"engine_id" codec:"engine_id"`
	SamplingInterval uint16 `json:"sampling_interval" codec:"sampling_interval"`
}

type RecordBaseV5 struct {
	InputSnmp     uint16 `json:"input_snmp" codec:"input_snmp"`
	OutputSnmp    uint16 `json:"output_snmp" codec:"output_snmp"`
	InPkts        uint32 `json:"in_pkts" codec:"in_pkts"`
	InBytes       uint32 `json:"in_bytes" codec:"in_bytes"`
	FirstSwitched uint32 `json:"first_switched" codec:"first_switched"`
	LastSwitched  uint32 `json:"last_switched" codec:"last_switched"`
	L4SrcPort     uint16 `json:"l4_src_port" codec:"l4_src_port"`
	L4DstPort     uint16 `json:"l4_dst_port" codec:"l4_dst_port"`
	_             uint8
	TcpFlags      uint8  `json:"tcp_flags" codec:"tcp_flags"`
	Protocol      uint8  `json:"protocol" codec:"protocol"`
	SrcTos        uint8  `json:"src_tos" codec:"src_tos"`
	SrcAs         uint16 `json:"src_as" codec:"src_as"`
	DstAs         uint16 `json:"dst_as" codec:"dst_as"`
	SrcMask       uint8  `json:"src_mask" codec:"src_mask"`
	DstMask       uint8  `json:"dst_mask" codec:"dst_mask"`
	_             uint16
}

type BinaryRecordV5 struct {
	Ipv4SrcAddrInt uint32 `json:"-" codec:"-"`
	Ipv4DstAddrInt uint32 `json:"-" codec:"-"`
	Ipv4NextHopInt uint32 `json:"-" codec:"-"`

	RecordBaseV5
}

type NetFlowV5 struct {
	HeaderV5
	BinaryRecordV5

	Host              string `json:"host" codec:"host"`
	SamplingAlgorithm uint8  `json:"sampling_algorithm" codec:"sampling_algorithm"`
	Ipv4SrcAddr       string `json:"ipv4_src_addr" codec:"ipv4_src_addr"`
	Ipv4DstAddr       string `json:"ipv4_dst_addr" codec:"ipv4_dst_addr"`
	Ipv4NextHop       string `json:"ipv4_next_hop" codec:"ipv4_next_hop"`
}

func (this *NetFlowV5) DecodeNetFlowV5(header *HeaderV5, binRecord *BinaryRecordV5, ip string) NetFlowV5 {
	netflow := NetFlowV5{
		Host:           ip,
		HeaderV5:       *header,
		BinaryRecordV5: *binRecord,
		Ipv4SrcAddr:    this.IntToIPv4Addr(binRecord.Ipv4SrcAddrInt).String(),
		Ipv4DstAddr:    this.IntToIPv4Addr(binRecord.Ipv4DstAddrInt).String(),
		Ipv4NextHop:    this.IntToIPv4Addr(binRecord.Ipv4NextHopInt).String(),
	}
	//Modify sampling settings
	netflow.SamplingAlgorithm = uint8(0x3 & (netflow.SamplingInterval >> 14))
	netflow.SamplingInterval = 0x3fff & netflow.SamplingInterval
	return netflow
}

func (this *NetFlowV5) IntToIPv4Addr(intAddr uint32) net.IP {
	return net.IPv4(
		byte(intAddr>>24),
		byte(intAddr>>16),
		byte(intAddr>>8),
		byte(intAddr),
	)
}

//func (this *NetFlowV5) PayLoadToNetFlowV5(data []byte, host string) []NetFlowV5 {
func (this *NetFlowV5) PayLoadToNetFlowV5(data []byte, host string) []string {
	//log.Critical(data)
	//datas := []NetFlowV5{}

	result := []string{}
	header := HeaderV5{}
	buf := bytes.NewBuffer(data)
	err := binary.Read(buf, binary.BigEndian, &header)
	if err != nil {
		log.Error("Error:", err)
		return nil
	}
	if header.Version == 5 {
		//NETFLOW5 解析
		//log.Error("NETFLOW5 Deteced")
		for i := 0; i < int(header.FlowRecords); i++ {
			record := BinaryRecordV5{}
			err = binary.Read(buf, binary.BigEndian, &record)
			if err != nil {
				log.Error(fmt.Printf("binary.Read failed: %v\n", err))
				return nil
			}
			//ip解析
			//src := this.IntToIPv4Addr(record.Ipv4SrcAddrInt).String()
			//log.Error(fmt.Sprintf("srouce %s %s", src, curl.HttpGet("http://127.0.0.1:8080/check/"+src)))
			//log.Error(fmt.Sprintf("srouce %s ", src),geoip.ParseStringIp(src))
			//dst := this.IntToIPv4Addr(record.Ipv4DstAddrInt).String()
			//log.Error(fmt.Sprintf("srouce %s %s", dst, curl.HttpGet("http://127.0.0.1:8080/check/"+dst)))
			//log.Error(fmt.Sprintf("srouce %s ", dst),geoip.ParseStringIp(dst))
			netflow := this.DecodeNetFlowV5(&header, &record, host)
			bufs, err := json.Marshal(netflow)
			if err != nil {
				log.Error(err)
				return nil
			}
			//fmt.Println(string(bufs))

			result = append(result, string(bufs))

			//datas = append(datas, netflow)
			//log.Informational(fmt.Sprintf("%v\n",string(bufs)))
		}
		//return datas
		return result
	} else {
		log.Error("Netflow version want 5 got %d", header.Version)
		return nil
	}
}
