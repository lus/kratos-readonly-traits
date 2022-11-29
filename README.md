# kratos-readonly-traits

`kratos-readonly-traits` is a simple service simplifying the implementation of read-only traits in
[Ory Kratos](https://github.com/ory/kratos).

It works by exposing an endpoint called by Kratos as a blocking web hook during the settings flow.

## Installation & Configuration

### Docker

A Docker image is automatically built and pushed to GHCR.
A `docker-compose.yml` file may look like this:

```yml
services:
    kratos-readonly-traits:
        image: ghcr.io/lus/kratos-readonly-traits:latest
        restart: unless-stopped
        ports:
            # I do not recommend actually exposing this port to the internet.
            # A better solution would be to put this service in a mutual Docker network with Kratos.
            - 8080:8080
        environment:
            # By default the application runs in development mode.
            # I highly recommend setting this value for production usage.
            ENVIRONMENT: prod
            # The minimum log level to print.
            # Available values: 'trace', 'debug', 'info', 'warn', 'error', 'fatal', 'panic' and 'disabled'
            # The default value is 'info'. I recommend leaving this as this will not clutter your console anyway.
            LOG_LEVEL: info
            # The address to bind the HTTP server to.
            # The default value is :8080 (= 0.0.0.0:8080).
            LISTEN_ADDRESS: :8080
            # The message that will appear below the input field when trying to manipulate a read-only trait.
            # The default value is 'This field is read-only.'.
            ERROR_MESSAGE: This field is read-only.
```

### From source

After making sure that [Go 1.19](https://go.dev/dl/) is installed, simply follow these steps:

1. Clone the repository and enter the directory it got cloned to:
    ```shell
    git clone https://github.com/lus/kratos-readonly-traits && cd kratos-readonly-traits
    ```
2. Build the binary:
    ```shell
    go build -o server cmd/server/main.go
    ```
3. Make sure to set the environment variables according to your configuration. A `.env` file is supported natively:
    ```
    # By default the application runs in development mode.
    # I highly recommend setting this value for production usage.
    ENVIRONMENT: prod
   
    # The minimum log level to print.
    # Available values: 'trace', 'debug', 'info', 'warn', 'error', 'fatal', 'panic' and 'disabled'
    # The default value is 'info'. I recommend leaving this as this will not clutter your console anyway.
    LOG_LEVEL: info
   
    # The address to bind the HTTP server to.
    # The default value is :8080 (= 0.0.0.0:8080).
    LISTEN_ADDRESS: :8080
   
    # The message that will appear below the input field when trying to manipulate a read-only trait.
    # The default value is 'This field is read-only.'.
    ERROR_MESSAGE: This field is read-only.
    ```
4. Run the application:
    ```shell
    ./server
    ```

## Getting started

### Configure Kratos webhook

In order for this service to be able to interrupt the settings flow of Kratos, you need to configure a blocking webhook.

1. Please create a `<something>.jsonnet` file somewhere where Kratos can access it. It must look like this:
    ```jsonnet
    function(ctx) {
        schema_url: ctx.identity.schema_url,
        old_traits: ctx.flow.identity.traits,
        new_traits: ctx.identity.traits
    }
    ```
2. Adjust the configuration file of Kratos so that it contains a webhook configuration like this:
    ```yml
    selfservice:
      flows:
        settings:
          after:
            profile:
              hooks:
                - hook: web_hook
                  config:
                    url: http://<kratos-readonly-traits>:8080
                    method: POST
                    body: file:///path/to/<something>.jsonnet
                    can_interrupt: true
    ```
   
### Adjust the identity schema

This service will not treat any trait as read-only by default.
To configure a read-only trait, simply add the `lus/kratos-readonly-traits.readonly` boolean to the identity schema:

```json
{
    "$schema": "http://json-schema.org/draft-07/schema#",
    "$id": "https://example.com/schemas/user.schema.json",
    "title": "User",
    "type": "object",
    "properties": {
        "traits": {
            "type": "object",
            "properties": {
                "username": {
                    "title": "Username",
                    "type": "string",
                    "minLength": 1,
                    "maxLength": 20,
                    "pattern": "^[a-z0-9]+$",
                    "ory.sh/kratos": {
                        "credentials": {
                            "password": {
                                "identifier": true
                            }
                        }
                    },
                    "lus/kratos-readonly-traits": {
                        "readonly": true
                    }
                },
                "email": {
                    "title": "E-Mail",
                    "type": "string",
                    "format": "email",
                    "ory.sh/kratos": {
                        "credentials": {
                            "password": {
                                "identifier": true
                            }
                        },
                        "recovery": {
                            "via": "email"
                        },
                        "verification": {
                            "via": "email"
                        }
                    }
                }
            },
            "required": [
                "username",
                "email"
            ],
            "additionalProperties": false
        }
    }
}
```

This schema will make the `username` trait read-only.

## Support

Feel free to open issues in this repository if you encounter any problem or want to suggest a feature.
If you want to ask a quick question, feel free to join my [Discord server](https://go.lus.pm/discord).
