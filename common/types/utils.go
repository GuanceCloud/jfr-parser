package types

import (
	"github.com/grafana/jfr-parser/common/types/jdk8"
	"strings"
)

const (
	BufferLostTypeId                  = "org.openjdk.jmc.flightrecorder.bufferlosttypeid"
	TypeIdentifierValueInterpretation = "org.openjdk.jmc.flightrecorder.value_interpretation.type_identifier"
)

var jdk8ToJDK11TypeMaps = map[string]string{
	jdk8.CpuLoad:                            CPULoad,
	jdk8.ExecutionSample:                    ExecutionSample,
	jdk8.ExecutionSamplingInfoEventId:       ExecutionSamplingInfoEventId,
	jdk8.Processes:                          Processes,
	jdk8.OsMemorySummary:                    OSMemorySummary,
	jdk8.OsInformation:                      OSInformation,
	jdk8.CpuInformation:                     CPUInformation,
	jdk8.ThreadAllocationStatistics:         ThreadAllocationStatistics,
	jdk8.HeapConf:                           HeapConf,
	jdk8.GcConf:                             GcConf,
	jdk8.HeapSummary:                        HeapSummary,
	jdk8.AllocInsideTlab:                    AllocInsideTlab,
	jdk8.AllocOutsideTlab:                   AllocOutsideTlab,
	jdk8.VmInfo:                             VmInfo,
	jdk8.ClassLoad:                          ClassLoad,
	jdk8.ClassUnload:                        ClassUnload,
	jdk8.ClassLoadStatistics:                ClassLoadStatistics,
	jdk8.Compilation:                        Compilation,
	jdk8.FileWrite:                          FileWrite,
	jdk8.FileRead:                           FileRead,
	jdk8.SocketWrite:                        SocketWrite,
	jdk8.SocketRead:                         SocketRead,
	jdk8.ThreadPark:                         ThreadPark,
	jdk8.ThreadSleep:                        ThreadSleep,
	jdk8.MonitorEnter:                       MonitorEnter,
	jdk8.MonitorWait:                        MonitorWait,
	jdk8.MetaspaceOom:                       MetaspaceOom,
	jdk8.CodeCacheFull:                      CodeCacheFull,
	jdk8.CodeCacheStatistics:                CodeCacheStatistics,
	jdk8.CodeSweeperStatistics:              CodeSweeperStatistics,
	jdk8.SweepCodeCache:                     SweepCodeCache,
	jdk8.EnvironmentVariable:                EnvironmentVariable,
	jdk8.SystemProperties:                   SystemProperties,
	jdk8.ObjectCount:                        ObjectCount,
	jdk8.GcReferenceStatistics:              GcReferenceStatistics,
	jdk8.OldObjectSample:                    OldObjectSample,
	jdk8.GcPauseL3:                          GCPauseL3,
	jdk8.GcPauseL2:                          GCPauseL2,
	jdk8.GcPauseL1:                          GCPauseL1,
	jdk8.GcPause:                            GCPause,
	jdk8.MetaspaceSummary:                   MetaspaceSummary,
	jdk8.GarbageCollection:                  GarbageCollection,
	jdk8.ConcurrentModeFailure:              ConcurrentModeFailure,
	jdk8.ThrowablesStatistics:               ThrowablesStatistics,
	jdk8.ErrorsThrown:                       ErrorsThrown,
	jdk8.ExceptionsThrown:                   ExceptionsThrown,
	jdk8.CompilerStats:                      CompilerStats,
	jdk8.CompilerFailure:                    CompilerFailure,
	jdk8.UlongFlag:                          UlongFlag,
	jdk8.BooleanFlag:                        BooleanFlag,
	jdk8.StringFlag:                         StringFlag,
	jdk8.DoubleFlag:                         DoubleFlag,
	jdk8.LongFlag:                           LongFlag,
	jdk8.IntFlag:                            IntFlag,
	jdk8.UintFlag:                           UintFlag,
	jdk8.UlongFlagChanged:                   UlongFlagChanged,
	jdk8.BooleanFlagChanged:                 BooleanFlagChanged,
	jdk8.StringFlagChanged:                  StringFlagChanged,
	jdk8.DoubleFlagChanged:                  DoubleFlagChanged,
	jdk8.LongFlagChanged:                    LongFlagChanged,
	jdk8.IntFlagChanged:                     IntFlagChanged,
	jdk8.UintFlagChanged:                    UintFlagChanged,
	jdk8.TimeConversion:                     TimeConversion,
	jdk8.ThreadDump:                         ThreadDump,
	BufferLostTypeId:                        JFRDataLost,
	jdk8.GcConfYoungGeneration:              GCConfYoungGeneration,
	jdk8.GcConfSurvivor:                     GCConfSurvivor,
	jdk8.GcConfTlab:                         GCConfTlab,
	jdk8.JavaThreadStart:                    JavaThreadStart,
	jdk8.JavaThreadEnd:                      JavaThreadEnd,
	jdk8.VmOperations:                       VmOperations,
	jdk8.ThreadStatistics:                   ThreadStatistics,
	jdk8.ContextSwitchRate:                  ContextSwitchRate,
	jdk8.CompilerConfig:                     CompilerConfig,
	jdk8.CodeCacheConfig:                    CodeCacheConfig,
	jdk8.CodeSweeperConfig:                  CodeSweeperConfig,
	jdk8.CompilerPhase:                      CompilerPhase,
	jdk8.GcCollectorG1GarbageCollection:     GCCollectorG1GarbageCollection,
	jdk8.GcCollectorOldGarbageCollection:    GCCollectorOldGarbageCollection,
	jdk8.GcCollectorParoldGarbageCollection: GCCollectorParoldGarbageCollection,
	jdk8.GcCollectorYoungGarbageCollection:  GCCollectorYoungGarbageCollection,
	jdk8.GcDetailedAllocationRequiringGc:    GCDetailedAllocationRequiringGC,
	jdk8.GcDetailedEvacuationFailed:         GCDetailedEvacuationFailed,
	jdk8.GcDetailedEvacuationInfo:           GCDetailedEvacuationInfo,
	jdk8.GcDetailedObjectCountAfterGc:       GCDetailedObjectCountAfterGC,
	jdk8.GcDetailedPromotionFailed:          GCDetailedPromotionFailed,
	jdk8.GcHeapPsSummary:                    GcHeapPsSummary,
	jdk8.GcMetaspaceAllocationFailure:       GcMetaspaceAllocationFailure,
	jdk8.GcMetaspaceChunkFreeListSummary:    GcMetaspaceChunkFreeListSummary,
	jdk8.GcMetaspaceGcThreshold:             GCMetaspaceGCThreshold,
	jdk8.RecordingSetting:                   RecordingSetting,
	jdk8.RECORDINGS:                         Recordings,
	jdk8.GcG1mmu:                            GCG1mmu,
}

func ConvertToJDK11Type(typeID string) string {
	if strings.HasPrefix(typeID, jdk8.JDK9And10Prefix) {
		switch {
		case strings.HasSuffix(typeID, "AllocationRequiringGc"):
			return GCDetailedAllocationRequiringGC
		case strings.HasSuffix(typeID, "GCG1MMU"):
			return GCG1mmu
		default:
			return jdkTypePrefix + typeID[len(jdk8.JDK9And10Prefix):]
		}
	}

	if newType, ok := jdk8ToJDK11TypeMaps[typeID]; ok {
		return newType
	}
	return typeID
}
