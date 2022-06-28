package main

import (
	"fmt"
	"simplego/pulumi"
	"simplego/s3"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		bucket, err := s3.NewBucket(ctx, "my-bucket", &s3.BucketArgs{
			Website: &s3.BucketWebsiteArgs{
				// Pass a plain *string in!
				IndexDocument: pulumi.Ptr("index.html"),
			},
		})
		if err != nil {
			return err
		}

		for _, filename := range []string{"index.html", "other.html"} {
			_, err = s3.NewBucketObject(ctx, filename, &s3.BucketObjectArgs{
				// Bucket is typed as `string`, and bucket.Bucket() returns a string!
				Bucket:  bucket.Bucket(),
				Content: pulumi.Ptr("<h1>Hello, world!</h1>"),
				Key:     &filename,
			})
			if err != nil {
				return err
			}
		}

		fmt.Printf("Bucket website endpoint is: %s\n", bucket.WebsiteEndpoint())

		// We can just concatenate the strings!
		ctx.Export("url", "http://"+bucket.WebsiteEndpoint())

		// And if we *really* want to Apply (to be able to do anything we can express in Pulumi today),
		// we still can
		ctx.Export("url2", pulumi.Apply(bucket.WebsiteEndpoint(), func(endpoint string) string {
			return "http://" + endpoint
		}))
		return nil
	})
}
