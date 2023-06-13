package routers

import (
	"context"
	"encoding/json"

	"github.com/diegovillarino/go/tree/victor_user/database"
	"github.com/diegovillarino/go/tree/victor_user/models"
)

/*ModificarPerfil modifica el perfil de usuario */
func ModificarPerfil(ctx context.Context, claim models.Claim) models.RespApi {

	var r models.RespApi
	r.Status = 400

	var t models.User

	body := ctx.Value(models.Key("body")).(string)
	err := json.Unmarshal([]byte(body), &t)
	if err != nil {
		r.Message = "Datos Incorrectos " + err.Error()
		return r
	}

	var status bool

	status, err = database.ModificoRegistro(t, claim.ID.Hex())
	if err != nil {
		r.Message = "Ocurrió un error al intentar modificar el registro. Reintente nuevamente " + err.Error()
		return r
	}

	if !status {
		r.Message = "No se ha logrado modificar el registro del usuario "
		return r
	}

	r.Status = 200
	r.Message = "Modificar Perfil OK !"
	return r
}
