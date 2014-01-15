package main

import (
	"flag"
	"fmt"
	"launchpad.net/goamz/aws"
	"launchpad.net/goamz/s3"
	"log"
)

var bucketName, prefix, marker string
var items = []string{}
var folders = []string{}

// AWS_ACCESS_KEY, AWS_SECRET_ACCESS_KEY env vars
var auth, err = aws.EnvAuth()

// new s3 connection + bucket obj
var s = s3.New(auth, aws.USEast)

func init() {
	flag.StringVar(&bucketName, "b", "", "bucket to list")
	flag.StringVar(&prefix, "p", "", "prefix")
}

func bucketList(name, prefix string, marker string) {
	log.Println("listing bucket at prefix", prefix)
	log.Println("items has ", len(items), " keys")
	bucket := s.Bucket(name)

	// list out bucket contents
	list, err := bucket.List(prefix, "/", marker, 2000)
	if err != nil {
		log.Println(err.Error())
		panic(err.Error())
	}

	// append any files to items
	if len(list.Contents) > 0 {
		log.Println("found", len(list.Contents), "files in prefix ", prefix)

		for _, k := range list.Contents {
			items = append(items, k.Key)
		}

		if list.IsTruncated {
			last := items[len(items)-1]
			bucketList(name, prefix, last)
		}
	}

	// recurse over each folder
	if len(list.CommonPrefixes) > 0 {
		log.Println("found", len(list.CommonPrefixes), "folders in prefix ", prefix)

		for _, p := range list.CommonPrefixes {
			log.Println(p)
			bucketList(name, p, "")
		}
	} else {
		log.Println("no more subfolders in ", prefix)
		return
	}
}

func main() {
	flag.Parse()
	bucketList(bucketName, prefix, marker)
}
