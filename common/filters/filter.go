package filters

import (
	"github.com/grafana/jfr-parser/common"
	"github.com/grafana/jfr-parser/common/attributes"
	"github.com/grafana/jfr-parser/parser"
	"log/slog"
	"reflect"
)

const (
	Parked    = "PARKED"
	Runnable  = "RUNNABLE"
	Waiting   = "WAITING"
	Contended = "CONTENDED"
)

var (
	JavaMonitorInflate              = Types(common.JavaMonitorInflate)
	DatadogProfilerSetting          = Types(common.ProfilerSetting)
	DatadogEndpoint                 = Types(common.DatadogEndpoint)
	DATADOG_SCOPE                   = Types(common.SCOPE)
	DATADOG_EXCEPTION_COUNT         = Types(common.ExceptionCount)
	DATADOG_EXCEPTION_SAMPLE        = Types(common.ExceptionSample)
	DATADOG_EXECUTION_SAMPLE        = Types(common.DatadogExecutionSample)
	DatadogMethodSample             = Types(common.DatadogMethodSample)
	ThreadParked                    = AttributeEqual(attributes.ThreadStat, Parked)
	ThreadWaiting                   = AttributeEqual(attributes.ThreadStat, Waiting)
	DD_METHOD_SAMPLE_THREAD_PARKED  = AndFilters(DatadogMethodSample, ThreadParked)
	DD_METHOD_SAMPLE_THREAD_WAITING = AndFilters(DatadogMethodSample, ThreadWaiting)
	FILE_IO                         = Types(common.FileRead, common.FileWrite)
	SOCKET_IO                       = Types(common.SocketRead, common.SocketWrite)
	StacktraceNotNull               = NotNull(attributes.EventStacktrace)
	ALLOCATION                      = AndFilters(ALLOC_ALL, StacktraceNotNull)
	OBJ_ALLOCATION                  = AndFilters(ObjAlloc, StacktraceNotNull)
	MONITOR_ENTER                   = AndFilters(MonitorEnter, StacktraceNotNull)
	MONITOR_WAIT                    = Types(common.MonitorWait)
	THREAD_PARK                     = Types(common.ThreadPark)
	THREAD_SLEEP                    = Types(common.ThreadSleep)
	ASYNC_PROFILER_LOCK             = Types(common.MonitorEnter, common.ThreadPark)
	GARBAGE_COLLECTION              = AndFilters(GarbageCollection, StacktraceNotNull)
	SYNCHRONIZATION                 = Types(common.MonitorWait, common.ThreadPark)
	OLD_OBJECT_SAMPLE               = Types(common.OldObjectSample)
	DATADOG_PROFILER_CONFIG         = Types(common.DatadogProfilerConfig)
	DATADOG_ALLOCATION_SAMPLE       = Types(common.AllocationSample)
	DATADOG_HEAP_LIVE_OBJECT        = Types(common.HeapLiveObject)
)

var (
	SocketRead           = Types(common.SocketRead)
	SocketWrite          = Types(common.SocketWrite)
	SOCKET_READ_OR_WRITE = OrFilters(SocketRead, SocketWrite)
	NO_RMI_SOCKET_READ   = AndFilters(SocketRead, NotFilter(MethodFilter("sun.rmi.transport.tcp.TCPTransport", "handleMessages")),
		NotFilter(MethodFilter("javax.management.remote.rmi.RMIConnector$RMINotifClient", "fetchNotifs")))
	NO_RMI_SOCKET_WRITE = AndFilters(SocketWrite,
		NotFilter(MethodFilter("sun.rmi.transport.tcp.TCPTransport$ConnectionHandler", "run")),
		NotFilter(MethodFilter("sun.rmi.transport.tcp.TCPTransport$ConnectionHandler", "run0")))
	ENVIRONMENT_VARIABLE     = Types(common.EnvironmentVariable)
	FILE_READ                = Types(common.FileRead)
	FILE_WRITE               = Types(common.FileWrite)
	CodeCacheFull            = Types(common.CodeCacheFull)
	CodeCacheStatistics      = Types(common.CodeCacheStatistics)
	CodeCacheConfig          = Types(common.CodeCacheConfig)
	SweepCodeCache           = Types(common.SweepCodeCache)
	CODE_CACHE               = OrFilters(CodeCacheFull, CodeCacheStatistics, SweepCodeCache, CodeCacheConfig)
	CPU_INFORMATION          = Types(common.CPUInformation)
	GC_CONFIG                = Types(common.GcConf)
	HEAP_CONFIG              = Types(common.HeapConf)
	BeforeGc                 = AttributeEqual(attributes.GcWhen, "Before GC") //$NON-NLS-1$
	AfterGc                  = AttributeEqual(attributes.GcWhen, "After GC")  //$NON-NLS-1$
	ALLOC_OUTSIDE_TLAB       = Types(common.AllocOutsideTlab)
	ALLOC_INSIDE_TLAB        = Types(common.AllocInsideTlab)
	ALLOC_ALL                = Types(common.AllocInsideTlab, common.AllocOutsideTlab)
	ObjAlloc                 = Types(common.ObjAllocSample)
	REFERENCE_STATISTICS     = Types(common.GcReferenceStatistics)
	GarbageCollection        = Types(common.GarbageCollection)
	OLD_GARBAGE_COLLECTION   = Types(common.GcCollectorOldGarbageCollection)
	YOUNG_GARBAGE_COLLECTION = Types(common.GcCollectorYoungGarbageCollection)
	CONCURRENT_MODE_FAILURE  = Types(common.ConcurrentModeFailure)
	ERRORS                   = Types(common.ErrorsThrown)
	EXCEPTIONS               = Types(common.ExceptionsThrown)
	THROWABLES               = OrFilters(EXCEPTIONS, ERRORS)
	THROWABLES_STATISTICS    = Types(common.ThrowableStatistics)
	ClassUnload              = Types(common.ClassUnload)
	CLASS_LOAD_STATISTICS    = Types(common.ClassLoadStatistics)
	ClassLoaderStatistics    = Types(common.ClassLoaderStatistics)
	ClassLoad                = Types(common.ClassLoad)
	CLASS_LOAD_OR_UNLOAD     = OrFilters(ClassLoad, ClassUnload)
	ClassDefine              = Types(common.ClassDefine)
	CLASS_LOADER_EVENTS      = OrFilters(ClassLoad, ClassUnload, ClassDefine, ClassLoaderStatistics)
	MonitorEnter             = Types(common.MonitorEnter)
	FILE_OR_SOCKET_IO        = Types(common.SocketRead, common.SocketWrite, common.FileRead, common.FileWrite) // NOTE: Are there more types to add (i.e. relevant types with duration)?
	THREAD_LATENCIES         = Types(common.MonitorEnter,
		common.MonitorWait, common.ThreadSleep, common.ThreadPark, common.SocketRead,
		common.SocketWrite, common.FileRead, common.FileWrite, common.ClassLoad,
		common.COMPILATION, common.ExecutionSamplingInfoEventId)
	FilterExecutionSample      = Types(common.ExecutionSample)
	DatadogExecutionSample     = Types(common.DatadogExecutionSample)
	CONTEXT_SWITCH_RATE        = Types(common.ContextSwitchRate)
	CPU_LOAD                   = Types(common.CpuLoad)
	GcPause                    = Types(common.GcPause)
	GC_PAUSE_PHASE             = Types(common.GcPauseL1, common.GcPauseL2, common.GcPauseL3, common.GcPauseL4)
	TIME_CONVERSION            = Types(common.TimeConversion)
	VM_INFO                    = Types(common.VmInfo)
	THREAD_DUMP                = Types(common.ThreadDump)
	SYSTEM_PROPERTIES          = Types(common.SystemProperties)
	JFR_DATA_LOST              = Types(common.JfrDataLost)
	PROCESSES                  = Types(common.Processes)
	OBJECT_COUNT               = Types(common.ObjectCount)
	METASPACE_OOM              = Types(common.MetaspaceOom)
	COMPILATION                = Types(common.COMPILATION)
	COMPILER_FAILURE           = Types(common.CompilerFailure)
	COMPILER_STATS             = Types(common.CompilerStats)
	OS_MEMORY_SUMMARY          = Types(common.OSMemorySummary)
	HeapSummary                = Types(common.HeapSummary)
	HEAP_SUMMARY_BEFORE_GC     = AndFilters(HeapSummary, BeforeGc)
	HEAP_SUMMARY_AFTER_GC      = AndFilters(HeapSummary, AfterGc)
	MetaspaceSummary           = Types(common.MetaspaceSummary)
	METASPACE_SUMMARY_AFTER_GC = AndFilters(MetaspaceSummary, AfterGc)
	RECORDINGS                 = Types(common.RECORDINGS)
	RECORDING_SETTING          = Types(common.RecordingSetting)
	SafePoints                 = Types(common.SafepointBegin,
		common.SafepointCleanup, common.SafepointCleanupTask, common.SafepointStateSync,
		common.SafepointWaitBlocked, common.SafepointEnd)
	VM_OPERATIONS                       = Types(common.VmOperations) // NOTE: Not sure if there are any VM events that are neither blocking nor safepoint, but just in case.
	VM_OPERATIONS_BLOCKING_OR_SAFEPOINT = AndFilters(
		Types(common.VmOperations), OrFilters(AttributeEqual(attributes.Blocking, true), AttributeEqual(attributes.Safepoint, true)))
	// VmOperationsSafepoint NOTE: Are there any VM operations that are blocking, but not safepoints. Should we include those in the VM Thread??
	VmOperationsSafepoint      = AndFilters(Types(common.VmOperations), AttributeEqual(attributes.Safepoint, true))
	APPLICATION_PAUSES         = OrFilters(GcPause, SafePoints, VmOperationsSafepoint)
	BIASED_LOCKING_REVOCATIONS = Types(common.BiasedLockClassRevocation, common.BiasedLockRevocation, common.BiasedLockSelfRevocation)
	THREAD_CPU_LOAD            = Types(common.ThreadCpuLoad)
	NATIVE_METHOD_SAMPLE       = Types(common.NativeMethodSample)
	THREAD_START               = Types(common.JavaThreadStart)
	THREAD_END                 = Types(common.JavaThreadEnd)
)

type EventFilterFunc func(metadata *parser.ClassMetadata) parser.Predicate[parser.Event]

func (e EventFilterFunc) GetPredicate(metadata *parser.ClassMetadata) parser.Predicate[parser.Event] {
	return e(metadata)
}

type AndPredicate[T any] []parser.PredicateFunc[T]

func (a AndPredicate[T]) Test(t T) bool {
	for _, fn := range a {
		if !fn(t) {
			return false
		}
	}
	return true
}

type OrPredicate[T any] []parser.PredicateFunc[T]

func (o OrPredicate[T]) Test(t T) bool {
	for _, fn := range o {
		if fn(t) {
			return true
		}
	}
	return false
}

type NotPredicate[T any] parser.PredicateFunc[T]

func (n NotPredicate[T]) Test(t T) bool {
	return !n(t)
}

func And[T any](p ...parser.Predicate[T]) parser.Predicate[T] {
	ap := make(AndPredicate[T], 0, len(p))
	for _, pp := range p {
		ap = append(ap, pp.Test)
	}
	return ap
}

func Or[T any](p ...parser.Predicate[T]) parser.Predicate[T] {
	op := make(OrPredicate[T], 0, len(p))
	for _, pp := range p {
		op = append(op, pp.Test)
	}
	return op
}

func Not[T any](p parser.Predicate[T]) parser.Predicate[T] {
	return NotPredicate[T](p.Test)
}

func AndAlways(p ...parser.Predicate[parser.Event]) parser.Predicate[parser.Event] {
	switch len(p) {
	case 0:
		return parser.AlwaysTrue
	case 1:
		return p[0]
	}

	notAlwaysPred := make([]parser.Predicate[parser.Event], 0)
	for _, pp := range p {
		if parser.IsAlwaysFalse(pp) {
			return parser.AlwaysFalse
		}
		if parser.IsAlwaysTrue(pp) {
			continue
		}
		notAlwaysPred = append(notAlwaysPred, pp)
	}
	if len(notAlwaysPred) == 0 {
		return parser.AlwaysTrue
	}
	return And(notAlwaysPred...)
}

func OrAlways(p ...parser.Predicate[parser.Event]) parser.Predicate[parser.Event] {
	switch len(p) {
	case 0:
		return parser.AlwaysFalse
	case 1:
		return p[0]
	}
	notAlwaysPred := make([]parser.Predicate[parser.Event], 0)
	for _, pp := range p {
		if parser.IsAlwaysTrue(pp) {
			return parser.AlwaysTrue
		}
		if parser.IsAlwaysFalse(pp) {
			continue
		}
		notAlwaysPred = append(notAlwaysPred, pp)
	}
	if len(notAlwaysPred) == 0 {
		return parser.AlwaysFalse
	}
	return Or(notAlwaysPred...)
}

func NotAlways(p parser.Predicate[parser.Event]) parser.Predicate[parser.Event] {
	switch {
	case parser.IsAlwaysTrue(p):
		return parser.AlwaysFalse
	case parser.IsAlwaysFalse(p):
		return parser.AlwaysTrue
	}
	return Not(p)
}

func Types(classNames ...string) parser.EventFilter {
	if len(classNames) == 1 {
		return EventFilterFunc(func(metadata *parser.ClassMetadata) parser.Predicate[parser.Event] {
			if classNames[0] == metadata.Name {
				return parser.AlwaysTrue
			}
			return parser.AlwaysFalse
		})
	} else {
		et := make(map[string]struct{}, len(classNames))
		for _, className := range classNames {
			et[className] = struct{}{}
		}
		return EventFilterFunc(func(metadata *parser.ClassMetadata) parser.Predicate[parser.Event] {
			if _, ok := et[metadata.Name]; ok {
				return parser.AlwaysTrue
			}
			return parser.AlwaysFalse
		})
	}
}

func AttributeEqual[T comparable](attr *attributes.Attribute[T], target T) parser.EventFilter {
	return EventFilterFunc(func(metadata *parser.ClassMetadata) parser.Predicate[parser.Event] {

		if parser.IsAlwaysFalse(HasAttribute(attr).GetPredicate(metadata)) {
			return parser.AlwaysFalse
		}

		return parser.PredicateFunc[parser.Event](func(e parser.Event) bool {
			value, err := attr.GetValue(e.(*parser.GenericEvent))
			if err != nil {
				slog.Warn("unable to get attribute", "attribute", attr.Name)
				return false
			}
			return value == target
		})
	})
}

func HasAttribute[T any](attr *attributes.Attribute[T]) parser.EventFilter {
	return EventFilterFunc(func(metadata *parser.ClassMetadata) parser.Predicate[parser.Event] {
		f := metadata.GetField(attr.Name)
		if f == nil {
			return parser.AlwaysFalse
		}

		if metadata.ClassMap[f.ClassID].Name != attr.ClassName {
			return parser.AlwaysFalse
		}

		return parser.AlwaysTrue
	})
}

func MethodFilter(typeName, method string) parser.EventFilter {
	return EventFilterFunc(func(metadata *parser.ClassMetadata) parser.Predicate[parser.Event] {
		methodFilter := EventFilterFunc(func(metadata *parser.ClassMetadata) parser.Predicate[parser.Event] {

			return parser.PredicateFunc[parser.Event](func(e parser.Event) bool {
				stacktrace, err := attributes.EventStacktrace.GetValue(e.(*parser.GenericEvent))
				if err != nil {
					return false
				}
				for _, frame := range stacktrace.Frames {
					// todo check type full name
					if frame.Method.Type.Name.String == typeName && frame.Method.Name.String == method {
						return true
					}
				}
				return false
			})
		})

		return AndFilters(HasAttribute[*parser.StackTrace](attributes.EventStacktrace), methodFilter).GetPredicate(metadata)
	})
}

func AndFilters(filters ...parser.EventFilter) parser.EventFilter {
	return EventFilterFunc(func(metadata *parser.ClassMetadata) parser.Predicate[parser.Event] {
		predicates := make([]parser.Predicate[parser.Event], 0, len(filters))

		for _, filter := range filters {
			predicates = append(predicates, filter.GetPredicate(metadata))
		}

		return AndAlways(predicates...)
	})
}

func OrFilters(filters ...parser.EventFilter) parser.EventFilter {
	return EventFilterFunc(func(metadata *parser.ClassMetadata) parser.Predicate[parser.Event] {

		predicates := make([]parser.Predicate[parser.Event], 0, len(filters))

		for _, filter := range filters {
			predicates = append(predicates, filter.GetPredicate(metadata))
		}

		return OrAlways(predicates...)
	})
}

func NotFilter(filter parser.EventFilter) parser.EventFilter {
	return EventFilterFunc(func(metadata *parser.ClassMetadata) parser.Predicate[parser.Event] {
		return NotAlways(filter.GetPredicate(metadata))
	})
}

func NotNull[T any](attr *attributes.Attribute[T]) parser.EventFilter {
	return EventFilterFunc(func(metadata *parser.ClassMetadata) parser.Predicate[parser.Event] {
		if parser.IsAlwaysFalse(HasAttribute[T](attr).GetPredicate(metadata)) {
			return parser.AlwaysFalse
		}

		return parser.PredicateFunc[parser.Event](func(e parser.Event) bool {
			value, err := attr.GetValue(e.(*parser.GenericEvent))
			if err != nil {
				return false
			}
			rv := reflect.ValueOf(value)
			switch rv.Kind() {
			case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Pointer, reflect.Slice:
				return rv.IsNil()
			default:
				return false
			}
		})
	})
}
