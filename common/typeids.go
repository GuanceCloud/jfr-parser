package common

const (
	jdkTypePrefix = "jdk."
	ddTypePrefix  = "datadog."
)

const (
	CpuLoad                      = jdkTypePrefix + "CPULoad"
	ThreadCpuLoad                = jdkTypePrefix + "ThreadCPULoad"
	ExecutionSample              = jdkTypePrefix + "ExecutionSample"
	ExecutionSamplingInfoEventId = jdkTypePrefix + "ExecutionSampling"
	NativeMethodSample           = jdkTypePrefix + "NativeMethodSample"
	Processes                    = jdkTypePrefix + "SystemProcess"
	OSMemorySummary              = jdkTypePrefix + "PhysicalMemory"
	OSInformation                = jdkTypePrefix + "OSInformation"
	CPUInformation               = jdkTypePrefix + "CPUInformation"
	THREAD_ALLOCATION_STATISTICS = jdkTypePrefix + "ThreadAllocationStatistics"
	HeapConf                     = jdkTypePrefix + "GCHeapConfiguration"
	GcConf                       = jdkTypePrefix + "GCConfiguration"
	HeapSummary                  = jdkTypePrefix + "GCHeapSummary"
	AllocInsideTlab              = jdkTypePrefix + "ObjectAllocationInNewTLAB"
	AllocOutsideTlab             = jdkTypePrefix + "ObjectAllocationOutsideTLAB"
	ObjAllocSample               = jdkTypePrefix + "ObjectAllocationSample"
	VmInfo                       = jdkTypePrefix + "JVMInformation"
	ClassDefine                  = jdkTypePrefix + "ClassDefine"
	ClassLoad                    = jdkTypePrefix + "ClassLoad"
	ClassUnload                  = jdkTypePrefix + "ClassUnload"
	ClassLoadStatistics          = jdkTypePrefix + "ClassLoadingStatistics"
	ClassLoaderStatistics        = jdkTypePrefix + "ClassLoaderStatistics"
	COMPILATION                  = jdkTypePrefix + "Compilation"
	FileWrite                    = jdkTypePrefix + "FileWrite"
	FileRead                     = jdkTypePrefix + "FileRead"
	SocketWrite                  = jdkTypePrefix + "SocketWrite"
	SocketRead                   = jdkTypePrefix + "SocketRead"
	ThreadPark                   = jdkTypePrefix + "ThreadPark"
	ThreadSleep                  = jdkTypePrefix + "ThreadSleep"
	MonitorEnter                 = jdkTypePrefix + "JavaMonitorEnter"
	MonitorWait                  = jdkTypePrefix + "JavaMonitorWait"
	MetaspaceOom                 = jdkTypePrefix + "MetaspaceOOM"
	CodeCacheFull                = jdkTypePrefix + "CodeCacheFull"
	CodeCacheStatistics          = jdkTypePrefix + "CodeCacheStatistics"
	CODE_SWEEPER_STATISTICS      = jdkTypePrefix + "CodeSweeperStatistics"
	SweepCodeCache               = jdkTypePrefix + "SweepCodeCache"
	EnvironmentVariable          = jdkTypePrefix + "InitialEnvironmentVariable"
	SystemProperties             = jdkTypePrefix + "InitialSystemProperty"
	ObjectCount                  = jdkTypePrefix + "ObjectCount"
	GcReferenceStatistics        = jdkTypePrefix + "GCReferenceStatistics"
	OldObjectSample              = jdkTypePrefix + "OldObjectSample"
	GcPauseL4                    = jdkTypePrefix + "GCPhasePauseLevel4"
	GcPauseL3                    = jdkTypePrefix + "GCPhasePauseLevel3"
	GcPauseL2                    = jdkTypePrefix + "GCPhasePauseLevel2"
	GcPauseL1                    = jdkTypePrefix + "GCPhasePauseLevel1"
	GcPause                      = jdkTypePrefix + "GCPhasePause"
	MetaspaceSummary             = jdkTypePrefix + "MetaspaceSummary"
	GarbageCollection            = jdkTypePrefix + "GarbageCollection"
	ConcurrentModeFailure        = jdkTypePrefix + "ConcurrentModeFailure"
	ThrowableStatistics          = jdkTypePrefix + "ExceptionStatistics"
	ErrorsThrown                 = jdkTypePrefix + "JavaErrorThrow"
	/*
	 * NOTE: The parser filters all JavaExceptionThrow events created from the Error constructor to
	 * avoid duplicates, so this event type represents 'non error throwables' rather than
	 * exceptions. See note in SyntheticAttributeExtension which does the duplicate filtering.
	 */
	ExceptionsThrown                       = jdkTypePrefix + "JavaExceptionThrow"
	CompilerStats                          = jdkTypePrefix + "CompilerStatistics"
	CompilerFailure                        = jdkTypePrefix + "CompilationFailure"
	ULONG_FLAG                             = jdkTypePrefix + "UnsignedLongFlag"
	BOOLEAN_FLAG                           = jdkTypePrefix + "BooleanFlag"
	STRING_FLAG                            = jdkTypePrefix + "StringFlag"
	DOUBLE_FLAG                            = jdkTypePrefix + "DoubleFlag"
	LONG_FLAG                              = jdkTypePrefix + "LongFlag"
	INT_FLAG                               = jdkTypePrefix + "IntFlag"
	UINT_FLAG                              = jdkTypePrefix + "UnsignedIntFlag"
	ULONG_FLAG_CHANGED                     = jdkTypePrefix + "UnsignedLongFlagChanged"
	BOOLEAN_FLAG_CHANGED                   = jdkTypePrefix + "BooleanFlagChanged"
	STRING_FLAG_CHANGED                    = jdkTypePrefix + "StringFlagChanged"
	DOUBLE_FLAG_CHANGED                    = jdkTypePrefix + "DoubleFlagChanged"
	LONG_FLAG_CHANGED                      = jdkTypePrefix + "LongFlagChanged"
	INT_FLAG_CHANGED                       = jdkTypePrefix + "IntFlagChanged"
	UINT_FLAG_CHANGED                      = jdkTypePrefix + "UnsignedIntFlagChanged"
	TimeConversion                         = jdkTypePrefix + "CPUTimeStampCounter"
	ThreadDump                             = jdkTypePrefix + "ThreadDump"
	JfrDataLost                            = jdkTypePrefix + "DataLoss"
	DUMP_REASON                            = jdkTypePrefix + "DumpReason"
	GC_CONF_YOUNG_GENERATION               = jdkTypePrefix + "YoungGenerationConfiguration"
	GC_CONF_SURVIVOR                       = jdkTypePrefix + "GCSurvivorConfiguration"
	GC_CONF_TLAB                           = jdkTypePrefix + "GCTLABConfiguration"
	JavaThreadStart                        = jdkTypePrefix + "ThreadStart"
	JavaThreadEnd                          = jdkTypePrefix + "ThreadEnd"
	VmOperations                           = jdkTypePrefix + "ExecuteVMOperation"
	VM_SHUTDOWN                            = jdkTypePrefix + "Shutdown"
	THREAD_STATISTICS                      = jdkTypePrefix + "JavaThreadStatistics"
	ContextSwitchRate                      = jdkTypePrefix + "ThreadContextSwitchRate"
	COMPILER_CONFIG                        = jdkTypePrefix + "CompilerConfiguration"
	CodeCacheConfig                        = jdkTypePrefix + "CodeCacheConfiguration"
	CODE_SWEEPER_CONFIG                    = jdkTypePrefix + "CodeSweeperConfiguration"
	COMPILER_PHASE                         = jdkTypePrefix + "CompilerPhase"
	GC_COLLECTOR_G1_GARBAGE_COLLECTION     = jdkTypePrefix + "G1GarbageCollection"
	GcCollectorOldGarbageCollection        = jdkTypePrefix + "OldGarbageCollection"
	GC_COLLECTOR_PAROLD_GARBAGE_COLLECTION = jdkTypePrefix + "ParallelOldGarbageCollection"
	GcCollectorYoungGarbageCollection      = jdkTypePrefix + "YoungGarbageCollection"
	GC_DETAILED_ALLOCATION_REQUIRING_GC    = jdkTypePrefix + "AllocationRequiringGC"
	GC_DETAILED_EVACUATION_FAILED          = jdkTypePrefix + "EvacuationFailed"
	GC_DETAILED_EVACUATION_INFO            = jdkTypePrefix + "EvacuationInformation"
	GC_DETAILED_OBJECT_COUNT_AFTER_GC      = jdkTypePrefix + "ObjectCountAfterGC"
	GC_DETAILED_PROMOTION_FAILED           = jdkTypePrefix + "PromotionFailed"
	GC_HEAP_PS_SUMMARY                     = jdkTypePrefix + "PSHeapSummary"
	GC_METASPACE_ALLOCATION_FAILURE        = jdkTypePrefix + "MetaspaceAllocationFailure"
	GC_METASPACE_CHUNK_FREE_LIST_SUMMARY   = jdkTypePrefix + "MetaspaceChunkFreeListSummary"
	GC_METASPACE_GC_THRESHOLD              = jdkTypePrefix + "MetaspaceGCThreshold"
	GC_G1MMU                               = jdkTypePrefix + "G1MMU"
	GC_G1_EVACUATION_YOUNG_STATS           = jdkTypePrefix + "G1EvacuationYoungStatistics"
	GC_G1_EVACUATION_OLD_STATS             = jdkTypePrefix + "G1EvacuationOldStatistics"
	GC_G1_BASIC_IHOP                       = jdkTypePrefix + "G1BasicIHOP"
	GC_G1_HEAP_REGION_TYPE_CHANGE          = jdkTypePrefix + "G1HeapRegionTypeChange"
	GC_G1_HEAP_REGION_INFORMATION          = jdkTypePrefix + "G1HeapRegionInformation"
	BiasedLockSelfRevocation               = jdkTypePrefix + "BiasedLockSelfRevocation"
	BiasedLockRevocation                   = jdkTypePrefix + "BiasedLockRevocation"
	BiasedLockClassRevocation              = jdkTypePrefix + "BiasedLockClassRevocation"
	GC_G1_ADAPTIVE_IHOP                    = jdkTypePrefix + "G1AdaptiveIHOP"
	RECORDINGS                             = jdkTypePrefix + "ActiveRecording"
	RecordingSetting                       = jdkTypePrefix + "ActiveSetting"

	// SafepointBegin Safepointing begin
	SafepointBegin = jdkTypePrefix + "SafepointBegin"
	// SafepointStateSync Synchronize run state of threads
	SafepointStateSync = jdkTypePrefix + "SafepointStateSynchronization"
	// SafepointWaitBlocked SAFEPOINT_WAIT_BLOCKED Safepointing begin waiting on running threads to block
	SafepointWaitBlocked = jdkTypePrefix + "SafepointWaitBlocked"
	// SafepointCleanup SAFEPOINT_CLEANUP Safepointing begin running cleanup (parent)
	SafepointCleanup = jdkTypePrefix + "SafepointCleanup"
	// SafepointCleanupTask SAFEPOINT_CLEANUP_TASK Safepointing begin running cleanup task, individual subtasks
	SafepointCleanupTask = jdkTypePrefix + "SafepointCleanupTask"
	// SafepointEnd Safepointing end
	SafepointEnd   = jdkTypePrefix + "SafepointEnd"
	MODULE_EXPORT  = jdkTypePrefix + "ModuleExport"
	MODULE_REQUIRE = jdkTypePrefix + "ModuleRequire"
	NATIVE_LIBRARY = jdkTypePrefix + "NativeLibrary"
	HEAP_DUMP      = jdkTypePrefix + "HeapDump"
	PROCESS_START  = jdkTypePrefix + "ProcessStart"

	JavaMonitorInflate = jdkTypePrefix + "JavaMonitorInflate"
)

const (
	ExceptionCount  = ddTypePrefix + "ExceptionCount"
	ExceptionSample = ddTypePrefix + "ExceptionSample"
	ProfilerSetting = ddTypePrefix + "ProfilerSetting"

	// DatadogProfilerConfig 标识符	名称	说明	内容类型
	//datadog.DatadogProfilerConfig	Datadog Profiler Configuration	[datadog.DatadogProfilerConfig]	null
	DatadogProfilerConfig = ddTypePrefix + "DatadogProfilerConfig"
	SCOPE                 = ddTypePrefix + "Scope"

	// DatadogExecutionSample EXECUTION_SAMPLE 标识符	名称	说明	内容类型
	//datadog.ExecutionSample	Method CPU Profiling Sample	[datadog.ExecutionSample]	null
	DatadogExecutionSample = ddTypePrefix + "ExecutionSample"

	// DatadogMethodSample MethodSample METHOD_SAMPLE 标识符	名称	说明	内容类型
	//datadog.MethodSample	Method Wall Profiling Sample	[datadog.MethodSample]	null
	DatadogMethodSample = ddTypePrefix + "MethodSample"

	// AllocationSample ALLOCATION_SAMPLE 标识符	名称	说明	内容类型
	//datadog.ObjectSample	Allocation sample	[datadog.ObjectSample]	null
	AllocationSample = ddTypePrefix + "ObjectSample"

	// HeapLiveObject 标识符	名称	说明	内容类型
	//datadog.HeapLiveObject	Heap Live Object	[datadog.HeapLiveObject]	null
	HeapLiveObject = ddTypePrefix + "HeapLiveObject"

	DatadogEndpoint = ddTypePrefix + "Endpoint"
)
