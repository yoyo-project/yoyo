# yoyo

[![MIT license](https://img.shields.io/badge/license-MIT-brightgreen.svg)](https://opensource.org/licenses/MIT)
[![Go Report Card](https://goreportcard.com/badge/github.com/yoyo-project/yoyo)](https://goreportcard.com/report/github.com/yoyo-project/yoyo)
[![Maintainability](https://api.codeclimate.com/v1/badges/1e0d4f34de5f07425ba5/maintainability)](https://codeclimate.com/github/yoyo-project/yoyo/maintainability)
[![codecov](https://codecov.io/gh/yoyo-project/yoyo/branch/main/graph/badge.svg)](https://codecov.io/gh/yoyo-project/yoyo)
[![CircleCI](https://circleci.com/gh/yoyo-project/yoyo/tree/main.svg?style=shield)](https://circleci.com/gh/yoyo-project/yoyo/tree/main)

A Migration Generator and Database Access Layer Generator for Go projects. Made with ❤️ to hopefully make your life a
little bit easier. 

## A Note on Using Databases in Code

Okay, accessing databases in code is annoying. Relational data is annoying. It doesn't matter which language you're using.
It doesn't matter if you're using PostgreSQL, MySQL, or Clickhouse, there's an innate complexity in relational data that
is simply not native to Go or most other general purpose programming languages. That's why many of us turn to tools like
ORMs, choosing to accept the cons of ORM magic because of the pros of having something else do the translation.

While yoyo is not an ORM, it is meant to address the same problem in a different way. By generating a database access
layer with entities for your project, tailor-made to your schema and without any reflection in sight.

## What does yoyo do?

yoyo is a code-generation tool which really does two things:

- Help manage your schema by generating migrations.
- Help you work with your schema by generating database access layer code in your project.

## Usage

### `yoyo generate`

[![Stability: Experimental](https://masterminds.github.io/stability/experimental.svg)](https://masterminds.github.io/stability/experimental.html)

yoyo's main function is generating database access code.

### `yoyo reverse`

[![Stability: Experimental](https://masterminds.github.io/stability/experimental.svg)](https://masterminds.github.io/stability/experimental.html)

Read an existing database and attempt to translate it to a schema in `yoyo.yml`

### `yoyo generate migration`

[![Stability: Experimental](https://masterminds.github.io/stability/experimental.svg)](https://masterminds.github.io/stability/experimental.html)

## Configuration

Configuration for yoyo is kept in your project's `yoyo.yml` file.

## Managing Database Connections

When running or generating migrations, Yoyo's connection to your database is environment-driven
and managed internally.

When running as a part of your app, Yoyo cedes control for your flexibility. Therefore, it needs
to be handed a connection in the form of a `*sql.DB`.  

## What Yoyo can't do

- Anything that crosses into another database. Yoyo is a single-database (single-schema) tool, so something like a MySQL
Foreign Key that references a table in a different schema won't work in Yoyo.
