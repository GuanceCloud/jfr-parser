package units

type IQuantity interface {
	Unit() *Unit
	In(unit *Unit) (IQuantity, error)
	FloatValue() float64
	IntValue() int64
	Add(q IQuantity) (IQuantity, error)
}

type IntQuantity struct {
	num  int64
	unit *Unit
}

func (i *IntQuantity) Add(q IQuantity) (IQuantity, error) {
	//TODO implement me
	panic("implement me")
}

func (i *IntQuantity) IntValue() int64 {
	return i.num
}

func (i *IntQuantity) FloatValue() float64 {
	return float64(i.num)
}

func (i *IntQuantity) Unit() *Unit {
	return i.unit
}

func (i *IntQuantity) In(unit *Unit) (IQuantity, error) {
	// todo
	return nil, nil
}

func NewIntQuantity(number int64, unit *Unit) IQuantity {
	return &IntQuantity{
		num:  number,
		unit: unit,
	}
}

type FloatQuantity struct {
	num  float64
	unit *Unit
}

func (f *FloatQuantity) Add(q IQuantity) (IQuantity, error) {
	//TODO implement me
	panic("implement me")
}

func (f *FloatQuantity) Unit() *Unit {
	return f.unit
}

func (f *FloatQuantity) In(unit *Unit) (IQuantity, error) {
	//TODO implement me
	return nil, nil
}

func (f *FloatQuantity) FloatValue() float64 {
	return f.num
}

func (f *FloatQuantity) IntValue() int64 {
	return int64(f.num)
}

func NewFloatQuantity(number float64, unit *Unit) IQuantity {
	return &FloatQuantity{
		num:  number,
		unit: unit,
	}
}
