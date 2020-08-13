package decimal

import (
    "fmt"
    "testing"
)

/*
警告：
1. 所有对象为指针，对象赋值时需要new一个新对象并copy其值。例如 b := new(Big).Copy(a) ;将a的值赋给b。
2. 无法定义全局默认小数位精度（默认所有位数为32位），当遇到除法运算后，可以使用 a.FloorN(n) ,将小数位数截取位n位。
*/

//创建|赋值
func TestNew(t *testing.T) {
    //直接new，值为0
    a := new(Big).FloorN(8)
    fmt.Println(a) //0

    //从int64创建
    b := New(123, 0)
    fmt.Println(b) //123

    //从int64创建float
    c := New(123, 1)
    fmt.Println(c) //12.3

    //从string创建
    d, _ := new(Big).SetString("123.46789")
    fmt.Println(d) //123.46789

    //赋值
    e := new(Big).Copy(d)
    fmt.Println(e) //123.46789
}

//运算
func TestYunSuan(t *testing.T) {
    a := New(456, 0)
    b := New(123, 0)

    //加
    c := new(Big).Add(a, b)
    fmt.Println(c) //579

    //减
    d := new(Big).Sub(a, b)
    fmt.Println(d) //333

    //乘
    e := new(Big).Mul(a, b)
    fmt.Println(e) //56088

    //除
    f := new(Big).Quo(a, b)
    fmt.Println(f) //3.7073170731707317073170731707318

    //向下取，n 为小数位数
    f1 := new(Big).Copy(f).FloorN(8)
    fmt.Println(f1) //3.70731707

    //向上取，n 为小数位数
    f2 := new(Big).Copy(f).CeilN(8)
    fmt.Println(f2) //3.70731708

    //四舍五入
    f3 := new(Big).Copy(f).Quantize(8)
    fmt.Println(f3) //3.70731707

    //取余
    g := new(Big).Rem(a, b)
    fmt.Println(g) //87

    //比较 1:大于；0：等于，-1：小于
    fmt.Println(a.Cmp(b)) //1
}
