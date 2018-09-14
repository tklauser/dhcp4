package dhcp4server

import (
	"net"
	"testing"
)

func TestNextIP(t *testing.T) {
	for _, tt := range []struct {
		ip   net.IP
		want net.IP
	}{
		{
			ip:   net.IP{192, 168, 0, 0},
			want: net.IP{192, 168, 0, 1},
		},
		{
			ip:   net.IP{192, 168, 0, 255},
			want: net.IP{192, 168, 1, 0},
		},
		{
			ip:   net.IP{192, 168, 1, 255},
			want: net.IP{192, 168, 2, 0},
		},
		{
			ip:   net.IP{192, 255, 255, 255},
			want: net.IP{193, 0, 0, 0},
		},
		{
			ip:   net.IP{255, 255, 255, 255},
			want: net.IP{0, 0, 0, 0},
		},
	} {
		nextIP(tt.ip)
		if !tt.ip.Equal(tt.want) {
			t.Errorf("got %q, want %q", tt.ip, tt.want)
		}
	}
}

func TestIPAlloc(t *testing.T) {
	_, subnet, err := net.ParseCIDR("192.168.1.0/24")
	if err != nil {
		t.Fatal(err)
	}

	ipa := newIPAllocator(subnet)

	want := net.IP{192, 168, 1, 0}
	if got := ipa.alloc(); !got.Equal(want) {
		t.Fatalf("first alloc() = %v, want %v", got, want)
	}

	want = net.IP{192, 168, 1, 1}
	if got := ipa.alloc(); !got.Equal(want) {
		t.Fatalf("second alloc() = %v, want %v", got, want)
	}

	ipa.free(net.IP{192, 168, 1, 0})
	want = net.IP{192, 168, 1, 0}
	if got := ipa.alloc(); !got.Equal(want) {
		t.Fatalf("third alloc() = %v, want %v", got, want)
	}

	want = net.IP{192, 168, 1, 2}
	if got := ipa.alloc(); !got.Equal(want) {
		t.Fatalf("fourth alloc() = %v, want %v", got, want)
	}

	ipa.free(net.IP{192, 168, 1, 3})

	want = net.IP{192, 168, 1, 3}
	if got := ipa.alloc(); !got.Equal(want) {
		t.Fatalf("fifth alloc() = %v, want %v", got, want)
	}
}
