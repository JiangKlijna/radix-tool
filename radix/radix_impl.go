package radix

import (
	"math/big"
)

// characterRadix 字符表进制
type characterRadix struct {
	bit *big.Int
}

func (ch *characterRadix) GetBaseNumber() *big.Int {
	return ch.bit
}
func (ch *characterRadix) GetRuneByInt(i int64) rune {
	return rune(i)
}
func (ch *characterRadix) GetIntByRune(r rune) *big.Int {
	return big.NewInt(int64(r))
}
func NewCharacterRadix(bit int) *Radix {
	// bit [2, 0xD800|55296, 0xDFFF|57343]
	return &Radix{&characterRadix{big.NewInt(int64(bit))}}
}

// stringRadix 自定义进制要素
type stringRadix struct {
	bit   *big.Int
	base  []rune
	table map[rune]*big.Int
}

func (sh *stringRadix) GetBaseNumber() *big.Int {
	return sh.bit
}
func (sh *stringRadix) GetRuneByInt(i int64) rune {
	return sh.base[i]
}
func (sh *stringRadix) GetIntByRune(r rune) *big.Int {
	return sh.table[r]
}

func NewRadixByBit(bit int) *Radix {
	const BASE = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	return NewRadixByString(BASE[0:bit])
}

func NewRadixByString(base string) *Radix {
	return NewRadixByRunes([]rune(base))
}

func NewRadixByRunes(base []rune) *Radix {
	radix := &stringRadix{big.NewInt(int64(len(base))), base, make(map[rune]*big.Int)}
	for i, v := range base {
		radix.table[v] = big.NewInt(int64(i))
	}
	return &Radix{radix}
}
