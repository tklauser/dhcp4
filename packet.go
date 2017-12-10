package dhcp4

import (
	"net"
	"strings"

	"github.com/u-root/dhcp4/util"
)

const (
	minPacketLen = 236
	MAClen       = 16
)

// Packet is a DHCPv4 packet as described in RFC 2131 Section 2.
type Packet struct {
	Op            OpCode
	HType         uint8
	Hops          uint8
	TransactionID [4]byte
	Secs          uint16
	Broadcast     bool

	// Client IP address.
	CIAddr net.IP

	// Your IP address.
	YIAddr net.IP

	// Server IP address.
	SIAddr net.IP

	// Gateway IP address.
	GIAddr net.IP

	// Client hardware address.
	CHAddr net.HardwareAddr

	ServerName string
	BootFile   string
	Options    Options
}

func NewPacket(op OpCode) *Packet {
	return &Packet{
		Op:    op,
		HType: 1, /* ethernet */
		Options: Options{
			End: []byte{},
		},
	}
}

func (p *Packet) writeIP(b *util.Buffer, ip net.IP) {
	var zeros [net.IPv4len]byte
	if ip == nil {
		b.WriteBytes(zeros[:])
	} else {
		b.WriteBytes(ip[:net.IPv4len])
	}
}

func (p *Packet) MarshalBinary() ([]byte, error) {
	b := util.NewBuffer(make([]byte, 0, minPacketLen))
	b.Write8(uint8(p.Op))
	b.Write8(p.HType)
	// HLen
	b.Write8(uint8(len(p.CHAddr)))
	b.Write8(p.Hops)
	b.WriteBytes(p.TransactionID[:])
	b.Write16(p.Secs)

	var flags uint16
	if p.Broadcast {
		flags |= 1 << 15
	}
	b.Write16(flags)

	p.writeIP(b, p.CIAddr)
	p.writeIP(b, p.YIAddr)
	p.writeIP(b, p.SIAddr)
	p.writeIP(b, p.GIAddr)
	b.WriteBytes(p.CHAddr)

	var sname [64]byte
	copy(sname[:], []byte(p.ServerName))
	sname[len(p.ServerName)] = 0
	b.WriteBytes(sname[:])

	var file [128]byte
	copy(file[:], []byte(p.BootFile))
	file[len(p.BootFile)] = 0
	b.WriteBytes(file[:])

	// The magic cookie.
	b.WriteBytes([]byte{99, 130, 83, 99})

	p.Options.Marshal(b)
	// TODO pad to 272 bytes for really old crap.
	return b.Data(), nil
}

func (p *Packet) UnmarshalBinary(q []byte) error {
	b := util.NewBuffer(q)
	if b.Len() < minPacketLen {
		return ErrInvalidPacket
	}

	p.Op = OpCode(b.Read8())
	p.HType = b.Read8()
	hlen := b.Read8()
	p.Hops = b.Read8()
	b.ReadBytes(p.TransactionID[:])
	p.Secs = b.Read16()

	flags := b.Read16()
	if flags&1<<15 != 0 {
		p.Broadcast = true
	}

	p.CIAddr = make(net.IP, net.IPv4len)
	b.ReadBytes(p.CIAddr)
	p.YIAddr = make(net.IP, net.IPv4len)
	b.ReadBytes(p.YIAddr)
	p.SIAddr = make(net.IP, net.IPv4len)
	b.ReadBytes(p.SIAddr)
	p.GIAddr = make(net.IP, net.IPv4len)
	b.ReadBytes(p.GIAddr)

	if hlen > MAClen {
		hlen = MAClen
	}
	p.CHAddr = make(net.HardwareAddr, hlen)
	b.ReadBytes(p.CHAddr)

	var sname [64]byte
	b.ReadBytes(sname[:])
	p.ServerName = string(sname[:strings.Index(string(sname[:]), "\x00")])

	var file [128]byte
	b.ReadBytes(file[:])
	p.BootFile = string(file[:strings.Index(string(file[:]), "\x00")])

	// Read the cookie and then fucking ignore it.
	var cookie [4]byte
	b.ReadBytes(cookie[:])

	return (&p.Options).Unmarshal(b)
}
