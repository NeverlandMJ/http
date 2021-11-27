package app

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/NeverlandMJ/http/pkg/banners"
)

// this is out logcal server
type Server struct {
	mux        *http.ServeMux
	bannersSvc *banners.Service
}

// it creates server
func NewServer(mux *http.ServeMux, bannersSvc *banners.Service) *Server {
	return &Server{mux: mux, bannersSvc: bannersSvc}
}

func (s *Server) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	s.mux.ServeHTTP(writer, request)
}

// serverni ishga solish
func (s *Server) Init() {
	s.mux.HandleFunc("/banners.getAll", s.handleGetAllBanners)
	s.mux.HandleFunc("/banners.getById", s.handleGetBannerByID)
	s.mux.HandleFunc("/banners.save", s.handleSaveBanner)
	s.mux.HandleFunc("/banners.removeById", s.handleremoveByID)
}

// get All
func (s *Server) handleGetAllBanners(w http.ResponseWriter, r *http.Request) {
	banners, err := s.bannersSvc.All(r.Context())
	if err != nil {
		log.Print(err)
		http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
		return
	}

	data, err := json.Marshal(banners)
	if err != nil {
		log.Print(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	if err != nil {
		log.Print(err)
	}
}

// get by id
func (s *Server) handleGetBannerByID(w http.ResponseWriter, r *http.Request) {
	idParam := r.URL.Query().Get("id")

	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		log.Print(err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	banner, err := s.bannersSvc.ByID(r.Context(), id)
	if err != nil {
		log.Print(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(banner)
	if err != nil {
		log.Print(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Contetn-Type", "applicatrion/json")
	_, err = w.Write(data)
	if err != nil {
		log.Print(err)
	}
}

// Save and update
// func (s *Server) handleSaveBanner(w http.ResponseWriter, r *http.Request) {
// 	idParam := r.URL.Query().Get("id")
// 	titleParam := r.URL.Query().Get("title")
// 	contentParam := r.URL.Query().Get("content")
// 	buttonParam := r.URL.Query().Get("button")
// 	linkParam := r.URL.Query().Get("link")

// 	id, err := strconv.ParseInt(idParam, 10, 64)
// 	if err != nil {
// 		log.Print(err)
// 		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
// 		return
// 	}

// 	banner := banners.Banner{
// 		ID:      id,
// 		Title:   titleParam,
// 		Content: contentParam,
// 		Button:  buttonParam,
// 		Link:    linkParam,
// 	}

// 	GotB, err := s.bannersSvc.Save(r.Context(), &banner)
// 	if err != nil {
// 		log.Print(err)
// 		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
// 		return
// 	}

// 	data, err := json.Marshal(GotB)
// 	if err != nil {
// 		log.Print(err)
// 		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
// 		return
// 	}
// 	w.Header().Set("Contetn-Type", "applicatrion/json")
// 	_, err = w.Write(data)
// 	if err != nil {
// 		log.Print(err)
// 	}

// }

// delete banner byID
func (s *Server) handleremoveByID(writer http.ResponseWriter, request *http.Request) {
	idParam := request.URL.Query().Get("id")

	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	dBanner, err := s.bannersSvc.RemoveByID(request.Context(), id)
	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(dBanner)
	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	writer.Header().Set("Contetn-Type", "applicatrion/json")
	_, err = writer.Write(data)
	if err != nil {
		log.Print(err)
	}
}

func (s *Server) handleUpload(w http.ResponseWriter, r *http.Request) {
	file, header, err := r.FormFile("image")
	if err != nil {
		log.Print(err)
		return
	}
	defer file.Close()
	newFile, err := os.Create("./web/banners/" + header.Filename)
	if err != nil {
		log.Print(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	defer newFile.Close()
}

func (s *Server) handleSaveBanner(w http.ResponseWriter, r *http.Request) {
	idp := r.FormValue("id")
	titlep := r.FormValue("title")
	contentp := r.FormValue("content")
	buttonp := r.FormValue("button")
	linkp := r.FormValue("link")

	id, err := strconv.ParseInt(idp, 10, 64)
	if err != nil {
		log.Print(err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	banner := &banners.Banner{
		ID:      id,
		Title:   titlep,
		Content: contentp,
		Button:  buttonp,
		Link:    linkp,
	}

	file, header, err := r.FormFile("image")
	if err == nil {
		im := strings.Split(header.Filename, ".")
		banner.Image = string(banner.ID) + "." + im[1]
	}

	Newbanner, err := s.bannersSvc.Save(r.Context(), banner, file)

	if err != nil {
		log.Print(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(Newbanner)
	if err != nil {
		log.Print(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Contetn-Type", "applicatrion/json")
	_, err = w.Write(data)
	if err != nil {
		log.Print(err)
	}

}
