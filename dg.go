package main

import (
	"dg/definitions"
	"dg/dgraph"
	"dg/env"
	"flag"
	"fmt"
	"os"
)

const(
	usage = `
Usage:

  dg exec -query '{ user(func: eq(type, "user")) { uid name } }'
  dg exec -f -query '{ user(func: eq(type, "user")) { uid name } }' (beautify JSON output)
  dg exec -schema "jobTitle: string @index(exact) . "
  dg exec -mutate '<0xa4> <jobTitle> "job1" .'
  dg exec -delete "<0xa4> * * ."
  
You can also use a yaml file with schema definitions or mutations:

  dg yaml -file "/users/usera/changes.yaml"
`

	ymlFormat = `
  The format of yaml file is :

	description : "schema changes for banner"
	
	schema :
	  - "jobTitle: string @index(exact) . "
	  - "jobParticipant: uid @reverse . "
	
	mutations :
	  - _:job <jobTitle> "job1" .
	  - _:job <jobParticipant> <0x11> .
	  - _:job <type> "job" .
	
	deletes :
	  - "<0x11> * * . "
	  - "<0x3f30> <mediaArtist> * . "
`
)


func main() {
	connStr := env.GetEnv("DGRAPH_CONNECTION", "localhost:9080")
	yamlCmd := flag.NewFlagSet("yaml", flag.ExitOnError)
	fileYaml := yamlCmd.String("file", "", "the YAML file containing mutations/schema")

	exeCmd := flag.NewFlagSet("exec", flag.ExitOnError)
	qryCmd := exeCmd.String("query",  "","the query to be executed")
	beautifyCmd := exeCmd.Bool("f",  false,"the query result (JSON) can be beautified")
	mutCmd := exeCmd.String("mutate", "", "the mutation to be executed")
	schCmd := exeCmd.String("schema", "", "the schema to be executed")
	delCmd := exeCmd.String("delete", "", "the delete to be executed")

	if len(os.Args) < 2 {
		fmt.Println("\n\n" + `ERROR: Parameters are missing !`)
		fmt.Println(usage)
		return
	}

	if os.Args[1] == "yaml" {
		yamlCmd.Parse(os.Args[2:])
	} else {
		flag.Parse()
	}

	if yamlCmd.Parsed() {
		if *fileYaml == "" {
			fmt.Println(`ERROR: You need to specify the YAML file (ie. -file="/home/usera/changes.yaml"\n\n`)
			fmt.Println(usage)
			fmt.Println(ymlFormat)
			return
		}
	}

	if os.Args[1] == "exec" {
		exeCmd.Parse(os.Args[2:])
	} else {
		flag.Parse()
	}

	if exeCmd.Parsed() {
		if *qryCmd == "" && *mutCmd == "" && *schCmd == "" && *delCmd == "" {
			fmt.Println(`ERROR: You need to specify the query/mutation/schema/delete to be executed\n\n`)
			fmt.Println(usage)
			return
		}
	}

	dgraph.Open(connStr)
	defer dgraph.Close()

	if yamlCmd.Parsed() {
		fmt.Println("\nsetting schema/mutations from file ...")
		definitions.SetSchemaAndDataFromFile(dgraph.Client, *fileYaml)
		return
	}

	if exeCmd.Parsed() {

		if len(*qryCmd) > 0 {
			fmt.Println("\nRunning query : ", *qryCmd)
			definitions.RunQuery(dgraph.Client, *qryCmd, *beautifyCmd)
			return
		}

		if len(*mutCmd) > 0 {
			fmt.Println("\nExecuting mutation : ", *mutCmd)
			definitions.SetMutation(dgraph.Client, *mutCmd)
			return
		}

		if len(*schCmd) > 0 {
			fmt.Println("\nExecuting schema : ", *schCmd)
			definitions.SetSchema(dgraph.Client, *schCmd)
			return
		}

		if len(*delCmd) > 0 {
			fmt.Println("\nExecuting delete : ", *delCmd)
			definitions.SetDelete(dgraph.Client, *delCmd)
			return
		}
	}
}
