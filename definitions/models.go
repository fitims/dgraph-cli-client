package definitions

import (
	"context"
	"fmt"
	"github.com/dgraph-io/dgo"
	"github.com/dgraph-io/dgo/protos/api"
	"github.com/tidwall/pretty"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"strings"
)

type FileSettings struct {
	Description string   `yaml:"description"`
	Schema      []string `yaml:"schema"`
	Mutations   []string `yaml:"mutations"`
	Deletes     []string `yaml:"deletes"`
	query       string   `yaml:"query"`
}

func SetSchemaAndDataFromFile(cl *dgo.Dgraph, yamlFile string) error {

	data, err := ioutil.ReadFile(yamlFile)
	if err != nil {
		fmt.Println("Error : ", err)
		return err
	}

	var t FileSettings
	err = yaml.Unmarshal([]byte(data), &t)
	if err != nil {
		fmt.Println("Error : ", err)
		return err
	}

	fmt.Println("Setting YAML : ", t.Description)
	if len(t.Schema) > 0 {
		for _, s := range t.Schema {
			err := SetSchema(cl, s)
			if err != nil {
				fmt.Println("Error setting schema from YAML file : ", err)
				return err
			}
		}
	}

	if len(t.Mutations) > 0 {
		mutation := strings.Join(t.Mutations, "\n ")
		fmt.Println("setting mutations : \n", mutation)
		err := SetMutation(cl, mutation)
		if err != nil {
			fmt.Println("Error setting mutations from YAML file : ", err)
			return err
		}

	}

	if (len(t.Deletes)) > 0 {
		for _, s := range t.Deletes {
			err := SetDelete(cl, s)
			if err != nil {
				fmt.Println("Error deleting edges from YAML file : ", err)
				return err
			}
		}
	}

	if len(t.query) > 0 {
		RunQuery(cl, t.query, false)
	}

	return nil
}

func RunQuery(cl *dgo.Dgraph, qry string, beutify bool) error {

	res, err := cl.NewTxn().Query(context.Background(), qry)
	if err != nil {
		fmt.Println("Error in getting response from server : ", err)
		return err
	}

	fmt.Println("response schema :", res.Schema)

	fmt.Println("result : \n-----------------------------\n")

	if beutify {
		fmt.Println(string(pretty.Color(pretty.Pretty(res.Json), nil)))
	} else {
		fmt.Println(string(res.Json), nil)
	}

	fmt.Println("\n-----------------------------\n")
	return nil
}

func SetSchema(cl *dgo.Dgraph, sch string) error {
	req := api.Operation{
		Schema: sch,
	}

	fmt.Println("executing : ", sch)
	err := cl.Alter(context.Background(), &req)
	if err != nil {
		fmt.Println("Error in getting response from server : ", err)
		return err
	}
	return nil
}

func SetMutation(cl *dgo.Dgraph, m string) error {
	req := &api.Mutation{
		CommitNow: true,
		SetNquads: []byte(m),
	}

	fmt.Println("populating/changing : \n", m)
	v, err := cl.NewTxn().Mutate(context.Background(), req)
	if err != nil {
		fmt.Println("Error in getting response from server : ", err)
		return err
	}
	fmt.Println("Assigned uids   : ", v.Uids)
	return nil
}

func SetDelete(cl *dgo.Dgraph, m string) error {
	req := &api.Mutation{
		CommitNow: true,
		DelNquads: []byte(m),
	}

	_, err := cl.NewTxn().Mutate(context.Background(), req)
	if err != nil {
		fmt.Println("Error in getting response from server : ", err)
		return err
	}

	return nil
}