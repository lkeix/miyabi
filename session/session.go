package session

import (
	"crypto/aes"
	"crypto/cipher"
	"math/rand"
	"miyabi"
	"net/http"
	"sync"
	"time"
)

type (
	// Sessions manage application's sessions
	Sessions struct {
		lock      sync.Mutex
		secretKey string
		sessions  []*Session
		option    Options
		block     cipher.Block
	}

	// Session session data structure
	Session struct {
		lake  map[string]interface{}
		plane []byte
	}
	// Options session options
	Options struct {
		sessLength int64
		MaxAge     int64
	}
)

var sessions *Sessions

// NewSessions create sessions instance
func NewSessions(secretKey string) {
	aesBlock, _ := aes.NewCipher(([]byte)(secretKey))
	sessions = &Sessions{
		secretKey: secretKey,
		block:     aesBlock,
	}
}

// Start start session, return session instance
func Start(ctx *miyabi.Context) *Session {
	if cookie, err := ctx.Request.Base.Cookie("sessionID"); err == nil {
		sessID := cookie.Value
		if i := sessions.evaluateSessions(([]byte)(sessID)); i >= 0 {
			return sessions.sessions[i]
		}
	}
	var session *Session
	session.plane = generateSessionID(4)
	for i := sessions.evaluateSessions(session.plane); i < 0; {
		session.plane = generateSessionID(4)
	}
	session.lake = make(map[string]interface{})
	cookie := &http.Cookie{
		Name:  "sessionID",
		Value: sessions.encrypt(session.plane),
	}
	http.SetCookie(*ctx.Response.Writer, cookie)
	return session
}

// Get get key data from lake.
func (sess *Session) Get(key string) interface{} {
	return sess.lake[key]
}

// Set set key data to lake.
func (sess *Session) Set(key string, value interface{}) {
	sess.lake[key] = value
}

// Save store data sessions
func (sess *Session) Save() {
	sessions.lock.Lock()
	for i := 0; i < len(sessions.sessions); i++ {
		block := sessions.block
		dst := make([]byte, len([]byte(sess.plane)))
		block.Decrypt(dst, ([]byte)(sess.plane))
		if string(sessions.sessions[i].plane) == string(dst) {
			sessions.sessions[i] = sess
		}
	}
	sessions.lock.Unlock()
}

// Destroy clear session data.
func (sess *Session) Destroy() {
	sessions.lock.Lock()
	for i := 0; i < len(sessions.sessions); i++ {
		block := sessions.block
		dst := make([]byte, len([]byte(sess.plane)))
		block.Decrypt(dst, ([]byte)(sess.plane))
		if string(sessions.sessions[i].plane) == string(dst) {
			sessions.sessions[i] = sess
			sessions.sessions = append(sessions.sessions[:i], sessions.sessions[i+1:]...)
		}
	}
	sessions.lock.Unlock()
}

func (sess *Sessions) evaluateSessions(sessID []byte) int {
	sess.lock.Lock()
	for i := 0; i < len(sess.sessions); i++ {
		block := sess.block
		dst := make([]byte, len([]byte(sessID)))
		block.Decrypt(dst, ([]byte)(sessID))
		if string(sess.sessions[i].plane) == string(dst) {
			sess.lock.Unlock()
			return i
		}
	}
	sess.lock.Unlock()
	return -1
}

func (sess *Sessions) encrypt(plane []byte) string {
	sess.lock.Lock()
	dst := make([]byte, len(plane))
	sess.block.Encrypt(dst, plane)
	sess.lock.Unlock()
	return string(dst)
}

func generateSessionID(length int) []byte {
	rand.Seed(time.Now().UnixNano())
	base := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789[+:/,.-~_]"
	result := ""
	for i := 0; i < length*16; i++ {
		result += string(base[rand.Intn(len(base))])
	}
	return ([]byte)(result)
}
