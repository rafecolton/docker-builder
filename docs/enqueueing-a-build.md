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

### Request Fields

Required Fields:

* `account / type: string` - the GitHub account for the repo being cloned
* `repo / type: string` - the name of the repo
* `ref / type: string` - the ref (can be any valid/unambiguous ref - a branch, tag, sha, etc)

Other Fields: 

* `api_token / type: string` - the GitHub api token (not required for public repos)
* `depth / type: string (must be int > 0)` - clone depth (default: no `--depth` argument passed to `git clone`)
