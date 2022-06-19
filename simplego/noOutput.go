package main

import (
	"simplego/pulumi"
	"simplego/s3"
)

func noOutput(ctx *pulumi.Context) error {
	bucket, err := s3.NewBucket(ctx, "my-bucket", &s3.BucketArgs{
		Website: &s3.BucketWebsiteArgs{
			// Pass a plain *string in!
			IndexDocument: pulumi.Ptr("index.html"),
		},
	})
	if err != nil {
		return err
	}

	_, err = s3.NewBucketObject(ctx, "my-obj", &s3.BucketObjectArgs{
		// Bucket is typed as `string`, and bucket.Bucket() returns a string!
		Bucket:  bucket.Bucket(),
		Content: pulumi.Ptr("<h1>Hello, world!</h1>"),
		Key:     pulumi.Ptr("index.html"),
	})
	if err != nil {
		return err
	}

	// We can just concatenate the strings!
	ctx.Export("url", "http://"+bucket.WebsiteEndpoint())
	return nil
}
