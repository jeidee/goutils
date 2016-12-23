package aws

import (
	"errors"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type S3 struct {
	session *session.Session
	svc     *s3.S3
	bucket  string
}

var awss3 *S3

func GetS3(params ...string) (*S3, error) {
	if len(params) == 0 {
		if awss3 == nil {
			return nil, errors.New("S3 Storage is not initialized")
		}
		return awss3, nil
	}

	region := ""
	credentialsProfile := ""
	bucket := ""

	if len(params) >= 1 {
		bucket = params[0]
	}

	if len(params) >= 2 {
		region = params[1]
	}

	if len(params) >= 3 {
		credentialsProfile = params[2]
	}

	awss3 = &S3{}

	// Create a new session
	sess, err := func() (*session.Session, error) {
		if credentialsProfile != "" {
			return session.NewSession(&aws.Config{
				Region:      aws.String(region),
				Credentials: credentials.NewSharedCredentials("", credentialsProfile),
			})
		}

		return session.NewSession(&aws.Config{
			Region: aws.String(region),
		})
	}()

	if err != nil {
		return nil, err
	}

	awss3.session = sess
	awss3.svc = s3.New(sess)
	awss3.bucket = bucket

	return awss3, nil
}

func (o *S3) GetBucketList() ([]string, error) {
	list, err := o.svc.ListBuckets(&s3.ListBucketsInput{})
	if err != nil {
		return nil, err
	}

	buckets := make([]string, len(list.Buckets))

	for i, b := range list.Buckets {
		buckets[i] = *b.Name
	}

	return buckets, nil
}

func (o *S3) PutObject(reader io.ReadSeeker, key string) error {

	ret, err := o.svc.PutObject(&s3.PutObjectInput{
		Body:   reader,
		Bucket: &o.bucket,
		Key:    &key,
	})
	if err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Println(ret)

	return nil
}

func (o *S3) DeleteObject(key string) error {

	_, err := o.svc.DeleteObject(&s3.DeleteObjectInput{
		Bucket: &o.bucket,
		Key:    &key,
	})
	if err != nil {
		fmt.Println(err)
		return err
	}

	//fmt.Println(ret)

	return nil
}

func (o *S3) GetObject(key string) (*s3.GetObjectOutput, error) {

	ret, err := o.svc.GetObject(&s3.GetObjectInput{
		Bucket: &o.bucket,
		Key:    &key,
	})
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func (o *S3) HeadObject(key string) (*s3.HeadObjectOutput, error) {

	ret, err := o.svc.HeadObject(&s3.HeadObjectInput{
		Bucket: &o.bucket,
		Key:    &key,
	})
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func (o *S3) PutObjectAcl(key string, acl string) (*s3.PutObjectAclOutput, error) {

	ret, err := o.svc.PutObjectAcl(&s3.PutObjectAclInput{
		Bucket: &o.bucket,
		Key:    &key,
		ACL:    &acl,
	})
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func (o *S3) Upload(reader io.Reader, key string, acl string) (string, error) {
	uploader := s3manager.NewUploader(o.session)

	ret, err := uploader.Upload(&s3manager.UploadInput{
		Body:   reader,
		Bucket: &o.bucket,
		Key:    &key,
		ACL:    &acl,
	})
	if err != nil {
		return "", err
	}

	return ret.Location, nil
}

func (o *S3) Download(w io.WriterAt, key string) (int64, error) {
	downloader := s3manager.NewDownloader(o.session)

	numBytes, err := downloader.Download(w, &s3.GetObjectInput{
		Bucket: &o.bucket,
		Key:    &key,
	})
	if err != nil {
		return 0, err
	}

	return numBytes, nil
}
