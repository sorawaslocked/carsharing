package client

import (
	"context"
	"log/slog"

	pkglog "carsharing/shared/pkg/log"
	"carsharing/shared/pkg/utils"

	usersvc "carsharing/protos/gen/service/user"

	"google.golang.org/grpc"
)

type DocumentAnalyzer struct {
	log    *slog.Logger
	client usersvc.DocumentAnalyzerServiceClient
}

func NewDocumentAnalyzer(log *slog.Logger, conn *grpc.ClientConn) *DocumentAnalyzer {
	return &DocumentAnalyzer{
		log:    pkglog.WithComponent(log, "adapter.grpc.client.DocumentAnalyzer"),
		client: usersvc.NewDocumentAnalyzerServiceClient(conn),
	}
}

func (c *DocumentAnalyzer) Analyze(ctx context.Context, documentID string, objectKey string, userID string) {
	log := pkglog.WithMetadata(pkglog.WithMethod(c.log, "Analyze"), utils.MetadataFromCtx(ctx))

	if _, err := c.client.Analyze(ctx, &usersvc.AnalyzeRequest{
		DocumentId: documentID,
		ObjectKey:  objectKey,
		UserId:     userID,
	}); err != nil {
		log.Error("grpc: analyzing document",
			slog.String("documentID", documentID),
			pkglog.Err(err),
		)
	}
}
