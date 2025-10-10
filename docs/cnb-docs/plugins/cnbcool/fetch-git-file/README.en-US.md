# fetch-git-file

Fetch CNB git files to local

## Parameters

### token

- type: String
- required: No
- default: `CNB_TOKEN`

Token required for accessing CNB open API

Requires read permission for file management `repo-contents:r`

Defaults to the environment variable `CNB_TOKEN`

### slug

- type: String
- required: Yes

Repository path

### ref

- type: String
- required: Yes

Repository branch, tag, or SHA

### files

- type: String
- required: Yes

List of files to fetch

Multiple lines of text, with one path per line, for example:

```shell
a/b/c.txt
d/e.ts
f.yml
```

### target

- type: String
- required: Yes
- default: _tmp_

Directory to store the files

Use `.` for the repository root directory

## Usage in Cloud Native Builds

```yaml
main:
  push:
    - stages:
        - name: fetch-git-file
          image: cnbcool/fetch-git-file
          settings:
            slug: xx/xx
            ref: master
            files: |
              a/b/c.txt
              d/e.ts
              f.yml
```
