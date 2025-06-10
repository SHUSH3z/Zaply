package main

import (
	"archive/zip"
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

func (a *App) SelectFolder() (string, error) {
	selectedDirectory, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		Title:                "Selecione uma Pasta",
		CanCreateDirectories: true,
		ShowHiddenFiles:      true,
	})

	if err != nil {
		fmt.Printf("Erro ao abrir diálogo de seleção de pasta: %s\n", err)
		return "", fmt.Errorf("erro ao selecionar pasta: %w", err)
	}

	if selectedDirectory == "" {
		fmt.Println("Seleção de pasta cancelada pelo usuário.")
		return "", nil
	}

	fmt.Printf("Pasta selecionada: %s\n", selectedDirectory)
	return selectedDirectory, nil
}

type Ponto struct {
	ID         int     `json:"id"`
	Coordenada string  `json:"coordenada"`
	Lat        float64 `json:"lat"`
	Long       float64 `json:"long"`
	Endereco   string  `json:"endereco"`
}

type KMZRequest struct {
	Pasta  string  `json:"pasta"`
	Pontos []Ponto `json:"pontos"`
}

func (a *App) GerarKMZ(req KMZRequest) (string, error) {
	kml := `<?xml version="1.0" encoding="UTF-8"?>
<kml xmlns="http://www.opengis.net/kml/2.2">
  <Document>
    <name>Mapa Gerado</name>
`

	var path []string

	for _, ponto := range req.Pontos {
		kml += fmt.Sprintf(`
      <Placemark>
        <name>%d</name>
        <description><![CDATA[%s<br/><img src="%d.jpg" width="300"/>]]></description>
        <Point><coordinates>%f,%f</coordinates></Point>
      </Placemark>
		`, ponto.ID, ponto.Endereco, ponto.ID, ponto.Long, ponto.Lat)

		path = append(path, fmt.Sprintf("%f,%f,0", ponto.Long, ponto.Lat))
	}

	// Linha conectando os pontos
	kml += `
    <Placemark>
      <name>Rota</name>
      <LineString>
        <coordinates>
` + strings.Join(path, " ") + `
        </coordinates>
      </LineString>
    </Placemark>
`

	kml += `</Document></kml>`

	// Criar diretório temporário para o KMZ
	tmpDir := "temp_kmz"
	os.MkdirAll(tmpDir, 0755)

	kmlPath := filepath.Join(tmpDir, "doc.kml")
	err := os.WriteFile(kmlPath, []byte(kml), 0644)
	if err != nil {
		return "", fmt.Errorf("erro ao salvar KML: %v", err)
	}

	// Copiar imagens para a pasta temporária
	for _, ponto := range req.Pontos {
		src := filepath.Join(req.Pasta, fmt.Sprintf("%d.jpg", ponto.ID))
		dst := filepath.Join(tmpDir, fmt.Sprintf("%d.jpg", ponto.ID))

		data, err := os.ReadFile(src)
		if err != nil {
			fmt.Printf("Erro ao copiar imagem %s: %v\n", src, err)
			continue // não bloqueia a geração por causa de uma imagem
		}
		os.WriteFile(dst, data, 0644)
	}

	// ⚠️ Verifica se a pasta 'uploads' existe
	uploadsDir := "uploads"
	if _, err := os.Stat(uploadsDir); os.IsNotExist(err) {
		err := os.MkdirAll(uploadsDir, 0755)
		if err != nil {
			return "", fmt.Errorf("erro ao criar diretório uploads: %v", err)
		}
	}

	// Compactar em .kmz
	kmzPath := filepath.Join(uploadsDir, "mapa.kmz")
	err = criarZip(kmzPath, tmpDir)
	if err != nil {
		return "", fmt.Errorf("erro ao criar KMZ: %v", err)
	}

	return kmzPath, nil
}
func criarZip(zipPath, folder string) error {
	out, err := os.Create(zipPath)
	if err != nil {
		return err
	}
	defer out.Close()

	w := zip.NewWriter(out)
	defer w.Close()

	return filepath.Walk(folder, func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		rel, _ := filepath.Rel(folder, path)
		f, err := w.Create(rel)
		if err != nil {
			return err
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		_, err = f.Write(data)
		return err
	})
}
func garantirDiretorioUploads() error {
	path := filepath.Join("uploads")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return os.MkdirAll(path, 0755)
	}
	return nil
}
