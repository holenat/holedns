package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"time"
)

func main() {
	flgConf := flag.String("c", "", "config file path")
	flag.Parse()

	log.SetFlags(log.Lshortfile)

	cfg, err := ParseConfig(*flgConf)
	if err != nil {
		fmt.Println(err)
		return
	}

	log.Printf("%v\n", cfg)

	rand.Seed(time.Now().Unix())

	r, err := NewResolve(&cfg.ResolverConfig)
	if err != nil {
		log.Printf("new resolver fail: %v\n", err)
		return
	}

	log.Printf("%v\n", r.Run())
}
