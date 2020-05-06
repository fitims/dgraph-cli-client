# dgraph cli client

This is a dgraph cli client that will execute commands on dgraph. I needed an app to execute queries and mutations on a production server that cannot run ratel or anything else on it.

It uses ```DGRAPH_CONNECTION``` environment variable to read the connection details to the instance of dgraph. If it is not provided,. it will try and connect to localhost.

The cli has the following features:
```
dg exec -query '{ user(func: eq(type, "user")) { uid name } }'
dg exec -f -query '{ user(func: eq(type, "user")) { uid name } }' (beautify JSON output)
dg exec -schema "jobTitle: string @index(exact) . "
dg exec -mutate '<0xa4> <jobTitle> "job1" .'
dg exec -delete "<0xa4> * * ."
```
You can also use a yaml file with schema definitions or mutations:
```
dg yaml -file "/users/usera/changes.yaml"
```

the format of the yaml files is :
```
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
  ```
