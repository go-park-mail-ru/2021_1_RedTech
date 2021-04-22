package repository

import (
	"Redioteka/internal/pkg/domain"
	"Redioteka/internal/pkg/utils/log"
	"Redioteka/internal/pkg/utils/randstring"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"io"
)

type s3AvatarRepository struct {
	region *string
	endpoint *string
	bucketName *string
}

func NewS3AvatarRepository() domain.AvatarRepository {
	s3rep := new(s3AvatarRepository)
	s3rep.region = aws.String("ru-msk")
	s3rep.endpoint = aws.String("http://hb.bizmrg.com")
	s3rep.bucketName = aws.String("redtech_static")
	return s3rep
}

func (s3rep *s3AvatarRepository) UploadAvatar(reader io.Reader, path, ext string) (string, error) {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState:session.SharedConfigEnable,
		Config: aws.Config{
			Region: s3rep.region,
			Endpoint: s3rep.endpoint,
		},
	}))

	filename := randstring.RandString(32) + ext
	log.Log.Info("Created file with name " + filename)

	uploader := s3manager.NewUploader(sess)

	avatar, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: s3rep.bucketName,
		Key:    aws.String(path + filename),
		Body: reader,
	})
	if err != nil {
		return "", fmt.Errorf("file uploading to s3rep error %s", err)
	}

	svc := s3.New(sess)
	_, err = svc.PutObjectAcl(&s3.PutObjectAclInput{
		ACL:    aws.String("public-read"),
		Bucket: s3rep.bucketName,
		Key:    aws.String(path + filename),
	})

	return avatar.Location, err
}
