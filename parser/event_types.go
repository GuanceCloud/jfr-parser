package parser

import (
	"fmt"
	"reflect"
)

var events = map[string]func() Event{
	"jdk.ActiveRecording":                      func() Event { return new(ActiveRecording) },
	"jdk.ActiveSetting":                        func() Event { return new(ActiveSetting) },
	"jdk.BooleanFlag":                          func() Event { return new(BooleanFlag) },
	"jdk.CPUInformation":                       func() Event { return new(CPUInformation) },
	"jdk.CPULoad":                              func() Event { return new(CPULoad) },
	"jdk.CPUTimeStampCounter":                  func() Event { return new(CPUTimeStampCounter) },
	"jdk.ClassLoaderStatistics":                func() Event { return new(ClassLoaderStatistics) },
	"jdk.ClassLoadingStatistics":               func() Event { return new(ClassLoadingStatistics) },
	"jdk.CodeCacheConfiguration":               func() Event { return new(CodeCacheConfiguration) },
	"jdk.CodeCacheStatistics":                  func() Event { return new(CodeCacheStatistics) },
	"jdk.CodeSweeperConfiguration":             func() Event { return new(CodeSweeperConfiguration) },
	"jdk.CodeSweeperStatistics":                func() Event { return new(CodeSweeperStatistics) },
	"jdk.CompilerConfiguration":                func() Event { return new(CompilerConfiguration) },
	"jdk.CompilerStatistics":                   func() Event { return new(CompilerStatistics) },
	"jdk.DoubleFlag":                           func() Event { return new(DoubleFlag) },
	"jdk.ExceptionStatistics":                  func() Event { return new(ExceptionStatistics) },
	"jdk.ExecutionSample":                      func() Event { return new(ExecutionSample) },
	"jdk.GCConfiguration":                      func() Event { return new(GCConfiguration) },
	"jdk.GCHeapConfiguration":                  func() Event { return new(GCHeapConfiguration) },
	"jdk.GCSurvivorConfiguration":              func() Event { return new(GCSurvivorConfiguration) },
	"jdk.GCTLABConfiguration":                  func() Event { return new(GCTLABConfiguration) },
	"jdk.InitialEnvironmentVariable":           func() Event { return new(InitialEnvironmentVariable) },
	"jdk.InitialSystemProperty":                func() Event { return new(InitialSystemProperty) },
	"jdk.IntFlag":                              func() Event { return new(IntFlag) },
	"jdk.JavaMonitorEnter":                     func() Event { return new(JavaMonitorEnter) },
	"jdk.JavaMonitorWait":                      func() Event { return new(JavaMonitorWait) },
	"jdk.JavaThreadStatistics":                 func() Event { return new(JavaThreadStatistics) },
	"jdk.JVMInformation":                       func() Event { return new(JVMInformation) },
	"jdk.LoaderConstraintsTableStatistics":     func() Event { return new(LoaderConstraintsTableStatistics) },
	"jdk.LongFlag":                             func() Event { return new(LongFlag) },
	"jdk.ModuleExport":                         func() Event { return new(ModuleExport) },
	"jdk.ModuleRequire":                        func() Event { return new(ModuleRequire) },
	"jdk.NativeLibrary":                        func() Event { return new(NativeLibrary) },
	"jdk.NetworkUtilization":                   func() Event { return new(NetworkUtilization) },
	"jdk.ObjectAllocationInNewTLAB":            func() Event { return new(ObjectAllocationInNewTLAB) },
	"jdk.ObjectAllocationOutsideTLAB":          func() Event { return new(ObjectAllocationOutsideTLAB) },
	"jdk.OSInformation":                        func() Event { return new(OSInformation) },
	"jdk.PhysicalMemory":                       func() Event { return new(PhysicalMemory) },
	"jdk.PlaceholderTableStatistics":           func() Event { return new(PlaceholderTableStatistics) },
	"jdk.ProtectionDomainCacheTableStatistics": func() Event { return new(ProtectionDomainCacheTableStatistics) },
	"jdk.StringFlag":                           func() Event { return new(StringFlag) },
	"jdk.StringTableStatistics":                func() Event { return new(StringTableStatistics) },
	"jdk.SymbolTableStatistics":                func() Event { return new(SymbolTableStatistics) },
	"jdk.SystemProcess":                        func() Event { return new(SystemProcess) },
	"jdk.ThreadAllocationStatistics":           func() Event { return new(ThreadAllocationStatistics) },
	"jdk.ThreadCPULoad":                        func() Event { return new(ThreadCPULoad) },
	"jdk.ThreadContextSwitchRate":              func() Event { return new(ThreadContextSwitchRate) },
	"jdk.ThreadDump":                           func() Event { return new(ThreadDump) },
	"jdk.ThreadPark":                           func() Event { return new(ThreadPark) },
	"jdk.ThreadStart":                          func() Event { return new(ThreadStart) },
	"jdk.UnsignedIntFlag":                      func() Event { return new(UnsignedIntFlag) },
	"jdk.UnsignedLongFlag":                     func() Event { return new(UnsignedLongFlag) },
	"jdk.VirtualizationInformation":            func() Event { return new(VirtualizationInformation) },
	"jdk.YoungGenerationConfiguration":         func() Event { return new(YoungGenerationConfiguration) },
}

func indirect(rv reflect.Value, isNil bool) reflect.Value {
	for {
		// Load value from interface, but only if the result will be
		// usefully addressable.
		if rv.Kind() == reflect.Interface && !rv.IsNil() {
			e := rv.Elem()
			if e.Kind() == reflect.Ptr && !e.IsNil() && (!isNil || e.Elem().Kind() == reflect.Ptr) {
				rv = e
				continue
			}
		}

		if rv.Kind() != reflect.Ptr {
			break
		}

		if isNil && rv.CanSet() {
			return rv
		}

		if rv.Elem().Kind() == reflect.Interface && rv.Elem().Elem() == rv {
			return rv.Elem()
		}

		if rv.IsNil() {
			rv.Set(reflect.New(rv.Type().Elem()))
		}

		rv = rv.Elem()
	}

	return rv
}

func dereference(v interface{}) reflect.Value {
	rv := reflect.ValueOf(v)

	for rv.Kind() == reflect.Ptr || rv.Kind() == reflect.Interface {

		if rv.Elem().Kind() == reflect.Interface && rv.Elem().Elem() == rv {
			return rv.Elem()
		}

		rv = rv.Elem()
	}

	return rv
}

func isNilValue(v interface{}) bool {
	rv := dereference(v)

	switch rv.Kind() {
	case reflect.Invalid:
		return true
	case reflect.Ptr, reflect.Interface, reflect.Map, reflect.Slice, reflect.Chan, reflect.Func:
		return rv.IsNil()
	}
	return false
}

type GenericEvent struct {
	EventBase
	TypeID     string
	Attributes map[string]ParseResolvable
}

func NewGenericEvent(typeID string) *GenericEvent {
	return &GenericEvent{
		TypeID:     typeID,
		Attributes: make(map[string]ParseResolvable),
	}
}

func (g *GenericEvent) parseField(name string, p ParseResolvable) (err error) {
	g.Attributes[name] = p
	return nil
}

func (g *GenericEvent) Parse(r Reader, classes ClassMap, cpools PoolMap, class *ClassMetadata) error {
	g.Metadata = class
	return parseFields(r, classes, cpools, class, nil, true, g.parseField)
}

func (g *GenericEvent) GetAttr(fieldName string, v interface{}) error {
	rv := reflect.ValueOf(v)

	if rv.Kind() != reflect.Ptr {
		return fmt.Errorf("v must be a pointer")
	}
	if rv.IsNil() {
		return fmt.Errorf("v is a nil pointer")
	}

	attr, ok := g.Attributes[fieldName]
	if !ok {
		return fmt.Errorf("field [%s] not exists in this event", fieldName)
	}

	nilValue := isNilValue(attr)

	rv = indirect(rv, nilValue)

	if nilValue {
		switch rv.Kind() {
		case reflect.Interface, reflect.Ptr, reflect.Map, reflect.Slice:
			rv.Set(reflect.Zero(rv.Type()))
		}
		return nil
	}

	switch rv.Kind() {
	case reflect.Bool:
		x, err := toBoolean(attr)
		if err != nil {
			return fmt.Errorf("unable to resolve boolean: %w", err)
		}
		rv.SetBool(x)

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		var x int64
		switch v := attr.(type) {
		case *Byte:
			x = int64(*v)
		case *Short:
			x = int64(*v)
		case *Int:
			x = int64(*v)
		case *Long:
			x = int64(*v)
		default:
			return fmt.Errorf("unable to assign %T to number", attr)
		}

		if rv.OverflowInt(x) {
			return fmt.Errorf("unable to assign value to %s: number overflow", rv.Type().Name())
		}
		rv.SetInt(x)

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		var x int64
		switch v := attr.(type) {
		case *Byte:
			x = int64(*v)
		case *Short:
			x = int64(*v)
		case *Int:
			x = int64(*v)
		case *Long:
			x = int64(*v)
		default:
			return fmt.Errorf("unable to assign %T to number", attr)
		}
		if x < 0 {
			return fmt.Errorf("unable to assign negative number to unsigned number")
		}
		if rv.OverflowUint(uint64(x)) {
			return fmt.Errorf("unable to assign value to %s: number overflow", rv.Type().Name())
		}
		rv.SetUint(uint64(x))

	case reflect.Float32, reflect.Float64:
		var f64 float64
		switch v := attr.(type) {
		case *Float:
			f64 = float64(*v)
		case *Double:
			f64 = float64(*v)
		default:
			return fmt.Errorf("unable to assign %T to float", attr)
		}
		if rv.OverflowFloat(f64) {
			return fmt.Errorf("unable to assign value to %s: number overflow", rv.Type().Name())
		}
		rv.SetFloat(f64)

	case reflect.String:
		x, err := toString(attr)
		if err != nil {
			return fmt.Errorf("unable to resolve string: %w", err)
		}
		rv.SetString(x)
	case reflect.Struct, reflect.Interface:
		attrValue := dereference(attr)
		if !attrValue.Type().AssignableTo(rv.Type()) {
			return fmt.Errorf("unable to assign value of type %s to type %s", attrValue.Type().Name(), rv.Type().Name())
		}
		rv.Set(attrValue)
	}

	return nil
}

type EventBase struct {
	Metadata *ClassMetadata
}

func (e *EventBase) SetMetadata(metadata *ClassMetadata) {
	e.Metadata = metadata
}

func (e *EventBase) GetMetadata() *ClassMetadata {
	return e.Metadata
}

func ParseEvent(r Reader, classes ClassMap, cpools PoolMap) (Event, error) {
	kind, err := r.VarLong()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve event type: %w", err)
	}
	if kind == MetadataEventType || kind == ConstantPoolEventType {
		return nil, nil
	}
	return parseEvent(r, classes, cpools, int(kind))
}

func parseEvent(r Reader, classes ClassMap, cpools PoolMap, classID int) (Event, error) {
	class, ok := classes[classID]
	if !ok {
		return nil, fmt.Errorf("unknown class %d", classID)
	}
	if class.SuperType != EventSuperType {
		return nil, nil
	}
	var v Event
	//if _, ok := events[class.Name]; ok {
	//	//v = typeFn()
	//	v = NewGenericEvent(class.Name)
	//} else {
	//	v = new(UnsupportedEvent)
	//}
	v = NewGenericEvent(class.Name)
	v.SetMetadata(class)
	if err := v.Parse(r, classes, cpools, class); err != nil {
		return nil, fmt.Errorf("unable to parse event %s: %w", class.Name, err)
	}
	return v, nil
}

type ActiveRecording struct {
	EventBase
	StartTime         int64
	Duration          int64
	EventThread       *Thread
	ID                int64
	Name              string
	Destination       string
	MaxAge            int64
	MaxSize           int64
	RecordingStart    int64
	RecordingDuration int64
}

func (ar *ActiveRecording) parseField(name string, p ParseResolvable) (err error) {
	switch name {
	case "startTime":
		ar.StartTime, err = toLong(p)
	case "duration":
		ar.Duration, err = toLong(p)
	case "eventThread":
		ar.EventThread, err = toThread(p)
	case "id":
		ar.ID, err = toLong(p)
	case "name":
		ar.Name, err = toString(p)
	case "destination":
		ar.Destination, err = toString(p)
	case "maxAge":
		ar.MaxAge, err = toLong(p)
	case "maxSize":
		ar.MaxSize, err = toLong(p)
	case "recordingStart":
		ar.RecordingStart, err = toLong(p)
	case "recordingDuration":
		ar.RecordingDuration, err = toLong(p)
	}
	return err
}

func (ar *ActiveRecording) Parse(r Reader, classes ClassMap, cpools PoolMap, class *ClassMetadata) error {
	return parseFields(r, classes, cpools, class, nil, true, ar.parseField)
}

type ActiveSetting struct {
	EventBase
	StartTime   int64
	Duration    int64
	EventThread *Thread
	ID          int64
	Name        string
	Value       string
}

func (as *ActiveSetting) parseField(name string, p ParseResolvable) (err error) {
	switch name {
	case "startTime":
		as.StartTime, err = toLong(p)
	case "duration":
		as.Duration, err = toLong(p)
	case "eventThread":
		as.EventThread, err = toThread(p)
	case "id":
		as.ID, err = toLong(p)
	case "name":
		as.Name, err = toString(p)
	case "value":
		as.Value, err = toString(p)
	}
	return err
}

func (as *ActiveSetting) Parse(r Reader, classes ClassMap, cpools PoolMap, class *ClassMetadata) error {
	return parseFields(r, classes, cpools, class, nil, true, as.parseField)
}

type BooleanFlag struct {
	EventBase
	StartTime int64
	Name      string
	Value     bool
	Origin    *FlagValueOrigin
}

func (bf *BooleanFlag) parseField(name string, p ParseResolvable) (err error) {
	switch name {
	case "startTime":
		bf.StartTime, err = toLong(p)
	case "name":
		bf.Name, err = toString(p)
	case "value":
		bf.Value, err = toBoolean(p)
	case "origin":
		bf.Origin, err = toFlagValueOrigin(p)
	}
	return err
}

func (bf *BooleanFlag) Parse(r Reader, classes ClassMap, cpools PoolMap, class *ClassMetadata) error {
	return parseFields(r, classes, cpools, class, nil, true, bf.parseField)
}

type CPUInformation struct {
	EventBase
	StartTime   int64
	CPU         string
	Description string
	Sockets     int32
	Cores       int32
	HWThreads   int32
}

func (ci *CPUInformation) parseField(name string, p ParseResolvable) (err error) {
	switch name {
	case "startTime":
		ci.StartTime, err = toLong(p)
	case "duration":
		ci.CPU, err = toString(p)
	case "eventThread":
		ci.Description, err = toString(p)
	case "sockets":
		ci.Sockets, err = toInt(p)
	case "cores":
		ci.Cores, err = toInt(p)
	case "hwThreads":
		ci.HWThreads, err = toInt(p)
	}
	return err
}

func (ci *CPUInformation) Parse(r Reader, classes ClassMap, cpools PoolMap, class *ClassMetadata) error {
	return parseFields(r, classes, cpools, class, nil, true, ci.parseField)
}

type CPULoad struct {
	EventBase
	StartTime    int64
	JVMUser      float32
	JVMSystem    float32
	MachineTotal float32
}

func (cl *CPULoad) parseField(name string, p ParseResolvable) (err error) {
	switch name {
	case "startTime":
		cl.StartTime, err = toLong(p)
	case "jvmUser":
		cl.JVMUser, err = toFloat(p)
	case "jvmSystem":
		cl.JVMSystem, err = toFloat(p)
	case "machineTotal":
		cl.MachineTotal, err = toFloat(p)
	}
	return err
}

func (cl *CPULoad) Parse(r Reader, classes ClassMap, cpools PoolMap, class *ClassMetadata) error {
	return parseFields(r, classes, cpools, class, nil, true, cl.parseField)
}

type CPUTimeStampCounter struct {
	EventBase
	StartTime           int64
	FastTimeEnabled     bool
	FastTimeAutoEnabled bool
	OSFrequency         int64
	FastTimeFrequency   int64
}

func (ctsc *CPUTimeStampCounter) parseField(name string, p ParseResolvable) (err error) {
	switch name {
	case "startTime":
		ctsc.StartTime, err = toLong(p)
	case "fastTimeEnabled":
		ctsc.FastTimeEnabled, err = toBoolean(p)
	case "fastTimeAutoEnabled":
		ctsc.FastTimeAutoEnabled, err = toBoolean(p)
	case "osFrequency":
		ctsc.OSFrequency, err = toLong(p)
	case "fastTimeFrequency":
		ctsc.FastTimeFrequency, err = toLong(p)
	}
	return err
}

func (ctsc *CPUTimeStampCounter) Parse(r Reader, classes ClassMap, cpools PoolMap, class *ClassMetadata) error {
	return parseFields(r, classes, cpools, class, nil, true, ctsc.parseField)
}

type ClassLoaderStatistics struct {
	EventBase
	StartTime                 int64
	ClassLoader               *ClassLoader
	ParentClassLoader         *ClassLoader
	ClassLoaderData           int64
	ClassCount                int64
	ChunkSize                 int64
	BlockSize                 int64
	AnonymousClassCount       int64
	AnonymousChunkSize        int64
	AnonymousBlockSize        int64
	UnsafeAnonymousClassCount int64
	UnsafeAnonymousChunkSize  int64
	UnsafeAnonymousBlockSize  int64
	HiddenClassCount          int64
	HiddenChunkSize           int64
	HiddenBlockSize           int64
}

func (cls *ClassLoaderStatistics) parseField(name string, p ParseResolvable) (err error) {
	switch name {
	case "startTime":
		cls.StartTime, err = toLong(p)
	case "classLoader":
		cls.ClassLoader, err = toClassLoader(p)
	case "parentClassLoader":
		cls.ParentClassLoader, err = toClassLoader(p)
	case "classLoaderData":
		cls.ClassLoaderData, err = toLong(p)
	case "classCount":
		cls.ClassCount, err = toLong(p)
	case "chunkSize":
		cls.ChunkSize, err = toLong(p)
	case "blockSize":
		cls.BlockSize, err = toLong(p)
	case "anonymousClassCount":
		cls.AnonymousClassCount, err = toLong(p)
	case "anonymousChunkSize":
		cls.AnonymousChunkSize, err = toLong(p)
	case "anonymousBlockSize":
		cls.AnonymousBlockSize, err = toLong(p)
	case "unsafeAnonymousClassCount":
		cls.UnsafeAnonymousClassCount, err = toLong(p)
	case "unsafeAnonymousChunkSize":
		cls.UnsafeAnonymousChunkSize, err = toLong(p)
	case "unsafeAnonymousBlockSize":
		cls.UnsafeAnonymousBlockSize, err = toLong(p)
	case "hiddenClassCount":
		cls.HiddenClassCount, err = toLong(p)
	case "hiddenChunkSize":
		cls.HiddenChunkSize, err = toLong(p)
	case "hiddenBlockSize":
		cls.HiddenBlockSize, err = toLong(p)
	}
	return err
}

func (cls *ClassLoaderStatistics) Parse(r Reader, classes ClassMap, cpools PoolMap, class *ClassMetadata) error {
	return parseFields(r, classes, cpools, class, nil, true, cls.parseField)
}

type ClassLoadingStatistics struct {
	EventBase
	StartTime          int64
	LoadedClassCount   int64
	UnloadedClassCount int64
}

func (cls *ClassLoadingStatistics) parseField(name string, p ParseResolvable) (err error) {
	switch name {
	case "startTime":
		cls.StartTime, err = toLong(p)
	case "loadedClassCount":
		cls.LoadedClassCount, err = toLong(p)
	case "unloadedClassCount":
		cls.UnloadedClassCount, err = toLong(p)
	}
	return err
}

func (cls *ClassLoadingStatistics) Parse(r Reader, classes ClassMap, cpools PoolMap, class *ClassMetadata) error {
	return parseFields(r, classes, cpools, class, nil, true, cls.parseField)
}

type CodeCacheConfiguration struct {
	EventBase
	StartTime          int64
	InitialSize        int64
	ReservedSize       int64
	NonNMethodSize     int64
	ProfiledSize       int64
	NonProfiledSize    int64
	ExpansionSize      int64
	MinBlockLength     int64
	StartAddress       int64
	ReservedTopAddress int64
}

func (ccc *CodeCacheConfiguration) parseField(name string, p ParseResolvable) (err error) {
	switch name {
	case "startTime":
		ccc.StartTime, err = toLong(p)
	case "initialSize":
		ccc.InitialSize, err = toLong(p)
	case "reservedSize":
		ccc.ReservedSize, err = toLong(p)
	case "nonNMethodSize":
		ccc.NonNMethodSize, err = toLong(p)
	case "profiledSize":
		ccc.ProfiledSize, err = toLong(p)
	case "NonProfiledSize":
		ccc.NonProfiledSize, err = toLong(p)
	case "ExpansionSize":
		ccc.ExpansionSize, err = toLong(p)
	case "MinBlockLength":
		ccc.MinBlockLength, err = toLong(p)
	case "StartAddress":
		ccc.StartAddress, err = toLong(p)
	case "ReservedTopAddress":
		ccc.ReservedTopAddress, err = toLong(p)
	}
	return err
}

func (ccc *CodeCacheConfiguration) Parse(r Reader, classes ClassMap, cpools PoolMap, class *ClassMetadata) error {
	return parseFields(r, classes, cpools, class, nil, true, ccc.parseField)
}

type CodeCacheStatistics struct {
	EventBase
	StartTime           int64
	CodeBlobType        *CodeBlobType
	StartAddress        int64
	ReservedTopAddress  int64
	EntryCount          int32
	MethodCount         int32
	AdaptorCount        int32
	UnallocatedCapacity int64
	FullCount           int32
}

func (ccs *CodeCacheStatistics) parseField(name string, p ParseResolvable) (err error) {
	switch name {
	case "startTime":
		ccs.StartTime, err = toLong(p)
	case "codeBlobType":
		ccs.CodeBlobType, err = toCodeBlobType(p)
	case "startAddress":
		ccs.StartAddress, err = toLong(p)
	case "reservedTopAddress":
		ccs.ReservedTopAddress, err = toLong(p)
	case "entryCount":
		ccs.EntryCount, err = toInt(p)
	case "methodCount":
		ccs.MethodCount, err = toInt(p)
	case "adaptorCount":
		ccs.AdaptorCount, err = toInt(p)
	case "unallocatedCapacity":
		ccs.UnallocatedCapacity, err = toLong(p)
	case "fullCount":
		ccs.FullCount, err = toInt(p)
	}
	return err
}

func (ccs *CodeCacheStatistics) Parse(r Reader, classes ClassMap, cpools PoolMap, class *ClassMetadata) error {
	return parseFields(r, classes, cpools, class, nil, true, ccs.parseField)
}

type CodeSweeperConfiguration struct {
	EventBase
	StartTime       int64
	SweeperEnabled  bool
	FlushingEnabled bool
	SweepThreshold  int64
}

func (csc *CodeSweeperConfiguration) parseField(name string, p ParseResolvable) (err error) {
	switch name {
	case "startTime":
		csc.StartTime, err = toLong(p)
	case "sweeperEnabled":
		csc.SweeperEnabled, err = toBoolean(p)
	case "flushingEnabled":
		csc.FlushingEnabled, err = toBoolean(p)
	case "sweepThreshold":
		csc.SweepThreshold, err = toLong(p)
	}
	return err
}

func (csc *CodeSweeperConfiguration) Parse(r Reader, classes ClassMap, cpools PoolMap, class *ClassMetadata) error {
	return parseFields(r, classes, cpools, class, nil, true, csc.parseField)
}

type CodeSweeperStatistics struct {
	EventBase
	StartTime            int64
	SweepCount           int32
	MethodReclaimedCount int32
	TotalSweepTime       int64
	PeakFractionTime     int64
	PeakSweepTime        int64
}

func (css *CodeSweeperStatistics) parseField(name string, p ParseResolvable) (err error) {
	switch name {
	case "startTime":
		css.StartTime, err = toLong(p)
	case "sweepCount":
		css.SweepCount, err = toInt(p)
	case "methodReclaimedCount":
		css.MethodReclaimedCount, err = toInt(p)
	case "totalSweepTime":
		css.TotalSweepTime, err = toLong(p)
	case "peakFractionTime":
		css.PeakFractionTime, err = toLong(p)
	case "peakSweepTime":
		css.PeakSweepTime, err = toLong(p)
	}
	return err
}

func (css *CodeSweeperStatistics) Parse(r Reader, classes ClassMap, cpools PoolMap, class *ClassMetadata) error {
	return parseFields(r, classes, cpools, class, nil, true, css.parseField)
}

type CompilerConfiguration struct {
	EventBase
	StartTime         int64
	ThreadCount       int32
	TieredCompilation bool
}

func (cc *CompilerConfiguration) parseField(name string, p ParseResolvable) (err error) {
	switch name {
	case "startTime":
		cc.StartTime, err = toLong(p)
	case "threadCount":
		cc.ThreadCount, err = toInt(p)
	case "tieredCompilation":
		cc.TieredCompilation, err = toBoolean(p)
	}
	return err
}

func (cc *CompilerConfiguration) Parse(r Reader, classes ClassMap, cpools PoolMap, class *ClassMetadata) error {
	return parseFields(r, classes, cpools, class, nil, true, cc.parseField)
}

type CompilerStatistics struct {
	EventBase
	StartTime             int64
	CompileCount          int32
	BailoutCount          int32
	InvalidatedCount      int32
	OSRCompileCount       int32
	StandardCompileCount  int32
	OSRBytesCompiled      int64
	StandardBytesCompiled int64
	NMethodsSize          int64
	NMethodCodeSize       int64
	PeakTimeSpent         int64
	TotalTimeSpent        int64
}

func (cs *CompilerStatistics) parseField(name string, p ParseResolvable) (err error) {
	switch name {
	case "startTime":
		cs.StartTime, err = toLong(p)
	case "compileCount":
		cs.CompileCount, err = toInt(p)
	case "bailoutCount":
		cs.BailoutCount, err = toInt(p)
	case "invalidatedCount":
		cs.InvalidatedCount, err = toInt(p)
	case "osrCompileCount":
		cs.OSRCompileCount, err = toInt(p)
	case "standardCompileCount":
		cs.StandardCompileCount, err = toInt(p)
	case "osrBytesCompiled":
		cs.OSRBytesCompiled, err = toLong(p)
	case "standardBytesCompiled":
		cs.StandardBytesCompiled, err = toLong(p)
	case "nmethodsSize":
		cs.NMethodsSize, err = toLong(p)
	case "nmethodCodeSize":
		cs.NMethodCodeSize, err = toLong(p)
	case "peakTimeSpent":
		cs.PeakTimeSpent, err = toLong(p)
	case "totalTimeSpent":
		cs.TotalTimeSpent, err = toLong(p)
	}
	return err
}

func (cs *CompilerStatistics) Parse(r Reader, classes ClassMap, cpools PoolMap, class *ClassMetadata) error {
	return parseFields(r, classes, cpools, class, nil, true, cs.parseField)
}

type DoubleFlag struct {
	EventBase
	StartTime int64
	Name      string
	Value     float64
	Origin    *FlagValueOrigin
}

func (df *DoubleFlag) parseField(name string, p ParseResolvable) (err error) {
	switch name {
	case "startTime":
		df.StartTime, err = toLong(p)
	case "name":
		df.Name, err = toString(p)
	case "value":
		df.Value, err = toDouble(p)
	case "origin":
		df.Origin, err = toFlagValueOrigin(p)
	}
	return err
}

func (df *DoubleFlag) Parse(r Reader, classes ClassMap, cpools PoolMap, class *ClassMetadata) error {
	return parseFields(r, classes, cpools, class, nil, true, df.parseField)
}

type ExceptionStatistics struct {
	EventBase
	StartTime   int64
	Duration    int64
	EventThread *Thread
	StackTrace  *StackTrace
	Throwable   int64
}

func (es *ExceptionStatistics) parseField(name string, p ParseResolvable) (err error) {
	switch name {
	case "startTime":
		es.StartTime, err = toLong(p)
	case "duration":
		es.Duration, err = toLong(p)
	case "eventThread":
		es.EventThread, err = toThread(p)
	case "stackTrace":
		es.StackTrace, err = toStackTrace(p)
	case "throwable":
		es.Throwable, err = toLong(p)
	}
	return err
}

func (es *ExceptionStatistics) Parse(r Reader, classes ClassMap, cpools PoolMap, class *ClassMetadata) error {
	return parseFields(r, classes, cpools, class, nil, true, es.parseField)
}

type ExecutionSample struct {
	EventBase
	StartTime     int64
	SampledThread *Thread
	StackTrace    *StackTrace
	State         *ThreadState
	ContextId     int64
}

func (es *ExecutionSample) parseField(name string, p ParseResolvable) (err error) {
	switch name {
	case "startTime":
		es.StartTime, err = toLong(p)
	case "sampledThread":
		es.SampledThread, err = toThread(p)
	case "stackTrace":
		es.StackTrace, err = toStackTrace(p)
	case "state":
		es.State, err = toThreadState(p)
	case "contextId":
		es.ContextId, err = toLong(p)
	}
	return err
}

func (es *ExecutionSample) Parse(r Reader, classes ClassMap, cpools PoolMap, class *ClassMetadata) error {
	return parseFields(r, classes, cpools, class, nil, true, es.parseField)
}

type GCConfiguration struct {
	EventBase
	StartTime              int64
	YoungCollector         *GCName
	OldCollector           *GCName
	ParallelGCThreads      int32
	ConcurrentGCThreads    int32
	UsesDynamicGCThreads   bool
	IsExplicitGCConcurrent bool
	IsExplicitGCDisabled   bool
	PauseTarget            int64
	GCTimeRatio            int32
}

func (gc *GCConfiguration) parseField(name string, p ParseResolvable) (err error) {
	switch name {
	case "startTime":
		gc.StartTime, err = toLong(p)
	case "youngCollector":
		gc.YoungCollector, err = toGCName(p)
	case "oldCollector":
		gc.OldCollector, err = toGCName(p)
	case "parallelGCThreads":
		gc.ParallelGCThreads, err = toInt(p)
	case "concurrentGCThreads":
		gc.ConcurrentGCThreads, err = toInt(p)
	case "usesDynamicGCThreads":
		gc.UsesDynamicGCThreads, err = toBoolean(p)
	case "isExplicitGCConcurrent":
		gc.IsExplicitGCConcurrent, err = toBoolean(p)
	case "isExplicitGCDisabled":
		gc.IsExplicitGCDisabled, err = toBoolean(p)
	case "pauseTarget":
		gc.PauseTarget, err = toLong(p)
	case "gcTimeRatio":
		gc.GCTimeRatio, err = toInt(p)
	}
	return err
}

func (gc *GCConfiguration) Parse(r Reader, classes ClassMap, cpools PoolMap, class *ClassMetadata) error {
	return parseFields(r, classes, cpools, class, nil, true, gc.parseField)
}

type GCHeapConfiguration struct {
	EventBase
	StartTime          int64
	MinSize            int64
	MaxSize            int64
	InitialSize        int64
	UsesCompressedOops bool
	CompressedOopsMode *NarrowOopMode
	ObjectAlignment    int64
	HeapAddressBits    int8
}

func (ghc *GCHeapConfiguration) parseField(name string, p ParseResolvable) (err error) {
	switch name {
	case "startTime":
		ghc.StartTime, err = toLong(p)
	case "minSize":
		ghc.MinSize, err = toLong(p)
	case "maxSize":
		ghc.MaxSize, err = toLong(p)
	case "initialSize":
		ghc.InitialSize, err = toLong(p)
	case "usesCompressedOops":
		ghc.UsesCompressedOops, err = toBoolean(p)
	case "compressedOopsMode":
		ghc.CompressedOopsMode, err = toNarrowOopMode(p)
	case "objectAlignment":
		ghc.ObjectAlignment, err = toLong(p)
	case "heapAddressBits":
		ghc.HeapAddressBits, err = toByte(p)
	}
	return err
}

func (ghc *GCHeapConfiguration) Parse(r Reader, classes ClassMap, cpools PoolMap, class *ClassMetadata) error {
	return parseFields(r, classes, cpools, class, nil, true, ghc.parseField)
}

type GCSurvivorConfiguration struct {
	EventBase
	StartTime                int64
	MaxTenuringThreshold     int8
	InitialTenuringThreshold int8
}

func (gcs *GCSurvivorConfiguration) parseField(name string, p ParseResolvable) (err error) {
	switch name {
	case "startTime":
		gcs.StartTime, err = toLong(p)
	case "maxTenuringThreshold":
		gcs.MaxTenuringThreshold, err = toByte(p)
	case "initialTenuringThreshold":
		gcs.InitialTenuringThreshold, err = toByte(p)
	}
	return err
}

func (gsc *GCSurvivorConfiguration) Parse(r Reader, classes ClassMap, cpools PoolMap, class *ClassMetadata) error {
	return parseFields(r, classes, cpools, class, nil, true, gsc.parseField)
}

type GCTLABConfiguration struct {
	EventBase
	StartTime            int64
	UsesTLABs            bool
	MinTLABSize          int64
	TLABRefillWasteLimit int64
}

func (gtc *GCTLABConfiguration) parseField(name string, p ParseResolvable) (err error) {
	switch name {
	case "startTime":
		gtc.StartTime, err = toLong(p)
	case "usesTLABs":
		gtc.UsesTLABs, err = toBoolean(p)
	case "minTLABSize":
		gtc.MinTLABSize, err = toLong(p)
	case "tlabRefillWasteLimit":
		gtc.TLABRefillWasteLimit, err = toLong(p)
	}
	return err
}

func (gtc *GCTLABConfiguration) Parse(r Reader, classes ClassMap, cpools PoolMap, class *ClassMetadata) error {
	return parseFields(r, classes, cpools, class, nil, true, gtc.parseField)
}

type InitialEnvironmentVariable struct {
	EventBase
	StartTime int64
	Key       string
	Value     string
}

func (iev *InitialEnvironmentVariable) parseField(name string, p ParseResolvable) (err error) {
	switch name {
	case "startTime":
		iev.StartTime, err = toLong(p)
	case "key":
		iev.Key, err = toString(p)
	case "value":
		iev.Value, err = toString(p)
	}
	return err
}

func (iev *InitialEnvironmentVariable) Parse(r Reader, classes ClassMap, cpools PoolMap, class *ClassMetadata) error {
	return parseFields(r, classes, cpools, class, nil, true, iev.parseField)
}

type InitialSystemProperty struct {
	EventBase
	StartTime int64
	Key       string
	Value     string
}

func (isp *InitialSystemProperty) parseField(name string, p ParseResolvable) (err error) {
	switch name {
	case "startTime":
		isp.StartTime, err = toLong(p)
	case "key":
		isp.Key, err = toString(p)
	case "value":
		isp.Value, err = toString(p)
	}
	return err
}

func (isp *InitialSystemProperty) Parse(r Reader, classes ClassMap, cpools PoolMap, class *ClassMetadata) error {
	return parseFields(r, classes, cpools, class, nil, true, isp.parseField)
}

type IntFlag struct {
	EventBase
	StartTime int64
	Name      string
	Value     int32
	Origin    *FlagValueOrigin
}

func (f *IntFlag) parseField(name string, p ParseResolvable) (err error) {
	switch name {
	case "startTime":
		f.StartTime, err = toLong(p)
	case "name":
		f.Name, err = toString(p)
	case "value":
		f.Value, err = toInt(p)
	case "origin":
		f.Origin, err = toFlagValueOrigin(p)
	}
	return err
}

func (f *IntFlag) Parse(r Reader, classes ClassMap, cpools PoolMap, class *ClassMetadata) error {
	return parseFields(r, classes, cpools, class, nil, true, f.parseField)
}

type JavaMonitorEnter struct {
	EventBase
	StartTime     int64
	Duration      int64
	EventThread   *Thread
	StackTrace    *StackTrace
	MonitorClass  *Class
	PreviousOwner *Thread
	Address       int64
	ContextId     int64
}

func (jme *JavaMonitorEnter) parseField(name string, p ParseResolvable) (err error) {
	switch name {
	case "startTime":
		jme.StartTime, err = toLong(p)
	case "duration":
		jme.Duration, err = toLong(p)
	case "eventThread":
		jme.EventThread, err = toThread(p)
	case "stackTrace":
		jme.StackTrace, err = toStackTrace(p)
	case "monitorClass":
		jme.MonitorClass, err = toClass(p)
	case "previousOwner":
		jme.PreviousOwner, err = toThread(p)
	case "address":
		jme.Address, err = toLong(p)
	case "contextId":
		jme.ContextId, err = toLong(p)
	}
	return err
}

func (jme *JavaMonitorEnter) Parse(r Reader, classes ClassMap, cpools PoolMap, class *ClassMetadata) error {
	return parseFields(r, classes, cpools, class, nil, true, jme.parseField)
}

type JavaMonitorWait struct {
	EventBase
	StartTime    int64
	Duration     int64
	EventThread  *Thread
	StackTrace   *StackTrace
	MonitorClass *Class
	Notifier     *Thread
	Timeout      int64
	TimedOut     bool
	Address      int64
}

func (jmw *JavaMonitorWait) parseField(name string, p ParseResolvable) (err error) {
	switch name {
	case "startTime":
		jmw.StartTime, err = toLong(p)
	case "duration":
		jmw.Duration, err = toLong(p)
	case "eventThread":
		jmw.EventThread, err = toThread(p)
	case "stackTrace":
		jmw.StackTrace, err = toStackTrace(p)
	case "monitorClass":
		jmw.MonitorClass, err = toClass(p)
	case "notifier":
		jmw.Notifier, err = toThread(p)
	case "timeout":
		jmw.Timeout, err = toLong(p)
	case "timedOut":
		jmw.TimedOut, err = toBoolean(p)
	case "address":
		jmw.Address, err = toLong(p)
	}
	return err
}

func (jmw *JavaMonitorWait) Parse(r Reader, classes ClassMap, cpools PoolMap, class *ClassMetadata) error {
	return parseFields(r, classes, cpools, class, nil, true, jmw.parseField)
}

type JavaThreadStatistics struct {
	EventBase
	StartTime        int64
	ActiveCount      int64
	DaemonCount      int64
	AccumulatedCount int64
	PeakCount        int64
}

func (jts *JavaThreadStatistics) parseField(name string, p ParseResolvable) (err error) {
	switch name {
	case "startTime":
		jts.StartTime, err = toLong(p)
	case "activeCount":
		jts.ActiveCount, err = toLong(p)
	case "daemonCount":
		jts.DaemonCount, err = toLong(p)
	case "accumulatedCount":
		jts.AccumulatedCount, err = toLong(p)
	case "peakCount":
		jts.PeakCount, err = toLong(p)
	}
	return err
}

func (jts *JavaThreadStatistics) Parse(r Reader, classes ClassMap, cpools PoolMap, class *ClassMetadata) error {
	return parseFields(r, classes, cpools, class, nil, true, jts.parseField)
}

type JVMInformation struct {
	EventBase
	StartTime     int64
	JVMName       string
	JVMVersion    string
	JVMArguments  string
	JVMFlags      string
	JavaArguments string
	JVMStartTime  int64
	PID           int64
}

func (ji *JVMInformation) parseField(name string, p ParseResolvable) (err error) {
	switch name {
	case "startTime":
		ji.StartTime, err = toLong(p)
	case "jvmName":
		ji.JVMName, err = toString(p)
	case "jvmVersion":
		ji.JVMVersion, err = toString(p)
	case "jvmArguments":
		ji.JVMArguments, err = toString(p)
	case "jvmFlags":
		ji.JVMFlags, err = toString(p)
	case "javaArguments":
		ji.JavaArguments, err = toString(p)
	case "jvmStartTime":
		ji.JVMStartTime, err = toLong(p)
	case "pid":
		ji.PID, err = toLong(p)
	}
	return err
}

func (ji *JVMInformation) Parse(r Reader, classes ClassMap, cpools PoolMap, class *ClassMetadata) error {
	return parseFields(r, classes, cpools, class, nil, true, ji.parseField)
}

type LoaderConstraintsTableStatistics struct {
	EventBase
	StartTime                    int64
	BucketCount                  int64
	EntryCount                   int64
	TotalFootprint               int64
	BucketCountMaximum           int64
	BucketCountAverage           float32
	BucketCountVariance          float32
	BucketCountStandardDeviation float32
	InsertionRate                float32
	RemovalRate                  float32
}

func (lcts *LoaderConstraintsTableStatistics) parseField(name string, p ParseResolvable) (err error) {
	switch name {
	case "startTime":
		lcts.StartTime, err = toLong(p)
	case "bucketCount":
		lcts.BucketCount, err = toLong(p)
	case "entryCount":
		lcts.EntryCount, err = toLong(p)
	case "totalFootprint":
		lcts.TotalFootprint, err = toLong(p)
	case "bucketCountMaximum":
		lcts.BucketCountMaximum, err = toLong(p)
	case "bucketCountAverage":
		lcts.BucketCountAverage, err = toFloat(p)
	case "bucketCountVariance":
		lcts.BucketCountVariance, err = toFloat(p)
	case "bucketCountStandardDeviation":
		lcts.BucketCountStandardDeviation, err = toFloat(p)
	case "insertionRate":
		lcts.InsertionRate, err = toFloat(p)
	case "removalRate":
		lcts.RemovalRate, err = toFloat(p)
	}
	return err
}

func (lcts *LoaderConstraintsTableStatistics) Parse(r Reader, classes ClassMap, cpools PoolMap, class *ClassMetadata) error {
	return parseFields(r, classes, cpools, class, nil, true, lcts.parseField)
}

type LongFlag struct {
	EventBase
	StartTime int64
	Name      string
	Value     int64
	Origin    *FlagValueOrigin
}

func (lf *LongFlag) parseField(name string, p ParseResolvable) (err error) {
	switch name {
	case "startTime":
		lf.StartTime, err = toLong(p)
	case "name":
		lf.Name, err = toString(p)
	case "value":
		lf.Value, err = toLong(p)
	case "origin":
		lf.Origin, err = toFlagValueOrigin(p)
	}
	return err
}

func (lf *LongFlag) Parse(r Reader, classes ClassMap, cpools PoolMap, class *ClassMetadata) error {
	return parseFields(r, classes, cpools, class, nil, true, lf.parseField)
}

type ModuleExport struct {
	EventBase
	StartTime       int64
	ExportedPackage *Package
	TargetModule    *Module
}

func (me *ModuleExport) parseField(name string, p ParseResolvable) (err error) {
	switch name {
	case "startTime":
		me.StartTime, err = toLong(p)
	case "exportedPackage":
		me.ExportedPackage, err = toPackage(p)
	case "targetModule":
		me.TargetModule, err = toModule(p)
	}
	return err
}

func (me *ModuleExport) Parse(r Reader, classes ClassMap, cpools PoolMap, class *ClassMetadata) error {
	return parseFields(r, classes, cpools, class, nil, true, me.parseField)
}

type ModuleRequire struct {
	EventBase
	StartTime      int64
	Source         *Module
	RequiredModule *Module
}

func (mr *ModuleRequire) parseField(name string, p ParseResolvable) (err error) {
	switch name {
	case "startTime":
		mr.StartTime, err = toLong(p)
	case "sourced":
		mr.Source, err = toModule(p)
	case "requiredModule":
		mr.RequiredModule, err = toModule(p)
	}
	return err
}

func (mr *ModuleRequire) Parse(r Reader, classes ClassMap, cpools PoolMap, class *ClassMetadata) error {
	return parseFields(r, classes, cpools, class, nil, true, mr.parseField)
}

type NativeLibrary struct {
	EventBase
	StartTime   int64
	Name        string
	BaseAddress int64
	TopAddress  int64
}

func (nl *NativeLibrary) parseField(name string, p ParseResolvable) (err error) {
	switch name {
	case "startTime":
		nl.StartTime, err = toLong(p)
	case "name":
		nl.Name, err = toString(p)
	case "baseAddress":
		nl.BaseAddress, err = toLong(p)
	case "topAddress":
		nl.TopAddress, err = toLong(p)
	}
	return err
}

func (nl *NativeLibrary) Parse(r Reader, classes ClassMap, cpools PoolMap, class *ClassMetadata) error {
	return parseFields(r, classes, cpools, class, nil, true, nl.parseField)
}

type NetworkUtilization struct {
	EventBase
	StartTime        int64
	NetworkInterface *NetworkInterfaceName
	ReadRate         int64
	WriteRate        int64
}

func (nu *NetworkUtilization) parseField(name string, p ParseResolvable) (err error) {
	switch name {
	case "startTime":
		nu.StartTime, err = toLong(p)
	case "networkInterface":
		nu.NetworkInterface, err = toNetworkInterfaceName(p)
	case "readRate":
		nu.ReadRate, err = toLong(p)
	case "writeRate":
		nu.WriteRate, err = toLong(p)
	}
	return err
}

func (nu *NetworkUtilization) Parse(r Reader, classes ClassMap, cpools PoolMap, class *ClassMetadata) error {
	return parseFields(r, classes, cpools, class, nil, true, nu.parseField)
}

type ObjectAllocationInNewTLAB struct {
	EventBase
	StartTime      int64
	EventThread    *Thread
	StackTrace     *StackTrace
	ObjectClass    *Class
	AllocationSize int64
	TLABSize       int64
	ContextId      int64
}

func (oa *ObjectAllocationInNewTLAB) parseField(name string, p ParseResolvable) (err error) {
	switch name {
	case "startTime":
		oa.StartTime, err = toLong(p)
	case "sampledThread":
		oa.EventThread, err = toThread(p)
	case "stackTrace":
		oa.StackTrace, err = toStackTrace(p)
	case "objectClass":
		oa.ObjectClass, err = toClass(p)
	case "allocationSize":
		oa.AllocationSize, err = toLong(p)
	case "tlabSize":
		oa.TLABSize, err = toLong(p)
	case "contextId":
		oa.ContextId, err = toLong(p)
	}

	return err
}

func (oa *ObjectAllocationInNewTLAB) Parse(r Reader, classes ClassMap, cpools PoolMap, class *ClassMetadata) error {
	return parseFields(r, classes, cpools, class, nil, true, oa.parseField)
}

type ObjectAllocationOutsideTLAB struct {
	EventBase
	StartTime      int64
	EventThread    *Thread
	StackTrace     *StackTrace
	ObjectClass    *Class
	AllocationSize int64
	ContextId      int64
}

func (oa *ObjectAllocationOutsideTLAB) parseField(name string, p ParseResolvable) (err error) {
	switch name {
	case "startTime":
		oa.StartTime, err = toLong(p)
	case "sampledThread":
		oa.EventThread, err = toThread(p)
	case "stackTrace":
		oa.StackTrace, err = toStackTrace(p)
	case "objectClass":
		oa.ObjectClass, err = toClass(p)
	case "allocationSize":
		oa.AllocationSize, err = toLong(p)
	case "contextId":
		oa.ContextId, err = toLong(p)
	}
	return err
}

func (oa *ObjectAllocationOutsideTLAB) Parse(r Reader, classes ClassMap, cpools PoolMap, class *ClassMetadata) error {
	return parseFields(r, classes, cpools, class, nil, true, oa.parseField)
}

type OSInformation struct {
	EventBase
	StartTime int64
	OSVersion string
}

func (os *OSInformation) parseField(name string, p ParseResolvable) (err error) {
	switch name {
	case "startTime":
		os.StartTime, err = toLong(p)
	case "osVersion":
		os.OSVersion, err = toString(p)
	}
	return err
}

func (os *OSInformation) Parse(r Reader, classes ClassMap, cpools PoolMap, class *ClassMetadata) error {
	return parseFields(r, classes, cpools, class, nil, true, os.parseField)
}

type PhysicalMemory struct {
	EventBase
	StartTime int64
	TotalSize int64
	UsedSize  int64
}

func (pm *PhysicalMemory) parseField(name string, p ParseResolvable) (err error) {
	switch name {
	case "startTime":
		pm.StartTime, err = toLong(p)
	case "totalSize":
		pm.TotalSize, err = toLong(p)
	case "usedSize":
		pm.UsedSize, err = toLong(p)
	}
	return err
}

func (pm *PhysicalMemory) Parse(r Reader, classes ClassMap, cpools PoolMap, class *ClassMetadata) error {
	return parseFields(r, classes, cpools, class, nil, true, pm.parseField)
}

type PlaceholderTableStatistics struct {
	EventBase
	StartTime                    int64
	BucketCount                  int64
	EntryCount                   int64
	TotalFootprint               int64
	BucketCountMaximum           int64
	BucketCountAverage           float32
	BucketCountVariance          float32
	BucketCountStandardDeviation float32
	InsertionRate                float32
	RemovalRate                  float32
}

func (pts *PlaceholderTableStatistics) parseField(name string, p ParseResolvable) (err error) {
	switch name {
	case "startTime":
		pts.StartTime, err = toLong(p)
	case "bucketCount":
		pts.BucketCount, err = toLong(p)
	case "entryCount":
		pts.EntryCount, err = toLong(p)
	case "totalFootprint":
		pts.TotalFootprint, err = toLong(p)
	case "bucketCountMaximum":
		pts.BucketCountMaximum, err = toLong(p)
	case "bucketCountAverage":
		pts.BucketCountAverage, err = toFloat(p)
	case "bucketCountVariance":
		pts.BucketCountVariance, err = toFloat(p)
	case "bucketCountStandardDeviation":
		pts.BucketCountStandardDeviation, err = toFloat(p)
	case "insertionRate":
		pts.InsertionRate, err = toFloat(p)
	case "removalRate":
		pts.RemovalRate, err = toFloat(p)
	}
	return err
}

func (pts *PlaceholderTableStatistics) Parse(r Reader, classes ClassMap, cpools PoolMap, class *ClassMetadata) error {
	return parseFields(r, classes, cpools, class, nil, true, pts.parseField)
}

type ProtectionDomainCacheTableStatistics struct {
	EventBase
	StartTime                    int64
	BucketCount                  int64
	EntryCount                   int64
	TotalFootprint               int64
	BucketCountMaximum           int64
	BucketCountAverage           float32
	BucketCountVariance          float32
	BucketCountStandardDeviation float32
	InsertionRate                float32
	RemovalRate                  float32
}

func (pdcts *ProtectionDomainCacheTableStatistics) parseField(name string, p ParseResolvable) (err error) {
	switch name {
	case "startTime":
		pdcts.StartTime, err = toLong(p)
	case "bucketCount":
		pdcts.BucketCount, err = toLong(p)
	case "entryCount":
		pdcts.EntryCount, err = toLong(p)
	case "totalFootprint":
		pdcts.TotalFootprint, err = toLong(p)
	case "bucketCountMaximum":
		pdcts.BucketCountMaximum, err = toLong(p)
	case "bucketCountAverage":
		pdcts.BucketCountAverage, err = toFloat(p)
	case "bucketCountVariance":
		pdcts.BucketCountVariance, err = toFloat(p)
	case "bucketCountStandardDeviation":
		pdcts.BucketCountStandardDeviation, err = toFloat(p)
	case "insertionRate":
		pdcts.InsertionRate, err = toFloat(p)
	case "removalRate":
		pdcts.RemovalRate, err = toFloat(p)
	}
	return err
}

func (pdcts *ProtectionDomainCacheTableStatistics) Parse(r Reader, classes ClassMap, cpools PoolMap, class *ClassMetadata) error {
	return parseFields(r, classes, cpools, class, nil, true, pdcts.parseField)
}

type StringFlag struct {
	EventBase
	StartTime int64
	Name      string
	Value     string
	Origin    *FlagValueOrigin
}

func (sf *StringFlag) parseField(name string, p ParseResolvable) (err error) {
	switch name {
	case "startTime":
		sf.StartTime, err = toLong(p)
	case "name":
		sf.Name, err = toString(p)
	case "value":
		sf.Value, err = toString(p)
	case "origin":
		sf.Origin, err = toFlagValueOrigin(p)
	}
	return err
}

func (sf *StringFlag) Parse(r Reader, classes ClassMap, cpools PoolMap, class *ClassMetadata) error {
	return parseFields(r, classes, cpools, class, nil, true, sf.parseField)
}

type StringTableStatistics struct {
	EventBase
	StartTime                    int64
	BucketCount                  int64
	EntryCount                   int64
	TotalFootprint               int64
	BucketCountMaximum           int64
	BucketCountAverage           float32
	BucketCountVariance          float32
	BucketCountStandardDeviation float32
	InsertionRate                float32
	RemovalRate                  float32
}

func (sts *StringTableStatistics) parseField(name string, p ParseResolvable) (err error) {
	switch name {
	case "startTime":
		sts.StartTime, err = toLong(p)
	case "bucketCount":
		sts.BucketCount, err = toLong(p)
	case "entryCount":
		sts.EntryCount, err = toLong(p)
	case "totalFootprint":
		sts.TotalFootprint, err = toLong(p)
	case "bucketCountMaximum":
		sts.BucketCountMaximum, err = toLong(p)
	case "bucketCountAverage":
		sts.BucketCountAverage, err = toFloat(p)
	case "bucketCountVariance":
		sts.BucketCountVariance, err = toFloat(p)
	case "bucketCountStandardDeviation":
		sts.BucketCountStandardDeviation, err = toFloat(p)
	case "insertionRate":
		sts.InsertionRate, err = toFloat(p)
	case "removalRate":
		sts.RemovalRate, err = toFloat(p)
	}
	return err
}

func (sts *StringTableStatistics) Parse(r Reader, classes ClassMap, cpools PoolMap, class *ClassMetadata) error {
	return parseFields(r, classes, cpools, class, nil, true, sts.parseField)
}

type SymbolTableStatistics struct {
	EventBase
	StartTime                    int64
	BucketCount                  int64
	EntryCount                   int64
	TotalFootprint               int64
	BucketCountMaximum           int64
	BucketCountAverage           float32
	BucketCountVariance          float32
	BucketCountStandardDeviation float32
	InsertionRate                float32
	RemovalRate                  float32
}

func (sts *SymbolTableStatistics) parseField(name string, p ParseResolvable) (err error) {
	switch name {
	case "startTime":
		sts.StartTime, err = toLong(p)
	case "bucketCount":
		sts.BucketCount, err = toLong(p)
	case "entryCount":
		sts.EntryCount, err = toLong(p)
	case "totalFootprint":
		sts.TotalFootprint, err = toLong(p)
	case "bucketCountMaximum":
		sts.BucketCountMaximum, err = toLong(p)
	case "bucketCountAverage":
		sts.BucketCountAverage, err = toFloat(p)
	case "bucketCountVariance":
		sts.BucketCountVariance, err = toFloat(p)
	case "bucketCountStandardDeviation":
		sts.BucketCountStandardDeviation, err = toFloat(p)
	case "insertionRate":
		sts.InsertionRate, err = toFloat(p)
	case "removalRate":
		sts.RemovalRate, err = toFloat(p)
	}
	return err
}

func (sts *SymbolTableStatistics) Parse(r Reader, classes ClassMap, cpools PoolMap, class *ClassMetadata) error {
	return parseFields(r, classes, cpools, class, nil, true, sts.parseField)
}

type SystemProcess struct {
	EventBase
	StartTime   int64
	PID         string
	CommandLine string
}

func (sp *SystemProcess) parseField(name string, p ParseResolvable) (err error) {
	switch name {
	case "startTime":
		sp.StartTime, err = toLong(p)
	case "pid":
		sp.PID, err = toString(p)
	case "commandLine":
		sp.CommandLine, err = toString(p)
	}
	return err
}

func (sp *SystemProcess) Parse(r Reader, classes ClassMap, cpools PoolMap, class *ClassMetadata) error {
	return parseFields(r, classes, cpools, class, nil, true, sp.parseField)
}

type ThreadAllocationStatistics struct {
	EventBase
	StartTime int64
	Allocated int64
	Thread    *Thread
}

func (tas *ThreadAllocationStatistics) parseField(name string, p ParseResolvable) (err error) {
	switch name {
	case "startTime":
		tas.StartTime, err = toLong(p)
	case "allocated":
		tas.Allocated, err = toLong(p)
	case "thread":
		tas.Thread, err = toThread(p)
	}
	return err
}

func (tas *ThreadAllocationStatistics) Parse(r Reader, classes ClassMap, cpools PoolMap, class *ClassMetadata) error {
	return parseFields(r, classes, cpools, class, nil, true, tas.parseField)
}

type ThreadCPULoad struct {
	EventBase
	StartTime   int64
	EventThread *Thread
	User        float32
	System      float32
}

func (tcl *ThreadCPULoad) parseField(name string, p ParseResolvable) (err error) {
	switch name {
	case "startTime":
		tcl.StartTime, err = toLong(p)
	case "eventThread":
		tcl.EventThread, err = toThread(p)
	case "user":
		tcl.User, err = toFloat(p)
	case "system":
		tcl.System, err = toFloat(p)
	}
	return err
}

func (tcl *ThreadCPULoad) Parse(r Reader, classes ClassMap, cpools PoolMap, class *ClassMetadata) error {
	return parseFields(r, classes, cpools, class, nil, true, tcl.parseField)
}

type ThreadContextSwitchRate struct {
	EventBase
	StartTime  int64
	SwitchRate float32
}

func (tcsr *ThreadContextSwitchRate) parseField(name string, p ParseResolvable) (err error) {
	switch name {
	case "startTime":
		tcsr.StartTime, err = toLong(p)
	case "switchRate":
		tcsr.SwitchRate, err = toFloat(p)
	}
	return err
}

func (tcsr *ThreadContextSwitchRate) Parse(r Reader, classes ClassMap, cpools PoolMap, class *ClassMetadata) error {
	return parseFields(r, classes, cpools, class, nil, true, tcsr.parseField)
}

type ThreadDump struct {
	EventBase
	StartTime int64
	Result    string
}

func (td *ThreadDump) parseField(name string, p ParseResolvable) (err error) {
	switch name {
	case "startTime":
		td.StartTime, err = toLong(p)
	case "result":
		td.Result, err = toString(p)
	}
	return err
}

func (td *ThreadDump) Parse(r Reader, classes ClassMap, cpools PoolMap, class *ClassMetadata) error {
	return parseFields(r, classes, cpools, class, nil, true, td.parseField)
}

type ThreadPark struct {
	EventBase
	StartTime   int64
	Duration    int64
	EventThread *Thread
	StackTrace  *StackTrace
	ParkedClass *Class
	Timeout     int64
	Until       int64
	Address     int64
	ContextId   int64
}

func (tp *ThreadPark) parseField(name string, p ParseResolvable) (err error) {
	switch name {
	case "startTime":
		tp.StartTime, err = toLong(p)
	case "duration":
		tp.Duration, err = toLong(p)
	case "eventThread":
		tp.EventThread, err = toThread(p)
	case "stackTrace":
		tp.StackTrace, err = toStackTrace(p)
	case "parkedClass":
		tp.ParkedClass, err = toClass(p)
	case "timeout":
		tp.Timeout, err = toLong(p)
	case "until":
		tp.Until, err = toLong(p)
	case "address":
		tp.Address, err = toLong(p)
	case "contextId": // todo this one seems to be unimplemented in the profiler yet
		tp.ContextId, err = toLong(p)
	}
	return err
}

func (tp *ThreadPark) Parse(r Reader, classes ClassMap, cpools PoolMap, class *ClassMetadata) error {
	return parseFields(r, classes, cpools, class, nil, true, tp.parseField)
}

type ThreadStart struct {
	EventBase
	StartTime    int64
	EventThread  *Thread
	StackTrace   *StackTrace
	Thread       *Thread
	ParentThread *Thread
}

func (ts *ThreadStart) parseField(name string, p ParseResolvable) (err error) {
	switch name {
	case "startTime":
		ts.StartTime, err = toLong(p)
	case "eventThread":
		ts.EventThread, err = toThread(p)
	case "stackTrace":
		ts.StackTrace, err = toStackTrace(p)
	case "thread":
		ts.Thread, err = toThread(p)
	case "parentThread":
		ts.ParentThread, err = toThread(p)
	}
	return err
}

func (ts *ThreadStart) Parse(r Reader, classes ClassMap, cpools PoolMap, class *ClassMetadata) error {
	return parseFields(r, classes, cpools, class, nil, true, ts.parseField)
}

type UnsignedIntFlag struct {
	EventBase
	StartTime int64
	Name      string
	Value     int32
	Origin    *FlagValueOrigin
}

func (uif *UnsignedIntFlag) parseField(name string, p ParseResolvable) (err error) {
	switch name {
	case "startTime":
		uif.StartTime, err = toLong(p)
	case "name":
		uif.Name, err = toString(p)
	case "value":
		uif.Value, err = toInt(p)
	case "origin":
		uif.Origin, err = toFlagValueOrigin(p)
	}
	return err
}

func (uif *UnsignedIntFlag) Parse(r Reader, classes ClassMap, cpools PoolMap, class *ClassMetadata) error {
	return parseFields(r, classes, cpools, class, nil, true, uif.parseField)
}

type UnsignedLongFlag struct {
	EventBase
	StartTime int64
	Name      string
	Value     int64
	Origin    *FlagValueOrigin
}

func (ulf *UnsignedLongFlag) parseField(name string, p ParseResolvable) (err error) {
	switch name {
	case "startTime":
		ulf.StartTime, err = toLong(p)
	case "name":
		ulf.Name, err = toString(p)
	case "value":
		ulf.Value, err = toLong(p)
	case "origin":
		ulf.Origin, err = toFlagValueOrigin(p)
	}
	return err
}

func (ulf *UnsignedLongFlag) Parse(r Reader, classes ClassMap, cpools PoolMap, class *ClassMetadata) error {
	return parseFields(r, classes, cpools, class, nil, true, ulf.parseField)
}

type VirtualizationInformation struct {
	EventBase
	StartTime int64
	Name      string
}

func (vi *VirtualizationInformation) parseField(name string, p ParseResolvable) (err error) {
	switch name {
	case "startTime":
		vi.StartTime, err = toLong(p)
	case "name":
		vi.Name, err = toString(p)
	}
	return err
}

func (vi *VirtualizationInformation) Parse(r Reader, classes ClassMap, cpools PoolMap, class *ClassMetadata) error {
	return parseFields(r, classes, cpools, class, nil, true, vi.parseField)
}

type YoungGenerationConfiguration struct {
	EventBase
	StartTime int64
	MinSize   int64
	MaxSize   int64
	NewRatio  int32
}

func (ygc *YoungGenerationConfiguration) parseField(name string, p ParseResolvable) (err error) {
	switch name {
	case "startTime":
		ygc.StartTime, err = toLong(p)
	case "minSize":
		ygc.MinSize, err = toLong(p)
	case "maxSize":
		ygc.MaxSize, err = toLong(p)
	case "newRatio":
		ygc.NewRatio, err = toInt(p)
	}
	return err
}

func (ygc *YoungGenerationConfiguration) Parse(r Reader, classes ClassMap, cpools PoolMap, class *ClassMetadata) error {
	return parseFields(r, classes, cpools, class, nil, true, ygc.parseField)
}

type UnsupportedEvent struct {
	EventBase
}

func (ue *UnsupportedEvent) parseField(name string, p ParseResolvable) error {
	return nil
}

func (ue *UnsupportedEvent) Parse(r Reader, classes ClassMap, cpools PoolMap, class *ClassMetadata) error {
	return parseFields(r, classes, cpools, class, nil, true, ue.parseField)
}
