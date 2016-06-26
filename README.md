[![Build Status](https://travis-ci.org/ryanmoran/piper.svg?branch=master)](https://travis-ci.org/ryanmoran/piper)

# piper
Like `fly execute` but running against a local docker daemon

## But why?
Have you ever run `fly execute` with a set of huge inputs
that take FOREVER to upload only to find out you had a typo
in your task script? This is why.

## But how?
`piper` is really just a wrapper around some called to the
`docker` cli. It pulls the image needed to run the task.
It then runs `docker run` with some volume mounts for your
inputs and outputs. Its as simple as that.

## Installation
`go get github.com/ryanmoran/piper/piper`
