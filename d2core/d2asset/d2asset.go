package d2asset

import (
	"log"

	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2data/d2datadict"
	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2fileformats/d2dat"
	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2fileformats/d2mpq"
	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2interface"

	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2fileformats/d2pl2"
	"github.com/OpenDiablo2/OpenDiablo2/d2core/d2config"
)

var singleton *assetManager

func Initialize(term d2interface.Terminal) error {
	verifyNotInit()

	var (
		config                  = d2config.Get()
		archiveManager          = createArchiveManager(config)
		fileManager             = createFileManager(config, archiveManager)
		paletteManager          = createPaletteManager()
		paletteTransformManager = createPaletteTransformManager()
		animationManager        = createAnimationManager()
		fontManager             = createFontManager()
	)

	singleton = &assetManager{
		archiveManager,
		fileManager,
		paletteManager,
		paletteTransformManager,
		animationManager,
		fontManager,
	}

	term.BindAction("assetspam", "display verbose asset manager logs", func(verbose bool) {
		if verbose {
			term.OutputInfo("asset manager verbose logging enabled")
		} else {
			term.OutputInfo("asset manager verbose logging disabled")
		}

		archiveManager.cache.SetVerbose(verbose)
		fileManager.cache.SetVerbose(verbose)
		paletteManager.cache.SetVerbose(verbose)
		paletteTransformManager.cache.SetVerbose(verbose)
		animationManager.cache.SetVerbose(verbose)
	})

	term.BindAction("assetstat", "display asset manager cache statistics", func() {
		term.OutputInfo("archive cache: %f", float64(archiveManager.cache.GetWeight())/float64(archiveManager.cache.GetBudget())*100.0)
		term.OutputInfo("file cache: %f", float64(fileManager.cache.GetWeight())/float64(fileManager.cache.GetBudget())*100.0)
		term.OutputInfo("palette cache: %f", float64(paletteManager.cache.GetWeight())/float64(paletteManager.cache.GetBudget())*100.0)
		term.OutputInfo("palette transform cache: %f", float64(paletteTransformManager.cache.GetWeight())/float64(paletteTransformManager.cache.GetBudget())*100.0)
		term.OutputInfo("animation cache: %f", float64(animationManager.cache.GetWeight())/float64(animationManager.cache.GetBudget())*100.0)
		term.OutputInfo("font cache: %f", float64(fontManager.cache.GetWeight())/float64(fontManager.cache.GetBudget())*100.0)
	})

	term.BindAction("assetclear", "clear asset manager cache", func() {
		archiveManager.cache.Clear()
		fileManager.cache.Clear()
		paletteManager.cache.Clear()
		paletteTransformManager.cache.Clear()
		animationManager.cache.Clear()
		fontManager.cache.Clear()
	})

	return nil
}

func Shutdown() {
	singleton = nil
}

func LoadArchive(archivePath string) (*d2mpq.MPQ, error) {
	verifyWasInit()
	return singleton.archiveManager.loadArchive(archivePath)
}

func LoadFileStream(filePath string) (*d2mpq.MpqDataStream, error) {
	verifyWasInit()

	data, err := singleton.fileManager.loadFileStream(filePath)
	if err != nil {
		log.Printf("error loading file stream %s (%v)", filePath, err.Error())
	}

	return data, err
}

func LoadFile(filePath string) ([]byte, error) {
	verifyWasInit()

	data, err := singleton.fileManager.loadFile(filePath)
	if err != nil {
		log.Printf("error loading file %s (%v)", filePath, err.Error())
	}

	return data, err
}

func FileExists(filePath string) (bool, error) {
	verifyWasInit()
	return singleton.fileManager.fileExists(filePath)
}

func LoadAnimation(animationPath, palettePath string) (*Animation, error) {
	verifyWasInit()
	return LoadAnimationWithTransparency(animationPath, palettePath, 255)
}

func LoadPaletteTransform(pl2Path string) (*d2pl2.PL2, error) {
	verifyWasInit()
	return singleton.paletteTransformManager.loadPaletteTransform(pl2Path)
}

func LoadAnimationWithTransparency(animationPath, palettePath string, transparency int) (*Animation, error) {
	verifyWasInit()
	return singleton.animationManager.loadAnimation(animationPath, palettePath, transparency)
}

func LoadComposite(object *d2datadict.ObjectLookupRecord, palettePath string) (*Composite, error) {
	verifyWasInit()
	return CreateComposite(object, palettePath), nil
}

func LoadFont(tablePath, spritePath, palettePath string) (*Font, error) {
	verifyWasInit()
	return singleton.fontManager.loadFont(tablePath, spritePath, palettePath)
}

func LoadPalette(palettePath string) (*d2dat.DATPalette, error) {
	verifyWasInit()
	return singleton.paletteManager.loadPalette(palettePath)
}

func verifyWasInit() {
	if singleton == nil {
		panic(ErrNotInit)
	}
}

func verifyNotInit() {
	if singleton != nil {
		panic(ErrWasInit)
	}
}
