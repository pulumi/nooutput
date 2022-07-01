# Outputless Pulumi Experiments

A few experiments for alternative programming model projections for working with the Pulumi resource model.

All examples here are in Go, but in principle any of these approaches could be applied to other Pulumi languages.

* [current](./current): The current Pulumi programming model
* [tokenization](./tokenization/): A model where there are no outputs, and "tokens" which have the correct underlying type as the true data as used to smuggle references to eventual values (+ dependencies).
* [blockingget](./blockingget/): A model where accessors block on outputs, and unknowns are embedded into tokens. __Note__: _This approach doesn't really work, as dependencies are not retained._
* [tokenization-blockingget](./tokenization-blockingget/): A combination of the above, where tokens are used in all cases, but accessors block on outputs being available, so that resolved values of outputs can appear in the token (or `UNKNOWN` for unknowns during preview).  Makes printing and/or stepping through programs substantially "simpler".  In practice, some risk of limiting paralelism by blocking "too soon".

Here's the end state for the `tokenization-blockingget` approach:

```go
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
            // Pass a plain map[string]string in!
			Tags: map[string]string{
				"Owner": "lukehoban",
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

		// Object valued outputs need to use `Lookup`
		ctx.Export("bucketOwner", bucket.TagsAll().Lookup("Owner"))
		return nil
	})
}
```