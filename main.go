package main

import (
	"context"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/sirupsen/logrus"
)

func main() {
	logger := logrus.New()
	logger.Info(len(os.Args))
	if len(os.Args) == 1 || len(os.Args) > 3 {
		logger.Fatal("expected at least 2 args and a maximum of three. Refer to documentation")
	}

	var directory = os.Args[1]
	var fmtstring string
	if len(os.Args) == 3 {
		fmtstring = os.Args[2]
		logger.Infof("using %s as format string for handler", fmtstring)
	}

	ctx := context.Background()
	awsCfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		logger.WithError(err).Fatal("failed to initialize aws config")
	}

	lambdaClient := lambda.NewFromConfig(awsCfg)

	files, err := os.ReadDir(directory)
	if err != nil {
		logger.WithError(err).Fatal("failed to list directory")
	}

	var wg = new(sync.WaitGroup)

	for _, f := range files {
		wg.Add(1)
		if !strings.HasSuffix(f.Name(), ".zip") {
			continue
		}
		go func(dirEnt fs.DirEntry, wg *sync.WaitGroup) {
			defer wg.Done()
			entry := logger.WithField("file", dirEnt.Name())
			file, err := os.Open(filepath.Join(directory, dirEnt.Name()))
			if err != nil {
				entry.WithError(err).Error("failed to open zip")
				return
			}

			data, err := io.ReadAll(file)
			if err != nil {
				entry.WithError(err).Error("failed to read zip")
				return
			}

			handler := strings.TrimSuffix(dirEnt.Name(), ".zip")
			if fmtstring != "" {
				handler = fmt.Sprintf(fmtstring, handler)
			}
			output, err := lambdaClient.UpdateFunctionCode(ctx, &lambda.UpdateFunctionCodeInput{
				FunctionName: aws.String(handler),
				ZipFile:      data,
			})
			if err != nil {
				entry.WithError(err).Error("failed to update function")
				return
			}

			entry.WithFields(logrus.Fields{
				"last_modified": aws.ToString(output.LastModified),
			}).Info("function updated successfully")

		}(f, wg)
	}

	logger.Info("updating functions")
	wg.Wait()

	logger.Info("done")

}
