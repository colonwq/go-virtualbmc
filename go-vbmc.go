package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"text/tabwriter"
	"time"
	zmq "github.com/pebbe/zmq4"
)

type short_struct struct {
	Domain_name  string   `json:"domain_name"`
	Domain_names []string `json:"domain_names"`
	Command      string   `json:"command"`
}

type add_struct struct {
	Address               string `json:"address"`
	Namespace             string `json:"namespace"`
	Name                  string `json:"name"`
	Domain_name           string `json:"domain_name"`
	Libvirt_sasl_password string `json:"libvirt_sasl_password"`
	Libvirt_sasl_username string `json:"libvirt_sasl_username"`
	Libvirt_uri           string `json:"libvirt_uri"`
	Password              string `json:"password"`
	Port                  uint `json:"port"`
	Username              string `json:"username"`
	Command               string `json:"command"`
}

const (
	REQUEST_TIMEOUT = 1000 * time.Millisecond
	//MAX_RETRIES     = 3 //  Before we abandon
)

type Response_msg struct {
    Rc   int      `json:"rc"`
    Msg []string  `json:"msg"`
}

type Response struct {
    Rc   int      `json:"rc"`
    Header []string `json:"header"`
    Rows [][]string `json:"rows"`
}

func send_recieve_message_string(outmessage []byte) ([]string, int, error) {
	var size int
	//creat a request client
	client, _ := zmq.NewSocket(zmq.REQ)

	err := client.Connect("tcp://127.0.0.1:50891")

	if err != nil {
		fmt.Println("Connection error: ", err)
	}

	size, err = client.SendMessage( outmessage )

	if err != nil {
		fmt.Println("Send error: ", err)
	}

	poller := zmq.NewPoller()
	poller.Add(client, zmq.POLLIN)
	polled, err := poller.Poll( REQUEST_TIMEOUT )
	reply := []string{}
	if len(polled) == 1 {
		reply, err = client.RecvMessage(0)
	}

	err = client.Close()

	if err != nil {
		fmt.Println("Close error: ", err)
	}

	return reply, size, err
}

func send_recieve_message(outmessage []byte) ([]byte, int, error) {
	var size int
	//creat a request client
	client, _ := zmq.NewSocket(zmq.REQ)

	err := client.Connect("tcp://192.168.39.17:50891")
	//err := client.Connect("tcp://127.0.0.1:50891")

	if err != nil {
		fmt.Println("Connection error: ", err)
	}

	size, err = client.SendMessage( outmessage )

	if err != nil {
		fmt.Println("Send error: ", err)
	}

	poller := zmq.NewPoller()
	poller.Add(client, zmq.POLLIN)
	polled, err := poller.Poll( REQUEST_TIMEOUT )
	reply := []byte{}
	if len(polled) == 1 {
		reply, err = client.RecvBytes(0)
	}

	err = client.Close()

	if err != nil {
		fmt.Println("Close error: ", err)
	}

	return reply, size, err
}

func main() {

	var rc int = 0
	var msg string
	if len(os.Args) < 2 {
		fmt.Println( "lots of nothing to do" )
		os.Exit(1)
	}

	addCmd                   := flag.NewFlagSet("add", flag.ExitOnError)
	AddAddress               := addCmd.String("address","::","Bind address" )
	AddNamespace             := addCmd.String("namespace","default","k8 Namespace")
	AddName                  := addCmd.String("name","default","k8 VM/VMI Name")
	AddLibvirt_sasl_password := addCmd.String("libvirt_sasl_password","","Libvirt SASL Password" )
	AddLibvirt_sasl_username := addCmd.String("libvirt_sasl_username","","Libvirt SASL Username" )
	AddLibvirt_uri           := addCmd.String("libvirt_uri", "qemu:///system","Libvirt URI connection string" )
	AddPassword              := addCmd.String("password","password","IPMI login password" )
	AddPort                  := addCmd.Uint("port",623,"IPMI conection port" )
	AddUsername              := addCmd.String("username","admin","IPMI login username" )

	showCmd := flag.NewFlagSet("show", flag.ExitOnError)
	ShowColumns := showCmd.String("columns","","specify the column(s) to include, can be repeated")
	ShowFit_width := showCmd.Bool("fit_width",false,"Fit the table to the display width. Implied if --max-\nwidth greater than 0. Set the environment variable\nCLIFF_FIT_WIDTH=1 to always enable")
	ShowFormatter := showCmd.String("formatter","table","{csv,json,table,value,yaml}, --format {csv,json,table,value,yaml}\nthe output format, defaults to table")
	ShowMax_width := showCmd.Int("max_width",-1,"Maximum display width, <1 to disable. You can also use\nthe CLIFF_MAX_TERM_WIDTH environment variable, but the\nparameter takes precedence.")
	ShowNoindent := showCmd.Bool("noindent",false,"whether to disable indenting the JSON")
	ShowPrint_empty := showCmd.Bool("print_empty",false,"Print empty table if there is no data to show.")
	ShowQuote_mode := showCmd.String("quote_mode","nonnumeric","when to include quotes, defaults to nonnumeric")
	ShowSort_columns := showCmd.String("sort_columns","","specify the column(s) to sort the data (columns\nspecified first have a priority, non-existing columns\nare ignored), can be repeated")


	switch os.Args[1] {
	case "list":
		err := showCmd.Parse(os.Args[2:])
		if err != nil {
			fmt.Println("Error parsing the command line: ", err )
			os.Exit(2)
		}
		rc, msg = list(
			*ShowColumns,
			*ShowFit_width,
			*ShowFormatter,
			*ShowMax_width,
			*ShowNoindent,
			*ShowPrint_empty,
			*ShowQuote_mode,
			*ShowSort_columns)
	case "add":
		err := addCmd.Parse(os.Args[2:])
		if err != nil {
			fmt.Println("Error parsing the command line: ", err )
			os.Exit(2)
		}
		rc, msg = add(addCmd.Arg(0),
			*AddAddress,
			*AddNamespace,
			*AddName,
			*AddLibvirt_sasl_username,
			*AddLibvirt_sasl_password,
			*AddLibvirt_uri,
			*AddUsername,
			*AddPassword,
			*AddPort)
	case "show":
		if len(os.Args) < 3 {
			fmt.Println("no host given")
			os.Exit(2)
		}
		err := showCmd.Parse(os.Args[2:])
		if err != nil {
			fmt.Println("Error parsing the command line: ", err )
			os.Exit(2)
		}
		rc, msg = show(showCmd.Arg(0),
			*ShowColumns,
			*ShowFit_width,
			*ShowFormatter,
			*ShowMax_width,
			*ShowNoindent,
			*ShowPrint_empty,
			*ShowQuote_mode,
			*ShowSort_columns)
	case "start":
		if len(os.Args) != 3 {
			fmt.Println("no host given")
			os.Exit(2)
		}
		rc, msg = simple_command("start", os.Args[2])
	case "stop":
		if len(os.Args) != 3 {
			fmt.Println("no host given")
			os.Exit(2)
		}
		rc, msg = simple_command("stop", os.Args[2])
	case "delete":
		if len(os.Args) != 3 {
			fmt.Println("no host given")
			os.Exit(2)
		}
		rc, msg = simple_command("delete", os.Args[2])
	default:
		msg = "Nothing to do. Move along"
		rc = 1
	}

	if rc != 0 {
		fmt.Println("Error: ", msg )
	}
	os.Exit(rc)
}

func add(host string,
			Address string,
			Namespace string,
			Name string,
			Libvirt_sasl_username string,
			Libvirt_sasl_password string,
			Libvirt_uri string,
			Username string,
			Password string,
			Port uint) (int, string){

	add_cmd := &add_struct{
		Command: "add",
		Domain_name: host,
		Port: Port,
		Password: Password,
		Username: Username,
		Address: Address,
		Libvirt_sasl_username: Libvirt_sasl_username,
		Libvirt_sasl_password: Libvirt_sasl_password,
		Libvirt_uri: Libvirt_uri,
		Namespace: Namespace,
		Name: Name,
	}

	add_cmd_json, json_err := json.Marshal(add_cmd)

	if json_err != nil {
		fmt.Println("json marshal error: ", json_err)
	}

	inmessage, _, err := send_recieve_message( add_cmd_json )
	if err != nil {
		fmt.Println( "Send/recieve error: ", err )
	}

	res := Response_msg{}
	bytes := []byte(inmessage)
	json.Unmarshal(bytes , &res )

	if res.Rc == 0 {
		return res.Rc, "Ok"
	} else {
		return res.Rc, res.Msg[0]
	}
}

func show(host string,
			Columns string,
			Fit_width bool,
			Formatter string,
			Max_width int,
			Noindent bool,
			Print_empty bool,
			Quote_mode string,
			Sort_columns string) (int, string) {

	var msg [1]string
	show_cmd := &short_struct{
		Command:     "show",
		Domain_name: host,
	}

	show_cmd_json, json_err := json.Marshal(show_cmd)
	if json_err != nil {
		fmt.Println("json marshal error: ", json_err)
	}

	inmessage, _, err := send_recieve_message( show_cmd_json )
	if err != nil {
		fmt.Println( "Send/recieve error: ", err )
	}

	res := Response{}
	bytes := []byte(inmessage)
	json.Unmarshal(bytes , &res )

	if res.Rc == 0 {
		w := new(tabwriter.Writer)
		w.Init(os.Stdout, 8, 8, 2, ' ', tabwriter.Debug)

		fmt.Fprintf( w, "%s\t%s\n", res.Header[0], res.Header[1] )

		for i := 0; i < len(res.Rows); i++ {
			fmt.Fprintf( w, "%s\t%s\n", res.Rows[i][0], res.Rows[i][1] )
		}

		w.Flush()
		msg[0] = "Ok"

	} else {
		res := Response_msg{}
		bytes := []byte(inmessage)
		json.Unmarshal(bytes , &res )
		msg[0] = res.Msg[0]
	}

	return res.Rc, msg[0]
}

func simple_command(command string, host string) (int, string) {
	start_cmd := &short_struct{
		Command: command,
	}
	start_cmd.Domain_names = append(start_cmd.Domain_names, host )

	start_cmd_json, json_err := json.Marshal(start_cmd)
	if json_err != nil {
		fmt.Println("json marshal error: ", json_err)
	}

	inmessage, _, err := send_recieve_message( start_cmd_json )
	if err != nil {
		fmt.Println( "Send/recieve error: ", err )
	}

	res := Response_msg{}
	bytes := []byte(inmessage)
	json.Unmarshal(bytes , &res )

	if res.Rc == 0 {
		return res.Rc, "Ok"
	} else {
		return res.Rc, res.Msg[0]
	}
}

func list(
			Columns string,
			Fit_width bool,
			Formatter string,
			Max_width int,
			Noindent bool,
			Print_empty bool,
			Quote_mode string,
			Sort_columns string) (int, string){
	var msg [1]string

	list_cmd := &short_struct{
		Command:     "list",
	}

	list_cmd_json, json_err := json.Marshal(list_cmd)
	if json_err != nil {
		fmt.Println("json marshal error: ", json_err)
	}

	inmessage, _, err := send_recieve_message( list_cmd_json )
	if err != nil {
		fmt.Println( "Send/recieve error: ", err )
	}

	res := Response{}
	bytes := []byte(inmessage)
	json.Unmarshal(bytes , &res )

	fmt.Println( res )
	if res.Rc == 0 {
		if len(res.Rows) > 0 {
			w := new(tabwriter.Writer)
			w.Init(os.Stdout, 8, 8, 2, ' ', tabwriter.Debug)

			fmt.Fprintf( w, "%s\t%s\t%s\n", res.Header[0], res.Header[1], res.Header[2] )

			for i := 0; i < len(res.Rows); i++ {
				fmt.Fprintf( w, "%s\t%s\t%s\n", res.Rows[i][0], res.Rows[i][1], res.Rows[i][2] )
			}

			w.Flush()
		}

		return res.Rc, "Ok"
	} else {
		res := Response_msg{}
		bytes := []byte(inmessage)
		json.Unmarshal(bytes , &res )
		msg[0] = res.Msg[0]
		return res.Rc, res.Msg[0]
	}
}
