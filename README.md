# Temporal Microservice

## Introduction

This project is created to understand <a href="https://docs.temporal.io/">Temporal</a>'s Microservice approachment. The repo contains:

- Example Temporal Microservice (Launch Microservice)
- Deployment files of Microservice
- API Gateway for triggering or sending a signal to Workflows. (Endpoint would be gRPC or Rest)

![Temporal flow implementation](https://github.com/esgungor/temporal-microservice/blob/main/docs/TemporalFlow.png?raw=true)

<table>
<thead>
<th>Package Name</th>
<th>Description</th>

</thead>
<tbody>
<tr>
<td>App</td>
<td>It contains defitinition of Workflow</td>
<tr>

<td>Server</td>
<td>API gateway to interact with temporal workflows which is contains GET, CREATE, DELETE, UPDATE Methods</td>
</tr>
<tr>

<td>Worker</td>
<td>Microservice which is reponsible for run Workflows Actions</td>
</tr>
</tbody>

</table>

## Endpoints

<table>
<thead>
<th>Endpoints</th>
<th>Method</th>

<th>Description</th>

</thead>
<tbody>
<tr>
<td>/</td>
<td>GET</td>
<td>List of launch workflows</td>
</tr>

<tr>
<td>/</td>
<td>POST</td>
<td>Create a predefined launch</td>
</tr>

<tr>
<td>/launch/{workflowId}/{runId}</td>
<td>DELETE</td>
<td>Delete a launch</td>
</tr>

<tr>
<td>/launch/{workflowId}/{runId}</td>
<td>UPDATE</td>
<td>Update a launch without paramater(Only for testing endpoint</td>
</tr>

<tr>
<td>/launch/{workflowId}/{runId}</td>
<td>GET</td>
<td>Get a launch properties. (Only state you could get from this endpoint)</td>
</tr>
</tbody>

</table>

## Run

Firstly you should set Temporal server ip as a environment variable.

```bash
export TEMPORAL_SERVER_IP=<TEMPORAL IP>
```

After that you should start <b>API Gateway Server</b> and <b>Temporal worker</b> with the following command

```bash
go run server/main.go
go run worker/main.go
```

API Gateway will serve from port 5000.
