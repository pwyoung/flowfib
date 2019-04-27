/* 
       
Copyright (C) 2019 Thomas W. Young, fbp@twyoung.com 

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file or its derivitaves except in compliance with the License.
You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

#flowfib.go

##Name

		 flowfib - computes Fibonacci numbers using a pipeline of goroutines
 
##Synopsys
 
		flowfib [N:default=10 [R:default=1 [F:default=1]	
will create and run N Go routines and N+1 channels. The extra channel allows main to communicate with the pipeline. 

##Description 
		
Flowfib was created to help understand and benchmark large Go pipeline
networks, goroutines, and channels.  Flowfib helps find the limitations 
of Go facilities and the hardware it runs on.  

Flowfib builds a pipeline of re-entrant goroutines and channels. 
The pipeline computes Fibonacci numbers, fib(N).  [fib(0) = 0.]
Flowfib imports Go's "math/big" to handle large integers.

Each routine reads a dataflow(Information Packet) containing the last two  
Fibonacci numbers from the previous routine's channel, then
computes the next Fibonacci number, and writes the updated dataflow
to the next channel.

Flowfib writes a flow containing the first two numbers(0,1) to the first 
channel, then reads the result from the last channel and displays it.

This procedure will be run R(round) times.  On each round, 
F(number of flows) 
will be injected into the pipeline. Setting F larger than one increases
parallelism on multi-core processors. Flowfib forces F to be no larger than 
N (else deadlock would result).   Increasing the number of rounds(R)
extends the running time without affecting the output.  Flowfib ignores 
all but the first pipeline result.

Flowfib does not require any special mutex's, signalling, semaphones, 
threading etc. to safely coordinate processing, 
just Go's builtin channels and goroutines.

##Timing

Flowfib computes the more than 200,000 digits of fib(1000000)
in less than 7.2 seconds on a laptop running Ubuntu-18.04.2, 
consuming about 14 seconds total cpu time on eight CPUs. 

fib(3000000) 
nearly maxes out16GB memory, running for about 48 seconds.  

fib(4000000) 
thrashes mightily, swapping out more than half of memory, finally 
struggling through to a finish after 3 minutes and eleven seconds.  
The lower order digits are 7091204718905974245428xxxxxx5.
x's mark some digits redacted to intrigue the reader.  (Note that fib(10^N) 
always ends in five.)

##Example
	$	./flowfib 1000

Outputs:
	fib(1000) = 43466557686937456435688527675040625802564660517371780402481729089536555417949051890403879840079255169295922593080322634775209689623239873322471161642996440906533187938298969649928516003704476137795166849228875
	 
##Author 
	Tom Young, fbp@twyoung.com 
*/
package main

import (
	"fmt"
	"math/big"
	"os"
	"strconv"
)

const S int = 12

type Flow struct {
	f0 big.Int
	f1 big.Int
	id int
}

// Number of goroutines and channels
var N int = 10 /* Number of goroutines (compute fib(N) ) */
var R int = 1   /* Number of rounds (compute fib(N) R*F times) */
var F int = 1    /* Number of flows injected by main into the pipeline
                    in each round */

// Install test string. ?? Fix this
var Foo string = "my test string"

func doit(i int, in chan Flow, out chan Flow) {

	a := 1
	
	for a == 1 {
		f, _ := <-in
		f0 := f.f0
		f1 := f.f1
		f0.Add(&f0, &f1)    // Add f1 to f0
		f.f0, f.f1 = f1, f0 // swap f0 and f1
		f.id = i
		out <- f
	}

}

/*      Set command line parameters
  
  NOTE: init() is special -- runs first(before main).
*/
func init() {
	
	a := len(os.Args)

	if a > 1 {
		N, _ = strconv.Atoi(os.Args[1])  		// Number of goroutines
		if a > 2 {
			R, _ = strconv.Atoi(os.Args[2])		// Number of rounds
			if a > 3 {
				F, _ = strconv.Atoi(os.Args[3]) // Number of flows 
			}
		}
	}

	if N < 1 { // 0 deadlocks.
		N = 1000
	}

	if R < 1 {
		R = 1
	}

	if F > N {
		F = N
	}
}

/* BUG(twy):
   If N is too large (>5500000 ?), many will thrash
*/

func main() {

	// Build a pipe line  of N+1 channels passing the Flow struct
	// ?? Change to pass the struct pointer??
	p := make([]chan Flow, 0, N+1)
	for i := 0; i <= N; i++ {
		p = append(p, make(chan Flow))
	}

	for i := 0; i < N; i++ { // Invoke N goroutines
		go doit(i+1, p[i], p[i+1])
	}
		
	for j := 0; j < R; j++ {
		for i := 0; i < F; i++ {
			f0 := new(Flow)
			f0.f0 = *big.NewInt(0)
			f0.f1 = *big.NewInt(1)
						// Send Flow IP to first routine
			p[0] <- *f0			
		}
		for i := 0; i < F; i++ {
			fl1, _ := <-p[N] 		// Get Flow IP from last routine
			if j == 0 && i == 0 {
				fmt.Printf("\nfib(%d) = ", fl1.id)
				res := &fl1.f0
				fmt.Println(res)
			}
		}
	}
	
	//fmt.Printf("%d routines processed %d flows in each of %d rounds.\n", N, F,  R)
	
	// panic("Show stack")
	// ? Shutdown all channels & routines??
}
