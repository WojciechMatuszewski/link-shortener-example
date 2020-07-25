# Serverless link shortener example

Each object in `S3` can have metadata specified that, when visited, redirect to some address.

This is a basic example of doing such a thing.

A `POST` is send to an endpoint, which creates 0 bytes `S3` object with correct metadata.
Returns the object URL.


