# Files (messenger)

This is file sharing component of the [Messenger](https://github.com/barpav/messenger) pet-project.

## Functions

* File uploading and downloading with file-level access (accessible for specified users or public access).

* Storing file data (MongoDB, [GridFS](https://www.mongodb.com/docs/manual/core/gridfs/)).

* Explicit or automatic (unused files garbage collector) file deletion.

* Maintaining and updating file usage statistics by reference counting (RabbitMQ).

See microservice [REST API](https://barpav.github.io/msg-api-spec/#/files) and [deployment diagram](https://github.com/barpav/messenger#deployment-diagram) for details.