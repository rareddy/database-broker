package broker

// details about MySQL datasource. it is an external data source and configurable

type SqlServer struct {
	datasource //ananymous field
}

// support externaldatasource interface
func (ds SqlServer) defaultPort() int {
	return 1433
}

func (ds SqlServer) name() string {
	return "sqlserver"
}

// support configurable interface
func (ds SqlServer) springboot(bindAlias string, multiSource bool) map[string]interface{} {
	props := make(map[string]interface{})
	props[sb(bindAlias, "driver-class-name", multiSource)] = "net.sourceforge.jtds.jdbc.Driver"
	props[sb(bindAlias, "url", multiSource)] = "jdbc:jtds:sqlserver://" + i2s(ds.Parameters["service-name"]) + ":1433/" + i2s(ds.Parameters["database-name"])
	props[sb(bindAlias, "username", multiSource)] = i2s(ds.Parameters["username"])
	props[sb(bindAlias, "password", multiSource)] = i2s(ds.Parameters["password"])
	props[sb(bindAlias, "validationQuery", multiSource)] = i2s(ds.Parameters["ping-query"])
	return props
}

func (ds SqlServer) wildflyswarm(bindAlias string, multiSource bool) map[string]interface{} {
	props := make(map[string]interface{})
	props[wfs(bindAlias, "driver-name")] = "sqlserver"
	props[wfs(bindAlias, "jndi-name")] = "java:/jboss/datasources/" + bindAlias
	props[wfs(bindAlias, "connection-url")] = "jdbc:jtds:sqlserver://" + i2s(ds.Parameters["service-name"]) + ":1433/" + i2s(ds.Parameters["database-name"])
	props[wfs(bindAlias, "user-name")] = i2s(ds.Parameters["username"])
	props[wfs(bindAlias, "password")] = i2s(ds.Parameters["password"])
	props[wfs(bindAlias, "check-valid-connection-sql")] = i2s(ds.Parameters["ping-query"])
	return props
}

func (ds SqlServer) nodejs(bindAlias string, multiSource bool) map[string]interface{} {
	// TODO: need to figure out what is it required in Node.JS
	return ds.Parameters
}

func (ds SqlServer) other(bindAlias string, multiSource bool) map[string]interface{} {
	// return with no modification
	return ds.Parameters
}
