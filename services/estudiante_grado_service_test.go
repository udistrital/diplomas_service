package services

import "testing"

func TestAgruparEstudiantesAprobadosPorFacultad(t *testing.T) {
	t.Parallel()

	result := agruparEstudiantesAprobadosPorFacultad([]EstudianteAprobadoGrado{
		{ID: 1, CodigoEstudiante: 20172007065, IDFacultadOikos: 14},
		{ID: 2, CodigoEstudiante: 20241117013, IDFacultadOikos: 14},
		{ID: 3, CodigoEstudiante: 20182395011, IDFacultadOikos: 2},
	})

	if result.Total != 3 {
		t.Fatalf("expected total 3, got %d", result.Total)
	}
	if len(result.Facultades) != 2 {
		t.Fatalf("expected 2 facultades, got %d", len(result.Facultades))
	}
	if result.Facultades[0].FacultadID != 2 || result.Facultades[0].Total != 1 {
		t.Fatalf("unexpected first group: %#v", result.Facultades[0])
	}
	if result.Facultades[1].FacultadID != 14 || result.Facultades[1].Total != 2 {
		t.Fatalf("unexpected second group: %#v", result.Facultades[1])
	}
}

func TestEstudianteAprobadoGradoDocumentoDigitalRequest(t *testing.T) {
	t.Parallel()

	estudiante := EstudianteAprobadoGrado{
		CodigoEstudiante: 20172007065,
		IDProyectoOikos:  73,
		AnioInsGrado:     2025,
		PerInsGrado:      1,
	}

	result := estudiante.DocumentoDigitalRequest(38, 1)

	if result.CodigoEstudiante != 20172007065 {
		t.Fatalf("expected codigo_estudiante, got %d", result.CodigoEstudiante)
	}
	if result.ProgramaAcademicoID != 73 {
		t.Fatalf("expected programa_academico_id from IdProyectoOikos, got %d", result.ProgramaAcademicoID)
	}
	if result.Vigencia != 2025 {
		t.Fatalf("expected vigencia from AnioInsGrado, got %d", result.Vigencia)
	}
	if result.PeriodoID != 1 {
		t.Fatalf("expected periodo_id from PerInsGrado, got %d", result.PeriodoID)
	}
	if !result.Activo {
		t.Fatal("expected activo true")
	}
}

func TestFiltrarEstudiantesAprobadosPorCodigo(t *testing.T) {
	t.Parallel()

	estudiantes := []EstudianteAprobadoGrado{
		{ID: 1, CodigoEstudiante: 20172007065},
		{ID: 2, CodigoEstudiante: 20161020507},
		{ID: 3, CodigoEstudiante: 20182395011},
	}

	result := filtrarEstudiantesAprobadosPorCodigo(estudiantes, 20161020507)

	if len(result) != 1 {
		t.Fatalf("expected one estudiante, got %d", len(result))
	}
	if result[0].CodigoEstudiante != 20161020507 {
		t.Fatalf("expected codigo 20161020507, got %d", result[0].CodigoEstudiante)
	}
}
