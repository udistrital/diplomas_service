package routers

import (
	beego "github.com/beego/beego/v2/server/web"

	"github.com/udistrital/diplomas_mid/controllers"
)

func init() {
	api := beego.NewNamespace("/v1",
		beego.NSNamespace("/documento_digital",
			beego.NSRouter("/:id/firmar", &controllers.FirmaDiplomaController{}, "post:Firmar"),
		),
		beego.NSNamespace("/diplomas",
			beego.NSRouter("/estudiantes-aprobados", &controllers.EstudianteGradoController{}, "get:AprobadosPorFacultad"),
			beego.NSRouter("/estudiantes-aprobados/documentos", &controllers.EstudianteGradoController{}, "post:CrearDocumentosAprobados"),
		),
		beego.NSNamespace("/firmantes",
			beego.NSRouter("/rol-activo", &controllers.FirmanteController{}, "get:RolActivo"),
			beego.NSRouter("/firma", &controllers.FirmaFirmanteController{}, "post:SubirFirma"),
		),
	)

	beego.AddNamespace(api)
}
