package util

import (
	"math/rand"
	"strings"
	"time"
)

func init(){
	rand.Seed(time.Now().UnixMicro());
}
var char string="abcdefghijklmnopqrstuvwxyz"
func RandomString(n int) string{
	sb:=strings.Builder{};
	for i:=0;i<n;i++{
		sb.WriteByte(char[rand.Intn(len(char))]);
	}
	return sb.String();
}
func RandomNum(min,max int64) int64{
	return min+rand.Int63n(max-min+1);
}
func RandomBalance()int64{
	return RandomNum(100000,1000000)
}

func RandomCurrency() string{
	currencies:=[]string{"USD","VND","GBP"}
	return currencies[rand.Intn(len(currencies))];
}

func RandomName() string{
	return RandomString(10);
}
func RandomEmail()string{
	return RandomString(10)+"@gmail.com"
}