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

func authS3(auth aws.Auth) *s3.S3 {
	var err error
	if auth.AccessKey == "" && auth.SecretKey == "" {
		auth, err = aws.EnvAuth()
		if err != nil {
			log.Println(err.Error())
			panic(err.Error())
		}
	}
	return s3.New(auth, aws.USEast)
}

var items = []s3.Key{}

func ListFiles(auth aws.Auth, name, prefix, marker string) []s3.Key {
	var s3Conn = authS3(auth)
	bucketWalker(s3Conn, name, prefix, marker)
	return items
}

func bucketWalker(s *s3.S3, name, prefix, marker string) {
	bucket := s.Bucket(name)

	// list out bucket contents
	list, err := bucket.List(prefix, "/", marker, 1000)
	if err != nil {
		log.Println(err.Error())
		return
	}

	// append any files to items
	if len(list.Contents) > 0 {
		for _, k := range list.Contents {
			items = append(items, k)
		}

		if list.IsTruncated {
			last := items[len(items)-1]
			bucketWalker(s, name, prefix, last.Key)
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
