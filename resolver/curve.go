package resolver

import (
	"crypto/rand"
	"encoding/binary"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"golang.org/x/crypto/curve25519"
	"golang.org/x/crypto/nacl/secretbox"
)

var public, private = &[32]byte{}, &[32]byte{}

var resolverIP = "208.67.220.220:443"
var serverPubKey = decodeKey("B735:1140:206F:225D:3E2B:D822:D7FD:691E:A1C3:3CC8:D666:8D0C:BE04:BFAB:CA43:FB79")

func init() {
	rand.Read(private[:])
	curve25519.ScalarBaseMult(public, private)
	fmt.Println(serverPubKey)
	fmt.Println(public)
}

var incr uint64 = 1

func genNonce(b []byte) (n uint64, err error) {
	if len(b) < 12 {
		return 0, errors.New("buffer too small, must greater than 12")
	}
	binary.LittleEndian.PutUint64(b, incr)
	incr++
	rand.Read(b[8:12])
	n = 12
	return
}

func decodeKey(key string) *[32]byte {
	r := new([32]byte)
	for k, v := range strings.Split(key, ":") {
		d, _ := strconv.ParseUint(v, 16, 16)
		r[k*2] = byte(d >> 8)
		r[k*2+1] = byte(d)
	}
	return r
}

func encodePacket(dst, src []byte) {
	copy(dst, public[:])
	dst = dst[32:]

	var nonce [24]byte
	genNonce(nonce[:])
	copy(dst, nonce[:])
	dst = dst[24:]

	secretbox.Open(dst)
	//	poly1305.Sum()
}
