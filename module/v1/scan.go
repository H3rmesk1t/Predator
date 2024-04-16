package v1

import (
	"Predator/module/v1/gopoc"
	"Predator/module/v1/utils"
	ymlPoc "Predator/module/v1/ymlpoc"
	"Predator/module/v1/ymlpoc/lib"
	"Predator/module/v1/ymlpoc/structs"
	"Predator/pkg/config"
	pkgUtils "Predator/pkg/utils"
	"Predator/pkg/xhttp"
	"embed"
	"fmt"
	"github.com/corpix/uarand"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type Task struct {
	Req *xhttp.Request
	Poc *structs.Poc
}

//go:embed files
var FS embed.FS
var once sync.Once
var AllPocs []*structs.Poc

func WebScan(info *config.HostInfo) {
	once.Do(initpoc)
	var pocinfo = config.PocInfo
	buf := strings.Split(info.Url, "/")
	pocinfo.Target = strings.Join(buf[:3], "/")

	if pocinfo.PocName != "" {
		YmlPocCheck(pocinfo)
	} else {
		for _, infostr := range info.InfoStr {
			pocinfo.PocName = CheckInfoPoc(infostr)
			if pocinfo.PocName != "" {
				gopoc.GoPocCheck(info.Url, pocinfo.PocName)
			} else {
				YmlPocCheck(pocinfo)
			}
		}
	}
}

func YmlPocCheck(PocInfo config.PocInfoStruct) {
	lib.InitReversePlatform(config.CeyeApi, config.CeyeDomain)
	TotalReqeusts := 0
	for _, poc := range AllPocs {
		ruleLens := len(poc.Rules)
		if poc.Transport == "tcp" || poc.Transport == "udp" {
			ruleLens += 1
		}
		TotalReqeusts += 1 * ruleLens
	}
	if TotalReqeusts == 0 {
		TotalReqeusts = 1
	}
	lib.InitCache(TotalReqeusts)
	YmlPocCheckStart(PocInfo, AllPocs)
}

func YmlPocCheckStart(PocInfo config.PocInfoStruct, pocs []*structs.Poc) {
	req, err := xhttp.NewRequest("GET", PocInfo.Target, nil)
	if err != nil {
		errlog := fmt.Sprintf("[-] web pocs init %v %v\n", PocInfo.Target, err)
		pkgUtils.LogError(errlog)
		return
	}
	req.GetHeaders().Set("User-agent", uarand.GetRandom())
	req.GetHeaders().Set("Accept", config.Accept)
	req.GetHeaders().Set("Accept-Language", "zh-CN,zh;q=0.9,en-US;q=0.8,en;q=0.7")
	if config.Cookie != "" {
		req.GetHeaders().Set("Cookie", config.Cookie)
	}
	req.GetHeaders().Set("Connection", "close")
	pocs = filterPoc(PocInfo.PocName)
	CheckMultiPoc(req, pocs, PocInfo.Target, config.PocNum)
}

func CheckMultiPoc(req *xhttp.Request, pocs []*structs.Poc, target string, workers int) {
	tasks := make(chan Task)
	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		go func() {
			for task := range tasks {
				isVul, _ := ymlPoc.ExecutePoc(task.Req, target, task.Poc)
				if isVul {
					result := fmt.Sprintf("[!] %s %s", task.Req.GetUrl(), task.Poc.Name)
					pkgUtils.LogSuccess(result)
				}
				wg.Done()
			}
		}()
	}

	for _, poc := range pocs {
		task := Task{
			Req: req,
			Poc: poc,
		}
		wg.Add(1)
		tasks <- task
	}
	wg.Wait()
	close(tasks)
}

func initpoc() {
	if config.PocFile == "" {
		entries, err := FS.ReadDir("files")
		if err != nil {
			fmt.Printf("[-] init pocs error: %v\n", err)
			return
		}
		for _, one := range entries {
			path := one.Name()
			if strings.HasSuffix(path, ".yaml") || strings.HasSuffix(path, ".yml") {
				if poc, _ := utils.LoadPoc(path, FS); poc != nil {
					AllPocs = append(AllPocs, poc)
				}
			}
		}
	} else {
		err := filepath.Walk(config.PocFile,
			func(path string, info os.FileInfo, err error) error {
				if err != nil || info == nil {
					return err
				}
				if !info.IsDir() {
					if strings.HasSuffix(path, ".yaml") || strings.HasSuffix(path, ".yml") {
						poc, _ := utils.LoadPocbyPath(path)
						if poc != nil {
							AllPocs = append(AllPocs, poc)
						}
					}
				}
				return nil
			})
		if err != nil {
			fmt.Printf("[-] init pocs error: %v\n", err)
		}
	}
}

func filterPoc(pocName string) (pocs []*structs.Poc) {
	if pocName == "" {
		return AllPocs
	}
	for _, poc := range AllPocs {
		if strings.Contains(poc.Name, pocName) {
			pocs = append(pocs, poc)
		}
	}
	return
}

func CheckInfoPoc(infostr string) string {
	for _, goPoc := range gopoc.GoPocDatas {
		if infostr == goPoc.Name {
			return goPoc.Alias
		}
	}
	return ""
}
