package main

import (
		"log"
	"os"

	"github.com/alunegov/kdlscan2"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
	}

	cmd := os.Args[1]
	switch cmd {
	case "scan":
		scan(os.Args)
	case "update":
		update(os.Args)
	case "generate":
		generate(os.Args)
	case "sync":
		sync(os.Args)
	default:
		log.Println("unsupported command")
		printUsage()
	}
}

func printUsage() {
	log.Fatalln("usage")
}

func scan(args []string) {
	if len(args) < 4 {
		log.Fatalln("scan usage: kdlscan2 scan target(proto) lng [pg [pg...]]")
	}
	log.Printf("scanning %s with %v to %s\r\n", args[3], args[4:], args[2])
	if err := kdlscan2.Scan(args[2], args[3], args[4:]); err != nil {
		log.Fatal(err)
	}
}

func update(args []string) {
	if len(args) < 4 {
		log.Fatalln("update usage: kdlscan2 update proto ref_proto [-!] [-x]")
	}
	markModified := false
	markDeleted := false
	for _, arg := range args[4:] {
		if arg == "-!" {
			markModified = true
		} else if arg == "-x" {
			markDeleted = true
		}
	}
	log.Printf("updating %s with ref %s, markModified = %v, markDeleted = %v \r\n", args[2], args[3],
		markModified, markDeleted)
	if err := kdlscan2.Update(args[2], args[3], markModified, markDeleted); err != nil {
		log.Fatal(err)
	}
}

func generate(args []string) {
	if len(args) < 4 {
		log.Fatalln("update usage: kdlscan2 generate target(lng) proto drc [drc_encoding]")
	}
	drcFileEnc := ""
	if len(args) > 5 {
		drcFileEnc = args[5]
	}
	log.Printf("generating on %s with %s (enc=%s) to %s", args[3], args[4], drcFileEnc, args[2])
	if err := kdlscan2.Generate(args[2], args[3], args[4], drcFileEnc); err != nil {
		log.Fatal(err)
	}
}

func sync(args []string) {
	if len(args) < 4 {
		log.Fatalln("sync usage: kdlscan2 sync proto|lng ref_proto|ref_proto")
	}
	log.Printf("syncing %s with %s", args[2], args[3])
	if err := kdlscan2.Sync(args[2], args[3]); err != nil {
		log.Fatal(err)
	}
}
