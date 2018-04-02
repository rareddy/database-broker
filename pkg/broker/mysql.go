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
func (ds MySql) springboot(bindAlias string, multiSource bool) map[string]interface{} {
	props := make(map[string]interface{})
	props[sb(bindAlias, "driver-class-name", multiSource)] = "com.mysql.jdbc.Driver"
	props[sb(bindAlias, "url", multiSource)] = "jdbc:mysql://" + i2s(ds.Parameters["service-name"]) + ":3306/" + i2s(ds.Parameters["database-name"])
	props[sb(bindAlias, "username", multiSource)] = i2s(ds.Parameters["username"])
	props[sb(bindAlias, "password", multiSource)] = i2s(ds.Parameters["password"])
	props[sb(bindAlias, "validationQuery", multiSource)] = i2s(ds.Parameters["ping-query"])
	return props
}

func (ds MySql) wildflyswarm(bindAlias string, multiSource bool) map[string]interface{} {
	props := make(map[string]interface{})
	props[wfs(bindAlias, "driver-name")] = "mysql"
	props[wfs(bindAlias, "jndi-name")] = "java:datasources/" + bindAlias
	props[wfs(bindAlias, "connection-url")] = "jdbc:mysql://" + i2s(ds.Parameters["service-name"]) + ":3306/" + i2s(ds.Parameters["database-name"])
	props[wfs(bindAlias, "user-name")] = i2s(ds.Parameters["username"])
	props[wfs(bindAlias, "password")] = i2s(ds.Parameters["password"])
	props[wfs(bindAlias, "check-valid-connection-sql")] = i2s(ds.Parameters["ping-query"])
	return props
}

func (ds MySql) nodejs(bindAlias string, multiSource bool) map[string]interface{} {
	// TODO: need to figure out what is it required in Node.JS
	return ds.Parameters
}

func (ds MySql) other(bindAlias string, multiSource bool) map[string]interface{} {
	// return with no modification
	return ds.Parameters
}
