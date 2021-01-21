package pkg

const mapping = `
{
  "mappings": {
    "sflow": {
      "properties": {
        "@timestamp": {
          "type": "date"
        },
        "Datagram": {
            "properties": {
                "IPLength": {
                    "ignore_above": 1024,
                    "type": "keyword"
                },
                "SrcIP": {
                    "ignore_above": 1024,
                    "type": "keyword"
                },
                "DstIP": {
                    "ignore_above": 1024,
                    "type": "keyword"
                },
                "SrcPort": {
                    "type": "long"
                },
                "DstPort": {
                    "type": "keyword"
                }
            }
        },
        "DatagramVersion": {
            "ignore_above": 1024,
            "type": "keyword"
          },
        "AgentAddress": {
            "ignore_above": 1024,
            "type": "keyword"
          },
        "SubAgentID": {
            "ignore_above": 1024,
            "type": "keyword"
          },
        "SequenceNumber": {
            "type": "long"
          },
        "AgentUptime": {
            "type": "long"
          },
        "SampleCount": {
            "type": "long"
          },
        "EnterpriseID": {
            "type": "keyword"
          },
        "Format": {
            "ignore_above": 1024,
            "type": "keyword"
          },
        "SourceIDClass": {
            "type": "keyword"
          },
        "SourceIDIndex": {
            "type": "long"
          },
        "SamplingRate": {
            "type": "long"
          },
        "SamplePool": {
            "type": "long"
          },
        "Dropped": {
            "type": "long"
          },
        "InputInterfaceFormat": {
            "ignore_above": 1024,
            "type": "keyword"
          },
        "InputInterface": {
            "type": "long"
          },
        "OutputInterfaceFormat": {
            "ignore_above": 1024,
            "type": "keyword"
          },
        "OutputInterface": {
            "type": "long"
          },
        "RecordCount": {
            "type": "long"
          },
        "SampleLength": {
            "type": "long"
          },

        "SFlowRawPacketFlowRecord": {
          "properties": {
            "SFlowBaseFlowRecord": {
              "properties": {
                "EnterpriseID": {
                  "ignore_above": 1024,
                  "type": "keyword"
                },
                "Format": {
                  "ignore_above": 1024,
                  "type": "keyword"
                },
                "FlowDataLength": {
                  "type": "long"
                }
              }
            },
            "Header": {
              "properties": {
                "FlowRecords": {
                  "type": "long"
                },
                "Packets": {
                  "type": "long"
                },
                "Bytes": {
                  "type": "long"
                },
                "RateBytes": {
                  "type": "long"
                },
                "SrcMac": {
                  "ignore_above": 1024,
                  "type": "keyword"
                },
                "DstMac": {
                  "ignore_above": 1024,
                  "type": "keyword"
                },
                "SrcIP": {
                  "ignore_above": 1024,
                  "type": "keyword"
                },
                "DstIP": {
                  "ignore_above": 1024,
                  "type": "keyword"
                },
                "Ipv4_version": {
                  "type": "long"
                },
                "Ipv4_ihl": {
                  "type": "long"
                },
                "Ipv4_tos": {
                  "type": "long"
                },
                "Ipv4_ttl": {
                  "type": "long"
                },
                "Ipv4_protocol": {
                  "ignore_above": 1024,
                  "type": "keyword"
                },
                "SrcPort": {
                  "type": "keyword"
                },
                "DstPort": {
                  "type": "keyword"
                }
              }
            },
            "HeaderProtocol": {
              "ignore_above": 1024,
              "type": "keyword"
            },
            "FrameLength": {
              "type": "long"
            },
            "PayloadRemoved": {
              "type": "long"
            },
            "HeaderLength": {
              "type": "long"
            }
          }
        },
        "SFlowExtendedSwitchFlowRecord": {
          "properties": {
            "SFlowBaseFlowRecord": {
              "properties": {
                "EnterpriseID": {
                  "ignore_above": 1024,
                  "type": "keyword"
                },
                "Format": {
                  "ignore_above": 1024,
                  "type": "keyword"
                },
                "FlowDataLength": {
                  "type": "long"
                }
              }
            },
            "IncomingVLAN": {
              "type": "long"
            },
            "IncomingVLANPriority": {
              "type": "long"
            },
            "OutgoingVLAN": {
              "type": "long"
            },
            "OutgoingVLAN": {
              "type": "long"
            }
          }
        },
        "SFlowExtendedRouterFlowRecord": {
          "properties": {
            "SFlowBaseFlowRecord": {
              "properties": {
                "EnterpriseID": {
                  "ignore_above": 1024,
                  "type": "keyword"
                },
                "Format": {
                  "ignore_above": 1024,
                  "type": "keyword"
                },
                "FlowDataLength": {
                  "type": "long"
                }
              }
            },
            "NextHop": {
              "ignore_above": 1024,
              "type": "keyword"
            },
            "NextHopSourceMask": {
              "ignore_above": 1024,
              "type": "keyword"
            },
            "NextHopDestinationMask": {
              "ignore_above": 1024,
              "type": "keyword"
            }
          }
        },
        "SFlowExtendedGatewayFlowRecord": {
          "properties": {
            "SFlowBaseFlowRecord": {
              "properties": {
                "EnterpriseID": {
                  "ignore_above": 1024,
                  "type": "keyword"
                },
                "Format": {
                  "ignore_above": 1024,
                  "type": "keyword"
                },
                "FlowDataLength": {
                  "type": "long"
                }
              }
            },
            "NextHop": {
              "ignore_above": 1024,
              "type": "keyword"
            },
            "AS": {
              "type": "long"
            },
            "SourceAS": {
              "type": "long"
            },
            "PeerAS": {
              "type": "long"
            },
            "ASPathCount": {
              "type": "long"
            },
            "ASPath": {
              "ignore_above": 1024,
              "type": "keyword"
            },
            "Communities": {
              "ignore_above": 1024,
              "type": "keyword"
            },
            "LocalPref": {
              "type": "long"
            }
          }
        },
        "SFlowExtendedUserFlow": {
          "properties": {
            "SFlowBaseFlowRecord": {
              "properties": {
                "EnterpriseID": {
                  "ignore_above": 1024,
                  "type": "keyword"
                },
                "Format": {
                  "ignore_above": 1024,
                  "type": "keyword"
                },
                "FlowDataLength": {
                  "type": "long"
                }
              }
            },
            "SourceCharSet": {
              "ignore_above": 1024,
              "type": "keyword"
            },
            "SourceUserID": {
              "ignore_above": 1024,
              "type": "keyword"
            },
            "DestinationCharSet": {
              "ignore_above": 1024,
              "type": "keyword"
            },
            "DestinationUserID": {
              "ignore_above": 1024,
              "type": "keyword"
            }
          }
        },
        "SFlowGenericInterfaceCounters": {
          "properties": {
            "SFlowBaseFlowRecord": {
              "properties": {
                "EnterpriseID": {
                  "ignore_above": 1024,
                  "type": "keyword"
                },
                "Format": {
                  "ignore_above": 1024,
                  "type": "keyword"
                },
                "FlowDataLength": {
                  "type": "long"
                }
              }
            },
            "IfIndex": {
              "type": "long"
            },
            "IfType": {
              "type": "long"
            },
            "IfSpeed": {
              "type": "long"
            },
            "IfDirection": {
              "type": "long"
            },
            "IfStatus": {
              "type": "long"
            },
            "IfInOctets": {
              "type": "long"
            },
            "IfInUcastPkts": {
              "type": "long"
            },
            "IfInMulticastPkts": {
              "type": "long"
            },
            "IfInBroadcastPkts": {
              "type": "long"
            },
            "IfInDiscards": {
              "type": "long"
            },
            "IfInErrors": {
              "type": "long"
            },
            "IfInUnknownProtos": {
              "type": "long"
            },
            "IfOutOctets": {
              "type": "long"
            },
            "IfOutUcastPkts": {
              "type": "long"
            },
            "IfOutMulticastPkts": {
              "type": "long"
            },
            "IfOutDiscards": {
              "type": "long"
            },
            "IfOutErrors": {
              "type": "long"
            },
            "IfPromiscuousMode": {
              "type": "long"
            }
          }
        },
        "SFlowEthernetCounters": {
          "properties": {
            "SFlowBaseFlowRecord": {
              "properties": {
                "EnterpriseID": {
                  "ignore_above": 1024,
                  "type": "keyword"
                },
                "Format": {
                  "ignore_above": 1024,
                  "type": "keyword"
                },
                "FlowDataLength": {
                  "type": "long"
                }
              }
            },
            "AlignmentErrors": {
              "type": "long"
            },
            "FCSErrors": {
              "type": "long"
            },
            "SingleCollisionFrames": {
              "type": "long"
            },
            "MultipleCollisionFrames": {
              "type": "long"
            },
            "SQETestErrors": {
              "type": "long"
            },
            "DeferredTransmissions": {
              "type": "long"
            },
            "LateCollisions": {
              "type": "long"
            },
            "ExcessiveCollisions": {
              "type": "long"
            },
            "InternalMacTransmitErrors": {
              "type": "long"
            },
            "CarrierSenseErrors": {
              "type": "long"
            },
            "FrameTooLongs": {
              "type": "long"
            },
            "InternalMacReceiveErrors": {
              "type": "long"
            },
            "SymbolErrors": {
              "type": "long"
            }
          }
        },
        "SFlowProcessorCounters": {
          "properties": {
            "SFlowBaseFlowRecord": {
              "properties": {
                "EnterpriseID": {
                  "ignore_above": 1024,
                  "type": "keyword"
                },
                "Format": {
                  "ignore_above": 1024,
                  "type": "keyword"
                },
                "FlowDataLength": {
                  "type": "long"
                }
              }
            },
            "FiveSecCpu": {
              "type": "long"
            },
            "OneMinCpu": {
              "type": "long"
            },
            "FiveMinCpu": {
              "type": "long"
            },
            "TotalMemory": {
              "type": "long"
            },
            "FreeMemory": {
              "type": "long"
            }
          }
        }
      }
    }
  },
  "settings": {
    "index.mapping.total_fields.limit": 10000,
    "index.refresh_interval": "5s",
    "number_of_shards": 1,
    "number_of_replicas": 0
  }
}`
