# yoyo

[![MIT license](https://img.shields.io/badge/license-MIT-brightgreen.svg)](https://opensource.org/licenses/MIT)
[![Go Report Card](https://goreportcard.com/badge/github.com/dotvezz/yoyo)](https://goreportcard.com/report/github.com/dotvezz/yoyo)
[![codecov](https://codecov.io/gh/dotvezz/yoyo/branch/main/graph/badge.svg)](https://codecov.io/gh/dotvezz/yoyo)
[![CircleCI](https://circleci.com/gh/dotvezz/yoyo/tree/main.svg?style=shield)](https://circleci.com/gh/dotvezz/yoyo/tree/main)

A Repository Generator for Go

## Usage

### yoyo generate repos

Status: Concept (Pre-WIP)

yoyo's main function is generating repositories

### yoyo generate migration

Status: WIP

## Configuration

Status: WIP

Configuration for yoyo is kept in your project's `yoyo.yml` file.

## Managing Database Connections

When running or generating migrations, Yoyo's connection to your database is environment-driven
and managed internally.

When running as a part of your app, Yoyo cedes control for your flexibility. Therefore, it needs
to be handed a connection in the form of a `*sql.DB`.  

## What Yoyo can't do

- Anything that crosses into another database. Yoyo is a single-database (single-schema) tool, so something like a MySQL
Foreign Key that references a table in a different schema won't work in Yoyo.
