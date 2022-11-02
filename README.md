# secretoenvs
secretoenvs (pronounced "secret to envs") is a utility to convert an Environment Variable stored as JSON to key=value format (shell environment variables).

## Contents
 - [Motivation](#motivation)
    - [Use cases](#use-cases)
      - [Retrieving AWS Secrets Manager values and generating a dotenv file](#retrieving-aws-secrets-manager-values-and-generating-a-dotenv-file)
      - [Retrieving AWS Secrets Manager values and generating a configuration file for php-fpm](#retrieving-aws-secrets-manager-values-and-generating-a-configuration-file-for-php-fpm)
 - [Installing](#installing)
 - [Use](#use)
    - [Example of use](#example-of-use)
        - [Basic](#basic)
        - [Adding quotes to avoid interpreting values](#adding-quotes-to-avoid-interpreting-values)
        - [Added prefix to avoid name conflicts with other ENVs](#added-prefix-to-avoid-name-conflicts-with-other-envs)
        - [Changing the key value separator](#changing-the-key-value-separator)
        - [Passing the JSON without needing an intermediate ENV](#passing-the-json-without-needing-an-intermediate-env)

## Motivation
AWS Secrets Manager lets you store and retrieve sensitive values. Many AWS services can already interact with it natively, **however using this service with existing daemons can be a little unpleasant, as these daemons often require a file format other than JSON for their settings.**

In this way, when it is not possible to use the values retrieved from AWS Secrets in JSON format (standard format of the API response of this service), it is necessary to create string manipulation commands with sed, awk and others to arrive at the desired format.

So I decided to create a small utility (less than 100 lines currently) to make it easy to **convert the secrets in JSON format to formats that require each Secret Value on a separate line** (like dotenv format).

### Use cases

#### Retrieving AWS Secrets Manager values and generating a dotenv file

We can retrieve the secret values and generate a dotenv file with the following command:
```sh
secretoenvs -raw-value "$(aws secretsmanager get-secret-value --secret-id arn:aws:secretsmanager:{region}:{account_id}:secret:{secret_id} | jq '.SecretString | fromjson')" > .env
```

#### Retrieving AWS Secrets Manager values and generating a configuration file for php-fpm

We can retrieve the secret values and generate a configuration file for php-fpm with the following command:
```sh
secretoenvs -key-prefix='env[' -key-value-separator='] = ' -raw-value "$(aws secretsmanager get-secret-value --secret-id arn:aws:secretsmanager:{region}:{account_id}:secret:{secret_id} | jq '.SecretString | fromjson')" > /usr/local/etc/php-fpm.d/secrets.conf
```

## Installing
There are two ways to install:

Through `go install`
```sh
go install github.com/RafaelYon/secretoenvs@0.2.0
```

Or downloading a [pre-compiled version](https://github.com/RafaelYon/secretoenvs/releases/tag/0.1.0):
```sh
wget https://github.com/RafaelYon/secretoenvs/releases/download/0.2.0/secretoenvs_linux_amd64 -O ~/bin/secretoenvs
```

After getting a binary you can start using it.

## Use
This utility has a basic explanation of the arguments and flags that can be specified:

```sh
secretoenvs -h
```
Output:
```
Usage of secretoenvs [RAW_JSON or ENVIRONMENT_VARIATION_NAME]:
  -key-prefix string
        Specifies a prefix for each key in the output.
  -key-value-separator string
        Specifies the characters to use to separate key values (default "=")
  -quotation-marks string
        Specifies a combination of characters to be placed before and after the ENV value in the output.
  -raw-value
        Specifies that JSON was passed as an argument. So it is not necessary to specify an ENV as input.
```

### Example of use
Assumed there is the following ENV with a JSON in your environment:
```sh
echo $AWS_SECRETS_JSON
```
Output:
```
{"API_HOST":"https://example.com/api","API_KEY":"some secrets value"}
```

#### Basic
We can convert this JSON to key=value notation with the following statement:
```sh
secretoenvs AWS_SECRETS_JSON
```
Output:
```
API_HOST=https://example.com/api
API_KEY=some secrets value
```

#### Adding quotes to avoid interpreting values
Quotes can be specified by adding the following `-quotation-mark=\'`" flag in front of positional arguments:
```sh
secretoenvs -quotation-marks=\' AWS_SECRETS_JSON
```
Output:
```
API_HOST='https://example.com/api'
API_KEY='some secrets value'
```

#### Added prefix to avoid name conflicts with other ENVs
Prefix in the name of ENVs can be specified by adding the `-key-prefix=SOME_PREFIX` flag in front of positional arguments:
```sh
secretoenvs -key-prefix=AWS_ AWS_SECRETS_JSON
```
Output:
```
AWS_API_HOST=https://example.com/api
AWS_API_KEY=some secrets value
```

#### Changing the key value separator
The characters used to separate the key from the value can be changed by adding the `-key-value-separator=SOME_SEPARATOR` flag in front of the positional arguments:
```sh
secretoenvs -key-value-separator=': ' AWS_SECRETS_JSON
```
Output:
```
API_HOST: https://example.com/api
API_KEY: some secrets value
```

#### Passing the JSON without needing an intermediate ENV
You can specify a JSON directly:
```sh
.secretoenvs -raw-value '{"API_HOST":"https://example.com/api","API_KEY":"some secrets value"}'
```
Output:
```
API_HOST=https://example.com/api
API_KEY=some secrets value
```
