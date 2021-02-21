package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/h2non/bimg"
)

var (
	dirFlag = flag.String("path", ".", "resize images")
)

func main() {
	flag.Parse()
	err := filepath.Walk(*dirFlag,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if strings.HasSuffix(path, ".jpg") || strings.HasSuffix(path, ".jpeg") {
				buffer, err := bimg.Read(path)
				if err != nil {
					fmt.Println(err)
					return err
				}
				options := bimg.Options{}
				size, err := bimg.NewImage(buffer).Size()
				if size.Width > 1024 {
					ratio := float64(size.Width) / float64(size.Height)
					if ratio < 1 {
						width := int(1024 * ratio)
						options.Height = 1024
						options.Width = width
					} else {
						height := int(1024 / ratio)
						options.Width = 1024
						options.Height = height
					}
					newImage, err := bimg.NewImage(buffer).Process(options)
					if err != nil {
						fmt.Println(err)
						return err
					}
					bimg.Write(path, newImage)
					fmt.Printf("%s rewrited \n", path)
				}
			}
			return nil
		})
	if err != nil {
		log.Println(err)
	}
}
