package secexport

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/aws/aws-sdk-go/service/ssm"
)

type secretsData struct {
	Value string
	Key   string
	Arn   string
}

type secretNormalized struct {
	Value string
	Arn   string
}

type AWSSecrets struct {
	Data map[string]secretNormalized
}

func (s *AWSSecrets) Values() *string {
	var res string

	for k, v := range s.Data {
		var buf bytes.Buffer
		buf.WriteString("export ")
		buf.WriteString(strings.ToUpper(strings.ReplaceAll(strings.ReplaceAll(k, "-", "_"), "/", "_")))
		buf.WriteString("='")
		buf.WriteString(strings.ReplaceAll(v.Value, "\"", ""))
		buf.WriteString("'\n")

		res += buf.String()
		buf.Reset()
	}

	return &res
}

func (s *AWSSecrets) Parse(d []byte) (*string, error) {
	var normalized map[string]secretNormalized

	err := json.Unmarshal(d, &normalized)
	if err != nil {
		return nil, err
	}

	s.Data = normalized
	return s.Values(), nil
}

func RetreiveSecrets(args []*string, pm bool, sc bool) (*AWSSecrets, error) {
	var err error

	sess, err := startSession()
	if err != nil {
		log.Printf("Got error while creating session: %v", err)
		return nil, err
	}
	secrets := make([]*secretsData, 0)

	if sc {
		secrets, err = secretsManager(sess, args)
		if err != nil {
			log.Printf("Got error while processing secrets from the secrets manager. \n%v", err)
			return nil, err
		}
	}

	parameters := make([]*secretsData, 0)
	if pm {
		parameters, err = parameterStore(sess, args)
		if err != nil {
			log.Printf("Got error while processing secrets from the parameter store. \n%v", err)
			return nil, err
		}
	}

	normalized, err := normalize(append(secrets, parameters...))
	if err != nil {
		log.Printf("Got error while normilizing secrets. \n%v", err)
		return nil, err
	}

	return &AWSSecrets{Data: normalized}, nil
}

func normalize(secrets []*secretsData) (map[string]secretNormalized, error) {
	if len(secrets) == 0 {
		return nil, nil
	}

	result := make(map[string]secretNormalized)

	for _, secret := range secrets {
		transformed := transform(secret)
		for key, value := range transformed {
			if _, exists := result[key]; exists {
				key = fmt.Sprintf("%s_%s", key, GetSHA1(&secret.Arn))
			}
			result[key] = value
		}
	}

	return result, nil
}

func startSession() (*session.Session, error) {
	s, err := session.NewSession(&aws.Config{Region: aws.String("us-east-1")})
	if err != nil {
		return nil, err
	}

	return s, nil
}

// Retreives secrets from the AWS SEcrets Manager.
func secretsManager(sess *session.Session, args []*string) ([]*secretsData, error) {
	svc := secretsmanager.New(sess)

	filters := []*secretsmanager.Filter{{Key: aws.String("all"), Values: args}}

	input := &secretsmanager.ListSecretsInput{Filters: filters}

	result, err := svc.ListSecrets(input)
	if err != nil {
		return nil, err
	}

	secrets := make([]*secretsData, 0, len(result.SecretList))

	if len(result.SecretList) <= 0 {
		return secrets, nil
	}

	for _, v := range result.SecretList {
		secret := &secretsData{}

		input := &secretsmanager.GetSecretValueInput{SecretId: aws.String(*v.ARN)}

		res, err := svc.GetSecretValue(input)
		if err != nil {
			return nil, err
		}

		secret.Arn = *v.ARN
		secret.Value = *res.SecretString
		secret.Key = *v.Name

		secrets = append(secrets, secret)
	}

	return secrets, nil
}

// Retreives secrets from the AWS Parameter Store
func parameterStore(sess *session.Session, args []*string) ([]*secretsData, error) {
	pm := ssm.New(sess)
	describeInput := &ssm.DescribeParametersInput{
		ParameterFilters: []*ssm.ParameterStringFilter{
			{
				Key:    aws.String("Name"),
				Option: aws.String("Equals"),
				Values: args,
			},
		},
	}

	parameters, err := pm.DescribeParameters(describeInput)
	if err != nil {
		return nil, err
	}

	names := make([]string, 0, len(parameters.Parameters))
	secrets := make([]*secretsData, 0, len(names))

	if len(parameters.Parameters) <= 0 {
		return secrets, nil
	}

	for _, p := range parameters.Parameters {
		names = append(names, *p.Name)
	}

	input := &ssm.GetParametersInput{
		Names: aws.StringSlice(names),
	}

	result, err := pm.GetParameters(input)
	if err != nil {
		return nil, err
	}

	if len(result.Parameters) <= 0 {
		return secrets, nil
	}

	for _, v := range result.Parameters {
		secret := &secretsData{
			Arn:   *v.ARN,
			Value: *v.Value,
			Key:   *v.Name,
		}

		secrets = append(secrets, secret)
	}

	return secrets, nil
}

func transform(secret *secretsData) map[string]secretNormalized {
	result := make(map[string]secretNormalized)

	if IsJSON(&secret.Value) {
		var data interface{}
		err := json.Unmarshal([]byte(secret.Value), &data)
		if err != nil {
			log.Printf("Error while unmarshaling JSON string: %v", err)
			return result
		}

		flatten("", data, secret.Arn, result)

	} else {
		result[secret.Key] = secretNormalized{Value: secret.Value, Arn: secret.Arn}
	}

	return result
}

func flatten(prefix string, value interface{}, arn string, result map[string]secretNormalized) {
	switch v := value.(type) {
	case map[string]interface{}:
		for k, val := range v {
			newKey := k
			if prefix != "" {
				newKey = prefix + "." + k
			}
			flatten(newKey, val, arn, result)
		}
	case []interface{}:
		jsonValue, _ := json.Marshal(v)
		result[prefix] = secretNormalized{Value: string(jsonValue), Arn: arn}
	default:
		jsonValue, _ := json.Marshal(v)
		result[prefix] = secretNormalized{Value: string(jsonValue), Arn: arn}
	}
}
