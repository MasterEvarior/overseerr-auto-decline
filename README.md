# I don't want your movie
![test workflow](https://github.com/MasterEvarior/i-dont-want-your-movie/actions/workflows/test.yaml/badge.svg) ![release workflow](https://github.com/MasterEvarior/i-dont-want-your-movie/actions/workflows/release.yaml/badge.svg)



## Build
To build the container yourself, simply clone the repository and then build the container with the provided docker file. You can the run it as described in the section below.
```shell
docker build . --tag i-dont-want-your-movie
```
Alternatively you can build the binary directly with Go.
```shell
go build -o ./i-dont-want-your-movie
```

## Run
To run the docker container, you have to give it a couple of environment variables.
```shell
docker run -d \
  -e URL=https://your.overseerr.com \
  -e API_KEY=eW91cl9hcGlfa2V5Cg== \
  -e MEDIA=8966,24021 \
  --name i-dont-want-your-movie \
  ghcr.io/masterevarior/i-dont-want-your-movie:latest
```
This will decline (but not delete) any request for media with the id `8966` or `24021`.

### Environment Variables
| Name            | Description                                                                                                                                                                                             | Example                    | Mandatory |
|-----------------|---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|----------------------------|-----------|
| URL             | The URL to your Overseerr instance.                                                                                                                                                                     | https://your.overseerr.com | ✅        |
| API_KEY         | Your API-Key to your Overseerr instance.                                                                                                                                                                | eW91cl9hcGlfa2V5Cg==       | ✅        |
| MEDIA           | A list of comma separated [TMDB](https://www.themoviedb.org/) or [TVDB](https://thetvdb.com/) id. Any movies or series that is included here will be declined and, depending on your settings, deleted. | 8966,24021                 | ✅        |
| DELETE_REQUESTS | Wether the requests should not only be declined but also be deleted. If this variable is set, they will also be deleted.                                                                                | true                       | ❌        |

## Development, improvements and more
Pull requests, improvements and issues are always welcome.