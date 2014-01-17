package gowalks3

// USAGE
//
// go gowalks3 -b name.of.my.bucket.com -p path/to/prefix -a ACCESS_KEY_ID -s SECRET_ACCESS_KEY
//
import (
	"flag"
	"launchpad.net/goamz/aws"
	"launchpad.net/goamz/s3"
	"log"
)

var bucketName, prefix, marker, accessKey, secretAccessKey string
var items = []string{}

func authS3(accessKey, secretKey string) *s3.S3 {
	var auth aws.Auth
	var err error
	if accessKey == "" && secretKey == "" {
		auth, err = aws.EnvAuth()
		if err != nil {
			log.Println(err.Error())
			panic(err.Error())
		}

	} else {
		auth = aws.Auth{accessKey, secretAccessKey}
	}
	return s3.New(auth, aws.USEast)
}

func BucketFiles(accessKey, secretKey, name, prefix, marker string) []string {
	var s3Conn = authS3(accessKey, secretKey)
	listBucket(s3Conn, name, prefix, marker)
	return items
}

func listBucket(s *s3.S3, name, prefix, marker string) {
	bucket := s.Bucket(name)

	// list out bucket contents
	list, err := bucket.List(prefix, "/", marker, 2000)
	if err != nil {
		log.Println(err.Error())
		panic(err.Error())
	}

	// append any files to items
	if len(list.Contents) > 0 {
		for _, k := range list.Contents {
			items = append(items, k.Key)
		}

		log.Println("items is", len(items))
		if list.IsTruncated {
			last := items[len(items)-1]
			listBucket(s, name, prefix, last)
		}
	}

	// recurse over each folder
	if len(list.CommonPrefixes) > 0 {
		for _, p := range list.CommonPrefixes {
			listBucket(s, name, p, "")
		}
	} else {
		return
	}
}

func init() {
	flag.StringVar(&bucketName, "b", "", "bucket to list")
	flag.StringVar(&prefix, "p", "", "prefix")
	flag.StringVar(&accessKey, "a", "", "aws access key id")
	flag.StringVar(&secretAccessKey, "s", "", "aws secret access key")
}

func main() {
	flag.Parse()
	log.Println("About to recursively list", bucketName, "at", prefix)
	items := BucketFiles(accessKey, secretAccessKey, bucketName, prefix, marker)
	log.Println("Found", len(items), "files")
}
