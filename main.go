package main

import (
	"database/sql"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/rpc"
	"path/filepath"
	"time"

	_ "github.com/lib/pq"
	"github.com/powerman/rpc-codec/jsonrpc2"
	"gopkg.in/yaml.v2"
)

var DBConn *sql.DB

type configFile struct {
	ListenPort         string `yaml:"listen_port"`
	DBConnectionString string `yaml:"db_connection_string"`
}

type HttpConn struct {
	in  io.Reader
	out io.Writer
}

func (c *HttpConn) Read(p []byte) (n int, err error) {
	return c.in.Read(p)
}

func (c *HttpConn) Write(d []byte) (n int, err error) {
	return c.out.Write(d)
}

func (c *HttpConn) Close() error {
	return nil
}

//User - RPC struct
type User struct {
	UUID             int    `json:"uuid"`
	Login            string `json:"login"`
	RegistrationDate int    `json:"registration_date"`
}

//Args - RPC args
type Args struct {
	ID    int    `json:"id"`
	Login string `json:"login"`
}

//Create - new user
func (user *User) Create(args *Args, result *string) error {

	_, err := DBConn.Exec("INSERT INTO users (login,registration_date) VALUES ($1,$2)", args.Login, time.Now().Unix())
	if err != nil {
		log.Println(err)
		return err
	}

	*result = "ok"
	return nil
}

//Get - get user by id
func (user *User) Get(args *Args, result *User) error {

	row := DBConn.QueryRow("SELECT * FROM users WHERE uuid = $1", args.ID)

	err := row.Scan(&user.UUID, &user.Login, &user.RegistrationDate)
	if err != nil {
		log.Println(err)
		return err
	}

	*result = *user
	return nil
}

//Update - update user login
func (user *User) Update(args *Args, result *string) error {

	_, err := DBConn.Exec("UPDATE users SET login = $1 WHERE uuid = $2", args.Login, args.ID)
	if err != nil {
		log.Println(err)
		return err
	}

	*result = "ok"
	return nil
}

func GetConfig() (configFile, error) {
	var conf configFile

	filename, _ := filepath.Abs("./configs.yaml")
	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Println(err)
		return conf, err
	}

	err = yaml.Unmarshal(yamlFile, &conf)
	if err != nil {
		log.Println(err)
		return conf, err
	}

	return conf, nil
}

func main() {
	var err error

	//Get configs
	conf, err := GetConfig()
	if err != nil {
		panic(err)
	}

	//Get DB connection
	DBConn, err = sql.Open("postgres", conf.DBConnectionString)
	if err != nil {
		panic(err)
	}
	defer DBConn.Close()

	server := rpc.NewServer()
	server.Register(&User{})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		//jsonRPC server
		serverCodec := jsonrpc2.NewServerCodec(&HttpConn{in: r.Body, out: w}, server)

		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusOK)

		if err1 := server.ServeRequest(serverCodec); err1 != nil {
			http.Error(w, "Error while serving JSON request", http.StatusInternalServerError)
			return

		}
	})
	http.ListenAndServe(conf.ListenPort, nil)
}
