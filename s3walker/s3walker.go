package s3walker

// USAGE
//
// go gowalks3 -b name.of.my.bucket.com -p path/to/prefix -a ACCESS_KEY_ID -s SECRET_ACCESS_KEY
//
import (
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

func ListFiles(accessKey, secretKey, name, prefix, marker string) []string {
	var s3Conn = authS3(accessKey, secretKey)
	bucketWalker(s3Conn, name, prefix, marker)
	return items
}

func bucketWalker(s *s3.S3, name, prefix, marker string) {
	bucket := s.Bucket(name)

	// list out bucket contents
	list, err := bucket.List(prefix, "/", marker, 1000)
	if err != nil {
		log.Println(err.Error())
		panic(err.Error())
	}

	// append any files to items
	if len(list.Contents) > 0 {
		for _, k := range list.Contents {
			items = append(items, k.Key)
		}

		if list.IsTruncated {
			last := items[len(items)-1]
			bucketWalker(s, name, prefix, last)
		}
	}

	// recurse over each folder
	if len(list.CommonPrefixes) > 0 {
		for _, p := range list.CommonPrefixes {
			bucketWalker(s, name, p, "")
		}
	} else {
		return
	}
}
