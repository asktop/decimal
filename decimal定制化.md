<h1> github.com/ericlagergren/decimal定制化 </h1>


[TOC]



### 1. 修改为decimal转string时不使用科学计数法

**github.com/ericlagergren/decimal/big.go**

```go
func (x *Big) String() string {
	...
    //修改 normal 为 plain
	f.format(x, plain, e)
	...
}

func (x *Big) Format(s fmt.State, c rune) {
    ...
    //修改 normal 为 plain
	case 'v':
		// %v == %s
		if !hash && !plus {
			f.format(x, plain, e)
			break
		}
	...
}

func (x *Big) MarshalText() ([]byte, error) {
    ...
    //修改 normal 为 plain
	f.format(x, plain, e)
	...
}
```



### 2. 修改decimal可以scan读取sql的数据

**github.com/ericlagergren/decimal/big.go**

```go
//添加导包
import "database/sql/driver"

//修改:
func (z *Big) Scan(state fmt.ScanState, verb rune) error {
	return z.scan(byteReader{state})
}

var _ fmt.Scanner = (*Big)(nil)

//为：
func (x *Big) Value()(driver.Value,error){
	return x.String(),nil
}

// Scan implements the sql.Scanner interface for database deserialization.
func (d *Big) Scan(value interface{}) error {
	// first try to see if the data is stored in database as a Numeric datatype
	switch v := value.(type) {

	case float32:
		d.SetFloat64(float64(v))
		return nil

	case float64:
		// numeric in sqlite3 sends us float64
		d.SetFloat64(v)
		return nil

	case int64:
		// at least in sqlite3 when the value is 0 in db, the data is sent
		// to us as an int64 instead of a float64 ...
		//d.SetString()
		// New(v, 0)
		d.SetMantScale(v,0)
		return nil

	default:
		// default is trying to interpret value stored as string
		str, err := unquoteIfQuoted(v)
		if err != nil {
			return err
		}
		err = d.scan(strings.NewReader(str))
		return err
	}
}

func unquoteIfQuoted(value interface{}) (string, error) {
	var bytes []byte

	switch v := value.(type) {
	case string:
		bytes = []byte(v)
	case []byte:
		bytes = v
	default:
		return "", fmt.Errorf("Could not convert value '%+v' to byte array of type '%T'",
			value, value)
	}

	// If the amount is quoted, strip the quotes
	if len(bytes) > 2 && bytes[0] == '"' && bytes[len(bytes)-1] == '"' {
		bytes = bytes[1 : len(bytes)-1]
	}
	return string(bytes), nil
}
```



### 3. 修改默认精度 

**github.com/ericlagergren/decimal/context.go**

```go
//修改:16为32
DefaultPrecision   = 32               // default precision for literals.
```



### 4. 添加向上下取小数方法 

**github.com/ericlagergren/decimal/big.go**

```go
//添加:
//向下取，n 为小数位数
 func (z *Big) FloorN(n int) *Big {
 	scale := New(1, n)
 	z.Quo(z, scale)
 	floor(z, z)
 	z.Mul(z, scale)
 	//位数补全
 	if n > 0 && z.Scale() >= 0 && z.Scale() < n {
 		y := z.String()
 		if z.Scale() == 0 {
 			y = y + "." + strings.Repeat("0", n)
 		} else {
 			y = y + strings.Repeat("0", n-z.Scale())
 		}
 		z.SetString(y)
 	}
 	return z
 }
 func floor(z, x *Big) *Big {
 	if z.CheckNaNs(x, nil) {
 		return z
 	}
 	ctx := z.Context
 	ctx.RoundingMode = ToNegativeInf
 	return ctx.RoundToInt(z.Copy(x))
 }

  //向上取，n 为小数位数
 func (z *Big) CeilN(n int) *Big {
 	scale := New(1, n)
 	z.Quo(z, scale)
 	ceil(z, z)
 	z.Mul(z, scale)
 	//位数补全
 	if n > 0 && z.Scale() >= 0 && z.Scale() < n {
 		y := z.String()
 		if z.Scale() == 0 {
 			y = y + "." + strings.Repeat("0", n)
 		} else {
 			y = y + strings.Repeat("0", n-z.Scale())
 		}
 		z.SetString(y)
 	}
 	return z
 }
 func ceil(z, x *Big) *Big {
 	// ceil(x) = -floor(-x)
 	return z.Neg(floor(z, copyNeg(z, x)))
 }
 func copyNeg(z, x *Big) *Big {
 	if x.Signbit() {
 		return z.CopySign(x, New(+1, 0))
 	}
 	return z.CopySign(x, New(-1, 0))
 }
```


