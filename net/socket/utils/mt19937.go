//梅森旋转算法（Mersenne twister）是一个伪随机数发生算法。由松本真和西村拓士[1]在1997年开发，
//基于有限二进制字段上的矩阵线性递归可以快速产生高质量的伪随机数，修正了古典随机数发生算法的很多缺陷。
//https://zh.wikipedia.org/wiki/梅森旋转算法
//用于与不同客户端统一随机数生成规律
//这里用于生成随机密钥
package utils

import "math"

type MT19937 struct {
	mt    [624]uint64
	index int32
}

func NewMT19937(seed uint64) *MT19937 {
	return new(MT19937).init(seed)
}

func (m *MT19937) init(seed uint64) *MT19937 {
	m.mt[0] = seed
	for i := 1; i < len(m.mt); i++ {
		m.mt[i] = (m.mt[i-1]^(m.mt[i-1]>>30))*0x6c078965 + uint64(i)
	}

	return m
}

func (m *MT19937) twist() {
	for i := 0; i < len(m.mt); i++ {
		y := m.mt[i]&0x80000000 + (m.mt[(i+1)%624] & 0x7fffffff)
		m.mt[i] = m.mt[(i+397)%624] ^ y>>1

		if y%2 != 0 {
			m.mt[i] = m.mt[i] ^ 0x9908b0df
		}
	}
}

func (m *MT19937) extractNumber() uint64 {
	if m.index == 0 {
		m.twist()
	}

	y := m.mt[m.index]
	y = y ^ (y >> 11)
	y = y ^ (y << 7 & 0x9d2c5680)
	y = y ^ (y << 15 & 0xefc60000)
	y = y ^ (y >> 18)

	m.index = (m.index + 1) % 624
	return y
}

//-----------

func (m *MT19937) Uint8() uint8 {
	return uint8(m.extractNumber() % math.MaxUint8)
}

func (m *MT19937) Int31n(n int32) int32 {
	if n <= 0 {
		return 0
	}

	return int32(m.extractNumber()%math.MaxInt32) % n
}

func (m *MT19937) Rand32BytesKey() (key []byte) {
	key = make([]byte, 32)

	for i := 0; i < len(key); i++ {
		key[i] = byte(m.Uint8())
	}

	return
}
