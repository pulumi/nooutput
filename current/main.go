package main

import (
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/s3"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		bucket, err := s3.NewBucket(ctx, "my-bucket", &s3.BucketArgs{
			Website: &s3.BucketWebsiteArgs{
				// Wrap a string into a pulumi.StringOutput (which implents pulumi.StringInput)
				IndexDocument: pulumi.String("index.html"),
			},
			// Must know to use `pulumi.StringMap` and other complex object wrappers
			Tags: pulumi.StringMap{
				"Owner": pulumi.String("lukehoban"),
			},
		})
		if err != nil {
			return err
		}

		_, err = s3.NewBucketObject(ctx, "my-obj", &s3.BucketObjectArgs{
			// Pass an Output as an Input
			Bucket:  bucket.Bucket,
			Content: pulumi.String("<h1>Hello, world!</h1>"),
			Key:     pulumi.String("index.html"),
		})
		if err != nil {
			return err
		}

		// Two ways to turn a pulumi.StringOutput into a concatenated string:
		// 1. Know to use special `pulumi.Sprintf`
		// 2. Use the more fundamental ApplyT
		ctx.Export("url", pulumi.Sprintf("http://%s", bucket.WebsiteEndpoint))

		ctx.Export("url2", bucket.WebsiteEndpoint.ApplyT(func(s string) string {
			return "http://" + s
		}).(pulumi.StringOutput))

		// Object valued outputs need to use `MapIndex` to lookup (or Apply)
		ctx.Export("bucketOwner", bucket.TagsAll.MapIndex(pulumi.String("Owner")))
		return nil
	})
}
