package units

type Kind int

var (
	UnknownKind Kind = 0
	TimeSpan    Kind = 1
	Memory      Kind = 2
	Number      Kind = 3
	TimeStamp   Kind = 4
	Frequency   Kind = 5
	Percentage  Kind = 6
)

type Numeric interface {
	Float() bool
	Int64() int64
	Float64() float64
	Multi(n Numeric) Numeric
}

type I64 int64

func (i I64) Float() bool {
	return false
}

func (i I64) Int64() int64 {
	return int64(i)
}

func (i I64) Float64() float64 {
	return float64(i)
}

func (i I64) Multi(n Numeric) Numeric {
	if n.Float() {
		return F64(n.Float64() * i.Float64())
	}

	return I64(i.Int64() * n.Int64())
}

type F64 float64

func (f F64) Float() bool {
	return true
}

func (f F64) Int64() int64 {
	return int64(f)
}

func (f F64) Float64() float64 {
	return float64(f)
}

func (f F64) Multi(n Numeric) Numeric {
	return F64(float64(f) * n.Float64())
}

var _ Numeric = F64(0)
var _ Numeric = I64(0)

type Unit struct {
	Name string
	Kind Kind
	Base Numeric
}

func newUnit(Name string, kind Kind, Base Numeric) *Unit {
	return &Unit{
		Kind: kind,
		Name: Name,
		Base: Base,
	}
}

func (u *Unit) Derived(Name string, times Numeric) *Unit {
	return &Unit{
		Kind: u.Kind,
		Name: Name,
		Base: u.Base.Multi(times),
	}
}

var (
	Unknown = newUnit("unknown", UnknownKind, I64(0))

	Nanosecond  = newUnit("ns", TimeSpan, I64(1))
	Microsecond = Nanosecond.Derived("μs", I64(1000))
	Millisecond = Microsecond.Derived("ms", I64(1000))
	Second      = Millisecond.Derived("s", I64(1000))
	Minute      = Second.Derived("min", I64(60))
	Hour        = Minute.Derived("h", I64(60))
	Day         = Hour.Derived("d", I64(24))
	Week        = Day.Derived("w", I64(7))

	Byte     = newUnit("B", Memory, I64(1))
	Kilobyte = Byte.Derived("KB", I64(1024))
	Megabyte = Kilobyte.Derived("MB", I64(1024))
	Gigabyte = Megabyte.Derived("GB", I64(1024))
	Terabyte = Gigabyte.Derived("TB", I64(1024))
	Petabyte = Terabyte.Derived("PB", I64(1024))

	UnixNano   = newUnit("epoch_ns", TimeStamp, I64(1))
	UnixMicro  = UnixNano.Derived("epoch_μs", I64(1000))
	UnixMilli  = UnixMicro.Derived("epoch_ms", I64(1000))
	UnixSecond = UnixMilli.Derived("epoch_s", I64(1000))

	Multiple = newUnit("", Percentage, I64(1))    // eg: 0.15
	Percent  = newUnit("%", Percentage, I64(100)) // eg: 15%

	Hertz = newUnit("hz", Frequency, I64(1))
)
