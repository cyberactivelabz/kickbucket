# Wordlists

A wordlist is a plain text file, one word per line, lowercase, no duplicates. BucketMap combines each word with your keyword to build candidate bucket names.

## Use your own

```sh
bucketmap scan acme --wordlist ./my-words.txt
```

## Share one

Add it to `wordlists/` and open a pull request. See [Contributing](/contributing).
