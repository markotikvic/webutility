package numbersutility

import (
	"fmt"
	"math"
	"math/big"
	"strings"
)

func RoundFloat64(f float64, dec int) float64 {
	p := math.Pow(10, float64(dec))
	return math.Round(f*p) / p
}

func NewBigFloat(f float64, prec uint) *big.Float {
	x := big.NewFloat(f)
	x.SetPrec(prec)
	return x
}

func AddBigFloat(x, y *big.Float) *big.Float {
	z := big.NewFloat(0.0)
	z.SetPrec(x.Prec())
	z.Add(x, y)
	return z
}

func SubBigFloat(x, y *big.Float) *big.Float {
	z := big.NewFloat(0.0)
	z.SetPrec(x.Prec())

	yneg := big.NewFloat(0.0)
	yneg.Neg(y)

	z.Add(x, yneg)
	return z
}

func MulBigFloat(x, y *big.Float) *big.Float {
	z := big.NewFloat(0.0)
	z.SetPrec(x.Prec())
	z.Mul(x, y)
	return z
}

func DivBigFloat(x, y *big.Float) *big.Float {
	z := big.NewFloat(0.0)
	z.SetPrec(x.Prec())
	z.Quo(x, y)
	return z
}

func BigFloattoFloat(f *big.Float) float64 {
	v, _ := f.Float64()
	return v
}

func Float64ToString(f float64) string {
	return fmt.Sprintf("%.2f", f)
}

func Float64PtrToString(f *float64) string {
	if f == nil {
		return ""
	}
	return fmt.Sprintf("%.2f", *f)
}

func FormatFloat64(f float64, dec int) string {
	res := ""
	f = RoundFloat64(f, dec)
	i := int64(f)
	format := fmt.Sprintf("%%.%df", dec)
	parts := strings.Split(fmt.Sprintf(format, f-float64(i)), ".")
	decimals := parts[1]
	res = FormatInt64Number(i) + "," + decimals

	return res
}

func MaxFloat64(vars ...float64) (max float64) {
	max = vars[0]
	for _, v := range vars {
		if v > max {
			max = v
		}
	}
	return max
}
