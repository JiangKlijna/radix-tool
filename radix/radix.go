package radix

import (
	"math/big"
)

// Radixer 代表进制要素
type Radixer interface {
	GetBaseNumber() *big.Int
	GetRuneByInt(int64) rune
	GetIntByRune(rune) *big.Int
}

// Radix 进制转换
type Radix struct {
	Radixer
}

// TenStrToX 10进制转x进制
func (radix *Radix) TenStrToX(ten string) string {
	t, b := new(big.Int).SetString(ten, 10)
	if !b {
		panic("invalid decimal string: " + ten)
	}
	return radix.TenToX(t)
}

// TenToX 10进制转x进制
func (radix *Radix) TenToX(i *big.Int) string {
	if i.Sign() == 0 {
        return string(radix.GetRuneByInt(0))
    }
	t := new(big.Int).Set(i)

	runes := make([]rune, 0, 64)
	bit := radix.GetBaseNumber()

	rem := new(big.Int)
    for t.Sign() > 0 {
        // 同时获得商和余数
        t.QuoRem(t, bit, rem)
        runes = append(runes, radix.GetRuneByInt(rem.Int64()))
    }
	for j, k := 0, len(runes)-1; j < k; j, k = j+1, k-1 {
        runes[j], runes[k] = runes[k], runes[j]
    }
	return string(runes)
}

// XStrToTenStr x进制转10进制
func (radix *Radix) XToTenStr(x string) string {
	return radix.XToTen(x).String()
}

// XToTen x进制转10进制
func (radix *Radix) XToTen(x string) *big.Int {
	runes := []rune(x)
	length := int64(len(runes))
	bit := radix.GetBaseNumber()

	num := new(big.Int)
	for i := length - 1; i >= 0; i-- {
		temp := new(big.Int).Exp(bit, big.NewInt(length - 1 - i), nil)
		num.Add(num, temp.Mul(temp, radix.GetIntByRune(runes[i])))
	}
	return num
}
