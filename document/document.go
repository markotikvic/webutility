package document

import (
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"time"

	web "git.to-net.rs/marko.tikvic/webutility"
)

// Document ...
type Document struct {
	ID               int64         `json:"id"`
	FileName         string        `json:"fileName"`
	Extension        string        `json:"extension"`
	ContentType      string        `json:"contentType"`
	Size             int64         `json:"fileSize"`
	UploadedBy       string        `json:"uploadedBy"`
	LastModifiedBy   string        `json:"lastModifiedBy"`
	TimeUploaded     int64         `json:"timeUploaded"`
	TimeLastModified int64         `json:"timeLastModified"`
	RoleAccessLevel  int64         `json:"accessLevel"`
	Description      string        `json:"description"`
	Download         *DownloadLink `json:"download"`
	Path             string        `json:"-"`
	directory        string
	data             []byte
}

// OpenFileAsDocument ...
func OpenFileAsDocument(path string) (*Document, error) {
	d := &Document{Path: path}

	f, err := os.Open(d.Path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	stats, err := f.Stat()
	if err != nil {
		return nil, err
	}

	d.FileName = stats.Name()
	d.Size = stats.Size()
	d.Extension = web.FileExtension(d.FileName)

	d.data = make([]byte, d.Size)
	if _, err = f.Read(d.data); err != nil {
		return nil, err
	}

	return d, err
}

// DownloadLink ...
type DownloadLink struct {
	Method string `json:"method"`
	URL    string `json:"url"`
}

// SetDownloadInfo ...
func (d *Document) SetDownloadInfo(method, url string) {
	d.Download = &DownloadLink{
		Method: method,
		URL:    url,
	}
}

// ServeDocument writes d's buffer to w and sets appropriate headers according to d's content type
// and downloadPrompt.
func ServeDocument(w http.ResponseWriter, d *Document, downloadPrompt bool) error {
	f, err := os.Open(d.Path)
	if err != nil {
		return err
	}
	defer f.Close()

	web.SetContentType(w, mime.TypeByExtension(d.Extension))
	web.SetResponseStatus(w, http.StatusOK)
	if downloadPrompt {
		w.Header().Set("Content-Disposition", "attachment; filename="+d.FileName)
	}

	buf := make([]byte, d.Size)
	if _, err := f.Read(buf); err != nil {
		return err
	}

	w.Header().Set("Content-Length", fmt.Sprintf("%d", d.Size))
	web.WriteResponse(w, buf)

	return nil
}

// ParseDocument ...
func ParseDocument(req *http.Request) (doc *Document, err error) {
	req.ParseMultipartForm(32 << 20)
	file, fheader, err := req.FormFile("document")
	if err != nil {
		return doc, err
	}

	claims, _ := web.GetTokenClaims(req)
	owner := claims.Username

	fname := fheader.Filename

	fsize := fheader.Size
	ftype := fmt.Sprintf("%v", fheader.Header["Content-Type"][0])

	fextn := web.FileExtension(fname)
	if fextn == "" {
		return doc, errors.New("invalid extension")
	}

	doc = new(Document)

	doc.FileName = fname
	doc.Size = fsize
	doc.ContentType = ftype
	doc.Extension = "." + fextn

	t := time.Now().Unix()
	doc.TimeUploaded = t
	doc.TimeLastModified = t

	doc.UploadedBy = owner
	doc.LastModifiedBy = owner
	doc.RoleAccessLevel = 0

	doc.data = make([]byte, doc.Size)
	if _, err = io.ReadFull(file, doc.data); err != nil {
		return doc, err
	}

	return doc, nil
}

// SaveToFile ...
func (d *Document) SaveToFile(path string) (f *os.File, err error) {
	d.Path = path

	if web.FileExists(path) {
		err = fmt.Errorf("file %s alredy exists", path)
		return nil, err
	}

	if parentDir := web.DirectoryFromPath(path); parentDir != "" {
		if err = os.MkdirAll(parentDir, os.ModePerm); err != nil {
			if !os.IsExist(err) {
				return nil, err
			}
		}
	}

	if f, err = os.Create(path); err != nil {
		return nil, err
	}

	if _, err = f.Write(d.data); err != nil {
		f.Close()
		d.DeleteFile()
		return nil, err
	}
	f.Close()

	return f, nil
}

func DeleteDocuments(docs []*Document) error {
	for _, d := range docs {
		if d == nil {
			continue
		}
		if err := d.DeleteFile(); err != nil {
			return err
		}
	}
	return nil
}

// DeleteFile ...
func (d *Document) DeleteFile() error {
	return os.Remove(d.Path)
}
