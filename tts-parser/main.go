package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"regexp"
	"syscall"
	"tts/internal/module"
	"tts/internal/zap"

	uberzap "go.uber.org/zap"
)

func main() {
	// // url := `https://drive.google.com/uc?export=download&id=1jEnm1nhmTG-8c1aKUCIFwA6S5jKWA3AS`
	// url := `https://drive.google.com/uc?export=download&id=0B6g_wRdiFw9TVGlha2djalJVX3M`
	// req, err := http.NewRequest(http.MethodGet, url, nil)
	// if err != nil {
	// 	panic(err)
	// }

	// resp, err := http.DefaultClient.Do(req)
	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Println(resp.StatusCode)
	// boo, err := io.ReadAll(resp.Body)
	// if err != nil {
	// 	panic(err)
	// }

	// os.WriteFile("kek.html", boo, 0777)

	// for k, v := range resp.Header {
	// 	fmt.Println(k, v)
	// }

	// return
	cfg := NewConfig()

	err := cfg.AutoLoadEnvs()
	if err != nil && (!errors.Is(err, ErrEnvFileIsNotFound) && !errors.Is(err, ErrConfigFileIsNotFound)) {
		panic(err)
	}

	if err != nil {
		fmt.Println(err)
	}

	err = cfg.Parse()
	if err != nil {
		panic(err)
	}

	logger, err := zap.NewProductionZaplogger("log.txt", cfg.LogLevel)
	if err != nil {
		panic(err)
	}

	logger = logger.Named("tts")

	defer func(logger *uberzap.Logger) {
		_ = logger.Sync()
	}(logger)

	b := NewBackend(cfg, logger)

	err = b.init()
	if err != nil {
		logger.Fatal("backend initialization", uberzap.Error(err))
	}
	logger.Error("download module", uberzap.Error(err))

	ctx, stopNotify := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

	b.Start(ctx)

	stopNotify()
}

var urlRegex = regexp.MustCompile(`"http://.*",$`)

func kek() {
	dir := `H:\Tabletop Simulator\Mods\Workshop`

	fs, err := os.ReadDir(dir)
	if err != nil {
		panic(err)
	}

	for _, fName := range fs {
		if filepath.Ext(fName.Name()) != ".json" {
			continue
		}

		fmt.Println(fName.Name())

		f, err := os.OpenFile(filepath.Join(dir, fName.Name()), os.O_RDONLY, 0o666)
		if err != nil {
			panic(err)
		}

		mod := new(module.Module)
		err = json.NewDecoder(f).Decode(mod)
		if err != nil {
			panic(err)
		}

		result := module.NewTTSModule()

		result.ScanModule(mod)

		// fmt.Println(len(result.Assets))
		// fmt.Println(len(result.Images))
		// fmt.Println(len(result.Models))
		// fmt.Println(len(result.PDFs))
		// fmt.Println(len(result.Audio))
		// fmt.Println(len(result.All))

		by, err := os.ReadFile("module.json")
		if err != nil {
			panic(err)
		}

		found := urlRegex.FindAllString(string(by), -1)

		for _, url := range found {
			if !result.Contains(url[1 : len(url)-2]) {
				url = module.FixURL(url)
				fmt.Println(url, module.FileNameFromURL(url))
			}
		}

		// for idx, img := range result.Images {
		// 	fmt.Println(idx, img)
		// }
	}
}

var idRegex = regexp.MustCompile(`.*\(([0-9]*)\).zip`)

func removeDownloaded() {
	fs, err := os.ReadDir(`I:\tabletop\TTS\BoredFull`)
	if err != nil {
		panic(err)
	}

	for _, f := range fs {
		subs := idRegex.FindAllStringSubmatch(f.Name(), -1)
		if len(subs) == 0 {
			fmt.Println("it cant be", f.Name())

			continue
		}

		err := os.Remove(fmt.Sprintf(`E:\Tabletop Simulator\Mods\Workshop\%s.json`, subs[0][1]))
		if err != nil {
			fmt.Println("not found", subs[0][1])
		}
	}
}
