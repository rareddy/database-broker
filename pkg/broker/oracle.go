package broker

// details about MySQL datasource. it is an external data source and configurable

type Oracle struct {
	datasource //ananymous field
}

// support externaldatasource interface
func (ds Oracle) defaultPort() int {
	return 1521
}

func (ds Oracle) name() string {
	return "oracle"
}

// support configurable interface
func (ds Oracle) springboot(bindAlias string, multiSource bool) map[string]interface{} {
	props := make(map[string]interface{})
	props[sb(bindAlias, "driver-class-name", multiSource)] = "oracle.jdbc.driver.OracleDriver"
	props[sb(bindAlias, "url", multiSource)] = "jdbc:oracle:thin:@//" + i2s(ds.Parameters["service-name"]) + ":1521/" + i2s(ds.Parameters["database-name"])
	props[sb(bindAlias, "username", multiSource)] = i2s(ds.Parameters["username"])
	props[sb(bindAlias, "password", multiSource)] = i2s(ds.Parameters["password"])
	props[sb(bindAlias, "validationQuery", multiSource)] = i2s(ds.Parameters["ping-query"])
	return props
}

func (ds Oracle) wildflyswarm(bindAlias string, multiSource bool) map[string]interface{} {
	props := make(map[string]interface{})
	props[wfs(bindAlias, "driver-name")] = "oracle"
	props[wfs(bindAlias, "jndi-name")] = "java:datasources/" + bindAlias
	props[wfs(bindAlias, "connection-url")] = "jdbc:oracle:thin:@//" + i2s(ds.Parameters["service-name"]) + ":1521/" + i2s(ds.Parameters["database-name"])
	props[wfs(bindAlias, "user-name")] = i2s(ds.Parameters["username"])
	props[wfs(bindAlias, "password")] = i2s(ds.Parameters["password"])
	props[wfs(bindAlias, "check-valid-connection-sql")] = i2s(ds.Parameters["ping-query"])
	return props
}

func (ds Oracle) nodejs(bindAlias string, multiSource bool) map[string]interface{} {
	// TODO: need to figure out what is it required in Node.JS
	return ds.Parameters
}

func (ds Oracle) other(bindAlias string, multiSource bool) map[string]interface{} {
	// return with no modification
	return ds.Parameters
}
