## Job Control

There are several routes available for managing the status of jobs being
processed by the build server.

### POST /jobs

Create a new job (equivalent to `/docker-build`)

Example Request:

```bash
curl -s -XPOST http://localhost:5000/jobs -d '{"account":"rafecolton","repo":"docker-builder","ref":"master"}'
```

Example Response:

```javascript
{
  "account": "rafecolton",
  "completed": "0001-01-01T00:00:00Z",
  "created": "2014-07-06T14:02:01.92446296-07:00",
  "id": "035c4ea0-d73b-5bde-7d6f-c806b04f2ec3",
  "info_route": "http://localhost:5000/jobs/035c4ea0-d73b-5bde-7d6f-c806b04f2ec3",
  "log_route": "http://localhost:5000/jobs/035c4ea0-d73b-5bde-7d6f-c806b04f2ec3/tail?n=100",
  "ref": "master",
  "repo": "docker-builder",
  "status": "created"
}
```

### GET /jobs

Get a list of all jobs

Example Request:

```bash
curl -s -XGET http://localhost:5000/jobs
```

Example Response:

```javascript
[
  {
    "account": "rafecolton",
    "completed": "0001-01-01T00:00:00Z",
    "created": "2014-07-06T14:02:01.92446296-07:00",
    "id": "035c4ea0-d73b-5bde-7d6f-c806b04f2ec3",
    "info_route": "http://localhost:5000/jobs/035c4ea0-d73b-5bde-7d6f-c806b04f2ec3",
    "log_route": "http://localhost:5000/jobs/035c4ea0-d73b-5bde-7d6f-c806b04f2ec3/tail?n=100",
    "ref": "master",
    "repo": "docker-builder",
    "status": "created"
  },
  //...
]
```

**NEW:** filter `/jobs` by adding a query string

The `/jobs` route may be filtered by the following fields:

* `account`
* `bobfile`
* `id`
* `ref`
* `repo`
* `status` - valid statuses include
  - `created`
  - `cloning`
  - `building`
  - `errored`
  - `completed`
  - `validating` (used for tests only)


It is also possible to filter by multiple fields.  For example, to get a
list of failed builds for my project "foo-project":

```bash
curl -s -XGET http://localhost:5000/jobs?repo=foo-project&status=errored
```

**Nuances:** If you specify a filter not on the list, it is ignored.  If
you specify a filter nultiple times, the first will be used and
the remaining will be inored.

For example, consider the following query string:

```
?repo=foo-project&repo=bar-project&baz=gorp
```

The returned job list will be filtered by `repo=foo-project`

### GET /jobs/:id

Get info job `:id`

Example Request:

```bash
curl -s -XGET http://localhost:5000/jobs/035c4ea0-d73b-5bde-7d6f-c806b04f2ec3
```

Example Response:

```javascript
{
  "account": "rafecolton",
  "completed": "0001-01-01T00:00:00Z",
  "created": "2014-07-06T14:02:01.92446296-07:00",
  "id": "035c4ea0-d73b-5bde-7d6f-c806b04f2ec3",
  "info_route": "http://localhost:5000/jobs/035c4ea0-d73b-5bde-7d6f-c806b04f2ec3",
  "log_route": "http://localhost:5000/jobs/035c4ea0-d73b-5bde-7d6f-c806b04f2ec3/tail?n=100",
  "ref": "master",
  "repo": "docker-builder",
  "status": "created"
}
```

### GET /jobs/:id/tail?n=100

Get the last `n` lines of the log from job `:id`

Example Request:

```bash
curl -s -XGET http://localhost:5000/jobs/035c4ea0-d73b-5bde-7d6f-c806b04f2ec3/tail?n=100
```

Example Response:

```javascript
// ... 100 lines worth of logs
```
