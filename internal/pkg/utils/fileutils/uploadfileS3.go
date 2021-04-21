package fileutils

import (
	"Redioteka/internal/pkg/utils/log"
	"Redioteka/internal/pkg/utils/randstring"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"io"
)

func UploadFileS3(reader io.Reader, path, ext string) (string, error) {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState:session.SharedConfigEnable,
		Config: aws.Config{
			Region:aws.String("ru-msk"),
			Endpoint:aws.String("http://hb.bizmrg.com"),
		},
	}))

	bucket := "redtech_static"

	filename := randstring.RandString(32) + ext
	log.Log.Info("Created file with name " + filename)

	uploader := s3manager.NewUploader(sess)

	avatar, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(path + filename),
		Body: reader,
	})
	if err != nil {
		return "", fmt.Errorf("file uploading to s3 error %s", err)
	}

	svc := s3.New(sess)
	_, err = svc.PutObjectAcl(&s3.PutObjectAclInput{
		ACL:    aws.String("public-read"),
		Bucket: aws.String(bucket),
		Key:    aws.String(path + filename),
	})

	return avatar.Location, err
}