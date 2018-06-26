package service

import (
	"cloud.google.com/go/storage"
	"golang.org/x/net/context"
	"google.golang.org/api/option"
	"time"
)

func UploadFromCloudStorage(object string, content []byte) (string, []string) {
	var errs []string

	base_url := "https://storage.googleapis.com/"
	var url string

	ctx := context.Background()

	// Creates a client.
	client, err := storage.NewClient(ctx, option.WithCredentialsFile("controlcenter-3e32a-firebase-adminsdk-6iyax-3f53416bc5.json"))
	if err != nil {
		errs = append(errs, err.Error())
	}

	if len(errs) == 0 {
		// Sets the name for the new bucket.
		bucketName := "controlcenter-3e32a.appspot.com"

		// Creates a Bucket instance.
		bucket := client.Bucket(bucketName)

		obj := bucket.Object(object)
		// Write something to obj.
		// w implements io.Writer.
		ww := obj.NewWriter(ctx)
		// Write some text to obj. This will overwrite whatever is there.

		ww.Deleted = time.Now()
		ww.Created = time.Now()
		// We can't check that the Write fails, since it depends on the write to the
		// underling fakeTransport failing which is racy.
		ww.CacheControl = "Cache-Control:private, max-age=0, no-transform"
		ww.Write(content)

		// Close, just like writing a file.
		if err := ww.Close(); err != nil {
			errs = append(errs, err.Error())
		}

		objAttrs, err := obj.Attrs(ctx)
		if err != nil {
			errs = append(errs, err.Error())
		}

		url = base_url + bucketName + "/" + objAttrs.Name
	}

	return url, errs
}