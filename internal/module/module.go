package module

import (
	"maps"
	"net/url"
	"regexp"
	"strings"

	service "github.com/ldmonster/tts-parser/internal"

	"github.com/Masterminds/semver"
)

type FileMapping[T comparable] map[string]T

func (fm FileMapping[T]) Merge(input FileMapping[T]) {
	maps.Copy(fm, input)
}

type ModuleFile struct {
	URL       string
	Type      service.FileType
	Extension string
}

func (mf ModuleFile) GetFilename() string {
	return FileNameFromURL(mf.URL)
}

func (mf ModuleFile) GetExtension() string {
	switch mf.Type {
	case service.FileTypeAsset:
		return ".unity3d"
	case service.FileTypeModel:
		return ".obj"
	case service.FileTypePDF:
		return ".PDF"
	default:
		return mf.Extension
	}
}

func (mf ModuleFile) GetFolder() string {
	switch mf.Type {
	case service.FileTypeAsset:
		return "Assetbundles"
	case service.FileTypeModel:
		return "Models"
	case service.FileTypeImage:
		return "Images"
	case service.FileTypePDF:
		return "PDF"
	case service.FileTypeAudio:
		return "Audio"
	default:
		return ""
	}
}

func NewTTSModule() *TTSModule {
	return &TTSModule{
		Assets: make(FileMapping[ModuleFile], 0),
		Models: make(FileMapping[ModuleFile], 0),
		Images: make(FileMapping[ModuleFile], 0),
		PDFs:   make(FileMapping[ModuleFile], 0),
		Audio:  make(FileMapping[ModuleFile], 0),
	}
}

type TTSModule struct {
	Assets FileMapping[ModuleFile] // .unity3d
	Models FileMapping[ModuleFile] // .obj
	Images FileMapping[ModuleFile]
	PDFs   FileMapping[ModuleFile] // .PDF
	Audio  FileMapping[ModuleFile]

	ID uint

	Name          string
	EpochTime     uint
	VersionNumber *semver.Version
}

func (m *TTSModule) Merge(input *TTSModule) {
	m.Assets.Merge(input.Assets)
	m.Models.Merge(input.Models)
	m.Images.Merge(input.Images)
	m.PDFs.Merge(input.PDFs)
	m.Audio.Merge(input.Audio)
}

// return orphans
func (m *TTSModule) MergeFiles(input []service.File) []service.File {
	orphans := make([]service.File, 0, 1)
	for _, f := range input {
		newf := ModuleFile{
			URL:       f.URL,
			Type:      f.Type,
			Extension: f.Extension,
		}

		switch f.Type {
		case service.FileTypeAsset:
			_, ok := m.Assets[f.URL]
			if !ok {
				orphans = append(orphans, f)
				break
			}

			m.Assets[f.URL] = newf
		case service.FileTypeModel:
			_, ok := m.Models[f.URL]
			if !ok {
				orphans = append(orphans, f)
				break
			}

			m.Models[f.URL] = newf
		case service.FileTypeImage:
			_, ok := m.Images[f.URL]
			if !ok {
				orphans = append(orphans, f)
				break
			}

			m.Images[f.URL] = newf
		case service.FileTypePDF:
			_, ok := m.PDFs[f.URL]
			if !ok {
				orphans = append(orphans, f)
				break
			}

			m.PDFs[f.URL] = newf
		case service.FileTypeAudio:
			_, ok := m.Audio[f.URL]
			if !ok {
				orphans = append(orphans, f)
				break
			}

			m.Audio[f.URL] = newf
		default:
			panic("type is empty")
		}
	}

	return orphans
}

func (m *TTSModule) AddAsset(url string) {
	url = FixURL(url)
	m.Assets[url] = ModuleFile{URL: url, Type: service.FileTypeAsset}
}

func (m *TTSModule) AddModel(url string) {
	url = FixURL(url)
	m.Models[url] = ModuleFile{URL: url, Type: service.FileTypeModel}
}

func (m *TTSModule) AddImage(url string) {
	url = FixURL(url)
	m.Images[url] = ModuleFile{URL: url, Type: service.FileTypeImage}
}

func (m *TTSModule) AddPDF(url string) {
	url = FixURL(url)
	m.PDFs[url] = ModuleFile{URL: url, Type: service.FileTypePDF}
}

func (m *TTSModule) AddAudio(url string) {
	url = FixURL(url)
	m.Audio[url] = ModuleFile{URL: url, Type: service.FileTypeAudio}
}

func (m *TTSModule) AddAssetsBundle(b *CustomAssetbundle) {
	if b == nil {
		return
	}

	m.AddAsset(b.AssetbundleURL)
}

func (m *TTSModule) AddMesh(b *CustomMesh) {
	if b == nil {
		return
	}

	m.AddImage(b.DiffuseURL)
	m.AddModel(b.MeshURL)
	m.AddModel(b.ColliderURL)
}

func (m *TTSModule) AddUIAssets(b CustomUIAssets) {
	if b == nil {
		return
	}

	for _, asset := range b {
		m.AddImage(asset.URL)
	}
}

func (m *TTSModule) AddImages(b *CustomImage) {
	if b == nil {
		return
	}

	m.AddImage(b.ImageURL)

	if b.ImageSecondaryURL != "" {
		m.AddImage(b.ImageSecondaryURL)
	}
}

func (m *TTSModule) AddPDFs(b *CustomPDF) {
	if b == nil {
		return
	}

	m.AddPDF(b.PDFURL)
}

func (m *TTSModule) AddDecals(b AttachedDecals) {
	if b == nil {
		return
	}

	for _, decal := range b {
		if decal.CustomDecal == nil {
			continue
		}

		m.AddImage(decal.CustomDecal.ImageURL)
	}
}

func (m *TTSModule) AddDeck(b CustomDeck) {
	if b == nil {
		return
	}

	for _, card := range b {
		m.AddImage(card.FaceURL)
		m.AddImage(card.BackURL)
	}
}

func (m *TTSModule) BatchAdd(state Object) {
	m.AddAssetsBundle(state.CustomAssetbundle)
	m.AddImages(state.CustomImage)
	m.AddMesh(state.CustomMesh)
	m.AddDeck(state.CustomDeck)
	m.AddUIAssets(state.CustomUIAssets)
	m.AddPDFs(state.CustomPDF)
	m.AddDecals(state.AttachedDecals)
}

func (m *TTSModule) GetAll() FileMapping[ModuleFile] {
	allLength := len(m.Assets) + len(m.Models) + len(m.Images) + len(m.PDFs) + len(m.Audio)

	all := make(FileMapping[ModuleFile], allLength)

	maps.Copy(all, m.Assets)
	maps.Copy(all, m.Models)
	maps.Copy(all, m.Images)
	maps.Copy(all, m.PDFs)
	maps.Copy(all, m.Audio)

	return all
}

func (m *TTSModule) Contains(url string) bool {
	url = FixURL(url)

	_, ok := m.GetAll()[url]

	return ok
}

func (m *TTSModule) ScanModule(mod *Module) {
	m.Name = mod.SaveName
	if mod.VersionNumber != "" {
		m.VersionNumber = semver.MustParse(mod.VersionNumber)
	} else {
		m.VersionNumber = semver.MustParse("0")
	}

	if mod.TableURL != "" {
		m.AddImage(mod.TableURL)
	}

	if mod.SkyURL != "" {
		m.AddImage(mod.SkyURL)
	}

	if mod.MusicPlayer != nil {
		if len(mod.MusicPlayer.AudioLibrary) > 0 {
			for _, audio := range mod.MusicPlayer.AudioLibrary {
				for _, val := range audio {
					url, err := url.Parse(val)
					if err == nil && (url.Scheme == "http" || url.Scheme == "https") {
						m.AddAudio(val)
					}
				}
			}
		}
	}

	m.AddUIAssets(mod.CustomUIAssets)

	for _, state := range mod.Objects {
		m.Merge(scanObject(state))
	}
}

func scanObject(state Object) *TTSModule {
	result := NewTTSModule()

	for _, state := range state.ContainedObjects {
		result.BatchAdd(state)

		result.Merge(scanObject(state))
	}

	for _, state := range state.States {
		result.BatchAdd(state)

		result.Merge(scanObject(state))
	}

	for _, state := range state.ChildObjects {
		result.BatchAdd(state)

		result.Merge(scanObject(state))
	}

	result.BatchAdd(state)

	return result
}

var urlStartRegex = regexp.MustCompile(`^http.*$`)

// they replace all cloud-3 links to akamaihd
// e.g.
// http://cloud-3.steamusercontent.com/ugc/2039602457422068571/78FE30D056AE0AF74047A9F4A8B68B481062F77B/
// to
// https://steamusercontent-a.akamaihd.net/ugc/2039602457422068571/78FE30D056AE0AF74047A9F4A8B68B481062F77B/
//
// also add http:// if scheme not found
func FixURL(url string) string {
	if url == "" {
		return ""
	}

	if !urlStartRegex.MatchString(url) {
		url = "http://" + url
	}

	return strings.Replace(url, `http://cloud-3.steamusercontent.com`, `https://steamusercontent-a.akamaihd.net`, 1)
}

var nonDigitRegex = regexp.MustCompile(`\W`)

func FileNameFromURL(url string) string {
	return nonDigitRegex.ReplaceAllString(url, "")
}
