package main

import (
	"fmt"
	"math"
	"math/big"
	"os"
	"strconv"
)

func main() {
	n := 160
	n0 := 40
	n, _ = strconv.Atoi(os.Args[1])

	M := 2 * int(math.Ceil(float64(n)/math.Pow(math.Log(float64(n)), float64(3))))
	N := int(math.Ceil(math.Log(float64(10)) / math.Log(math.E*float64(2*M)) * float64(n+n0+1)))
	fmt.Printf("N: %d\n", N)

	ch_1 := make(chan *big.Float, 1)
	ch_2 := make(chan *big.Float, 1)

	ch_b := make(chan *big.Float, (M+1)*N)
	for k := 0; k < (M+1)*N; k++ {
		go func(k int, ch chan<- *big.Float) {
			p := big.NewInt(int64(2*k + 1))
			x := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(n)), p)
			x.Mul(x, big.NewInt(4))
			x.Mod(x, p)
			tmp1 := new(big.Float).SetInt(x).SetPrec(uint(n0 * 2))
			tmp2 := new(big.Float).SetInt(p).SetPrec(uint(n0 * 2))
			tmp1.Quo(tmp1, tmp2)
			if k&1 == 1 {
				tmp1.Neg(tmp1)
			}
			ch <- tmp1
		}(k, ch_b)
	}

	go func() {
		b := new(big.Float).SetInt64(0).SetPrec(uint(n0 * 2))
		for k := 0; k < (M+1)*N; k++ {
			tmp := <-ch_b
			b.Add(b, tmp)
		}
		ch_1 <- b
	}()

	ch_c := make(chan *big.Float, N)
	for k := 0; k < N; k++ {
		go func(k int, ch chan<- *big.Float) {
			p := big.NewInt(int64(2*M*N + 2*k + 1))
			x := big.NewInt(1)
			tmpC := big.NewInt(1)
			for j := 1; j <= k; j++ {
				tmpC.Mul(tmpC, big.NewInt(int64(N-j+1)))
				tmpC.Div(tmpC, big.NewInt(int64(j)))
				x.Add(x, tmpC)
				x.Mod(x, p)
			}
			y := x
			y1 := new(big.Int).Exp(big.NewInt(5), big.NewInt(int64(N-2)), p)
			y2 := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(n-N+2)), p)
			y.Mul(y, y1)
			y.Mod(y, p)
			y.Mul(y, y2)
			y.Mod(y, p)
			tmp1 := new(big.Float).SetInt(y).SetPrec(uint(n0 * 2))
			tmp2 := new(big.Float).SetInt(p).SetPrec(uint(n0 * 2))
			tmp1.Quo(tmp1, tmp2)
			if k&1 == 1 {
				tmp1.Neg(tmp1)
			}
			ch <- tmp1
		}(k, ch_c)
	}

	go func() {
		c := new(big.Float).SetInt64(0).SetPrec(uint(n0 * 2))
		for k := 0; k < N; k++ {
			tmp := <-ch_c
			c.Add(c, tmp)
		}
		ch_2 <- c
	}()

	b := <-ch_1
	c := <-ch_2
	result := new(big.Float).Sub(b, c)
	tmp := new(big.Float).Copy(result)
	tmp1, _ := tmp.Int64()
	tmp = tmp.SetInt64(tmp1)
	if result.Sign() < 0 {
		tmp = tmp.Sub(tmp, new(big.Float).SetInt64(1))
	}
	result = result.Sub(result, tmp)
	fmt.Println(result.Text('f', n0))
}
