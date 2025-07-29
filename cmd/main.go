package main

import (
	"fmt"
	"jfr-parser/pprof"
	"os"
	"time"
)

func Format(buf []byte) {
	pi := &pprof.ParseInput{
		StartTime:  time.Now(),
		EndTime:    time.Now(),
		SampleRate: 100,
	}
	profiles, err := pprof.ParseJFR(buf, pi, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	for i := 0; i < len(profiles.Profiles); i++ {
		fmt.Println(profiles.Profiles[i].Metric)
		fmt.Println(len(profiles.Profiles[i].Profile.Location))
		fmt.Println(len(profiles.Profiles[i].Profile.Sample))
		fmt.Printf("%+v \n", profiles.Profiles[i].Profile.StringTable)
	}
	fmt.Printf("len of profiles is %d \n", len(profiles.Profiles))
	fmt.Println("----------------end------------------------")
}

func main() {
	//buf, err := os.ReadFile("/home/songlq/opt/async-profile/async-profiler-2.8.3-linux-x64/tmall.jfr")
	buf, err := os.ReadFile("/home/songlq/tmp/dkjfr/metric012025-07-25-10-54.jfr")
	if err != nil {
		panic(err)
	}
	Format(buf)
}
