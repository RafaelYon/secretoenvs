# secretoenvs
secretoenvs (pronounced "secret to envs") is a utility to convert an Environment Variable stored as JSON to key=value format (shell environment variables).

## Motivation
In Amazon Web Service (AWS) Elastic Container Service (ECS) you can inject secret values stored in AWS Secrets Manager by specifying the following information in the container definition:
  - Accessible key name for the container
  - Amazon Resource Name from Secret or from a value stored in a given Secret

Example:
```json
{
  "containerDefinitions": [{
    "secrets": [{
      "name": "environment_variable_name",
      "valueFrom": "arn:aws:secretsmanager:region:aws_account_id:secret:appauthexample-AbCdEf:username1::"
    }]
  }]
}
```
> [Specifying sensitive data using Secrets Manager](https://docs.aws.amazon.com/AmazonECS/latest/developerguide/specifying-sensitive-data-secrets.html)


However it **is not possible to reference the same ARN of a Secret for different values** So the **only way to inject more than one value stored in the same Secrets in ECS is to inject all Secret values in a single ENV of the container in JSON format**.

Therefore, **the application needs to pre-process the ENV** (in JSON format) before starting to use the relevant values for it. And for applications written in languages like PHP this pre-processing is **repeated each time the script is executed**.

So, I would like to process this ENV with all Secrets values (decode the JSON) just once, preferably at the beginning of the container execution. When researching a bit, I didn't find any ready-made solution, so I decided to create this utility to make a solution that would meet my needs.

With this utility it is possible to create the following script for the container's entrpoint:
```bash
#!/bin/sh
set -e

# Export AWS Secret Manager Values stored in AWS_SECRETS_VALUES ENV
for s in $(secretoenv AWS_SECRETS_VALUES); do
    export $s
done

exec "php-fpm $@"
```

With this, **we pre-process the ENV only when the container starts** and the values we are interested in are exposed to all PHP-FPM child processes and our sensitive values are not versioned in the docker image.

## Use
This utility has a basic explanation of the arguments and flags that can be specified:

```
$ secretoenvs -h
```
Output:
```
Usage of secretoenvs [RAW_JSON or ENVIRONMENT_VARIATION_NAME]:
  -key-prefix string
        Specifies a prefix for each key in the output.
  -quotation-marks string
        Specifies a combination of characters to be placed before and after the ENV value in the output.
  -raw-value
        Specifies that JSON was passed as an argument. So it is not necessary to specify an ENV as input.
```

### Example of use
Assumed there is the following ENV with a JSON in your environment:
```
$ echo $AWS_SECRETS_JSON
```
Output:
```
{"API_HOST":"https://example.com/api","API_KEY":"some secrets value"}
```

#### Basic
We can convert this JSON to key=value notation with the following statement:
```
$ secretoenvs AWS_SECRETS_JSON
```
Output:
```
API_HOST=https://example.com/api
API_KEY=some secrets value
```

#### Adding quotes to avoid interpreting values
Quotes can be specified by adding the following `-quotation-mark=\'`" flag in front of positional arguments:
```
$ .secretoenvs -quotation-marks=\' AWS_SECRETS_JSON
```
Output:
```
API_HOST='https://example.com/api'
API_KEY='some secrets value'
```

#### Added prefix to avoid name conflicts with other ENVs
Prefix in the name of ENVs can be specified by adding the `-key-prefix=SOME_PREFIX` flag in front of positional arguments:
```
$ .secretoenvs -key-prefix=AWS_ AWS_SECRETS_JSON
```
Output:
```
AWS_API_HOST=https://example.com/api
AWS_API_KEY=some secrets value
```

#### Passing the JSON without needing an intermediate ENV
You can specify a JSON directly:
```
$ .secretoenvs -raw-value '{"API_HOST":"https://example.com/api","API_KEY":"some secrets value"}'
```
Output:
```
API_HOST=https://example.com/api
API_KEY=some secrets value
```
