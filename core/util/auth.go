package util

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/go-faster/errors"
	"github.com/medama-io/medama/model"
	"go.jetpack.io/typeid"
)

const (
	// DefaultCipherKeySize is the default size of the cipher key.
	DefaultCipherKeySize = 32
)

type AuthService struct {
	// Cache used to store session tokens.
	Cache *Cache
	// Key used to encrypt session tokens.
	aes32Key []byte
}

// NewAuthService returns a new instance of AuthService.
func NewAuthService(ctx context.Context) (*AuthService, error) {
	// Generate a new random key for encrypting session tokens.
	// Since we store sessions in an in-memory cache, it doesn't
	// matter if the key doesn't persist as sessions will be
	// invalidated when the server restarts.
	key := make([]byte, DefaultCipherKeySize)
	_, err := rand.Read(key)
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate cipher key")
	}

	return &AuthService{
		aes32Key: key,
		Cache:    NewCache(ctx, model.SessionDuration),
	}, nil
}

// HashPassword hashes a password using argon.
func (a *AuthService) HashPassword(password string) (string, error) {
	hash, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		return "", errors.Wrap(err, "failed to hash password")
	}

	return hash, nil
}

// ComparePasswords compares a supplied password with a stored hash.
func (a *AuthService) ComparePasswords(suppliedPassword string, storedHash string) (bool, error) {
	match, err := argon2id.ComparePasswordAndHash(suppliedPassword, storedHash)
	if err != nil {
		return false, errors.Wrap(err, "failed to compare passwords")
	}

	return match, nil
}

// EncryptSession encrypts a session token and stores it in the cache.
func (a *AuthService) EncryptSession(ctx context.Context, sessionId string, duration time.Duration) (string, error) {
	// Create a new AES cipher block.
	block, err := aes.NewCipher(a.aes32Key)
	if err != nil {
		return "", errors.Wrap(err, "auth: encrypt")
	}

	// Wrap the block in a GCM cipher.
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", errors.Wrap(err, "auth: encrypt")
	}

	// Create a random 12 byte nonce.
	nonce := make([]byte, aesgcm.NonceSize())
	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return "", errors.Wrap(err, "auth: encrypt")
	}

	// Authenticate cookie name and value with {name:value} format.
	plaintext := fmt.Sprintf("%s:%s", model.SessionCookieName, sessionId)

	// Encrypt with nonce for variable ciphertext.
	encryptedValue := aesgcm.Seal(nonce, nonce, []byte(plaintext), nil)

	// Return encrypted session token.
	return string(encryptedValue), nil
}

// DecryptSession decrypts a session cookie to return the session token.
func (a *AuthService) DecryptSession(ctx context.Context, session string) (string, error) {
	// Create a new AES cipher block.
	block, err := aes.NewCipher(a.aes32Key)
	if err != nil {
		return "", errors.Wrap(err, "auth: decrypt")
	}

	// Wrap the block in a GCM cipher.
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", errors.Wrap(err, "auth: decrypt")
	}

	// Check for potential index out of range error.
	nonceSize := aesgcm.NonceSize()
	if len(session) < nonceSize {
		return "", model.ErrInvalidSession
	}

	// Extract nonce from session token.
	nonce, ciphertext := session[:nonceSize], session[nonceSize:]

	// Decrypt session token.
	plaintext, err := aesgcm.Open(nil, []byte(nonce), []byte(ciphertext), nil)
	if err != nil {
		return "", model.ErrInvalidSession
	}

	// Split plaintext into name and value.
	name, value, ok := strings.Cut(string(plaintext), ":")
	if !ok {
		return "", model.ErrInvalidSession
	}

	// Check if cookie name is valid.
	if name != model.SessionCookieName {
		return "", model.ErrInvalidSession
	}

	return value, nil
}

// CreateSession creates a new session token and stores it in the cache.
// This returns an encrypted session token as a cookie.
func (a *AuthService) CreateSession(ctx context.Context, userId string) (*http.Cookie, error) {
	// Generate session token.
	sessionIdType, err := typeid.WithPrefix("sess")
	if err != nil {
		return nil, errors.Wrap(err, "auth: session")
	}
	sessionId := sessionIdType.String()

	// Create session cookie.
	cookie := &http.Cookie{
		Name:     model.SessionCookieName,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	}

	// Encrypt session token.
	encryptedSession, err := a.EncryptSession(ctx, sessionId, model.SessionDuration)

	// Update cookie value with encrypted base64 enoded session token.
	encodedSession := base64.URLEncoding.EncodeToString([]byte(encryptedSession))
	cookie.Value = encodedSession

	// Set session token in cache.
	a.Cache.Set(sessionId, userId, model.SessionDuration)

	return cookie, err
}

// ReadSession decrypts and reads a session token from the cache.
// This returns the user id associated with the encrypted session value.
func (a *AuthService) ReadSession(ctx context.Context, session string) (string, error) {
	// Decode base64 encoded session token.
	encryptedSession, err := base64.URLEncoding.DecodeString(session)
	if err != nil {
		return "", model.ErrInvalidSession
	}

	// Decrypt session token.
	sessionId, err := a.DecryptSession(ctx, string(encryptedSession))
	if err != nil {
		return "", errors.Wrap(err, "session")
	}

	// Check if session exists.
	userId, err := a.Cache.Get(ctx, sessionId)
	if err != nil {
		return "", model.ErrSessionNotFound
	}

	return userId.(string), nil
}

// RevokeSession deletes a session token from the cache.
func (a *AuthService) RevokeSession(ctx context.Context, sessionId string) {
	a.Cache.Delete(sessionId)
}
