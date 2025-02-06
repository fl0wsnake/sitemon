// package deploy

// import (
// 	gcf "github.com/pulumi/pulumi-gcp/sdk/v7/go/gcp/cloudfunctions"
// 	gcs "github.com/pulumi/pulumi-gcp/sdk/v7/go/gcp/storage"
// 	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
// )

// const bucketName = "sitemon"
// const functionName = "sitemon"

// func main() {
// 	pulumi.Run(func(ctx *pulumi.Context) error {
// 		_, err := gcs.NewBucket(ctx, bucketName, nil)
// 		if err != nil {
// 			return err
// 		}

// 		_, err = gcf.NewFunction(ctx, functionName, &gcf.FunctionArgs{})
// 		if err != nil {
// 			return err
// 		}
// 		return nil
// 	})
// }
