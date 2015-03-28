s3up [![Build Status](https://travis-ci.org/tscolari/s3up.svg?branch=master)](https://travis-ci.org/tscolari/s3up)
==============

Usage
-----

```
Usage:
  s3up [flags]
Flags:
  -i, --access-key="": AWS Access Key
  -b, --bucket-name="": Target S3 bucket
  -e, --endpoint-url="https://s3.amazonaws.com": the s3 region endpoint url (see http://docs.aws.amazon.com/general/latest/gr/rande.html#s3_region)
  -n, --file-name="": How the file will be called on s3
  -h, --help=false: help for s3up
  -s, --secret-key="": AWS Secret Key
  -v, --verbose=false: Verbose mode
  -k, --versions-to-keep=5: Number of versions to keep
```

e.g: the command:

```
  pg_dump | bzip2 -c | s3up --access-key X --secret-key Y --bucket-name Z --file-name my-pg-bkp
```

will result in storing the file on S3 as:

```
  Z/my-backup/timestamp
```


