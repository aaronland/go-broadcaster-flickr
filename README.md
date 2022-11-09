# go-broadcaster-flickr

Go package implementing the `aaronland/go-broadcaster` interfaces for broadcasting messages to Flickr.

## Documentation

Documentation is incomplete at this time.

## Tools

```
$> make cli
go build -mod vendor -o bin/broadcast cmd/broadcast/main.go
```

### broadcast

```
$> ./bin/broadcast \
	-broadcaster 'flickr://?credentials=file:///usr/local/flickr-broadcaster.txt' \
	-title testing \
	-body 'this is a test' \
	-image /usr/local/test.jpg
	
Int64UID#{FLICKR_PHOTO_ID}
```

Where `/usr/local/flickr-broadcaster.txt` is a valid [aaronland/go-flickr-api](https://github.com/aaronland/go-flickr-api#clients) client URI. For example, something like:

```
oauth1://?consumer_key={CONSUMER_KEY}&consumer_secret={CONSUMER_SECRET}&oauth_token={OAUTH_TOKEN}&oauth_token_secret={OAUTH_TOKEN_SECRET}
```

## See also

* https://github.com/aaronland/go-broadcaster
* https://github.com/aaronland/go-flickr-api