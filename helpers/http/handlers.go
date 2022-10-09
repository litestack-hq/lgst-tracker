package http

import (
	"net/http"
)

func (m *HttpModule) Ping(w http.ResponseWriter, r *http.Request) {
	// req := LoginReqest{}
	// if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
	// 	msg := "Failed to parse request"
	// 	m.Logger.Warn().Err(err).Msg(msg)
	// }

	// req.EmailAddress = utils.NormalizeEmail(req.EmailAddress)
	// req.Password = strings.TrimSpace(req.Password)

	// if err := m.Validator.Validate.Struct(req); err != nil {
	// 	if _, ok := err.(*validator.InvalidValidationError); ok {
	// 		msg := "validation failed to run"
	// 		m.Logger.Err(err).Msg(msg)
	// 		render.Status(r, http.StatusInternalServerError)
	// 		render.JSON(w, r, LoginResponse{
	// 			Message: msg,
	// 		})
	// 		return
	// 	}

	// 	msg := "Input validation failed."
	// 	render.Status(r, http.StatusBadRequest)
	// 	render.JSON(w, r, LoginResponse{
	// 		Message: msg,
	// 		Error:   m.Validator.TranslateErrors(err),
	// 	})
	// 	return
	// }

	// user, err := m.DataStore.FindUserByEmailAddress(req.EmailAddress)
	// if err != nil {
	// 	m.Logger.Error().Err(err).Msg("find user by email failed")
	// }

	// if user == nil {
	// 	render.Status(r, http.StatusUnauthorized)
	// 	render.JSON(w, r, LoginResponse{
	// 		Message: "User not found.",
	// 	})
	// 	return
	// }

	// if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
	// 	render.Status(r, http.StatusUnauthorized)
	// 	render.JSON(w, r, LoginResponse{
	// 		Message: "Password mismatch.",
	// 	})
	// 	return
	// }

	// token, err := jwt.CreateToken(user.Id, m.Conf.APP_KEY)
	// if err != nil {
	// 	msg := "JWT generation failed"
	// 	m.Logger.Err(err).Msg(msg)
	// 	render.Status(r, http.StatusInternalServerError)
	// 	render.JSON(w, r, LoginResponse{
	// 		Message: msg,
	// 	})
	// 	return
	// }

	// render.JSON(w, r, LoginResponse{
	// 	Message: "Login successful.",
	// 	Token:   &token,
	// 	User:    user,
	// })
}
