#!/bin/bash

Copyright() {
		cat << EOF
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
EOF
}

# export GOMAXPROCS=8;  Sets max nbr of CPUs. 
# See: https://github.com/golang-standards/project-layout
#      https://thenewstack.io/understanding-golang-packages/
#      https://blog.golang.org/godoc-documenting-go-code
#      https://dave.cheney.net/2015/05/12/introducing-gb
#      https://blog.golang.org/pipelines0
#	   https://magefile.org/
#	   https://segment.com/blog/5-advanced-testing-techniques-in-go/
#	   https://play.golang.org/p/9tHAm_zm0L	
#	   https://appliedgo.net/flow/
#	   https://github.com/ryanpeach/goflow
#	   https://gobyexample.com/command-line-flags
#	   https://github.com/golang/go/wiki/PackagePublishing
#      https://github.com/golang-standards/project-layout
	
export GOMAXPROCS=8
export  GOPATH=/home/tyoung3/go
export pgm=flowfib

case $1 in
	b)shift; go build  ${pgm}.go;; #strip ${pgm};;
	d)shift; go doc  $*|pandoc -t html5 --toc >/tmp/$pgm.html&&$BROWSER /tmp/$pgm.html 
			go doc  $*|pandoc -t markdown --toc > README.md 
		;;
	g)shift; godoc -http=:6060 &
			 $BROWSER taos:6060 &
			 ;; 
	ix)shift; go install -ldflags='-X ${pgm}.Foo="Test string generation"'  ;;
	i)shift; go install   ;;
	m)shift; gnome-system-monitor &;;
	q)shift; taskset -c $* time  go run ${pgm}.go;;
	r)shift; go run -race ${pgm}.go $*;;
	t)shift;taskset -c 6,7 time go run ${pgm}.go  | wc ;;
	v|vet)shift;go tool vet -v *.go;;
	x)$EDITOR $0 *go&;;
	*)cat <<- EOF
	
	$0 USAGE:
	 	b		. Build ${pgm} from ${pgm}.go
	 	d [-cmd] [-u] [-c]		. Document ${pgm}.go
	 	i		. Install ${pgm}.go w/variable string
	 	q CPU_list		. run, but limit CPUs
	 	r 	N	. run ${pgm}.go.  N=number of goroutines
	 	t  		. Run with CPUs 6 and 7 
	 	v		. Vet all *.go files 
	 	x		. Edit $0
	 	*		. Show this usage
				
EOF
		;;

esac
