package radix

import (
	"math/big"
)

// BigIntPow x^y
func BigIntPow(x *big.Int, y int) *big.Int {
	if y == 0 {
		return big.NewInt(1)
	} else if y < 0 {
		panic("y < 0")
	}
	sum := new(big.Int).Set(x)
	if y == 1 {
		return sum
	}
	for i, n := 0, y-1; i < n; i++ {
		sum.Mul(sum, x)
	}
	return sum
}

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

// TenStrToX 十进制转x进制
func (radix *Radix) TenStrToX(ten string) string {
	t, b := new(big.Int).SetString(ten, 10)
	if !b {
		panic(ten + " is invalid")
	}
	return radix.TenToX(t)
}

// TenToX 十进制转x进制
func (radix *Radix) TenToX(i *big.Int) string {
	t := big.NewInt(0)
	t.Set(i)
	bit := radix.GetBaseNumber()
	buf := ""
	rem := new(big.Int)
	for t.Sign() > 0 {
		rem.Mod(t, bit)
		buf = string(radix.GetRuneByInt(rem.Int64())) + buf
		t.Div(t, bit)
	}
	// 特别处理0
	if buf == "" {
		buf = string(radix.GetRuneByInt(0))
		// 修复问题：当输入为0时，不返回空字符串
		if buf == "" && bit.Int64() == 10 {
			buf = "0" // 返回'0'的标准字符表示
		}
	}
	return buf
}

// XStrToTenStr x进制转10进制
func (radix *Radix) XStrToTenStr(x string) string {
	return radix.XStrToTen(x).String()
}

// XStrToTen x进制转十进制
func (radix *Radix) XStrToTen(x string) *big.Int {
	return radix.XToTen([]rune(x))
}

// XToTen x进制转十进制
func (radix *Radix) XToTen(xx []rune) *big.Int {
	bit := radix.GetBaseNumber()

	num := new(big.Int)
	for i := len(xx) - 1; i >= 0; i-- {

		temp := BigIntPow(bit, len(xx)-1-i)
		num.Add(num, temp.Mul(temp, radix.GetIntByRune(xx[i])))
	}
	return num
}
