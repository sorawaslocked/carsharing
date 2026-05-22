//go:build integration

package postgres_test

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
	"time"

	sharedmodel "carsharing/shared/model"
	pgadapter "carsharing/user-service/internal/adapter/postgres"
	"carsharing/user-service/internal/model"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testcontainers/testcontainers-go"
	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

var testPool *pgxpool.Pool

func TestMain(m *testing.M) {
	os.Exit(run(m))
}

func run(m *testing.M) int {
	ctx := context.Background()

	pgc, err := tcpostgres.Run(ctx, "postgres:16-alpine",
		tcpostgres.WithDatabase("testdb"),
		tcpostgres.WithUsername("testuser"),
		tcpostgres.WithPassword("testpass"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(60*time.Second),
		),
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "start postgres container: %v\n", err)
		return 1
	}
	defer pgc.Terminate(ctx) //nolint:errcheck

	connStr, err := pgc.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		fmt.Fprintf(os.Stderr, "get connection string: %v\n", err)
		return 1
	}

	testPool, err = pgxpool.New(ctx, connStr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "create pool: %v\n", err)
		return 1
	}
	defer testPool.Close()

	// Wait for the pool to accept connections (defensive retry after the log-based wait).
	for i := range 10 {
		if err = testPool.Ping(ctx); err == nil {
			break
		}
		if i == 9 {
			fmt.Fprintf(os.Stderr, "ping postgres: %v\n", err)
			return 1
		}
		time.Sleep(500 * time.Millisecond)
	}

	if err := applyMigrations(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "apply migrations: %v\n", err)
		return 1
	}

	return m.Run()
}

func applyMigrations(ctx context.Context) error {
	entries, err := os.ReadDir("../../../migrations")
	if err != nil {
		return fmt.Errorf("read migrations dir: %w", err)
	}

	var upFiles []string
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".up.sql") {
			upFiles = append(upFiles, filepath.Join("../../../migrations", e.Name()))
		}
	}
	sort.Strings(upFiles)

	for _, path := range upFiles {
		sql, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("read %s: %w", path, err)
		}
		if _, err := testPool.Exec(ctx, string(sql)); err != nil {
			return fmt.Errorf("exec %s: %w", path, err)
		}
	}

	return nil
}

// truncate clears users (cascades to user_roles and documents) before each test.
func truncate(t *testing.T) {
	t.Helper()
	if _, err := testPool.Exec(context.Background(), "TRUNCATE users CASCADE"); err != nil {
		t.Fatalf("truncate: %v", err)
	}
}

func discardLog() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

func newUserRepo() *pgadapter.UserRepository {
	return pgadapter.NewUserRepository(discardLog(), testPool)
}

func newDocRepo() *pgadapter.DocumentRepository {
	return pgadapter.NewDocumentRepository(discardLog(), testPool)
}

func ptr[T any](v T) *T { return &v }

// testUser returns a valid model.User ready to be inserted.
func testUser() model.User {
	phone := "+12345678901"
	return model.User{
		Email:        "john@example.com",
		PhoneNumber:  &phone,
		FirstName:    "John",
		LastName:     "Doe",
		BirthDate:    time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
		PasswordHash: []byte("testhash"),
		Roles:        []sharedmodel.Role{sharedmodel.RoleUser},
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
}

func testUserWithEmail(email string) model.User {
	u := testUser()
	u.Email = email
	u.PhoneNumber = nil
	return u
}

// mustInsertUser inserts a user and returns its generated ID.
func mustInsertUser(t *testing.T, u model.User) string {
	t.Helper()
	id, err := newUserRepo().Insert(context.Background(), u)
	if err != nil {
		t.Fatalf("mustInsertUser: %v", err)
	}
	return id
}

func testDoc(userID string) model.Document {
	return testDocOfType(userID, model.DocumentImageTypeIDFront)
}

func testDocOfType(userID string, imageType model.DocumentImageType) model.Document {
	return model.Document{
		UserID:    userID,
		ImageType: imageType,
		Status:    model.DocumentStatusPending,
		Image:     sharedmodel.Image{Key: "documents/" + string(imageType) + "/key"},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func testDocAt(userID string, imageType model.DocumentImageType, ts time.Time) model.Document {
	return model.Document{
		UserID:    userID,
		ImageType: imageType,
		Status:    model.DocumentStatusPending,
		Image:     sharedmodel.Image{Key: "documents/" + string(imageType) + "/key"},
		CreatedAt: ts,
		UpdatedAt: ts,
	}
}
