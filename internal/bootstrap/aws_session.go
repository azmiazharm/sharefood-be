// Package bootstrap
// @author Daud Valentino
package bootstrap

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"

	"sharefood/internal/appctx"
	"sharefood/pkg/logger"
)

// RegistryAWSSession initialize aws session
func RegistryAWSSession(appCtx *appctx.Config) *session.Session {
	var (
		awsConfig aws.Config
	)

	awsConfig.CredentialsChainVerboseErrors = aws.Bool(true)
	region := appCtx.AWS.Region
	if len(region) != 0 {
		awsConfig.Region = aws.String(region)
	}

	accessKeyID := appCtx.AWS.AccessKey
	secretAccessKey := appCtx.AWS.AccessSecret
	if len(accessKeyID) != 0 && len(secretAccessKey) != 0 {
		awsConfig.Credentials = credentials.NewStaticCredentials(accessKeyID, secretAccessKey, "")
	}

	sess, err := session.NewSession(&awsConfig)

	if err != nil {
		logger.Fatal(err, logger.EventName("aws-session"))
	}

	return sess
}
