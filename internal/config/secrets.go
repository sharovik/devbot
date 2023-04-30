package config

// Use this code snippet in your app.
// If you need more information about configurations or implementing the sample code, visit the AWS docs:
// https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/setting-up.html
import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
)

// SecretConfigValues the secret values object
type SecretConfigValues struct {
	//MessagesAPIOAuthToken slack oauth token
	MessagesAPIOAuthToken string `json:"MESSAGES_API_OAUTH_TOKEN"`

	//MessagesAPIWebAPIOAuthToken slack web-oauth token
	MessagesAPIWebAPIOAuthToken string `json:"MESSAGES_API_WEB_API_OAUTH_TOKEN"`

	//BitBucketClientID the client id for bitbucket api
	BitBucketClientID string `json:"BITBUCKET_CLIENT_ID"`

	//BitBucketClientID the client id for bitbucket api
	BitBucketClientSecret string `json:"BITBUCKET_CLIENT_SECRET"`

	//GoogleClientID the client id for a Google oauth2
	GoogleClientID string `json:"DEVBOT_GOOGLE_CLIENT_ID"`

	//GoogleClientSecret the client secret for Google oauth2
	GoogleClientSecret string `json:"DEVBOT_GOOGLE_CLIENT_SECRET"`
}

// GetSecret method retrieves the secrets from the vault
func GetSecret(secretName string, region string) (secrets SecretConfigValues, err error) {
	awsSession, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})

	if err != nil {
		return SecretConfigValues{}, err
	}

	//Create a Secrets Manager client
	svc := secretsmanager.New(awsSession, aws.NewConfig().WithRegion(region))
	input := &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(secretName),
		VersionStage: aws.String("AWSCURRENT"), // VersionStage defaults to AWSCURRENT if unspecified
	}
	// In this sample we only handle the specific exceptions for the 'GetSecretValue' API.
	// See https://docs.aws.amazon.com/secretsmanager/latest/apireference/API_GetSecretValue.html
	result, err := svc.GetSecretValue(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case secretsmanager.ErrCodeDecryptionFailure:
				// Secrets Manager can't decrypt the protected secret text using the provided KMS key.
				fmt.Println(secretsmanager.ErrCodeDecryptionFailure, aerr.Error())
			case secretsmanager.ErrCodeInternalServiceError:
				// An error occurred on the server side.
				fmt.Println(secretsmanager.ErrCodeInternalServiceError, aerr.Error())
			case secretsmanager.ErrCodeInvalidParameterException:
				// You provided an invalid value for a parameter.
				fmt.Println(secretsmanager.ErrCodeInvalidParameterException, aerr.Error())
			case secretsmanager.ErrCodeInvalidRequestException:
				// You provided a parameter value that is not valid for the current state of the resource.
				fmt.Println(secretsmanager.ErrCodeInvalidRequestException, aerr.Error())
			case secretsmanager.ErrCodeResourceNotFoundException:
				// We can't find the resource that you asked for.
				fmt.Println(secretsmanager.ErrCodeResourceNotFoundException, aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return
	}
	// Decrypts secret using the associated KMS CMK.
	// Depending on whether the secret is a string or binary, one of these fields will be populated.
	if result.SecretString == nil {
		fmt.Println("The secrets string is empty")
		return
	}

	if err = json.Unmarshal([]byte(*result.SecretString), &secrets); err != nil {
		fmt.Println("Failed to unmarshal the secrets values string")
		return
	}

	return secrets, err
}
