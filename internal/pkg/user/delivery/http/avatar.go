package http

import (
	"Redioteka/internal/pkg/domain"
	"Redioteka/internal/pkg/user"
	"Redioteka/internal/pkg/utils/jsonerrors"
	"Redioteka/internal/pkg/utils/log"
	"Redioteka/internal/pkg/utils/session"
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/gorilla/mux"
)

const (
	root    = "./static"
	urlRoot = "https://redioteka.com/static"
	path    = "/media/img/users/"
)

//Avatar - handler for uploading user avatar
func (handler *UserHandler) Avatar(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	vars := mux.Vars(r)
	idString := vars["id"]
	urlID, err := strconv.Atoi(idString)
	if err != nil {
		log.Log.Warn(fmt.Sprintf("Cannot get user id: %s", err))
		http.Error(w, jsonerrors.URLParams, http.StatusBadRequest)
		return
	}

	sess, err := session.GetSession(r)
	if err != nil || session.Manager.Check(sess) != nil {
		http.Error(w, jsonerrors.Session, http.StatusUnauthorized)
		return
	}
	if uint(urlID) != sess.UserID {
		log.Log.Warn("User try update wrong avatar")
		http.Error(w, jsonerrors.JSONMessage("wrong id"), http.StatusForbidden)
		return
	}

	r.ParseMultipartForm(5 * 1024 * 1024)
	uploaded, header, err := r.FormFile("avatar")
	if err != nil {
		log.Log.Warn(fmt.Sprintf("error while parsing file: %s", err))
		http.Error(w, jsonerrors.JSONMessage("file parsing"), http.StatusBadRequest)
	}
	defer uploaded.Close()

	filename, err := handler.UUsecase.UploadAvatar(uploaded, path, filepath.Ext(header.Filename))
	if err != nil {
		log.Log.Error(err)
		http.Error(w, jsonerrors.JSONMessage("upload"), http.StatusInternalServerError)
	}

	err = handler.UUsecase.Update(&domain.User{
		ID:     sess.UserID,
		Avatar: filename,
	})
	if err != nil {
		http.Error(w, jsonerrors.JSONMessage("update fail"), user.CodeFromError(err))
		return
	}
	if err := json.NewEncoder(w).Encode(map[string]string{"user_avatar": filename}); err != nil {
		log.Log.Error(err)
		http.Error(w, jsonerrors.JSONEncode, http.StatusInternalServerError)
		return

	}
}
