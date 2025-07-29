package pprof

import (
	"fmt"
	"io"

	"jfr-parser/parser"
)

func ParseJFR(body []byte, pi *ParseInput, jfrLabels *LabelsSnapshot) (res *Profiles, err error) {
	/*	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("jfr parser panic: %v", r)
		}
	}()*/
	p := parser.NewParser(body, parser.Options{
		SymbolProcessor: parser.ProcessSymbols,
	})
	return parse(p, pi, jfrLabels)
}

func parse(parser *parser.Parser, piOriginal *ParseInput, jfrLabels *LabelsSnapshot) (result *Profiles, err error) {
	var event string

	builders := newJfrPprofBuilders(parser, jfrLabels, piOriginal)

	var values = [2]int64{1, 0}

	for {
		typ, err := parser.ParseEvent()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, fmt.Errorf("jfr parser ParseEvent error: %w", err)
		}

		switch typ {
		case parser.TypeMap.T_EXECUTION_SAMPLE, parser.TypeMap.Datadog_ExecutionSample:
			fmt.Printf("T_EXECUTION_SAMPLE Datadog_ExecutionSample %d \n", typ)
			ts := parser.GetThreadState(parser.ExecutionSample.State)
			correlation := StacktraceCorrelation{
				ContextId: parser.ExecutionSample.ContextId,
				SpanId:    parser.ExecutionSample.SpanId,
				SpanName:  parser.ExecutionSample.SpanName,
			}
			if ts != nil && ts.Name != "STATE_SLEEPING" {
				builders.addStacktrace(sampleTypeCPU, correlation, parser.ExecutionSample.StackTrace, values[:1])
			}
			if event == "wall" {
				builders.addStacktrace(sampleTypeWall, correlation, parser.ExecutionSample.StackTrace, values[:1])
			}
		case parser.TypeMap.T_WALL_CLOCK_SAMPLE:
			fmt.Printf("T_WALL_CLOCK_SAMPLE %d \n", typ)
			values[0] = int64(parser.WallClockSample.Samples)
			builders.addStacktrace(sampleTypeWall, StacktraceCorrelation{}, parser.WallClockSample.StackTrace, values[:1])
		case parser.TypeMap.T_ALLOC_IN_NEW_TLAB:
			fmt.Printf("T_ALLOC_IN_NEW_TLAB %d \n", typ)
			values[1] = int64(parser.ObjectAllocationInNewTLAB.TlabSize)
			correlation := StacktraceCorrelation{
				ContextId: parser.ObjectAllocationInNewTLAB.ContextId,
				SpanId:    parser.ObjectAllocationInNewTLAB.SpanId,
				SpanName:  parser.ObjectAllocationInNewTLAB.SpanName,
			}
			builders.addStacktrace(sampleTypeInTLAB, correlation, parser.ObjectAllocationInNewTLAB.StackTrace, values[:2])
		case parser.TypeMap.T_ALLOC_OUTSIDE_TLAB:
			fmt.Printf("T_ALLOC_OUTSIDE_TLAB %d \n", typ)
			values[1] = int64(parser.ObjectAllocationOutsideTLAB.AllocationSize)
			correlation := StacktraceCorrelation{
				ContextId: parser.ObjectAllocationOutsideTLAB.ContextId,
				SpanId:    parser.ObjectAllocationOutsideTLAB.SpanId,
				SpanName:  parser.ObjectAllocationOutsideTLAB.SpanName,
			}
			builders.addStacktrace(sampleTypeOutTLAB, correlation, parser.ObjectAllocationOutsideTLAB.StackTrace, values[:2])
		case parser.TypeMap.T_ALLOC_SAMPLE, parser.TypeMap.Datadog_HeapUsage:
			fmt.Printf("T_ALLOC_SAMPLE Datadog_HeapUsage %d \n", typ)
			values[1] = int64(parser.ObjectAllocationSample.Weight)
			builders.addStacktrace(sampleTypeAllocSample, StacktraceCorrelation{}, parser.ObjectAllocationSample.StackTrace, values[:2])
		case parser.TypeMap.T_MONITOR_ENTER:
			fmt.Printf("T_MONITOR_ENTER %d \n", typ)
			values[1] = int64(parser.JavaMonitorEnter.Duration)
			correlation := StacktraceCorrelation{
				ContextId: parser.JavaMonitorEnter.ContextId,
				SpanId:    parser.JavaMonitorEnter.SpanId,
				SpanName:  parser.JavaMonitorEnter.SpanName,
			}
			builders.addStacktrace(sampleTypeLock, correlation, parser.JavaMonitorEnter.StackTrace, values[:2])
		case parser.TypeMap.T_THREAD_PARK:
			fmt.Printf("T_THREAD_PARK %d \n", typ)
			values[1] = int64(parser.ThreadPark.Duration)
			builders.addStacktrace(sampleTypeThreadPark, StacktraceCorrelation{}, parser.ThreadPark.StackTrace, values[:2])
		case parser.TypeMap.T_LIVE_OBJECT:
			fmt.Printf("T_LIVE_OBJECT %d \n", typ)
			builders.addStacktrace(sampleTypeLiveObject, StacktraceCorrelation{}, parser.LiveObject.StackTrace, values[:1])
		case parser.TypeMap.T_MALLOC:
			fmt.Printf("T_MALLOC %d \n", typ)
			values[1] = int64(parser.Malloc.Size)
			builders.addStacktrace(sampleTypeMalloc, StacktraceCorrelation{}, parser.Malloc.StackTrace, values[:2])
		case parser.TypeMap.T_ObjectSample:
			values[0] = int64(parser.ObjectSample.Weight)
			values[1] = int64(parser.ObjectSample.Size)
			builders.addStacktrace(sampleObjectSample, StacktraceCorrelation{}, parser.ObjectSample.StackTrace, values[:2])
		case parser.TypeMap.T_ACTIVE_SETTING:
			//fmt.Printf("T_ACTIVE_SETTING %d \n", typ)
			if parser.ActiveSetting.Name == "event" {
				event = parser.ActiveSetting.Value
			}
			//fmt.Printf("%+v \n", parser.ActiveSetting)
		default:
			fmt.Printf("skip tyid =%d \n", typ)
		}
	}
	fmt.Printf("parser %+v", parser.TypeMap)
	result = builders.build(event)

	return result, nil
}
