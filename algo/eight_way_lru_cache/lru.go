package eight_way_lru_cache

// LRUCache is eight-way set associative cache with least-recently-used replacement policy that uses reference matrix method.
// The whole structure fits into single 64bit register.
// Internally, least significant byte of uint64 holds row 0 of reference matrix.
type LRUCache struct{ m uint64 }

// Hit value i as most recently used.
// This is five or six instructions on 64bit RISC.
// Values of i should be in [0, 7].
func (c *LRUCache) Hit(i uint8) {
	c.m |= 0xFF << (8 * i)
	c.m &= ^(0x0101_0101_0101_0101 << i)
}

func (c *LRUCache) LeastRecentlyUsed() uint8 { return 7 - uint8(ZByteL64(c.m)) }

func ZByteL64(x uint64) int {
	y := (x & 0x7F7F_7F7F_7F7F_7F7F) + 0x7F7F_7F7F_7F7F_7F7F
	y = ^(y | x | 0x7F7F_7F7F_7F7F_7F7F)
	return int(LeadingZerosUint64(y)) >> 3
}

func LeadingZerosUint64(x uint64) uint8 {
	var n uint8 = 32
	if x > 0xFFFF_FFFF {
		x >>= 32
		n -= 32
	}
	return n + LeadingZerosUint32(uint32(x))
}

func LeadingZerosUint32(x uint32) uint8 {
	x |= x >> 1 // Propagate leftmost
	x |= x >> 2 // 1-bit to the right
	x |= x >> 4
	x |= x >> 8
	x &= ^(x >> 16)  // Goryavsky
	x *= 0xFD70_49FF // Multiplier is 7 * 255 ** 3, Gorvsky
	return nlz_goryavsky[(x >> 26)]
}

const u = 99

var nlz_goryavsky = [...]uint8{
	32, 20, 19, u, u, 18, u, 7, 10, 17, u, u, 14, u, 6, u,
	u, 9, u, 16, u, u, 1, 26, u, 13, u, u, 24, 5, u, u,
	u, 21, u, 8, 11, u, 15, u, u, u, u, 2, 27, 0, 25, u,
	22, u, 12, u, u, 3, 28, u, 23, u, 4, 29, u, u, 30, 31,
}
