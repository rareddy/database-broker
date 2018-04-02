package broker

// details about MySQL datasource. it is an external data source and configurable

type MySql struct {
	datasource //ananymous field
}

// support externaldatasource interface
func (ds MySql) defaultPort() int {
	return 3306
}

func (ds MySql) name() string {
	return "mysql"
}

// support configurable interface
func (ds MySql) springboot(multiSource bool) map[string]interface{} {
	props := make(map[string]interface{})
	props[sb(ds.SourceName, "driver-class-name", multiSource)] = "com.mysql.jdbc.Driver"
	props[sb(ds.SourceName, "url", multiSource)] = "jdbc:mysql://"+ds.SourceName+":3306/"+i2s(ds.Parameters["database-name"])
	props[sb(ds.SourceName, "username", multiSource)] = i2s(ds.Parameters["username"])
	props[sb(ds.SourceName, "password", multiSource)] = i2s(ds.Parameters["password"])
	props[sb(ds.SourceName, "validationQuery", multiSource)] = i2s(ds.Parameters["ping-query"])
	return props
}

func (ds MySql) wildflyswarm(multiSource bool) map[string]interface{} {
	props := make(map[string]interface{})
	props[wfs(ds.SourceName, "driver-name")] = "mysql"
	props[wfs(ds.SourceName, "jndi-name")] = "java:datasources/"+ds.SourceName
	props[wfs(ds.SourceName, "connection-url")] = "jdbc:mysql://"+ds.SourceName+":3306/"+i2s(ds.Parameters["database-name"])
	props[wfs(ds.SourceName, "user-name")] = i2s(ds.Parameters["username"])
	props[wfs(ds.SourceName, "password")] = i2s(ds.Parameters["password"])
	props[wfs(ds.SourceName, "check-valid-connection-sql")] = i2s(ds.Parameters["ping-query"])
	return props
}

func (ds MySql) nodejs(multiSource bool) map[string]interface{} {
	// TODO: need to figure out what is it required in Node.JS
	return ds.Parameters
}

func (ds MySql) other(multiSource bool) map[string]interface{} {
	// return with no modification
	return ds.Parameters
}

