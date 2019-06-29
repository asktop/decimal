## decimal定制化

## github.com/ericlagergren/decimal

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



