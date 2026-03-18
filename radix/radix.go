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
		temp := new(big.Int).Exp(bit, big.NewInt(length-1-i), nil)
		num.Add(num, temp.Mul(temp, radix.GetIntByRune(runes[i])))
	}
	return num
}

// TenToXParallel 并行10进制转X进制

// TenToXParallel 并行10进制转X进制
func (radix *Radix) TenToXParallel(i *big.Int) string {
	// 1. 基础情况：0值处理
	if i.Sign() == 0 {
		return string(radix.GetRuneByInt(0))
	}

	// 2. 阈值判断：提高阈值以抵消 Goroutine 调度开销
	// big.Int 运算在 1024 bit 以下非常快，并行反而变慢
	if i.BitLen() < 1024 {
		return radix.TenToX(i)
	}

	base := radix.GetBaseNumber()

	// 3. 优化的分割策略：目标是将数字近似对半切分
	// 计算 n 使得 base^n 大致等于 sqrt(i)
	// 公式推导：log_base(i) ~ bitLen(i)/bitLen(base) -> n = log_base(i) / 2
	n := i.BitLen() / 2 / base.BitLen()
	if n < 1 {
		n = 1
	}

	// 计算分割基数 split = base^n
	split := new(big.Int).Exp(base, big.NewInt(int64(n)), nil)

	// 执行除法获取高位和低位
	high := new(big.Int)
	low := new(big.Int)
	high.QuoRem(i, split, low)

	// 4. 并发模型优化：主协程处理低位，新协程处理高位
	ch := make(chan string, 1)

	go func() {
		// 高位任务交给新协程
		ch <- radix.TenToXParallel(high)
	}()

	// 当前协程直接计算低位，利用等待时间
	lowStr := radix.TenToXParallel(low)

	// 接收高位结果
	highStr := <-ch

	// 5. 合并逻辑：处理前导零对齐
	if high.Sign() > 0 {
		lowRunes := []rune(lowStr)
		padding := n - len(lowRunes)

		if padding > 0 {
			// 构造前导零切片
			zeroRune := radix.GetRuneByInt(0)
			prefix := make([]rune, padding)
			for i := range prefix {
				prefix[i] = zeroRune
			}
			// 拼接：前导零 + 低位字符串
			lowStr = string(prefix) + lowStr
		}
	}

	// 如果高位为0，直接返回低位（避免无意义的字符串拼接）
	if high.Sign() == 0 {
		return lowStr
	}

	return highStr + lowStr
}

// XToTenParallel 并行X进制转10进制
func (radix *Radix) XToTenParallel(x string) *big.Int {
	runes := []rune(x)
	// 定义内部递归函数，避免反复类型转换和内存分配
	var parallelRec func(rs []rune) *big.Int
	parallelRec = func(rs []rune) *big.Int {
		n := len(rs)
		if n == 0 {
			return big.NewInt(0)
		}

		// 阈值可以适当调大，big.Int 运算本身有优化，太细的切分得不偿失
		if n < 256 {
			return radix.XToTen(string(rs))
		}

		mid := n / 2

		// 优化：当前线程处理右半部分，新线程处理左半部分
		// 这样避免了当前线程空等，也减少了一半的 Goroutine 创建
		ch := make(chan *big.Int, 1)
		go func() {
			ch <- parallelRec(rs[:mid])
		}()

		lowNum := parallelRec(rs[mid:])
		highNum := <-ch

		// 合并逻辑
		base := radix.GetBaseNumber()

		// 修复：使用 rune 的长度计算指数，而不是 string 的字节数
		lowLen := int64(len(rs[mid:]))

		// Exp 计算
		multiplier := new(big.Int).Exp(base, big.NewInt(lowLen), nil)

		res := new(big.Int).Mul(highNum, multiplier)
		res.Add(res, lowNum)
		return res
	}

	return parallelRec(runes)
}
