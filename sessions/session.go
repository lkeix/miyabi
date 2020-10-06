package sessions

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"fmt"
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
		iv        []byte
		sessions  []*Session
		options   *Options
		block     cipher.Block
	}

	// Session session data structure
	Session struct {
		lake  map[string]interface{}
		plane []byte
		last  time.Time
	}

	// Options session options
	Options struct {
		HTTPOnly bool
		SameSite http.SameSite
		MaxAge   int
	}
)

var sessions *Sessions

// NewSessions create sessions instance
func NewSessions() {
	secretKey := string(generateString(2 * 16))
	aesBlock, err := aes.NewCipher(([]byte)(secretKey))
	if err != nil {
		panic(err)
	}
	sessions = &Sessions{
		secretKey: secretKey,
		block:     aesBlock,
		iv:        make([]byte, aesBlock.BlockSize()),
	}
	sessions.iv = ([]byte)(generateString(sessions.block.BlockSize()))
}

// SetOptions set options
func SetOptions(opts *Options) {
	sessions.options = opts
}

// Start start session, return session instance
func Start(ctx *miyabi.Context) *Session {
	if cookie, err := ctx.Request.Base.Cookie("sess"); err == nil {
		sessID := cookie.Value
		if i := sessions.varifySession(sessID); i >= 0 {
			return sessions.sessions[i]
		}
	}
	session := newSession(3)
	cookie := &http.Cookie{
		Name:     "sess",
		Value:    sessions.encrypt(session.plane),
		Secure:   ctx.IsTSL,
		SameSite: sessions.options.SameSite,
		HttpOnly: sessions.options.HTTPOnly,
		MaxAge:   sessions.options.MaxAge,
	}
	http.SetCookie(*ctx.Response.Writer, cookie)
	sessions.sessions = append(sessions.sessions, session)
	return session
}

func newSession(length int) *Session {
	session := &Session{
		plane: make([]byte, length*16),
		lake:  make(map[string]interface{}),
		last:  time.Now(),
	}
	session.plane = generateString(length * 16)
	for i := sessions.evaluateSessions(session.plane); i >= 0; {
		session.plane = generateString(length * 16)
	}
	session.lake = make(map[string]interface{})
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

// Destroy clear session data.
func (sess *Session) Destroy() {
	sessions.lock.Lock()
	for i := 0; i < len(sessions.sessions); i++ {
		if string(sessions.sessions[i].plane) == string(sess.plane) {
			sessions.sessions = append(sessions.sessions[:i], sessions.sessions[i+1:]...)
		}
	}
	sessions.lock.Unlock()
}

func (sess *Sessions) evaluateSessions(sessID []byte) int {
	sess.lock.Lock()
	for i := 0; i < len(sess.sessions); i++ {
		if string(sess.sessions[i].plane) == string(sessID) {
			sess.lock.Unlock()
			return i
		}
	}
	sess.lock.Unlock()
	return -1
}

func (sess *Sessions) varifySession(sessID string) int {
	for i := 0; i < len(sess.sessions); i++ {
		if string(sess.sessions[i].plane) == sess.decrypt(sessID) {
			return i
		}
	}
	return -1
}

func (sess *Sessions) encrypt(plane []byte) string {
	sess.lock.Lock()
	dst := make([]byte, len(plane))
	cbc := cipher.NewCBCEncrypter(sess.block, sess.iv)
	cbc.CryptBlocks(dst, plane)
	sess.lock.Unlock()
	return hex.EncodeToString(dst)
}

func (sess *Sessions) decrypt(encryptStr string) string {
	encrypt, _ := hex.DecodeString(encryptStr)
	cbc := cipher.NewCBCDecrypter(sess.block, sess.iv)
	dst := make([]byte, len(encrypt))
	cbc.CryptBlocks(dst, encrypt)
	return string(dst)
}

func generateString(length int) []byte {
	rand.Seed(time.Now().UnixNano())
	base := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789[+:/,.-~_];."
	result := ""
	for i := 0; i < length; i++ {
		result += string(base[rand.Intn(len(base))])
	}
	return ([]byte)(result)
}

func destroyOldSession() {
	for i := 0; i < len(sessions.sessions); i++ {
		now := time.Now()
		d := now.Sub(sessions.sessions[i].last)
		dif := int(d.Hours())*60*60 + int(d.Minutes())*60 + int(d.Seconds())
		fmt.Println(dif)
		if dif > sessions.options.MaxAge {
			sessions.sessions = append(sessions.sessions[:i], sessions.sessions[i+1:]...)
		}
	}
}
