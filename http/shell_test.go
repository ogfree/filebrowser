package http

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/asdine/storm/v3"
	"github.com/gorilla/websocket"

	"github.com/filebrowser/filebrowser/v2/auth"
	"github.com/filebrowser/filebrowser/v2/settings"
	"github.com/filebrowser/filebrowser/v2/storage/bolt"
	"github.com/filebrowser/filebrowser/v2/users"
)

func TestShell(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "db")
	db, err := storm.Open(dbPath)
	if err != nil {
		t.Fatalf("failed to open db: %v", err)
	}
	defer db.Close()

	storage, err := bolt.NewStorage(db)
	if err != nil {
		t.Fatalf("failed to get storage: %v", err)
	}

	user := &users.User{
		ID:       1,
		Username: "test",
		Password: "password",
		Perm: users.Permissions{
			Execute: true,
		},
	}
	if err := storage.Users.Save(user); err != nil {
		t.Fatalf("failed to save user: %v", err)
	}

	sett := &settings.Settings{
		Key: []byte("test"),
	}
	if err := storage.Settings.Save(sett); err != nil {
		t.Fatalf("failed to save settings: %v", err)
	}

	handler := handle(shellHandler, "", storage, &settings.Server{EnableShell: true})
	ts := httptest.NewServer(handler)
	defer ts.Close()

	auther, err := auth.New(auth.MethodJSONAuth, storage)
	if err != nil {
		t.Fatalf("failed to create auther: %v", err)
	}

	token, err := auther.GetToken(user, "", time.Hour)
	if err != nil {
		t.Fatalf("failed to get token: %v", err)
	}

	// Convert http:// to ws://
	u := "ws" + strings.TrimPrefix(ts.URL, "http")

	// Connect to the server
	header := http.Header{}
	header.Add("X-AUTH-TOKEN", token)
	ws, _, err := websocket.DefaultDialer.Dial(u, header)
	if err != nil {
		t.Fatalf("%v", err)
	}
	defer ws.Close()

	// Send a message
	if err := ws.WriteMessage(websocket.TextMessage, []byte("echo hello\n")); err != nil {
		t.Fatalf("%v", err)
	}

	// Read the response
	_, p, err := ws.ReadMessage()
	if err != nil {
		t.Fatalf("%v", err)
	}

	fmt.Println(string(p))

	// Check if the response is correct
	if !strings.Contains(string(p), "hello") {
		t.Fatalf("unexpected response: %s", p)
	}
}
