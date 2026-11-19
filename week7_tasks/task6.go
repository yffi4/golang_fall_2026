package main

import (
	"fmt"
	"runtime"
)

// Task 6: Watch the GC
//
// Run (from week7_tasks/):
//   go build -o /tmp/gcdemo task6.go   # or any writable dir
//   GOGC=50   /tmp/gcdemo
//   GOGC=100  /tmp/gcdemo
//   GOGC=200  /tmp/gcdemo
//   GOGC=400  /tmp/gcdemo
//   GOGC=off  /tmp/gcdemo
//   GOGC=off GOMEMLIMIT=100MiB /tmp/gcdemo
//   GODEBUG=gctrace=1 /tmp/gcdemo   # to see individual GC cycle lines

type Record struct {
	ID   int
	Data [256]byte
}

var live []*Record  // stays reachable: the live heap
var discard *Record // a global: storing into it forces every Record to escape
var sum int

//go:noinline
func consume(r *Record) int { return r.ID }

func main() {
	for i := 0; i < 80_000; i++ {
		live = append(live, &Record{ID: i})
	}
	for i := 0; i < 10_000_000; i++ {
		discard = &Record{ID: i} // storing into a global leaks it -> escapes;
		sum += consume(discard)  // garbage again the moment the next iteration overwrites discard
	}

	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Println("sum:            ", sum)
	fmt.Println("GC cycles:      ", m.NumGC)
	fmt.Println("heap from OS:   ", m.Sys>>20, "MB")
	fmt.Println("total pause:    ", float64(m.PauseTotalNs)/1e6, "ms")
	runtime.KeepAlive(live)
}

// Table actually measured by running the built binary with each GOGC
// setting (Apple M5 Pro, 15 logical CPUs, go1.27.1; "heap from OS" = Sys,
// "total pause" = PauseTotalNs):
//
//   GOGC    GC cycles   heap from OS   total pause
//   50        236          50 MB         3.9 ms
//   100       118          65 MB         2.4 ms
//   200        59          87 MB         1.1 ms
//   400        30         144 MB         0.6 ms
//   off          0        2894 MB          0 ms
//   off+GOMEMLIMIT=100MiB  45  103 MB     1.6 ms
//
// Same shape as the handout's table: lower GOGC -> more cycles, less
// memory, more total pause time; GOGC=off never collects and the heap
// grows to the full ~2.7 GB of garbage churned through the second loop
// (80,000 x ~272 B live + 10,000,000 x ~272 B short-lived); GOGC=off with
// a GOMEMLIMIT caps the heap back down by triggering GC on memory
// pressure instead of on the GOGC ratio.
//
// A real GODEBUG=gctrace=1 line from this run (GOGC=100), field by field:
//   gc 115 @0.514s 0%: 0.008+0.39+0.014 ms clock, 0.12+.../15 P, 47->48->24 MB, 49 MB goal, 0 MB stacks, 0 MB globals, 15 P
//     gc 115       - the 115th GC cycle since the program started
//     @0.514s      - 0.514 seconds since the program started
//     0%           - cumulative % of CPU time spent in GC so far
//     0.008+0.39+0.014 ms clock
//                  - stop-the-world pause (mark start) + concurrent mark
//                    time + stop-the-world pause (mark termination)
//     47->48->24 MB
//                  - heap size when this GC started -> heap size at the
//                    end of marking (it can grow while marking runs
//                    concurrently with your program) -> live heap left
//                    after sweeping
//     49 MB goal   - the heap goal for the *next* cycle: live heap (24 MB)
//                    x (1 + GOGC/100) = 24 x 2 = 48 MB (~49, rounding)
//     0 MB stacks, 0 MB globals
//                  - extra memory this cycle counted from goroutine stacks
//                    / global variables (here: negligible)
//     15 P         - the GC used up to 15 logical processors for this cycle
//
// Question: with GOGC=off there is no GC at all - so why is it the
// slowest run?
//   On THIS machine it actually wasn't: three repeated wall-clock runs
//   gave GOGC=100 ~0.52 s real time vs GOGC=off ~0.40 s - GOGC=off was
//   faster here, not slower. "No GC" means every one of the 10,000,000
//   throwaway Records is allocated normally and never freed, so the heap
//   grows to ~2.7 GB - but this machine has plenty of spare RAM and a
//   fast multi-core concurrent marker (15 P), so paying for frequent,
//   cheap GC cycles under GOGC=100 cost MORE wall-clock time overall than
//   just letting the allocator run free never touched a scarce resource.
//   The handout's claim ("off is slowest") holds in the situation it's
//   actually about: a machine where growing to ~2.7 GB causes real memory
//   pressure - paging, swapping, or hitting a container's memory limit.
//   Once the OS has to do extra work to back that much memory (or the
//   process gets killed and restarted), GOGC=off loses badly - which is
//   also exactly why GOMEMLIMIT exists: it's the setting that protects you
//   from that failure mode without giving up GOGC=off's lower CPU
//   overhead when memory isn't actually tight. This is the handout's
//   "measure, don't assume" lesson applying to itself: the right choice
//   is workload- and machine-dependent, not a universal ranking.
//
// Why does GOMEMLIMIT=100MiB together with GOGC=off work so well?
//   GOGC=off disables the "double the live heap" trigger, but
//   GOMEMLIMIT gives the runtime a second, independent trigger: "run a GC
//   when total memory gets close to this limit." With GOGC=off and
//   GOMEMLIMIT=100MiB, the GC runs purely to stay under 100 MiB, which
//   keeps the heap small (so allocation stays cheap) while still skipping
//   the "double on every cycle" heuristic that made GOGC=50 do far more
//   cycles than necessary for this particular workload. It gets the best
//   of both: a size ceiling without GOGC's aggressive doubling schedule.
//
// Which GOGC for a low-memory service vs. a batch job on a big server?
//   Low-memory service: a smaller GOGC (e.g. 50-100), or better, set
//   GOMEMLIMIT to a fixed byte budget slightly under the container's
//   memory limit - predictable memory use matters more than CPU here, and
//   GOMEMLIMIT reacts directly to actual memory pressure instead of only
//   to the live-heap ratio.
//   Batch job on a big server: a larger GOGC (e.g. 200-400) or even
//   GOGC=off if the job is short-lived and memory is plentiful - trading
//   memory for fewer, cheaper GC cycles minimizes total run time when
//   there's no tight memory ceiling to respect.
