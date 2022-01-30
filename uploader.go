package main

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func main() {
	const (
		YOURACCESSKEYID     = "US1WYQ83Z8I86JZPZ3ZPP"
		YOURSECRETACCESSKEY = "X1jVAsxxbjcnxFdqIVuoZLiidK7h2sYm8JG0rEuD"
		YOURENDPOINT        = "172.17.0.5:9000"
		YOURBUCKET          = "h5p"
	)

	minioClient, err := minio.New(YOURENDPOINT, &minio.Options{
		Creds:  credentials.NewStaticV4(YOURACCESSKEYID, YOURSECRETACCESSKEY, ""),
		Secure: false,
	})
	if err != nil {
		log.Fatalln(err)
	}
	minioClient.TraceOn(os.Stdout)

	input := make(chan minio.SnowballObject, 1)
	opts := minio.SnowballOptions{
		Opts:     minio.PutObjectOptions{},
		InMemory: true,
		Compress: true,
	}

	//rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	// prefix := []byte("aaaaaaaaaaaaaaa")
	// for i := range prefix {
	// 	prefix[i] += byte(rng.Intn(25))
	// }

	// Generate
	go func() {
		defer close(input)

		files, err := ioutil.ReadDir("./test")
		if err != nil {
			log.Fatal(err)
		}

		for _, f := range files {
			fmt.Println(f.Size())
			size := f.Size()
			key := fmt.Sprintf(f.Name())
			input <- minio.SnowballObject{
				Key:     key,
				Size:    int64(size),
				ModTime: time.Now(),
				Content: bytes.NewBuffer(make([]byte, size)),
				Close: func() {
					fmt.Println(key, "Close function called")
				},
			}
		}

		// // Create 100 objects
		// for i := 0; i < 100; i++ {
		// 	// With random size 0 -> 10000
		// 	size := rng.Intn(10000)
		// 	key := fmt.Sprintf("%s/%d-%d.bin", string(prefix), i, size)
		// 	input <- minio.SnowballObject{
		// 		// Create path to store objects within the bucket.
		// 		Key:     key,
		// 		Size:    int64(size),
		// 		ModTime: time.Now(),
		// 		Content: bytes.NewBuffer(make([]byte, size)),
		// 		Close: func() {
		// 			fmt.Println(key, "Close function called")
		// 		},
		// 	}
		// }
	}()

	err = minioClient.PutObjectsSnowball(context.TODO(), YOURBUCKET, opts, input)
	if err != nil {
		log.Fatalln(err)
	}
	// Objects successfully uploaded.

}
