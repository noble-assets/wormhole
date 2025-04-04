// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: wormhole/v1/wormhole.proto

package types

import (
	fmt "fmt"
	_ "github.com/cosmos/cosmos-sdk/types/tx/amino"
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

// Config is an object used to store the module configuration.
type Config struct {
	ChainId          uint16 `protobuf:"varint,1,opt,name=chain_id,json=chainId,proto3,customtype=uint16" json:"chain_id"`
	GuardianSetIndex uint32 `protobuf:"varint,2,opt,name=guardian_set_index,json=guardianSetIndex,proto3" json:"guardian_set_index,omitempty"`
	GovChain         uint16 `protobuf:"varint,3,opt,name=gov_chain,json=govChain,proto3,customtype=uint16" json:"gov_chain"`
	GovAddress       []byte `protobuf:"bytes,4,opt,name=gov_address,json=govAddress,proto3" json:"gov_address,omitempty"`
}

func (m *Config) Reset()         { *m = Config{} }
func (m *Config) String() string { return proto.CompactTextString(m) }
func (*Config) ProtoMessage()    {}
func (*Config) Descriptor() ([]byte, []int) {
	return fileDescriptor_98ced89aa7840992, []int{0}
}
func (m *Config) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Config) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Config.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Config) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Config.Merge(m, src)
}
func (m *Config) XXX_Size() int {
	return m.Size()
}
func (m *Config) XXX_DiscardUnknown() {
	xxx_messageInfo_Config.DiscardUnknown(m)
}

var xxx_messageInfo_Config proto.InternalMessageInfo

func (m *Config) GetGuardianSetIndex() uint32 {
	if m != nil {
		return m.GuardianSetIndex
	}
	return 0
}

func (m *Config) GetGovAddress() []byte {
	if m != nil {
		return m.GovAddress
	}
	return nil
}

// GuardianSet is an object used to store a specific guardian set in state.
type GuardianSet struct {
	Addresses      [][]byte `protobuf:"bytes,1,rep,name=addresses,proto3" json:"addresses,omitempty"`
	ExpirationTime uint64   `protobuf:"varint,2,opt,name=expiration_time,json=expirationTime,proto3" json:"expiration_time,omitempty"`
}

func (m *GuardianSet) Reset()         { *m = GuardianSet{} }
func (m *GuardianSet) String() string { return proto.CompactTextString(m) }
func (*GuardianSet) ProtoMessage()    {}
func (*GuardianSet) Descriptor() ([]byte, []int) {
	return fileDescriptor_98ced89aa7840992, []int{1}
}
func (m *GuardianSet) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *GuardianSet) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_GuardianSet.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *GuardianSet) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GuardianSet.Merge(m, src)
}
func (m *GuardianSet) XXX_Size() int {
	return m.Size()
}
func (m *GuardianSet) XXX_DiscardUnknown() {
	xxx_messageInfo_GuardianSet.DiscardUnknown(m)
}

var xxx_messageInfo_GuardianSet proto.InternalMessageInfo

func (m *GuardianSet) GetAddresses() [][]byte {
	if m != nil {
		return m.Addresses
	}
	return nil
}

func (m *GuardianSet) GetExpirationTime() uint64 {
	if m != nil {
		return m.ExpirationTime
	}
	return 0
}

func init() {
	proto.RegisterType((*Config)(nil), "wormhole.v1.Config")
	proto.RegisterType((*GuardianSet)(nil), "wormhole.v1.GuardianSet")
}

func init() { proto.RegisterFile("wormhole/v1/wormhole.proto", fileDescriptor_98ced89aa7840992) }

var fileDescriptor_98ced89aa7840992 = []byte{
	// 335 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x6c, 0x91, 0xb1, 0x4e, 0x02, 0x31,
	0x1c, 0xc6, 0xaf, 0x82, 0x08, 0x05, 0x51, 0x1b, 0x87, 0x0b, 0x31, 0x07, 0x61, 0x30, 0xa8, 0x91,
	0x0b, 0x21, 0x71, 0x56, 0x18, 0x0c, 0xeb, 0xe9, 0xa4, 0xc3, 0xa5, 0x70, 0xb5, 0x34, 0xe1, 0xfa,
	0x27, 0xd7, 0x72, 0xe2, 0x5b, 0xf8, 0x18, 0x8e, 0x4e, 0x3e, 0x03, 0x23, 0xa3, 0x71, 0x20, 0x06,
	0x06, 0x5f, 0xc3, 0x5c, 0x11, 0x6e, 0x71, 0x69, 0xbe, 0x7e, 0xbf, 0xaf, 0xfd, 0x7f, 0x69, 0x71,
	0xe5, 0x19, 0xa2, 0x70, 0x08, 0x23, 0xe6, 0xc6, 0x2d, 0x77, 0xa3, 0x9b, 0xe3, 0x08, 0x34, 0x90,
	0xe2, 0x76, 0x1f, 0xb7, 0x2a, 0x47, 0x34, 0x14, 0x12, 0x5c, 0xb3, 0xae, 0x79, 0xe5, 0x98, 0x03,
	0x07, 0x23, 0xdd, 0x44, 0xad, 0xdd, 0xfa, 0x07, 0xc2, 0xb9, 0x2e, 0xc8, 0x27, 0xc1, 0xc9, 0x19,
	0xce, 0x0f, 0x86, 0x54, 0x48, 0x5f, 0x04, 0x36, 0xaa, 0xa1, 0xc6, 0x7e, 0xa7, 0x3c, 0x5b, 0x54,
	0xad, 0xaf, 0x45, 0x35, 0x37, 0x11, 0x52, 0xb7, 0xae, 0xbc, 0x3d, 0xc3, 0x7b, 0x01, 0x69, 0x63,
	0xc2, 0x27, 0x34, 0x0a, 0x04, 0x95, 0xbe, 0x62, 0xda, 0x17, 0x32, 0x60, 0x53, 0x7b, 0xc7, 0x1c,
	0xda, 0x7d, 0xfb, 0x79, 0x3f, 0x47, 0xde, 0xe1, 0x26, 0x70, 0xc7, 0x74, 0x2f, 0xc1, 0xe4, 0x02,
	0x17, 0x38, 0xc4, 0xbe, 0xb9, 0xc3, 0xce, 0xfc, 0x3b, 0x20, 0xcf, 0x21, 0xee, 0x26, 0x9c, 0x54,
	0x71, 0x31, 0x09, 0xd3, 0x20, 0x88, 0x98, 0x52, 0x76, 0xb6, 0x86, 0x1a, 0x25, 0x0f, 0x73, 0x88,
	0x6f, 0xd6, 0x4e, 0xfd, 0x11, 0x17, 0x6f, 0xd3, 0x09, 0xe4, 0x04, 0x17, 0xfe, 0xb2, 0x4c, 0xd9,
	0xa8, 0x96, 0x69, 0x94, 0xbc, 0xd4, 0x20, 0x4d, 0x7c, 0xc0, 0xa6, 0x63, 0x11, 0x51, 0x2d, 0x40,
	0xfa, 0x5a, 0x84, 0xcc, 0x94, 0xcd, 0x6e, 0xca, 0x96, 0x53, 0x7a, 0x2f, 0x42, 0xd6, 0xb9, 0x9e,
	0x2d, 0x1d, 0x34, 0x5f, 0x3a, 0xe8, 0x7b, 0xe9, 0xa0, 0xd7, 0x95, 0x63, 0xcd, 0x57, 0x8e, 0xf5,
	0xb9, 0x72, 0xac, 0x87, 0x53, 0x2e, 0xf4, 0x70, 0xd2, 0x6f, 0x0e, 0x20, 0x74, 0x25, 0xf4, 0x47,
	0xec, 0x92, 0x2a, 0xc5, 0xb4, 0xda, 0xfe, 0x86, 0xab, 0x5f, 0xc6, 0x4c, 0xf5, 0x73, 0xe6, 0x79,
	0xdb, 0xbf, 0x01, 0x00, 0x00, 0xff, 0xff, 0x1d, 0x26, 0x40, 0xf6, 0xb2, 0x01, 0x00, 0x00,
}

func (m *Config) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Config) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Config) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.GovAddress) > 0 {
		i -= len(m.GovAddress)
		copy(dAtA[i:], m.GovAddress)
		i = encodeVarintWormhole(dAtA, i, uint64(len(m.GovAddress)))
		i--
		dAtA[i] = 0x22
	}
	if m.GovChain != 0 {
		i = encodeVarintWormhole(dAtA, i, uint64(m.GovChain))
		i--
		dAtA[i] = 0x18
	}
	if m.GuardianSetIndex != 0 {
		i = encodeVarintWormhole(dAtA, i, uint64(m.GuardianSetIndex))
		i--
		dAtA[i] = 0x10
	}
	if m.ChainId != 0 {
		i = encodeVarintWormhole(dAtA, i, uint64(m.ChainId))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *GuardianSet) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *GuardianSet) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *GuardianSet) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.ExpirationTime != 0 {
		i = encodeVarintWormhole(dAtA, i, uint64(m.ExpirationTime))
		i--
		dAtA[i] = 0x10
	}
	if len(m.Addresses) > 0 {
		for iNdEx := len(m.Addresses) - 1; iNdEx >= 0; iNdEx-- {
			i -= len(m.Addresses[iNdEx])
			copy(dAtA[i:], m.Addresses[iNdEx])
			i = encodeVarintWormhole(dAtA, i, uint64(len(m.Addresses[iNdEx])))
			i--
			dAtA[i] = 0xa
		}
	}
	return len(dAtA) - i, nil
}

func encodeVarintWormhole(dAtA []byte, offset int, v uint64) int {
	offset -= sovWormhole(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *Config) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.ChainId != 0 {
		n += 1 + sovWormhole(uint64(m.ChainId))
	}
	if m.GuardianSetIndex != 0 {
		n += 1 + sovWormhole(uint64(m.GuardianSetIndex))
	}
	if m.GovChain != 0 {
		n += 1 + sovWormhole(uint64(m.GovChain))
	}
	l = len(m.GovAddress)
	if l > 0 {
		n += 1 + l + sovWormhole(uint64(l))
	}
	return n
}

func (m *GuardianSet) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.Addresses) > 0 {
		for _, b := range m.Addresses {
			l = len(b)
			n += 1 + l + sovWormhole(uint64(l))
		}
	}
	if m.ExpirationTime != 0 {
		n += 1 + sovWormhole(uint64(m.ExpirationTime))
	}
	return n
}

func sovWormhole(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozWormhole(x uint64) (n int) {
	return sovWormhole(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *Config) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowWormhole
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
			return fmt.Errorf("proto: Config: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Config: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field ChainId", wireType)
			}
			m.ChainId = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowWormhole
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.ChainId |= uint16(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field GuardianSetIndex", wireType)
			}
			m.GuardianSetIndex = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowWormhole
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.GuardianSetIndex |= uint32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field GovChain", wireType)
			}
			m.GovChain = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowWormhole
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.GovChain |= uint16(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field GovAddress", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowWormhole
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
				return ErrInvalidLengthWormhole
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthWormhole
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.GovAddress = append(m.GovAddress[:0], dAtA[iNdEx:postIndex]...)
			if m.GovAddress == nil {
				m.GovAddress = []byte{}
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipWormhole(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthWormhole
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
func (m *GuardianSet) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowWormhole
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
			return fmt.Errorf("proto: GuardianSet: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: GuardianSet: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Addresses", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowWormhole
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
				return ErrInvalidLengthWormhole
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthWormhole
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Addresses = append(m.Addresses, make([]byte, postIndex-iNdEx))
			copy(m.Addresses[len(m.Addresses)-1], dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field ExpirationTime", wireType)
			}
			m.ExpirationTime = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowWormhole
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.ExpirationTime |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipWormhole(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthWormhole
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
func skipWormhole(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowWormhole
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
					return 0, ErrIntOverflowWormhole
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
					return 0, ErrIntOverflowWormhole
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
				return 0, ErrInvalidLengthWormhole
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupWormhole
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthWormhole
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthWormhole        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowWormhole          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupWormhole = fmt.Errorf("proto: unexpected end of group")
)
