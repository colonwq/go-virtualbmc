package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
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

func send_recieve_message(outmessage []byte) ([]string, int, error) {
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
	//fmt.Println("Reply: ", reply)

	err = client.Close()

	if err != nil {
		fmt.Println("Close error: ", err)
	}

	return reply, size, err
}

func main() {

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
		list(
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
		add(addCmd.Arg(0),
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
		show(showCmd.Arg(0),
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
		return_string := simple_command("start", os.Args[2])
		fmt.Println("Retrun string: " , return_string)
	case "stop":
		if len(os.Args) != 3 {
			fmt.Println("no host given")
			os.Exit(2)
		}
		return_string := simple_command("stop", os.Args[2])
		fmt.Println("Retrun string: " , return_string)
	case "delete":
		if len(os.Args) != 3 {
			fmt.Println("no host given")
			os.Exit(2)
		}
		return_string := simple_command("delete", os.Args[2])
		fmt.Println("Retrun string: " , return_string)
	default:
		fmt.Println("Nothing to do. Move along")
		os.Exit(1)
	}
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
			Port uint) {

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
	//fmt.Println("json cmd: ", string(add_cmd_json))

	inmessage, size, err := send_recieve_message( add_cmd_json )
	if err != nil {
		fmt.Println( "Send/recieve error: ", err )
	}

	fmt.Println("json reply(", size, "): " , inmessage, "." )
}

func show(host string,
			Columns string,
			Fit_width bool,
			Formatter string,
			Max_width int,
			Noindent bool,
			Print_empty bool,
			Quote_mode string,
			Sort_columns string) {
	show_cmd := &short_struct{
		Command:     "show",
		Domain_name: host,
	}
	//these are needed for only the printing aspects. 
	//Columns:     Columns,
	//Sort_columns:     Sort_columns,
	//Fit_width:   Fit_width,
	//Formatter:   Formatter,
	//Max_width:   Max_width,
	//Noindent:    Noindent,
	//Print_empty: Print_empty,
	//Quote_mode:  Quote_mode,

	show_cmd_json, json_err := json.Marshal(show_cmd)
	if json_err != nil {
		fmt.Println("json marshal error: ", json_err)
	}
	//fmt.Println("json cmd: ", string(show_cmd_json))

	inmessage, size, err := send_recieve_message( show_cmd_json )
	if err != nil {
		fmt.Println( "Send/recieve error: ", err )
	}

	fmt.Println("json reply(", size, "): " , inmessage )
}

func simple_command(command string, host string) ([]string) {
	start_cmd := &short_struct{
		Command: command,
	}
	start_cmd.Domain_names = append(start_cmd.Domain_names, host )
	//fmt.Println("simple command: " , start_cmd.Command )

	start_cmd_json, json_err := json.Marshal(start_cmd)
	if json_err != nil {
		fmt.Println("json marshal error: ", json_err)
	}
	//fmt.Println("json cmd: ", string(start_cmd_json))

	inmessage, _, err := send_recieve_message( start_cmd_json )
	if err != nil {
		fmt.Println( "Send/recieve error: ", err )
	}

	//fmt.Println("json reply(", size, "): " , inmessage )

	return inmessage ;
}

func list(
			Columns string,
			Fit_width bool,
			Formatter string,
			Max_width int,
			Noindent bool,
			Print_empty bool,
			Quote_mode string,
			Sort_columns string) {

	list_cmd := &short_struct{
		Command:     "list",
	}
	//these are only needed for the printing aspect
	//Columns:     []string{},
	//Sort_columns:     []string{},
	//Fit_width:   false,
	//Formatter:   "table",
	//Max_width:   0,
	//Noindent:    false,
	//Print_empty: false,
	//Quote_mode:  "nonnumeric",


	list_cmd_json, json_err := json.Marshal(list_cmd)
	if json_err != nil {
		fmt.Println("json marshal error: ", json_err)
	}
	//fmt.Println("json cmd: ", string(list_cmd_json))

	inmessage, size, err := send_recieve_message( list_cmd_json )
	if err != nil {
		fmt.Println( "Send/recieve error: ", err )
	}

	fmt.Println("json reply(", size, "): " , inmessage )

}
