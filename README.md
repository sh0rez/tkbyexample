# [Tanka by Example](https://tkbyexample.netlify.com)

Tanka by Example is a hands-on introduction to [Tanka](https://tanka.dev) and
[Jsonnet](https://jsonnet.org) using annotated example programs, inspired by the
popular [gobyexample.com](https://gobyexample.com).

## Creating examples

Example source is located in `src/examples/<example_path>`. To create a new one,
make a new folder in that directory and add a `.x.yml` file to it:

```yaml
title: Hello World # <title>, also shown in index and <h1>
description: A warm welcome to tkbyexample # <meta name=description>
```

To write the actual example, add a `main.jsonnet` file and annotate it using
comments.

## Running locally

This page consists of two parts:

1. Golang based example generator, ported from https://gobyexample.com
2. GatbsyJS based static site using React, rendering the markdown generated by Go

For development purposes (file watching, etc):

```bash
$ make dev
```

This will run the generator in watching mode (`go run ./gen dev`) and `gatsby develop` for you.

## Publishing

The site is built and published by Netlify. Pushing to the `master` branch will
automatically update the page.
