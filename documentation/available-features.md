# Available internal features
Here you can see the list potential useful features of the project, which can be useful during your custom events build.

## Table of contents
- [BitBucket API client](#bitbucket-api-client)
- [Slack API client](#slack-api-client)
- [Http client](#http-client)
- [Migrations service](migrations.md)
- [Available helper functions](#available-helper-functions)
- [Logger](#logger)
- [Scenarios](scenarios.md)
- [Database query builder](query-builder.md)

## BitBucket API client
The client, which can be used for custom requests to the BitBucket API. The good examples of usage of that client are [bb release](https://github.com/sharovik/bitbucket-release-event) and [start pipeline](https://github.com/sharovik/bitbucket-run-pipeline) events. Feel free to check the source of these events.

### Source
You can find the source here: `internal/client/bitbucket.go`

### Basic example
Imagine you have your own custom event. This bitbucket API client is injected by default into container. To use the features of the BitBucket client you can do the next steps:
1. select the method from available public methods
2. use selected method in your custom event
```
//... some code of your event here
accessToken, err := container.C.BibBucketClient.GetAuthToken()
if err != nil {
	return
}
//... now you have access to the access token of bitbucket API
```

## Slack API client
With this client you can send the messages to a specific channel/user, attache files to a specific channel or do a custom requests to the slack events API

### Source
You can find the source here: `internal/client/slack.go`

### Basic example
Let's imagine you have your own custom event, where you need to send the message to a specific user. 
``` 
//... some code of your event here
response, statusCode, err := container.C.MessageClient.SendMessage(dto.SlackRequestChatPostMessage{
    Channel:           "SLACK_UID",
    Text:              "This is the test message",
    AsUser:            true,
    Ts:                time.Time{},
    DictionaryMessage: dto.DictionaryMessage{},
    OriginalMessage:   dto.SlackResponseEventMessage{},
})
if err != nil {
    log.Logger().AddError(err).
        Interface("response", response).
        Interface("status", statusCode).
        Msg("Failed to sent answer message")
}
//If there was no error, that means the message was sent to the "SLACK_UID"
```

## Http Client
This is a simple http client, which can be used for requests to your custom API gateways.

### Source
You can find the source here: `internal/client/http.go`

### Basic example
With the Http client you can trigger your custom API endpoints. The client supports next type of the requests: `GET`, `POST`, `PUT` and of course there is always possibility to use your custom request type via `Request` method.

Request via available `GET`, `POST`, `PUT` methods.
```
//... some custom code of your event
container.C.HTTPClient.SetBaseURL({YOUR_BASE_URL})

//We define the form which we need to sent to the API endpoint
form := url.Values{}
form.Add("some_key", "this is some key value")
form.Add("some_other_key", "this is some other key value")

//Actual API POST request using form request
response, statusCode, err := container.C.HTTPClient.Post({YOUR_ENDPOINT_HERE}, form, map[string]string{})
if err != nil {
    log.Logger().AddError(err).
        Str("endpoint", {YOUR_ENDPOINT_HERE}).
        Int("status_code", statusCode).
        RawJSON("response", response).
        Msg("Error received during API request")
}
```

Or if your API supports the json requests
```
//... some custom code of your event
type Request struct {
	Name   string                `json:"name"`
}

var request := Request{
    Name: "test",
}

byteRequest, err := json.Marshal(request)
if err != nil {
    log.Logger().AddError(err).
        Interface("request", request).
        Msg("Error during marshal of the request")
}

container.C.HTTPClient.SetBaseURL({YOUR_BASE_URL})
response, statusCode, err := b.client.Post({YOUR_ENDPOINT_HERE}, byteRequest, map[string]string{})
if err != nil {
    log.Logger().AddError(err).
        Int("status_code", statusCode).
        RawJSON("response", response).
        Msg("The request failed")
}
```

Now let's try to make DELETE request via `Request` method
``` 
//... some custom code of your event

//In that case you don't need to define the base url 
//container.C.HTTPClient.SetBaseURL({YOUR_BASE_URL})
//so you can write the full endpoint url including the base url
response, statusCode, err := client.Request(http.MethodDelete, {YOUR_FULL_API_ENDPOINT_URL}, interface{}, map[string]string{})
if err != nil {
    log.Logger().AddError(err).
        Int("status_code", statusCode).
        RawJSON("response", response).
        Msg("The request failed")
}
```

The `Request` method currently supports next types of `body` parameter:
1. `string` - when we need to send a simple string to selected endpoint
2. `url.Values` - when we need to send the form-data
3. `[]byte` - when we need to send the json requests

## Available helper functions
These functions can be used in different cases and can help you. All these functions are public and available in `github.com/sharovik/devbot/internal/helper` package

### Source
You can find the source here: `internal/helper/helper.go`

### Basic example
``` 
import "github.com/sharovik/devbot/internal/helper"

//... some custom code of your event
matches := helper.FindMatches(`(?i)(help)`, "Can you find the help word in this string?")
if len(matches) != 0 {
    fmt.Println("Yes I can.")
}
```

## Logger
The logger is based on the zerolog library, which you can find [here](https://github.com/rs/zerolog).

There is a several method which can be using for the logging actions in your custom event.
### Source
You can find the source here: `internal/log/logger.go`

### Basic example
This is the simple usage of the logger
``` 
//... some code of your event
log.Logger().Info().Str("some_field", "The value").Msg("This is the test log")
```
But once you need to have a global context in your project, you can use `AppendGlobalContext` method to append the values to a global context. After this log trigger, these values will be available in `context` attribute of all logs.
``` 
//... some code of your event
numberOfRetries := 1
log.Logger().AppendGlobalContext(map[string]interface{}{
    "number_retries": numberOfRetries,
})
```
