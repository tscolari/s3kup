s3up [![Build Status](https://travis-ci.org/tscolari/s3up.svg?branch=master)](https://travis-ci.org/tscolari/s3up)
==============

Usage
-----

```
Usage:
  s3up [flags]
  s3up [command]

Available Commands:
  push        Pushes the piped input to s3
  list        List remote stored versions
  pull        Get remote version contents
  help        Help about any command

Flags:
  -i, --access-key="": AWS Access Key
  -b, --bucket-name="": Target S3 bucket
  -e, --endpoint-url="https://s3.amazonaws.com": the s3 region endpoint url (see http://docs.aws.amazon.com/general/latest/gr/rande.html#s3_region)
  -n, --file-name="": How the file will be called on s3
  -h, --help=false: help for s3up
  -s, --secret-key="": AWS Secret Key
  -v, --verbose=false: Verbose mode
```

Pushing backups
---------------

```
s3up help push                                                                                                                           s3up/git/master
Pushes the pipped input to s3, as a versioned backup

Usage:
  s3up push [flags]
Flags:
  -h, --help=false: help for push
  -k, --versions-to-keep=5: Number of versions to keep

Global Flags:
  -i, --access-key="": AWS Access Key
  -b, --bucket-name="": Target S3 bucket
  -e, --endpoint-url="https://s3.amazonaws.com": the s3 region endpoint url (see http://docs.aws.amazon.com/general/latest/gr/rande.html#s3_region)
  -n, --file-name="": How the file will be called on s3
  -s, --secret-key="": AWS Secret Key
  -v, --verbose=false: Verbose mode
```

It will always push the piped input as the content for the backup

e.g:

```
  pg_dump | bzip2 -c | s3up push --access-key X --secret-key Y --bucket-name Z --file-name my-pg-bkp
```

will the input on S3 as:

```
  s3://Z/my-backup/unixnanotimestamp
```

Listing backups
---------------

```
s3up help list                                                                                                                          s3up/git/master !
List remote stored versions

Usage:
  s3up list [flags]
Flags:
  -h, --help=false: help for list

Global Flags:
  -i, --access-key="": AWS Access Key
  -b, --bucket-name="": Target S3 bucket
  -e, --endpoint-url="https://s3.amazonaws.com": the s3 region endpoint url (see http://docs.aws.amazon.com/general/latest/gr/rande.html#s3_region)
  -n, --file-name="": How the file will be called on s3
  -s, --secret-key="": AWS Secret Key
  -v, --verbose=false: Verbose mode
```

e.g:

```
  s3up list --access-key X --secret-key Y --bucket-name Z --file-name my-pg-bkp

  * 1427554100187348642 [10B at 2015-03-28T14:48:21.000Z]
  * 1427571015905296950 [123MB at 2015-03-28T19:30:17.000Z]
  * 1427835207908555851 [130MB at 2015-03-31T20:53:29.000Z]
```

Fetching a backup
-----------------

```
s3up help pull
Get remote version and print it's contents to STDOUT

Usage:
  s3up pull [flags]
Flags:
  -h, --help=false: help for pull

Global Flags:
  -i, --access-key="": AWS Access Key
  -b, --bucket-name="": Target S3 bucket
  -e, --endpoint-url="https://s3.amazonaws.com": the s3 region endpoint url (see http://docs.aws.amazon.com/general/latest/gr/rande.html#s3_region)
  -n, --file-name="": How the file will be called on s3
  -s, --secret-key="": AWS Secret Key
  -v, --verbose=false: Verbose mode
```

e.g:

1. Fetching the last backup

```
  s3up pull --access-key X --secret-key Y --bucket-name Z --file-name my-pg-bkp > dump.bz2
```

2. Fetching a specific version

```
  s3up pull 1427571015905296950 --access-key X --secret-key Y --bucket-name Z --file-name my-pg-bkp > dump.bz2
```

LICENSE
=======

Copyright 2015 Tiago Scolari, under Apache License.
