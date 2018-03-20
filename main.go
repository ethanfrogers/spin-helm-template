package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"k8s.io/helm/pkg/chartutil"
	"k8s.io/helm/pkg/downloader"
	"k8s.io/helm/pkg/engine"
	"k8s.io/helm/pkg/getter"
	"k8s.io/helm/pkg/helm/environment"
	"k8s.io/helm/pkg/helm/helmpath"
	"k8s.io/helm/pkg/proto/hapi/chart"
	"k8s.io/helm/pkg/timeconv"
)

type ArtifactResponse struct {
	Artifacts []SpinnakerArtifact `json:"artifacts"`
}

type SpinnakerArtifact struct {
	Type      string `json:"type"`
	Reference string `json:"reference"`
	Name      string `json:"name"`
}

type TemplateRequest struct {
	Chart       string `json:"chart"`
	Version     string `json:"version"`
	ReleaseName string `json:"releaseName"`
	Namespace   string `json:"namespace"`
}

func tpl(chartName, version, releaseName, namespace string) (map[string]string, error) {
	// download the requested chart
	dl := downloader.ChartDownloader{
		Out:      os.Stdout,
		HelmHome: helmpath.Home(os.Getenv("HELM_HOME")),
		Getters:  getter.All(environment.EnvSettings{}),
	}

	tmpDir, err := ioutil.TempDir("", "")
	if err != nil {
		return nil, err
	}

	filename := ""
	fname, _, err := dl.DownloadTo(chartName, version, tmpDir)

	if err != nil {
		return nil, err
	}

	if err == nil {
		lname, err := filepath.Abs(fname)
		if err != nil {
			return nil, err
		}
		filename = lname
		// return nil, nil
	}
	// render the chart
	config := &chart.Config{Raw: string([]byte{}), Values: map[string]*chart.Value{}}

	options := chartutil.ReleaseOptions{
		Name:      releaseName,
		Time:      timeconv.Now(),
		Namespace: namespace,
	}
	c, err := chartutil.Load(filename)
	if err != nil {
		return nil, err
	}
	renderer := engine.New()
	vals, err := chartutil.ToRenderValues(c, config, options)
	if err != nil {
		return nil, err
	}
	return renderer.Render(c, vals)
	// respond with artifacts
}

// https://github.com/spinnaker/spinnaker/issues/2330
func main() {

	mux := http.NewServeMux()
	mux.HandleFunc("/template", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("not found"))
			return
		}

		cht := &TemplateRequest{}
		if err := json.NewDecoder(r.Body).Decode(cht); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		templated, err := tpl(cht.Chart, cht.Version, cht.ReleaseName, cht.Namespace)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		artifacts := []SpinnakerArtifact{}

		for key, v := range templated {
			e := base64.StdEncoding.EncodeToString([]byte(v))
			artifacts = append(artifacts, SpinnakerArtifact{
				Type:      "embedded/base64",
				Reference: e,
				Name:      key,
			})
		}

		wrapper := ArtifactResponse{Artifacts: artifacts}

		b, _ := json.Marshal(wrapper)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(b)
	})

	server := &http.Server{
		Addr:    ":3005",
		Handler: mux,
	}

	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
