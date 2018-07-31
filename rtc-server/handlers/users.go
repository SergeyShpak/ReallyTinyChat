package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"

	"github.com/SergeyShpak/ReallyTinyChat/rtc-server/errors"
	"github.com/SergeyShpak/ReallyTinyChat/rtc-server/types"
)

func (h *Handler) StoreUser(ws *websocket.Conn, u *types.User, secret []byte) error {
	key := h.getCacheKey(u)
	if err := h.cache.Set(key, secret); err != nil {
		return err
	}
	users.Store(ws, u)
	return nil
}

func (h *Handler) GetUserWithConn(ws *websocket.Conn) (*types.User, []byte, error) {
	u, err := h.getUserFromConnection(ws)
	if err != nil {
		return nil, nil, err
	}
	secret, err := h.GetUserSecret(u)
	if err != nil {
		return nil, nil, err
	}
	return u, secret, nil
}

func (h *Handler) GetUserSecret(u *types.User) ([]byte, error) {
	key := h.getCacheKey(u)
	var secret []byte
	if err := h.cache.Get(key, &secret); err != nil {
		return nil, err
	}
	return secret, nil
}

func (h *Handler) GetUserSecretWithConn(ws *websocket.Conn) ([]byte, error) {
	_, secret, err := h.GetUserWithConn(ws)
	return secret, err
}

func (h *Handler) RemoveUserWithConn(ws *websocket.Conn) error {
	u, err := h.getUserFromConnection(ws)
	if err != nil {
		return err
	}
	key := h.getCacheKey(u)
	if err := h.cache.Remove(key); err != nil {
		return err
	}
	users.Delete(ws)
	return nil
}

func (h *Handler) getCacheKey(u *types.User) string {
	return strings.Join([]string{"userSecret:", u.Login, "@", u.Room}, "")
}

func (h *Handler) getUserFromConnection(ws *websocket.Conn) (*types.User, error) {
	uIface, ok := users.Load(ws)
	if !ok {
		return nil, errors.NewServerError(http.StatusNotFound, "no user is associated with the current connection")
	}
	u, ok := uIface.(*types.User)
	if !ok {
		return nil, errors.NewServerError(http.StatusInternalServerError, fmt.Sprintf("could not cast %v to *types.User", uIface))
	}
	return u, nil
}
