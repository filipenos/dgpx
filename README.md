# dgpx
Challenge for dgpx

[![Build Status](https://travis-ci.org/filipenos/dgpx.svg?branch=master)](https://travis-ci.org/filipenos/dgpx)

    go get -u github.com/filipenos/dgpx
    dgpx -h
    Usage of dgpx:
      -port int
        	port to listen (default 8080)
      -token string
        	token used in authorization

    docker run --rm --publish 8080:8080 -e TOKEN=codigo_do_token filipenos/dgpx:1.0
