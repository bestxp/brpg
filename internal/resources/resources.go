package resources

//go:generate pkger -o internal/resources

import (
	"bytes"
	"image"
	"image/png"
	"io"
	"os"

	_ "github.com/bestxp/brpg"
	"github.com/markbates/pkger"
)

type Frames struct {
	Frames []image.Image
	image.Config
}

type Audio struct {
	Stream io.ReadSeeker
	Len    int
}

func LoadAudios() (map[string]*Audio, error) {
	out := map[string]*Audio{}

	prefix := "/resources/audio"

	err := pkger.Walk(prefix, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() == false {
			filename := prefix + "/" + info.Name()
			file, err := pkger.Open(filename)
			if err != nil {
				return err
			}
			defer file.Close()

			byt, err := io.ReadAll(file)
			if err != nil {
				return err
			}
			out[filename] = &Audio{
				Stream: bytes.NewReader(byt),
				Len:    len(byt),
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}
	return out, nil
}

func LoadResources() (map[string]Frames, error) {
	images := map[string]image.Image{}
	cfgs := map[string]image.Config{}
	sprites := map[string]Frames{}

	prefix := "/resources/sprites"
	err := pkger.Walk(prefix, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() == false {
			filename := prefix + "/" + info.Name()
			file, err := pkger.Open(filename)
			if err != nil {
				return err
			}
			defer file.Close()

			img, err := png.Decode(file)
			if err != nil {
				return err
			}

			fileCfg, err := pkger.Open(filename)
			if err != nil {
				return err
			}
			defer fileCfg.Close()

			cfg, err := png.DecodeConfig(fileCfg)
			if err != nil {
				return err
			}

			images[info.Name()] = img
			cfgs[info.Name()] = cfg
		}

		return nil
	})
	if err != nil {
		return sprites, err
	}

	sprites["big_demon_idle"] = Frames{
		Frames: []image.Image{
			images["big_demon_idle_anim_f0.png"],
			images["big_demon_idle_anim_f1.png"],
			images["big_demon_idle_anim_f2.png"],
			images["big_demon_idle_anim_f3.png"],
		},
		Config: cfgs["big_demon_idle_anim_f0.png"],
	}
	sprites["big_demon_run"] = Frames{
		Frames: []image.Image{
			images["big_demon_run_anim_f0.png"],
			images["big_demon_run_anim_f1.png"],
			images["big_demon_run_anim_f2.png"],
			images["big_demon_run_anim_f3.png"],
		},
		Config: cfgs["big_demon_run_anim_f0.png"],
	}

	sprites["big_zombie_idle"] = Frames{
		Frames: []image.Image{
			images["big_zombie_idle_anim_f0.png"],
			images["big_zombie_idle_anim_f1.png"],
			images["big_zombie_idle_anim_f2.png"],
			images["big_zombie_idle_anim_f3.png"],
		},
		Config: cfgs["big_zombie_idle_anim_f0.png"],
	}
	sprites["big_zombie_run"] = Frames{
		Frames: []image.Image{
			images["big_zombie_run_anim_f0.png"],
			images["big_zombie_run_anim_f1.png"],
			images["big_zombie_run_anim_f2.png"],
			images["big_zombie_run_anim_f3.png"],
		},
		Config: cfgs["big_zombie_run_anim_f0.png"],
	}

	sprites["floor_1"] = Frames{
		Frames: []image.Image{images["floor_1.png"]},
		Config: cfgs["floor_1.png"],
	}
	sprites["floor_2"] = Frames{
		Frames: []image.Image{images["floor_2.png"]},
		Config: cfgs["floor_2.png"],
	}
	sprites["floor_3"] = Frames{
		Frames: []image.Image{images["floor_3.png"]},
		Config: cfgs["floor_3.png"],
	}
	sprites["floor_4"] = Frames{
		Frames: []image.Image{images["floor_4.png"]},
		Config: cfgs["floor_4.png"],
	}
	sprites["floor_5"] = Frames{
		Frames: []image.Image{images["floor_5.png"]},
		Config: cfgs["floor_5.png"],
	}
	sprites["floor_6"] = Frames{
		Frames: []image.Image{images["floor_6.png"]},
		Config: cfgs["floor_6.png"],
	}
	sprites["floor_7"] = Frames{
		Frames: []image.Image{images["floor_7.png"]},
		Config: cfgs["floor_7.png"],
	}
	sprites["floor_8"] = Frames{
		Frames: []image.Image{images["floor_8.png"]},
		Config: cfgs["floor_8.png"],
	}
	sprites["floor_ladder"] = Frames{
		Frames: []image.Image{images["floor_ladder.png"]},
		Config: cfgs["floor_ladder.png"],
	}

	return sprites, nil
}
