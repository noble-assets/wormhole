// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: wormhole/v1/genesis.proto

package types

import (
	fmt "fmt"
	_ "github.com/cosmos/gogoproto/gogoproto"
	proto "github.com/cosmos/gogoproto/proto"
	io "io"
	math "math"
	math_bits "math/bits"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

type GenesisState struct {
	Config           Config                 `protobuf:"bytes,1,opt,name=config,proto3" json:"config"`
	WormchainChannel string                 `protobuf:"bytes,2,opt,name=wormchain_channel,json=wormchainChannel,proto3" json:"wormchain_channel,omitempty"`
	GuardianSets     map[uint32]GuardianSet `protobuf:"bytes,3,rep,name=guardian_sets,json=guardianSets,proto3" json:"guardian_sets" protobuf_key:"varint,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	Sequences        map[string]uint64      `protobuf:"bytes,4,rep,name=sequences,proto3" json:"sequences,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"varint,2,opt,name=value,proto3"`
	VaaArchive       [][]byte               `protobuf:"bytes,5,rep,name=vaa_archive,json=vaaArchive,proto3" json:"vaa_archive,omitempty"`
}

func (m *GenesisState) Reset()         { *m = GenesisState{} }
func (m *GenesisState) String() string { return proto.CompactTextString(m) }
func (*GenesisState) ProtoMessage()    {}
func (*GenesisState) Descriptor() ([]byte, []int) {
	return fileDescriptor_51ad1415a98c53ea, []int{0}
}
func (m *GenesisState) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *GenesisState) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_GenesisState.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *GenesisState) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GenesisState.Merge(m, src)
}
func (m *GenesisState) XXX_Size() int {
	return m.Size()
}
func (m *GenesisState) XXX_DiscardUnknown() {
	xxx_messageInfo_GenesisState.DiscardUnknown(m)
}

var xxx_messageInfo_GenesisState proto.InternalMessageInfo

func (m *GenesisState) GetConfig() Config {
	if m != nil {
		return m.Config
	}
	return Config{}
}

func (m *GenesisState) GetWormchainChannel() string {
	if m != nil {
		return m.WormchainChannel
	}
	return ""
}

func (m *GenesisState) GetGuardianSets() map[uint32]GuardianSet {
	if m != nil {
		return m.GuardianSets
	}
	return nil
}

func (m *GenesisState) GetSequences() map[string]uint64 {
	if m != nil {
		return m.Sequences
	}
	return nil
}

func (m *GenesisState) GetVaaArchive() [][]byte {
	if m != nil {
		return m.VaaArchive
	}
	return nil
}

func init() {
	proto.RegisterType((*GenesisState)(nil), "wormhole.v1.GenesisState")
	proto.RegisterMapType((map[uint32]GuardianSet)(nil), "wormhole.v1.GenesisState.GuardianSetsEntry")
	proto.RegisterMapType((map[string]uint64)(nil), "wormhole.v1.GenesisState.SequencesEntry")
}

func init() { proto.RegisterFile("wormhole/v1/genesis.proto", fileDescriptor_51ad1415a98c53ea) }

var fileDescriptor_51ad1415a98c53ea = []byte{
	// 380 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x74, 0x92, 0x4f, 0x4f, 0xf2, 0x40,
	0x10, 0xc6, 0x5b, 0x0a, 0x24, 0x6c, 0xe1, 0x0d, 0xec, 0xcb, 0xa1, 0xf6, 0x50, 0x1a, 0x0f, 0xa6,
	0x09, 0xb1, 0x0d, 0x78, 0x31, 0xc6, 0x83, 0x42, 0x94, 0x7b, 0x49, 0x4c, 0xf4, 0x42, 0x96, 0xba,
	0xb6, 0x8d, 0x65, 0x17, 0xbb, 0xdb, 0x1a, 0xbe, 0x85, 0x1f, 0x8b, 0x23, 0x47, 0x4f, 0xc6, 0xc0,
	0x17, 0x31, 0xdd, 0xf2, 0xa7, 0xc4, 0x78, 0x9b, 0xcc, 0xf3, 0xcc, 0x6f, 0x9e, 0x4c, 0x06, 0x9c,
	0xbc, 0xd3, 0x78, 0x16, 0xd0, 0x08, 0x3b, 0x69, 0xcf, 0xf1, 0x31, 0xc1, 0x2c, 0x64, 0xf6, 0x3c,
	0xa6, 0x9c, 0x42, 0x75, 0x27, 0xd9, 0x69, 0x4f, 0x6f, 0xfb, 0xd4, 0xa7, 0xa2, 0xef, 0x64, 0x55,
	0x6e, 0xd1, 0xf5, 0xe2, 0xf4, 0xde, 0x2e, 0xb4, 0xd3, 0xa5, 0x02, 0xea, 0xa3, 0x1c, 0x38, 0xe6,
	0x88, 0x63, 0xd8, 0x03, 0x55, 0x8f, 0x92, 0x97, 0xd0, 0xd7, 0x64, 0x53, 0xb6, 0xd4, 0xfe, 0x7f,
	0xbb, 0xb0, 0xc0, 0x1e, 0x0a, 0x69, 0x50, 0x5e, 0x7e, 0x75, 0x24, 0x77, 0x6b, 0x84, 0x5d, 0xd0,
	0xca, 0x3c, 0x5e, 0x80, 0x42, 0x32, 0xf1, 0x02, 0x44, 0x08, 0x8e, 0xb4, 0x92, 0x29, 0x5b, 0x35,
	0xb7, 0xb9, 0x17, 0x86, 0x79, 0x1f, 0x3e, 0x80, 0x86, 0x9f, 0xa0, 0xf8, 0x39, 0x44, 0x64, 0xc2,
	0x30, 0x67, 0x9a, 0x62, 0x2a, 0x96, 0xda, 0xef, 0x1e, 0xad, 0x29, 0x26, 0xb2, 0x47, 0x5b, 0xfb,
	0x18, 0x73, 0x76, 0x47, 0x78, 0xbc, 0xd8, 0xae, 0xaf, 0xfb, 0x05, 0x01, 0xde, 0x83, 0x1a, 0xc3,
	0x6f, 0x09, 0x26, 0x1e, 0x66, 0x5a, 0x59, 0x30, 0xad, 0xbf, 0x99, 0xe3, 0x9d, 0x55, 0x00, 0xdd,
	0xc3, 0x28, 0xec, 0x00, 0x35, 0x45, 0x68, 0x82, 0x62, 0x2f, 0x08, 0x53, 0xac, 0x55, 0x4c, 0xc5,
	0xaa, 0xbb, 0x20, 0x45, 0xe8, 0x36, 0xef, 0xe8, 0x8f, 0xa0, 0xf5, 0x2b, 0x11, 0x6c, 0x02, 0xe5,
	0x15, 0x2f, 0xc4, 0xc9, 0x1a, 0x6e, 0x56, 0x42, 0x1b, 0x54, 0x52, 0x14, 0x25, 0x58, 0x1c, 0x42,
	0xed, 0x6b, 0xc7, 0x59, 0x0e, 0x00, 0x37, 0xb7, 0x5d, 0x95, 0x2e, 0x65, 0xfd, 0x1a, 0xfc, 0x3b,
	0x0e, 0x56, 0xe4, 0xd6, 0x72, 0x6e, 0xbb, 0xc8, 0x2d, 0x17, 0xa6, 0x07, 0x37, 0xcb, 0xb5, 0x21,
	0xaf, 0xd6, 0x86, 0xfc, 0xbd, 0x36, 0xe4, 0x8f, 0x8d, 0x21, 0xad, 0x36, 0x86, 0xf4, 0xb9, 0x31,
	0xa4, 0xa7, 0x33, 0x3f, 0xe4, 0x41, 0x32, 0xb5, 0x3d, 0x3a, 0x73, 0x08, 0x9d, 0x46, 0xf8, 0x1c,
	0xb1, 0xec, 0xf6, 0xfb, 0x67, 0x70, 0xf8, 0x62, 0x8e, 0xd9, 0xb4, 0x2a, 0x7e, 0xe2, 0xe2, 0x27,
	0x00, 0x00, 0xff, 0xff, 0x70, 0x2e, 0x82, 0x46, 0x6f, 0x02, 0x00, 0x00,
}

func (m *GenesisState) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *GenesisState) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *GenesisState) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.VaaArchive) > 0 {
		for iNdEx := len(m.VaaArchive) - 1; iNdEx >= 0; iNdEx-- {
			i -= len(m.VaaArchive[iNdEx])
			copy(dAtA[i:], m.VaaArchive[iNdEx])
			i = encodeVarintGenesis(dAtA, i, uint64(len(m.VaaArchive[iNdEx])))
			i--
			dAtA[i] = 0x2a
		}
	}
	if len(m.Sequences) > 0 {
		for k := range m.Sequences {
			v := m.Sequences[k]
			baseI := i
			i = encodeVarintGenesis(dAtA, i, uint64(v))
			i--
			dAtA[i] = 0x10
			i -= len(k)
			copy(dAtA[i:], k)
			i = encodeVarintGenesis(dAtA, i, uint64(len(k)))
			i--
			dAtA[i] = 0xa
			i = encodeVarintGenesis(dAtA, i, uint64(baseI-i))
			i--
			dAtA[i] = 0x22
		}
	}
	if len(m.GuardianSets) > 0 {
		for k := range m.GuardianSets {
			v := m.GuardianSets[k]
			baseI := i
			{
				size, err := (&v).MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintGenesis(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x12
			i = encodeVarintGenesis(dAtA, i, uint64(k))
			i--
			dAtA[i] = 0x8
			i = encodeVarintGenesis(dAtA, i, uint64(baseI-i))
			i--
			dAtA[i] = 0x1a
		}
	}
	if len(m.WormchainChannel) > 0 {
		i -= len(m.WormchainChannel)
		copy(dAtA[i:], m.WormchainChannel)
		i = encodeVarintGenesis(dAtA, i, uint64(len(m.WormchainChannel)))
		i--
		dAtA[i] = 0x12
	}
	{
		size, err := m.Config.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintGenesis(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0xa
	return len(dAtA) - i, nil
}

func encodeVarintGenesis(dAtA []byte, offset int, v uint64) int {
	offset -= sovGenesis(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *GenesisState) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = m.Config.Size()
	n += 1 + l + sovGenesis(uint64(l))
	l = len(m.WormchainChannel)
	if l > 0 {
		n += 1 + l + sovGenesis(uint64(l))
	}
	if len(m.GuardianSets) > 0 {
		for k, v := range m.GuardianSets {
			_ = k
			_ = v
			l = v.Size()
			mapEntrySize := 1 + sovGenesis(uint64(k)) + 1 + l + sovGenesis(uint64(l))
			n += mapEntrySize + 1 + sovGenesis(uint64(mapEntrySize))
		}
	}
	if len(m.Sequences) > 0 {
		for k, v := range m.Sequences {
			_ = k
			_ = v
			mapEntrySize := 1 + len(k) + sovGenesis(uint64(len(k))) + 1 + sovGenesis(uint64(v))
			n += mapEntrySize + 1 + sovGenesis(uint64(mapEntrySize))
		}
	}
	if len(m.VaaArchive) > 0 {
		for _, b := range m.VaaArchive {
			l = len(b)
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	return n
}

func sovGenesis(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozGenesis(x uint64) (n int) {
	return sovGenesis(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *GenesisState) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowGenesis
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: GenesisState: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: GenesisState: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Config", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Config.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field WormchainChannel", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.WormchainChannel = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field GuardianSets", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.GuardianSets == nil {
				m.GuardianSets = make(map[uint32]GuardianSet)
			}
			var mapkey uint32
			mapvalue := &GuardianSet{}
			for iNdEx < postIndex {
				entryPreIndex := iNdEx
				var wire uint64
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowGenesis
					}
					if iNdEx >= l {
						return io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					wire |= uint64(b&0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				fieldNum := int32(wire >> 3)
				if fieldNum == 1 {
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowGenesis
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						mapkey |= uint32(b&0x7F) << shift
						if b < 0x80 {
							break
						}
					}
				} else if fieldNum == 2 {
					var mapmsglen int
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowGenesis
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						mapmsglen |= int(b&0x7F) << shift
						if b < 0x80 {
							break
						}
					}
					if mapmsglen < 0 {
						return ErrInvalidLengthGenesis
					}
					postmsgIndex := iNdEx + mapmsglen
					if postmsgIndex < 0 {
						return ErrInvalidLengthGenesis
					}
					if postmsgIndex > l {
						return io.ErrUnexpectedEOF
					}
					mapvalue = &GuardianSet{}
					if err := mapvalue.Unmarshal(dAtA[iNdEx:postmsgIndex]); err != nil {
						return err
					}
					iNdEx = postmsgIndex
				} else {
					iNdEx = entryPreIndex
					skippy, err := skipGenesis(dAtA[iNdEx:])
					if err != nil {
						return err
					}
					if (skippy < 0) || (iNdEx+skippy) < 0 {
						return ErrInvalidLengthGenesis
					}
					if (iNdEx + skippy) > postIndex {
						return io.ErrUnexpectedEOF
					}
					iNdEx += skippy
				}
			}
			m.GuardianSets[mapkey] = *mapvalue
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Sequences", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Sequences == nil {
				m.Sequences = make(map[string]uint64)
			}
			var mapkey string
			var mapvalue uint64
			for iNdEx < postIndex {
				entryPreIndex := iNdEx
				var wire uint64
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowGenesis
					}
					if iNdEx >= l {
						return io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					wire |= uint64(b&0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				fieldNum := int32(wire >> 3)
				if fieldNum == 1 {
					var stringLenmapkey uint64
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowGenesis
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						stringLenmapkey |= uint64(b&0x7F) << shift
						if b < 0x80 {
							break
						}
					}
					intStringLenmapkey := int(stringLenmapkey)
					if intStringLenmapkey < 0 {
						return ErrInvalidLengthGenesis
					}
					postStringIndexmapkey := iNdEx + intStringLenmapkey
					if postStringIndexmapkey < 0 {
						return ErrInvalidLengthGenesis
					}
					if postStringIndexmapkey > l {
						return io.ErrUnexpectedEOF
					}
					mapkey = string(dAtA[iNdEx:postStringIndexmapkey])
					iNdEx = postStringIndexmapkey
				} else if fieldNum == 2 {
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowGenesis
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						mapvalue |= uint64(b&0x7F) << shift
						if b < 0x80 {
							break
						}
					}
				} else {
					iNdEx = entryPreIndex
					skippy, err := skipGenesis(dAtA[iNdEx:])
					if err != nil {
						return err
					}
					if (skippy < 0) || (iNdEx+skippy) < 0 {
						return ErrInvalidLengthGenesis
					}
					if (iNdEx + skippy) > postIndex {
						return io.ErrUnexpectedEOF
					}
					iNdEx += skippy
				}
			}
			m.Sequences[mapkey] = mapvalue
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field VaaArchive", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.VaaArchive = append(m.VaaArchive, make([]byte, postIndex-iNdEx))
			copy(m.VaaArchive[len(m.VaaArchive)-1], dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipGenesis(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthGenesis
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipGenesis(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowGenesis
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthGenesis
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupGenesis
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthGenesis
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthGenesis        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowGenesis          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupGenesis = fmt.Errorf("proto: unexpected end of group")
)
