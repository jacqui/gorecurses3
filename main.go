package gowalks3

// USAGE
//
// go gowalks3 -b name.of.my.bucket.com -p path/to/prefix -a ACCESS_KEY_ID -s SECRET_ACCESS_KEY
//
import (
	"flag"
	"github.com/jacqui/gowalks3"
	"launchpad.net/goamz/aws"
	"launchpad.net/goamz/s3"
	"log"
)

var accessKey, secretAccessKey string

func init() {
	flag.StringVar(&bucketName, "b", "", "bucket to list")
	flag.StringVar(&prefix, "p", "", "prefix")
	flag.StringVar(&accessKey, "a", "", "aws access key id")
	flag.StringVar(&secretAccessKey, "s", "", "aws secret access key")
}

func main() {
	flag.Parse()
	log.Println("About to recursively list", bucketName, "at", prefix)
	items := gowalks3.BucketFiles(accessKey, secretAccessKey, bucketName, prefix, marker)
	log.Println("Found", len(items), "files")
}
