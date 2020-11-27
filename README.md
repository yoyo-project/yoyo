# yoyo

A Repository Generator for Go

## Managing Database Connections

When running or generating migrations, Yoyo's connection to your database is environment-driven
and managed internally.

When running as a part of your app, Yoyo cedes control for your flexibility. Therefore, it needs
to be handed a connection in the form of a `*sql.DB`.  

## What Yoyo can't do

- Anything that crosses into another database. Yoyo is a single-database (single-schema) tool, so something like a MySQL
Foreign Key that references a table in a different schema won't work in Yoyo.
