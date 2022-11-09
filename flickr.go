package flickr

import (
	"bytes"
	"context"
	"fmt"
	"github.com/aaronland/go-broadcaster"
	"github.com/aaronland/go-flickr-api/client"
	"github.com/aaronland/go-image-encode"
	"github.com/aaronland/go-uid"
	"github.com/sfomuseum/runtimevar"
	_ "image"
	"log"
	"net/url"
	"time"
)

func init() {
	ctx := context.Background()
	broadcaster.RegisterBroadcaster(ctx, "flickr", NewFlickrBroadcaster)
}

type FlickrBroadcaster struct {
	broadcaster.Broadcaster
	flickr_client client.Client
	encoder       encode.Encoder
	logger        *log.Logger
}

func NewFlickrBroadcaster(ctx context.Context, uri string) (broadcaster.Broadcaster, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse URI, %w", err)
	}

	q := u.Query()

	creds_uri := q.Get("credentials")

	if creds_uri == "" {
		return nil, fmt.Errorf("Missing ?credentials= parameter")
	}

	rt_ctx, rt_cancel := context.WithTimeout(ctx, 5*time.Second)
	defer rt_cancel()

	client_uri, err := runtimevar.StringVar(rt_ctx, creds_uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to derive URI from credentials, %w", err)
	}

	cl, err := client.NewClient(ctx, client_uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to create new Flickr client, %w", err)
	}

	enc, err := encode.NewEncoder(ctx, "png://")

	if err != nil {
		return nil, fmt.Errorf("Failed to create image encoder, %w", err)
	}

	logger := log.Default()

	br := &FlickrBroadcaster{
		flickr_client: cl,
		encoder:       enc,
		logger:        logger,
	}

	return br, nil
}

func (b *FlickrBroadcaster) BroadcastMessage(ctx context.Context, msg *broadcaster.Message) (uid.UID, error) {

	switch len(msg.Images) {
	case 0:
		return nil, fmt.Errorf("Missing image")
	case 1:
		// pass
	default:
		return nil, fmt.Errorf("You can only broadcast one image at a time.")
	}

	r := new(bytes.Buffer)

	err := b.encoder.Encode(ctx, msg.Images[0], r)

	if err != nil {
		return nil, fmt.Errorf("Failed to encode image, %w", err)
	}

	args := &url.Values{}

	if msg.Title != "" {
		args.Set("title", msg.Title)
	}

	if msg.Body != "" {
		args.Set("description", msg.Body)
	}
	
	photo_id, err := client.UploadAsyncWithClient(ctx, b.flickr_client, r, args)

	if err != nil {
		return nil, fmt.Errorf("Failed to upload image, %w", err)
	}

	return uid.NewInt64UID(ctx, photo_id)
}

func (b *FlickrBroadcaster) SetLogger(ctx context.Context, logger *log.Logger) error {
	b.logger = logger
	return nil
}
