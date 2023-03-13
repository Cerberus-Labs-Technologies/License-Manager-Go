package license

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"strconv"
	"strings"
	"time"
)

// IntArray is an implementation of a integer array for the MySQL type
type IntArray []int

// Value implements the driver.Valuer interface,
// and turns the IntArray into a bitfield for MySQL storage.
func (b IntArray) Value() (driver.Value, error) {
	return []int{}, nil
}

// Scan implements the sql.Scanner interface,
// and returns an IntArray from mysql string splitted by comma
func (b *IntArray) Scan(src interface{}) error {
	if src == nil {
		*b = []int{}
		return nil
	}
	switch src.(type) {
	case []byte:
		stringArray := strings.Split(string(src.([]uint8)), ",")
		var t2 []int
		for _, i := range stringArray {
			j, err := strconv.Atoi(i)
			if err != nil {
				panic(err)
			}
			t2 = append(t2, j)
		}
		*b = t2
	case string:
		*b = IntArray{}
	default:
		return nil
	}

	return nil
}

// StringArray is an implementation of a string array for the MySQL type
type StringArray []string

// Value implements the driver.Valuer interface,
// and turns the StringArray into a bitfield (BIT(1)) for MySQL storage.
func (b StringArray) Value() (driver.Value, error) {
	return []string{}, nil
}

// Scan implements the sql.Scanner interface,
// and returns an StringArray from mysql string splitted by comma
func (b *StringArray) Scan(src interface{}) error {
	if src == nil {
		*b = []string{}
		return nil
	}
	switch src.(type) {
	case []byte:
		stringArray := strings.Split(string(src.([]uint8)), ",")
		// remove empty strings
		for i := 0; i < len(stringArray); i++ {
			if stringArray[i] == "" {
				stringArray = append(stringArray[:i], stringArray[i+1:]...)
			}
		}
		*b = (stringArray)
	case string:
		*b = StringArray{}
	default:
		return nil
	}

	return nil
}

type IntBool bool

// Value implements the driver.Valuer interface,
// and turns the BitBool into a bitfield (BIT(1)) for MySQL storage.
func (b IntBool) Value() (driver.Value, error) {
	if b {
		return strconv.Itoa(1), nil
	}
	return strconv.Itoa(0), nil
}

func (b *IntBool) ToBoolean() bool {
	return bool(*b)
}

// Scan implements the sql.Scanner interface,
// and turns the bitfield incoming from MySQL into a BitBool
func (b *IntBool) Scan(src interface{}) error {
	switch src.(type) {
	case int64:
		*b = IntBool(src.(int64) == 1)
	case int:
		*b = IntBool(src.(int) == 1)
	case int32:
		*b = IntBool(src.(int32) == 1)
	case int16:
		*b = IntBool(src.(int16) == 1)
	case int8:
		*b = IntBool(src.(int8) == 1)
	case uint64:
		*b = IntBool(src.(uint64) == 1)
	case uint:
		*b = IntBool(src.(uint) == 1)
	case uint32:
		*b = IntBool(src.(uint32) == 1)
	case uint16:
		*b = IntBool(src.(uint16) == 1)
	case uint8:
		*b = IntBool(src.(uint8) == 1)
	case float64:
		*b = IntBool(src.(float64) == 1)
	case float32:
		*b = IntBool(src.(float32) == 1)
	case string:
		*b = IntBool(src.(string) == "1")
	case []byte:
		*b = IntBool(string(src.([]uint8)) == "1")
	case nil:
		*b = false
	default:
		// print out the type of src
		return errors.New("incompatible type for IntBool")
	}
	return nil
}

type JsonFloat64 struct {
	sql.NullFloat64
}

func (f JsonFloat64) MarshalJSON() ([]byte, error) {
	if f.Valid {
		return []byte(strconv.FormatFloat(f.Float64, 'f', -1, 64)), nil
	}
	return []byte("null"), nil
}

func (f *JsonFloat64) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		f.Valid = false
		return nil
	}
	f.Valid = true
	var err error
	f.Float64, err = strconv.ParseFloat(string(data), 64)
	return err
}

type JsonString struct {
	sql.NullString
}

func (f JsonString) MarshalJSON() ([]byte, error) {
	if f.Valid {
		return []byte("\"" + f.String + "\""), nil
	}
	return []byte("null"), nil
}

func (f *JsonString) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		f.String = ""
		f.Valid = false
		return nil
	}
	f.Valid = true
	f.String = strings.Trim(string(data), "\"")
	return nil
}

type JsonInt32 struct {
	sql.NullInt32
}

func (f JsonInt32) MarshalJSON() ([]byte, error) {
	if f.Valid {
		return []byte(strconv.FormatInt(int64(f.Int32), 10)), nil
	}
	return []byte("null"), nil
}

func (f *JsonInt32) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		f.Valid = false
		return nil
	}
	f.Valid = true
	var err error

	i, err := strconv.ParseInt(string(data), 10, 8)

	if err != nil {
		return err
	}

	f.Int32 = int32(i)
	return err
}

type TimeStamp struct {
	time.Time
}

func (t TimeStamp) MarshalJSON() ([]byte, error) {
	return []byte("\"" + t.Format("2006-01-02 15:04:05") + "\""), nil
}

func (t *TimeStamp) UnmarshalJSON(data []byte) error {
	a, err := time.Parse("2006-01-02 15:04:05", strings.Trim(string(data), "\""))
	t.Time = a
	return err
}

func (dateT *TimeStamp) Scan(src interface{}) error {
	if t, ok := src.(time.Time); ok {
		dateT.Time = t
	}
	return nil
}

func (dateT TimeStamp) Value() (driver.Value, error) {
	return dateT.Time, nil
}
