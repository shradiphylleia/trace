package httpapi

import (
	"encoding/json"
	"errors"
	"html/template"
	"io"
	"net/http"
	"strconv"
	"strings"
	"github.com/go-chi/chi/v5"
	"traceshare/internal/app"
	"traceshare/internal/domain"
)

const maxMultipartMemory=32<<20

type Handler struct {
	service *app.Service
}

func NewHandler(service *app.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Register(router chi.Router) {
	router.Get("/", h.home)
	router.Get("/healthz", h.health)
	router.Post("/api/artifacts", h.createArtifact)
	router.Get("/api/artifacts/{shortCode}", h.getArtifact)
	router.Get("/api/artifacts/{shortCode}/download", h.downloadArtifact)
	router.Get("/api/search", h.searchArtifacts)
	router.Get("/t/{shortCode}", h.sharePage)
}

func (h *Handler) home(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_ =homeTemplate.Execute(w, nil)
}

func (h *Handler) health(w http.ResponseWriter, _*http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *Handler) createArtifact(w http.ResponseWriter,r *http.Request) {
	r.Body=http.MaxBytesReader(w,r.Body,26<<20)
	if err:=r.ParseMultipartForm(maxMultipartMemory);err!=nil {
		writeError(w, http.StatusBadRequest, "invalid multipart form")
		return
	}

	file,header,err:=r.FormFile("file")
	if err!=nil {
		writeError(w, http.StatusBadRequest, "file is required")
		return
	}
	defer file.Close()

	contentType:=header.Header.Get("Content-Type")
	if contentType==""{
		contentType="application/octet-stream"
	}

	input:=domain.CreateArtifactInput{
		Title:r.FormValue("title"),
		Description: r.FormValue("description"),
		Type:domain.ArtifactType(r.FormValue("artifact_type")),
		ServiceName: r.FormValue("service_name"),
		Environment: r.FormValue("environment"),
		Tags:splitTags(r.FormValue("tags")),
		Creator:r.FormValue("creator"),
		Expiration:r.FormValue("expiration"),
		FileName:header.Filename,
		ContentType:contentType,
		SizeBytes:header.Size,
	}

	artifact,err:=h.service.Create(r.Context(), input, file)
	if err!=nil{
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, map[string]any{
		"id":artifact.ID,
		"short_code":artifact.ShortCode,
		"short_url":artifact.ShareURL,
		"artifact":artifact,
	})
}

func (h *Handler) getArtifact(w http.ResponseWriter, r *http.Request) {
	artifact, err := h.service.Get(r.Context(), chi.URLParam(r, "shortCode"))
	if err!=nil{
		handleServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, artifact)
}

func (h *Handler) downloadArtifact(w http.ResponseWriter,r *http.Request) {
	artifact,body,err:=h.service.Download(r.Context(),chi.URLParam(r, "shortCode"))
	if err!=nil{
		handleServiceError(w, err)
		return
	}
	defer body.Close()

	w.Header().Set("Content-Type",artifact.ContentType)
	w.Header().Set("Content-Disposition", `attachment; filename="`+strings.ReplaceAll(artifact.FileName,`"`,"")+`"`)
	w.Header().Set("Content-Length", strconv.FormatInt(artifact.SizeBytes, 10))
	w.WriteHeader(http.StatusOK)
	_, _ = io.Copy(w, body)
}

func (h *Handler) searchArtifacts(w http.ResponseWriter, r *http.Request) {
	limit, _:=strconv.Atoi(r.URL.Query().Get("limit"))
	artifacts,err:=h.service.Search(r.Context(), app.SearchFilter{
		Query:r.URL.Query().Get("q"),
		Service: r.URL.Query().Get("service"),
		Tag:r.URL.Query().Get("tag"),
		Limit:limit,
	})
	if err!=nil {
		writeError(w, http.StatusInternalServerError, "search failed")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": artifacts})
}

func (h *Handler) sharePage(w http.ResponseWriter,r *http.Request) {
	artifact,err:=h.service.Get(r.Context(), chi.URLParam(r,"shortCode"))
	if err!=nil{
		handleServiceError(w, err)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_ = shareTemplate.Execute(w, artifact)
}

func splitTags(raw string) []string {
	parts:=strings.Split(raw,",")
	tags:=make([]string, 0, len(parts))
	for _,part:=range parts {
		part=strings.TrimSpace(part)
		if part!=""{
			tags=append(tags, part)
		}
	}
	return tags
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_=json.NewEncoder(w).Encode(value)
}

func writeError(w http.ResponseWriter,status int,message string) {
	writeJSON(w, status, map[string]string{"error":message})
}

func handleServiceError(w http.ResponseWriter, err error) {
	if errors.Is(err, app.ErrNotFound) {
		writeError(w, http.StatusNotFound, "artifact not found")
		return
	}
	writeError(w, http.StatusInternalServerError, "request failed")
}

var homeTemplate = template.Must(template.New("home").Parse(`<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>TraceShare</title>
  <style>
    :root { color-scheme: light; font-family: Inter, ui-sans-serif, system-ui, -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif; color: #17202a; background: #f6f7f9; }
    body { margin: 0; }
    main { max-width: 1120px; margin: 0 auto; padding: 32px 20px 48px; }
    header { display: flex; align-items: flex-end; justify-content: space-between; gap: 24px; margin-bottom: 24px; }
    h1 { font-size: 34px; margin: 0 0 6px; letter-spacing: 0; }
    p { margin: 0; color: #526070; line-height: 1.5; }
    .layout { display: grid; grid-template-columns: minmax(0, 1.1fr) minmax(320px, .9fr); gap: 20px; align-items: start; }
    section, form { background: #fff; border: 1px solid #d9dee7; border-radius: 8px; padding: 20px; box-shadow: 0 1px 2px rgba(21, 31, 43, .04); }
    label { display: grid; gap: 6px; font-size: 13px; font-weight: 650; color: #314154; }
    input, textarea, select { width: 100%; box-sizing: border-box; border: 1px solid #c8d0db; border-radius: 6px; padding: 10px 11px; font: inherit; color: #17202a; background: #fff; }
    textarea { min-height: 84px; resize: vertical; }
    .grid { display: grid; grid-template-columns: repeat(2, minmax(0, 1fr)); gap: 14px; }
    .stack { display: grid; gap: 14px; }
    button { border: 0; border-radius: 6px; background: #1f6feb; color: #fff; font-weight: 700; padding: 11px 14px; cursor: pointer; }
    button.secondary { background: #2f3b4a; }
    code, pre { font-family: "SFMono-Regular", Consolas, monospace; }
    pre { white-space: pre-wrap; background: #101820; color: #edf2f7; border-radius: 6px; padding: 14px; overflow: auto; }
    .result { display: grid; gap: 12px; }
    .pill { display: inline-flex; align-items: center; min-height: 24px; padding: 0 8px; border-radius: 999px; background: #eef3f8; color: #314154; font-size: 12px; font-weight: 700; }
    .items { display: grid; gap: 10px; }
    .item { border: 1px solid #d9dee7; border-radius: 6px; padding: 12px; }
    @media (max-width: 780px) { header, .layout, .grid { display: grid; grid-template-columns: 1fr; } }
  </style>
</head>
<body>
  <main>
    <header>
      <div>
        <h1>TraceShare</h1>
        <p>Upload debugging evidence once. Share one short link in Slack.</p>
      </div>
    </header>
    <div class="layout">
      <form id="upload" class="stack">
        <div class="grid">
          <label>Title <input name="title" required placeholder="Checkout API 500"></label>
          <label>Service <input name="service_name" required placeholder="payments"></label>
          <label>Type
            <select name="artifact_type" required>
              <option value="stack_trace">Stack trace</option>
              <option value="log">Log</option>
              <option value="api_payload">API payload</option>
              <option value="validation_report">Validation report</option>
              <option value="screenshot">Screenshot</option>
            </select>
          </label>
          <label>Environment <input name="environment" required placeholder="staging"></label>
          <label>Creator <input name="creator" required placeholder="shraddha@sharma.com"></label>
          <label>Expiration
            <select name="expiration">
              <option value="7d">7 days</option>
              <option value="14d">14 days</option>
              <option value="never">Never</option>
            </select>
          </label>
        </div>
        <label>Description <textarea name="description" placeholder="What changed, who observed it, and what flow failed"></textarea></label>
        <label>Tags <input name="tags" placeholder="checkout, sev2, qa"></label>
        <label>Artifact file <input name="file" type="file" required></label>
        <button type="submit">Upload and create short link</button>
        <div id="upload-result" class="result"></div>
      </form>
      <section class="stack">
        <form id="search" class="stack">
          <label>Search <input name="q" placeholder="error text, title, or tag"></label>
          <div class="grid">
            <label>Service <input name="service" placeholder="payments"></label>
            <label>Tag <input name="tag" placeholder="sev2"></label>
          </div>
          <button class="secondary" type="submit">Search artifacts</button>
        </form>
        <div id="search-result" class="items"></div>
      </section>
    </div>
  </main>
  <script>
    const upload = document.querySelector('#upload');
    const uploadResult = document.querySelector('#upload-result');
    upload.addEventListener('submit', async (event) => {
      event.preventDefault();
      uploadResult.textContent = 'Uploading...';
      const response = await fetch('/api/artifacts', { method: 'POST', body: new FormData(upload) });
      const body = await response.json();
      if (!response.ok) {
        uploadResult.innerHTML = '<pre>' + body.error + '</pre>';
        return;
      }
      uploadResult.innerHTML = '<p><strong>Short URL</strong></p><p><a href="' + body.short_url + '">' + body.short_url + '</a></p><pre>' + JSON.stringify(body.artifact, null, 2) + '</pre>';
      upload.reset();
    });

    const search = document.querySelector('#search');
    const searchResult = document.querySelector('#search-result');
    search.addEventListener('submit', async (event) => {
      event.preventDefault();
      const params = new URLSearchParams(new FormData(search));
      const response = await fetch('/api/search?' + params.toString());
      const body = await response.json();
      searchResult.innerHTML = body.items.map((item) => '<article class="item"><p><strong><a href="' + item.share_url + '">' + item.title + '</a></strong></p><p>' + item.service_name + ' / ' + item.environment + ' / ' + item.artifact_type + '</p><p>' + item.tags.map((tag) => '<span class="pill">' + tag + '</span>').join(' ') + '</p></article>').join('') || '<p>No matching artifacts.</p>';
    });
  </script>
</body>
</html>`))

var shareTemplate = template.Must(template.New("share").Parse(`<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>{{.Title}} - TraceShare</title>
  <style>
    :root { font-family: Inter, ui-sans-serif, system-ui, -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif; color: #17202a; background: #f6f7f9; }
    body { margin: 0; }
    main { max-width: 1080px; margin: 0 auto; padding: 32px 20px 48px; }
    header { margin-bottom: 20px; }
    h1 { font-size: 30px; margin: 0 0 8px; letter-spacing: 0; }
    p { margin: 0 0 10px; color: #526070; line-height: 1.5; }
    .panel { background: #fff; border: 1px solid #d9dee7; border-radius: 8px; padding: 20px; box-shadow: 0 1px 2px rgba(21, 31, 43, .04); }
    dl { display: grid; grid-template-columns: 150px 1fr; gap: 10px 16px; margin: 0 0 20px; }
    dt { font-weight: 750; color: #314154; }
    dd { margin: 0; color: #17202a; overflow-wrap: anywhere; }
    a.button { display: inline-flex; align-items: center; min-height: 40px; padding: 0 14px; border-radius: 6px; background: #1f6feb; color: #fff; font-weight: 700; text-decoration: none; }
    .pill { display: inline-flex; align-items: center; min-height: 24px; padding: 0 8px; border-radius: 999px; background: #eef3f8; color: #314154; font-size: 12px; font-weight: 700; margin-right: 6px; }
    pre { white-space: pre-wrap; background: #101820; color: #edf2f7; border-radius: 6px; padding: 16px; overflow: auto; }
    img { max-width: 100%; border: 1px solid #d9dee7; border-radius: 6px; }
    @media (max-width: 680px) { dl { grid-template-columns: 1fr; } }
  </style>
</head>
<body>
  <main>
    <header>
      <p>TraceShare artifact</p>
      <h1>{{.Title}}</h1>
      <p>{{.Description}}</p>
    </header>
    <section class="panel">
      <dl>
        <dt>Short code</dt><dd>{{.ShortCode}}</dd>
        <dt>Type</dt><dd>{{.Type}}</dd>
        <dt>Service</dt><dd>{{.ServiceName}}</dd>
        <dt>Environment</dt><dd>{{.Environment}}</dd>
        <dt>Creator</dt><dd>{{.Creator}}</dd>
        <dt>Created</dt><dd>{{.CreatedAt.Format "2006-01-02 15:04 UTC"}}</dd>
        <dt>Expires</dt><dd>{{if .ExpiresAt}}{{.ExpiresAt.Format "2006-01-02 15:04 UTC"}}{{else}}Never{{end}}</dd>
        <dt>Tags</dt><dd>{{range .Tags}}<span class="pill">{{.}}</span>{{else}}None{{end}}</dd>
      </dl>
      <p><a class="button" href="{{.DownloadURL}}">Download artifact</a></p>
      {{if .Preview}}<pre>{{.Preview}}</pre>{{else if eq .Type "screenshot"}}<img src="{{.DownloadURL}}" alt="{{.Title}}">{{end}}
    </section>
  </main>
</body>
</html>`))
