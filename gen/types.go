package main

import (
	"fmt"

	"github.com/grafana/jfr-parser/parser/types/def"
)

var (
	T_METADATA                     = def.TypeID(0)
	T_CPOOL                        = def.TypeID(1)
	T_BOOLEAN                      = def.TypeID(4)
	T_CHAR                         = def.TypeID(5)
	T_FLOAT                        = def.TypeID(6)
	T_DOUBLE                       = def.TypeID(7)
	T_BYTE                         = def.TypeID(8)
	T_SHORT                        = def.TypeID(9)
	T_INT                          = def.TypeID(10)
	T_LONG                         = def.TypeID(11)
	T_STRING                       = def.TypeID(20)
	T_CLASS                        = def.TypeID(21)
	T_THREAD                       = def.TypeID(22)
	T_CLASS_LOADER                 = def.TypeID(23)
	T_FRAME_TYPE                   = def.TypeID(24)
	T_THREAD_STATE                 = def.TypeID(25)
	T_STACK_TRACE                  = def.TypeID(26)
	T_STACK_FRAME                  = def.TypeID(27)
	T_METHOD                       = def.TypeID(28)
	T_PACKAGE                      = def.TypeID(29)
	T_SYMBOL                       = def.TypeID(30)
	T_LOG_LEVEL                    = def.TypeID(31)
	T_ExecutionMode                = def.TypeID(33)
	T_CounterName                  = def.TypeID(34)
	T_EVENT                        = def.TypeID(100)
	T_EXECUTION_SAMPLE             = def.TypeID(101)
	T_ALLOC_IN_NEW_TLAB            = def.TypeID(102)
	T_ALLOC_OUTSIDE_TLAB           = def.TypeID(103)
	T_MONITOR_ENTER                = def.TypeID(104)
	T_THREAD_PARK                  = def.TypeID(105)
	T_CPU_LOAD                     = def.TypeID(106)
	T_ACTIVE_RECORDING             = def.TypeID(107)
	T_ACTIVE_SETTING               = def.TypeID(108)
	T_OS_INFORMATION               = def.TypeID(109)
	T_CPU_INFORMATION              = def.TypeID(110)
	T_JVM_INFORMATION              = def.TypeID(111)
	T_INITIAL_SYSTEM_PROPERTY      = def.TypeID(112)
	T_NATIVE_LIBRARY               = def.TypeID(113)
	T_LOG                          = def.TypeID(114)
	T_LIVE_OBJECT                  = def.TypeID(115)
	T_WALL_CLOCK_SAMPLE            = def.TypeID(118)
	T_MALLOC                       = def.TypeID(119)
	T_FREE                         = def.TypeID(120)
	T_HeapUsage                    = def.TypeID(121)
	T_QueueTime                    = def.TypeID(123)
	T_DatadogProfilerConfig        = def.TypeID(122)
	T_DatadogProfilerClassRefCache = def.TypeID(124)
	T_ProfilerCounter              = def.TypeID(125)
	T_ANNOTATION                   = def.TypeID(200)
	T_LABEL                        = def.TypeID(201)
	T_CATEGORY                     = def.TypeID(202)
	T_TIMESTAMP                    = def.TypeID(203)
	T_TIMESPAN                     = def.TypeID(204)
	T_DATA_AMOUNT                  = def.TypeID(205)
	T_MEMORY_ADDRESS               = def.TypeID(206)
	T_UNSIGNED                     = def.TypeID(207)
	T_PERCENTAGE                   = def.TypeID(208)
	T_ALLOC_SAMPLE                 = def.TypeID(209)

	T_ObjectSample                 = def.TypeID(103)
	T_Datadog_MethodSample         = def.TypeID(102)
	T_Datadog_Excepion_Sample      = def.TypeID(5818)
	T_GCName                       = def.TypeID(169)
	T_GCCause                      = def.TypeID(170)
	T_G1YCType                     = def.TypeID(173)
	T_WallClockSamplingEpoch       = def.TypeID(118)
	T_GarbageCollection            = def.TypeID(35)
	T_SystemGC                     = def.TypeID(36)
	T_ParallelOldGarbageCollection = def.TypeID(37)
	T_YoungGarbageCollection       = def.TypeID(38)
	T_OldGarbageCollection         = def.TypeID(39)
	T_G1GarbageCollection          = def.TypeID(40)
)

func TypeID2Sym(id def.TypeID) string {
	switch id {
	case T_METADATA:
		return "T_METADATA"
	case T_CPOOL:
		return "T_CPOOL"
	case T_BOOLEAN:
		return "T_BOOLEAN"
	case T_CHAR:
		return "T_CHAR"
	case T_FLOAT:
		return "T_FLOAT"
	case T_DOUBLE:
		return "T_DOUBLE"
	case T_BYTE:
		return "T_BYTE"
	case T_SHORT:
		return "T_SHORT"
	case T_INT:
		return "T_INT"
	case T_LONG:
		return "T_LONG"
	case T_STRING:
		return "T_STRING"
	case T_CLASS:
		return "T_CLASS"
	case T_THREAD:
		return "T_THREAD"
	case T_CLASS_LOADER:
		return "T_CLASS_LOADER"
	case T_FRAME_TYPE:
		return "T_FRAME_TYPE"
	case T_THREAD_STATE:
		return "T_THREAD_STATE"
	case T_STACK_TRACE:
		return "T_STACK_TRACE"
	case T_STACK_FRAME:
		return "T_STACK_FRAME"
	case T_METHOD:
		return "T_METHOD"
	case T_PACKAGE:
		return "T_PACKAGE"
	case T_SYMBOL:
		return "T_SYMBOL"
	case T_LOG_LEVEL:
		return "T_LOG_LEVEL"
	case T_EVENT:
		return "T_EVENT"
	case T_EXECUTION_SAMPLE:
		return "T_EXECUTION_SAMPLE"
	case T_ALLOC_IN_NEW_TLAB:
		return "T_ALLOC_IN_NEW_TLAB"
	case T_ALLOC_OUTSIDE_TLAB:
		return "T_ALLOC_OUTSIDE_TLAB"
	case T_ALLOC_SAMPLE:
		return "T_ALLOC_SAMPLE"
	case T_MONITOR_ENTER:
		return "T_MONITOR_ENTER"
	case T_THREAD_PARK:
		return "T_THREAD_PARK"
	case T_CPU_LOAD:
		return "T_CPU_LOAD"
	case T_ACTIVE_RECORDING:
		return "T_ACTIVE_RECORDING"
	case T_ACTIVE_SETTING:
		return "T_ACTIVE_SETTING"
	case T_OS_INFORMATION:
		return "T_OS_INFORMATION"
	case T_CPU_INFORMATION:
		return "T_CPU_INFORMATION"
	case T_JVM_INFORMATION:
		return "T_JVM_INFORMATION"
	case T_INITIAL_SYSTEM_PROPERTY:
		return "T_INITIAL_SYSTEM_PROPERTY"
	case T_NATIVE_LIBRARY:
		return "T_NATIVE_LIBRARY"
	case T_LOG:
		return "T_LOG"
	case T_LIVE_OBJECT:
		return "T_LIVE_OBJECT"
	case T_ANNOTATION:
		return "T_ANNOTATION"
	case T_LABEL:
		return "T_LABEL"
	case T_CATEGORY:
		return "T_CATEGORY"
	case T_TIMESTAMP:
		return "T_TIMESTAMP"
	case T_TIMESPAN:
		return "T_TIMESPAN"
	case T_DATA_AMOUNT:
		return "T_DATA_AMOUNT"
	case T_MEMORY_ADDRESS:
		return "T_MEMORY_ADDRESS"
	case T_UNSIGNED:
		return "T_UNSIGNED"
	case T_PERCENTAGE:
		return "T_PERCENTAGE"
	/*	T_HeapUsage                    = def.TypeID(121)
		T_QueueTime                    = def.TypeID(123)
		T_DatadogProfilerConfig        = def.TypeID(122)
		T_DatadogProfilerClassRefCache = def.TypeID(124)
		T_ProfilerCounter              = def.TypeID(125)
		T_ExecutionMode                = def.TypeID(33)
		T_CountName                    = def.TypeID(34)*/
	case T_ExecutionMode:
		return "T_ExecutionMode"
	case T_CounterName:
		return "T_CounterName"
	case T_HeapUsage:
		return "T_HeapUsage"
	case T_QueueTime:
		return "T_QueueTime"
	case T_DatadogProfilerConfig:
		return "T_DatadogProfilerConfig"
	case T_DatadogProfilerClassRefCache:
		return "T_DatadogProfilerClassRefCache"
	case T_ProfilerCounter:
		return "T_ProfilerCounter"
		/*
			T_ObjectSample                 = def.TypeID(103)
			T_Datadog_MethodSample         = def.TypeID(102)
			T_Datadog_Excepion_Sample      = def.TypeID(5818)
			T_GCName                       = def.TypeID(169)
			T_GCCause                      = def.TypeID(170)
			T_G1YCType                     = def.TypeID(173)
			T_WallClockSamplingEpoch       = def.TypeID(118)
			T_GarbageCollection            = def.TypeID(35)
			T_SystemGC                     = def.TypeID(36)
			T_ParallelOldGarbageCollection = def.TypeID(37)
			T_YoungGarbageCollection       = def.TypeID(38)
			T_OldGarbageCollection         = def.TypeID(39)
			T_G1GarbageCollection          = def.TypeID(40)

		*/
	case T_ObjectSample:
		return "T_ObjectSample"
	case T_Datadog_MethodSample:
		return "T_Datadog_MethodSample"
	case T_Datadog_Excepion_Sample:
		return "T_Datadog_Excepion_Sample"
	case T_GCName:
		return "T_GCName"
	case T_GCCause:
		return "T_GCCause"
	case T_G1YCType:
		return "T_G1YCType"
	case T_WallClockSamplingEpoch:
		return "T_WallClockSamplingEpoch"
	case T_GarbageCollection:
		return "T_GarbageCollection"
	case T_SystemGC:
		return "T_SystemGC"
	case T_ParallelOldGarbageCollection:
		return "T_ParallelOldGarbageCollection"
	case T_YoungGarbageCollection:
		return "T_YoungGarbageCollection"
	case T_OldGarbageCollection:
		return "T_OldGarbageCollection"
	case T_G1GarbageCollection:
		return "T_G1GarbageCollection"
	default:
		return fmt.Sprintf("unknown type %d", id)
	}
}

var Type_boolean = def.Class{
	Name:   "boolean",
	ID:     T_BOOLEAN,
	Fields: []def.Field{},
}
var Type_char = def.Class{
	Name:   "char",
	ID:     T_CHAR,
	Fields: []def.Field{},
}
var Type_float = def.Class{
	Name:   "float",
	ID:     T_FLOAT,
	Fields: []def.Field{},
}
var Type_double = def.Class{
	Name:   "double",
	ID:     T_DOUBLE,
	Fields: []def.Field{},
}
var Type_byte = def.Class{
	Name:   "byte",
	ID:     T_BYTE,
	Fields: []def.Field{},
}
var Type_short = def.Class{
	Name:   "short",
	ID:     T_SHORT,
	Fields: []def.Field{},
}
var Type_int = def.Class{
	Name:   "int",
	ID:     T_INT,
	Fields: []def.Field{},
}
var Type_long = def.Class{
	Name:   "long",
	ID:     T_LONG,
	Fields: []def.Field{},
}
var Type_java_lang_String = def.Class{
	Name:   "java.lang.String",
	ID:     T_STRING,
	Fields: []def.Field{},
}
var Type_java_lang_Class = def.Class{
	Name: "java.lang.Class",
	ID:   T_CLASS,
	Fields: []def.Field{
		{Name: "classLoader", Type: T_CLASS_LOADER, ConstantPool: true},
		{Name: "name", Type: T_SYMBOL, ConstantPool: true},
		{Name: "package", Type: T_PACKAGE, ConstantPool: true},
		{Name: "modifiers", Type: T_INT, ConstantPool: false},
	},
}
var Type_java_lang_Thread = def.Class{
	Name: "java.lang.Thread",
	ID:   T_THREAD,
	Fields: []def.Field{
		{Name: "osName", Type: T_STRING, ConstantPool: false},
		{Name: "osThreadId", Type: T_LONG, ConstantPool: false},
		{Name: "javaName", Type: T_STRING, ConstantPool: false},
		{Name: "javaThreadId", Type: T_LONG, ConstantPool: false},
	},
}
var Type_jdk_types_ClassLoader = def.Class{
	Name: "jdk.types.ClassLoader",
	ID:   T_CLASS_LOADER,
	Fields: []def.Field{
		{Name: "type", Type: T_CLASS, ConstantPool: true},
		{Name: "name", Type: T_SYMBOL, ConstantPool: true},
	},
}
var Type_jdk_types_FrameType = def.Class{
	Name: "jdk.types.FrameType",
	ID:   T_FRAME_TYPE,
	Fields: []def.Field{
		{Name: "description", Type: T_STRING, ConstantPool: false},
	},
}
var Type_jdk_types_ThreadState = def.Class{
	Name: "jdk.types.ThreadState",
	ID:   T_THREAD_STATE,
	Fields: []def.Field{
		{Name: "name", Type: T_STRING, ConstantPool: false},
	},
}
var Type_jdk_types_StackTrace = def.Class{
	Name: "jdk.types.StackTrace",
	ID:   T_STACK_TRACE,
	Fields: []def.Field{
		{Name: "truncated", Type: T_BOOLEAN, ConstantPool: false},
		{Name: "frames", Type: T_STACK_FRAME, ConstantPool: false, Array: true},
	},
}
var Type_jdk_types_StackFrame = def.Class{
	Name: "jdk.types.StackFrame",
	ID:   T_STACK_FRAME,
	Fields: []def.Field{
		{Name: "method", Type: T_METHOD, ConstantPool: true},
		{Name: "lineNumber", Type: T_INT, ConstantPool: false},
		{Name: "bytecodeIndex", Type: T_INT, ConstantPool: false},
		{Name: "type", Type: T_FRAME_TYPE, ConstantPool: true},
	},
}
var Type_jdk_types_Method = def.Class{
	Name: "jdk.types.Method",
	ID:   T_METHOD,
	Fields: []def.Field{
		{Name: "type", Type: T_CLASS, ConstantPool: true},
		{Name: "name", Type: T_SYMBOL, ConstantPool: true},
		{Name: "descriptor", Type: T_SYMBOL, ConstantPool: true},
		{Name: "modifiers", Type: T_INT, ConstantPool: false},
		{Name: "hidden", Type: T_BOOLEAN, ConstantPool: false},
	},
}
var Type_jdk_types_Package = def.Class{
	Name: "jdk.types.Package",
	ID:   T_PACKAGE,
	Fields: []def.Field{
		{Name: "name", Type: T_SYMBOL, ConstantPool: true},
	},
}
var Type_jdk_types_Symbol = def.Class{
	Name: "jdk.types.Symbol",
	ID:   T_SYMBOL,
	Fields: []def.Field{
		{Name: "string", Type: T_STRING, ConstantPool: false},
	},
}
var Type_profiler_types_LogLevel = def.Class{
	Name: "profiler.types.LogLevel",
	ID:   T_LOG_LEVEL,
	Fields: []def.Field{
		{Name: "name", Type: T_STRING, ConstantPool: false},
	},
}
var Type_jdk_ExecutionSample = def.Class{
	Name: "jdk.ExecutionSample",
	ID:   T_EXECUTION_SAMPLE,
	Fields: []def.Field{
		{Name: "startTime", Type: T_LONG, ConstantPool: false},
		{Name: "sampledThread", Type: T_THREAD, ConstantPool: true},
		{Name: "stackTrace", Type: T_STACK_TRACE, ConstantPool: true},
		{Name: "state", Type: T_THREAD_STATE, ConstantPool: true},
		{Name: "spanId", Type: T_LONG, ConstantPool: false},
		{Name: "spanName", Type: T_LONG, ConstantPool: false},
		{Name: "contextId", Type: T_LONG, ConstantPool: false},
	},
}

var Type_profiler_WallClockSample = def.Class{
	Name: "profiler.WallClockSample",
	ID:   T_WALL_CLOCK_SAMPLE,
	Fields: []def.Field{
		{Name: "startTime", Type: T_LONG, ConstantPool: false},
		{Name: "sampledThread", Type: T_THREAD, ConstantPool: true},
		{Name: "stackTrace", Type: T_STACK_TRACE, ConstantPool: true},
		{Name: "state", Type: T_THREAD_STATE, ConstantPool: true},
		{Name: "samples", Type: T_INT, ConstantPool: false},
	},
}

var Type_jdk_ObjectAllocationInNewTLAB = def.Class{
	Name: "jdk.ObjectAllocationInNewTLAB",
	ID:   T_ALLOC_IN_NEW_TLAB,
	Fields: []def.Field{
		{Name: "startTime", Type: T_LONG, ConstantPool: false},
		{Name: "eventThread", Type: T_THREAD, ConstantPool: true},
		{Name: "stackTrace", Type: T_STACK_TRACE, ConstantPool: true},
		{Name: "objectClass", Type: T_CLASS, ConstantPool: true},
		{Name: "allocationSize", Type: T_LONG, ConstantPool: false},
		{Name: "tlabSize", Type: T_LONG, ConstantPool: false},
		{Name: "contextId", Type: T_LONG, ConstantPool: false},
		{Name: "spanId", Type: T_LONG, ConstantPool: false},
		{Name: "spanName", Type: T_LONG, ConstantPool: false},
	},
}
var Type_jdk_ObjectAllocationOutsideTLAB = def.Class{
	Name: "jdk.ObjectAllocationOutsideTLAB",
	ID:   T_ALLOC_OUTSIDE_TLAB,
	Fields: []def.Field{
		{Name: "startTime", Type: T_LONG, ConstantPool: false},
		{Name: "eventThread", Type: T_THREAD, ConstantPool: true},
		{Name: "stackTrace", Type: T_STACK_TRACE, ConstantPool: true},
		{Name: "objectClass", Type: T_CLASS, ConstantPool: true},
		{Name: "allocationSize", Type: T_LONG, ConstantPool: false},
		{Name: "contextId", Type: T_LONG, ConstantPool: false},
		{Name: "spanId", Type: T_LONG, ConstantPool: false},
		{Name: "spanName", Type: T_LONG, ConstantPool: false},
	},
}
var Type_jdk_ObjectAllocationSample = def.Class{
	Name: "jdk.ObjectAllocationSample",
	ID:   T_ALLOC_SAMPLE,
	Fields: []def.Field{
		{Name: "startTime", Type: T_LONG, ConstantPool: false},
		{Name: "eventThread", Type: T_THREAD, ConstantPool: true},
		{Name: "stackTrace", Type: T_STACK_TRACE, ConstantPool: true},
		{Name: "objectClass", Type: T_CLASS, ConstantPool: true},
		{Name: "weight", Type: T_LONG, ConstantPool: false},
	},
}

var Type_jdk_JavaMonitorEnter = def.Class{
	Name: "jdk.JavaMonitorEnter",
	ID:   T_MONITOR_ENTER,
	Fields: []def.Field{
		{Name: "startTime", Type: T_LONG, ConstantPool: false},
		{Name: "duration", Type: T_LONG, ConstantPool: false},
		{Name: "eventThread", Type: T_THREAD, ConstantPool: true},
		{Name: "stackTrace", Type: T_STACK_TRACE, ConstantPool: true},
		{Name: "monitorClass", Type: T_CLASS, ConstantPool: true},
		{Name: "previousOwner", Type: T_THREAD, ConstantPool: true},
		{Name: "address", Type: T_LONG, ConstantPool: false},
		{Name: "contextId", Type: T_LONG, ConstantPool: false},
		{Name: "spanId", Type: T_LONG, ConstantPool: false},
		{Name: "spanName", Type: T_LONG, ConstantPool: false},
	},
}

var Type_jdk_ThreadPark = def.Class{
	Name: "jdk.ThreadPark",
	ID:   T_THREAD_PARK,
	Fields: []def.Field{
		{Name: "startTime", Type: T_LONG, ConstantPool: false},
		{Name: "duration", Type: T_LONG, ConstantPool: false},
		{Name: "eventThread", Type: T_THREAD, ConstantPool: true},
		{Name: "stackTrace", Type: T_STACK_TRACE, ConstantPool: true},
		{Name: "parkedClass", Type: T_CLASS, ConstantPool: true},
		{Name: "timeout", Type: T_LONG, ConstantPool: false},
		{Name: "until", Type: T_LONG, ConstantPool: false},
		{Name: "address", Type: T_LONG, ConstantPool: false},
	},
}
var Type_jdk_CPULoad = def.Class{
	Name: "jdk.CPULoad",
	ID:   T_CPU_LOAD,
	Fields: []def.Field{
		{Name: "startTime", Type: T_LONG, ConstantPool: false},
		{Name: "jvmUser", Type: T_FLOAT, ConstantPool: false},
		{Name: "jvmSystem", Type: T_FLOAT, ConstantPool: false},
		{Name: "machineTotal", Type: T_FLOAT, ConstantPool: false},
	},
}
var Type_jdk_ActiveRecording = def.Class{
	Name: "jdk.ActiveRecording",
	ID:   T_ACTIVE_RECORDING,
	Fields: []def.Field{
		{Name: "startTime", Type: T_LONG, ConstantPool: false},
		{Name: "duration", Type: T_LONG, ConstantPool: false},
		{Name: "eventThread", Type: T_THREAD, ConstantPool: true},
		{Name: "id", Type: T_LONG, ConstantPool: false},
		{Name: "name", Type: T_STRING, ConstantPool: false},
		{Name: "destination", Type: T_STRING, ConstantPool: false},
		{Name: "maxAge", Type: T_LONG, ConstantPool: false},
		{Name: "maxSize", Type: T_LONG, ConstantPool: false},
		{Name: "recordingStart", Type: T_LONG, ConstantPool: false},
		{Name: "recordingDuration", Type: T_LONG, ConstantPool: false},
	},
}
var Type_jdk_ActiveSetting = def.Class{
	Name: "jdk.ActiveSetting",
	ID:   T_ACTIVE_SETTING,
	Fields: []def.Field{
		{Name: "startTime", Type: T_LONG, ConstantPool: false},
		{Name: "duration", Type: T_LONG, ConstantPool: false},
		{Name: "eventThread", Type: T_THREAD, ConstantPool: true},
		{Name: "stackTrace", Type: T_STACK_TRACE, ConstantPool: true},
		{Name: "id", Type: T_LONG, ConstantPool: false},
		{Name: "name", Type: T_STRING, ConstantPool: false},
		{Name: "value", Type: T_STRING, ConstantPool: false},
	},
}
var Type_jdk_OSInformation = def.Class{
	Name: "jdk.OSInformation",
	ID:   T_OS_INFORMATION,
	Fields: []def.Field{
		{Name: "startTime", Type: T_LONG, ConstantPool: false},
		{Name: "osVersion", Type: T_STRING, ConstantPool: false},
	},
}
var Type_jdk_CPUInformation = def.Class{
	Name: "jdk.CPUInformation",
	ID:   T_CPU_INFORMATION,
	Fields: []def.Field{
		{Name: "startTime", Type: T_LONG, ConstantPool: false},
		{Name: "cpu", Type: T_STRING, ConstantPool: false},
		{Name: "description", Type: T_STRING, ConstantPool: false},
		{Name: "sockets", Type: T_INT, ConstantPool: false},
		{Name: "cores", Type: T_INT, ConstantPool: false},
		{Name: "hwThreads", Type: T_INT, ConstantPool: false},
	},
}
var Type_jdk_JVMInformation = def.Class{
	Name: "jdk.JVMInformation",
	ID:   T_JVM_INFORMATION,
	Fields: []def.Field{
		{Name: "startTime", Type: T_LONG, ConstantPool: false},
		{Name: "jvmName", Type: T_STRING, ConstantPool: false},
		{Name: "jvmVersion", Type: T_STRING, ConstantPool: false},
		{Name: "jvmArguments", Type: T_STRING, ConstantPool: false},
		{Name: "jvmFlags", Type: T_STRING, ConstantPool: false},
		{Name: "javaArguments", Type: T_STRING, ConstantPool: false},
		{Name: "jvmStartTime", Type: T_LONG, ConstantPool: false},
		{Name: "pid", Type: T_LONG, ConstantPool: false},
	},
}
var Type_jdk_InitialSystemProperty = def.Class{
	Name: "jdk.InitialSystemProperty",
	ID:   T_INITIAL_SYSTEM_PROPERTY,
	Fields: []def.Field{
		{Name: "startTime", Type: T_LONG, ConstantPool: false},
		{Name: "key", Type: T_STRING, ConstantPool: false},
		{Name: "value", Type: T_STRING, ConstantPool: false},
	},
}
var Type_jdk_NativeLibrary = def.Class{
	Name: "jdk.NativeLibrary",
	ID:   T_NATIVE_LIBRARY,
	Fields: []def.Field{
		{Name: "startTime", Type: T_LONG, ConstantPool: false},
		{Name: "name", Type: T_STRING, ConstantPool: false},
		{Name: "baseAddress", Type: T_LONG, ConstantPool: false},
		{Name: "topAddress", Type: T_LONG, ConstantPool: false},
	},
}
var Type_profiler_Log = def.Class{
	Name: "profiler.Log",
	ID:   T_LOG,
	Fields: []def.Field{
		{Name: "startTime", Type: T_LONG, ConstantPool: false},
		{Name: "level", Type: T_LOG_LEVEL, ConstantPool: true},
		{Name: "message", Type: T_STRING, ConstantPool: false},
	},
}

var Type_profiler_LiveObject = def.Class{
	Name: "profiler.LiveObject",
	ID:   T_LIVE_OBJECT,
	Fields: []def.Field{
		{Name: "startTime", Type: T_LONG, ConstantPool: false},
		{Name: "eventThread", Type: T_THREAD, ConstantPool: true},
		{Name: "stackTrace", Type: T_STACK_TRACE, ConstantPool: true},
		{Name: "objectClass", Type: T_CLASS, ConstantPool: true},
		{Name: "allocationSize", Type: T_LONG, ConstantPool: false},
		{Name: "allocationTime", Type: T_LONG, ConstantPool: false},
	},
}
var Type_jdk_jfr_Label = def.Class{
	Name: "jdk.jfr.Label",
	ID:   T_LABEL,
	Fields: []def.Field{
		{Name: "value", Type: T_STRING, ConstantPool: false},
	},
}
var Type_jdk_jfr_Category = def.Class{
	Name: "jdk.jfr.Category",
	ID:   T_CATEGORY,
	Fields: []def.Field{
		{Name: "value", Type: T_STRING, ConstantPool: false, Array: true},
	},
}
var Type_jdk_jfr_Timestamp = def.Class{
	Name: "jdk.jfr.Timestamp",
	ID:   T_TIMESTAMP,
	Fields: []def.Field{
		{Name: "value", Type: T_STRING, ConstantPool: false},
	},
}
var Type_jdk_jfr_Timespan = def.Class{
	Name: "jdk.jfr.Timespan",
	ID:   T_TIMESPAN,
	Fields: []def.Field{
		{Name: "value", Type: T_STRING, ConstantPool: false},
	},
}
var Type_jdk_jfr_DataAmount = def.Class{
	Name: "jdk.jfr.DataAmount",
	ID:   T_DATA_AMOUNT,
	Fields: []def.Field{
		{Name: "value", Type: T_STRING, ConstantPool: false},
	},
}
var Type_jdk_jfr_MemoryAddress = def.Class{
	Name:   "jdk.jfr.MemoryAddress",
	ID:     T_MEMORY_ADDRESS,
	Fields: []def.Field{},
}
var Type_jdk_jfr_Unsigned = def.Class{
	Name:   "jdk.jfr.Unsigned",
	ID:     T_UNSIGNED,
	Fields: []def.Field{},
}
var Type_jdk_jfr_Percentage = def.Class{
	Name:   "jdk.jfr.Percentage",
	ID:     T_PERCENTAGE,
	Fields: []def.Field{},
}

var Type_profiler_Malloc = def.Class{
	Name: "profiler.Malloc",
	ID:   T_MALLOC,
	Fields: []def.Field{
		{Name: "startTime", Type: T_LONG, ConstantPool: false},
		{Name: "eventThread", Type: T_THREAD, ConstantPool: true},
		{Name: "stackTrace", Type: T_STACK_TRACE, ConstantPool: true},
		{Name: "address", Type: T_LONG, ConstantPool: false},
		{Name: "size", Type: T_LONG, ConstantPool: false},
	},
}

var Type_profiler_Free = def.Class{
	Name: "profiler.Free",
	ID:   T_FREE,
	Fields: []def.Field{
		{Name: "startTime", Type: T_LONG, ConstantPool: false},
		{Name: "eventThread", Type: T_THREAD, ConstantPool: true},
		{Name: "stackTrace", Type: T_STACK_TRACE, ConstantPool: true},
		{Name: "address", Type: T_LONG, ConstantPool: false},
	},
}

var Type_ExecutionMode = def.Class{
	Name: "datadog.types.ExecutionMode",
	ID:   T_ExecutionMode,
	Fields: []def.Field{
		{Name: "name", Type: T_STRING, ConstantPool: false},
	},
}

//datadog.ProfilerCounter:class{name: datadog.ProfilerCounter, id: 125, fields: [
//{Name:startTime Type:11 ConstantPool:false Array:false}
//{Name:name Type:34 ConstantPool:false Array:false}
//{Name:count Type:11 ConstantPool:false Array:false}]}

var Type_ProfilerCounter = def.Class{
	Name: "datadog.ProfilerCounter",
	ID:   T_ProfilerCounter,
	Fields: []def.Field{
		{Name: "startTime", Type: T_LONG, ConstantPool: false},
		{Name: "name", Type: T_CounterName, ConstantPool: false},
		{Name: "count", Type: T_LONG, ConstantPool: false},
	},
}

//profiler.types.CounterName:class{name: profiler.types.CounterName, id: 34, fields: [
//{Name:value Type:20 ConstantPool:false Array:false}]}

var Type_CounterName = def.Class{
	Name: "profiler.types.CounterName",
	ID:   T_CounterName,
	Fields: []def.Field{
		{Name: "value", Type: T_STRING, ConstantPool: false},
	},
}

//class{name: datadog.ObjectSample, id: 103, fields: [{Name:startTime Type:11 ConstantPool:false Array:false}
//{Name:eventThread Type:22 ConstantPool:true Array:false} {Name:stackTrace Type:26 ConstantPool:true Array:false}
//{Name:objectClass Type:21 ConstantPool:true Array:false} {Name:size Type:11 ConstantPool:false Array:false} {Name:weight Type:6 ConstantPool:false Array:false}
//{Name:spanId Type:11 ConstantPool:false Array:false} {Name:localRootSpanId Type:11 ConstantPool:false Array:false} {Name:_dd.trace.operation Type:32 ConstantPool:true Array:false}]}

var Type_DatadogObjectSample = def.Class{
	Name: "datadog.ObjectSample",
	ID:   T_ObjectSample,
	Fields: []def.Field{
		{Name: "startTime", Type: T_LONG, ConstantPool: false},
		{Name: "eventThread", Type: T_THREAD, ConstantPool: true},
		{Name: "stackTrace", Type: T_STACK_TRACE, ConstantPool: true},
		{Name: "objectClass", Type: T_CLASS, ConstantPool: true},
		{Name: "size", Type: T_LONG, ConstantPool: false},
		{Name: "weight", Type: T_FLOAT, ConstantPool: false},
	},
}
var Type_datadogProfilerClassRefCache = def.Class{
	Name: "datadog.DatadogProfilerClassRefCache",
	ID:   def.TypeID(124),
	// {Name:startTime Type:11 ConstantPool:false Array:false} {Name:size Type:11 ConstantPool:false Array:false}]
	Fields: []def.Field{
		{Name: "startTime", Type: T_LONG, ConstantPool: false},
		{Name: "size", Type: T_LONG, ConstantPool: false, Array: false},
	},
}

//datadog.DatadogProfilerConfig:class{name: datadog.DatadogProfilerConfig, id: 122, fields: [
//{Name:startTime Type:11 ConstantPool:false Array:false}
//{Name:duration Type:11 ConstantPool:false Array:false}
//{Name:eventThread Type:22 ConstantPool:true Array:false}
//{Name:cpuInterval Type:11 ConstantPool:false Array:false}
//{Name:wallInterval Type:11 ConstantPool:false Array:false}
//{Name:allocInterval Type:11 ConstantPool:false Array:false}
//{Name:memleakInterval Type:11 ConstantPool:false Array:false}
//{Name:memleakCapacity Type:11 ConstantPool:false Array:false}
//{Name:memleakTrackPercent Type:11 ConstantPool:false Array:false}
//{Name:gcGenerations Type:4 ConstantPool:false Array:false}
//{Name:modeMask Type:10 ConstantPool:false Array:false}
//{Name:version Type:20 ConstantPool:false Array:false}
//{Name:cpuEngine Type:20 ConstantPool:false Array:false}]}

var Type_DatadogProfilerConfig = def.Class{
	Name: "datadog.DatadogProfilerConfig",
	ID:   T_DatadogProfilerConfig,
	Fields: []def.Field{
		{Name: "startTime", Type: T_LONG, ConstantPool: false},
		{Name: "duration", Type: T_LONG, ConstantPool: false},
		{Name: "eventThread", Type: T_THREAD, ConstantPool: true},
		{Name: "cpuInterval", Type: T_LONG, ConstantPool: false},
		{Name: "wallInterval", Type: T_LONG, ConstantPool: false},
		{Name: "allocInterval", Type: T_LONG, ConstantPool: false},
		{Name: "memleakInterval", Type: T_LONG, ConstantPool: false},
		{Name: "memleakCapacity", Type: T_LONG, ConstantPool: false},
		{Name: "memleakTrackPercent", Type: T_LONG, ConstantPool: false},
		{Name: "gcGenerations", Type: T_BOOLEAN, ConstantPool: false},
		{Name: "modeMask", Type: T_INT, ConstantPool: false},
		{Name: "version", Type: T_STRING, ConstantPool: false},
		{Name: "cpuEngine", Type: T_STRING, ConstantPool: false},
	},
}

//datadog.HeapLiveObject:class{name: datadog.HeapLiveObject, id: 105, fields: [
//{Name:startTime Type:11 ConstantPool:false Array:false}
//{Name:eventThread Type:22 ConstantPool:true Array:false}
//{Name:stackTrace Type:26 ConstantPool:true Array:false}
//{Name:objectClass Type:21 ConstantPool:true Array:false}
//{Name:age Type:11 ConstantPool:false Array:false}
//{Name:size Type:11 ConstantPool:false Array:false}
//{Name:weight Type:6 ConstantPool:false Array:false}
//{Name:spanId Type:11 ConstantPool:false Array:false}
//{Name:localRootSpanId Type:11 ConstantPool:false Array:false}
//{Name:_dd.trace.operation Type:32 ConstantPool:true Array:false}]}

var Type_Datadog_HeapLiveObject = def.Class{
	Name: "datadog.HeapLiveObject",
	ID:   def.TypeID(105),
	Fields: []def.Field{
		{Name: "startTime", Type: T_LONG, ConstantPool: false},
		{Name: "eventThread", Type: T_THREAD, ConstantPool: true},
		{Name: "stackTrace", Type: T_STACK_TRACE, ConstantPool: true},
		{Name: "objectClass", Type: T_CLASS, ConstantPool: true},
		{Name: "age", Type: T_LONG, ConstantPool: false},
		{Name: "size", Type: T_LONG, ConstantPool: false},
		{Name: "weight", Type: T_FLOAT, ConstantPool: false},
	},
}

//datadog.HeapUsage:class{name: datadog.HeapUsage, id: 121, fields: [
//{Name:startTime Type:11 ConstantPool:false Array:false}
//{Name:size Type:11 ConstantPool:false Array:false}
//{Name:isLive Type:4 ConstantPool:false Array:false}]}

var Type_Datadog_HeapUsage = def.Class{
	Name: "datadog.HeapUsage",
	ID:   def.TypeID(121),
	Fields: []def.Field{
		{Name: "startTime", Type: T_LONG, ConstantPool: false},
		{Name: "size", Type: T_LONG, ConstantPool: false},
		{Name: "isLive", Type: T_BOOLEAN, ConstantPool: false},
	},
}

// datadog.MethodSample:class{name: datadog.MethodSample, id: 102, fields: [
//{Name:startTime Type:11 ConstantPool:false Array:false}
//{Name:eventThread Type:22 ConstantPool:true Array:false}
//{Name:stackTrace Type:26 ConstantPool:true Array:false}
//{Name:state Type:25 ConstantPool:true Array:false}
//{Name:mode Type:33 ConstantPool:true Array:false}
//{Name:weight Type:11 ConstantPool:false Array:false}
//{Name:spanId Type:11 ConstantPool:false Array:false}
//{Name:localRootSpanId Type:11 ConstantPool:false Array:false}
//{Name:_dd.trace.operation Type:32 ConstantPool:true Array:false}]}

var Type_Datadog_MethodSample = def.Class{
	Name: "datadog.MethodSample",
	ID:   T_Datadog_MethodSample,
	Fields: []def.Field{
		{Name: "startTime", Type: T_LONG, ConstantPool: false},
		{Name: "eventThread", Type: T_THREAD, ConstantPool: true},
		{Name: "stackTrace", Type: T_STACK_TRACE, ConstantPool: true},
		{Name: "state", Type: T_THREAD_STATE, ConstantPool: true},
		//{Name: "mode", Type: T_ExecutionMode, ConstantPool: true},
		{Name: "weight", Type: T_FLOAT, ConstantPool: false},
	},
}

//datadog.QueueTime:class{name: datadog.QueueTime, id: 123, fields: [
//{Name:startTime Type:11 ConstantPool:false Array:false}
//{Name:duration Type:11 ConstantPool:false Array:false}
//{Name:eventThread Type:22 ConstantPool:true Array:false}
//{Name:origin Type:22 ConstantPool:true Array:false}
//{Name:task Type:21 ConstantPool:true Array:false}
//{Name:scheduler Type:21 ConstantPool:true Array:false}
//{Name:spanId Type:11 ConstantPool:false Array:false}
//{Name:localRootSpanId Type:11 ConstantPool:false Array:false}
//{Name:_dd.trace.operation Type:32 ConstantPool:true Array:false}]}

var Type_QueueTime = def.Class{
	Name: "datadog.QueueTime",
	ID:   def.TypeID(123),
	Fields: []def.Field{
		{Name: "startTime", Type: T_LONG, ConstantPool: false},
		{Name: "duration", Type: T_LONG, ConstantPool: false},
		{Name: "eventThread", Type: T_THREAD, ConstantPool: true},
		{Name: "origin", Type: T_THREAD, ConstantPool: true},
		{Name: "task", Type: T_CLASS, ConstantPool: true},
		{Name: "scheduler", Type: T_CLASS, ConstantPool: true},
	},
}

//datadog.WallClockSamplingEpoch:class{name: datadog.WallClockSamplingEpoch, id: 118, fields: [
//{Name:startTime Type:11 ConstantPool:false Array:false}
//{Name:duration Type:11 ConstantPool:false Array:false}
//{Name:samplePoolSize Type:10 ConstantPool:false Array:false}
//{Name:numSuccessfulSamples Type:10 ConstantPool:false Array:false}
//{Name:numFailedSamples Type:10 ConstantPool:false Array:false}
//{Name:numExitedThreads Type:10 ConstantPool:false Array:false}
//{Name:numPermissionDenied Type:10 ConstantPool:false Array:false}]}

var Type_WallClockSamplingEpoch = def.Class{
	Name: "datadog.WallClockSamplingEpoch",
	ID:   T_WallClockSamplingEpoch,
	Fields: []def.Field{
		{Name: "startTime", Type: T_LONG, ConstantPool: false},
		{Name: "duration", Type: T_LONG, ConstantPool: false},
		{Name: "samplePoolSize", Type: T_INT, ConstantPool: false},
		{Name: "numSuccessfulSamples", Type: T_INT, ConstantPool: false},
		{Name: "numFailedSamples", Type: T_INT, ConstantPool: false},
		{Name: "numExitedThreads", Type: T_INT, ConstantPool: false},
		{Name: "numPermissionDenied", Type: T_INT, ConstantPool: false},
	},
}

//datadog.ExecutionSample:class{name: datadog.ExecutionSample, id: 101, fields: [
//{Name:startTime Type:11 ConstantPool:false Array:false}
//{Name:eventThread Type:22 ConstantPool:true Array:false}
//{Name:stackTrace Type:26 ConstantPool:true Array:false}
//{Name:state Type:25 ConstantPool:true Array:false}
//{Name:mode Type:33 ConstantPool:true Array:false}
//{Name:weight Type:11 ConstantPool:false Array:false}
//{Name:spanId Type:11 ConstantPool:false Array:false}
//{Name:localRootSpanId Type:11 ConstantPool:false Array:false}
//{Name:_dd.trace.operation Type:32 ConstantPool:true Array:false}]}

var Type_datadogExecutionSample = def.Class{
	Name: "datadog.ExecutionSample",
	ID:   T_EXECUTION_SAMPLE,
	Fields: []def.Field{
		{Name: "startTime", Type: T_LONG, ConstantPool: false},
		{Name: "sampledThread", Type: T_THREAD, ConstantPool: true},
		{Name: "stackTrace", Type: T_STACK_TRACE, ConstantPool: true},
		{Name: "state", Type: T_THREAD_STATE, ConstantPool: true},
		//{Name: "mode", Type: T_ExecutionMode, ConstantPool: true}, skip
		{Name: "weight", Type: T_LONG, ConstantPool: false},
	},
}

//class{name: datadog.ExceptionSample, id: 5818, fields: [
//{Name:startTime Type:204 ConstantPool:false Array:false}
//{Name:duration Type:204 ConstantPool:false Array:false}
//{Name:eventThread Type:162 ConstantPool:true Array:false}
//{Name:stackTrace Type:186 ConstantPool:true Array:false}
//{Name:type Type:212 ConstantPool:false Array:false}
//{Name:message Type:212 ConstantPool:false Array:false}
//{Name:sampled Type:210 ConstantPool:false Array:false}
//{Name:firstOccurrence Type:210 ConstantPool:false Array:false}
//{Name:localRootSpanId Type:204 ConstantPool:false Array:false}
//{Name:spanId Type:204 ConstantPool:false Array:false}
//]}

var Type_DatadogException_sample = def.Class{
	Name: "datadog.ExceptionSample",
	ID:   T_Datadog_Excepion_Sample,
	Fields: []def.Field{
		{Name: "startTime", Type: T_LONG, ConstantPool: false},
		{Name: "duration", Type: T_LONG, ConstantPool: false},
		{Name: "eventThread", Type: T_THREAD, ConstantPool: true},
		{Name: "stackTrace", Type: T_STACK_TRACE, ConstantPool: true},
		{Name: "type", Type: T_STRING, ConstantPool: false},
		{Name: "message", Type: T_STRING, ConstantPool: false},
		{Name: "sampled", Type: T_BOOLEAN, ConstantPool: false},
		{Name: "firstOccurrence", Type: T_BOOLEAN, ConstantPool: false},
	},
}

//class{name: jdk.types.GCName, id: 169, fields: [{Name:name Type:212 ConstantPool:false Array:false}]}

var Type_Types_GCName = def.Class{
	Name: "jdk.types.GCName",
	ID:   T_GCName,
	Fields: []def.Field{
		{Name: "name", Type: T_STRING, ConstantPool: false},
	},
}

//class{name: jdk.types.GCCause, id: 170, fields: [{Name:cause Type:212 ConstantPool:false Array:false}]}

var Type_Types_GCCause = def.Class{
	Name: "jdk.types.GCCause",
	ID:   T_GCCause,
	Fields: []def.Field{
		{Name: "name", Type: T_STRING, ConstantPool: false},
	},
}

// class{name: jdk.types.G1YCType, id: 173, fields: [{Name:type Type:212 ConstantPool:false Array:false}]}
var Type_G1YCType = def.Class{
	Name: "jdk.types.G1YCType",
	ID:   T_G1YCType,
	Fields: []def.Field{
		{Name: "type", Type: T_STRING, ConstantPool: false},
	},
}

//class{name: jdk.GarbageCollection, id: 35, fields: [{Name:startTime Type:204 ConstantPool:false Array:false}
//{Name:duration Type:204 ConstantPool:false Array:false}
//{Name:gcId Type:205 ConstantPool:false Array:false}{Name:name Type:169 ConstantPool:true Array:false}
//{Name:cause Type:170 ConstantPool:true Array:false}
//{Name:sumOfPauses Type:204 ConstantPool:false Array:false}
//{Name:longestPause Type:204 ConstantPool:false Array:false}]}

var Type_GarbageCollection = def.Class{
	Name: "jdk.GarbageCollection",
	ID:   T_GarbageCollection,
	Fields: []def.Field{
		{Name: "startTime", Type: T_LONG, ConstantPool: false},
		{Name: "duration", Type: T_LONG, ConstantPool: false},
		{Name: "gcId", Type: T_INT, ConstantPool: false},
		{Name: "name", Type: T_GCName, ConstantPool: true},
		{Name: "cause", Type: T_GCCause, ConstantPool: true},
		{Name: "sumOfPauses", Type: T_LONG, ConstantPool: false},
		{Name: "longestPause", Type: T_LONG, ConstantPool: false},
	},
}

//class{name: jdk.SystemGC, id: 36, fields: [{Name:startTime Type:204 ConstantPool:false Array:false} {Name:duration Type:204 ConstantPool:false Array:false}
//{Name:eventThread Type:162 ConstantPool:true Array:false}
//{Name:stackTrace Type:186 ConstantPool:true Array:false} {Name:invokedConcurrent Type:210 ConstantPool:false Array:false}]}

var Type_SystemGC = def.Class{
	Name: "jdk.SystemGC",
	ID:   T_SystemGC,
	Fields: []def.Field{
		{Name: "startTime", Type: T_LONG, ConstantPool: false},
		{Name: "duration", Type: T_LONG, ConstantPool: false},
		{Name: "eventThread", Type: T_THREAD, ConstantPool: true},
		{Name: "stackTrace", Type: T_STACK_TRACE, ConstantPool: true},
		{Name: "invokedConcurrent", Type: T_BOOLEAN, ConstantPool: false},
	},
}

//class{name: jdk.ParallelOldGarbageCollection, id: 37, fields: [{Name:startTime Type:204 ConstantPool:false Array:false}
//{Name:duration Type:204 ConstantPool:false Array:false}
//{Name:gcId Type:205 ConstantPool:false Array:false} {Name:densePrefix Type:204 ConstantPool:false Array:false}]}

var Type_ParallelOldGarbageCollection = def.Class{
	Name: "jdk.ParallelOldGarbageCollection",
	ID:   T_ParallelOldGarbageCollection,
	Fields: []def.Field{
		{Name: "startTime", Type: T_LONG, ConstantPool: false},
		{Name: "duration", Type: T_LONG, ConstantPool: false},
		{Name: "gcId", Type: T_INT, ConstantPool: false},
		{Name: "densePrefix", Type: T_LONG, ConstantPool: false},
	},
}

//class{name: jdk.YoungGarbageCollection, id: 38, fields: [{Name:startTime Type:204 ConstantPool:false Array:false}
//{Name:duration Type:204 ConstantPool:false Array:false}
//{Name:gcId Type:205 ConstantPool:false Array:false} {Name:tenuringThreshold Type:205 ConstantPool:false Array:false}]}

var Type_YoungGarbageCollection = def.Class{
	Name: "jdk.YoungGarbageCollection",
	ID:   T_YoungGarbageCollection,
	Fields: []def.Field{
		{Name: "startTime", Type: T_LONG, ConstantPool: false},
		{Name: "duration", Type: T_LONG, ConstantPool: false},
		{Name: "gcId", Type: T_INT, ConstantPool: false},
		{Name: "tenuringThreshold", Type: T_INT, ConstantPool: false},
	},
}

//class{name: jdk.OldGarbageCollection, id: 39, fields: [{Name:startTime Type:204 ConstantPool:false Array:false}
//{Name:duration Type:204 ConstantPool:false Array:false} {Name:gcId Type:205 ConstantPool:false Array:false}]}

var Type_OldGarbageCollection = def.Class{
	Name: "jdk.OldGarbageCollection",
	ID:   T_OldGarbageCollection,
	Fields: []def.Field{
		{Name: "startTime", Type: T_LONG, ConstantPool: false},
		{Name: "duration", Type: T_LONG, ConstantPool: false},
		{Name: "gcId", Type: T_INT, ConstantPool: false},
	},
}

//class{name: jdk.G1GarbageCollection, id: 40, fields: [{Name:startTime Type:204 ConstantPool:false Array:false}
//{Name:duration Type:204 ConstantPool:false Array:false}
//{Name:gcId Type:205 ConstantPool:false Array:false}
//{Name:type Type:173 ConstantPool:true Array:false}]}

var Type_G1GarbageCollection = def.Class{
	Name: "jdk.G1GarbageCollection",
	ID:   T_G1GarbageCollection,
	Fields: []def.Field{
		{Name: "startTime", Type: T_LONG, ConstantPool: false},
		{Name: "duration", Type: T_LONG, ConstantPool: false},
		{Name: "gcId", Type: T_INT, ConstantPool: false},
		{Name: "type", Type: T_G1YCType, ConstantPool: true},
	},
}
