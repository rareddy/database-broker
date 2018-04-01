package broker

// details about MySQL datasource. it is an external data source and configurable

type PostgreSQL struct {
	datasource //ananymous field
}

// support externaldatasource interface
func (ds PostgreSQL) defaultPort() int {
	return 5432
}

func (ds PostgreSQL) name() string {
	return "postgresql"
}

// support configurable interface
func (ds PostgreSQL) springboot() map[string]interface{} {
	props := make(map[string]interface{})
	props[sb(ds.SourceName, "driver-class-name")] = "org.postgresql.Driver"
	props[sb(ds.SourceName, "url")] = "jdbc:postgresql://"+ds.SourceName+":5432/"+i2s(ds.Parameters["database-name"])
	props[sb(ds.SourceName, "username")] = i2s(ds.Parameters["username"])
	props[sb(ds.SourceName, "password")] = i2s(ds.Parameters["password"])
	props[sb(ds.SourceName, "validationQuery")] = i2s(ds.Parameters["ping-query"])
	return props
}

func (ds PostgreSQL) wildflyswarm() map[string]interface{} {
	props := make(map[string]interface{})
	props[wfs(ds.SourceName, "driver-name")] = "postgresql"
	props[wfs(ds.SourceName, "jndi-name")] = "java:datasources/"+ds.SourceName
	props[wfs(ds.SourceName, "connection-url")] = "jdbc:postgresql://"+ds.SourceName+":5432/"+i2s(ds.Parameters["database-name"])
	props[wfs(ds.SourceName, "user-name")] = i2s(ds.Parameters["username"])
	props[wfs(ds.SourceName, "password")] = i2s(ds.Parameters["password"])
	props[wfs(ds.SourceName, "check-valid-connection-sql")] = i2s(ds.Parameters["ping-query"])
	return props
}

func (ds PostgreSQL) nodejs() map[string]interface{} {
	// TODO: need to figure out what is it required in Node.JS
	return ds.Parameters
}

func (ds PostgreSQL) other() map[string]interface{} {
	// return with no modification
	return ds.Parameters
}

