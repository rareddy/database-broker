package broker

// details about MySQL datasource. it is an external data source and configurable

type SqlServer struct {
	datasource //ananymous field
}

// support externaldatasource interface
func (ds SqlServer) defaultPort() int {
	return 1433
}

// support configurable interface
func (ds SqlServer) springboot() map[string]interface{} {
	props := make(map[string]interface{})
	props[sb(ds.SourceName, "driver-class-name")] = "net.sourceforge.jtds.jdbc.Driver"
	props[sb(ds.SourceName, "jndi-name")] = "java:datasources/"+ds.SourceName
	props[sb(ds.SourceName, "url")] = "jdbc:jtds:sqlserver://"+ds.SourceName+":1433/"+i2s(ds.Parameters["database-name"])
	props[sb(ds.SourceName, "username")] = i2s(ds.Parameters["username"])
	props[sb(ds.SourceName, "password")] = i2s(ds.Parameters["password"])
	props[sb(ds.SourceName, "validationQuery")] = i2s(ds.Parameters["ping-query"])
	return props
}

func (ds SqlServer) wildflyswarm() map[string]interface{} {
	props := make(map[string]interface{})
	props[wfs(ds.SourceName, "driver-name")] = "sqlserver"
	props[wfs(ds.SourceName, "jndi-name")] = "java:datasources/"+ds.SourceName
	props[wfs(ds.SourceName, "connection-url")] = "jdbc:jtds:sqlserver://"+ds.SourceName+":1433/"+i2s(ds.Parameters["database-name"])
	props[wfs(ds.SourceName, "user-name")] = i2s(ds.Parameters["username"])
	props[wfs(ds.SourceName, "password")] = i2s(ds.Parameters["password"])
	props[wfs(ds.SourceName, "check-valid-connection-sql")] = i2s(ds.Parameters["ping-query"])
	return props
}

func (ds SqlServer) nodejs() map[string]interface{} {
	// TODO: need to figure out what is it required in Node.JS
	return ds.Parameters
}

func (ds SqlServer) other() map[string]interface{} {
	// return with no modification
	return ds.Parameters
}

