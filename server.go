package webutility

import (
	"database/sql"
	"fmt"
	"net/http"

	"git.to-net.rs/marko.tikvic/webutility/logger"
	"github.com/gorilla/mux"
)

type Server struct {
	DB     *sql.DB
	Router *mux.Router
	Logger *logger.Logger
	Port   string
	DBs    map[string]*sql.DB
	dsn    map[string]string
}

func NewODBCServer(dsn, port, logDir string) (s *Server, err error) {
	s = new(Server)

	s.Port = port

	if s.DB, err = sql.Open("odbc", fmt.Sprintf("DSN=%s;", dsn)); err != nil {
		return nil, err
	}

	s.Router = mux.NewRouter()

	if s.Logger, err = logger.New("err", logDir, logger.MaxLogSize1MB); err != nil {
		return nil, fmt.Errorf("can't create logger: %s", err.Error())
	}

	s.DBs = make(map[string]*sql.DB)
	s.DBs["default"] = s.DB

	s.dsn = make(map[string]string)
	s.dsn["default"] = dsn

	return s, nil
}

func (s *Server) Run() {
	s.Logger.Print("Server listening on %s", s.Port)
	s.Logger.PrintAndTrace(http.ListenAndServe(s.Port, s.Router).Error())
}

func (s *Server) Cleanup() {
	if s.DB != nil {
		s.DB.Close()
	}

	if s.Logger != nil {
		s.Logger.Close()
	}
}

func (s *Server) StartTransaction() (*sql.Tx, error) {
	return s.DB.Begin()
}

func CommitChanges(tx *sql.Tx, err *error, opt ...error) {
	if *err != nil {
		tx.Rollback()
		return
	}

	for _, e := range opt {
		if e != nil {
			tx.Rollback()
			return
		}
	}

	if *err = tx.Commit(); *err != nil {
		tx.Rollback()
	}
}
