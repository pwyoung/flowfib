Copyright (C) 2019 Thomas W. Young, fbp@twyoung.com

Licensed under the Apache License, Version 2.0 (the "License"); you may
not use this file or its derivitaves except in compliance with the
License. You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

flowfib.go
==========

Name
----

    flowfib - computes Fibonacci numbers using a pipeline of goroutines

Synopsys
--------

    flowfib [N:default=10 [R:default=1 [F:default=1]

will create and run N Go routines and N+1 channels. The extra channel
allows main to communicate with the pipeline.

Description
-----------

Flowfib was created to help understand and benchmark large Go pipeline
networks, goroutines, and channels. Flowfib helps find the limitations
of Go facilities and the hardware it runs on.

Flowfib builds a pipeline of re-entrant goroutines and channels. The
pipeline computes Fibonacci numbers, fib(N). \[fib(0) = 0.\] Flowfib
imports Go's "math/big" to handle large integers.

Each routine reads a dataflow(Information Packet) containing the last
two Fibonacci numbers from the previous routine's channel, then computes
the next Fibonacci number, and writes the updated dataflow to the next
channel.

Flowfib writes a flow containing the first two numbers(0,1) to the first
channel, then reads the result from the last channel and displays it.

This procedure will be run R(round) times. On each round, F(number of
flows) will be injected into the pipeline. Setting F larger than one
increases parallelism on multi-core processors. Flowfib forces F to be
no larger than N (else deadlock would result). Increasing the number of
rounds(R) extends the running time without affecting the output. Flowfib
ignores all but the first pipeline result.

Flowfib does not require any special mutex's, signalling, semaphones,
threading etc. to safely coordinate processing, just Go's builtin
channels and goroutines.

Timing Notes
------------

Flowfib computes the more than 200,000 digits of fib(1000000) in less
than 7.2 seconds on a laptop running Ubuntu-18.04.2, consuming about 14
seconds total cpu time on eight CPUs. (NOTE: tested on Linux 64bit
only.)

fib(3000000) nearly maxes out16GB memory, running for about 48 seconds.

fib(4000000) thrashes mightily, swapping out more than half of memory,
finally struggling through to a finish after 3 minutes and eleven
seconds. The lower order digits are 7091204718905974245428xxxxxx5. x's
mark some digits redacted to intrigue the reader. (Note that fib(10\^N)
always ends in five.) fib(300000) ran in less than 21 seconds on a
Raspberry Pi(very nice cross compilation) ( 1.3 seconds on the eight
core laptop). fib(400000) never finished on the Pi.

Example
-------

    $   ./flowfib 1000

Outputs:

    fib(1000) = 43466557686937456435688527675040625802564660517371780402481729089536555417949051890403879840079255169295922593080322634775209689623239873322471161642996440906533187938298969649928516003704476137795166849228875

Author
------

    Tom Young, fbp@twyoung.com
