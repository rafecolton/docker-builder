## Enqueueing a Build

### Example Request

```bash
#!/bin/bash

curl -XPOST -H 'Content-Type: application/json' 'http://localhost:5000/docker-build' -d '
{
  "account": "my-account",
  "repo": "my-repo",
  "ref": "master"
}
'
```

### Example Response

```javascript
{
  "account": "my-account",
  "completed": "0001-01-01T00:00:00Z",
  "created": "2014-07-06T14:02:01.92446296-07:00",
  "id": "035c4ea0-d73b-5bde-7d6f-c806b04f2ec3",
  "log_route": "/jobs/035c4ea0-d73b-5bde-7d6f-c806b04f2ec3/tail?n=100",
  "ref": "master",
  "repo": "my-repo",
  "status": "created"
}
```

### Request Fields

Required Fields:

* `account / type: string` - the GitHub account for the repo being cloned
* `repo / type: string` - the name of the repo
* `ref / type: string` - the ref (can be any valid/unambiguous ref - a branch, tag, sha, etc)

Other Fields:

* `api_token / type: string` - the GitHub api token (not required for public repos)
* `depth / type: string (must be int > 0)` - clone depth (default: no `--depth` argument passed to `git clone`)
* `sync / type: bool` - sets whether the server should respond to a
  request immediately or block until build completes. The default
behavior is to respond immediately (value: `false`). Set this to `true`
to wait for a response until after the build, tag, and push phases are
all complete.
