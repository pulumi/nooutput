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

		_, err = s3.NewBucketObject(ctx, "my-obj", &s3.BucketObjectArgs{
			// Bucket is typed as `string`, and bucket.Bucket() returns a string!
			Bucket:  bucket.Bucket(),
			Content: pulumi.Ptr("<h1>Hello, world!</h1>"),
			Key:     pulumi.Ptr("index.html"),
		})
		if err != nil {
			return err
		}

		url := "http://" + bucket.WebsiteEndpoint()
		fmt.Printf("%s\n", url)

		// We can just concatenate the strings!
		pulumi.Printf("url %s", url)
		ctx.Export("url", url)
		return nil
	})
}
