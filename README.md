# gorecurses3

## Description
recursively walks a bucket:prefix in s3, building up an array of keys (files)

## Usage

Command line:

```
git clone https://github.com/jacqui/gorecurses3.git
cd gorecurses3
go run main.go -b BUCKET_NAME -p PREFIX/TO/WALK -a AWS_ACCESS_KEY_ID -s AWS_SECRET_ACCESS_KEY
```

Library:

```
import "github.com/jacqui/gorecurses3/s3walker"
items := s3walker.ListFiles(accessKey, secretAccessKey, bucketName, prefix, marker)
```

