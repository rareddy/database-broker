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
func (ds Oracle) springboot() map[string]interface{} {
	props := make(map[string]interface{})
	props[sb(ds.SourceName, "driver-class-name")] = "oracle.jdbc.driver.OracleDriver"
	props[sb(ds.SourceName, "url")] = "jdbc:oracle:thin:@//"+ds.SourceName+":1521/"+i2s(ds.Parameters["database-name"])
	props[sb(ds.SourceName, "username")] = i2s(ds.Parameters["username"])
	props[sb(ds.SourceName, "password")] = i2s(ds.Parameters["password"])
	props[sb(ds.SourceName, "validationQuery")] = i2s(ds.Parameters["ping-query"])
	return props
}

func (ds Oracle) wildflyswarm() map[string]interface{} {
	props := make(map[string]interface{})
	props[wfs(ds.SourceName, "driver-name")] = "oracle"
	props[wfs(ds.SourceName, "jndi-name")] = "java:datasources/"+ds.SourceName
	props[wfs(ds.SourceName, "connection-url")] = "jdbc:oracle:thin:@//"+ds.SourceName+":1521/"+i2s(ds.Parameters["database-name"])
	props[wfs(ds.SourceName, "user-name")] = i2s(ds.Parameters["username"])
	props[wfs(ds.SourceName, "password")] = i2s(ds.Parameters["password"])
	props[wfs(ds.SourceName, "check-valid-connection-sql")] = i2s(ds.Parameters["ping-query"])
	return props
}

func (ds Oracle) nodejs() map[string]interface{} {
	// TODO: need to figure out what is it required in Node.JS
	return ds.Parameters
}

func (ds Oracle) other() map[string]interface{} {
	// return with no modification
	return ds.Parameters
}

