package main

// USAGE
//
// go gowalks3 -b name.of.my.bucket.com -p path/to/prefix -a ACCESS_KEY_ID -s SECRET_ACCESS_KEY
//
import (
	"flag"
	"github.com/jacqui/gorecurses3/s3walker"
	"launchpad.net/goamz/aws"
	"log"
)

var bucketName, prefix, marker, accessKey, secretAccessKey string

func init() {
	flag.StringVar(&bucketName, "b", "", "bucket to list")
	flag.StringVar(&prefix, "p", "", "prefix")
	flag.StringVar(&accessKey, "a", "", "aws access key id")
	flag.StringVar(&secretAccessKey, "s", "", "aws secret access key")
}

func main() {
	flag.Parse()
	log.Println("About to recursively list", bucketName, "at", prefix)
	auth := aws.Auth{accessKey, secretAccessKey}
	items := s3walker.ListFiles(auth, bucketName, prefix, marker)
	log.Println("Found", len(items), "files")
}
