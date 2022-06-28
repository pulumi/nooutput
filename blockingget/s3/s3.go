package s3

import (
	"context"
	"fmt"
	"simplego/pulumi"

	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/s3"
	p "github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Bucket struct {
	inner s3.Bucket
}

func (b *Bucket) WebsiteEndpoint() string {
	// ch := make(chan string)
	// b.inner.WebsiteEndpoint.ApplyT(func(s string) interface{} {
	// 	ch <- s
	// 	return nil
	// })
	// return <-ch

	ch := make(chan string)
	go func() {
		v, known, secret, deps, err := b.inner.WebsiteEndpoint.OutputState.Await(context.Background())
		fmt.Printf("WebsiteEndpoint():%v,%v,%v,%v,%v\n", v, known, secret, deps, err)
		if !known {
			ch <- "UNKNOWN"
		} else {
			ch <- v.(string)
		}
	}()
	ret := <-ch
	return ret
}

func (b *Bucket) Bucket() string {
	ch := make(chan string)
	go func() {
		v, known, secret, deps, err := b.inner.Bucket.OutputState.Await(context.Background())
		fmt.Printf("Bucket():%v,%v,%v,%v,%v\n", v, known, secret, deps, err)
		if !known {
			ch <- "UNKNOWN"
		} else {
			ch <- v.(string)
		}
	}()
	ret := <-ch
	return ret
}

type BucketWebsiteArgs struct {
	IndexDocument *string
}

type BucketArgs struct {
	Website *BucketWebsiteArgs
}

func NewBucket(ctx *pulumi.Context,
	name string, args *BucketArgs, opts ...p.ResourceOption) (*Bucket, error) {
	if args == nil {
		args = &BucketArgs{}
	}

	innerArgs := s3.BucketArgs{
		Website: &s3.BucketWebsiteArgs{
			IndexDocument: pulumi.PtrTokenToInput(args.Website.IndexDocument),
		},
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
