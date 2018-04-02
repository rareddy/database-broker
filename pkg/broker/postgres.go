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
func (ds PostgreSQL) springboot(bindAlias string, multiSource bool) map[string]interface{} {
	props := make(map[string]interface{})
	props[sb(bindAlias, "driver-class-name", multiSource)] = "org.postgresql.Driver"
	props[sb(bindAlias, "url", multiSource)] = "jdbc:postgresql://" + i2s(ds.Parameters["service-name"]) + ":5432/" + i2s(ds.Parameters["database-name"])
	props[sb(bindAlias, "username", multiSource)] = i2s(ds.Parameters["username"])
	props[sb(bindAlias, "password", multiSource)] = i2s(ds.Parameters["password"])
	props[sb(bindAlias, "validationQuery", multiSource)] = i2s(ds.Parameters["ping-query"])
	return props
}

func (ds PostgreSQL) wildflyswarm(bindAlias string, multiSource bool) map[string]interface{} {
	props := make(map[string]interface{})
	props[wfs(bindAlias, "driver-name")] = "postgresql"
	props[wfs(bindAlias, "jndi-name")] = "java:/jboss/datasources/" + bindAlias
	props[wfs(bindAlias, "connection-url")] = "jdbc:postgresql://" + i2s(ds.Parameters["service-name"]) + ":5432/" + i2s(ds.Parameters["database-name"])
	props[wfs(bindAlias, "user-name")] = i2s(ds.Parameters["username"])
	props[wfs(bindAlias, "password")] = i2s(ds.Parameters["password"])
	props[wfs(bindAlias, "check-valid-connection-sql")] = i2s(ds.Parameters["ping-query"])
	return props
}

func (ds PostgreSQL) nodejs(bindAlias string, multiSource bool) map[string]interface{} {
	// TODO: need to figure out what is it required in Node.JS
	return ds.Parameters
}

func (ds PostgreSQL) other(bindAlias string, multiSource bool) map[string]interface{} {
	// return with no modification
	return ds.Parameters
}
