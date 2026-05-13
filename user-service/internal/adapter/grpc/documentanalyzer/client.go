package documentanalyzer

import (
	"context"
	"log/slog"

	usersvc "github.com/sorawaslocked/car-rental-protos/gen/service/user"
	pkglog "github.com/sorawaslocked/car-rental-user-service/internal/pkg/log"
	"github.com/sorawaslocked/car-rental-user-service/internal/pkg/utils"
	"google.golang.org/grpc"
)

type Client struct {
	log    *slog.Logger
	client usersvc.DocumentAnalyzerServiceClient
}

func NewClient(log *slog.Logger, conn *grpc.ClientConn) *Client {
	return &Client{
		log:    pkglog.WithComponent(log, "adapter.grpc.DocumentAnalyzerClient"),
		client: usersvc.NewDocumentAnalyzerServiceClient(conn),
	}
}

func (c *Client) Analyze(ctx context.Context, documentID string, storageURL string) {
	logger := pkglog.WithMetadata(pkglog.WithMethod(c.log, "Analyze"), utils.MetadataFromCtx(ctx))

	if _, err := c.client.Analyze(ctx, &usersvc.AnalyzeRequest{
		DocumentId: documentID,
		StorageUrl: storageURL,
	}); err != nil {
		logger.Error("grpc: analyzing document",
			slog.String("documentID", documentID),
			pkglog.Err(err),
		)
	}
}
