package s3

import (
	"simplego/pulumi"

	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/s3"
	p "github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Bucket struct {
	inner s3.Bucket
}

func (b *Bucket) WebsiteEndpoint() string {
	return pulumi.OutputToToken(b.inner.WebsiteEndpoint)
}

func (b *Bucket) Bucket() string {
	return pulumi.OutputToToken(b.inner.Bucket)
}

func (b *Bucket) TagsAll() pulumi.StringMap {
	return pulumi.StringMap(b.inner.TagsAll)
}

type BucketWebsiteArgs struct {
	IndexDocument *string
}

type BucketArgs struct {
	Website *BucketWebsiteArgs
	Tags    map[string]string
}

func NewBucket(ctx *pulumi.Context,
	name string, args *BucketArgs, opts ...p.ResourceOption) (*Bucket, error) {
	if args == nil {
		args = &BucketArgs{}
	}

	// TODO: Does this need to be able to accept a token representation of a whole map?
	var tags map[string]p.StringOutput
	if args.Tags != nil {
		tags = map[string]p.StringOutput{}
		for k, v := range args.Tags {
			tags[k] = pulumi.TokenToInput(v).ToStringOutput()
		}
	}

	innerArgs := s3.BucketArgs{
		Website: &s3.BucketWebsiteArgs{
			IndexDocument: pulumi.PtrTokenToInput(args.Website.IndexDocument),
		},
		Tags: p.ToStringMapOutput(tags),
	}

	var resource Bucket
	err := ctx.RegisterResource("aws:s3/bucket:Bucket", name, innerArgs, &resource.inner, opts...)
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

type BucketObject struct {
	inner s3.BucketObject
}

type BucketObjectArgs struct {
	Bucket  string
	Key     *string
	Content *string
}

func NewBucketObject(ctx *pulumi.Context,
	name string, args *BucketObjectArgs, opts ...p.ResourceOption) (*BucketObject, error) {
	if args == nil {
		args = &BucketObjectArgs{}
	}

	innerArgs := s3.BucketObjectArgs{
		Bucket:  pulumi.TokenToInput(args.Bucket),
		Key:     pulumi.PtrTokenToInput(args.Key),
		Content: pulumi.PtrTokenToInput(args.Content),
	}

	var resource BucketObject
	err := ctx.RegisterResource("aws:s3/bucketObject:BucketObject", name, innerArgs, &resource.inner, opts...)
	if err != nil {
		return nil, err
	}
	return &resource, nil
}
