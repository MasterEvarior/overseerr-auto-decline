# Overseerr Auto Decline
![quality workflow](https://github.com/MasterEvarior/overseerr-auto-decline/actions/workflows/quality.yaml/badge.svg) ![release workflow](https://github.com/MasterEvarior/overseerr-auto-decline/actions/workflows/publish.yaml/badge.svg)

[Overseerr](https://overseerr.dev/) is a fantastic application, which allows your friends to request movies and TV series for your [Plex Server](https://www.plex.tv/). It is however not possible to create a deny-list of movies and TV series, requests of which are automatically declined. 

With this application this is now possible. Download the container. configure the media you ~~absolutely hate~~ you want to automatically decline and voila: never tell your friends "no" manually again.

## Build
To build the container yourself, simply clone the repository and then build the container with the provided docker file. You can the run it as described in the section below.
```shell
docker build . --tag overseerr-auto-decline
```
Alternatively you can build the binary directly with Go.
```shell
go build -o ./overseerr-auto-decline
```

## Run
To run the docker container, you have to give it a couple of environment variables.
```shell
docker run -d \
  -e URL=https://your.overseerr.com \
  -e API_KEY=eW91cl9hcGlfa2V5Cg== \
  -e MEDIA=8966,24021 \
  --name overseerr-auto-decline \
  -p 8080:8080 \
  ghcr.io/masterevarior/overseerr-auto-decline:latest
```
This will decline (but not delete) any request for media with the id `8966` or `24021`.

### Environment Variables
| Name            | Description                                                                                                                                                                                             | Example                    | Mandatory |
|-----------------|---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|----------------------------|-----------|
| URL             | The URL to your Overseerr instance.                                                                                                                                                                     | https://your.overseerr.com | ✅        |
| API_KEY         | Your API-Key to your Overseerr instance.                                                                                                                                                                | eW91cl9hcGlfa2V5Cg==       | ✅        |
| MEDIA           | A list of comma separated [TMDB](https://www.themoviedb.org/) or [TVDB](https://thetvdb.com/) id. Any movies or series that is included here will be declined and, depending on your settings, deleted. | 8966,24021                 | ✅        |
| DELETE_REQUESTS | Wether the requests should not only be declined but also be deleted. If this variable is set, they will also be deleted.                                                                                | true                       | ❌        |

## Configure Webhook
For this application to work, you will have to configure [the webhook](https://docs.overseerr.dev/using-overseerr/notifications/webhooks) inside of Overseerr.

Point the URL to whereever your Docker container is running. Choose `Request Pending Approval` as notification type. Do not choose any other notification type. Finally add this as your JSON payload:
```json
{
    "request_id": "{{request_id}}",
    "tmdbid": "{{media_tmdbid}}",
    "tvdbid": "{{media_tvdbid}}"
}
```

## Development

### Linting
[golangci-lint](https://golangci-lint.run/) is used for linting. Note that the command does not install the required formatters.

```shell
golangci-lint run
```

### Git Hooks

There are some hooks for formatting and the like. To use those, execute the following command:

```shell
git config --local core.hooksPath .githooks/
```

## Improvements, issues and more
Pull requests, improvements and issues are always welcome.